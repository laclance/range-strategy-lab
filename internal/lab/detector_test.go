package lab

import (
	"math"
	"testing"
	"time"
)

func TestRollingPriorPercentileExcludesCurrentBar(t *testing.T) {
	got := rollingPriorPercentile([]float64{1, 100, 2}, 2, 0.5)
	if !math.IsNaN(got[0]) || !math.IsNaN(got[1]) {
		t.Fatalf("warmup thresholds=%v, want NaN", got[:2])
	}
	if got[2] != 1 {
		t.Fatalf("threshold=%v, want prior-window median 1", got[2])
	}
}

func TestMarkConsecutiveActiveDoesNotBackfill(t *testing.T) {
	got := markConsecutiveActive([]bool{true, true, false, true, true, true}, 3)
	want := []bool{false, false, false, false, false, true}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("active[%d]=%v, want %v; full=%v", i, got[i], want[i], got)
		}
	}
}

func TestSummarizeDetectorSplitsEpisodeMetrics(t *testing.T) {
	candles := make([]Candle, 7)
	for i := range candles {
		candles[i] = testCandle(i, 100, 101, 99, 100)
	}
	classifications := testClassifications(false, true, true, false, true, true, true)

	rows, episodes := SummarizeDetectorSplits(candles, classifications, []Split{{Name: "full_2021_2026"}})
	if len(rows) != 1 {
		t.Fatalf("rows=%d, want 1", len(rows))
	}
	row := rows[0]
	if row.ActiveBars != 5 || row.TotalBars != 7 || row.Episodes != 2 {
		t.Fatalf("bad row counts: %+v", row)
	}
	if !almostEqual(row.DutyCycle, 5.0/7.0) ||
		!almostEqual(row.AvgEpisodeLength, 2.5) ||
		!almostEqual(row.MedianEpisodeLength, 2.5) ||
		row.LongestEpisodeLength != 3 {
		t.Fatalf("bad row metrics: %+v", row)
	}
	if len(episodes) != 2 || episodes[0].LengthBars != 2 || episodes[1].LengthBars != 3 {
		t.Fatalf("bad episodes: %+v", episodes)
	}
}

func TestSummarizeDetectorSplitsClipsEpisodesToSplitBoundaries(t *testing.T) {
	candles := []Candle{
		dayCandle(2021, 1, 1),
		dayCandle(2021, 1, 2),
		dayCandle(2021, 1, 3),
		dayCandle(2021, 1, 4),
	}
	classifications := testClassifications(true, true, true, true)
	split := Split{
		Name:  "middle",
		Start: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC),
	}

	rows, episodes := SummarizeDetectorSplits(candles, classifications, []Split{split})
	if rows[0].ActiveBars != 2 || rows[0].TotalBars != 2 || rows[0].Episodes != 1 {
		t.Fatalf("bad clipped row: %+v", rows[0])
	}
	if len(episodes) != 1 {
		t.Fatalf("episodes=%d, want 1", len(episodes))
	}
	episode := episodes[0]
	if episode.StartIndex != 1 || episode.EndIndex != 2 || episode.LengthBars != 2 {
		t.Fatalf("bad clipped episode: %+v", episode)
	}
}

func TestSummarizeDetectorSplitsHandlesZeroEpisodes(t *testing.T) {
	candles := make([]Candle, 3)
	for i := range candles {
		candles[i] = testCandle(i, 100, 101, 99, 100)
	}
	rows, episodes := SummarizeDetectorSplits(candles, testClassifications(false, false, false), []Split{{Name: "full_2021_2026"}})
	if len(episodes) != 0 {
		t.Fatalf("episodes=%d, want 0", len(episodes))
	}
	row := rows[0]
	if row.ActiveBars != 0 || row.Episodes != 0 || row.AvgEpisodeLength != 0 || row.MedianEpisodeLength != 0 || row.LongestEpisodeLength != 0 {
		t.Fatalf("bad zero-episode row: %+v", row)
	}
}

func testClassifications(active ...bool) []RangeClassification {
	out := make([]RangeClassification, len(active))
	for i, isActive := range active {
		out[i] = RangeClassification{Index: i, Active: isActive}
	}
	return out
}

func dayCandle(year int, month time.Month, day int) Candle {
	t := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return Candle{
		OpenTime:  t.Add(-5 * time.Minute),
		CloseTime: t,
		Open:      100,
		High:      101,
		Low:       99,
		Close:     100,
		Volume:    1,
	}
}
