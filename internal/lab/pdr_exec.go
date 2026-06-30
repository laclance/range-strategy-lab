package lab

import "strings"

func (s *BTC15MPreviousDayRangeReversionStrategy) markExecuted(trades []Trade) {
	executed := map[string]bool{}
	for _, tr := range trades {
		executed[tr.Signal] = true
	}
	for i := range s.rows {
		s.rows[i].Executed = executed[s.rows[i].SignalID]
	}
}

func (s *BTC15MPreviousDayRangeReversionStrategy) skipRows() []BTC15MPreviousDayRangeReversionSkipRow {
	rows := make([]BTC15MPreviousDayRangeReversionSkipRow, 0, len(s.skips))
	for key, count := range s.skips {
		parts := strings.SplitN(key, "|", 2)
		rows = append(rows, BTC15MPreviousDayRangeReversionSkipRow{CandidateID: BTC15MPreviousDayRangeReversionCandidateID, Split: parts[0], Reason: parts[1], Count: count, MissingDataPolicy: "skip_explicit_missing_or_invalid_rows_no_fill"})
	}
	return rows
}
