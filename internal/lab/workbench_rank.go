package lab

import "sort"

func rangeWorkbenchCandidates(runID string, trialRows []RangeOptimizationWorkbenchTrialRow) ([]RangeOptimizationWorkbenchCandidate, []RangeOptimizationWorkbenchRejected) {
	passing := []RangeOptimizationWorkbenchTrialRow{}
	rejected := []RangeOptimizationWorkbenchRejected{}
	for _, row := range trialRows {
		if row.FailureReasons == "" {
			passing = append(passing, row)
		} else {
			rejected = append(rejected, RangeOptimizationWorkbenchRejected{RunID: runID, TrialID: row.TrialID, FamilyID: row.FamilyID, Timeframe: row.Timeframe, EntryArchetype: row.EntryArchetype, FailureReasons: row.FailureReasons})
		}
	}
	sort.SliceStable(passing, func(i, j int) bool {
		return passing[i].RobustnessScore > passing[j].RobustnessScore
	})
	candidates := []RangeOptimizationWorkbenchCandidate{}
	for i, row := range passing {
		candidateID := "range_workbench_candidate_" + row.TrialID
		candidates = append(candidates, RangeOptimizationWorkbenchCandidate{RunID: runID, Rank: i + 1, TrialID: row.TrialID, CandidateID: candidateID, FamilyID: row.FamilyID, Timeframe: row.Timeframe, EntryArchetype: row.EntryArchetype, RobustnessScore: row.RobustnessScore, FullTrades: row.FullTrades, FullGrossPnL: row.FullGrossPnL, FullNetPnL: row.FullNetPnL, FullProfitFactor: row.FullProfitFactor, FullMaxDrawdown: row.FullMaxDrawdown, SelectionReason: "passed_minimum_workbench_filters_sorted_by_robustness"})
	}
	return candidates, rejected
}

func rangeWorkbenchMarkSelection(trialRows []RangeOptimizationWorkbenchTrialRow, selectedTrialID string) {
	for i := range trialRows {
		if trialRows[i].TrialID == selectedTrialID {
			trialRows[i].SelectedForLocking = true
		}
	}
}
