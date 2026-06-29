package lab

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func derivTestEraStart() int64 {
	return time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC).UnixMilli()
}

// derivTestDir returns a repo-local temp dir (not under /tmp, which the audit
// rejects as an ephemeral cache path).
func derivTestDir(t *testing.T) string {
	t.Helper()
	d, err := os.MkdirTemp(".", "derivaudit")
	if err != nil {
		t.Fatalf("mkdir temp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(d) })
	abs, err := filepath.Abs(d)
	if err != nil {
		t.Fatalf("abs: %v", err)
	}
	return abs
}

func writeDerivAnchorCSV(t *testing.T, path string, startMs int64, n int) {
	t.Helper()
	var b strings.Builder
	b.WriteString("open_time,open,high,low,close,volume,close_time\n")
	for i := 0; i < n; i++ {
		ot := startMs + int64(i)*derivativesIntervalMs
		ct := ot + derivativesIntervalMs - 1
		fmt.Fprintf(&b, "%d,100,101,99,100,0,%d\n", ot, ct)
	}
	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
		t.Fatalf("write anchor csv: %v", err)
	}
}

type derivCSVOptions struct {
	rows         int  // number of contiguous rows (default n)
	skipEvery    int  // if >0, skip rows where index%skipEvery==0
	badCloseTime bool // emit close_time != open+5m-1ms
	negativeLast bool // emit a negative close on the last row
}

func writeDerivSourceCSV(t *testing.T, path string, startMs, n int64, opts derivCSVOptions) {
	t.Helper()
	rows := opts.rows
	if rows == 0 {
		rows = int(n)
	}
	var b strings.Builder
	b.WriteString("open_time,open,high,low,close,close_time,source_object_id\n")
	for i := 0; i < rows; i++ {
		if opts.skipEvery > 0 && i%opts.skipEvery == 0 && i != 0 {
			continue
		}
		ot := startMs + int64(i)*derivativesIntervalMs
		ct := ot + derivativesIntervalMs - 1
		if opts.badCloseTime {
			ct = ot // wrong close_time
		}
		closePx := "100"
		if opts.negativeLast && i == rows-1 {
			closePx = "-0.5"
		}
		fmt.Fprintf(&b, "%d,100,101,99,%s,%d,obj-%d\n", ot, closePx, ct, i)
	}
	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
		t.Fatalf("write deriv csv: %v", err)
	}
}

func derivTestConfig(t *testing.T, dir string, n int) FuturesDerivativesContextSourceAuditConfig {
	t.Helper()
	start := derivTestEraStart()
	end := start + int64(n-1)*derivativesIntervalMs
	anchorPath := filepath.Join(dir, "btcusdt_futures_um_5m_test.csv")
	writeDerivAnchorCSV(t, anchorPath, start, n)
	markPath := filepath.Join(dir, "binance_usdm_mark_price_klines_5m_BTCUSDT_test.csv")
	indexPath := filepath.Join(dir, "binance_usdm_index_price_klines_5m_BTCUSDT_test.csv")
	premiumPath := filepath.Join(dir, "binance_usdm_premium_index_klines_5m_BTCUSDT_test.csv")
	writeDerivSourceCSV(t, markPath, start, int64(n), derivCSVOptions{})
	writeDerivSourceCSV(t, indexPath, start, int64(n), derivCSVOptions{})
	writeDerivSourceCSV(t, premiumPath, start, int64(n), derivCSVOptions{negativeLast: true})
	return FuturesDerivativesContextSourceAuditConfig{
		DerivativeSources: []FuturesDerivativesContextSourceFileConfig{
			{Symbol: "BTCUSDT", SourceFamily: "mark_price_klines", ArchiveFamily: "markPriceKlines", Path: markPath, Required: true},
			{Symbol: "BTCUSDT", SourceFamily: "index_price_klines", ArchiveFamily: "indexPriceKlines", Path: indexPath, Required: true},
			{Symbol: "BTCUSDT", SourceFamily: "premium_index_klines", ArchiveFamily: "premiumIndexKlines", Path: premiumPath, Required: false, AllowNonPositive: true},
		},
		Anchors: []FuturesRangeUniverseSourceConfig{
			{Symbol: "BTCUSDT", Path: anchorPath, ApprovedPath: anchorPath, SkipSplitEligibilityCheck: true},
		},
		MinAlignedCoverage:       0.5,
		ConservativeLagIntervals: 1,
		EraStartMs:               start,
		EraEndMs:                 end,
	}
}

func TestRunFuturesDerivativesContextSourceAuditPasses(t *testing.T) {
	dir := derivTestDir(t)
	cfg := derivTestConfig(t, dir, 24)
	res, err := RunFuturesDerivativesContextSourceAudit(cfg, DefaultSplits())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.StopState != DerivativesContextSourceAuditStopStatePassedNeedsContextBrief {
		t.Fatalf("stop_state=%q want pass", res.StopState)
	}
	if res.RequiredStreams != 2 || res.AlignedRequiredStreams != 2 {
		t.Fatalf("required=%d aligned=%d want 2/2", res.RequiredStreams, res.AlignedRequiredStreams)
	}
	if len(res.SourceRows) != 3 || len(res.CandleAnchorRows) != 1 || len(res.CoverageRows) != 3 {
		t.Fatalf("rows: sources=%d anchors=%d coverage=%d", len(res.SourceRows), len(res.CandleAnchorRows), len(res.CoverageRows))
	}
	if len(res.TimestampAlignmentRows) != 3 || len(res.PublicationLagRows) != 3 || len(res.ProvenanceRows) != 3 {
		t.Fatalf("alignment=%d lag=%d provenance=%d", len(res.TimestampAlignmentRows), len(res.PublicationLagRows), len(res.ProvenanceRows))
	}
	if len(res.SummaryRows) != 4 { // 3 streams + overall
		t.Fatalf("summary rows=%d want 4", len(res.SummaryRows))
	}
	for _, sr := range res.SummaryRows {
		if sr.Trades != 0 || !sr.ZeroTradeCompatible {
			t.Fatalf("summary not zero-trade: %+v", sr)
		}
	}
	// Anti-lookahead: alignment never uses future rows.
	for _, ar := range res.TimestampAlignmentRows {
		if ar.UsesFutureRows {
			t.Fatalf("alignment row uses future rows: %+v", ar)
		}
	}
	// Warmup boundary surfaced as a skip.
	foundWarmup := false
	for _, sk := range res.SkipRows {
		if sk.Reason == "lag_warmup_boundary" {
			foundWarmup = true
		}
	}
	if !foundWarmup {
		t.Fatalf("expected lag_warmup_boundary skip row")
	}
	// Provenance recorded a sha256 for every stream.
	for _, sr := range res.SourceRows {
		if sr.FileSHA256 == "" {
			t.Fatalf("missing sha256 for %s/%s", sr.SourceFamily, sr.Symbol)
		}
	}
}

func TestRunFuturesDerivativesContextSourceAuditRejectsSourceGap(t *testing.T) {
	dir := derivTestDir(t)
	cfg := derivTestConfig(t, dir, 24)
	cfg.DerivativeSources[0].Path = filepath.Join(dir, "missing_mark.csv") // required mark missing
	res, _ := RunFuturesDerivativesContextSourceAudit(cfg, DefaultSplits())
	if res.StopState != DerivativesContextSourceAuditStopStateRejectedSourceGap {
		t.Fatalf("stop_state=%q want source gap", res.StopState)
	}
}

func TestRunFuturesDerivativesContextSourceAuditRejectsTimestampGap(t *testing.T) {
	dir := derivTestDir(t)
	cfg := derivTestConfig(t, dir, 24)
	start := derivTestEraStart()
	// Rewrite required mark file with bad close_time.
	writeDerivSourceCSV(t, cfg.DerivativeSources[0].Path, start, 24, derivCSVOptions{badCloseTime: true})
	res, _ := RunFuturesDerivativesContextSourceAudit(cfg, DefaultSplits())
	if res.StopState != DerivativesContextSourceAuditStopStateRejectedTimestampOrFinality {
		t.Fatalf("stop_state=%q want timestamp/finality gap", res.StopState)
	}
}

func TestRunFuturesDerivativesContextSourceAuditRejectsAlignmentGap(t *testing.T) {
	dir := derivTestDir(t)
	cfg := derivTestConfig(t, dir, 24)
	cfg.MinAlignedCoverage = 0.99
	start := derivTestEraStart()
	// Make required index sparse so lag coverage drops below 0.99.
	writeDerivSourceCSV(t, cfg.DerivativeSources[1].Path, start, 24, derivCSVOptions{skipEvery: 2})
	res, _ := RunFuturesDerivativesContextSourceAudit(cfg, DefaultSplits())
	if res.StopState != DerivativesContextSourceAuditStopStateRejectedAlignmentGap {
		t.Fatalf("stop_state=%q want alignment gap", res.StopState)
	}
}

func TestRunFuturesDerivativesContextSourceAuditRejectsClosedFamilyRescue(t *testing.T) {
	dir := derivTestDir(t)
	cfg := derivTestConfig(t, dir, 24)
	cfg.RescueClosedFamily = true
	res, err := RunFuturesDerivativesContextSourceAudit(cfg, DefaultSplits())
	if err == nil {
		t.Fatalf("expected error for closed-family rescue")
	}
	if res.StopState != DerivativesContextSourceAuditStopStateRejectedClosedFamilyRescue {
		t.Fatalf("stop_state=%q want closed family rescue", res.StopState)
	}
}

func TestRunFuturesDerivativesContextSourceAuditRejectsLivePath(t *testing.T) {
	dir := derivTestDir(t)
	cfg := derivTestConfig(t, dir, 24)
	cfg.DerivativeSources[0].Path = "https://data.binance.vision/some/object.csv"
	res, err := RunFuturesDerivativesContextSourceAudit(cfg, DefaultSplits())
	if err == nil {
		t.Fatalf("expected error for live/private path")
	}
	if res.StopState != DerivativesContextSourceAuditStopStateRejectedLiveOrPrivateAPI {
		t.Fatalf("stop_state=%q want live/private path", res.StopState)
	}
}
