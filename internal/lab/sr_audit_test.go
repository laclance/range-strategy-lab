package lab

import (
	"reflect"
	"testing"
	"time"

	sr "github.com/laclance/go-sr"
)

func TestToSRCandlesPreservesOHLCVTimestamps(t *testing.T) {
	openTime := time.Date(2024, 4, 1, 12, 0, 0, 0, time.UTC)
	candles := []Candle{{
		OpenTime:  openTime,
		CloseTime: openTime.Add(5*time.Minute - time.Millisecond),
		Open:      100,
		High:      103,
		Low:       98,
		Close:     101.5,
		Volume:    42,
	}}

	got := toSRCandles(candles)
	if len(got) != 1 {
		t.Fatalf("converted candles=%d, want 1", len(got))
	}
	if !got[0].OpenTime.Equal(candles[0].OpenTime) ||
		!got[0].CloseTime.Equal(candles[0].CloseTime) ||
		got[0].Open != candles[0].Open ||
		got[0].High != candles[0].High ||
		got[0].Low != candles[0].Low ||
		got[0].Close != candles[0].Close ||
		got[0].Volume != candles[0].Volume {
		t.Fatalf("converted candle mismatch: got %+v want %+v", got[0], candles[0])
	}
}

func TestRunSRAuditStartsAfterWarmupAndEmitsDeterministicMetadata(t *testing.T) {
	candles := make([]Candle, 25)
	for i := range candles {
		base := 100 + float64(i%5)
		candles[i] = testCandle(i, base, base+1.5, base-1.5, base+0.25)
	}
	cfg := SRAuditConfig{
		Timeframe:    "5m",
		Mode:         string(sr.ModeZones),
		LookbackBars: 3,
		MinStrength:  1,
	}

	rows, err := RunSRAudit(candles, cfg, DefaultSplits())
	if err != nil {
		t.Fatalf("RunSRAudit error: %v", err)
	}
	again, err := RunSRAudit(candles, cfg, DefaultSplits())
	if err != nil {
		t.Fatalf("second RunSRAudit error: %v", err)
	}
	if !reflect.DeepEqual(rows, again) {
		t.Fatalf("SR audit rows are not deterministic")
	}

	warmup := sr.WarmupCandles(cfg.LookbackBars, sr.ModeZones)
	if len(rows) != len(candles)-warmup {
		t.Fatalf("rows=%d, want %d", len(rows), len(candles)-warmup)
	}
	first := rows[0]
	if first.Index != warmup || first.WarmupBars != warmup {
		t.Fatalf("first row index/warmup = %d/%d, want %d/%d", first.Index, first.WarmupBars, warmup, warmup)
	}
	if first.Timeframe != "5m" || first.Mode != string(sr.ModeZones) || first.LookbackBars != 3 || first.MinStrength != 1 {
		t.Fatalf("bad SR metadata: %+v", first)
	}
	if first.DetectorProfileID != BalancedDetectorProfileID {
		t.Fatalf("detector profile=%q, want %q", first.DetectorProfileID, BalancedDetectorProfileID)
	}
	if first.Split != "2021_2022_stress" {
		t.Fatalf("split=%q, want 2021_2022_stress", first.Split)
	}
	if first.NearestSupportSourcePivots == nil || first.NearestResistanceSourcePivots == nil {
		t.Fatalf("source pivot slices should be empty slices, not nil: %+v", first)
	}
}

func TestSRAuditNearestMetadataMirrorsGoSRSideRules(t *testing.T) {
	levels := sr.Levels{
		Timeframe:      "5m",
		NearSupport:    true,
		NearResistance: true,
		Levels: []sr.Level{
			{
				Price:              100,
				Top:                101,
				Bottom:             99,
				Strength:           2,
				Score:              1.25,
				IsHigh:             false,
				LastTouchIndex:     7,
				SourcePivotIndexes: []int{3, 7},
			},
			{
				Price:              100,
				Top:                101.5,
				Bottom:             98.5,
				Strength:           3,
				Score:              2.5,
				IsHigh:             true,
				LastTouchIndex:     9,
				SourcePivotIndexes: []int{5, 9},
			},
		},
		RawZones: []sr.Level{{Price: 100}, {Price: 100}},
	}
	c := testCandle(0, 100, 102, 98, 100)
	row := newSRAuditRow(c, 12, SRAuditConfig{
		Timeframe:         "5m",
		Mode:              string(sr.ModeZones),
		LookbackBars:      120,
		MinStrength:       2,
		DetectorProfileID: BalancedDetectorProfileID,
	}, 138, levels, RangeClassification{RawActive: true, Active: true}, []Split{{Name: "full_2021_2026"}})

	if !row.HasSupport || !row.HasResistance || !row.NearSupport || !row.NearResistance {
		t.Fatalf("expected both near support and near resistance: %+v", row)
	}
	if row.NearestSupport != 100 ||
		row.NearestSupportDistance != 0 ||
		row.NearestSupportStrength != 2 ||
		row.NearestSupportScore != 1.25 ||
		row.NearestSupportTop != 101 ||
		row.NearestSupportBottom != 99 ||
		row.NearestSupportLastTouchIndex != 7 ||
		!reflect.DeepEqual(row.NearestSupportSourcePivots, []int{3, 7}) {
		t.Fatalf("bad support metadata: %+v", row)
	}
	if row.NearestResistance != 100 ||
		row.NearestResistanceDistance != 0 ||
		row.NearestResistanceStrength != 3 ||
		row.NearestResistanceScore != 2.5 ||
		row.NearestResistanceTop != 101.5 ||
		row.NearestResistanceBottom != 98.5 ||
		row.NearestResistanceLastTouchIndex != 9 ||
		!reflect.DeepEqual(row.NearestResistanceSourcePivots, []int{5, 9}) {
		t.Fatalf("bad resistance metadata: %+v", row)
	}
	if row.QualifiedZoneCount != 2 || row.RawZoneCount != 2 {
		t.Fatalf("bad zone counts: %+v", row)
	}
}

func TestRunSRAuditRejectsInvalidConfig(t *testing.T) {
	candles := []Candle{testCandle(0, 100, 101, 99, 100)}

	badMode := DefaultSRAuditConfig()
	badMode.Mode = string(sr.ModeLegacy)
	if _, err := RunSRAudit(candles, badMode, nil); err == nil {
		t.Fatalf("expected invalid mode error")
	}

	badLookback := DefaultSRAuditConfig()
	badLookback.LookbackBars = -1
	if _, err := RunSRAudit(candles, badLookback, nil); err == nil {
		t.Fatalf("expected invalid lookback error")
	}

	badTimeframe := DefaultSRAuditConfig()
	badTimeframe.Timeframe = "15m"
	if _, err := RunSRAudit(candles, badTimeframe, nil); err == nil {
		t.Fatalf("expected invalid timeframe error")
	}
}
