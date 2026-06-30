package lab

import (
	"math"
	"testing"
	"time"
)

func TestBTCPostCompressionRangeFactsUsePriorWindowsAndReferences(t *testing.T) {
	candles := []Candle{
		postCompression15MTestCandle(0, 100, 110, 90, 100, 10),
		postCompression15MTestCandle(1, 100, 108, 92, 100, 10),
		postCompression15MTestCandle(2, 100, 500, 1, 100, 10),
		postCompression15MTestCandle(3, 100, 104, 96, 100, 10),
		postCompression15MTestCandle(4, 100, 900, 1, 100, 10),
	}
	facts := btcPostCompressionBuildRangeFacts(candles, 2, []float64{0.5}, 2)

	if !postCompressionAlmostEqual(facts.high[2], 110) || !postCompressionAlmostEqual(facts.low[2], 90) {
		t.Fatalf("decision index 2 should use only prior bars [0,1], got high=%v low=%v", facts.high[2], facts.low[2])
	}
	if !postCompressionAlmostEqual(facts.widthPct[2], 0.20) {
		t.Fatalf("width at index 2=%v, want 0.20", facts.widthPct[2])
	}
	if !postCompressionAlmostEqual(facts.thresholds[0.5][4], facts.widthPct[2]) {
		t.Fatalf("threshold at index 4 should use previous valid widths and exclude current: got %v want %v", facts.thresholds[0.5][4], facts.widthPct[2])
	}
}

func TestBTCPostCompressionForwardLabelUsesNextOpenAndSideSymmetry(t *testing.T) {
	candles := []Candle{
		postCompression15MTestCandle(0, 100, 101, 99, 100, 10),
		postCompression15MTestCandle(1, 100, 101, 99, 100, 10),
		postCompression15MTestCandle(2, 100, 110, 95, 100, 10),
		postCompression15MTestCandle(3, 100, 108, 90, 105, 10),
	}

	longLabel, ok := btcPostCompressionForwardLabelFor(candles, 1, BTC15MPostCompressionSideLong, 2)
	if !ok {
		t.Fatalf("expected long label")
	}
	shortLabel, ok := btcPostCompressionForwardLabelFor(candles, 1, BTC15MPostCompressionSideShort, 2)
	if !ok {
		t.Fatalf("expected short label")
	}
	if longLabel.anchorOpenTime != candles[2].OpenTime.UTC().Format(timeLayout) || !postCompressionAlmostEqual(longLabel.anchorOpen, 100) {
		t.Fatalf("label should anchor at open[d+1], got %+v", longLabel)
	}
	if !postCompressionAlmostEqual(longLabel.forwardCloseReturnBP, 500) || !postCompressionAlmostEqual(shortLabel.forwardCloseReturnBP, -500) {
		t.Fatalf("bad side-adjusted close returns: long=%v short=%v", longLabel.forwardCloseReturnBP, shortLabel.forwardCloseReturnBP)
	}
	if !postCompressionAlmostEqual(longLabel.favorableBP, 1000) || !postCompressionAlmostEqual(longLabel.adverseBP, 1000) {
		t.Fatalf("bad long favorable/adverse excursions: %+v", longLabel)
	}
	if !postCompressionAlmostEqual(shortLabel.favorableBP, 1000) || !postCompressionAlmostEqual(shortLabel.adverseBP, 1000) {
		t.Fatalf("bad short favorable/adverse excursions: %+v", shortLabel)
	}
	if _, ok := btcPostCompressionForwardLabelFor(candles, 2, BTC15MPostCompressionSideLong, 2); ok {
		t.Fatalf("expected missing future label rejection")
	}
}

func TestBTCPostCompressionSyntheticRunDedupsAndFindsAdjacentPass(t *testing.T) {
	candles := []Candle{
		postCompression15MTestCandle(0, 100, 110, 90, 100, 10),
		postCompression15MTestCandle(1, 100, 110, 90, 100, 10),
		postCompression15MTestCandle(2, 100, 100.5, 99.5, 100, 10),
		postCompression15MTestCandle(3, 100, 100.5, 99.5, 100, 10),
		postCompression15MTestCandle(4, 100, 102.2, 100, 102, 12),
		postCompression15MTestCandle(5, 102, 110, 102, 109, 12),
		postCompression15MTestCandle(6, 109, 119, 109, 109.5, 10),
		postCompression15MTestCandle(7, 109.5, 110, 99, 100, 10),
		postCompression15MTestCandle(8, 100, 101, 99, 100, 10),
		postCompression15MTestCandle(9, 100, 101, 99, 100, 10),
	}
	cfg := DefaultFuturesBTC15MPostCompressionDirectionalExpansionAuditConfig()
	cfg.CompressionLookbacks = []int{2}
	cfg.CompressionPercentiles = []float64{0.8}
	cfg.PercentileReferenceBars = 2
	cfg.BreakoutATRMultiples = []float64{0.1, 0.2}
	cfg.VolumeModes = []string{BTC15MPostCompressionVolumeNone}
	cfg.VolumeLookbackBars = 2
	cfg.ATRPeriod = 2
	cfg.HorizonsBars = []int{1}
	cfg.MinFullDedupCandidates = 1
	cfg.MinSplitDedupCandidates = 1
	splits := []Split{{Name: fullSplitName}}
	result := btcPostCompressionRunFrom15M(candles, []BTC15MPostCompressionSourceRow{{SourceFactsPass: true, ValidationStatus: "accepted"}}, []BTC15MPostCompressionCoverageRow{{SourceResamplePass: true, ValidationStatus: "accepted", Complete: true}}, cfg, splits)

	if len(result.ParameterCells) != 2 {
		t.Fatalf("parameter cells=%d, want 2", len(result.ParameterCells))
	}
	for _, cell := range result.ParameterCells {
		if cell.DedupEventRows != 2 {
			t.Fatalf("cell %s dedup events=%d, want 2", cell.CellID, cell.DedupEventRows)
		}
	}
	if len(result.DedupEvents) != 2 {
		t.Fatalf("dedup events=%d, want 2: %+v", len(result.DedupEvents), result.DedupEvents)
	}
	for _, event := range result.DedupEvents {
		if event.MatchedCellCount != 2 {
			t.Fatalf("dedup event matched cells=%d, want 2: %+v", event.MatchedCellCount, event)
		}
	}
	if len(result.CandidateRows) != 4 {
		t.Fatalf("candidate rows=%d, want 4", len(result.CandidateRows))
	}
	if result.Falsification.StopState != BTC15MPostCompressionDirectionalExpansionStopStatePassedNeedsReview {
		t.Fatalf("stop state=%s failures=%v", result.Falsification.StopState, result.Falsification.FailureReasons)
	}
	if result.Falsification.AdjacentPassingCellSideHorizons != 2 {
		t.Fatalf("adjacent passing cells=%d, want 2", result.Falsification.AdjacentPassingCellSideHorizons)
	}
}

func TestBTCPostCompressionFalsificationStopStatePrecedence(t *testing.T) {
	cfg := DefaultFuturesBTC15MPostCompressionDirectionalExpansionAuditConfig()
	passing := BTC15MPostCompressionFalsification{
		SourceResamplePass:               true,
		LeakagePass:                      true,
		CandidateSizePass:                true,
		SplitSizePass:                    true,
		BaselineSeparationPass:           true,
		AdjacentCellClusterPass:          true,
		SplitStabilityPass:               true,
		ClosedFamilyProtectionPass:       true,
		DerivativesVetoContaminationPass: true,
		CommonOutputsZeroTrade:           true,
	}
	if got := btcPostCompressionFalsification(passing, cfg).StopState; got != BTC15MPostCompressionDirectionalExpansionStopStatePassedNeedsReview {
		t.Fatalf("passing stop state=%s", got)
	}
	passing.CandidateSizePass = false
	if got := btcPostCompressionFalsification(passing, cfg).StopState; got != BTC15MPostCompressionDirectionalExpansionStopStateFailedNoPremise {
		t.Fatalf("candidate-size failure stop state=%s", got)
	}
	passing.ClosedFamilyProtectionPass = false
	if got := btcPostCompressionFalsification(passing, cfg).StopState; got != BTC15MPostCompressionDirectionalExpansionStopStateClosedReslice {
		t.Fatalf("closed-family stop state=%s", got)
	}
}

func postCompression15MTestCandle(i int, open, high, low, close, volume float64) Candle {
	openTime := time.Date(2024, 1, 1, 0, i*15, 0, 0, time.UTC)
	return Candle{
		OpenTime:  openTime,
		CloseTime: openTime.Add(15*time.Minute - time.Millisecond),
		Open:      open,
		High:      high,
		Low:       low,
		Close:     close,
		Volume:    volume,
	}
}

func postCompressionAlmostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}
