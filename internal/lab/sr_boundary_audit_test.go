package lab

import (
	"math"
	"testing"
)

func TestSRBoundarySupportEventForwardMetrics(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 98, 99.2),
		testCandle(2, 99.2, 104, 96, 97.5),
		testCandle(3, 97.5, 105, 97, 100.5),
	}
	row := SRAuditRow{
		Index:                     0,
		OpenTime:                  candles[0].OpenTime.Format(timeLayout),
		CloseTime:                 candles[0].CloseTime.Format(timeLayout),
		Split:                     "2021_2022_stress",
		Close:                     100,
		DetectorProfileID:         BalancedDetectorProfileID,
		DetectorRawActive:         true,
		DetectorActive:            true,
		HasSupport:                true,
		NearSupport:               true,
		NearestSupport:            100,
		NearestSupportTop:         101,
		NearestSupportBottom:      98,
		NearestSupportDistancePct: 0.0004,
		NearestSupportStrength:    2,
		NearestSupportScore:       3.5,
	}

	events, summary, err := RunSRBoundaryAudit(candles, []SRAuditRow{row}, SRBoundaryAuditConfig{HorizonsBars: []int{3}, DetectorActiveOnly: true})
	if err != nil {
		t.Fatalf("RunSRBoundaryAudit error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("events=%d, want 1", len(events))
	}
	if len(summary) != 1 {
		t.Fatalf("summary rows=%d, want 1", len(summary))
	}

	event := events[0]
	if event.Side != SRBoundarySideSupport || event.HorizonBars != 3 {
		t.Fatalf("bad event side/horizon: %+v", event)
	}
	if event.FavorableMove != 5 || event.AdverseMove != 4 ||
		!boundaryAlmostEqual(event.FavorableMovePct, 0.05) || !boundaryAlmostEqual(event.AdverseMovePct, 0.04) {
		t.Fatalf("bad support moves: %+v", event)
	}
	if !event.WickBreak || !event.CloseBreak || !event.ReclaimedAfterBreak || event.Rejected {
		t.Fatalf("bad support break/reclaim flags: %+v", event)
	}
	if !event.FavorableGreaterThanAdverse {
		t.Fatalf("expected favorable greater than adverse: %+v", event)
	}
	if event.StrengthBucket != "2" || event.DistanceBucket != "0_5bp" {
		t.Fatalf("bad buckets: %+v", event)
	}
}

func TestSRBoundaryResistanceEventForwardMetrics(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 104, 94, 103),
		testCandle(2, 103, 103, 93, 99.5),
	}
	row := SRAuditRow{
		Index:                        0,
		OpenTime:                     candles[0].OpenTime.Format(timeLayout),
		CloseTime:                    candles[0].CloseTime.Format(timeLayout),
		Split:                        "2021_2022_stress",
		Close:                        100,
		DetectorProfileID:            BalancedDetectorProfileID,
		DetectorRawActive:            true,
		DetectorActive:               true,
		HasResistance:                true,
		NearResistance:               true,
		NearestResistance:            100,
		NearestResistanceTop:         102,
		NearestResistanceBottom:      99,
		NearestResistanceDistancePct: 0.0008,
		NearestResistanceStrength:    3,
		NearestResistanceScore:       4.5,
	}

	events, _, err := RunSRBoundaryAudit(candles, []SRAuditRow{row}, SRBoundaryAuditConfig{HorizonsBars: []int{2}, DetectorActiveOnly: true})
	if err != nil {
		t.Fatalf("RunSRBoundaryAudit error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("events=%d, want 1", len(events))
	}

	event := events[0]
	if event.Side != SRBoundarySideResistance || event.HorizonBars != 2 {
		t.Fatalf("bad event side/horizon: %+v", event)
	}
	if event.FavorableMove != 7 || event.AdverseMove != 4 ||
		!boundaryAlmostEqual(event.FavorableMovePct, 0.07) || !boundaryAlmostEqual(event.AdverseMovePct, 0.04) {
		t.Fatalf("bad resistance moves: %+v", event)
	}
	if !event.WickBreak || !event.CloseBreak || !event.ReclaimedAfterBreak || event.Rejected {
		t.Fatalf("bad resistance break/reclaim flags: %+v", event)
	}
	if event.StrengthBucket != "3" || event.DistanceBucket != "5_10bp" {
		t.Fatalf("bad buckets: %+v", event)
	}
}

func TestSRBoundaryAuditFiltersDetectorActiveAndSkipsMissingFuture(t *testing.T) {
	candles := []Candle{
		testCandle(0, 100, 101, 99, 100),
		testCandle(1, 100, 102, 99, 101),
		testCandle(2, 101, 103, 100, 102),
	}
	inactive := SRAuditRow{
		Index:                     0,
		OpenTime:                  candles[0].OpenTime.Format(timeLayout),
		CloseTime:                 candles[0].CloseTime.Format(timeLayout),
		Split:                     "2021_2022_stress",
		Close:                     100,
		DetectorActive:            false,
		HasSupport:                true,
		NearSupport:               true,
		NearestSupport:            100,
		NearestSupportTop:         100,
		NearestSupportBottom:      100,
		NearestSupportDistancePct: 0.0001,
		NearestSupportStrength:    2,
	}
	active := inactive
	active.Index = 1
	active.OpenTime = candles[1].OpenTime.Format(timeLayout)
	active.CloseTime = candles[1].CloseTime.Format(timeLayout)
	active.Close = 101
	active.DetectorActive = true
	active.DetectorRawActive = true

	events, _, err := RunSRBoundaryAudit(candles, []SRAuditRow{inactive, active}, SRBoundaryAuditConfig{HorizonsBars: []int{1, 2}, DetectorActiveOnly: true})
	if err != nil {
		t.Fatalf("RunSRBoundaryAudit error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("events=%d, want 1", len(events))
	}
	if events[0].Index != 1 || events[0].HorizonBars != 1 {
		t.Fatalf("unexpected event after filter/skip: %+v", events[0])
	}
	if !boundaryAlmostEqual(events[0].RejectionThreshold, 0.101) {
		t.Fatalf("zero-width rejection threshold=%v, want 0.101", events[0].RejectionThreshold)
	}
}

func TestSummarizeSRBoundaryQualityGroupsAndRates(t *testing.T) {
	events := []SRBoundaryEventRow{
		{
			Split:                       "2021_2022_stress",
			Side:                        SRBoundarySideSupport,
			HorizonBars:                 3,
			StrengthBucket:              "2",
			DistanceBucket:              "0_5bp",
			Score:                       2,
			DistancePct:                 0.0004,
			FavorableMovePct:            0.02,
			AdverseMovePct:              0.01,
			CloseBreak:                  true,
			WickBreak:                   true,
			ReclaimedAfterBreak:         true,
			FavorableGreaterThanAdverse: true,
		},
		{
			Split:            "2021_2022_stress",
			Side:             SRBoundarySideSupport,
			HorizonBars:      3,
			StrengthBucket:   "2",
			DistanceBucket:   "0_5bp",
			Score:            4,
			DistancePct:      0.0002,
			FavorableMovePct: 0.04,
			AdverseMovePct:   0.03,
			Rejected:         true,
		},
	}

	rows := SummarizeSRBoundaryQuality(events)
	if len(rows) != 1 {
		t.Fatalf("rows=%d, want 1", len(rows))
	}
	row := rows[0]
	if row.EventCount != 2 || !boundaryAlmostEqual(row.AvgScore, 3) || !boundaryAlmostEqual(row.AvgDistancePct, 0.0003) {
		t.Fatalf("bad averages: %+v", row)
	}
	if !boundaryAlmostEqual(row.AvgFavorablePct, 0.03) || !boundaryAlmostEqual(row.MedianFavorablePct, 0.03) ||
		!boundaryAlmostEqual(row.AvgAdversePct, 0.02) || !boundaryAlmostEqual(row.MedianAdversePct, 0.02) {
		t.Fatalf("bad move pct aggregates: %+v", row)
	}
	if !boundaryAlmostEqual(row.CloseBreakRate, 0.5) || !boundaryAlmostEqual(row.WickBreakRate, 0.5) ||
		!boundaryAlmostEqual(row.ReclaimAfterBreakRate, 1) || !boundaryAlmostEqual(row.RejectionRate, 0.5) ||
		!boundaryAlmostEqual(row.FavorableGreaterThanAdverseRate, 0.5) {
		t.Fatalf("bad rates: %+v", row)
	}
}

func boundaryAlmostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-12
}
