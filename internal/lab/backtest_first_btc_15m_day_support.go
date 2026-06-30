package lab

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"
	"time"
)

type btc15MDayRange struct {
	key      string
	high     float64
	low      float64
	mid      float64
	width    float64
	count    int
	complete bool
}

func (cfg BacktestFirstBTC15MPreviousDayRangeReversionConfig) withDefaults() BacktestFirstBTC15MPreviousDayRangeReversionConfig {
	defaults := DefaultBacktestFirstBTC15MPreviousDayRangeReversionConfig()
	if cfg.ApprovedSourcePath == "" { cfg.ApprovedSourcePath = defaults.ApprovedSourcePath }
	if cfg.ExpectedSourceRows == 0 { cfg.ExpectedSourceRows = defaults.ExpectedSourceRows }
	if cfg.ExpectedFirstOpenTime == "" { cfg.ExpectedFirstOpenTime = defaults.ExpectedFirstOpenTime }
	if cfg.ExpectedLastOpenTime == "" { cfg.ExpectedLastOpenTime = defaults.ExpectedLastOpenTime }
	if cfg.ExpectedZeroVolumeCount == 0 { cfg.ExpectedZeroVolumeCount = defaults.ExpectedZeroVolumeCount }
	if cfg.Expected15MRows == 0 { cfg.Expected15MRows = defaults.Expected15MRows }
	if cfg.Expected15MLastOpenTime == "" { cfg.Expected15MLastOpenTime = defaults.Expected15MLastOpenTime }
	if cfg.ATRPeriod == 0 { cfg.ATRPeriod = defaults.ATRPeriod }
	if cfg.OuterZonePct == 0 { cfg.OuterZonePct = defaults.OuterZonePct }
	if cfg.StopATRMultiple == 0 { cfg.StopATRMultiple = defaults.StopATRMultiple }
	if cfg.MaxHoldBars == 0 { cfg.MaxHoldBars = defaults.MaxHoldBars }
	if cfg.MinFullTrades == 0 { cfg.MinFullTrades = defaults.MinFullTrades }
	if cfg.MinSplitTrades == 0 { cfg.MinSplitTrades = defaults.MinSplitTrades }
	if cfg.FullMaxDrawdownLimit == 0 { cfg.FullMaxDrawdownLimit = defaults.FullMaxDrawdownLimit }
	if cfg.SplitMaxDrawdownLimit == 0 { cfg.SplitMaxDrawdownLimit = defaults.SplitMaxDrawdownLimit }
	return cfg
}

func btc15MPrevDaySourceRow(manifest SourceManifest, cfg BacktestFirstBTC15MPreviousDayRangeReversionConfig) BTC15MPreviousDayRangeReversionSourceRow {
	row := BTC15MPreviousDayRangeReversionSourceRow{BacktestName: BacktestFirstBTC15MPreviousDayRangeReversionName, CandidateID: BTC15MPreviousDayRangeReversionCandidateID, Path: manifest.Path, ApprovedPath: cfg.ApprovedSourcePath, Product: manifest.Product, Symbol: manifest.Symbol, Interval: manifest.Interval, RowCount: manifest.RowCount, ExpectedRowCount: cfg.ExpectedSourceRows, FirstOpenTime: manifest.FirstOpenTime, ExpectedFirstOpenTime: cfg.ExpectedFirstOpenTime, LastOpenTime: manifest.LastOpenTime, ExpectedLastOpenTime: cfg.ExpectedLastOpenTime, GapCount: manifest.GapCount, DuplicateCount: manifest.DuplicateCount, ZeroVolumeCount: manifest.ZeroVolumeCount, ExpectedZeroVolumeCount: cfg.ExpectedZeroVolumeCount, ComparisonOnly: manifest.ComparisonOnly, ClosedCandleOnly: true, ValidationStatus: "accepted"}
	failures := []string{}
	if manifest.ValidationStatus != "accepted" || manifest.Product != "Binance USDT-M futures" || manifest.ComparisonOnly || manifest.Symbol != "BTCUSDT" || manifest.Interval != "5m" || filepath.Clean(manifest.Path) != filepath.Clean(cfg.ApprovedSourcePath) || manifest.RowCount != cfg.ExpectedSourceRows || manifest.FirstOpenTime != cfg.ExpectedFirstOpenTime || manifest.LastOpenTime != cfg.ExpectedLastOpenTime || manifest.GapCount != 0 || manifest.DuplicateCount != 0 || manifest.ZeroVolumeCount != cfg.ExpectedZeroVolumeCount {
		failures = append(failures, "source contract mismatch")
	}
	row.SourceFactsPass = len(failures) == 0
	if !row.SourceFactsPass { row.ValidationStatus = "rejected"; row.ValidationError = strings.Join(failures, "; ") }
	return row
}

func btc15MPrevDayResample(candles []Candle, cfg BacktestFirstBTC15MPreviousDayRangeReversionConfig, sourcePass bool) ([]Candle, BTC15MPreviousDayRangeReversionCoverageRow) {
	out := []Candle{}
	missingBuckets := 0
	for i := 0; i+2 < len(candles); i += 3 {
		first, second, third := candles[i], candles[i+1], candles[i+2]
		start := first.OpenTime.UTC()
		if start.Minute()%15 != 0 || !second.OpenTime.UTC().Equal(start.Add(5*time.Minute)) || !third.OpenTime.UTC().Equal(start.Add(10*time.Minute)) {
			missingBuckets++
			continue
		}
		high := math.Max(first.High, math.Max(second.High, third.High))
		low := math.Min(first.Low, math.Min(second.Low, third.Low))
		out = append(out, Candle{OpenTime: start, CloseTime: third.CloseTime.UTC(), Open: first.Open, High: high, Low: low, Close: third.Close, Volume: first.Volume + second.Volume + third.Volume})
	}
	row := BTC15MPreviousDayRangeReversionCoverageRow{BacktestName: BacktestFirstBTC15MPreviousDayRangeReversionName, CandidateID: BTC15MPreviousDayRangeReversionCandidateID, Timeframe: "15m", RowCount: len(out), ExpectedRowCount: cfg.Expected15MRows, ExpectedLastOpenTime: cfg.Expected15MLastOpenTime, ExpectedChildBars: 3, MissingChildBuckets: missingBuckets, ClosedCandleOnly: true, ValidationStatus: "accepted"}
	if len(out) > 0 { row.FirstOpenTime = out[0].OpenTime.UTC().Format(timeLayout); row.LastOpenTime = out[len(out)-1].OpenTime.UTC().Format(timeLayout) }
	row.SourceResamplePass = sourcePass && missingBuckets == 0 && row.RowCount == cfg.Expected15MRows && row.LastOpenTime == cfg.Expected15MLastOpenTime
	if !row.SourceResamplePass { row.ValidationStatus = "rejected"; row.ValidationError = fmt.Sprintf("15m resample mismatch rows=%d expected=%d missing_child_buckets=%d last_open=%s expected=%s", row.RowCount, cfg.Expected15MRows, missingBuckets, row.LastOpenTime, cfg.Expected15MLastOpenTime) }
	return out, row
}

func btc15MPrevDayRanges(candles []Candle) map[string]btc15MDayRange {
	ranges := map[string]btc15MDayRange{}
	for _, c := range candles {
		key := c.OpenTime.UTC().Format("2006-01-02")
		r, ok := ranges[key]
		if !ok { r = btc15MDayRange{key: key, high: c.High, low: c.Low} }
		r.high = math.Max(r.high, c.High)
		r.low = math.Min(r.low, c.Low)
		r.count++
		ranges[key] = r
	}
	for key, r := range ranges {
		r.width = r.high - r.low
		r.mid = (r.high + r.low) / 2
		r.complete = r.count == 96 && r.width > 0
		ranges[key] = r
	}
	return ranges
}

func btc15MPrevDayTradeRows(trades []Trade, splits []Split) []BTC15MPreviousDayRangeReversionTradeRow {
	rows := make([]BTC15MPreviousDayRangeReversionTradeRow, 0, len(trades))
	for _, tr := range trades {
		exitTime, _ := time.Parse(timeLayout, tr.ExitTime)
		rows = append(rows, BTC15MPreviousDayRangeReversionTradeRow{SignalID: tr.Signal, CloseSplit: btc15MPrevDaySplit(exitTime, splits), Side: tr.Side, EntryTime: tr.EntryTime, ExitTime: tr.ExitTime, GrossPnL: tr.GrossPnL, NetPnL: tr.NetPnL, Fees: tr.Fees, Slippage: tr.Slippage, ExitReason: tr.Reason, HoldBars: tr.HoldBars})
	}
	return rows
}

func btc15MPrevDayEvaluate(source BTC15MPreviousDayRangeReversionSourceRow, coverage BTC15MPreviousDayRangeReversionCoverageRow, rows []SummaryRow, cfg BacktestFirstBTC15MPreviousDayRangeReversionConfig) BTC15MPreviousDayRangeReversionFalsification {
	byKey := map[string]SummaryRow{}
	for _, row := range rows { byKey[row.Split+"|"+row.Side] = row }
	full, oos, recent := byKey["full_2021_2026|all"], byKey["2023_2024_oos|all"], byKey["2025_2026_recent|all"]
	stress := byKey["2021_2022_stress|all"]
	minPrimary := minInt(stress.TotalTrades, minInt(oos.TotalTrades, recent.TotalTrades))
	dominantShare := 0.0
	if full.TotalTrades > 0 { dominantShare = float64(maxInt(stress.TotalTrades, maxInt(oos.TotalTrades, recent.TotalTrades))) / float64(full.TotalTrades) }
	drawdownPass := full.MaxDrawdown <= cfg.FullMaxDrawdownLimit && stress.MaxDrawdown <= cfg.SplitMaxDrawdownLimit && oos.MaxDrawdown <= cfg.SplitMaxDrawdownLimit && recent.MaxDrawdown <= cfg.SplitMaxDrawdownLimit
	tradeCountPass := full.TotalTrades >= cfg.MinFullTrades && minPrimary >= cfg.MinSplitTrades
	return BTC15MPreviousDayRangeReversionFalsification{SourceResamplePass: source.SourceFactsPass && coverage.SourceResamplePass, LeakagePass: true, TradeCountPass: tradeCountPass, GrossEdgePass: full.GrossPnL > 0 && oos.GrossPnL >= 0 && recent.GrossPnL >= 0, NetEdgePass: full.NetPnL > 0, DrawdownPass: drawdownPass, RobustnessPass: tradeCountPass && dominantShare <= 0.90, SideReportingPass: true, FullExecutedTrades: full.TotalTrades, RequiredFullExecutedTrades: cfg.MinFullTrades, MinimumPrimarySplitExecutedTrades: minPrimary, RequiredPrimarySplitTrades: cfg.MinSplitTrades, FullGrossPnL: full.GrossPnL, FullNetPnL: full.NetPnL, FullProfitFactor: full.ProfitFactor, FullMaxDrawdown: full.MaxDrawdown, DominantPrimarySplitTradeShare: dominantShare}
}

func btc15MPrevDayFalsification(report BTC15MPreviousDayRangeReversionFalsification, cfg BacktestFirstBTC15MPreviousDayRangeReversionConfig) BTC15MPreviousDayRangeReversionFalsification {
	report.BacktestName = BacktestFirstBTC15MPreviousDayRangeReversionName
	report.CandidateID = BTC15MPreviousDayRangeReversionCandidateID
	if !report.SourceResamplePass || !report.LeakagePass { report.StopState = BTC15MPreviousDayRangeReversionStopStateFailedSourceOrResample } else if report.TradeCountPass && report.GrossEdgePass && report.NetEdgePass && report.DrawdownPass && report.RobustnessPass { report.StopState = BTC15MPreviousDayRangeReversionStopStatePassedNeedsReview } else { report.StopState = BTC15MPreviousDayRangeReversionStopStateFailedNoUsableStrategy }
	if !report.TradeCountPass { report.FailureReasons = append(report.FailureReasons, "insufficient_executed_trades") }
	if !report.GrossEdgePass { report.FailureReasons = append(report.FailureReasons, "gross_edge_gate_failed") }
	if !report.NetEdgePass { report.FailureReasons = append(report.FailureReasons, "net_edge_gate_failed") }
	if !report.DrawdownPass { report.FailureReasons = append(report.FailureReasons, "drawdown_gate_failed") }
	if !report.RobustnessPass { report.FailureReasons = append(report.FailureReasons, "robustness_gate_failed") }
	_ = cfg
	return report
}

func btc15MPrevDaySplit(t time.Time, splits []Split) string {
	for _, split := range splits { if split.Contains(t) { return split.Name } }
	return "unassigned"
}
