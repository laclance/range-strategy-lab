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

func TestRollingPriorPercentileRejectsInvalidInputs(t *testing.T) {
	for _, got := range [][]float64{
		rollingPriorPercentile([]float64{1, 2}, 0, 0.5),
		rollingPriorPercentile([]float64{1, 2}, 1, 0),
		rollingPriorPercentile([]float64{1, 2}, 1, 1),
	} {
		for i, value := range got {
			if !math.IsNaN(value) {
				t.Fatalf("threshold[%d]=%v, want NaN for invalid input", i, value)
			}
		}
	}
}

func TestPercentileFromSortedBounds(t *testing.T) {
	if got := percentileFromSorted([]float64{}, 0.5); !math.IsNaN(got) {
		t.Fatalf("empty percentile=%v, want NaN", got)
	}
	if got := percentileFromSorted([]float64{10, 20, 30}, -0.5); got != 10 {
		t.Fatalf("low percentile=%v, want first value", got)
	}
	if got := percentileFromSorted([]float64{10, 20, 30}, 2); got != 30 {
		t.Fatalf("high percentile=%v, want last value", got)
	}
}

func TestInsertAndRemoveSorted(t *testing.T) {
	values := []float64{1, 3}
	values = insertSorted(values, 2)
	want := []float64{1, 2, 3}
	for i := range want {
		if values[i] != want[i] {
			t.Fatalf("inserted values=%v, want %v", values, want)
		}
	}
	values = removeSorted(values, 9)
	if len(values) != 3 {
		t.Fatalf("remove missing changed values=%v", values)
	}
	values = removeSorted(values, 2)
	want = []float64{1, 3}
	for i := range want {
		if values[i] != want[i] {
			t.Fatalf("removed values=%v, want %v", values, want)
		}
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

func TestCompressionRangeDetectorDefaultsNameAndValidation(t *testing.T) {
	cfg := RangeDetectorConfig{}.withDefaults()
	defaults := DefaultCompressionRangeDetectorConfig()
	if cfg != defaults {
		t.Fatalf("defaults=%+v, want %+v", cfg, defaults)
	}

	detector := CompressionRangeDetector{}
	if detector.Name() != "compression_range" {
		t.Fatalf("detector name=%q", detector.Name())
	}

	badConfigs := []RangeDetectorConfig{
		{ATRPeriod: -1},
		{DonchianPeriod: -1},
		{BollingerPeriod: -1},
		{ADXPeriod: -1},
		{LookbackDays: -1},
		{BarsPerDay: -1},
		{Percentile: -0.1},
		{Percentile: 1},
		{MinConsecutiveBars: -1},
	}
	for _, bad := range badConfigs {
		cfg := bad.withDefaults()
		if err := cfg.validate(); err == nil {
			t.Fatalf("expected validation error for %+v", bad)
		}
	}
}

func TestCompressionRangeDetectorClassifyCoversOptionalMetrics(t *testing.T) {
	candles := make([]Candle, 40)
	for i := range candles {
		base := 100 + float64(i%4)
		candles[i] = testCandle(i, base, base+1, base-1, base+0.5)
	}
	cfg := RangeDetectorConfig{
		ATRPeriod:          2,
		DonchianPeriod:     2,
		BollingerPeriod:    2,
		ADXPeriod:          2,
		LookbackDays:       1,
		BarsPerDay:         2,
		Percentile:         0.5,
		MinConsecutiveBars: 1,
		UseBollinger:       false,
		UseADX:             true,
	}

	rows, err := (CompressionRangeDetector{Config: cfg}).Classify(candles)
	if err != nil {
		t.Fatalf("Classify error: %v", err)
	}
	if len(rows) != len(candles) {
		t.Fatalf("rows=%d, want %d", len(rows), len(candles))
	}
	for i, row := range rows {
		if row.Index != i || row.CloseTime == "" {
			t.Fatalf("bad classification row %d: %+v", i, row)
		}
	}

	cfg.ATRPeriod = -1
	if _, err := (CompressionRangeDetector{Config: cfg}).Classify(candles); err == nil {
		t.Fatalf("expected invalid detector config error")
	}
}

func TestCompressionRangeDetectorClassifyEvaluatesBollingerFilter(t *testing.T) {
	candles := make([]Candle, 40)
	for i := range candles {
		candles[i] = testCandle(i, 100, 101, 99, 100)
	}
	cfg := RangeDetectorConfig{
		ATRPeriod:          2,
		DonchianPeriod:     2,
		BollingerPeriod:    2,
		ADXPeriod:          2,
		LookbackDays:       1,
		BarsPerDay:         2,
		Percentile:         0.5,
		MinConsecutiveBars: 1,
		UseBollinger:       true,
		UseADX:             false,
	}
	rows, err := (CompressionRangeDetector{Config: cfg}).Classify(candles)
	if err != nil {
		t.Fatalf("Classify error: %v", err)
	}
	active := 0
	for _, row := range rows {
		if row.Active {
			active++
		}
	}
	if active == 0 {
		t.Fatalf("expected at least one active row with Bollinger filter")
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

func TestSummarizeDetectorSplitsClipsLongClassificationInput(t *testing.T) {
	candles := []Candle{testCandle(0, 100, 101, 99, 100)}
	classifications := testClassifications(true, true)
	rows, episodes := SummarizeDetectorSplits(candles, classifications, []Split{{Name: "full_2021_2026"}})
	if len(rows) != 1 || rows[0].ActiveBars != 1 || len(episodes) != 1 {
		t.Fatalf("unexpected clipped summary rows=%+v episodes=%+v", rows, episodes)
	}
}

func TestMedianIntHandlesEmptyOddAndEvenInputs(t *testing.T) {
	if got := medianInt(nil); got != 0 {
		t.Fatalf("empty median=%v, want 0", got)
	}
	if got := medianInt([]int{5, 1, 3}); got != 3 {
		t.Fatalf("odd median=%v, want 3", got)
	}
	if got := medianInt([]int{10, 2}); got != 6 {
		t.Fatalf("even median=%v, want 6", got)
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
