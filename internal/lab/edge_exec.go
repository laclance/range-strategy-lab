package lab

import "strings"

func (s *BTC15MRangeEdgeExhaustionFadeStrategy) markExecuted(trades []Trade) {
	executed := map[string]bool{}
	for _, tr := range trades {
		executed[tr.Signal] = true
	}
	for i := range s.rows {
		s.rows[i].Executed = executed[s.rows[i].SignalID]
	}
}

func (s *BTC15MRangeEdgeExhaustionFadeStrategy) skipRows() []BTC15MRangeEdgeExhaustionFadeSkipRow {
	rows := make([]BTC15MRangeEdgeExhaustionFadeSkipRow, 0, len(s.skips))
	for key, count := range s.skips {
		parts := strings.SplitN(key, "|", 2)
		rows = append(rows, BTC15MRangeEdgeExhaustionFadeSkipRow{CandidateID: BTC15MRangeEdgeExhaustionFadeCandidateID, Split: parts[0], Reason: parts[1], Count: count, MissingDataPolicy: "skip_explicit_missing_or_invalid_rows_no_fill"})
	}
	return rows
}
