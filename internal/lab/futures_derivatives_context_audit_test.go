package lab

import (
	"strings"
	"testing"
	"time"
)

func TestDerivativesContextFeatureUsesLaggedClosedSourceRows(t *testing.T) {
	decisionClose := time.Date(2024, 1, 1, 0, 14, 59, 0, time.UTC)
	sourceOpen := time.Date(2024, 1, 1, 0, 5, 0, 0, time.UTC).UnixMilli()
	futureOpen := time.Date(2024, 1, 1, 0, 10, 0, 0, time.UTC).UnixMilli()
	priorOpen := sourceOpen - 12*derivativesIntervalMs
	state := routerTestState(7, RangeDiscoveryTimeframe15m, "range_state_v1::15m::geometry_balanced_orderly::vol_normal::trend_flat::impulse_none::participation_normal")
	state.DecisionCloseTime = decisionClose.Format(timeLayout)
	state.Timestamp = state.DecisionCloseTime
	state.GeometryBucket = "geometry_balanced_orderly"
	state.VolBucket = "vol_normal"
	state.TrendBucket = "trend_flat"
	state.ImpulseBucket = "impulse_none"
	state.ParticipationBucket = "participation_normal"
	streams := map[string]*derivStreamData{
		derivativesContextStreamKey(RangeUniverseSymbolBTCUSDT, "mark_price_klines"): {
			openClose: map[int64]float64{
				priorOpen:  100.00,
				sourceOpen: 100.20,
				futureOpen: 110.00,
			},
		},
		derivativesContextStreamKey(RangeUniverseSymbolBTCUSDT, "index_price_klines"): {
			openClose: map[int64]float64{
				priorOpen:  100.00,
				sourceOpen: 100.00,
				futureOpen: 100.00,
			},
		},
		derivativesContextStreamKey(RangeUniverseSymbolBTCUSDT, "premium_index_klines"): {
			openClose: map[int64]float64{sourceOpen: 0.0001},
		},
	}
	cfg := DefaultFuturesDerivativesContextAuditConfig()
	feature, ok, reason := derivativesContextFeatureForState(1, RangeUniverseSymbolBTCUSDT, state, streams, cfg)
	if !ok {
		t.Fatalf("feature rejected: %s", reason)
	}
	if feature.SourceOpenTime != "2024-01-01T00:05:00Z" {
		t.Fatalf("source_open=%s, want lagged 00:05 source", feature.SourceOpenTime)
	}
	if feature.MarkClose != 100.20 || feature.BasisBPS < 19.9 || feature.BasisBPS > 20.1 {
		t.Fatalf("feature used wrong source row: %+v", feature)
	}
	if feature.UsesFutureRows || feature.ForwardLabelsAsFeatureInput {
		t.Fatalf("feature leaked future/label input: %+v", feature)
	}
	if strings.Contains(feature.ContextStateID, RangeContextTriageLabelContainedRotation) {
		t.Fatalf("context state leaked label: %s", feature.ContextStateID)
	}
}

func TestDerivativesContextCohortRequiresDerivativeImprovement(t *testing.T) {
	cfg := DefaultFuturesDerivativesContextAuditConfig()
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
	cfg.MaxBucketShareOfLocalBaselineFull = 0.75
	cfg.MaxBucketShareOfLocalBaselineSplit = 0.75
	labels := []FuturesDerivativesContextAuditLabelRow{}
	for _, split := range rangeDiscoveryPeriodSplits(DefaultSplits()) {
		for i := 0; i < 40; i++ {
			labels = append(labels, derivativesContextTestLabel(split.Name, "basis_premium_wide", "basis_expanding", RangeContextTriageLabelContainedRotation))
		}
		for i := 0; i < 40; i++ {
			labels = append(labels, derivativesContextTestLabel(split.Name, "basis_discount_wide", "basis_contracting", RangeContextTriageLabelBoundaryChop))
		}
	}
	cohorts := derivativesContextCohortRows(labels, cfg, DefaultSplits(), true)
	rankings := derivativesContextRankingRows(cohorts)
	passing := false
	for _, row := range rankings {
		if row.PassesGate && row.DerivativesBucketID == "basis=basis_premium_wide" && row.RouteCandidate == RangeStateConstructionLoopRouteRotation {
			passing = true
		}
		if row.PassesGate && strings.Contains(row.DerivativesBucketID, "basis_discount_wide") && row.RouteCandidate == RangeStateConstructionLoopRouteRotation {
			t.Fatalf("bad basis bucket should not pass rotation: %+v", row)
		}
	}
	if !passing {
		t.Fatalf("expected passing derivatives-improved rotation cohort, top=%+v", rankings[:minInt(3, len(rankings))])
	}
	result := FuturesDerivativesContextAuditResult{
		CoverageRows:   []FuturesDerivativesContextAuditCoverageRow{{Scope: "lagged_basis_context", CoverageGatePass: true}},
		LocalStateRows: []FuturesDerivativesContextAuditLocalStateRow{{ClosedCandleOnly: true}},
		RankingRows:    rankings,
	}
	if got := FuturesDerivativesContextAuditStopState(result); got != DerivativesContextAuditStopStatePassedNeedsSpec {
		t.Fatalf("stop=%s, want passed needs spec", got)
	}
	result.LabelRows = []FuturesDerivativesContextAuditLabelRow{{ForwardLabelUsedAsFeature: true}}
	if got := FuturesDerivativesContextAuditStopState(result); got != DerivativesContextAuditStopStateRejectedFutureLabelLeak {
		t.Fatalf("stop=%s, want future leak rejection", got)
	}
}

func TestDerivativesContextStopStateRejectsClosedFamilyRescue(t *testing.T) {
	result := FuturesDerivativesContextAuditResult{
		CoverageRows:   []FuturesDerivativesContextAuditCoverageRow{{Scope: "lagged_basis_context", CoverageGatePass: true}},
		LocalStateRows: []FuturesDerivativesContextAuditLocalStateRow{{ClosedFamilyRescue: true}},
	}
	if got := FuturesDerivativesContextAuditStopState(result); got != DerivativesContextAuditStopStateRejectedClosedFamilyRescue {
		t.Fatalf("stop=%s, want closed-family rescue rejection", got)
	}
	cfg := DefaultFuturesDerivativesContextAuditConfig()
	cfg.RescueClosedFamily = true
	res, err := RunFuturesDerivativesContextAudit(cfg, DefaultSplits())
	if err != nil {
		t.Fatal(err)
	}
	if res.StopState != DerivativesContextAuditStopStateRejectedClosedFamilyRescue {
		t.Fatalf("run stop=%s, want closed-family rescue rejection", res.StopState)
	}
}

func derivativesContextTestLabel(split, basisLevel, basisChange, forwardLabel string) FuturesDerivativesContextAuditLabelRow {
	row := FuturesDerivativesContextAuditLabelRow{
		Symbol:                    RangeUniverseSymbolBTCUSDT,
		Timestamp:                 time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Format(timeLayout),
		Timeframe:                 RangeDiscoveryTimeframe1h,
		Split:                     split,
		HorizonBars:               2,
		LocalRangeBucketID:        "1h::geometry_balanced_orderly::vol_normal::trend_flat::impulse_none",
		BasisLevelBucket:          basisLevel,
		BasisChangeBucket:         basisChange,
		BasisChangePresent:        true,
		PremiumLevelBucket:        "premium_flat",
		PremiumPresent:            true,
		ContextStateID:            "derivatives_context_v1::BTCUSDT::1h::local::basis=" + basisLevel,
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
