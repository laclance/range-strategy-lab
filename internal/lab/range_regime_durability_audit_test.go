package lab

import (
	"sort"
	"testing"
	"time"
)

func TestRangeRegimeDurabilityEpisodesUseRawRunThatEventuallyBecomesActive(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 101),
		testCandle(2, 101, 103, 97, 102),
		testCandle(3, 102, 104, 96, 103),
		testCandle(4, 103, 105, 95, 104),
	}
	classifications := testCompressionClassifications(
		[]bool{true, true, false, true, true},
		[]bool{false, true, false, false, false},
	)

	episodes := rangeRegimeDurabilityEpisodes(candles, classifications, nil, []Split{{Name: fullSplitName}}, BalancedDetectorProfileID)
	if len(episodes) != 1 {
		t.Fatalf("episodes=%d, want 1", len(episodes))
	}
	episode := episodes[0]
	if episode.EpisodeID != 1 ||
		episode.StartIndex != 0 ||
		episode.EndIndex != 1 ||
		episode.High != 102 ||
		episode.Low != 98 ||
		episode.RawLengthBars != 2 ||
		episode.ActiveLengthBars != 1 ||
		episode.RawLengthBucket != "lt_12" ||
		episode.ActiveLengthBucket != "lt_12" {
		t.Fatalf("bad episode: %+v", episode)
	}
}

func TestRangeRegimeDurabilityLabelsStartAfterEpisodeEnd(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 101),
		testCandle(2, 101, 103, 97, 102),
		testCandle(3, 102, 200, 50, 102),
		testCandle(4, 102, 103, 100, 101),
	}
	episode := rangeRegimeDurabilityEpisode{
		EpisodeID:          7,
		Split:              fullSplitName,
		StartIndex:         0,
		EndIndex:           2,
		High:               103,
		Low:                97,
		EndClose:           102,
		RawLengthBars:      3,
		ActiveLengthBars:   1,
		RawLengthBucket:    "lt_12",
		ActiveLengthBucket: "lt_12",
		WidthBucket:        "gt_50bp",
		DetectorProfileID:  BalancedDetectorProfileID,
	}

	row, ok := newRangeRegimeDurabilityEpisodeRow(candles, episode, 1, 3)
	if !ok {
		t.Fatalf("expected label row")
	}
	if row.LabelWindowStartIndex != 3 || row.LabelWindowEndIndex != 3 {
		t.Fatalf("bad label window: %+v", row)
	}
	if !row.LabelReenteredRange || !row.LabelPersistedInsideRange ||
		row.LabelInvalidatedUp || row.LabelInvalidatedDown || row.LabelQuickInvalidated {
		t.Fatalf("label used episode candle instead of future close window: %+v", row)
	}
	wantUpMove := (200.0 - 102.0) / 102.0
	wantDownMove := (102.0 - 50.0) / 102.0
	if !boundaryAlmostEqual(row.LabelMaxUpMovePct, wantUpMove) ||
		!boundaryAlmostEqual(row.LabelMaxDownMovePct, wantDownMove) {
		t.Fatalf("label move window mismatch: %+v", row)
	}
}

func TestRunRangeRegimeDurabilityAuditSkipsMissingFutureAndInvalidConfig(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 101),
		testCandle(2, 101, 103, 97, 102),
		testCandle(3, 102, 103, 98, 101),
		testCandle(4, 101, 102, 99, 100),
	}
	classifications := testCompressionClassifications(
		[]bool{true, true, true, false, true},
		[]bool{false, false, true, false, true},
	)

	rows, summary, err := runRangeRegimeDurabilityAuditFromClassifications(candles, classifications, RangeRegimeDurabilityAuditConfig{
		HorizonsBars:          []int{1, 3},
		QuickInvalidationBars: 3,
		DetectorProfileID:     BalancedDetectorProfileID,
	}, []Split{{Name: fullSplitName}})
	if err != nil {
		t.Fatalf("runRangeRegimeDurabilityAuditFromClassifications error: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("rows=%d, want 1 after missing-future skip", len(rows))
	}
	if len(summary) != 1 || summary[0].EpisodeCount != 1 {
		t.Fatalf("bad summary after missing-future skip: %+v", summary)
	}

	defaults := RangeRegimeDurabilityAuditConfig{}.withDefaults()
	if len(defaults.HorizonsBars) != 4 ||
		defaults.QuickInvalidationBars != 3 ||
		defaults.DetectorProfileID != BalancedDetectorProfileID {
		t.Fatalf("bad defaults: %+v", defaults)
	}
	if _, _, err := runRangeRegimeDurabilityAuditFromClassifications(nil, nil, RangeRegimeDurabilityAuditConfig{
		HorizonsBars:          []int{0},
		QuickInvalidationBars: 3,
	}, nil); err == nil {
		t.Fatalf("expected invalid horizon error")
	}
	if _, _, err := runRangeRegimeDurabilityAuditFromClassifications(nil, nil, RangeRegimeDurabilityAuditConfig{
		HorizonsBars:          []int{1},
		QuickInvalidationBars: -1,
	}, nil); err == nil {
		t.Fatalf("expected invalid quick invalidation error")
	}
}

func TestRangeRegimeDurabilitySplitHandlingAndFullSummary(t *testing.T) {
	candles := []Candle{
		testDatedCandle(2022, 12, 31, 23, 50, 100, 101, 99, 100),
		testDatedCandle(2022, 12, 31, 23, 55, 100, 102, 98, 101),
		testDatedCandle(2023, 1, 1, 0, 0, 101, 103, 97, 102),
		testDatedCandle(2023, 1, 1, 0, 5, 102, 103, 100, 101),
	}
	classifications := testCompressionClassifications(
		[]bool{true, true, true, false},
		[]bool{false, false, true, false},
	)

	rows, summary, err := runRangeRegimeDurabilityAuditFromClassifications(candles, classifications, RangeRegimeDurabilityAuditConfig{
		HorizonsBars:          []int{1},
		QuickInvalidationBars: 1,
		DetectorProfileID:     BalancedDetectorProfileID,
	}, DefaultSplits())
	if err != nil {
		t.Fatalf("runRangeRegimeDurabilityAuditFromClassifications error: %v", err)
	}
	if len(rows) != 1 || rows[0].Split != "2023_2024_oos" {
		t.Fatalf("bad episode split rows=%+v", rows)
	}
	foundPeriod := false
	foundFull := false
	for _, row := range summary {
		if row.Split == "2023_2024_oos" && row.EpisodeCount == 1 {
			foundPeriod = true
		}
		if row.Split == fullSplitName && row.EpisodeCount == 1 {
			foundFull = true
		}
	}
	if !foundPeriod || !foundFull {
		t.Fatalf("missing period/full summary rows: %+v", summary)
	}
}

func TestRangeRegimeDurabilitySummaryDenominatorsAndLabels(t *testing.T) {
	rows := []RangeRegimeDurabilityEpisodeRow{
		{
			Split:                     "2021_2022_stress",
			HorizonBars:               3,
			RawLengthBucket:           "12_23",
			ActiveLengthBucket:        "lt_12",
			EpisodeWidthBucket:        "10_25bp",
			WidthToATRBucket:          "1_2x",
			DetectorProfileID:         BalancedDetectorProfileID,
			RawLengthBars:             12,
			ActiveLengthBars:          2,
			EpisodeWidthPct:           0.002,
			AvgNormalizedATR:          0.001,
			EndNormalizedATR:          0.0012,
			WidthToATRRatio:           2,
			LabelReenteredRange:       true,
			LabelPersistedInsideRange: true,
			LabelCloseDriftPct:        0.01,
			LabelMaxUpMovePct:         0.02,
			LabelMaxDownMovePct:       0.005,
		},
		{
			Split:                 "2021_2022_stress",
			HorizonBars:           3,
			RawLengthBucket:       "12_23",
			ActiveLengthBucket:    "lt_12",
			EpisodeWidthBucket:    "10_25bp",
			WidthToATRBucket:      "1_2x",
			DetectorProfileID:     BalancedDetectorProfileID,
			RawLengthBars:         24,
			ActiveLengthBars:      4,
			EpisodeWidthPct:       0.004,
			AvgNormalizedATR:      0.002,
			EndNormalizedATR:      0.0022,
			WidthToATRRatio:       2,
			LabelQuickInvalidated: true,
			LabelInvalidatedUp:    true,
			LabelTrendedUp:        true,
			LabelCloseDriftPct:    0.03,
			LabelMaxUpMovePct:     0.04,
			LabelMaxDownMovePct:   0.006,
		},
	}

	summary := summarizeRangeRegimeDurability(rows)
	var row RangeRegimeDurabilitySummaryRow
	for _, candidate := range summary {
		if candidate.Split == "2021_2022_stress" {
			row = candidate
			break
		}
	}
	if row.EpisodeCount != 2 ||
		row.LabelReenteredRangeCount != 1 ||
		row.LabelPersistedInsideRangeCount != 1 ||
		row.LabelQuickInvalidatedCount != 1 ||
		row.LabelInvalidatedUpCount != 1 ||
		row.LabelTrendedUpCount != 1 ||
		!boundaryAlmostEqual(row.AvgRawLengthBars, 18) ||
		!boundaryAlmostEqual(row.AvgActiveLengthBars, 3) ||
		!boundaryAlmostEqual(row.AvgEpisodeWidthPct, 0.003) ||
		!boundaryAlmostEqual(row.LabelReenteredRangeRate, 0.5) ||
		!boundaryAlmostEqual(row.LabelAvgCloseDriftPct, 0.02) ||
		!boundaryAlmostEqual(row.LabelAvgMaxUpMovePct, 0.03) {
		t.Fatalf("bad summary row: %+v", row)
	}
}

func TestRangeRegimeDurabilityRowsSortDeterministically(t *testing.T) {
	rows := []RangeRegimeDurabilityEpisodeRow{
		{Split: fullSplitName, EpisodeID: 3, HorizonBars: 12},
		{Split: "2021_2022_stress", EpisodeID: 2, HorizonBars: 3},
		{Split: "2021_2022_stress", EpisodeID: 2, HorizonBars: 1},
	}
	sort.Slice(rows, func(i, j int) bool {
		return lessRangeRegimeDurabilityEpisodeRow(rows[i], rows[j])
	})
	if rows[0].Split != "2021_2022_stress" ||
		rows[0].HorizonBars != 1 ||
		rows[1].HorizonBars != 3 ||
		rows[2].Split != fullSplitName {
		t.Fatalf("bad row sort: %+v", rows)
	}

	summary := summarizeRangeRegimeDurability([]RangeRegimeDurabilityEpisodeRow{
		{
			Split:              "2023_2024_oos",
			HorizonBars:        3,
			RawLengthBucket:    "24_47",
			ActiveLengthBucket: "lt_12",
			EpisodeWidthBucket: "10_25bp",
			WidthToATRBucket:   "1_2x",
			DetectorProfileID:  BalancedDetectorProfileID,
		},
		{
			Split:              "2021_2022_stress",
			HorizonBars:        1,
			RawLengthBucket:    "12_23",
			ActiveLengthBucket: "lt_12",
			EpisodeWidthBucket: "10_25bp",
			WidthToATRBucket:   "0_5_1x",
			DetectorProfileID:  BalancedDetectorProfileID,
		},
	})
	if len(summary) != 4 ||
		summary[0].Split != "2021_2022_stress" ||
		summary[1].Split != "2023_2024_oos" ||
		summary[2].Split != fullSplitName ||
		summary[2].HorizonBars != 1 ||
		summary[3].Split != fullSplitName ||
		summary[3].HorizonBars != 3 {
		t.Fatalf("bad summary sort: %+v", summary)
	}
}

func testDatedCandle(year int, month time.Month, day, hour, minute int, open, high, low, close float64) Candle {
	openTime := time.Date(year, month, day, hour, minute, 0, 0, time.UTC)
	return Candle{
		OpenTime:  openTime,
		CloseTime: openTime.Add(5*time.Minute - time.Millisecond),
		Open:      open,
		High:      high,
		Low:       low,
		Close:     close,
		Volume:    1,
	}
}
