package lab

import "testing"

func TestDerivativesNoTradeFilterPremiseEvaluationDeduplicatesNestedTrendDown(t *testing.T) {
	cfg := DefaultFuturesDerivativesNoTradeFilterPremiseAuditConfig()
	cfg.RequireExactMetricReproduction = false
	cfg.MinFullCandidateRows = 1
	cfg.MinSplitCandidateRows = 1
	cfg.MinFullToxicRate = 0.50
	cfg.MinWorstSplitToxicRate = 0.50
	cfg.MinFullToxicImprovement = -1

	defs := derivativesNoTradeFilterPremiseDefinitions(cfg)
	downLocal := defs[0].LocalRangeBucketID
	labels := []derivativesNoTradeFilterPremiseLabelWithFeature{}
	id := 0
	for _, split := range rangeDiscoveryPeriodSplits(DefaultSplits()) {
		id++
		labels = append(labels, premiseTestLabel(id, split.Name, downLocal, "basis_discount_small", "basis_change_flat", "premium_discount_small", true))
		id++
		labels = append(labels, premiseTestLabel(id, split.Name, downLocal, "basis_discount_small", "basis_change_flat", "premium_flat", false))
	}

	result := FuturesDerivativesNoTradeFilterPremiseAuditResult{FilterDefinitionRows: defs}
	derivativesNoTradeFilterPremiseEvaluate(&result, labels, cfg, DefaultSplits())

	full := premiseFindUnion(t, result.CanonicalUnionRows, fullSplitName)
	if full.SumExactCandidateRows != 9 || full.DeduplicatedRows != 6 || full.OverlapRows != 3 {
		t.Fatalf("bad union accounting: %+v", full)
	}
	if full.NestedTrendDownPremiumOverlapRows != 3 {
		t.Fatalf("nested trend-down overlap=%d want 3", full.NestedTrendDownPremiumOverlapRows)
	}
	if !full.DoubleCountingProtectionPass {
		t.Fatalf("expected double-counting protection to pass: %+v", full)
	}

	overlap := premiseFindOverlap(t, result.OverlapRows, fullSplitName, "exact_1_trend_down_basis_discount_premium_discount", "exact_2_trend_down_basis_discount")
	if overlap.OverlapRows != 3 || !overlap.NestedOverlap || overlap.DoubleCountedInUnion {
		t.Fatalf("bad nested overlap row: %+v", overlap)
	}
}

func TestDerivativesNoTradeFilterPremiseStopStatesProtectBoundaries(t *testing.T) {
	closedCfg := DefaultFuturesDerivativesNoTradeFilterPremiseAuditConfig()
	closedCfg.RescueClosedFamily = true
	closed, err := RunFuturesDerivativesNoTradeFilterPremiseAudit(closedCfg, DefaultSplits())
	if err != nil {
		t.Fatal(err)
	}
	if closed.StopState != DerivativesNoTradeFilterPremiseStopStateRejectedClosedFamilyRescue {
		t.Fatalf("closed rescue stop=%s", closed.StopState)
	}

	rotationCfg := DefaultFuturesDerivativesNoTradeFilterPremiseAuditConfig()
	rotationCfg.RescueRotationEntry = true
	rotation, err := RunFuturesDerivativesNoTradeFilterPremiseAudit(rotationCfg, DefaultSplits())
	if err != nil {
		t.Fatal(err)
	}
	if rotation.StopState != DerivativesNoTradeFilterPremiseStopStateRejectedRotationEntry {
		t.Fatalf("rotation rescue stop=%s", rotation.StopState)
	}

	result := FuturesDerivativesNoTradeFilterPremiseAuditResult{
		CoverageRows: premiseCoveragePassRows(),
		VetoCandidateRows: []FuturesDerivativesNoTradeFilterPremiseVetoCandidateRow{{
			ForwardLabelUsedAsFeature: true,
		}},
	}
	if got := FuturesDerivativesNoTradeFilterPremiseAuditStopState(result); got != DerivativesNoTradeFilterPremiseStopStateRejectedFutureLabelLeak {
		t.Fatalf("future leak stop=%s", got)
	}
}

func TestDerivativesNoTradeFilterPremiseDefinitionsStayNoTradeOnly(t *testing.T) {
	for _, def := range derivativesNoTradeFilterPremiseDefinitions(DefaultFuturesDerivativesNoTradeFilterPremiseAuditConfig()) {
		if def.RouteCandidate != RangeStateConstructionLoopRouteNoTradeToxic {
			t.Fatalf("definition route=%s, want no-trade toxic", def.RouteCandidate)
		}
		if def.RotationEntryAllowed || def.ClosedFamilyRescueAllowed || def.ForwardLabelsAsFilterInput {
			t.Fatalf("definition breaks premise boundaries: %+v", def)
		}
	}
}

func TestDerivativesNoTradeFilterPremiseKeepsFullLocalStateConstructionContext(t *testing.T) {
	cfg := DefaultFuturesDerivativesNoTradeFilterPremiseAuditConfig()
	if !stringSliceContains(cfg.StateConfig.Timeframes, RangeDiscoveryTimeframe15m) ||
		!stringSliceContains(cfg.StateConfig.Timeframes, RangeDiscoveryTimeframe1h) ||
		!stringSliceContains(cfg.StateConfig.Timeframes, RangeDiscoveryTimeframe4h) {
		t.Fatalf("state construction timeframes=%v, want 15m/1h/4h for exact local bucket reproduction", cfg.StateConfig.Timeframes)
	}
	if !intSliceContains(cfg.StateConfig.Horizons15M, 12) ||
		!intSliceContains(cfg.StateConfig.Horizons15M, 24) ||
		!intSliceContains(cfg.StateConfig.Horizons15M, 48) {
		t.Fatalf("15m horizons=%v, want full context-audit horizons for exact local bucket reproduction", cfg.StateConfig.Horizons15M)
	}
}

func premiseTestLabel(id int, split, local, basis, change, premium string, toxic bool) derivativesNoTradeFilterPremiseLabelWithFeature {
	label := FuturesDerivativesContextAuditLabelRow{
		LabelRowID:               id,
		Symbol:                   RangeUniverseSymbolBTCUSDT,
		LocalStateRowID:          id,
		FeatureRowID:             id,
		Timestamp:                "2024-01-01T00:00:00Z",
		Timeframe:                RangeDiscoveryTimeframe15m,
		Split:                    split,
		HorizonBars:              48,
		LocalRangeBucketID:       local,
		BasisLevelBucket:         basis,
		BasisChangeBucket:        change,
		BasisChangePresent:       true,
		PremiumLevelBucket:       premium,
		PremiumPresent:           true,
		ContextStateID:           "test",
		ForwardLabelMetadataOnly: true,
	}
	if toxic {
		label.ForwardLabel = RangeContextTriageLabelBoundaryChop
		label.NoTradeToxic = true
	} else {
		label.ForwardLabel = RangeContextTriageLabelContainedRotation
		label.RotationUseful = true
	}
	feature := FuturesDerivativesContextAuditBasisFeatureRow{
		FeatureRowID:       id,
		Symbol:             RangeUniverseSymbolBTCUSDT,
		LocalStateRowID:    id,
		Timeframe:          RangeDiscoveryTimeframe15m,
		Split:              split,
		DecisionCloseTime:  label.Timestamp,
		SourceCloseTime:    "2023-12-31T23:59:59Z",
		BasisBPS:           -2,
		PremiumBPS:         -2,
		ClosedCandleOnly:   true,
		UsesFutureRows:     false,
		BasisLevelBucket:   basis,
		BasisChangeBucket:  change,
		PremiumLevelBucket: premium,
	}
	return derivativesNoTradeFilterPremiseLabelWithFeature{label: label, feature: feature}
}

func premiseCoveragePassRows() []FuturesDerivativesNoTradeFilterPremiseCoverageRow {
	return []FuturesDerivativesNoTradeFilterPremiseCoverageRow{
		{Scope: "source_lag_alignment", CoverageGatePass: true, NoFillPolicy: true},
		{Scope: "local_range_state_resample", CoverageGatePass: true, NoFillPolicy: true},
		{Scope: "lagged_basis_context", CoverageGatePass: true, NoFillPolicy: true},
	}
}

func premiseFindUnion(t *testing.T, rows []FuturesDerivativesNoTradeFilterPremiseCanonicalUnionRow, split string) FuturesDerivativesNoTradeFilterPremiseCanonicalUnionRow {
	t.Helper()
	for _, row := range rows {
		if row.Split == split {
			return row
		}
	}
	t.Fatalf("missing union split %s", split)
	return FuturesDerivativesNoTradeFilterPremiseCanonicalUnionRow{}
}

func premiseFindOverlap(t *testing.T, rows []FuturesDerivativesNoTradeFilterPremiseOverlapRow, split, a, b string) FuturesDerivativesNoTradeFilterPremiseOverlapRow {
	t.Helper()
	for _, row := range rows {
		if row.Split != split {
			continue
		}
		if (row.CandidateID == a && row.OtherCandidateID == b) || (row.CandidateID == b && row.OtherCandidateID == a) {
			return row
		}
	}
	t.Fatalf("missing overlap %s/%s in %s", a, b, split)
	return FuturesDerivativesNoTradeFilterPremiseOverlapRow{}
}

func stringSliceContains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func intSliceContains(values []int, want int) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
