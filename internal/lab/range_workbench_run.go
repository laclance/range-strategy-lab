package lab

import "fmt"

func RunRangeOptimizationWorkbench(candles []Candle, manifest SourceManifest, cfg RangeOptimizationWorkbenchConfig, btCfg BacktestConfig, splits []Split, runID string) (RangeOptimizationWorkbenchResult, error) {
	cfg = cfg.withDefaults()
	if cfg.MaxTrials > 2500 {
		return RangeOptimizationWorkbenchResult{}, fmt.Errorf("range optimization workbench max trials %d exceeds 2500", cfg.MaxTrials)
	}
	if runID == "" {
		return RangeOptimizationWorkbenchResult{}, fmt.Errorf("range optimization workbench requires non-empty run id")
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	source := rangeWorkbenchSourceRow(manifest, cfg, runID)
	resampled15M, coverage := rangeWorkbenchResample15M(candles, cfg, source.SourceFactsPass, runID)
	result := RangeOptimizationWorkbenchResult{RunID: runID, SourceRows: []RangeOptimizationWorkbenchSourceRow{source}, CoverageRows: []RangeOptimizationWorkbenchCoverageRow{coverage}}
	grid := rangeWorkbenchGrid(cfg)
	result.GridRows = grid
	if !source.SourceFactsPass || !coverage.SourceResamplePass {
		result.Falsification = RangeOptimizationWorkbenchFalsification{BacktestName: RangeOptimizationWorkbenchName, RunID: runID, StopState: RangeOptimizationWorkbenchStopStateFailedSourceOrResample, SourceResamplePass: false, TotalTrials: len(grid), MaxTrials: cfg.MaxTrials, FailureReasons: []string{"source_or_resample_gate_failed"}}
		result.StopState = result.Falsification.StopState
		result.RobustnessSummary = RangeOptimizationWorkbenchRobustness{RunID: runID, TotalTrials: len(grid), MaxTrials: cfg.MaxTrials, StopState: result.StopState}
		return result, nil
	}
	for _, spec := range grid {
		trialCandles := candles
		if spec.Timeframe == "15m" {
			trialCandles = resampled15M
		}
		trialCfg := btCfg
		trialCfg.MaxHoldBars = spec.TimeStopBars
		strategy := NewRangeOptimizationWorkbenchTrialStrategy(trialCandles, spec, splits)
		bt := RunBacktest(trialCandles, strategy, trialCfg)
		summaries := SummarizeSplits(bt.Trades, btCfg.StartBalance, splits)
		trialRow, trialSummaryRows, _ := rangeWorkbenchBuildTrialRow(runID, manifest.Path, spec, summaries)
		result.TrialRows = append(result.TrialRows, trialRow)
		result.TrialSummaryRows = append(result.TrialSummaryRows, trialSummaryRows...)
	}
	candidates, rejected := rangeWorkbenchCandidates(runID, result.TrialRows)
	result.TopCandidates = candidates
	result.RejectedCandidates = rejected
	stopState := RangeOptimizationWorkbenchStopStateFailedNoCandidate
	selectedTrialID := ""
	selectedCandidateID := ""
	if len(candidates) > 0 {
		stopState = RangeOptimizationWorkbenchStopStateCandidateSelectedNeedsValidation
		selectedTrialID = candidates[0].TrialID
		selectedCandidateID = candidates[0].CandidateID
		rangeWorkbenchMarkSelection(result.TrialRows, selectedTrialID)
	}
	result.StopState = stopState
	result.RobustnessSummary = RangeOptimizationWorkbenchRobustness{RunID: runID, TotalTrials: len(result.TrialRows), MaxTrials: cfg.MaxTrials, PassingCandidates: len(candidates), RejectedCandidates: len(rejected), SelectedTrialID: selectedTrialID, SelectedCandidateID: selectedCandidateID, StopState: stopState}
	result.Falsification = RangeOptimizationWorkbenchFalsification{BacktestName: RangeOptimizationWorkbenchName, RunID: runID, StopState: stopState, SourceResamplePass: true, TotalTrials: len(result.TrialRows), MaxTrials: cfg.MaxTrials, PassingCandidates: len(candidates), SelectedTrialID: selectedTrialID, SelectedCandidateID: selectedCandidateID}
	if len(candidates) == 0 {
		result.Falsification.FailureReasons = []string{"no_candidate_passed_minimum_workbench_filters"}
	}
	return result, nil
}
