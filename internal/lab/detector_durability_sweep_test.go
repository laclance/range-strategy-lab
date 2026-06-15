package lab

import "testing"

func TestRunDetectorDurabilitySweepDefaultsAndRowCounts(t *testing.T) {
	candles := make([]Candle, 80)
	for i := range candles {
		base := 100 + float64(i%5)
		candles[i] = testCandle(i, base, base+1, base-1, base+0.5)
	}

	broadRows, _, stabilityRows, err := RunDetectorDurabilitySweep(candles, RangeDetectorConfig{}, RangeRegimeDurabilityAuditConfig{}, nil)
	if err != nil {
		t.Fatalf("RunDetectorDurabilitySweep error: %v", err)
	}
	profiles := DefaultDetectorSweepProfiles(DefaultCompressionRangeDetectorConfig().LookbackDays)
	wantBroadRows := len(profiles) * len(DefaultSplits()) * len(DefaultRangeRegimeDurabilityAuditConfig().HorizonsBars)
	if len(broadRows) != wantBroadRows {
		t.Fatalf("broad rows=%d, want %d", len(broadRows), wantBroadRows)
	}
	wantStabilityRows := len(profiles) * len(DefaultRangeRegimeDurabilityAuditConfig().HorizonsBars)
	if len(stabilityRows) != wantStabilityRows {
		t.Fatalf("stability rows=%d, want %d", len(stabilityRows), wantStabilityRows)
	}

	validCfg := RangeDetectorConfig{
		ATRPeriod:          2,
		DonchianPeriod:     2,
		BollingerPeriod:    2,
		ADXPeriod:          2,
		LookbackDays:       1,
		BarsPerDay:         1,
		MinConsecutiveBars: 1,
	}
	if _, _, _, err := RunDetectorDurabilitySweep(candles, RangeDetectorConfig{ATRPeriod: -1}, RangeRegimeDurabilityAuditConfig{}, nil); err == nil {
		t.Fatalf("expected invalid detector config error")
	}
	if _, _, _, err := RunDetectorDurabilitySweep(candles, validCfg, RangeRegimeDurabilityAuditConfig{
		HorizonsBars:          []int{0},
		QuickInvalidationBars: 1,
	}, nil); err == nil {
		t.Fatalf("expected invalid durability config error")
	}
}

func TestDetectorDurabilityRowsPreserveProfileAndFullAggregation(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 101),
		testCandle(2, 101, 103, 97, 102),
		testCandle(3, 102, 105, 101, 104),
	}
	classifications := testCompressionClassifications(
		[]bool{true, true, true, false},
		[]bool{false, false, true, false},
	)
	splits := []Split{
		{Name: "2021_2022_stress", Start: candles[0].CloseTime.Add(-1), End: candles[3].CloseTime.Add(1)},
		{Name: fullSplitName},
	}
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	cfg := RangeRegimeDurabilityAuditConfig{
		HorizonsBars:          []int{1},
		QuickInvalidationBars: 1,
		DetectorProfileID:     profile.ProfileID,
	}

	episodeRows := detectorDurabilityEpisodeRows(candles, classifications, nil, cfg, splits, profile.ProfileID)
	if len(episodeRows) != 1 {
		t.Fatalf("episode rows=%d, want 1", len(episodeRows))
	}
	if episodeRows[0].LabelWindowStartIndex != 3 || episodeRows[0].LabelWindowEndIndex != 3 {
		t.Fatalf("label window should start after episode end: %+v", episodeRows[0])
	}
	if !episodeRows[0].LabelQuickInvalidated || !episodeRows[0].LabelInvalidatedUp {
		t.Fatalf("expected future close to quick invalidate upward: %+v", episodeRows[0])
	}

	dutyRows, _ := SummarizeDetectorSplits(candles, classifications, splits)
	broadRows := detectorDurabilitySweepRows(profile, dutyRows, episodeRows, cfg.HorizonsBars, splits)
	if len(broadRows) != 2 {
		t.Fatalf("broad rows=%d, want 2", len(broadRows))
	}
	for _, row := range broadRows {
		if row.ProfileID != profile.ProfileID ||
			!row.IsBalancedBaseline ||
			row.IsADXComparison ||
			row.Percentile != profile.Percentile ||
			row.MinConsecutiveBars != profile.MinConsecutiveBars ||
			!row.UseBollinger ||
			row.UseADX {
			t.Fatalf("profile metadata not preserved: %+v", row)
		}
		if row.DurabilityEpisodeCount != 1 ||
			row.LabelQuickInvalidatedCount != 1 ||
			row.LabelQuickInvalidatedRate != 1 {
			t.Fatalf("bad durability aggregation: %+v", row)
		}
	}
	if broadRows[0].Split != "2021_2022_stress" || broadRows[1].Split != fullSplitName {
		t.Fatalf("bad split order/full aggregation: %+v", broadRows)
	}
}

func TestDetectorDurabilitySliceRowsPreserveSummaryAndProfile(t *testing.T) {
	profile := newDetectorSweepProfile(0.20, 24, false, false, 10, false)
	summaryRows := []RangeRegimeDurabilitySummaryRow{{
		Split:                          fullSplitName,
		HorizonBars:                    12,
		RawLengthBucket:                "48plus",
		ActiveLengthBucket:             "24_47",
		EpisodeWidthBucket:             "gt_50bp",
		WidthToATRBucket:               "gt_4x",
		EpisodeCount:                   7,
		AvgRawLengthBars:               50,
		AvgActiveLengthBars:            30,
		LabelPersistedInsideRangeCount: 2,
		LabelQuickInvalidatedCount:     4,
		LabelPersistedInsideRangeRate:  2.0 / 7.0,
		LabelQuickInvalidatedRate:      4.0 / 7.0,
	}}

	rows := detectorDurabilitySliceRows(profile, summaryRows)
	if len(rows) != 1 {
		t.Fatalf("rows=%d, want 1", len(rows))
	}
	row := rows[0]
	if row.ProfileID != profile.ProfileID ||
		row.Percentile != profile.Percentile ||
		row.MinConsecutiveBars != profile.MinConsecutiveBars ||
		row.UseBollinger != profile.UseBollinger ||
		row.EpisodeCount != 7 ||
		row.RawLengthBucket != "48plus" ||
		row.WidthToATRBucket != "gt_4x" ||
		!boundaryAlmostEqual(row.LabelQuickInvalidatedRate, 4.0/7.0) {
		t.Fatalf("slice row did not preserve profile/summary fields: %+v", row)
	}
}

func TestDetectorDurabilityStabilityRowsCalculateMinMaxDelta(t *testing.T) {
	profile := newDetectorSweepProfile(0.30, 12, true, true, 20, true)
	broadRows := []DetectorDurabilitySweepRow{
		{
			ProfileID:                     profile.ProfileID,
			Split:                         "2021_2022_stress",
			HorizonBars:                   12,
			DutyCycle:                     0.10,
			DurabilityEpisodeCount:        10,
			LabelPersistedInsideRangeRate: 0.20,
			LabelQuickInvalidatedRate:     0.60,
			LabelChoppedRate:              0.10,
			LabelTrendedUpRate:            0.20,
			LabelTrendedDownRate:          0.10,
			LabelReenteredRangeRate:       0.70,
			LabelAvgCloseDriftPct:         -0.01,
		},
		{
			ProfileID:                     profile.ProfileID,
			Split:                         "2023_2024_oos",
			HorizonBars:                   12,
			DutyCycle:                     0.15,
			DurabilityEpisodeCount:        20,
			LabelPersistedInsideRangeRate: 0.30,
			LabelQuickInvalidatedRate:     0.50,
			LabelChoppedRate:              0.20,
			LabelTrendedUpRate:            0.10,
			LabelTrendedDownRate:          0.20,
			LabelReenteredRangeRate:       0.80,
			LabelAvgCloseDriftPct:         0.02,
		},
		{
			ProfileID:                     profile.ProfileID,
			Split:                         "2025_2026_recent",
			HorizonBars:                   12,
			DutyCycle:                     0.12,
			DurabilityEpisodeCount:        15,
			LabelPersistedInsideRangeRate: 0.25,
			LabelQuickInvalidatedRate:     0.70,
			LabelChoppedRate:              0.05,
			LabelTrendedUpRate:            0.30,
			LabelTrendedDownRate:          0.15,
			LabelReenteredRangeRate:       0.60,
			LabelAvgCloseDriftPct:         0.01,
		},
	}

	rows := detectorDurabilityStabilityRows([]DetectorSweepProfile{profile}, broadRows, []int{12}, DefaultSplits())
	if len(rows) != 1 {
		t.Fatalf("rows=%d, want 1", len(rows))
	}
	row := rows[0]
	if !row.IsADXComparison ||
		row.PeriodSplits != 3 ||
		row.PeriodEpisodeCount != 45 ||
		row.EpisodeCountMin != 10 ||
		row.EpisodeCountMax != 20 ||
		row.EpisodeCountDelta != 10 ||
		!boundaryAlmostEqual(row.DutyCycleDelta, 0.05) ||
		!boundaryAlmostEqual(row.LabelPersistedInsideRangeRateMin, 0.20) ||
		!boundaryAlmostEqual(row.LabelPersistedInsideRangeRateMax, 0.30) ||
		!boundaryAlmostEqual(row.LabelQuickInvalidatedRateDelta, 0.20) ||
		!boundaryAlmostEqual(row.LabelTrendedRateMax, 0.45) ||
		!boundaryAlmostEqual(row.LabelAvgCloseDriftPctDelta, 0.03) {
		t.Fatalf("bad stability row: %+v", row)
	}
}
