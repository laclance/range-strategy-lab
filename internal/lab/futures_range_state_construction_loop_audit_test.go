package lab

import (
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRangeStateConstructionLoopSourceAndCoverageAcceptance(t *testing.T) {
	parent := make([]Candle, 96)
	for i := range parent {
		parent[i] = testCandle(i, 100, 101, 99, 100)
		if i < rangeStateConstructionLoopExpectedZeroVol {
			parent[i].Volume = 0
		}
	}
	manifest, err := ValidateResearchSource("btcusdt_futures_um_5m_state.csv", []string{"open_time", "open", "high", "low", "close", "volume", "close_time"}, parent, SourceValidationOptions{Product: SourceProductBinanceUSDMFutures})
	if err != nil {
		t.Fatal(err)
	}
	cfg := testRangeStateConstructionLoopConfig()
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

	result, err := RunFuturesRangeStateConstructionLoopAudit(parent, manifest, cfg, []Split{{Name: fullSplitName}})
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
		if !row.CoverageFactsPass || row.ValidationStatus != "accepted" || !row.Complete {
			t.Fatalf("bad coverage row: %+v", row)
		}
	}

	manifest.Product = "Binance spot"
	manifest.ComparisonOnly = true
	result, err = RunFuturesRangeStateConstructionLoopAudit(parent, manifest, cfg, []Split{{Name: fullSplitName}})
	if err != nil {
		t.Fatal(err)
	}
	if result.StopState != RangeStateConstructionLoopStopStateSourceGap {
		t.Fatalf("stop=%s, want source gap", result.StopState)
	}
}

func TestRangeStateConstructionLoopSourcePathRejection(t *testing.T) {
	dir := t.TempDir()
	parent := make([]Candle, 24)
	for i := range parent {
		parent[i] = testCandle(i, 100, 101, 99, 100)
	}
	path := writeRangeUniverseTestCSV(t, dir, "btcusdt_futures_um_5m_state.csv", parent)
	candles, manifest, err := LoadResearchSourceCSV(path, SourceValidationOptions{Product: SourceProductBinanceUSDMFutures})
	if err != nil {
		t.Fatal(err)
	}
	cfg := testRangeStateConstructionLoopConfig()
	cfg.ApprovedSourcePath = filepath.Join(dir, "other_btcusdt_futures_um_5m.csv")
	result, err := RunFuturesRangeStateConstructionLoopAudit(candles, manifest, cfg, []Split{{Name: fullSplitName}})
	if err != nil {
		t.Fatal(err)
	}
	if result.StopState != RangeStateConstructionLoopStopStateSourceGap || !strings.Contains(result.SourceRows[0].ValidationError, "approved path") {
		t.Fatalf("expected approved-path source gap, got stop=%s source=%+v", result.StopState, result.SourceRows)
	}
}

func TestRangeStateConstructionLoopStateUsesClosedDecisionFeaturesOnly(t *testing.T) {
	candles := make([]Candle, 70)
	for i := range candles {
		close := 105.0
		if i%2 == 0 {
			close = 106
		}
		candles[i] = testCandle(i, close-1, close+2, close-2, close)
		candles[i].Volume = 100 + float64(i)
	}
	candles[61].High = 200
	cfg := testRangeStateConstructionLoopConfig()
	cfg.ShortWindowBars = 4
	cfg.MediumWindowBars = 8
	cfg.FeatureLookbackBars = 8
	data := rangeStateFrameData{
		frame:   rangeDiscoveryFrameDef{timeframe: RangeDiscoveryTimeframe1h, barsPerDay: 24},
		candles: candles,
		metrics: rangeStateMetrics{
			normalizedATR:         rangeStateTestFloatSlice(len(candles), 0.01),
			atr:                   rangeStateTestFloatSlice(len(candles), 1),
			adx:                   rangeStateTestFloatSlice(len(candles), 20),
			atrPercentile:         rangeStateTestFloatSlice(len(candles), 0.5),
			realizedVolPercentile: rangeStateTestFloatSlice(len(candles), 0.5),
			volumePercentile:      rangeStateTestFloatSlice(len(candles), 0.5),
			trueRange:             rangeStateTestFloatSlice(len(candles), 1),
		},
		horizons:   []int{2},
		maxHorizon: 2,
	}
	row, ok, reason := rangeStateConstructionLoopStateRow(data, map[string]rangeStateFrameData{RangeDiscoveryTimeframe1h: data}, cfg, []Split{{Name: fullSplitName}}, 1, 1, 50, 60, 110, 100)
	if !ok {
		t.Fatalf("state row rejected: %s", reason)
	}
	if row.RangeHigh != 110 {
		t.Fatalf("state row used future high, got range high %.2f", row.RangeHigh)
	}
	if row.StateID == "" || !strings.Contains(row.StateID, row.GeometryBucket) || !strings.Contains(row.StateID, row.VolBucket) || !strings.Contains(row.StateID, row.TrendBucket) || !strings.Contains(row.StateID, row.ImpulseBucket) || !strings.Contains(row.StateID, row.ParticipationBucket) {
		t.Fatalf("state id missing buckets: %+v", row)
	}
	label := rangeStateConstructionLoopLabelRow(candles, row, 2, cfg)
	if label.ForwardLabel == "" {
		t.Fatalf("expected forward label")
	}
	if strings.Contains(row.StateID, label.ForwardLabel) || label.ForwardLabelUsedAsFeature {
		t.Fatalf("forward label leaked into feature state: state=%s label=%s", row.StateID, label.ForwardLabel)
	}
}

func TestRangeStateConstructionLoopCohortGatesAndStopStates(t *testing.T) {
	cfg := testRangeStateConstructionLoopConfig()
	cfg.MinFullCohortCount = 100
	cfg.MinSplitCohortCount = 30
	splits := DefaultSplits()
	labels := []FuturesRangeStateConstructionLoopLabelRow{}
	for _, split := range rangeDiscoveryPeriodSplits(splits) {
		for i := 0; i < 25; i++ {
			labels = append(labels, rangeStateTestLabel(split.Name, RangeContextTriageLabelContainedRotation))
		}
		for i := 0; i < 15; i++ {
			labels = append(labels, rangeStateTestLabel(split.Name, RangeContextTriageLabelFalseBreakReentryUp))
		}
		for i := 0; i < 10; i++ {
			labels = append(labels, rangeStateTestLabel(split.Name, RangeContextTriageLabelNoResolution))
		}
	}
	cohorts := rangeStateConstructionLoopCohortRows(labels, rangeStateAcceptedSources(), rangeStateAcceptedCoverage(), cfg, splits)
	rankings := rangeStateConstructionLoopRankingRows(cohorts, splits)
	passingRotation := false
	for _, row := range rankings {
		if row.PassesGate && row.RouteCandidate == RangeStateConstructionLoopRouteRotation {
			passingRotation = true
		}
		if row.PassesGate && row.RouteCandidate == RangeStateConstructionLoopRouteDiagnosticOnly {
			t.Fatalf("diagnostic route must not pass: %+v", row)
		}
	}
	if !passingRotation {
		t.Fatalf("expected a passing rotation cohort, rankings=%+v", rankings[:minInt(3, len(rankings))])
	}

	result := FuturesRangeStateConstructionLoopAuditResult{
		SourceRows:   rangeStateAcceptedSources(),
		CoverageRows: rangeStateAcceptedCoverage(),
		StateRows:    []FuturesRangeStateConstructionLoopStateRow{{Eligible: true}},
		RankingRows:  []FuturesRangeStateConstructionLoopRankingRow{{PassesGate: true, RouteCandidate: RangeStateConstructionLoopRouteRotation}},
	}
	if got := FuturesRangeStateConstructionLoopAuditStopState(result); got != RangeStateConstructionLoopStopStatePassedNeedsStrategySpec {
		t.Fatalf("stop=%s, want strategy premise", got)
	}
	result.RankingRows = []FuturesRangeStateConstructionLoopRankingRow{{PassesGate: true, RouteCandidate: RangeStateConstructionLoopRouteNoTradeToxic}}
	if got := FuturesRangeStateConstructionLoopAuditStopState(result); got != RangeStateConstructionLoopStopStatePassedNoTradeFilterOnly {
		t.Fatalf("stop=%s, want no-trade filter", got)
	}
	result.StateRows = []FuturesRangeStateConstructionLoopStateRow{{SkippedReason: "closed_family_reslice"}}
	if got := FuturesRangeStateConstructionLoopAuditStopState(result); got != RangeStateConstructionLoopStopStateRejectedClosedReslice {
		t.Fatalf("stop=%s, want closed-family reslice", got)
	}
}

func testRangeStateConstructionLoopConfig() FuturesRangeStateConstructionLoopAuditConfig {
	cfg := DefaultFuturesRangeStateConstructionLoopAuditConfig()
	cfg.SkipSourceFactCheck = true
	cfg.SkipCoverageCountCheck = true
	cfg.DetectorLookbackBarsOverride = 1
	cfg.DetectorMinConsecutiveBars = 1
	cfg.ShortWindowBars = 2
	cfg.MediumWindowBars = 4
	cfg.FeatureLookbackBars = 4
	cfg.LongLookbackDays = 1
	cfg.Horizons15M = []int{2}
	cfg.Horizons1H = []int{2}
	cfg.Horizons4H = []int{1}
	cfg.MinFullCohortCount = 1
	cfg.MinSplitCohortCount = 1
	cfg.MinNoTradeFullCohortCount = 1
	cfg.MinNoTradeSplitCohortCount = 1
	return cfg
}

func rangeStateTestFloatSlice(n int, value float64) []float64 {
	out := make([]float64, n)
	for i := range out {
		out[i] = value
	}
	return out
}

func rangeStateTestLabel(split string, label string) FuturesRangeStateConstructionLoopLabelRow {
	row := FuturesRangeStateConstructionLoopLabelRow{
		Timeframe:                 RangeDiscoveryTimeframe1h,
		Split:                     split,
		HorizonBars:               2,
		ForwardLabel:              label,
		GeometryVolID:             "1h::geometry_balanced_orderly::vol_normal",
		GeometryTrendID:           "1h::geometry_balanced_orderly::trend_flat",
		GeometryImpulseID:         "1h::geometry_balanced_orderly::impulse_none",
		GeometryParticipationID:   "1h::geometry_balanced_orderly::participation_normal",
		GeometryVolTrendID:        "1h::geometry_balanced_orderly::vol_normal::trend_flat",
		GeometryVolTrendImpulseID: "1h::geometry_balanced_orderly::vol_normal::trend_flat::impulse_none",
		AllFamiliesID:             "range_state_v1::1h::geometry_balanced_orderly::vol_normal::trend_flat::impulse_none::participation_normal",
		Timestamp:                 time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Format(timeLayout),
	}
	return rangeStateFinalizeLabelRow(row)
}

func rangeStateAcceptedSources() []FuturesRangeStateConstructionLoopSourceRow {
	return []FuturesRangeStateConstructionLoopSourceRow{{SourceFactsPass: true, ValidationStatus: "accepted"}}
}

func rangeStateAcceptedCoverage() []FuturesRangeStateConstructionLoopCoverageRow {
	return []FuturesRangeStateConstructionLoopCoverageRow{{
		CoverageFactsPass: true,
		FuturesRangeDiscoveryCoverageRow: FuturesRangeDiscoveryCoverageRow{
			ValidationStatus: "accepted",
			Complete:         true,
		},
	}}
}
