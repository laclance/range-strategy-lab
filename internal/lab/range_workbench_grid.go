package lab

import "fmt"

func rangeWorkbenchGrid(cfg RangeOptimizationWorkbenchConfig) []RangeOptimizationWorkbenchTrialSpec {
	rows := []RangeOptimizationWorkbenchTrialSpec{}
	add := func(spec RangeOptimizationWorkbenchTrialSpec) {
		if len(rows) >= cfg.MaxTrials {
			return
		}
		spec.TrialID = fmt.Sprintf("range_workbench_trial_%04d", len(rows)+1)
		spec = rangeWorkbenchTrialSpecDefaults(spec)
		rows = append(rows, spec)
	}

	for _, tf := range []string{"5m", "15m"} {
		for _, lookback := range rangeWorkbenchLookbacks(tf) {
			for _, edge := range []float64{0.10, 0.15} {
				for _, minRangeATR := range []float64{3, 6} {
					for _, stopMode := range []string{"range_edge_plus_0.25_atr", "range_edge_plus_0.50_atr"} {
						base := rangeWorkbenchBaseSpec(tf, lookback, edge, minRangeATR, stopMode)
						add(rangeWorkbenchFamilySpec(base, "rolling_midpoint_reversion", "midpoint_reversion", "midpoint"))
						add(rangeWorkbenchFamilySpec(base, "rolling_vwap_reversion", "vwap_reversion", "vwap"))
						add(rangeWorkbenchFamilySpec(base, "range_edge_exhaustion", "interior_to_edge_exhaustion_fade", "midpoint"))
					}
				}
			}
		}
	}

	for _, edge := range []float64{0.10, 0.15} {
		for _, minRangeATR := range []float64{3, 6} {
			for _, stopMode := range []string{"range_edge_plus_0.25_atr", "range_edge_plus_0.50_atr"} {
				base := rangeWorkbenchBaseSpec("15m", 96, edge, minRangeATR, stopMode)
				base.EntryArchetype = "previous_day_edge_reversion"
				add(rangeWorkbenchFamilySpec(base, "previous_day_range", "previous_day_edge_reversion", "midpoint"))
				add(rangeWorkbenchFamilySpec(base, "previous_day_vwap", "previous_day_edge_reversion", "vwap"))
			}
		}
	}
	return rows
}

func rangeWorkbenchLookbacks(timeframe string) []int {
	if timeframe == "15m" {
		return []int{48, 96}
	}
	return []int{144, 288}
}

func rangeWorkbenchTimeStop(timeframe string) int {
	if timeframe == "15m" {
		return 16
	}
	return 36
}

func rangeWorkbenchBaseSpec(timeframe string, lookback int, edge float64, minRangeATR float64, stopMode string) RangeOptimizationWorkbenchTrialSpec {
	return RangeOptimizationWorkbenchTrialSpec{
		Timeframe:            timeframe,
		RangeLookback:        lookback,
		EdgeZonePct:          edge,
		InteriorLowPct:       0.40,
		InteriorHighPct:      0.60,
		VWAPDistanceRangePct: 0.15,
		ProgressATRMultiple:  0.35,
		MinRangeATR:          minRangeATR,
		StopMode:             stopMode,
		TimeStopBars:         rangeWorkbenchTimeStop(timeframe),
	}
}

func rangeWorkbenchFamilySpec(base RangeOptimizationWorkbenchTrialSpec, family string, archetype string, target string) RangeOptimizationWorkbenchTrialSpec {
	base.FamilyID = family
	base.EntryArchetype = archetype
	base.TargetMode = target
	return base
}

func rangeWorkbenchTrialSpecDefaults(spec RangeOptimizationWorkbenchTrialSpec) RangeOptimizationWorkbenchTrialSpec {
	if spec.ATRPercentileMin == "" {
		spec.ATRPercentileMin = "none"
	}
	if spec.ATRPercentileMax == "" {
		spec.ATRPercentileMax = "none"
	}
	if spec.TrendMode == "" {
		spec.TrendMode = "none"
	}
	if spec.ImpulseATRMin == "" {
		spec.ImpulseATRMin = "none"
	}
	if spec.ImpulseATRMax == "" {
		spec.ImpulseATRMax = "none"
	}
	if spec.VolumeMode == "" {
		spec.VolumeMode = "none"
	}
	if spec.VolumeLookback == 0 {
		spec.VolumeLookback = 96
	}
	return spec
}
