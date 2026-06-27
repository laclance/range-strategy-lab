package lab

import (
	"math"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRangeContextTriageSourceAndCoverageAcceptance(t *testing.T) {
	parent := make([]Candle, 96)
	for i := range parent {
		parent[i] = testCandle(i, 100, 101, 99, 100)
		if i < rangeContextTriageExpectedZeroVol {
			parent[i].Volume = 0
		}
	}
	manifest, err := ValidateResearchSource("btcusdt_futures_um_5m_triage.csv", []string{"open_time", "open", "high", "low", "close", "volume", "close_time"}, parent, SourceValidationOptions{Product: SourceProductBinanceUSDMFutures})
	if err != nil {
		t.Fatal(err)
	}
	cfg := testRangeContextTriageConfig()
	cfg.ApprovedSourcePath = manifest.Path
	cfg.SkipSourceFactCheck = false
	cfg.ExpectedSourceRows = len(parent)
	cfg.ExpectedFirstOpenTime = parent[0].OpenTime.UTC().Format(timeLayout)
	cfg.ExpectedLastOpenTime = parent[len(parent)-1].OpenTime.UTC().Format(timeLayout)
	cfg.SkipCoverageCountCheck = false
	cfg.Expected15MRows = 32
	cfg.Expected15MLastOpenTime = parent[93].OpenTime.UTC().Format(timeLayout)
	cfg.Expected1HRows = 8
	cfg.Expected1HLastOpenTime = parent[84].OpenTime.UTC().Format(timeLayout)
	cfg.Expected4HRows = 2
	cfg.Expected4HLastOpenTime = parent[48].OpenTime.UTC().Format(timeLayout)

	result, err := RunFuturesRangeContextTriageAudit(parent, manifest, cfg, []Split{{Name: fullSplitName}})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.SourceRows) != 1 || !result.SourceRows[0].SourceFactsPass {
		t.Fatalf("bad source rows: %+v", result.SourceRows)
	}
	if len(result.CoverageRows) != 3 {
		t.Fatalf("coverage rows=%d, want 3", len(result.CoverageRows))
	}
	for _, row := range result.CoverageRows {
		if !row.CoverageFactsPass || row.ValidationStatus != "accepted" {
			t.Fatalf("bad coverage row: %+v", row)
		}
	}

	manifest.ComparisonOnly = true
	manifest.Product = "Binance spot"
	result, err = RunFuturesRangeContextTriageAudit(parent, manifest, cfg, []Split{{Name: fullSplitName}})
	if err != nil {
		t.Fatal(err)
	}
	if result.StopState != RangeContextTriageStopStateSourceGap {
		t.Fatalf("stop=%s, want source gap", result.StopState)
	}
}

func TestRangeContextTriageSourcePathRejection(t *testing.T) {
	dir := t.TempDir()
	parent := make([]Candle, 24)
	for i := range parent {
		parent[i] = testCandle(i, 100, 101, 99, 100)
	}
	path := writeRangeUniverseTestCSV(t, dir, "btcusdt_futures_um_5m_triage.csv", parent)
	candles, manifest, err := LoadResearchSourceCSV(path, SourceValidationOptions{Product: SourceProductBinanceUSDMFutures})
	if err != nil {
		t.Fatal(err)
	}
	cfg := testRangeContextTriageConfig()
	cfg.SkipSourceFactCheck = false
	cfg.ApprovedSourcePath = filepath.Join(dir, "other_btcusdt_futures_um_5m.csv")
	result, err := RunFuturesRangeContextTriageAudit(candles, manifest, cfg, []Split{{Name: fullSplitName}})
	if err != nil {
		t.Fatal(err)
	}
	if result.StopState != RangeContextTriageStopStateSourceGap || !strings.Contains(result.SourceRows[0].ValidationError, "approved path") {
		t.Fatalf("expected approved-path source gap, got stop=%s source=%+v", result.StopState, result.SourceRows)
	}
}

func TestRangeContextTriageMatureFreezeAndInvalidSkips(t *testing.T) {
	candles := []Candle{
		testCandle(0, 105, 110, 100, 105),
		testCandle(1, 106, 112, 99, 106),
		testCandle(2, 107, 111, 101, 107),
		testCandle(3, 108, 130, 90, 108),
		testCandle(4, 109, 111, 107, 109),
	}
	classifications := []RangeClassification{
		{Index: 0, RawActive: true},
		{Index: 1, RawActive: true},
		{Index: 2, RawActive: true, Active: true},
		{Index: 3, RawActive: true, Active: true},
		{Index: 4},
	}
	atr := []float64{0.01, 0.01, 0.01, 0.01, 0.01}
	rows := rangeContextTriageEpisodeRows(candles, classifications, atr, RangeDiscoveryTimeframe1h, testRangeContextTriageConfig(), []Split{{Name: fullSplitName}}, 0)
	if len(rows) != 1 {
		t.Fatalf("episodes=%d, want 1", len(rows))
	}
	row := rows[0]
	if row.High != 112 || row.Low != 99 || row.MatureIndex != 2 || row.RawEndIndex != 3 {
		t.Fatalf("freeze should stop at first mature active close: %+v", row)
	}
	if !row.Eligible {
		t.Fatalf("row should be eligible: %+v", row)
	}

	noMature := rangeContextTriageEpisodeRows(candles[:2], []RangeClassification{{RawActive: true}, {RawActive: true}}, atr[:2], RangeDiscoveryTimeframe1h, testRangeContextTriageConfig(), []Split{{Name: fullSplitName}}, 0)
	if len(noMature) != 1 || noMature[0].SkippedReason != "no_mature_active_bar" {
		t.Fatalf("expected no mature active skip, got %+v", noMature)
	}
	flat := []Candle{testCandle(0, 100, 100, 100, 100), testCandle(1, 100, 100, 100, 100)}
	flatRows := rangeContextTriageEpisodeRows(flat, []RangeClassification{{RawActive: true, Active: true}, {}}, []float64{0.01, 0.01}, RangeDiscoveryTimeframe1h, testRangeContextTriageConfig(), []Split{{Name: fullSplitName}}, 0)
	if len(flatRows) != 1 || flatRows[0].SkippedReason != "non_positive_width" {
		t.Fatalf("expected non-positive width skip, got %+v", flatRows)
	}
	missingATR := rangeContextTriageEpisodeRows(candles[:2], []RangeClassification{{RawActive: true, Active: true}, {}}, []float64{math.NaN(), math.NaN()}, RangeDiscoveryTimeframe1h, testRangeContextTriageConfig(), []Split{{Name: fullSplitName}}, 0)
	if len(missingATR) != 1 || missingATR[0].SkippedReason != "missing_atr" {
		t.Fatalf("expected missing ATR skip, got %+v", missingATR)
	}
}

func TestRangeContextTriageSessionAndQualityBuckets(t *testing.T) {
	cases := []struct {
		at   time.Time
		want string
	}{
		{time.Date(2024, 1, 1, 7, 59, 0, 0, time.UTC), RangeContextTriageSessionAsia},
		{time.Date(2024, 1, 1, 8, 0, 0, 0, time.UTC), RangeContextTriageSessionEurope},
		{time.Date(2024, 1, 1, 13, 0, 0, 0, time.UTC), RangeContextTriageSessionUSOverlap},
		{time.Date(2024, 1, 1, 17, 0, 0, 0, time.UTC), RangeContextTriageSessionUSLate},
	}
	for _, tc := range cases {
		if got := RangeContextTriageUTCSession(tc.at); got != tc.want {
			t.Fatalf("session=%s, want %s", got, tc.want)
		}
	}
	cfg := testRangeContextTriageConfig()
	row := FuturesRangeContextTriageEpisodeRow{WidthPct: 0.001, WidthToATRRatio: 0.5, PreMatureCloseInsideRate: 1}
	if got := rangeContextTriageQualityBucket(row, cfg); got != RangeContextTriageQualityTooNarrowNoise {
		t.Fatalf("quality=%s, want too narrow", got)
	}
	row = FuturesRangeContextTriageEpisodeRow{WidthPct: 0.02, WidthToATRRatio: 5, PreMatureCloseInsideRate: 1}
	if got := rangeContextTriageQualityBucket(row, cfg); got != RangeContextTriageQualityWideVolatile {
		t.Fatalf("quality=%s, want wide", got)
	}
	row = FuturesRangeContextTriageEpisodeRow{WidthPct: 0.005, WidthToATRRatio: 2, PreMatureMidCrossCount: 3, PreMatureBoundaryTouchCount: 4, PreMatureCloseInsideRate: 1}
	if got := rangeContextTriageQualityBucket(row, cfg); got != RangeContextTriageQualityChoppy {
		t.Fatalf("quality=%s, want choppy", got)
	}
	row = FuturesRangeContextTriageEpisodeRow{WidthPct: 0.002, WidthToATRRatio: 1.2, PreMatureCloseInsideRate: 0.95}
	if got := rangeContextTriageQualityBucket(row, cfg); got != RangeContextTriageQualityNarrowOrderly {
		t.Fatalf("quality=%s, want narrow orderly", got)
	}
	row = FuturesRangeContextTriageEpisodeRow{WidthPct: 0.005, WidthToATRRatio: 2, PreMatureCloseInsideRate: 0.95}
	if got := rangeContextTriageQualityBucket(row, cfg); got != RangeContextTriageQualityBalancedOrderly {
		t.Fatalf("quality=%s, want balanced", got)
	}
}

func TestRangeContextTriagePrimaryLabelsAndPrecedence(t *testing.T) {
	cfg := testRangeContextTriageConfig()
	base := rangeContextTriageTestEpisode()
	tests := []struct {
		name    string
		candles []Candle
		horizon int
		quality string
		want    string
	}{
		{
			name: "contained",
			candles: []Candle{
				testCandle(0, 105, 106, 104, 105),
				testCandle(1, 104, 106, 103, 104),
				testCandle(2, 106, 107, 105, 106),
				testCandle(3, 104, 106, 103, 104),
			},
			horizon: 3,
			want:    RangeContextTriageLabelContainedRotation,
		},
		{
			name: "clean up",
			candles: []Candle{
				testCandle(0, 105, 106, 104, 105),
				testCandle(1, 111, 119, 110, 118),
				testCandle(2, 118, 120, 117, 119),
			},
			horizon: 2,
			want:    RangeContextTriageLabelCleanExpansionUp,
		},
		{
			name: "clean down",
			candles: []Candle{
				testCandle(0, 105, 106, 104, 105),
				testCandle(1, 99, 100, 91, 92),
				testCandle(2, 92, 93, 90, 91),
			},
			horizon: 2,
			want:    RangeContextTriageLabelCleanExpansionDown,
		},
		{
			name: "false break",
			candles: []Candle{
				testCandle(0, 105, 106, 104, 105),
				testCandle(1, 111, 112, 110, 111),
				testCandle(2, 106, 107, 105, 106),
			},
			horizon: 2,
			want:    RangeContextTriageLabelFalseBreakReentryUp,
		},
		{
			name: "boundary chop",
			candles: []Candle{
				testCandle(0, 105, 106, 104, 105),
				testCandle(1, 111, 112, 110, 111),
				testCandle(2, 106, 107, 105, 106),
				testCandle(3, 112, 113, 111, 112),
				testCandle(4, 104, 105, 103, 104),
			},
			horizon: 4,
			want:    RangeContextTriageLabelBoundaryChop,
		},
		{
			name: "drift",
			candles: []Candle{
				testCandle(0, 105, 106, 104, 105),
				testCandle(1, 111, 112, 110, 111),
				testCandle(2, 112, 113, 111, 112),
			},
			horizon: 2,
			want:    RangeContextTriageLabelDriftThroughUp,
		},
		{
			name: "low width precedence",
			candles: []Candle{
				testCandle(0, 105, 106, 104, 105),
				testCandle(1, 104, 106, 103, 104),
				testCandle(2, 106, 107, 105, 106),
				testCandle(3, 104, 106, 103, 104),
			},
			horizon: 3,
			quality: RangeContextTriageQualityTooNarrowNoise,
			want:    RangeContextTriageLabelLowWidthNoise,
		},
		{
			name: "missing",
			candles: []Candle{
				testCandle(0, 105, 106, 104, 105),
				testCandle(1, 104, 106, 103, 104),
			},
			horizon: 5,
			want:    RangeContextTriageLabelMissingFuture,
		},
	}
	for _, tc := range tests {
		episode := base
		if tc.quality != "" {
			episode.QualityBucket = tc.quality
		}
		row := rangeContextTriageFailureRow(tc.candles, episode, tc.horizon, cfg)
		if row.PrimaryContextLabel != tc.want {
			t.Fatalf("%s label=%s, want %s row=%+v", tc.name, row.PrimaryContextLabel, tc.want, row)
		}
	}
}

func TestRangeContextTriageCohortGatesRankingAndStopStates(t *testing.T) {
	cfg := testRangeContextTriageConfig()
	cfg.MinFullCohortCount = 3
	cfg.MinSplitCohortCount = 1
	cfg.MinSessionSplitCohortCount = 1
	source := []FuturesRangeContextTriageSourceRow{{SourceFactsPass: true, ValidationStatus: "accepted"}}
	coverage := []FuturesRangeContextTriageCoverageRow{{CoverageFactsPass: true, FuturesRangeDiscoveryCoverageRow: FuturesRangeDiscoveryCoverageRow{Complete: true, ValidationStatus: "accepted"}}}
	failures := []FuturesRangeContextTriageFailureModeRow{
		rangeContextTriageGateFailure("2021_2022_stress", RangeDiscoveryTimeframe15m, RangeContextTriageLabelContainedRotation),
		rangeContextTriageGateFailure("2023_2024_oos", RangeDiscoveryTimeframe15m, RangeContextTriageLabelCleanExpansionUp),
		rangeContextTriageGateFailure("2025_2026_recent", RangeDiscoveryTimeframe15m, RangeContextTriageLabelCleanExpansionDown),
		rangeContextTriageGateFailure("2021_2022_stress", RangeDiscoveryTimeframe1h, RangeContextTriageLabelContainedRotation),
		rangeContextTriageGateFailure("2023_2024_oos", RangeDiscoveryTimeframe1h, RangeContextTriageLabelCleanExpansionUp),
		rangeContextTriageGateFailure("2025_2026_recent", RangeDiscoveryTimeframe1h, RangeContextTriageLabelCleanExpansionDown),
	}
	cohorts := rangeContextTriageCohortRows(failures, source, coverage, cfg, DefaultSplits())
	rankings := rangeContextTriageRankingRows(cohorts, cfg, DefaultSplits())
	if len(rankings) == 0 || !rankings[0].PassesGate {
		t.Fatalf("expected passing rankings, got %+v", rankings)
	}
	if rankings[0].Timeframe != RangeDiscoveryTimeframe1h {
		t.Fatalf("1h tie-break should win, got %+v", rankings[:2])
	}
	result := FuturesRangeContextTriageAuditResult{
		SourceRows:      source,
		CoverageRows:    coverage,
		EpisodeRows:     []FuturesRangeContextTriageEpisodeRow{{Eligible: true}},
		FailureModeRows: failures,
		CohortRows:      cohorts,
		RankingRows:     rankings,
	}
	if got := FuturesRangeContextTriageAuditStopState(result); got != RangeContextTriageStopStateReadyForStrategySpec {
		t.Fatalf("stop=%s, want ready", got)
	}

	toxic := []FuturesRangeContextTriageFailureModeRow{
		rangeContextTriageGateFailure("2021_2022_stress", RangeDiscoveryTimeframe15m, RangeContextTriageLabelNoResolution),
		rangeContextTriageGateFailure("2023_2024_oos", RangeDiscoveryTimeframe15m, RangeContextTriageLabelNoResolution),
		rangeContextTriageGateFailure("2025_2026_recent", RangeDiscoveryTimeframe15m, RangeContextTriageLabelNoResolution),
	}
	cohorts = rangeContextTriageCohortRows(toxic, source, coverage, cfg, DefaultSplits())
	rankings = rangeContextTriageRankingRows(cohorts, cfg, DefaultSplits())
	result.CohortRows = cohorts
	result.RankingRows = rankings
	if got := FuturesRangeContextTriageAuditStopState(result); got != RangeContextTriageStopStateFailedNoStrategyPremise {
		t.Fatalf("stop=%s, want failed", got)
	}

	empty := FuturesRangeContextTriageAuditResult{SourceRows: source, CoverageRows: coverage}
	if got := FuturesRangeContextTriageAuditStopState(empty); got != RangeContextTriageStopStateNoRangeEpisodes {
		t.Fatalf("empty stop=%s, want no range episodes", got)
	}
	noCohorts := FuturesRangeContextTriageAuditResult{SourceRows: source, CoverageRows: coverage, EpisodeRows: []FuturesRangeContextTriageEpisodeRow{{Eligible: true}}}
	if got := FuturesRangeContextTriageAuditStopState(noCohorts); got != RangeContextTriageStopStateNoUsableCohorts {
		t.Fatalf("no cohort stop=%s, want no usable cohorts", got)
	}
}

func testRangeContextTriageConfig() FuturesRangeContextTriageAuditConfig {
	cfg := DefaultFuturesRangeContextTriageAuditConfig()
	cfg.SkipSourceFactCheck = true
	cfg.SkipCoverageCountCheck = true
	cfg.DetectorLookbackBarsOverride = 1
	cfg.DetectorMinConsecutiveBars = 1
	cfg.HorizonsBars = []int{2}
	cfg.QuickFailureHorizonBars = 1
	cfg.ReentryWindowBars = 2
	cfg.CleanExpansionThreshold = 0.75
	cfg.DriftThreshold = 0.50
	cfg.BoundaryChopTransitions = 3
	cfg.MinFullCohortCount = 3
	cfg.MinSplitCohortCount = 1
	cfg.MinSessionSplitCohortCount = 1
	return cfg
}

func rangeContextTriageTestEpisode() FuturesRangeContextTriageEpisodeRow {
	return FuturesRangeContextTriageEpisodeRow{
		EpisodeID:       1,
		Timeframe:       RangeDiscoveryTimeframe1h,
		Split:           fullSplitName,
		MatureIndex:     0,
		RawEndIndex:     0,
		MatureCloseTime: testCandle(0, 105, 106, 104, 105).CloseTime.UTC().Format(timeLayout),
		MatureSession:   RangeContextTriageSessionAsia,
		High:            110,
		Low:             100,
		Mid:             105,
		Width:           10,
		WidthPct:        0.10,
		WidthToATRRatio: 2,
		QualityBucket:   RangeContextTriageQualityBalancedOrderly,
		Eligible:        true,
	}
}

func rangeContextTriageGateFailure(split string, timeframe string, label string) FuturesRangeContextTriageFailureModeRow {
	return FuturesRangeContextTriageFailureModeRow{
		EpisodeID:           1,
		Timeframe:           timeframe,
		Split:               split,
		HorizonBars:         12,
		MatureCloseTime:     "2023-01-01T00:00:00Z",
		MatureSession:       RangeContextTriageSessionAsia,
		QualityBucket:       RangeContextTriageQualityBalancedOrderly,
		PrimaryContextLabel: label,
		ConstructiveContext: rangeContextTriageConstructiveLabel(label),
		ToxicContext:        rangeContextTriageToxicLabel(label),
	}
}
