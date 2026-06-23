package lab

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

const (
	FuturesRangeUniverseStructuredCompressionBaselineName = "futures_range_universe_structured_compression_baseline"

	StructuredCompressionCandidate4HAllH6  = "structured_compression_4h_all_h6"
	StructuredCompressionCandidate1HAllH12 = "structured_compression_1h_all_h12"

	StructuredCompressionStopStateSourceGap            = "structured_compression_baseline_source_gap"
	StructuredCompressionStopStateCodegenOrTestBlocked = "structured_compression_baseline_codegen_or_test_blocked"
	StructuredCompressionStopStateFailedNoPromotion    = "structured_compression_baseline_failed_no_promotion"
	StructuredCompressionStopStatePassedNeedsOptimize  = "structured_compression_baseline_passed_needs_optimization_brief"
	StructuredCompressionStopStateMixedPortfolioReview = "structured_compression_baseline_mixed_needs_portfolio_stream_review"
	StructuredCompressionStopStateReviewOnlyNoStrategy = "structured_compression_baseline_review_only_no_strategy_change"
	StructuredCompressionSummaryAggregateSymbol        = "all"
	StructuredCompressionConfirmationWindowDefaultBars = 3
	StructuredCompressionEventDelayDefaultBars         = 24
)

type FuturesRangeUniverseStructuredCompressionBaselineConfig struct {
	Sources                    []FuturesRangeUniverseSourceConfig
	EventDelayBars             int
	ConfirmationWindowBars     int
	DetectorLookbackDays       int
	DetectorPercentile         float64
	DetectorMinConsecutiveBars int
	MinFullTrades              int
	MinKeySplitTrades          int
	NearFlatNetPct             float64
	Candidates                 []FuturesRangeUniverseStructuredCompressionCandidateConfig
}

type FuturesRangeUniverseStructuredCompressionCandidateConfig struct {
	CandidateID                  string
	Timeframe                    string
	Side                         string
	MaxHoldBars                  int
	TargetRangeWidthMultiple     float64
	StopBoundaryBufferRangeWidth float64
}

type FuturesRangeUniverseStructuredCompressionBaselineResult struct {
	SourceRows   []FuturesRangeUniverseSourceRow
	CoverageRows []FuturesRangeUniverseCoverageRow
	SignalRows   []FuturesRangeUniverseStructuredCompressionSignalRow
	TradeRows    []FuturesRangeUniverseStructuredCompressionTradeRow
	SummaryRows  []FuturesRangeUniverseStructuredCompressionSummaryRow
	Trades       []Trade
}

type FuturesRangeUniverseStructuredCompressionStrategy struct {
	symbol      string
	candidate   FuturesRangeUniverseStructuredCompressionCandidateConfig
	signals     []FuturesRangeUniverseStructuredCompressionSignalRow
	signalsByID map[string]FuturesRangeUniverseStructuredCompressionSignalRow
	byIndex     map[int]Signal
}

type FuturesRangeUniverseStructuredCompressionSignalRow struct {
	SignalID                string    `json:"signal_id"`
	Symbol                  string    `json:"symbol"`
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
	ConfirmationIndex       int       `json:"confirmation_index"`
	ConfirmationOpenTime    string    `json:"confirmation_open_time"`
	ConfirmationCloseTime   string    `json:"confirmation_close_time"`
	ConfirmationDelayBars   int       `json:"confirmation_delay_bars"`
	ConfirmationOpen        float64   `json:"confirmation_open"`
	ConfirmationHigh        float64   `json:"confirmation_high"`
	ConfirmationLow         float64   `json:"confirmation_low"`
	ConfirmationClose       float64   `json:"confirmation_close"`
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

type FuturesRangeUniverseStructuredCompressionTradeRow struct {
	SignalID              string    `json:"signal_id"`
	Symbol                string    `json:"symbol"`
	CandidateID           string    `json:"candidate_id"`
	Timeframe             string    `json:"timeframe"`
	Family                string    `json:"family"`
	EpisodeID             int       `json:"episode_id"`
	EpisodeStartIndex     int       `json:"episode_start_index"`
	EpisodeEndIndex       int       `json:"episode_end_index"`
	EpisodeStartTime      string    `json:"episode_start_time"`
	EpisodeEndTime        string    `json:"episode_end_time"`
	EpisodeHigh           float64   `json:"episode_high"`
	EpisodeLow            float64   `json:"episode_low"`
	EpisodeWidth          float64   `json:"episode_width"`
	BreakoutIndex         int       `json:"breakout_index"`
	BreakoutCloseTime     string    `json:"breakout_close_time"`
	BreakoutDelayBars     int       `json:"breakout_delay_bars"`
	BreakoutClose         float64   `json:"breakout_close"`
	ConfirmationIndex     int       `json:"confirmation_index"`
	ConfirmationCloseTime string    `json:"confirmation_close_time"`
	ConfirmationDelayBars int       `json:"confirmation_delay_bars"`
	ConfirmationClose     float64   `json:"confirmation_close"`
	EntrySplit            string    `json:"entry_split"`
	CloseSplit            string    `json:"close_split"`
	Side                  Direction `json:"side"`
	EntryTime             string    `json:"entry_time"`
	ExitTime              string    `json:"exit_time"`
	OpenIndex             int       `json:"open_index"`
	CloseIndex            int       `json:"close_index"`
	EntryPrice            float64   `json:"entry_price"`
	ExitPrice             float64   `json:"exit_price"`
	Stop                  float64   `json:"stop"`
	Target                float64   `json:"target"`
	Size                  float64   `json:"size"`
	InitialRisk           float64   `json:"initial_risk"`
	GrossPnL              float64   `json:"gross_pnl"`
	NetPnL                float64   `json:"net_pnl"`
	Fees                  float64   `json:"fees"`
	Slippage              float64   `json:"slippage"`
	GrossR                float64   `json:"gross_r"`
	NetR                  float64   `json:"net_r"`
	ExitReason            string    `json:"exit_reason"`
	HoldBars              int       `json:"hold_bars"`
}

type FuturesRangeUniverseStructuredCompressionSummaryRow struct {
	CandidateID            string  `json:"candidate_id"`
	Symbol                 string  `json:"symbol"`
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

func DefaultFuturesRangeUniverseStructuredCompressionBaselineConfig() FuturesRangeUniverseStructuredCompressionBaselineConfig {
	return FuturesRangeUniverseStructuredCompressionBaselineConfig{
		Sources:                    DefaultFuturesRangeUniverseDiscoveryAuditConfig().Sources,
		EventDelayBars:             StructuredCompressionEventDelayDefaultBars,
		ConfirmationWindowBars:     StructuredCompressionConfirmationWindowDefaultBars,
		DetectorLookbackDays:       20,
		DetectorPercentile:         0.30,
		DetectorMinConsecutiveBars: 12,
		MinFullTrades:              100,
		MinKeySplitTrades:          25,
		NearFlatNetPct:             0.005,
		Candidates: []FuturesRangeUniverseStructuredCompressionCandidateConfig{
			{CandidateID: StructuredCompressionCandidate4HAllH6, Timeframe: RangeDiscoveryTimeframe4h, Side: RangeDiscoverySideAll, MaxHoldBars: 6},
			{CandidateID: StructuredCompressionCandidate1HAllH12, Timeframe: RangeDiscoveryTimeframe1h, Side: RangeDiscoverySideAll, MaxHoldBars: 12},
		},
	}
}

func RunFuturesRangeUniverseStructuredCompressionBaselineBacktest(cfg FuturesRangeUniverseStructuredCompressionBaselineConfig, btCfg BacktestConfig, splits []Split) (FuturesRangeUniverseStructuredCompressionBaselineResult, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return FuturesRangeUniverseStructuredCompressionBaselineResult{}, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}

	result := FuturesRangeUniverseStructuredCompressionBaselineResult{}
	for _, source := range cfg.Sources {
		candles, sourceRow, err := LoadFuturesRangeUniverseSource(source, splits)
		result.SourceRows = append(result.SourceRows, sourceRow)
		if err != nil {
			return result, err
		}

		frameCandles := map[string][]Candle{}
		coverageByFrame := map[string]FuturesRangeUniverseCoverageRow{}
		for _, candidate := range cfg.Candidates {
			if _, ok := frameCandles[candidate.Timeframe]; ok {
				continue
			}
			frame, ok := structuredCompressionFrameDef(candidate.Timeframe)
			if !ok {
				return result, fmt.Errorf("unsupported structured compression timeframe %q", candidate.Timeframe)
			}
			resampled, coverage, err := resampleRangeDiscoveryFrame(candles, frame)
			uCoverage := FuturesRangeUniverseCoverageRow{Symbol: sourceRow.Symbol, FuturesRangeDiscoveryCoverageRow: coverage}
			result.CoverageRows = append(result.CoverageRows, uCoverage)
			coverageByFrame[candidate.Timeframe] = uCoverage
			if err != nil {
				return result, err
			}
			frameCandles[candidate.Timeframe] = resampled
		}

		for _, candidate := range cfg.Candidates {
			coverage := coverageByFrame[candidate.Timeframe]
			if !coverage.Complete || coverage.ValidationStatus != "accepted" {
				return result, fmt.Errorf("%s %s structured compression resample rejected: %s", sourceRow.Symbol, candidate.Timeframe, coverage.ValidationError)
			}
			candlesForFrame := frameCandles[candidate.Timeframe]
			if len(candlesForFrame) == 0 {
				continue
			}
			strategy, err := NewFuturesRangeUniverseStructuredCompressionStrategy(candlesForFrame, sourceRow.Symbol, candidate, cfg, btCfg, splits)
			if err != nil {
				return result, err
			}
			run := RunBacktest(candlesForFrame, strategy, btCfg)
			result.SignalRows = append(result.SignalRows, strategy.SignalRows()...)
			result.TradeRows = append(result.TradeRows, strategy.TradeRows(run.Trades, splits)...)
			result.Trades = append(result.Trades, run.Trades...)
		}
	}

	sort.Slice(result.CoverageRows, func(i, j int) bool {
		if result.CoverageRows[i].Symbol != result.CoverageRows[j].Symbol {
			return rangeUniverseSymbolSortKey(result.CoverageRows[i].Symbol) < rangeUniverseSymbolSortKey(result.CoverageRows[j].Symbol)
		}
		return structuredCompressionTimeframeSortKey(result.CoverageRows[i].Timeframe) < structuredCompressionTimeframeSortKey(result.CoverageRows[j].Timeframe)
	})
	sort.Slice(result.SignalRows, func(i, j int) bool {
		if result.SignalRows[i].Symbol != result.SignalRows[j].Symbol {
			return rangeUniverseSymbolSortKey(result.SignalRows[i].Symbol) < rangeUniverseSymbolSortKey(result.SignalRows[j].Symbol)
		}
		if result.SignalRows[i].Timeframe != result.SignalRows[j].Timeframe {
			return structuredCompressionTimeframeSortKey(result.SignalRows[i].Timeframe) < structuredCompressionTimeframeSortKey(result.SignalRows[j].Timeframe)
		}
		if result.SignalRows[i].ConfirmationIndex != result.SignalRows[j].ConfirmationIndex {
			return result.SignalRows[i].ConfirmationIndex < result.SignalRows[j].ConfirmationIndex
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
	result.SummaryRows = SummarizeFuturesRangeUniverseStructuredCompressionBaseline(result.SignalRows, result.TradeRows, cfg, btCfg.StartBalance, splits)
	return result, nil
}

func NewFuturesRangeUniverseStructuredCompressionStrategy(candles []Candle, symbol string, candidate FuturesRangeUniverseStructuredCompressionCandidateConfig, cfg FuturesRangeUniverseStructuredCompressionBaselineConfig, btCfg BacktestConfig, splits []Split) (FuturesRangeUniverseStructuredCompressionStrategy, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validateCandidate(candidate); err != nil {
		return FuturesRangeUniverseStructuredCompressionStrategy{}, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	frame, ok := structuredCompressionFrameDef(candidate.Timeframe)
	if !ok {
		return FuturesRangeUniverseStructuredCompressionStrategy{}, fmt.Errorf("unsupported structured compression timeframe %q", candidate.Timeframe)
	}
	detectorCfg := structuredCompressionDetectorConfig(cfg, frame.barsPerDay)
	classifications, err := (CompressionRangeDetector{Config: detectorCfg}).Classify(candles)
	if err != nil {
		return FuturesRangeUniverseStructuredCompressionStrategy{}, err
	}
	return newFuturesRangeUniverseStructuredCompressionStrategyFromClassifications(candles, symbol, candidate, cfg, btCfg, classifications, splits)
}

func newFuturesRangeUniverseStructuredCompressionStrategyFromClassifications(candles []Candle, symbol string, candidate FuturesRangeUniverseStructuredCompressionCandidateConfig, cfg FuturesRangeUniverseStructuredCompressionBaselineConfig, btCfg BacktestConfig, classifications []RangeClassification, splits []Split) (FuturesRangeUniverseStructuredCompressionStrategy, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validateCandidate(candidate); err != nil {
		return FuturesRangeUniverseStructuredCompressionStrategy{}, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	if len(classifications) > len(candles) {
		classifications = classifications[:len(candles)]
	}
	detectorCfg := structuredCompressionDetectorConfig(cfg, structuredCompressionBarsPerDay(candidate.Timeframe))
	normalizedATR := NormalizedATR(candles, detectorCfg.ATRPeriod)
	episodes := rangeRegimeDurabilityEpisodes(candles, classifications, normalizedATR, splits, BalancedDetectorProfileID)
	strategy := FuturesRangeUniverseStructuredCompressionStrategy{
		symbol:      strings.ToUpper(symbol),
		candidate:   candidate,
		signalsByID: map[string]FuturesRangeUniverseStructuredCompressionSignalRow{},
		byIndex:     map[int]Signal{},
	}
	for _, episode := range episodes {
		if episode.High <= episode.Low {
			continue
		}
		if event, ok := firstStructuredCompressionSignalEvent(candles, episode, candidate, cfg); ok {
			row := newFuturesRangeUniverseStructuredCompressionSignalRow(candles, strategy.symbol, candidate, episode, event.breakoutIndex, event.confirmationIndex, len(strategy.signals)+1, cfg, btCfg, splits)
			if row.SignalID != "" && row.SkippedReason == "" {
				if _, exists := strategy.byIndex[row.ConfirmationIndex]; exists {
					row.SkippedReason = "duplicate_confirmation_index"
				} else {
					strategy.byIndex[row.ConfirmationIndex] = Signal{
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

func (s FuturesRangeUniverseStructuredCompressionStrategy) Name() string {
	return s.candidate.CandidateID
}

func (s FuturesRangeUniverseStructuredCompressionStrategy) OnCandle(ctx StrategyContext) (Signal, bool) {
	sig, ok := s.byIndex[ctx.Index]
	return sig, ok
}

func (s FuturesRangeUniverseStructuredCompressionStrategy) SignalRows() []FuturesRangeUniverseStructuredCompressionSignalRow {
	return append([]FuturesRangeUniverseStructuredCompressionSignalRow(nil), s.signals...)
}

func (s FuturesRangeUniverseStructuredCompressionStrategy) TradeRows(trades []Trade, splits []Split) []FuturesRangeUniverseStructuredCompressionTradeRow {
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	rows := make([]FuturesRangeUniverseStructuredCompressionTradeRow, 0, len(trades))
	for _, trade := range trades {
		signal, ok := s.signalsByID[trade.Signal]
		if !ok {
			continue
		}
		entryTime, _ := parseTime(trade.EntryTime)
		exitTime, _ := parseTime(trade.ExitTime)
		initialRisk := math.Abs(trade.EntryPrice-trade.Stop) * trade.Size
		row := FuturesRangeUniverseStructuredCompressionTradeRow{
			SignalID:              signal.SignalID,
			Symbol:                signal.Symbol,
			CandidateID:           signal.CandidateID,
			Timeframe:             signal.Timeframe,
			Family:                signal.Family,
			EpisodeID:             signal.EpisodeID,
			EpisodeStartIndex:     signal.EpisodeStartIndex,
			EpisodeEndIndex:       signal.EpisodeEndIndex,
			EpisodeStartTime:      signal.EpisodeStartTime,
			EpisodeEndTime:        signal.EpisodeEndTime,
			EpisodeHigh:           signal.EpisodeHigh,
			EpisodeLow:            signal.EpisodeLow,
			EpisodeWidth:          signal.EpisodeWidth,
			BreakoutIndex:         signal.BreakoutIndex,
			BreakoutCloseTime:     signal.BreakoutCloseTime,
			BreakoutDelayBars:     signal.BreakoutDelayBars,
			BreakoutClose:         signal.BreakoutClose,
			ConfirmationIndex:     signal.ConfirmationIndex,
			ConfirmationCloseTime: signal.ConfirmationCloseTime,
			ConfirmationDelayBars: signal.ConfirmationDelayBars,
			ConfirmationClose:     signal.ConfirmationClose,
			EntrySplit:            splitNameForCloseTime(entryTime, splits),
			CloseSplit:            splitNameForCloseTime(exitTime, splits),
			Side:                  trade.Side,
			EntryTime:             trade.EntryTime,
			ExitTime:              trade.ExitTime,
			OpenIndex:             trade.OpenIndex,
			CloseIndex:            trade.CloseIndex,
			EntryPrice:            trade.EntryPrice,
			ExitPrice:             trade.ExitPrice,
			Stop:                  trade.Stop,
			Target:                trade.Target,
			Size:                  trade.Size,
			InitialRisk:           initialRisk,
			GrossPnL:              trade.GrossPnL,
			NetPnL:                trade.NetPnL,
			Fees:                  trade.Fees,
			Slippage:              trade.Slippage,
			ExitReason:            trade.Reason,
			HoldBars:              trade.HoldBars,
		}
		if initialRisk > 0 {
			row.GrossR = trade.GrossPnL / initialRisk
			row.NetR = trade.NetPnL / initialRisk
		}
		rows = append(rows, row)
	}
	return rows
}

type structuredCompressionSignalEvent struct {
	breakoutIndex     int
	confirmationIndex int
}

func firstStructuredCompressionSignalEvent(candles []Candle, episode rangeRegimeDurabilityEpisode, candidate FuturesRangeUniverseStructuredCompressionCandidateConfig, cfg FuturesRangeUniverseStructuredCompressionBaselineConfig) (structuredCompressionSignalEvent, bool) {
	for breakIndex := episode.EndIndex + 1; breakIndex < len(candles) && breakIndex <= episode.EndIndex+cfg.EventDelayBars; breakIndex++ {
		breakout := candles[breakIndex]
		if breakout.Close > episode.High && (candidate.Side == RangeDiscoverySideAll || candidate.Side == RangeDiscoverySideUp) {
			for confirmIndex := breakIndex + 1; confirmIndex < len(candles) && confirmIndex <= breakIndex+cfg.ConfirmationWindowBars; confirmIndex++ {
				confirm := candles[confirmIndex]
				if confirm.Close > episode.High && confirm.High > breakout.High {
					return structuredCompressionSignalEvent{breakoutIndex: breakIndex, confirmationIndex: confirmIndex}, true
				}
			}
		}
		if breakout.Close < episode.Low && (candidate.Side == RangeDiscoverySideAll || candidate.Side == RangeDiscoverySideDown) {
			for confirmIndex := breakIndex + 1; confirmIndex < len(candles) && confirmIndex <= breakIndex+cfg.ConfirmationWindowBars; confirmIndex++ {
				confirm := candles[confirmIndex]
				if confirm.Close < episode.Low && confirm.Low < breakout.Low {
					return structuredCompressionSignalEvent{breakoutIndex: breakIndex, confirmationIndex: confirmIndex}, true
				}
			}
		}
	}
	return structuredCompressionSignalEvent{}, false
}

func newFuturesRangeUniverseStructuredCompressionSignalRow(candles []Candle, symbol string, candidate FuturesRangeUniverseStructuredCompressionCandidateConfig, episode rangeRegimeDurabilityEpisode, breakoutIndex int, confirmationIndex int, sequence int, cfg FuturesRangeUniverseStructuredCompressionBaselineConfig, btCfg BacktestConfig, splits []Split) FuturesRangeUniverseStructuredCompressionSignalRow {
	breakout := candles[breakoutIndex]
	confirmation := candles[confirmationIndex]
	width := episode.High - episode.Low
	side := Long
	if breakout.Close < episode.Low {
		side = Short
	}
	signalID := fmt.Sprintf("%s_%s_%06d", strings.ToLower(symbol), candidate.CandidateID, sequence)
	entryIndex := confirmationIndex + 1
	row := FuturesRangeUniverseStructuredCompressionSignalRow{
		SignalID:                signalID,
		Symbol:                  strings.ToUpper(symbol),
		CandidateID:             candidate.CandidateID,
		Timeframe:               candidate.Timeframe,
		Family:                  RangeUniverseFamilyStructuredCompressionBreak,
		Split:                   splitNameForCloseTime(confirmation.CloseTime, splits),
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
		BreakoutIndex:           breakoutIndex,
		BreakoutOpenTime:        breakout.OpenTime.UTC().Format(timeLayout),
		BreakoutCloseTime:       breakout.CloseTime.UTC().Format(timeLayout),
		BreakoutDelayBars:       breakoutIndex - episode.EndIndex,
		BreakoutOpen:            breakout.Open,
		BreakoutHigh:            breakout.High,
		BreakoutLow:             breakout.Low,
		BreakoutClose:           breakout.Close,
		ConfirmationIndex:       confirmationIndex,
		ConfirmationOpenTime:    confirmation.OpenTime.UTC().Format(timeLayout),
		ConfirmationCloseTime:   confirmation.CloseTime.UTC().Format(timeLayout),
		ConfirmationDelayBars:   confirmationIndex - breakoutIndex,
		ConfirmationOpen:        confirmation.Open,
		ConfirmationHigh:        confirmation.High,
		ConfirmationLow:         confirmation.Low,
		ConfirmationClose:       confirmation.Close,
		Side:                    side,
		EntryIndex:              entryIndex,
		MaxHoldBars:             candidate.MaxHoldBars,
	}
	if width <= 0 || !finitePositive(width) {
		row.SkippedReason = "non_positive_range_width"
		return row
	}
	if entryIndex >= len(candles) {
		row.SkippedReason = "missing_entry_candle"
		return row
	}
	entry := candles[entryIndex]
	targetMultiple := candidate.TargetRangeWidthMultiple
	if targetMultiple == 0 {
		targetMultiple = 1
	}
	stopBuffer := candidate.StopBoundaryBufferRangeWidth
	row.EntryOpenTime = entry.OpenTime.UTC().Format(timeLayout)
	row.EntryOpen = entry.Open
	row.ExpectedEntryPrice = applySlippage(entry.Open, btCfg.SlippagePct, side, true)
	if side == Long {
		row.Stop = episode.High - width*stopBuffer
		row.Target = row.ExpectedEntryPrice + width*targetMultiple
	} else {
		row.Stop = episode.Low + width*stopBuffer
		row.Target = row.ExpectedEntryPrice - width*targetMultiple
	}
	if !finitePositive(row.Stop) || !finitePositive(row.Target) || !finitePositive(row.ExpectedEntryPrice) {
		row.SkippedReason = "non_positive_trade_price"
		return row
	}
	row.EntryGeometryValid = validEntryGeometry(Signal{Side: row.Side, Stop: row.Stop, Target: row.Target}, row.ExpectedEntryPrice)
	if !row.EntryGeometryValid {
		row.SkippedReason = "invalid_entry_geometry"
	}
	return row
}

func SummarizeFuturesRangeUniverseStructuredCompressionBaseline(signals []FuturesRangeUniverseStructuredCompressionSignalRow, trades []FuturesRangeUniverseStructuredCompressionTradeRow, cfg FuturesRangeUniverseStructuredCompressionBaselineConfig, startBalance float64, splits []Split) []FuturesRangeUniverseStructuredCompressionSummaryRow {
	cfg = cfg.withDefaults()
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	symbols := structuredCompressionSummarySymbols(cfg)
	rows := make([]FuturesRangeUniverseStructuredCompressionSummaryRow, 0, len(cfg.Candidates)*len(symbols)*len(splits)*3)
	for _, candidate := range cfg.Candidates {
		for _, symbol := range symbols {
			for _, split := range splits {
				for _, side := range []string{"all", string(Long), string(Short)} {
					filteredTrades := filterStructuredCompressionTrades(trades, candidate.CandidateID, symbol, split, side)
					filteredSignals := filterStructuredCompressionSignals(signals, candidate.CandidateID, symbol, split, side)
					row := summarizeStructuredCompressionTrades(filteredTrades, startBalance)
					row.CandidateID = candidate.CandidateID
					row.Symbol = symbol
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
	}
	markStructuredCompressionReviewFlags(rows, cfg, startBalance, splits)
	return rows
}

func FuturesRangeUniverseStructuredCompressionBaselineStopState(summary []FuturesRangeUniverseStructuredCompressionSummaryRow, cfg FuturesRangeUniverseStructuredCompressionBaselineConfig, startBalance float64, splits []Split) string {
	cfg = cfg.withDefaults()
	if len(summary) == 0 {
		return StructuredCompressionStopStateFailedNoPromotion
	}
	passing := 0
	nearViable := 0
	for _, candidate := range cfg.Candidates {
		if structuredCompressionCandidatePasses(summary, candidate, cfg, startBalance, splits) {
			passing++
		}
		if structuredCompressionCandidateNearViable(summary, candidate, cfg, startBalance, splits) {
			nearViable++
		}
	}
	if passing > 0 {
		return StructuredCompressionStopStatePassedNeedsOptimize
	}
	if nearViable >= 2 {
		return StructuredCompressionStopStateMixedPortfolioReview
	}
	return StructuredCompressionStopStateFailedNoPromotion
}

func (cfg FuturesRangeUniverseStructuredCompressionBaselineConfig) withDefaults() FuturesRangeUniverseStructuredCompressionBaselineConfig {
	defaults := DefaultFuturesRangeUniverseStructuredCompressionBaselineConfig()
	if len(cfg.Sources) == 0 {
		cfg.Sources = defaults.Sources
	}
	if cfg.EventDelayBars == 0 {
		cfg.EventDelayBars = defaults.EventDelayBars
	}
	if cfg.ConfirmationWindowBars == 0 {
		cfg.ConfirmationWindowBars = defaults.ConfirmationWindowBars
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
		cfg.Candidates = append([]FuturesRangeUniverseStructuredCompressionCandidateConfig(nil), defaults.Candidates...)
	}
	return cfg
}

func (cfg FuturesRangeUniverseStructuredCompressionBaselineConfig) validate() error {
	if len(cfg.Sources) == 0 {
		return fmt.Errorf("structured compression baseline source list cannot be empty")
	}
	if cfg.EventDelayBars <= 0 {
		return fmt.Errorf("structured compression event delay bars must be positive")
	}
	if cfg.ConfirmationWindowBars <= 0 {
		return fmt.Errorf("structured compression confirmation window bars must be positive")
	}
	if cfg.DetectorLookbackDays <= 0 {
		return fmt.Errorf("structured compression detector lookback days must be positive")
	}
	if cfg.DetectorPercentile <= 0 || cfg.DetectorPercentile >= 1 {
		return fmt.Errorf("structured compression detector percentile must be between 0 and 1")
	}
	if cfg.DetectorMinConsecutiveBars <= 0 {
		return fmt.Errorf("structured compression detector min consecutive bars must be positive")
	}
	if cfg.MinFullTrades <= 0 {
		return fmt.Errorf("structured compression min full trades must be positive")
	}
	if cfg.MinKeySplitTrades <= 0 {
		return fmt.Errorf("structured compression min key split trades must be positive")
	}
	if cfg.NearFlatNetPct <= 0 {
		return fmt.Errorf("structured compression near-flat net pct must be positive")
	}
	seen := map[string]bool{}
	for _, source := range cfg.Sources {
		symbol := strings.ToUpper(strings.TrimSpace(source.Symbol))
		if !rangeUniverseApprovedSymbol(symbol) {
			return fmt.Errorf("structured compression source symbol %q is not approved", symbol)
		}
		if seen[symbol] {
			return fmt.Errorf("structured compression source symbol %q is duplicated", symbol)
		}
		seen[symbol] = true
	}
	if len(cfg.Candidates) == 0 {
		return fmt.Errorf("structured compression requires at least one candidate")
	}
	for _, candidate := range cfg.Candidates {
		if err := cfg.validateCandidate(candidate); err != nil {
			return err
		}
	}
	return nil
}

func (cfg FuturesRangeUniverseStructuredCompressionBaselineConfig) validateCandidate(candidate FuturesRangeUniverseStructuredCompressionCandidateConfig) error {
	if candidate.CandidateID == "" {
		return fmt.Errorf("structured compression candidate id must not be empty")
	}
	if candidate.CandidateID != StructuredCompressionCandidate4HAllH6 && candidate.CandidateID != StructuredCompressionCandidate1HAllH12 {
		return fmt.Errorf("structured compression candidate %q is not approved for this baseline", candidate.CandidateID)
	}
	if candidate.Timeframe != RangeDiscoveryTimeframe1h && candidate.Timeframe != RangeDiscoveryTimeframe4h {
		return fmt.Errorf("structured compression candidate %s unsupported timeframe %q", candidate.CandidateID, candidate.Timeframe)
	}
	if candidate.Side != RangeDiscoverySideAll {
		return fmt.Errorf("structured compression candidate %s must use all sides", candidate.CandidateID)
	}
	if candidate.MaxHoldBars <= 0 {
		return fmt.Errorf("structured compression candidate %s max hold bars must be positive", candidate.CandidateID)
	}
	if candidate.TargetRangeWidthMultiple < 0 {
		return fmt.Errorf("structured compression candidate %s target range-width multiple must be non-negative", candidate.CandidateID)
	}
	if candidate.StopBoundaryBufferRangeWidth < 0 {
		return fmt.Errorf("structured compression candidate %s stop boundary buffer must be non-negative", candidate.CandidateID)
	}
	return nil
}

func structuredCompressionDetectorConfig(cfg FuturesRangeUniverseStructuredCompressionBaselineConfig, barsPerDay int) RangeDetectorConfig {
	detectorCfg := DefaultCompressionRangeDetectorConfig()
	detectorCfg.BarsPerDay = barsPerDay
	detectorCfg.LookbackDays = cfg.DetectorLookbackDays
	detectorCfg.Percentile = cfg.DetectorPercentile
	detectorCfg.MinConsecutiveBars = cfg.DetectorMinConsecutiveBars
	detectorCfg.UseBollinger = true
	detectorCfg.UseADX = false
	return detectorCfg
}

func structuredCompressionFrameDef(timeframe string) (rangeDiscoveryFrameDef, bool) {
	for _, frame := range rangeUniverseFrameDefs() {
		if frame.timeframe == timeframe {
			return frame, true
		}
	}
	return rangeDiscoveryFrameDef{}, false
}

func structuredCompressionBarsPerDay(timeframe string) int {
	frame, ok := structuredCompressionFrameDef(timeframe)
	if !ok {
		return 24
	}
	return frame.barsPerDay
}

func filterStructuredCompressionSignals(signals []FuturesRangeUniverseStructuredCompressionSignalRow, candidateID string, symbol string, split Split, side string) []FuturesRangeUniverseStructuredCompressionSignalRow {
	out := make([]FuturesRangeUniverseStructuredCompressionSignalRow, 0, len(signals))
	for _, signal := range signals {
		if signal.CandidateID != candidateID {
			continue
		}
		if symbol != StructuredCompressionSummaryAggregateSymbol && signal.Symbol != symbol {
			continue
		}
		eventTime, err := parseTime(signal.ConfirmationCloseTime)
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

func filterStructuredCompressionTrades(trades []FuturesRangeUniverseStructuredCompressionTradeRow, candidateID string, symbol string, split Split, side string) []FuturesRangeUniverseStructuredCompressionTradeRow {
	out := make([]FuturesRangeUniverseStructuredCompressionTradeRow, 0, len(trades))
	for _, trade := range trades {
		if trade.CandidateID != candidateID {
			continue
		}
		if symbol != StructuredCompressionSummaryAggregateSymbol && trade.Symbol != symbol {
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

func summarizeStructuredCompressionTrades(trades []FuturesRangeUniverseStructuredCompressionTradeRow, startBalance float64) FuturesRangeUniverseStructuredCompressionSummaryRow {
	row := FuturesRangeUniverseStructuredCompressionSummaryRow{TotalTrades: len(trades)}
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

func markStructuredCompressionReviewFlags(rows []FuturesRangeUniverseStructuredCompressionSummaryRow, cfg FuturesRangeUniverseStructuredCompressionBaselineConfig, startBalance float64, splits []Split) {
	for _, candidate := range cfg.Candidates {
		passes := structuredCompressionCandidatePasses(rows, candidate, cfg, startBalance, splits)
		near := structuredCompressionCandidateNearViable(rows, candidate, cfg, startBalance, splits)
		for _, symbol := range structuredCompressionSummarySymbols(cfg) {
			worstIndex := -1
			for i := range rows {
				if rows[i].CandidateID != candidate.CandidateID || rows[i].Symbol != symbol {
					continue
				}
				if rows[i].Split != fullSplitName && rows[i].Side == "all" {
					if worstIndex == -1 || rows[i].NetPnL < rows[worstIndex].NetPnL {
						worstIndex = i
					}
				}
				if rows[i].Symbol == StructuredCompressionSummaryAggregateSymbol && rows[i].Side == "all" {
					rows[i].PassesOptimizationGate = passes
					rows[i].NearViableForPortfolio = near
				}
			}
			if worstIndex >= 0 {
				rows[worstIndex].IsWorstPeriodSplit = true
			}
		}
	}
}

func structuredCompressionCandidatePasses(rows []FuturesRangeUniverseStructuredCompressionSummaryRow, candidate FuturesRangeUniverseStructuredCompressionCandidateConfig, cfg FuturesRangeUniverseStructuredCompressionBaselineConfig, startBalance float64, splits []Split) bool {
	byKey := structuredCompressionSummaryByKey(rows)
	aggregateFull := byKey[structuredCompressionSummaryKey(candidate.CandidateID, StructuredCompressionSummaryAggregateSymbol, fullSplitName, "all")]
	if aggregateFull.TotalTrades < cfg.MinFullTrades || aggregateFull.NetPnL <= 0 || aggregateFull.ProfitFactor < 1.2 {
		return false
	}
	btcPass := structuredCompressionSymbolPasses(byKey, candidate.CandidateID, RangeUniverseSymbolBTCUSDT, cfg, startBalance, splits)
	transferPasses := 0
	symbolPositiveFull := 0
	for _, source := range cfg.Sources {
		symbol := strings.ToUpper(source.Symbol)
		full := byKey[structuredCompressionSummaryKey(candidate.CandidateID, symbol, fullSplitName, "all")]
		if full.TotalTrades >= cfg.MinKeySplitTrades && full.NetPnL > 0 && full.ProfitFactor >= 1.0 {
			symbolPositiveFull++
		}
		if symbol == RangeUniverseSymbolBTCUSDT {
			continue
		}
		if structuredCompressionSymbolPasses(byKey, candidate.CandidateID, symbol, cfg, startBalance, splits) {
			transferPasses++
		}
	}
	if btcPass && transferPasses >= 1 {
		return true
	}
	if symbolPositiveFull >= 2 && structuredCompressionAggregateSplitsPass(byKey, candidate.CandidateID, cfg, startBalance, splits) {
		return true
	}
	return false
}

func structuredCompressionSymbolPasses(byKey map[string]FuturesRangeUniverseStructuredCompressionSummaryRow, candidateID string, symbol string, cfg FuturesRangeUniverseStructuredCompressionBaselineConfig, startBalance float64, splits []Split) bool {
	full := byKey[structuredCompressionSummaryKey(candidateID, symbol, fullSplitName, "all")]
	if full.TotalTrades < cfg.MinFullTrades || full.NetPnL <= 0 || full.ProfitFactor < 1.2 {
		return false
	}
	for _, splitName := range []string{"2023_2024_oos", "2025_2026_recent"} {
		row := byKey[structuredCompressionSummaryKey(candidateID, symbol, splitName, "all")]
		if row.TotalTrades < cfg.MinKeySplitTrades || row.NetPnL < 0 {
			return false
		}
	}
	for _, side := range []string{string(Long), string(Short)} {
		row := byKey[structuredCompressionSummaryKey(candidateID, symbol, fullSplitName, side)]
		if row.TotalTrades >= cfg.MinKeySplitTrades && row.NetPnL < -startBalance*cfg.NearFlatNetPct {
			return false
		}
	}
	return true
}

func structuredCompressionAggregateSplitsPass(byKey map[string]FuturesRangeUniverseStructuredCompressionSummaryRow, candidateID string, cfg FuturesRangeUniverseStructuredCompressionBaselineConfig, startBalance float64, splits []Split) bool {
	for _, splitName := range []string{"2023_2024_oos", "2025_2026_recent"} {
		row := byKey[structuredCompressionSummaryKey(candidateID, StructuredCompressionSummaryAggregateSymbol, splitName, "all")]
		if row.TotalTrades < cfg.MinKeySplitTrades || row.NetPnL < 0 {
			return false
		}
	}
	for _, side := range []string{string(Long), string(Short)} {
		row := byKey[structuredCompressionSummaryKey(candidateID, StructuredCompressionSummaryAggregateSymbol, fullSplitName, side)]
		if row.TotalTrades >= cfg.MinKeySplitTrades && row.NetPnL < -startBalance*cfg.NearFlatNetPct*float64(len(cfg.Sources)) {
			return false
		}
	}
	return true
}

func structuredCompressionCandidateNearViable(rows []FuturesRangeUniverseStructuredCompressionSummaryRow, candidate FuturesRangeUniverseStructuredCompressionCandidateConfig, cfg FuturesRangeUniverseStructuredCompressionBaselineConfig, startBalance float64, splits []Split) bool {
	byKey := structuredCompressionSummaryByKey(rows)
	full := byKey[structuredCompressionSummaryKey(candidate.CandidateID, StructuredCompressionSummaryAggregateSymbol, fullSplitName, "all")]
	if full.TotalTrades < cfg.MinFullTrades || full.GrossPnL <= 0 || full.ProfitFactor < 1.0 {
		return false
	}
	if full.NetPnL < -startBalance*cfg.NearFlatNetPct*float64(len(cfg.Sources)) {
		return false
	}
	for _, splitName := range []string{"2023_2024_oos", "2025_2026_recent"} {
		row := byKey[structuredCompressionSummaryKey(candidate.CandidateID, StructuredCompressionSummaryAggregateSymbol, splitName, "all")]
		if row.TotalTrades < cfg.MinKeySplitTrades || !structuredCompressionSplitNonNegativeOrNearFlat(row, startBalance*float64(len(cfg.Sources)), cfg.NearFlatNetPct) {
			return false
		}
	}
	return true
}

func structuredCompressionSplitNonNegativeOrNearFlat(row FuturesRangeUniverseStructuredCompressionSummaryRow, startBalance float64, nearFlatPct float64) bool {
	if row.TotalTrades == 0 {
		return false
	}
	if row.NetPnL >= 0 {
		return true
	}
	return row.NetPnL >= -startBalance*nearFlatPct && row.GrossPnL > 0
}

func structuredCompressionSummarySymbols(cfg FuturesRangeUniverseStructuredCompressionBaselineConfig) []string {
	symbols := []string{}
	for _, source := range cfg.Sources {
		symbols = append(symbols, strings.ToUpper(source.Symbol))
	}
	sort.Slice(symbols, func(i, j int) bool {
		return rangeUniverseSymbolSortKey(symbols[i]) < rangeUniverseSymbolSortKey(symbols[j])
	})
	symbols = append(symbols, StructuredCompressionSummaryAggregateSymbol)
	return symbols
}

func structuredCompressionSummaryKey(candidateID, symbol, split, side string) string {
	return candidateID + "|" + symbol + "|" + split + "|" + side
}

func structuredCompressionSummaryByKey(rows []FuturesRangeUniverseStructuredCompressionSummaryRow) map[string]FuturesRangeUniverseStructuredCompressionSummaryRow {
	byKey := map[string]FuturesRangeUniverseStructuredCompressionSummaryRow{}
	for _, row := range rows {
		byKey[structuredCompressionSummaryKey(row.CandidateID, row.Symbol, row.Split, row.Side)] = row
	}
	return byKey
}

func structuredCompressionTimeframeSortKey(timeframe string) int {
	switch timeframe {
	case RangeDiscoveryTimeframe1h:
		return 1
	case RangeDiscoveryTimeframe4h:
		return 2
	default:
		return 99
	}
}

func finitePositive(v float64) bool {
	return v > 0 && !math.IsNaN(v) && !math.IsInf(v, 0)
}
