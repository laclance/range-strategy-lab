package lab

import (
	"strings"
	"testing"
	"time"
)

func TestRangeContextRouterSourceAndCoverageAcceptanceAndSpotRejection(t *testing.T) {
	parent := make([]Candle, 96)
	for i := range parent {
		parent[i] = testCandle(i, 100, 101, 99, 100)
		if i < rangeStateConstructionLoopExpectedZeroVol {
			parent[i].Volume = 0
		}
	}
	manifest, err := ValidateResearchSource("btcusdt_futures_um_5m_router.csv", []string{"open_time", "open", "high", "low", "close", "volume", "close_time"}, parent, SourceValidationOptions{Product: SourceProductBinanceUSDMFutures})
	if err != nil {
		t.Fatal(err)
	}
	cfg := testRangeContextRouterConfig()
	cfg.StateAuditConfig.ApprovedSourcePath = manifest.Path
	cfg.StateAuditConfig.SkipSourceFactCheck = false
	cfg.StateAuditConfig.ExpectedSourceRows = len(parent)
	cfg.StateAuditConfig.ExpectedFirstOpenTime = parent[0].OpenTime.UTC().Format(timeLayout)
	cfg.StateAuditConfig.ExpectedLastOpenTime = parent[len(parent)-1].OpenTime.UTC().Format(timeLayout)
	cfg.StateAuditConfig.SkipCoverageCountCheck = false
	cfg.StateAuditConfig.Expected15MRows = 32
	cfg.StateAuditConfig.Expected15MLastOpenTime = parent[93].OpenTime.UTC().Format(timeLayout)
	cfg.StateAuditConfig.Expected1HRows = 8
	cfg.StateAuditConfig.Expected1HLastOpenTime = parent[84].OpenTime.UTC().Format(timeLayout)
	cfg.StateAuditConfig.Expected4HRows = 2
	cfg.StateAuditConfig.Expected4HLastOpenTime = parent[48].OpenTime.UTC().Format(timeLayout)

	result, err := RunFuturesRangeContextRouterAudit(parent, manifest, cfg, []Split{{Name: fullSplitName}})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.SourceRows) != 1 || !result.SourceRows[0].SourceFactsPass || result.SourceRows[0].RouterUsesForwardLabels {
		t.Fatalf("bad router source row: %+v", result.SourceRows)
	}
	if len(result.CoverageRows) != 3 {
		t.Fatalf("coverage rows=%d, want 3", len(result.CoverageRows))
	}

	manifest.Product = "Binance spot"
	manifest.ComparisonOnly = true
	result, err = RunFuturesRangeContextRouterAudit(parent, manifest, cfg, []Split{{Name: fullSplitName}})
	if err != nil {
		t.Fatal(err)
	}
	if result.StopState != RangeContextRouterStopStateSourceGap {
		t.Fatalf("stop=%s, want source gap", result.StopState)
	}
}

func TestRangeContextRouterRulesUseStateRollupsAndDoNotUseForwardLabels(t *testing.T) {
	noTradeState := routerTestState(1, "15m", "range_state_v1::15m::geometry_balanced_orderly::vol_compressed::trend_down_pressure::impulse_none::participation_low")
	rotationState := routerTestState(2, "15m", "range_state_v1::15m::geometry_wide_volatile::vol_compressed::trend_flat::impulse_stale::participation_low")
	conflictState := routerTestState(3, "15m", "range_state_v1::15m::geometry_conflict::vol_compressed::trend_flat::impulse_none::participation_low")
	conflictState.GeometryVolID = "15m::geometry_conflict::vol_compressed"

	rankings := []FuturesRangeStateConstructionLoopRankingRow{
		{
			Rank:                     1,
			CohortID:                 "toxic",
			RouteCandidate:           RangeStateConstructionLoopRouteNoTradeToxic,
			RollupType:               RangeStateConstructionLoopRollupAllFamilies,
			RollupID:                 noTradeState.AllFamiliesID,
			Timeframe:                "15m",
			HorizonBars:              48,
			PassesGate:               true,
			FullToxicRate:            0.7,
			WorstSplitToxicRate:      0.6,
			FutureLeakProtectionPass: true,
		},
		{
			Rank:                     2,
			CohortID:                 "rotation",
			RouteCandidate:           RangeStateConstructionLoopRouteRotation,
			RollupType:               RangeStateConstructionLoopRollupAllFamilies,
			RollupID:                 rotationState.AllFamiliesID,
			Timeframe:                "15m",
			HorizonBars:              24,
			PassesGate:               true,
			FullUsefulRate:           0.7,
			WeakestSplitUsefulRate:   0.6,
			FutureLeakProtectionPass: true,
		},
		{
			Rank:                     3,
			CohortID:                 "conflict-toxic",
			RouteCandidate:           RangeStateConstructionLoopRouteNoTradeToxic,
			RollupType:               RangeStateConstructionLoopRollupGeometryVol,
			RollupID:                 conflictState.GeometryVolID,
			Timeframe:                "15m",
			HorizonBars:              48,
			PassesGate:               true,
			FutureLeakProtectionPass: true,
		},
		{
			Rank:                     4,
			CohortID:                 "conflict-rotation",
			RouteCandidate:           RangeStateConstructionLoopRouteRotation,
			RollupType:               RangeStateConstructionLoopRollupAllFamilies,
			RollupID:                 conflictState.AllFamiliesID,
			Timeframe:                "15m",
			HorizonBars:              24,
			PassesGate:               true,
			FutureLeakProtectionPass: true,
		},
	}
	rules := rangeContextRouterRuleRows(rankings)
	rows := rangeContextRouterRows([]FuturesRangeStateConstructionLoopStateRow{noTradeState, rotationState, conflictState}, rules)
	byState := map[int]FuturesRangeContextRouterRow{}
	for _, row := range rows {
		byState[row.StateRowID] = row
		if row.ForwardLabelsAsRouterInput || row.ForwardLabelColumnsPresent {
			t.Fatalf("router row uses forward labels: %+v", row)
		}
		if strings.Contains(row.MatchedRuleIDs, RangeContextTriageLabelContainedRotation) || strings.Contains(row.RouterReason, RangeContextTriageLabelBoundaryChop) {
			t.Fatalf("forward label leaked into router fields: %+v", row)
		}
	}
	if byState[1].RouterLabel != RangeContextRouterLabelNoTrade {
		t.Fatalf("state 1 label=%s, want no_trade", byState[1].RouterLabel)
	}
	if byState[2].RouterLabel != RangeContextRouterLabelTradableRotation {
		t.Fatalf("state 2 label=%s, want tradable_rotation", byState[2].RouterLabel)
	}
	if byState[3].RouterLabel != RangeContextRouterLabelDiagnosticOnly || !byState[3].ConflictingRuleMatch {
		t.Fatalf("state 3 should be diagnostic conflict, got %+v", byState[3])
	}
}

func TestRangeContextRouterCohortGatesAndStopStates(t *testing.T) {
	cfg := testRangeContextRouterConfig()
	cfg.StateAuditConfig.Timeframes = []string{RangeDiscoveryTimeframe1h}
	cfg.StateAuditConfig.Horizons1H = []int{2}
	labels := []FuturesRangeStateConstructionLoopLabelRow{}
	rows := []FuturesRangeContextRouterRow{}
	stateID := 1
	for _, split := range rangeDiscoveryPeriodSplits(DefaultSplits()) {
		for i := 0; i < 30; i++ {
			label := rangeStateTestLabel(split.Name, RangeContextTriageLabelContainedRotation)
			label.StateRowID = stateID
			labels = append(labels, label)
			rows = append(rows, routerTestRow(stateID, split.Name, RangeContextRouterLabelTradableRotation))
			stateID++
		}
		for i := 0; i < 5; i++ {
			label := rangeStateTestLabel(split.Name, RangeContextTriageLabelNoResolution)
			label.StateRowID = stateID
			labels = append(labels, label)
			rows = append(rows, routerTestRow(stateID, split.Name, RangeContextRouterLabelTradableRotation))
			stateID++
		}
	}
	cohorts := rangeContextRouterCohortRows(rows, labels, cfg, DefaultSplits(), true)
	rankings := rangeContextRouterRankingRows(cohorts)
	passingRotation := false
	for _, row := range rankings {
		if row.PassesGate && row.RouterLabel == RangeContextRouterLabelTradableRotation {
			passingRotation = true
		}
		if row.PassesGate && row.RouterLabel == RangeContextRouterLabelDiagnosticOnly {
			t.Fatalf("diagnostic router cohort must not pass: %+v", row)
		}
	}
	if !passingRotation {
		t.Fatalf("expected passing rotation router cohort, rankings=%+v", rankings)
	}

	result := FuturesRangeContextRouterAuditResult{
		SourceRows:   rangeContextRouterAcceptedSources(),
		CoverageRows: rangeContextRouterAcceptedCoverage(),
		Rows:         []FuturesRangeContextRouterRow{{RouterLabel: RangeContextRouterLabelTradableRotation}},
		RankingRows:  []FuturesRangeContextRouterRankingRow{{PassesGate: true, RouterLabel: RangeContextRouterLabelTradableRotation}},
	}
	if got := FuturesRangeContextRouterAuditStopState(result); got != RangeContextRouterStopStatePassedNeedsRotationSpec {
		t.Fatalf("stop=%s, want rotation premise", got)
	}
	result.RankingRows = []FuturesRangeContextRouterRankingRow{{PassesGate: true, RouterLabel: RangeContextRouterLabelNoTrade}}
	if got := FuturesRangeContextRouterAuditStopState(result); got != RangeContextRouterStopStatePassedNoTradeFilterOnly {
		t.Fatalf("stop=%s, want no-trade filter", got)
	}
	result.RankingRows = []FuturesRangeContextRouterRankingRow{{PassesGate: true, RouterLabel: RangeContextRouterLabelTrendContinuation}}
	if got := FuturesRangeContextRouterAuditStopState(result); got != RangeContextRouterStopStatePassedNeedsContinuation {
		t.Fatalf("stop=%s, want continuation premise", got)
	}
	result.RankingRows = nil
	if got := FuturesRangeContextRouterAuditStopState(result); got != RangeContextRouterStopStateFailedNoActionableRoute {
		t.Fatalf("stop=%s, want failed no actionable", got)
	}
	result.RuleRows = []FuturesRangeContextRouterRuleRow{{ClosedFamilyReslice: true}}
	if got := FuturesRangeContextRouterAuditStopState(result); got != RangeContextRouterStopStateRejectedClosedReslice {
		t.Fatalf("stop=%s, want closed-family rejection", got)
	}
}

func testRangeContextRouterConfig() FuturesRangeContextRouterAuditConfig {
	cfg := DefaultFuturesRangeContextRouterAuditConfig()
	cfg.StateAuditConfig = testRangeStateConstructionLoopConfig()
	cfg.MinFullRouterRows = 1
	cfg.MinSplitRouterRows = 1
	cfg.MinNoTradeFullRouterRows = 1
	cfg.MinNoTradeSplitRouterRows = 1
	return cfg
}

func routerTestState(id int, timeframe string, allFamiliesID string) FuturesRangeStateConstructionLoopStateRow {
	return FuturesRangeStateConstructionLoopStateRow{
		StateRowID:                id,
		Timestamp:                 time.Date(2024, 1, 1, 0, id, 0, 0, time.UTC).Format(timeLayout),
		Timeframe:                 timeframe,
		Split:                     fullSplitName,
		RangeEpisodeID:            1,
		RangeStartIndex:           1,
		DecisionIndex:             id,
		RangeStartTime:            time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Format(timeLayout),
		DecisionCloseTime:         time.Date(2024, 1, 1, 0, id, 0, 0, time.UTC).Format(timeLayout),
		StateID:                   allFamiliesID,
		GeometryVolID:             timeframe + "::geometry_wide_volatile::vol_compressed",
		GeometryTrendID:           timeframe + "::geometry_wide_volatile::trend_flat",
		GeometryImpulseID:         timeframe + "::geometry_wide_volatile::impulse_none",
		GeometryParticipationID:   timeframe + "::geometry_wide_volatile::participation_low",
		GeometryVolTrendID:        timeframe + "::geometry_wide_volatile::vol_compressed::trend_flat",
		GeometryVolTrendImpulseID: timeframe + "::geometry_wide_volatile::vol_compressed::trend_flat::impulse_none",
		AllFamiliesID:             allFamiliesID,
		Eligible:                  true,
	}
}

func routerTestRow(stateID int, split string, routerLabel string) FuturesRangeContextRouterRow {
	return FuturesRangeContextRouterRow{
		RouterRowID:                stateID,
		StateRowID:                 stateID,
		Timeframe:                  RangeDiscoveryTimeframe1h,
		Split:                      split,
		DecisionIndex:              stateID,
		RouterLabel:                routerLabel,
		Actionable:                 routerLabel != RangeContextRouterLabelDiagnosticOnly,
		ClosedCandleOnly:           true,
		ForwardLabelsAsRouterInput: false,
	}
}

func rangeContextRouterAcceptedSources() []FuturesRangeContextRouterSourceRow {
	return []FuturesRangeContextRouterSourceRow{{
		SourceResamplePass:      true,
		RouterUsesForwardLabels: false,
		FuturesRangeStateConstructionLoopSourceRow: FuturesRangeStateConstructionLoopSourceRow{
			SourceFactsPass:  true,
			ValidationStatus: "accepted",
		},
	}}
}

func rangeContextRouterAcceptedCoverage() []FuturesRangeContextRouterCoverageRow {
	return []FuturesRangeContextRouterCoverageRow{{
		SourceResamplePass: true,
		FuturesRangeStateConstructionLoopCoverageRow: FuturesRangeStateConstructionLoopCoverageRow{
			CoverageFactsPass: true,
			FuturesRangeDiscoveryCoverageRow: FuturesRangeDiscoveryCoverageRow{
				ValidationStatus: "accepted",
				Complete:         true,
			},
		},
	}}
}
