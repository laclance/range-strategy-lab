package lab

import (
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestBTCRegimeETHSOLContextSourceGapOnUnapprovedPath(t *testing.T) {
	dir := t.TempDir()
	candles := make([]Candle, 12)
	for i := range candles {
		candles[i] = testCandle(i, 100, 101, 99, 100)
	}
	btc := writeRangeUniverseTestCSV(t, dir, "btcusdt_futures_um_5m_context.csv", candles)
	eth := writeRangeUniverseTestCSV(t, dir, "ethusdt_futures_um_5m_context.csv", candles)
	sol := writeRangeUniverseTestCSV(t, dir, "solusdt_futures_um_5m_context.csv", candles)

	cfg := testBTCRegimeETHSOLContextConfig(btc, eth, sol)
	cfg.Sources[1].ApprovedPath = filepath.Join(dir, "other_ethusdt_futures_um_5m_context.csv")
	result, err := RunFuturesBTCRegimeETHSOLContextAudit(cfg, []Split{{Name: fullSplitName}})
	if err != nil {
		t.Fatal(err)
	}
	if result.StopState != BTCRegimeETHSOLContextStopStateSourceGap {
		t.Fatalf("stop=%s, want source gap", result.StopState)
	}
	if len(result.SourceRows) != 3 || result.SourceRows[1].SourceFactsPass || !strings.Contains(result.SourceRows[1].ValidationError, "approved local path") {
		t.Fatalf("expected ETH approved-path rejection, sources=%+v", result.SourceRows)
	}
}

func TestBTCRegimeETHSOLContextStateIDsUseClosedInputsOnly(t *testing.T) {
	localState := routerTestState(1, RangeDiscoveryTimeframe1h, "range_state_v1::1h::geometry_balanced_orderly::vol_normal::trend_flat::impulse_none::participation_normal")
	localState.GeometryBucket = "geometry_balanced_orderly"
	localState.VolBucket = "vol_normal"
	localState.TrendBucket = "trend_flat"
	localState.ImpulseBucket = "impulse_none"
	localState.ParticipationBucket = "participation_normal"
	localState.ReturnShortPct = 0.01
	localState.ReturnMediumPct = 0.02
	btc := FuturesBTCRegimeETHSOLContextBTCStateRow{
		BTCStateRowID:     9,
		Timeframe:         RangeDiscoveryTimeframe1h,
		DecisionCloseTime: localState.DecisionCloseTime,
		BTCRegimeID:       "1h::geometry_balanced_orderly::vol_normal::trend_flat::impulse_none::participation_normal::close_mid::age_mature",
		ReturnShortPct:    0.00,
		ReturnMediumPct:   0.00,
	}
	rs := btcRegimeETHSOLContextRelativeStrengthRow(1, RangeUniverseSymbolETHUSDT, localState, btc)
	local := btcRegimeETHSOLContextLocalStateRow(RangeUniverseSymbolETHUSDT, localState, btc, rs.RelativeStrengthBucket)
	label := rangeStateTestLabel(fullSplitName, RangeContextTriageLabelContainedRotation)
	contextLabel := btcRegimeETHSOLContextLabelRow(1, RangeUniverseSymbolETHUSDT, local, label)

	if !local.ClosedCandleOnly || local.ForwardLabelsAsStateInput {
		t.Fatalf("local state should be closed-candle only with no forward inputs: %+v", local)
	}
	if strings.Contains(local.ContextStateID, label.ForwardLabel) {
		t.Fatalf("context state ID leaked forward label: state_id=%s label=%s", local.ContextStateID, label.ForwardLabel)
	}
	if !contextLabel.ForwardLabelMetadataOnly || contextLabel.ForwardLabelUsedAsFeature {
		t.Fatalf("label row should keep forward label as metadata only: %+v", contextLabel)
	}
}

func TestBTCRegimeETHSOLContextCohortRequiresBTCImprovement(t *testing.T) {
	cfg := DefaultFuturesBTCRegimeETHSOLContextAuditConfig()
	cfg.MinFullCohortRows = 90
	cfg.MinSplitCohortRows = 30
	cfg.MinUsefulRateFull = 0.70
	cfg.MinUsefulRateSplit = 0.70
	cfg.MaxToxicRateFull = 0.20
	cfg.MaxToxicRateSplit = 0.20
	cfg.MinUsefulMinusToxicMarginFull = 0.60
	cfg.MinUsefulMinusToxicMarginSplit = 0.60
	cfg.MinContextUsefulImprovementFull = 0.20
	cfg.MinContextUsefulImprovementSplit = 0.20
	cfg.MinContextMarginImprovementFull = 0.20
	cfg.MinContextMarginImprovementSplit = 0.20
	labels := []FuturesBTCRegimeETHSOLContextLabelRow{}
	for _, split := range rangeDiscoveryPeriodSplits(DefaultSplits()) {
		for i := 0; i < 40; i++ {
			labels = append(labels, btcRegimeETHSOLContextTestLabel(split.Name, "btc_good", RangeContextTriageLabelContainedRotation))
		}
		for i := 0; i < 40; i++ {
			labels = append(labels, btcRegimeETHSOLContextTestLabel(split.Name, "btc_bad", RangeContextTriageLabelBoundaryChop))
		}
	}
	cohorts := btcRegimeETHSOLContextCohortRows(labels, cfg, DefaultSplits(), true)
	rankings := btcRegimeETHSOLContextRankingRows(cohorts)
	passing := false
	for _, row := range rankings {
		if row.PassesGate && row.BTCRegimeID == "btc_good" && row.RouteCandidate == RangeStateConstructionLoopRouteRotation {
			passing = true
		}
		if row.PassesGate && row.BTCRegimeID == "btc_bad" && row.RouteCandidate == RangeStateConstructionLoopRouteRotation {
			t.Fatalf("bad BTC regime should not pass rotation: %+v", row)
		}
	}
	if !passing {
		t.Fatalf("expected passing BTC-improved rotation cohort, top rankings=%+v", rankings[:minInt(3, len(rankings))])
	}

	result := FuturesBTCRegimeETHSOLContextAuditResult{
		SourceRows:     btcRegimeETHSOLContextAcceptedSources(),
		CoverageRows:   btcRegimeETHSOLContextAcceptedCoverage(),
		LocalStateRows: []FuturesBTCRegimeETHSOLContextLocalStateRow{{ClosedCandleOnly: true}},
		RankingRows:    rankings,
	}
	if got := FuturesBTCRegimeETHSOLContextAuditStopState(result); got != BTCRegimeETHSOLContextStopStatePassedNeedsSpec {
		t.Fatalf("stop=%s, want passed needs strategy premise spec", got)
	}
	result.LabelRows = []FuturesBTCRegimeETHSOLContextLabelRow{{ForwardLabelUsedAsFeature: true}}
	if got := FuturesBTCRegimeETHSOLContextAuditStopState(result); got != BTCRegimeETHSOLContextStopStateRejectedFutureLeak {
		t.Fatalf("stop=%s, want future leak rejection", got)
	}
}

func testBTCRegimeETHSOLContextConfig(btc, eth, sol string) FuturesBTCRegimeETHSOLContextAuditConfig {
	cfg := DefaultFuturesBTCRegimeETHSOLContextAuditConfig()
	cfg.Sources = []FuturesBTCRegimeETHSOLContextSourceConfig{
		{Symbol: RangeUniverseSymbolBTCUSDT, Path: btc, ApprovedPath: btc, SkipSourceFactCheck: true, SkipSplitEligibilityCheck: true},
		{Symbol: RangeUniverseSymbolETHUSDT, Path: eth, ApprovedPath: eth, SkipSourceFactCheck: true, SkipSplitEligibilityCheck: true},
		{Symbol: RangeUniverseSymbolSOLUSDT, Path: sol, ApprovedPath: sol, SkipSourceFactCheck: true, SkipSplitEligibilityCheck: true},
	}
	cfg.StateConfig = testRangeStateConstructionLoopConfig()
	cfg.StateConfig.Timeframes = []string{RangeDiscoveryTimeframe15m}
	cfg.StateConfig.Horizons15M = []int{2}
	cfg.MinFullCohortRows = 1
	cfg.MinSplitCohortRows = 1
	return cfg
}

func btcRegimeETHSOLContextTestLabel(split, btcRegimeID, forwardLabel string) FuturesBTCRegimeETHSOLContextLabelRow {
	row := FuturesBTCRegimeETHSOLContextLabelRow{
		Symbol:                    RangeUniverseSymbolETHUSDT,
		Timestamp:                 time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Format(timeLayout),
		Timeframe:                 RangeDiscoveryTimeframe1h,
		Split:                     split,
		HorizonBars:               2,
		LocalRangeBucketID:        "1h::geometry_balanced_orderly::vol_normal::trend_flat::impulse_none",
		BTCRegimeID:               btcRegimeID,
		RelativeStrengthBucket:    "relative_strength_strong",
		ContextStateID:            "btc_regime_eth_sol_v1::ETHUSDT::1h::" + btcRegimeID,
		ForwardLabel:              forwardLabel,
		ForwardLabelMetadataOnly:  true,
		ForwardLabelUsedAsFeature: false,
	}
	row.RotationUseful = forwardLabel == RangeContextTriageLabelContainedRotation || forwardLabel == RangeContextTriageLabelFalseBreakReentryUp || forwardLabel == RangeContextTriageLabelFalseBreakReentryDown
	row.RotationToxic = forwardLabel == RangeContextTriageLabelBoundaryChop || forwardLabel == RangeContextTriageLabelCleanExpansionUp || forwardLabel == RangeContextTriageLabelCleanExpansionDown
	row.ContinuationUseful = forwardLabel == RangeContextTriageLabelCleanExpansionUp || forwardLabel == RangeContextTriageLabelCleanExpansionDown || forwardLabel == RangeContextTriageLabelDriftThroughUp || forwardLabel == RangeContextTriageLabelDriftThroughDown
	row.ContinuationToxic = forwardLabel == RangeContextTriageLabelFalseBreakReentryUp || forwardLabel == RangeContextTriageLabelFalseBreakReentryDown || forwardLabel == RangeContextTriageLabelBoundaryChop
	row.NoTradeToxic = forwardLabel == RangeContextTriageLabelBoundaryChop || forwardLabel == RangeContextTriageLabelLowWidthNoise || forwardLabel == RangeContextTriageLabelNoResolution
	return row
}

func btcRegimeETHSOLContextAcceptedSources() []FuturesBTCRegimeETHSOLContextSourceRow {
	return []FuturesBTCRegimeETHSOLContextSourceRow{
		{SourceFactsPass: true, FuturesRangeUniverseSourceRow: FuturesRangeUniverseSourceRow{ValidationStatus: "accepted"}},
		{SourceFactsPass: true, FuturesRangeUniverseSourceRow: FuturesRangeUniverseSourceRow{ValidationStatus: "accepted"}},
		{SourceFactsPass: true, FuturesRangeUniverseSourceRow: FuturesRangeUniverseSourceRow{ValidationStatus: "accepted"}},
	}
}

func btcRegimeETHSOLContextAcceptedCoverage() []FuturesBTCRegimeETHSOLContextCoverageRow {
	return []FuturesBTCRegimeETHSOLContextCoverageRow{{
		CoverageFactsPass: true,
		FuturesRangeDiscoveryCoverageRow: FuturesRangeDiscoveryCoverageRow{
			ValidationStatus: "accepted",
			Complete:         true,
		},
	}}
}
