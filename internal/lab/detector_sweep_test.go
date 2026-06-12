package lab

import "testing"

func TestDefaultDetectorSweepProfiles(t *testing.T) {
	profiles := DefaultDetectorSweepProfiles(20)
	if len(profiles) != 19 {
		t.Fatalf("profiles=%d, want 19", len(profiles))
	}

	ids := map[string]bool{}
	balancedBaselines := 0
	adxComparisons := 0
	adxOffProfiles := 0
	for _, profile := range profiles {
		if ids[profile.ProfileID] {
			t.Fatalf("duplicate profile id %q", profile.ProfileID)
		}
		ids[profile.ProfileID] = true
		if profile.LookbackDays != 20 {
			t.Fatalf("lookback=%d, want 20 for profile %+v", profile.LookbackDays, profile)
		}
		if profile.IsBalancedBaseline {
			balancedBaselines++
			if profile.ProfileID != "p30_c12_bollinger_on_adx_off" ||
				profile.Percentile != 0.30 ||
				profile.MinConsecutiveBars != 12 ||
				!profile.UseBollinger ||
				profile.UseADX {
				t.Fatalf("bad balanced baseline: %+v", profile)
			}
		}
		if profile.IsADXComparison {
			adxComparisons++
			if profile.ProfileID != "p30_c12_bollinger_on_adx_on" ||
				profile.Percentile != 0.30 ||
				profile.MinConsecutiveBars != 12 ||
				!profile.UseBollinger ||
				!profile.UseADX {
				t.Fatalf("bad ADX comparison: %+v", profile)
			}
		}
		if !profile.UseADX {
			adxOffProfiles++
		}
		if profile.IsBalancedBaseline && profile.IsADXComparison {
			t.Fatalf("profile cannot be both baseline and ADX comparison: %+v", profile)
		}
	}

	if balancedBaselines != 1 {
		t.Fatalf("balanced baselines=%d, want 1", balancedBaselines)
	}
	if adxComparisons != 1 {
		t.Fatalf("ADX comparisons=%d, want 1", adxComparisons)
	}
	if adxOffProfiles != 18 {
		t.Fatalf("ADX-off profiles=%d, want 18", adxOffProfiles)
	}
}

func TestRunDetectorSweepIncludesEverySplitForEachProfile(t *testing.T) {
	candles := make([]Candle, 64)
	for i := range candles {
		base := 100 + float64(i%7)
		candles[i] = testCandle(i, base, base+1, base-1, base+0.5)
	}
	splits := []Split{
		{Name: "early", Start: candles[0].CloseTime.Add(-1), End: candles[32].CloseTime},
		{Name: "full_2021_2026"},
	}
	cfg := RangeDetectorConfig{
		ATRPeriod:          2,
		DonchianPeriod:     2,
		BollingerPeriod:    2,
		ADXPeriod:          2,
		LookbackDays:       1,
		BarsPerDay:         1,
		MinConsecutiveBars: 1,
	}

	rows, err := RunDetectorSweep(candles, cfg, splits)
	if err != nil {
		t.Fatalf("RunDetectorSweep error: %v", err)
	}
	if len(rows) != 19*len(splits) {
		t.Fatalf("rows=%d, want %d", len(rows), 19*len(splits))
	}

	byProfile := map[string]map[string]bool{}
	for _, row := range rows {
		if row.ProfileID == "" {
			t.Fatalf("empty profile id in row %+v", row)
		}
		if row.LookbackDays != 1 {
			t.Fatalf("lookback=%d, want 1 in row %+v", row.LookbackDays, row)
		}
		if row.TotalBars == 0 {
			t.Fatalf("total bars should be nonzero in row %+v", row)
		}
		if byProfile[row.ProfileID] == nil {
			byProfile[row.ProfileID] = map[string]bool{}
		}
		byProfile[row.ProfileID][row.Split] = true
	}
	if len(byProfile) != 19 {
		t.Fatalf("profiles in rows=%d, want 19", len(byProfile))
	}
	for profileID, gotSplits := range byProfile {
		for _, split := range splits {
			if !gotSplits[split.Name] {
				t.Fatalf("profile %s missing split %s", profileID, split.Name)
			}
		}
	}
}

func TestDetectorSweepRowsPreserveProfileAndDutyMetrics(t *testing.T) {
	profile := newDetectorSweepProfile(0.30, 12, true, false, 20, false)
	dutyRows := []DetectorDutyCycleRow{{
		Split:                "full_2021_2026",
		ActiveBars:           10,
		TotalBars:            100,
		DutyCycle:            0.10,
		Episodes:             2,
		AvgEpisodeLength:     5,
		MedianEpisodeLength:  4.5,
		LongestEpisodeLength: 6,
	}}

	rows := detectorSweepRows(profile, dutyRows)
	if len(rows) != 1 {
		t.Fatalf("rows=%d, want 1", len(rows))
	}
	row := rows[0]
	if !row.IsBalancedBaseline || row.IsADXComparison {
		t.Fatalf("bad profile markers in row %+v", row)
	}
	if row.ProfileID != profile.ProfileID ||
		row.Percentile != profile.Percentile ||
		row.MinConsecutiveBars != profile.MinConsecutiveBars ||
		row.UseBollinger != profile.UseBollinger ||
		row.UseADX != profile.UseADX ||
		row.LookbackDays != profile.LookbackDays {
		t.Fatalf("profile metadata not preserved: %+v", row)
	}
	if row.ActiveBars != 10 ||
		row.TotalBars != 100 ||
		row.DutyCycle != 0.10 ||
		row.Episodes != 2 ||
		row.AvgEpisodeLength != 5 ||
		row.MedianEpisodeLength != 4.5 ||
		row.LongestEpisodeLength != 6 {
		t.Fatalf("duty metrics not preserved: %+v", row)
	}
}
