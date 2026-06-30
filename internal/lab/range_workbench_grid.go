package lab

import "fmt"

func rangeWorkbenchGrid(cfg RangeOptimizationWorkbenchConfig) []RangeOptimizationWorkbenchTrialSpec {
	rows := []RangeOptimizationWorkbenchTrialSpec{}
	add := func(spec RangeOptimizationWorkbenchTrialSpec) {
		if len(rows) >= cfg.MaxTrials {
			return
		}
		spec.TrialID = fmt.Sprintf("range_workbench_trial_%04d", len(rows)+1)
		if spec.ATRPercentileMin == "" { spec.ATRPercentileMin = "none" }
		if spec.ATRPercentileMax == "" { spec.ATRPercentileMax = "none" }
		if spec.TrendMode == "" { spec.TrendMode = "none" }
		if spec.ImpulseATRMin == "" { spec.ImpulseATRMin = "none" }
		if spec.ImpulseATRMax == "" { spec.ImpulseATRMax = "none" }
		if spec.VolumeMode == "" { spec.VolumeMode = "none" }
		if spec.VolumeLookback == 0 { spec.VolumeLookback = 96 }
		rows = append(rows, spec)
	}

	for _, tf := range []string{"5m", "15m"} {
		lookbacks := []int{144, 288}
		timeStop := 36
		if tf == "15m" {
			lookbacks = []int{48, 96}
			timeStop = 16
		}
		for _, lookback := range lookbacks {
			for _, edge := range []float64{0.10, 0.15} {
				for _, minRangeATR := range []float64{3, 6} {
					for _, stopMode := range []string{"range_edge_plus_0.25_atr", "range_edge_plus_0.50_atr"} {
						base := RangeOptimizationWorkbenchTrialSpec{Timeframe: tf, RangeLookback: lookback, EdgeZonePct: edge, InteriorLowPct: 0.40, InteriorHighPct: 0.60, VWAPDistanceRangePct: 0.15, ProgressATRMultiple: 0.35, MinRangeATR: minRangeATR, StopMode: stopMode, TimeStopBars: timeStop}
						mid := base
						mid.FamilyID = "rolling_midpoint_reversion"
						mid.EntryArchetype = "midpoint_reversion"
						mid.TargetMode = "midpoint"
						add(mid)
						vwap := base
						vwap.FamilyID = "rolling_vwap_reversion"
						vwap.EntryArchetype = "vwap_reversion"
						vwap.TargetMode = "vwap"
						add(vwap)
						edgeFade := base
						edgeFade.FamilyID = "range_edge_exhaustion"
						edgeFade.EntryArchetype = "interior_to_edge_exhaustion_fade"
						edgeFade.TargetMode = "midpoint"
						add(edgeFade)
					}
				}
			}
		}
	}

	for _, edge := range []float64{0.10, 0.15} {
		for _, minRangeATR := range []float64{3, 6} {
			for _, stopMode := range []string{"range_edge_plus_0.25_atr", "range_edge_plus_0.50_atr"} {
				base := RangeOptimizationWorkbenchTrialSpec{FamilyID: "previous_day_range", Timeframe: "15m", EntryArchetype: "previous_day_edge_reversion", RangeLookback: 96, EdgeZonePct: edge, InteriorLowPct: 0.40, InteriorHighPct: 0.60, VWAPDistanceRangePct: 0.15, ProgressATRMultiple: 0.35, MinRangeATR: minRangeATR, TargetMode: "midpoint", StopMode: stopMode, TimeStopBars: 16}
				add(base)
				vwap := base
				vwap.FamilyID = "previous_day_vwap"
				vwap.TargetMode = "vwap"
				add(vwap)
			}
		}
	}
	return rows
}
