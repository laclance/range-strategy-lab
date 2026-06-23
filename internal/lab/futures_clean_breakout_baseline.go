package lab

import (
	"fmt"
	"math"
	"sort"
)

const (
	FuturesCleanBreakoutBaselineName = "futures_clean_breakout_baseline"

	CleanBreakoutCandidate4HUp  = "clean_breakout_4h_up_h12"
	CleanBreakoutCandidate1HAll = "clean_breakout_1h_all_h12"

	CleanBreakoutStopStateSourceOrResampleGap       = "clean_breakout_source_or_resample_gap"
	CleanBreakoutStopStateCodegenOrTestBlocked      = "clean_breakout_codegen_or_test_blocked"
	CleanBreakoutStopStateFailedNoPromotion         = "clean_breakout_baseline_failed_no_promotion"
	CleanBreakoutStopStatePassedNeedsOptimization   = "clean_breakout_baseline_passed_needs_optimization_brief"
	CleanBreakoutStopStateMixedNeedsPortfolioReview = "clean_breakout_baseline_mixed_needs_portfolio_stream_review"
	CleanBreakoutStopStateReviewOnlyNoStrategy      = "clean_breakout_review_only_no_strategy_change"
)

type FuturesCleanBreakoutBaselineConfig struct {
	EventDelayBars             int
	MaxHoldBars                int
	DetectorLookbackDays       int
	DetectorPercentile         float64
	DetectorMinConsecutiveBars int
	MinFullTrades              int
	MinKeySplitTrades          int
	NearFlatNetPct             float64
	Candidates                 []FuturesCleanBreakoutBaselineCandidateConfig
}

type FuturesCleanBreakoutBaselineCandidateConfig struct {
	CandidateID string
	Timeframe   string
	Side        string
	MaxHoldBars int
}

type FuturesCleanBreakoutBaselineResult struct {
	CoverageRows []FuturesRangeDiscoveryCoverageRow
	SignalRows   []FuturesCleanBreakoutBaselineSignalRow
	TradeRows    []FuturesCleanBreakoutBaselineTradeRow
	SummaryRows  []FuturesCleanBreakoutBaselineSummaryRow
	Trades       []Trade
}

type FuturesCleanBreakoutBaselineStrategy struct {
	candidate   FuturesCleanBreakoutBaselineCandidateConfig
	signals     []FuturesCleanBreakoutBaselineSignalRow
	signalsByID map[string]FuturesCleanBreakoutBaselineSignalRow
	byIndex     map[int]Signal
}

type FuturesCleanBreakoutBaselineSignalRow struct {
	SignalID                string    `json:"signal_id"`
	CandidateID             string    `json:"candidate_id"`
	Timeframe               string    `json:"timeframe"`
	Family                  string    `json:"family"`
	SkippedReason           string    `json:"skipped_reason"`
	Split                   string    `json:"split"`
	EpisodeID               int       `json:"episode_id"`
	EpisodeStartIndex       int       `json:"episode_start_index"`
	EpisodeEndIndex         int       `json:"episode_end_index"`
	EpisodeStartTime        string    `json:"episode_start_time"`
	EpisodeEndTime          string    `json:"episode_end_time"`
	EpisodeHigh             float64   `json:"episode_high"`
	EpisodeLow              float64   `json:"episode_low"`
	EpisodeMid              float64   `json:"episode_mid"`
	EpisodeWidth            float64   `json:"episode_width"`
	EpisodeWidthPct         float64   `json:"episode_width_pct"`
	EpisodeRawLengthBars    int       `json:"episode_raw_length_bars"`
	EpisodeActiveLengthBars int       `json:"episode_active_length_bars"`
	BreakoutIndex           int       `json:"breakout_index"`
	BreakoutOpenTime        string    `json:"breakout_open_time"`
	BreakoutCloseTime       string    `json:"breakout_close_time"`
	BreakoutDelayBars       int       `json:"breakout_delay_bars"`
	BreakoutOpen            float64   `json:"breakout_open"`
	BreakoutHigh            float64   `json:"breakout_high"`
	BreakoutLow             float64   `json:"breakout_low"`
	BreakoutClose           float64   `json:"breakout_close"`
	Side                    Direction `json:"side"`
	EntryIndex              int       `json:"entry_index"`
	EntryOpenTime           string    `json:"entry_open_time"`
	EntryOpen               float64   `json:"entry_open"`
	ExpectedEntryPrice      float64   `json:"expected_entry_price"`
	Stop                    float64   `json:"stop"`
	Target                  float64   `json:"target"`
	MaxHoldBars             int       `json:"max_hold_bars"`
	EntryGeometryValid      bool      `json:"entry_geometry_valid"`
}

type FuturesCleanBreakoutBaselineTradeRow struct {
	SignalID          string    `json:"signal_id"`
	CandidateID       string    `json:"candidate_id"`
	Timeframe         string    `json:"timeframe"`
	Family            string    `json:"family"`
	EpisodeID         int       `json:"episode_id"`
	EpisodeStartIndex int       `json:"episode_start_index"`
	EpisodeEndIndex   int       `json:"episode_end_index"`
	EpisodeStartTime  string    `json:"episode_start_time"`
	EpisodeEndTime    string    `json:"episode_end_time"`
	EpisodeHigh       float64   `json:"episode_high"`
	EpisodeLow        float64   `json:"episode_low"`
	EpisodeWidth      float64   `json:"episode_width"`
	BreakoutIndex     int       `json:"breakout_index"`
	BreakoutCloseTime string    `json:"breakout_close_time"`
	BreakoutDelayBars int       `json:"breakout_delay_bars"`
	BreakoutClose     float64   `json:"breakout_close"`
	EntrySplit        string    `json:"entry_split"`
	CloseSplit        string    `json:"close_split"`
	Side              Direction `json:"side"`
	EntryTime         string    `json:"entry_time"`
	ExitTime          string    `json:"exit_time"`
	OpenIndex         int       `json:"open_index"`
	CloseIndex        int       `json:"close_index"`
	EntryPrice        float64   `json:"entry_price"`
	ExitPrice         float64   `json:"exit_price"`
	Stop              float64   `json:"stop"`
	Target            float64   `json:"target"`
	Size              float64   `json:"size"`
	InitialRisk       float64   `json:"initial_risk"`
	GrossPnL          float64   `json:"gross_pnl"`
	NetPnL            float64   `json:"net_pnl"`
	Fees              float64   `json:"fees"`
	Slippage          float64   `json:"slippage"`
	GrossR            float64   `json:"gross_r"`
	NetR              float64   `json:"net_r"`
	ExitReason        string    `json:"exit_reason"`
	HoldBars          int       `json:"hold_bars"`
}

type FuturesCleanBreakoutBaselineSummaryRow struct {
	CandidateID            string  `json:"candidate_id"`
	Timeframe              string  `json:"timeframe"`
	Split                  string  `json:"split"`
	Side                   string  `json:"side"`
	SignalCount            int     `json:"signal_count"`
	SkippedSignalCount     int     `json:"skipped_signal_count"`
	TotalTrades            int     `json:"total_trades"`
	Wins                   int     `json:"wins"`
	Losses                 int     `json:"losses"`
	WinRate                float64 `json:"win_rate"`
	GrossPnL               float64 `json:"gross_pnl"`
	NetPnL                 float64 `json:"net_pnl"`
	TotalCosts             float64 `json:"total_costs"`
	ProfitFactor           float64 `json:"profit_factor"`
	GrossProfitFactor      float64 `json:"gross_profit_factor"`
	MaxDrawdown            float64 `json:"max_drawdown"`
	AvgGrossR              float64 `json:"avg_gross_r"`
	AvgNetR                float64 `json:"avg_net_r"`
	AvgInitialRisk         float64 `json:"avg_initial_risk"`
	AvgHoldBars            float64 `json:"avg_hold_bars"`
	PassesOptimizationGate bool    `json:"passes_optimization_gate"`
	NearViableForPortfolio bool    `json:"near_viable_for_portfolio"`
	IsWorstPeriodSplit     bool    `json:"is_worst_period_split"`
}

func DefaultFuturesCleanBreakoutBaselineConfig() FuturesCleanBreakoutBaselineConfig {
	return FuturesCleanBreakoutBaselineConfig{
		EventDelayBars:             24,
		MaxHoldBars:                12,
		DetectorLookbackDays:       20,
		DetectorPercentile:         0.30,
		DetectorMinConsecutiveBars: 12,
		MinFullTrades:              100,
		MinKeySplitTrades:          25,
		NearFlatNetPct:             0.005,
		Candidates: []FuturesCleanBreakoutBaselineCandidateConfig{
			{CandidateID: CleanBreakoutCandidate4HUp, Timeframe: RangeDiscoveryTimeframe4h, Side: RangeDiscoverySideUp, MaxHoldBars: 12},
			{CandidateID: CleanBreakoutCandidate1HAll, Timeframe: RangeDiscoveryTimeframe1h, Side: RangeDiscoverySideAll, MaxHoldBars: 12},
		},
	}
}

func RunFuturesCleanBreakoutBaselineBacktest(parent []Candle, cfg FuturesCleanBreakoutBaselineConfig, btCfg BacktestConfig, splits []Split) (FuturesCleanBreakoutBaselineResult, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return FuturesCleanBreakoutBaselineResult{}, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}

	frameCandles := map[string][]Candle{}
	coverageByFrame := map[string]FuturesRangeDiscoveryCoverageRow{}
	coverageRows := []FuturesRangeDiscoveryCoverageRow{}
	for _, candidate := range cfg.Candidates {
		if _, ok := frameCandles[candidate.Timeframe]; ok {
			continue
		}
		frame, ok := cleanBreakoutFrameDef(candidate.Timeframe)
		if !ok {
			return FuturesCleanBreakoutBaselineResult{}, fmt.Errorf("unsupported clean breakout timeframe %q", candidate.Timeframe)
		}
		candles, coverage, err := resampleRangeDiscoveryFrame(parent, frame)
		coverageRows = append(coverageRows, coverage)
		coverageByFrame[candidate.Timeframe] = coverage
		if err != nil {
			return FuturesCleanBreakoutBaselineResult{CoverageRows: coverageRows}, err
		}
		frameCandles[candidate.Timeframe] = candles
	}

	result := FuturesCleanBreakoutBaselineResult{CoverageRows: coverageRows}
	for _, candidate := range cfg.Candidates {
		candles := frameCandles[candidate.Timeframe]
		coverage := coverageByFrame[candidate.Timeframe]
		if !coverage.Complete || coverage.ValidationStatus != "accepted" {
			return result, fmt.Errorf("clean breakout %s resample rejected: %s", candidate.Timeframe, coverage.ValidationError)
		}
		if len(candles) == 0 {
			continue
		}
		strategy, err := NewFuturesCleanBreakoutBaselineStrategy(candles, candidate, cfg, btCfg, splits)
		if err != nil {
			return result, err
		}
		run := RunBacktest(candles, strategy, btCfg)
		result.SignalRows = append(result.SignalRows, strategy.SignalRows()...)
		result.TradeRows = append(result.TradeRows, strategy.TradeRows(run.Trades, splits)...)
		result.Trades = append(result.Trades, run.Trades...)
	}
	sort.Slice(result.SignalRows, func(i, j int) bool {
		if result.SignalRows[i].Timeframe != result.SignalRows[j].Timeframe {
			return cleanBreakoutTimeframeSortKey(result.SignalRows[i].Timeframe) < cleanBreakoutTimeframeSortKey(result.SignalRows[j].Timeframe)
		}
		if result.SignalRows[i].BreakoutIndex != result.SignalRows[j].BreakoutIndex {
			return result.SignalRows[i].BreakoutIndex < result.SignalRows[j].BreakoutIndex
		}
		return result.SignalRows[i].SignalID < result.SignalRows[j].SignalID
	})
	sort.Slice(result.TradeRows, func(i, j int) bool {
		if result.TradeRows[i].EntryTime != result.TradeRows[j].EntryTime {
			return result.TradeRows[i].EntryTime < result.TradeRows[j].EntryTime
		}
		return result.TradeRows[i].SignalID < result.TradeRows[j].SignalID
	})
	sort.Slice(result.Trades, func(i, j int) bool {
		if result.Trades[i].EntryTime != result.Trades[j].EntryTime {
			return result.Trades[i].EntryTime < result.Trades[j].EntryTime
		}
		return result.Trades[i].Signal < result.Trades[j].Signal
	})
	result.SummaryRows = SummarizeFuturesCleanBreakoutBaseline(result.SignalRows, result.TradeRows, cfg, btCfg.StartBalance, splits)
	return result, nil
}

func NewFuturesCleanBreakoutBaselineStrategy(candles []Candle, candidate FuturesCleanBreakoutBaselineCandidateConfig, cfg FuturesCleanBreakoutBaselineConfig, btCfg BacktestConfig, splits []Split) (FuturesCleanBreakoutBaselineStrategy, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validateCandidate(candidate); err != nil {
		return FuturesCleanBreakoutBaselineStrategy{}, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	frame, ok := cleanBreakoutFrameDef(candidate.Timeframe)
	if !ok {
		return FuturesCleanBreakoutBaselineStrategy{}, fmt.Errorf("unsupported clean breakout timeframe %q", candidate.Timeframe)
	}
	detectorCfg := cleanBreakoutDetectorConfig(cfg, frame.barsPerDay)
	classifications, err := (CompressionRangeDetector{Config: detectorCfg}).Classify(candles)
	if err != nil {
		return FuturesCleanBreakoutBaselineStrategy{}, err
	}
	return newFuturesCleanBreakoutBaselineStrategyFromClassifications(candles, candidate, cfg, btCfg, classifications, splits)
}

func newFuturesCleanBreakoutBaselineStrategyFromClassifications(candles []Candle, candidate FuturesCleanBreakoutBaselineCandidateConfig, cfg FuturesCleanBreakoutBaselineConfig, btCfg BacktestConfig, classifications []RangeClassification, splits []Split) (FuturesCleanBreakoutBaselineStrategy, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validateCandidate(candidate); err != nil {
		return FuturesCleanBreakoutBaselineStrategy{}, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	if len(classifications) > len(candles) {
		classifications = classifications[:len(candles)]
	}
	normalizedATR := NormalizedATR(candles, DefaultCompressionRangeDetectorConfig().ATRPeriod)
	episodes := rangeRegimeDurabilityEpisodes(candles, classifications, normalizedATR, splits, BalancedDetectorProfileID)
	strategy := FuturesCleanBreakoutBaselineStrategy{
		candidate:   candidate,
		signalsByID: map[string]FuturesCleanBreakoutBaselineSignalRow{},
		byIndex:     map[int]Signal{},
	}
	for _, episode := range episodes {
		for eventIndex := episode.EndIndex + 1; eventIndex < len(candles) && eventIndex <= episode.EndIndex+cfg.EventDelayBars; eventIndex++ {
			if !cleanBreakoutCandidateAllowsEvent(candles[eventIndex], episode, candidate) {
				continue
			}
			row := newFuturesCleanBreakoutBaselineSignalRow(candles, candidate, episode, eventIndex, len(strategy.signals)+1, cfg, btCfg, splits)
			if row.SignalID != "" && row.SkippedReason == "" {
				if _, exists := strategy.byIndex[row.BreakoutIndex]; exists {
					row.SkippedReason = "duplicate_breakout_index"
				} else {
					strategy.byIndex[row.BreakoutIndex] = Signal{
						Side:        row.Side,
						Stop:        row.Stop,
						Target:      row.Target,
						MaxHoldBars: row.MaxHoldBars,
						Reason:      row.SignalID,
					}
					strategy.signalsByID[row.SignalID] = row
				}
			}
			strategy.signals = append(strategy.signals, row)
		}
	}
	return strategy, nil
}

func (s FuturesCleanBreakoutBaselineStrategy) Name() string {
	return s.candidate.CandidateID
}

func (s FuturesCleanBreakoutBaselineStrategy) OnCandle(ctx StrategyContext) (Signal, bool) {
	sig, ok := s.byIndex[ctx.Index]
	return sig, ok
}

func (s FuturesCleanBreakoutBaselineStrategy) SignalRows() []FuturesCleanBreakoutBaselineSignalRow {
	return append([]FuturesCleanBreakoutBaselineSignalRow(nil), s.signals...)
}

func (s FuturesCleanBreakoutBaselineStrategy) TradeRows(trades []Trade, splits []Split) []FuturesCleanBreakoutBaselineTradeRow {
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	rows := make([]FuturesCleanBreakoutBaselineTradeRow, 0, len(trades))
	for _, trade := range trades {
		signal, ok := s.signalsByID[trade.Signal]
		if !ok {
			continue
		}
		entryTime, _ := parseTime(trade.EntryTime)
		exitTime, _ := parseTime(trade.ExitTime)
		initialRisk := math.Abs(trade.EntryPrice-trade.Stop) * trade.Size
		row := FuturesCleanBreakoutBaselineTradeRow{
			SignalID:          signal.SignalID,
			CandidateID:       signal.CandidateID,
			Timeframe:         signal.Timeframe,
			Family:            signal.Family,
			EpisodeID:         signal.EpisodeID,
			EpisodeStartIndex: signal.EpisodeStartIndex,
			EpisodeEndIndex:   signal.EpisodeEndIndex,
			EpisodeStartTime:  signal.EpisodeStartTime,
			EpisodeEndTime:    signal.EpisodeEndTime,
			EpisodeHigh:       signal.EpisodeHigh,
			EpisodeLow:        signal.EpisodeLow,
			EpisodeWidth:      signal.EpisodeWidth,
			BreakoutIndex:     signal.BreakoutIndex,
			BreakoutCloseTime: signal.BreakoutCloseTime,
			BreakoutDelayBars: signal.BreakoutDelayBars,
			BreakoutClose:     signal.BreakoutClose,
			EntrySplit:        splitNameForCloseTime(entryTime, splits),
			CloseSplit:        splitNameForCloseTime(exitTime, splits),
			Side:              trade.Side,
			EntryTime:         trade.EntryTime,
			ExitTime:          trade.ExitTime,
			OpenIndex:         trade.OpenIndex,
			CloseIndex:        trade.CloseIndex,
			EntryPrice:        trade.EntryPrice,
			ExitPrice:         trade.ExitPrice,
			Stop:              trade.Stop,
			Target:            trade.Target,
			Size:              trade.Size,
			InitialRisk:       initialRisk,
			GrossPnL:          trade.GrossPnL,
			NetPnL:            trade.NetPnL,
			Fees:              trade.Fees,
			Slippage:          trade.Slippage,
			ExitReason:        trade.Reason,
			HoldBars:          trade.HoldBars,
		}
		if initialRisk > 0 {
			row.GrossR = trade.GrossPnL / initialRisk
			row.NetR = trade.NetPnL / initialRisk
		}
		rows = append(rows, row)
	}
	return rows
}

func newFuturesCleanBreakoutBaselineSignalRow(candles []Candle, candidate FuturesCleanBreakoutBaselineCandidateConfig, episode rangeRegimeDurabilityEpisode, eventIndex int, sequence int, cfg FuturesCleanBreakoutBaselineConfig, btCfg BacktestConfig, splits []Split) FuturesCleanBreakoutBaselineSignalRow {
	event := candles[eventIndex]
	width := episode.High - episode.Low
	side := Long
	if event.Close < episode.Low {
		side = Short
	}
	signalID := fmt.Sprintf("%s_%06d", candidate.CandidateID, sequence)
	entryIndex := eventIndex + 1
	row := FuturesCleanBreakoutBaselineSignalRow{
		SignalID:                signalID,
		CandidateID:             candidate.CandidateID,
		Timeframe:               candidate.Timeframe,
		Family:                  RangeDiscoveryFamilyCleanBreakout,
		Split:                   splitNameForCloseTime(event.CloseTime, splits),
		EpisodeID:               episode.EpisodeID,
		EpisodeStartIndex:       episode.StartIndex,
		EpisodeEndIndex:         episode.EndIndex,
		EpisodeStartTime:        candles[episode.StartIndex].CloseTime.UTC().Format(timeLayout),
		EpisodeEndTime:          candles[episode.EndIndex].CloseTime.UTC().Format(timeLayout),
		EpisodeHigh:             episode.High,
		EpisodeLow:              episode.Low,
		EpisodeMid:              (episode.High + episode.Low) / 2,
		EpisodeWidth:            width,
		EpisodeWidthPct:         episode.WidthPct,
		EpisodeRawLengthBars:    episode.RawLengthBars,
		EpisodeActiveLengthBars: episode.ActiveLengthBars,
		BreakoutIndex:           eventIndex,
		BreakoutOpenTime:        event.OpenTime.UTC().Format(timeLayout),
		BreakoutCloseTime:       event.CloseTime.UTC().Format(timeLayout),
		BreakoutDelayBars:       eventIndex - episode.EndIndex,
		BreakoutOpen:            event.Open,
		BreakoutHigh:            event.High,
		BreakoutLow:             event.Low,
		BreakoutClose:           event.Close,
		Side:                    side,
		EntryIndex:              entryIndex,
		MaxHoldBars:             candidateMaxHoldBars(candidate, cfg),
	}
	if width <= 0 {
		row.SkippedReason = "non_positive_range_width"
		return row
	}
	if entryIndex >= len(candles) {
		row.SkippedReason = "missing_entry_candle"
		return row
	}
	entry := candles[entryIndex]
	row.EntryOpenTime = entry.OpenTime.UTC().Format(timeLayout)
	row.EntryOpen = entry.Open
	row.ExpectedEntryPrice = applySlippage(entry.Open, btCfg.SlippagePct, side, true)
	if side == Long {
		row.Stop = episode.High
		row.Target = row.ExpectedEntryPrice + width
	} else {
		row.Stop = episode.Low
		row.Target = row.ExpectedEntryPrice - width
	}
	if row.Stop <= 0 || row.Target <= 0 || row.ExpectedEntryPrice <= 0 {
		row.SkippedReason = "non_positive_trade_price"
		return row
	}
	row.EntryGeometryValid = validEntryGeometry(Signal{Side: row.Side, Stop: row.Stop, Target: row.Target}, row.ExpectedEntryPrice)
	if !row.EntryGeometryValid {
		row.SkippedReason = "invalid_entry_geometry"
	}
	return row
}

func SummarizeFuturesCleanBreakoutBaseline(signals []FuturesCleanBreakoutBaselineSignalRow, trades []FuturesCleanBreakoutBaselineTradeRow, cfg FuturesCleanBreakoutBaselineConfig, startBalance float64, splits []Split) []FuturesCleanBreakoutBaselineSummaryRow {
	cfg = cfg.withDefaults()
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	rows := make([]FuturesCleanBreakoutBaselineSummaryRow, 0, len(cfg.Candidates)*len(splits)*3)
	for _, candidate := range cfg.Candidates {
		for _, split := range splits {
			for _, side := range []string{"all", string(Long), string(Short)} {
				filteredTrades := filterCleanBreakoutTrades(trades, candidate.CandidateID, split, side)
				filteredSignals := filterCleanBreakoutSignals(signals, candidate.CandidateID, split, side)
				row := summarizeCleanBreakoutTrades(filteredTrades, startBalance)
				row.CandidateID = candidate.CandidateID
				row.Timeframe = candidate.Timeframe
				row.Split = split.Name
				row.Side = side
				row.SignalCount = len(filteredSignals)
				for _, signal := range filteredSignals {
					if signal.SkippedReason != "" {
						row.SkippedSignalCount++
					}
				}
				rows = append(rows, row)
			}
		}
	}
	markCleanBreakoutReviewFlags(rows, cfg, startBalance, splits)
	return rows
}

func FuturesCleanBreakoutBaselineStopState(summary []FuturesCleanBreakoutBaselineSummaryRow, cfg FuturesCleanBreakoutBaselineConfig, startBalance float64, splits []Split) string {
	cfg = cfg.withDefaults()
	if len(summary) == 0 {
		return CleanBreakoutStopStateFailedNoPromotion
	}
	passing := 0
	nearViable := 0
	for _, candidate := range cfg.Candidates {
		if cleanBreakoutCandidatePasses(summary, candidate, cfg, startBalance, splits) {
			passing++
		}
		if cleanBreakoutCandidateNearViable(summary, candidate, cfg, startBalance, splits) {
			nearViable++
		}
	}
	if passing > 0 {
		return CleanBreakoutStopStatePassedNeedsOptimization
	}
	if nearViable >= 2 {
		return CleanBreakoutStopStateMixedNeedsPortfolioReview
	}
	return CleanBreakoutStopStateFailedNoPromotion
}

func (cfg FuturesCleanBreakoutBaselineConfig) withDefaults() FuturesCleanBreakoutBaselineConfig {
	defaults := DefaultFuturesCleanBreakoutBaselineConfig()
	if cfg.EventDelayBars == 0 {
		cfg.EventDelayBars = defaults.EventDelayBars
	}
	if cfg.MaxHoldBars == 0 {
		cfg.MaxHoldBars = defaults.MaxHoldBars
	}
	if cfg.DetectorLookbackDays == 0 {
		cfg.DetectorLookbackDays = defaults.DetectorLookbackDays
	}
	if cfg.DetectorPercentile == 0 {
		cfg.DetectorPercentile = defaults.DetectorPercentile
	}
	if cfg.DetectorMinConsecutiveBars == 0 {
		cfg.DetectorMinConsecutiveBars = defaults.DetectorMinConsecutiveBars
	}
	if cfg.MinFullTrades == 0 {
		cfg.MinFullTrades = defaults.MinFullTrades
	}
	if cfg.MinKeySplitTrades == 0 {
		cfg.MinKeySplitTrades = defaults.MinKeySplitTrades
	}
	if cfg.NearFlatNetPct == 0 {
		cfg.NearFlatNetPct = defaults.NearFlatNetPct
	}
	if len(cfg.Candidates) == 0 {
		cfg.Candidates = append([]FuturesCleanBreakoutBaselineCandidateConfig(nil), defaults.Candidates...)
	}
	for i := range cfg.Candidates {
		if cfg.Candidates[i].MaxHoldBars == 0 {
			cfg.Candidates[i].MaxHoldBars = cfg.MaxHoldBars
		}
	}
	return cfg
}

func (cfg FuturesCleanBreakoutBaselineConfig) validate() error {
	if cfg.EventDelayBars <= 0 {
		return fmt.Errorf("clean breakout event delay bars must be positive")
	}
	if cfg.MaxHoldBars <= 0 {
		return fmt.Errorf("clean breakout max hold bars must be positive")
	}
	if cfg.DetectorLookbackDays <= 0 {
		return fmt.Errorf("clean breakout detector lookback days must be positive")
	}
	if cfg.DetectorPercentile <= 0 || cfg.DetectorPercentile >= 1 {
		return fmt.Errorf("clean breakout detector percentile must be between 0 and 1")
	}
	if cfg.DetectorMinConsecutiveBars <= 0 {
		return fmt.Errorf("clean breakout detector min consecutive bars must be positive")
	}
	if cfg.MinFullTrades <= 0 {
		return fmt.Errorf("clean breakout min full trades must be positive")
	}
	if cfg.MinKeySplitTrades <= 0 {
		return fmt.Errorf("clean breakout min key split trades must be positive")
	}
	if cfg.NearFlatNetPct <= 0 {
		return fmt.Errorf("clean breakout near-flat net pct must be positive")
	}
	if len(cfg.Candidates) == 0 {
		return fmt.Errorf("clean breakout requires at least one candidate")
	}
	for _, candidate := range cfg.Candidates {
		if err := cfg.validateCandidate(candidate); err != nil {
			return err
		}
	}
	return nil
}

func (cfg FuturesCleanBreakoutBaselineConfig) validateCandidate(candidate FuturesCleanBreakoutBaselineCandidateConfig) error {
	if candidate.CandidateID == "" {
		return fmt.Errorf("clean breakout candidate id must not be empty")
	}
	if candidate.Timeframe != RangeDiscoveryTimeframe1h && candidate.Timeframe != RangeDiscoveryTimeframe4h {
		return fmt.Errorf("clean breakout candidate %s unsupported timeframe %q", candidate.CandidateID, candidate.Timeframe)
	}
	if candidate.Side != RangeDiscoverySideAll && candidate.Side != RangeDiscoverySideUp && candidate.Side != RangeDiscoverySideDown {
		return fmt.Errorf("clean breakout candidate %s unsupported side %q", candidate.CandidateID, candidate.Side)
	}
	if candidate.MaxHoldBars < 0 {
		return fmt.Errorf("clean breakout candidate %s max hold bars cannot be negative", candidate.CandidateID)
	}
	return nil
}

func cleanBreakoutDetectorConfig(cfg FuturesCleanBreakoutBaselineConfig, barsPerDay int) RangeDetectorConfig {
	detectorCfg := DefaultCompressionRangeDetectorConfig()
	detectorCfg.BarsPerDay = barsPerDay
	detectorCfg.LookbackDays = cfg.DetectorLookbackDays
	detectorCfg.Percentile = cfg.DetectorPercentile
	detectorCfg.MinConsecutiveBars = cfg.DetectorMinConsecutiveBars
	detectorCfg.UseBollinger = true
	detectorCfg.UseADX = false
	return detectorCfg
}

func cleanBreakoutFrameDef(timeframe string) (rangeDiscoveryFrameDef, bool) {
	for _, frame := range rangeDiscoveryFrameDefs() {
		if frame.timeframe == timeframe {
			return frame, true
		}
	}
	return rangeDiscoveryFrameDef{}, false
}

func cleanBreakoutCandidateAllowsEvent(candle Candle, episode rangeRegimeDurabilityEpisode, candidate FuturesCleanBreakoutBaselineCandidateConfig) bool {
	switch {
	case candle.Close > episode.High:
		return candidate.Side == RangeDiscoverySideAll || candidate.Side == RangeDiscoverySideUp
	case candle.Close < episode.Low:
		return candidate.Side == RangeDiscoverySideAll || candidate.Side == RangeDiscoverySideDown
	default:
		return false
	}
}

func candidateMaxHoldBars(candidate FuturesCleanBreakoutBaselineCandidateConfig, cfg FuturesCleanBreakoutBaselineConfig) int {
	if candidate.MaxHoldBars > 0 {
		return candidate.MaxHoldBars
	}
	return cfg.MaxHoldBars
}

func filterCleanBreakoutSignals(signals []FuturesCleanBreakoutBaselineSignalRow, candidateID string, split Split, side string) []FuturesCleanBreakoutBaselineSignalRow {
	out := make([]FuturesCleanBreakoutBaselineSignalRow, 0, len(signals))
	for _, signal := range signals {
		if signal.CandidateID != candidateID {
			continue
		}
		eventTime, err := parseTime(signal.BreakoutCloseTime)
		if err != nil || !split.Contains(eventTime) {
			continue
		}
		if side != "all" && string(signal.Side) != side {
			continue
		}
		out = append(out, signal)
	}
	return out
}

func filterCleanBreakoutTrades(trades []FuturesCleanBreakoutBaselineTradeRow, candidateID string, split Split, side string) []FuturesCleanBreakoutBaselineTradeRow {
	out := make([]FuturesCleanBreakoutBaselineTradeRow, 0, len(trades))
	for _, trade := range trades {
		if trade.CandidateID != candidateID {
			continue
		}
		exitTime, err := parseTime(trade.ExitTime)
		if err != nil || !split.Contains(exitTime) {
			continue
		}
		if side != "all" && string(trade.Side) != side {
			continue
		}
		out = append(out, trade)
	}
	return out
}

func summarizeCleanBreakoutTrades(trades []FuturesCleanBreakoutBaselineTradeRow, startBalance float64) FuturesCleanBreakoutBaselineSummaryRow {
	row := FuturesCleanBreakoutBaselineSummaryRow{TotalTrades: len(trades)}
	balance := startBalance
	equity := []float64{startBalance}
	netProfit, netLoss := 0.0, 0.0
	grossProfit, grossLoss := 0.0, 0.0
	holdBars := 0
	for _, trade := range trades {
		row.GrossPnL += trade.GrossPnL
		row.NetPnL += trade.NetPnL
		row.TotalCosts += trade.Fees + trade.Slippage
		row.AvgGrossR += trade.GrossR
		row.AvgNetR += trade.NetR
		row.AvgInitialRisk += trade.InitialRisk
		holdBars += trade.HoldBars
		if trade.NetPnL > 0 {
			row.Wins++
			netProfit += trade.NetPnL
		} else if trade.NetPnL < 0 {
			row.Losses++
			netLoss += -trade.NetPnL
		}
		if trade.GrossPnL > 0 {
			grossProfit += trade.GrossPnL
		} else if trade.GrossPnL < 0 {
			grossLoss += -trade.GrossPnL
		}
		balance += trade.NetPnL
		equity = append(equity, balance)
	}
	if row.TotalTrades > 0 {
		count := float64(row.TotalTrades)
		row.WinRate = float64(row.Wins) / count
		row.AvgGrossR /= count
		row.AvgNetR /= count
		row.AvgInitialRisk /= count
		row.AvgHoldBars = float64(holdBars) / count
	}
	if netLoss > 0 {
		row.ProfitFactor = netProfit / netLoss
	} else if netProfit > 0 {
		row.ProfitFactor = 999.99
	}
	if grossLoss > 0 {
		row.GrossProfitFactor = grossProfit / grossLoss
	} else if grossProfit > 0 {
		row.GrossProfitFactor = 999.99
	}
	row.MaxDrawdown = MaxDrawdown(equity)
	return row
}

func markCleanBreakoutReviewFlags(rows []FuturesCleanBreakoutBaselineSummaryRow, cfg FuturesCleanBreakoutBaselineConfig, startBalance float64, splits []Split) {
	for _, candidate := range cfg.Candidates {
		passes := cleanBreakoutCandidatePasses(rows, candidate, cfg, startBalance, splits)
		near := cleanBreakoutCandidateNearViable(rows, candidate, cfg, startBalance, splits)
		worstIndex := -1
		for i := range rows {
			if rows[i].CandidateID != candidate.CandidateID {
				continue
			}
			if rows[i].Split != fullSplitName && rows[i].Side == "all" {
				if worstIndex == -1 || rows[i].NetPnL < rows[worstIndex].NetPnL {
					worstIndex = i
				}
			}
			if rows[i].Side == "all" {
				rows[i].PassesOptimizationGate = passes
				rows[i].NearViableForPortfolio = near
			}
		}
		if worstIndex >= 0 {
			rows[worstIndex].IsWorstPeriodSplit = true
		}
	}
}

func cleanBreakoutCandidatePasses(rows []FuturesCleanBreakoutBaselineSummaryRow, candidate FuturesCleanBreakoutBaselineCandidateConfig, cfg FuturesCleanBreakoutBaselineConfig, startBalance float64, splits []Split) bool {
	byKey := cleanBreakoutSummaryByKey(rows)
	full := byKey[candidate.CandidateID+"|"+fullSplitName+"|all"]
	if full.TotalTrades < cfg.MinFullTrades || full.GrossPnL <= 0 || full.NetPnL <= 0 || full.ProfitFactor < 1.2 {
		return false
	}
	for _, split := range periodSplits(splits) {
		row := byKey[candidate.CandidateID+"|"+split.Name+"|all"]
		if (split.Name == "2023_2024_oos" || split.Name == "2025_2026_recent") && row.TotalTrades < cfg.MinKeySplitTrades {
			return false
		}
		if !cleanBreakoutSplitNonNegativeOrNearFlat(row, startBalance, cfg.NearFlatNetPct) {
			return false
		}
	}
	if candidate.Side == RangeDiscoverySideAll {
		for _, side := range []string{string(Long), string(Short)} {
			row := byKey[candidate.CandidateID+"|"+fullSplitName+"|"+side]
			if row.TotalTrades < cfg.MinKeySplitTrades || !cleanBreakoutSplitNonNegativeOrNearFlat(row, startBalance, cfg.NearFlatNetPct) {
				return false
			}
		}
	}
	return true
}

func cleanBreakoutCandidateNearViable(rows []FuturesCleanBreakoutBaselineSummaryRow, candidate FuturesCleanBreakoutBaselineCandidateConfig, cfg FuturesCleanBreakoutBaselineConfig, startBalance float64, splits []Split) bool {
	byKey := cleanBreakoutSummaryByKey(rows)
	full := byKey[candidate.CandidateID+"|"+fullSplitName+"|all"]
	if full.TotalTrades < cfg.MinFullTrades || full.GrossPnL <= 0 || full.ProfitFactor < 1.0 {
		return false
	}
	if full.NetPnL < -startBalance*cfg.NearFlatNetPct {
		return false
	}
	for _, splitName := range []string{"2023_2024_oos", "2025_2026_recent"} {
		row := byKey[candidate.CandidateID+"|"+splitName+"|all"]
		if row.TotalTrades < cfg.MinKeySplitTrades || !cleanBreakoutSplitNonNegativeOrNearFlat(row, startBalance, cfg.NearFlatNetPct) {
			return false
		}
	}
	return true
}

func cleanBreakoutSplitNonNegativeOrNearFlat(row FuturesCleanBreakoutBaselineSummaryRow, startBalance float64, nearFlatPct float64) bool {
	if row.TotalTrades == 0 {
		return false
	}
	if row.NetPnL >= 0 {
		return true
	}
	return row.NetPnL >= -startBalance*nearFlatPct && row.GrossPnL > 0
}

func cleanBreakoutSummaryByKey(rows []FuturesCleanBreakoutBaselineSummaryRow) map[string]FuturesCleanBreakoutBaselineSummaryRow {
	byKey := map[string]FuturesCleanBreakoutBaselineSummaryRow{}
	for _, row := range rows {
		byKey[row.CandidateID+"|"+row.Split+"|"+row.Side] = row
	}
	return byKey
}

func periodSplits(splits []Split) []Split {
	out := make([]Split, 0, len(splits))
	for _, split := range splits {
		if split.Name == fullSplitName {
			continue
		}
		out = append(out, split)
	}
	return out
}

func cleanBreakoutTimeframeSortKey(timeframe string) int {
	switch timeframe {
	case RangeDiscoveryTimeframe1h:
		return 1
	case RangeDiscoveryTimeframe4h:
		return 2
	default:
		return 99
	}
}
