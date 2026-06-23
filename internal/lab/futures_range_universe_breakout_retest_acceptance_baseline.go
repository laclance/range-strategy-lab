package lab

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

const (
	FuturesRangeUniverseBreakoutRetestAcceptanceBaselineName = "futures_range_universe_breakout_retest_acceptance_baseline"

	BreakoutRetestAcceptanceCandidate15MAllH12 = "breakout_retest_acceptance_15m_all_h12"
	BreakoutRetestAcceptanceCandidate1HAllH12  = "breakout_retest_acceptance_1h_all_h12"

	BreakoutRetestAcceptanceStopStateSourceGap           = "breakout_retest_acceptance_baseline_source_gap"
	BreakoutRetestAcceptanceStopStateNoRankedCandidate   = "breakout_retest_acceptance_baseline_no_ranked_candidate"
	BreakoutRetestAcceptanceStopStateFailedNoPromotion   = "breakout_retest_acceptance_baseline_failed_no_promotion"
	BreakoutRetestAcceptanceStopStatePassedNeedsRobust   = "breakout_retest_acceptance_baseline_passed_needs_robustness_review"
	BreakoutRetestAcceptanceSummaryAggregateSymbol       = "all"
	BreakoutRetestAcceptanceDefaultTargetWidthMultiple   = 1.0
	BreakoutRetestAcceptanceDefaultStopBoundaryBuffer    = 0.0
	BreakoutRetestAcceptanceDefaultMaxSelectedCandidates = 2
)

type FuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig struct {
	DiscoveryConfig              FuturesRangeUniverseDiscoveryAuditConfig
	MaxSelectedCandidates        int
	TargetRangeWidthMultiple     float64
	StopBoundaryBufferRangeWidth float64
	MinFullTrades                int
	MinKeySplitTrades            int
	NearFlatNetPct               float64
}

type FuturesRangeUniverseBreakoutRetestAcceptanceCandidateConfig struct {
	CandidateID                  string
	SelectedOrder                int
	Timeframe                    string
	Side                         string
	MaxHoldBars                  int
	TargetRangeWidthMultiple     float64
	StopBoundaryBufferRangeWidth float64
	SourceRank                   int
	RankScore                    float64
}

type FuturesRangeUniverseBreakoutRetestAcceptanceBaselineResult struct {
	SourceRows    []FuturesRangeUniverseSourceRow
	CoverageRows  []FuturesRangeUniverseCoverageRow
	SelectionRows []FuturesRangeUniverseBreakoutRetestAcceptanceSelectionRow
	SignalRows    []FuturesRangeUniverseBreakoutRetestAcceptanceSignalRow
	TradeRows     []FuturesRangeUniverseBreakoutRetestAcceptanceTradeRow
	SummaryRows   []FuturesRangeUniverseBreakoutRetestAcceptanceSummaryRow
	Trades        []Trade
}

type FuturesRangeUniverseBreakoutRetestAcceptanceSelectionRow struct {
	SelectedOrder                        int     `json:"selected_order"`
	CandidateID                          string  `json:"candidate_id"`
	SourceRank                           int     `json:"source_rank"`
	Timeframe                            string  `json:"timeframe"`
	Family                               string  `json:"family"`
	Side                                 string  `json:"side"`
	HorizonBars                          int     `json:"horizon_bars"`
	PassesGate                           bool    `json:"passes_gate"`
	RankScore                            float64 `json:"rank_score"`
	BTCUSDTGatePass                      bool    `json:"btcusdt_gate_pass"`
	TransferSymbolGatePassCount          int     `json:"transfer_symbol_gate_pass_count"`
	SymbolGatePassCount                  int     `json:"symbol_gate_pass_count"`
	SymbolSpecificException              bool    `json:"symbol_specific_exception"`
	SymbolSpecificBestSymbol             string  `json:"symbol_specific_best_symbol,omitempty"`
	BTCUSDTWeakestSplitCandidateCount    int     `json:"btcusdt_weakest_split_candidate_count"`
	BTCUSDTWeakestFavorableRate          float64 `json:"btcusdt_weakest_favorable_rate"`
	BTCUSDTWorstAdverseRate              float64 `json:"btcusdt_worst_adverse_rate"`
	BTCUSDTWorstQuickInvalidationRate    float64 `json:"btcusdt_worst_quick_invalidation_rate"`
	BTCUSDTWeakestFavorableMinusAdverse  float64 `json:"btcusdt_weakest_favorable_minus_adverse_pct"`
	BTCUSDTWeakestCostBufferPct          float64 `json:"btcusdt_weakest_cost_buffer_pct"`
	BestTransferSymbol                   string  `json:"best_transfer_symbol,omitempty"`
	BestTransferWeakestSplitCount        int     `json:"best_transfer_weakest_split_candidate_count"`
	BestTransferWeakestFavorableRate     float64 `json:"best_transfer_weakest_favorable_rate"`
	BestTransferWorstAdverseRate         float64 `json:"best_transfer_worst_adverse_rate"`
	BestTransferWorstQuickInvalidation   float64 `json:"best_transfer_worst_quick_invalidation_rate"`
	BestTransferWeakestFavorableMinusAdv float64 `json:"best_transfer_weakest_favorable_minus_adverse_pct"`
	BestTransferWeakestCostBufferPct     float64 `json:"best_transfer_weakest_cost_buffer_pct"`
	CombinedFullCandidateCount           int     `json:"combined_full_candidate_count"`
	CombinedWeakestSplitCandidateCount   int     `json:"combined_weakest_split_candidate_count"`
	CombinedWeakestFavorableRate         float64 `json:"combined_weakest_favorable_rate"`
	CombinedWorstAdverseRate             float64 `json:"combined_worst_adverse_rate"`
	CombinedWorstQuickInvalidationRate   float64 `json:"combined_worst_quick_invalidation_rate"`
	CombinedWeakestFavorableMinusAdverse float64 `json:"combined_weakest_favorable_minus_adverse_pct"`
	CombinedWeakestCostBufferPct         float64 `json:"combined_weakest_cost_buffer_pct"`
	PeriodSplitsRequired                 int     `json:"period_splits_required"`
	MinCandidatesPerSplit                int     `json:"min_candidates_per_split"`
	SymbolSpecificMinCandidatesPerSplit  int     `json:"symbol_specific_min_candidates_per_split"`
	SymbolSpecificCostBufferMultiple     float64 `json:"symbol_specific_cost_buffer_multiple"`
}

type FuturesRangeUniverseBreakoutRetestAcceptanceStrategy struct {
	symbol      string
	candidate   FuturesRangeUniverseBreakoutRetestAcceptanceCandidateConfig
	signals     []FuturesRangeUniverseBreakoutRetestAcceptanceSignalRow
	signalsByID map[string]FuturesRangeUniverseBreakoutRetestAcceptanceSignalRow
	byIndex     map[int]Signal
}

type FuturesRangeUniverseBreakoutRetestAcceptanceSignalRow struct {
	SignalID                string    `json:"signal_id"`
	Symbol                  string    `json:"symbol"`
	CandidateID             string    `json:"candidate_id"`
	SelectedOrder           int       `json:"selected_order"`
	SourceRank              int       `json:"source_rank"`
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
	RetestIndex             int       `json:"retest_index"`
	RetestOpenTime          string    `json:"retest_open_time"`
	RetestCloseTime         string    `json:"retest_close_time"`
	RetestDelayBars         int       `json:"retest_delay_bars"`
	ReentryDelayBars        int       `json:"reentry_delay_bars"`
	RetestOpen              float64   `json:"retest_open"`
	RetestHigh              float64   `json:"retest_high"`
	RetestLow               float64   `json:"retest_low"`
	RetestClose             float64   `json:"retest_close"`
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

type FuturesRangeUniverseBreakoutRetestAcceptanceTradeRow struct {
	SignalID          string    `json:"signal_id"`
	Symbol            string    `json:"symbol"`
	CandidateID       string    `json:"candidate_id"`
	SelectedOrder     int       `json:"selected_order"`
	SourceRank        int       `json:"source_rank"`
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
	RetestIndex       int       `json:"retest_index"`
	RetestCloseTime   string    `json:"retest_close_time"`
	RetestDelayBars   int       `json:"retest_delay_bars"`
	ReentryDelayBars  int       `json:"reentry_delay_bars"`
	RetestClose       float64   `json:"retest_close"`
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

type FuturesRangeUniverseBreakoutRetestAcceptanceSummaryRow struct {
	CandidateID        string  `json:"candidate_id"`
	Symbol             string  `json:"symbol"`
	Timeframe          string  `json:"timeframe"`
	Split              string  `json:"split"`
	Side               string  `json:"side"`
	SignalCount        int     `json:"signal_count"`
	SkippedSignalCount int     `json:"skipped_signal_count"`
	TotalTrades        int     `json:"total_trades"`
	Wins               int     `json:"wins"`
	Losses             int     `json:"losses"`
	WinRate            float64 `json:"win_rate"`
	GrossPnL           float64 `json:"gross_pnl"`
	NetPnL             float64 `json:"net_pnl"`
	TotalCosts         float64 `json:"total_costs"`
	ProfitFactor       float64 `json:"profit_factor"`
	GrossProfitFactor  float64 `json:"gross_profit_factor"`
	MaxDrawdown        float64 `json:"max_drawdown"`
	AvgGrossR          float64 `json:"avg_gross_r"`
	AvgNetR            float64 `json:"avg_net_r"`
	AvgInitialRisk     float64 `json:"avg_initial_risk"`
	AvgHoldBars        float64 `json:"avg_hold_bars"`
	PassesReviewGate   bool    `json:"passes_review_gate"`
	IsWorstPeriodSplit bool    `json:"is_worst_period_split"`
}

func DefaultFuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig() FuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig {
	return FuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig{
		DiscoveryConfig:              DefaultFuturesRangeUniverseDiscoveryAuditConfig(),
		MaxSelectedCandidates:        BreakoutRetestAcceptanceDefaultMaxSelectedCandidates,
		TargetRangeWidthMultiple:     BreakoutRetestAcceptanceDefaultTargetWidthMultiple,
		StopBoundaryBufferRangeWidth: BreakoutRetestAcceptanceDefaultStopBoundaryBuffer,
		MinFullTrades:                100,
		MinKeySplitTrades:            25,
		NearFlatNetPct:               0.005,
	}
}

func RunFuturesRangeUniverseBreakoutRetestAcceptanceBaselineBacktest(cfg FuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig, btCfg BacktestConfig, splits []Split) (FuturesRangeUniverseBreakoutRetestAcceptanceBaselineResult, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return FuturesRangeUniverseBreakoutRetestAcceptanceBaselineResult{}, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}

	discoveryResult, err := RunFuturesRangeUniverseDiscoveryAudit(cfg.DiscoveryConfig, splits)
	result := FuturesRangeUniverseBreakoutRetestAcceptanceBaselineResult{SourceRows: discoveryResult.SourceRows}
	if err != nil {
		return result, err
	}
	result.SelectionRows = SelectFuturesRangeUniverseBreakoutRetestAcceptanceCandidates(discoveryResult.RankingRows, cfg)
	if len(result.SelectionRows) == 0 {
		return result, nil
	}
	candidates := breakoutRetestAcceptanceCandidatesFromSelection(result.SelectionRows, cfg)

	for _, source := range cfg.DiscoveryConfig.Sources {
		candles, sourceRow, err := LoadFuturesRangeUniverseSource(source, splits)
		if err != nil {
			return result, err
		}

		frameCandles := map[string][]Candle{}
		coverageByFrame := map[string]FuturesRangeUniverseCoverageRow{}
		for _, candidate := range candidates {
			if _, ok := frameCandles[candidate.Timeframe]; ok {
				continue
			}
			frame, ok := breakoutRetestAcceptanceFrameDef(candidate.Timeframe)
			if !ok {
				return result, fmt.Errorf("unsupported breakout retest acceptance timeframe %q", candidate.Timeframe)
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

		for _, candidate := range candidates {
			coverage := coverageByFrame[candidate.Timeframe]
			if !coverage.Complete || coverage.ValidationStatus != "accepted" {
				return result, fmt.Errorf("%s %s breakout retest acceptance resample rejected: %s", sourceRow.Symbol, candidate.Timeframe, coverage.ValidationError)
			}
			candlesForFrame := frameCandles[candidate.Timeframe]
			if len(candlesForFrame) == 0 {
				continue
			}
			strategy, err := NewFuturesRangeUniverseBreakoutRetestAcceptanceStrategy(candlesForFrame, sourceRow.Symbol, candidate, cfg, btCfg, splits)
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
		return breakoutRetestAcceptanceTimeframeSortKey(result.CoverageRows[i].Timeframe) < breakoutRetestAcceptanceTimeframeSortKey(result.CoverageRows[j].Timeframe)
	})
	sort.Slice(result.SignalRows, func(i, j int) bool {
		if result.SignalRows[i].Symbol != result.SignalRows[j].Symbol {
			return rangeUniverseSymbolSortKey(result.SignalRows[i].Symbol) < rangeUniverseSymbolSortKey(result.SignalRows[j].Symbol)
		}
		if result.SignalRows[i].SelectedOrder != result.SignalRows[j].SelectedOrder {
			return result.SignalRows[i].SelectedOrder < result.SignalRows[j].SelectedOrder
		}
		if result.SignalRows[i].RetestIndex != result.SignalRows[j].RetestIndex {
			return result.SignalRows[i].RetestIndex < result.SignalRows[j].RetestIndex
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
	result.SummaryRows = SummarizeFuturesRangeUniverseBreakoutRetestAcceptanceBaseline(result.SignalRows, result.TradeRows, cfg, btCfg.StartBalance, splits)
	return result, nil
}

func SelectFuturesRangeUniverseBreakoutRetestAcceptanceCandidates(rankings []FuturesRangeUniverseRankingRow, cfg FuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig) []FuturesRangeUniverseBreakoutRetestAcceptanceSelectionRow {
	cfg = cfg.withDefaults()
	rows := append([]FuturesRangeUniverseRankingRow(nil), rankings...)
	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].Rank > 0 && rows[j].Rank > 0 && rows[i].Rank != rows[j].Rank {
			return rows[i].Rank < rows[j].Rank
		}
		if rows[i].PassesGate != rows[j].PassesGate {
			return rows[i].PassesGate
		}
		if rows[i].RankScore != rows[j].RankScore {
			return rows[i].RankScore > rows[j].RankScore
		}
		if rows[i].Timeframe != rows[j].Timeframe {
			return breakoutRetestAcceptanceTimeframeSortKey(rows[i].Timeframe) < breakoutRetestAcceptanceTimeframeSortKey(rows[j].Timeframe)
		}
		return rows[i].HorizonBars > rows[j].HorizonBars
	})

	selected := []FuturesRangeUniverseBreakoutRetestAcceptanceSelectionRow{}
	seen := map[string]bool{}
	for _, row := range rows {
		if !row.PassesGate || row.Family != RangeUniverseFamilyBreakoutRetestAcceptance {
			continue
		}
		if row.Side != RangeDiscoverySideAll {
			continue
		}
		if row.Timeframe == "" || row.HorizonBars <= 0 {
			continue
		}
		key := row.Family + "|" + row.Timeframe + "|" + row.Side
		if seen[key] {
			continue
		}
		seen[key] = true
		selected = append(selected, breakoutRetestAcceptanceSelectionRow(row, len(selected)+1))
		if len(selected) >= cfg.MaxSelectedCandidates {
			break
		}
	}
	return selected
}

func NewFuturesRangeUniverseBreakoutRetestAcceptanceStrategy(candles []Candle, symbol string, candidate FuturesRangeUniverseBreakoutRetestAcceptanceCandidateConfig, cfg FuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig, btCfg BacktestConfig, splits []Split) (FuturesRangeUniverseBreakoutRetestAcceptanceStrategy, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validateCandidate(candidate); err != nil {
		return FuturesRangeUniverseBreakoutRetestAcceptanceStrategy{}, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	frame, ok := breakoutRetestAcceptanceFrameDef(candidate.Timeframe)
	if !ok {
		return FuturesRangeUniverseBreakoutRetestAcceptanceStrategy{}, fmt.Errorf("unsupported breakout retest acceptance timeframe %q", candidate.Timeframe)
	}
	detectorCfg := breakoutRetestAcceptanceDetectorConfig(cfg, frame.barsPerDay)
	classifications, err := (CompressionRangeDetector{Config: detectorCfg}).Classify(candles)
	if err != nil {
		return FuturesRangeUniverseBreakoutRetestAcceptanceStrategy{}, err
	}
	return newFuturesRangeUniverseBreakoutRetestAcceptanceStrategyFromClassifications(candles, symbol, candidate, cfg, btCfg, classifications, splits)
}

func newFuturesRangeUniverseBreakoutRetestAcceptanceStrategyFromClassifications(candles []Candle, symbol string, candidate FuturesRangeUniverseBreakoutRetestAcceptanceCandidateConfig, cfg FuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig, btCfg BacktestConfig, classifications []RangeClassification, splits []Split) (FuturesRangeUniverseBreakoutRetestAcceptanceStrategy, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validateCandidate(candidate); err != nil {
		return FuturesRangeUniverseBreakoutRetestAcceptanceStrategy{}, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	if len(classifications) > len(candles) {
		classifications = classifications[:len(candles)]
	}
	frame, ok := breakoutRetestAcceptanceFrameDef(candidate.Timeframe)
	if !ok {
		return FuturesRangeUniverseBreakoutRetestAcceptanceStrategy{}, fmt.Errorf("unsupported breakout retest acceptance timeframe %q", candidate.Timeframe)
	}
	detectorCfg := breakoutRetestAcceptanceDetectorConfig(cfg, frame.barsPerDay)
	normalizedATR := NormalizedATR(candles, detectorCfg.ATRPeriod)
	episodes := rangeRegimeDurabilityEpisodes(candles, classifications, normalizedATR, splits, BalancedDetectorProfileID)
	strategy := FuturesRangeUniverseBreakoutRetestAcceptanceStrategy{
		symbol:      strings.ToUpper(symbol),
		candidate:   candidate,
		signalsByID: map[string]FuturesRangeUniverseBreakoutRetestAcceptanceSignalRow{},
		byIndex:     map[int]Signal{},
	}
	nextEventID := 0
	for _, episode := range episodes {
		if episode.High <= episode.Low {
			continue
		}
		events := rangeUniverseBreakoutRetestEvents(candles, episode, frame, cfg.DiscoveryConfig.Discovery, &nextEventID)
		for _, event := range events {
			if !breakoutRetestAcceptanceCandidateAllowsSide(candidate, event.side) {
				continue
			}
			row := newFuturesRangeUniverseBreakoutRetestAcceptanceSignalRow(candles, strategy.symbol, candidate, event, len(strategy.signals)+1, cfg, btCfg, splits)
			if row.SignalID != "" && row.SkippedReason == "" {
				if _, exists := strategy.byIndex[row.RetestIndex]; exists {
					row.SkippedReason = "duplicate_retest_index"
				} else {
					strategy.byIndex[row.RetestIndex] = Signal{
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

func (s FuturesRangeUniverseBreakoutRetestAcceptanceStrategy) Name() string {
	return s.candidate.CandidateID
}

func (s FuturesRangeUniverseBreakoutRetestAcceptanceStrategy) OnCandle(ctx StrategyContext) (Signal, bool) {
	sig, ok := s.byIndex[ctx.Index]
	return sig, ok
}

func (s FuturesRangeUniverseBreakoutRetestAcceptanceStrategy) SignalRows() []FuturesRangeUniverseBreakoutRetestAcceptanceSignalRow {
	return append([]FuturesRangeUniverseBreakoutRetestAcceptanceSignalRow(nil), s.signals...)
}

func (s FuturesRangeUniverseBreakoutRetestAcceptanceStrategy) TradeRows(trades []Trade, splits []Split) []FuturesRangeUniverseBreakoutRetestAcceptanceTradeRow {
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	rows := make([]FuturesRangeUniverseBreakoutRetestAcceptanceTradeRow, 0, len(trades))
	for _, trade := range trades {
		signal, ok := s.signalsByID[trade.Signal]
		if !ok {
			continue
		}
		entryTime, _ := parseTime(trade.EntryTime)
		exitTime, _ := parseTime(trade.ExitTime)
		initialRisk := math.Abs(trade.EntryPrice-trade.Stop) * trade.Size
		row := FuturesRangeUniverseBreakoutRetestAcceptanceTradeRow{
			SignalID:          signal.SignalID,
			Symbol:            signal.Symbol,
			CandidateID:       signal.CandidateID,
			SelectedOrder:     signal.SelectedOrder,
			SourceRank:        signal.SourceRank,
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
			RetestIndex:       signal.RetestIndex,
			RetestCloseTime:   signal.RetestCloseTime,
			RetestDelayBars:   signal.RetestDelayBars,
			ReentryDelayBars:  signal.ReentryDelayBars,
			RetestClose:       signal.RetestClose,
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

func newFuturesRangeUniverseBreakoutRetestAcceptanceSignalRow(candles []Candle, symbol string, candidate FuturesRangeUniverseBreakoutRetestAcceptanceCandidateConfig, event rangeDiscoveryEvent, sequence int, cfg FuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig, btCfg BacktestConfig, splits []Split) FuturesRangeUniverseBreakoutRetestAcceptanceSignalRow {
	episode := event.episode
	retestIndex := event.eventIndex
	breakoutIndex := episode.EndIndex + event.breakoutDelayBars
	if breakoutIndex < 0 || breakoutIndex >= len(candles) || retestIndex < 0 || retestIndex >= len(candles) {
		return FuturesRangeUniverseBreakoutRetestAcceptanceSignalRow{
			SignalID:      fmt.Sprintf("%s_%s_%06d", strings.ToLower(symbol), candidate.CandidateID, sequence),
			Symbol:        strings.ToUpper(symbol),
			CandidateID:   candidate.CandidateID,
			SelectedOrder: candidateSelectedOrder(candidate),
			Timeframe:     candidate.Timeframe,
			Family:        RangeUniverseFamilyBreakoutRetestAcceptance,
			SkippedReason: "missing_event_candle",
		}
	}
	breakout := candles[breakoutIndex]
	retest := candles[retestIndex]
	width := episode.High - episode.Low
	side := Long
	if event.side == RangeDiscoverySideDown {
		side = Short
	}
	signalID := fmt.Sprintf("%s_%s_%06d", strings.ToLower(symbol), candidate.CandidateID, sequence)
	entryIndex := retestIndex + 1
	row := FuturesRangeUniverseBreakoutRetestAcceptanceSignalRow{
		SignalID:                signalID,
		Symbol:                  strings.ToUpper(symbol),
		CandidateID:             candidate.CandidateID,
		SelectedOrder:           candidateSelectedOrder(candidate),
		SourceRank:              candidate.SourceRank,
		Timeframe:               candidate.Timeframe,
		Family:                  RangeUniverseFamilyBreakoutRetestAcceptance,
		Split:                   splitNameForCloseTime(retest.CloseTime, splits),
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
		BreakoutDelayBars:       event.breakoutDelayBars,
		BreakoutOpen:            breakout.Open,
		BreakoutHigh:            breakout.High,
		BreakoutLow:             breakout.Low,
		BreakoutClose:           breakout.Close,
		RetestIndex:             retestIndex,
		RetestOpenTime:          retest.OpenTime.UTC().Format(timeLayout),
		RetestCloseTime:         retest.CloseTime.UTC().Format(timeLayout),
		RetestDelayBars:         event.eventDelayBars,
		ReentryDelayBars:        event.reentryDelayBars,
		RetestOpen:              retest.Open,
		RetestHigh:              retest.High,
		RetestLow:               retest.Low,
		RetestClose:             retest.Close,
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
		targetMultiple = BreakoutRetestAcceptanceDefaultTargetWidthMultiple
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

func SummarizeFuturesRangeUniverseBreakoutRetestAcceptanceBaseline(signals []FuturesRangeUniverseBreakoutRetestAcceptanceSignalRow, trades []FuturesRangeUniverseBreakoutRetestAcceptanceTradeRow, cfg FuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig, startBalance float64, splits []Split) []FuturesRangeUniverseBreakoutRetestAcceptanceSummaryRow {
	cfg = cfg.withDefaults()
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	candidates := breakoutRetestAcceptanceCandidatesFromSelection(SelectFuturesRangeUniverseBreakoutRetestAcceptanceCandidatesFromSignals(signals), cfg)
	if len(candidates) == 0 {
		candidates = breakoutRetestAcceptanceCandidatesFromSelectionFromSignals(signals, trades, cfg)
	}
	symbols := breakoutRetestAcceptanceSummarySymbols(cfg)
	rows := make([]FuturesRangeUniverseBreakoutRetestAcceptanceSummaryRow, 0, len(candidates)*len(symbols)*len(splits)*3)
	for _, candidate := range candidates {
		for _, symbol := range symbols {
			for _, split := range splits {
				for _, side := range []string{"all", string(Long), string(Short)} {
					filteredTrades := filterBreakoutRetestAcceptanceTrades(trades, candidate.CandidateID, symbol, split, side)
					filteredSignals := filterBreakoutRetestAcceptanceSignals(signals, candidate.CandidateID, symbol, split, side)
					row := summarizeBreakoutRetestAcceptanceTrades(filteredTrades, startBalance)
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
	markBreakoutRetestAcceptanceReviewFlags(rows, cfg, startBalance, splits)
	return rows
}

func FuturesRangeUniverseBreakoutRetestAcceptanceBaselineStopState(result FuturesRangeUniverseBreakoutRetestAcceptanceBaselineResult, cfg FuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig, startBalance float64, splits []Split) string {
	cfg = cfg.withDefaults()
	if len(result.SelectionRows) == 0 {
		return BreakoutRetestAcceptanceStopStateNoRankedCandidate
	}
	if len(result.SummaryRows) == 0 {
		return BreakoutRetestAcceptanceStopStateFailedNoPromotion
	}
	passing := 0
	for _, candidate := range breakoutRetestAcceptanceCandidatesFromSelection(result.SelectionRows, cfg) {
		if breakoutRetestAcceptanceCandidatePasses(result.SummaryRows, candidate, cfg, startBalance, splits) {
			passing++
		}
	}
	if passing > 0 {
		return BreakoutRetestAcceptanceStopStatePassedNeedsRobust
	}
	return BreakoutRetestAcceptanceStopStateFailedNoPromotion
}

func (cfg FuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig) withDefaults() FuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig {
	defaults := DefaultFuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig()
	cfg.DiscoveryConfig = cfg.DiscoveryConfig.withDefaults()
	if len(cfg.DiscoveryConfig.Sources) == 0 {
		cfg.DiscoveryConfig.Sources = defaults.DiscoveryConfig.Sources
	}
	if cfg.MaxSelectedCandidates == 0 {
		cfg.MaxSelectedCandidates = defaults.MaxSelectedCandidates
	}
	if cfg.TargetRangeWidthMultiple == 0 {
		cfg.TargetRangeWidthMultiple = defaults.TargetRangeWidthMultiple
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
	return cfg
}

func (cfg FuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig) validate() error {
	if err := cfg.DiscoveryConfig.validate(); err != nil {
		return err
	}
	if cfg.MaxSelectedCandidates <= 0 {
		return fmt.Errorf("breakout retest acceptance max selected candidates must be positive")
	}
	if cfg.TargetRangeWidthMultiple <= 0 {
		return fmt.Errorf("breakout retest acceptance target range-width multiple must be positive")
	}
	if cfg.StopBoundaryBufferRangeWidth < 0 {
		return fmt.Errorf("breakout retest acceptance stop boundary buffer must be non-negative")
	}
	if cfg.MinFullTrades <= 0 {
		return fmt.Errorf("breakout retest acceptance min full trades must be positive")
	}
	if cfg.MinKeySplitTrades <= 0 {
		return fmt.Errorf("breakout retest acceptance min key split trades must be positive")
	}
	if cfg.NearFlatNetPct <= 0 {
		return fmt.Errorf("breakout retest acceptance near-flat net pct must be positive")
	}
	return nil
}

func (cfg FuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig) validateCandidate(candidate FuturesRangeUniverseBreakoutRetestAcceptanceCandidateConfig) error {
	if candidate.CandidateID == "" {
		return fmt.Errorf("breakout retest acceptance candidate id must not be empty")
	}
	if candidate.Timeframe != RangeDiscoveryTimeframe15m && candidate.Timeframe != RangeDiscoveryTimeframe1h && candidate.Timeframe != RangeDiscoveryTimeframe4h && candidate.Timeframe != RangeDiscoveryTimeframe5m {
		return fmt.Errorf("breakout retest acceptance candidate %s unsupported timeframe %q", candidate.CandidateID, candidate.Timeframe)
	}
	if candidate.Side != RangeDiscoverySideAll && candidate.Side != RangeDiscoverySideUp && candidate.Side != RangeDiscoverySideDown {
		return fmt.Errorf("breakout retest acceptance candidate %s unsupported side %q", candidate.CandidateID, candidate.Side)
	}
	if candidate.MaxHoldBars <= 0 {
		return fmt.Errorf("breakout retest acceptance candidate %s max hold bars must be positive", candidate.CandidateID)
	}
	if candidate.TargetRangeWidthMultiple < 0 {
		return fmt.Errorf("breakout retest acceptance candidate %s target range-width multiple must be non-negative", candidate.CandidateID)
	}
	if candidate.StopBoundaryBufferRangeWidth < 0 {
		return fmt.Errorf("breakout retest acceptance candidate %s stop boundary buffer must be non-negative", candidate.CandidateID)
	}
	return nil
}

func breakoutRetestAcceptanceSelectionRow(row FuturesRangeUniverseRankingRow, selectedOrder int) FuturesRangeUniverseBreakoutRetestAcceptanceSelectionRow {
	return FuturesRangeUniverseBreakoutRetestAcceptanceSelectionRow{
		SelectedOrder:                        selectedOrder,
		CandidateID:                          breakoutRetestAcceptanceCandidateID(row.Timeframe, row.Side, row.HorizonBars),
		SourceRank:                           row.Rank,
		Timeframe:                            row.Timeframe,
		Family:                               row.Family,
		Side:                                 row.Side,
		HorizonBars:                          row.HorizonBars,
		PassesGate:                           row.PassesGate,
		RankScore:                            row.RankScore,
		BTCUSDTGatePass:                      row.BTCUSDTGatePass,
		TransferSymbolGatePassCount:          row.TransferSymbolGatePassCount,
		SymbolGatePassCount:                  row.SymbolGatePassCount,
		SymbolSpecificException:              row.SymbolSpecificException,
		SymbolSpecificBestSymbol:             row.SymbolSpecificBestSymbol,
		BTCUSDTWeakestSplitCandidateCount:    row.BTCUSDTWeakestSplitCandidateCount,
		BTCUSDTWeakestFavorableRate:          row.BTCUSDTWeakestFavorableRate,
		BTCUSDTWorstAdverseRate:              row.BTCUSDTWorstAdverseRate,
		BTCUSDTWorstQuickInvalidationRate:    row.BTCUSDTWorstQuickInvalidationRate,
		BTCUSDTWeakestFavorableMinusAdverse:  row.BTCUSDTWeakestFavorableMinusAdverse,
		BTCUSDTWeakestCostBufferPct:          row.BTCUSDTWeakestCostBufferPct,
		BestTransferSymbol:                   row.BestTransferSymbol,
		BestTransferWeakestSplitCount:        row.BestTransferWeakestSplitCount,
		BestTransferWeakestFavorableRate:     row.BestTransferWeakestFavorableRate,
		BestTransferWorstAdverseRate:         row.BestTransferWorstAdverseRate,
		BestTransferWorstQuickInvalidation:   row.BestTransferWorstQuickInvalidation,
		BestTransferWeakestFavorableMinusAdv: row.BestTransferWeakestFavorableMinusAdv,
		BestTransferWeakestCostBufferPct:     row.BestTransferWeakestCostBufferPct,
		CombinedFullCandidateCount:           row.CombinedFullCandidateCount,
		CombinedWeakestSplitCandidateCount:   row.CombinedWeakestSplitCandidateCount,
		CombinedWeakestFavorableRate:         row.CombinedWeakestFavorableRate,
		CombinedWorstAdverseRate:             row.CombinedWorstAdverseRate,
		CombinedWorstQuickInvalidationRate:   row.CombinedWorstQuickInvalidationRate,
		CombinedWeakestFavorableMinusAdverse: row.CombinedWeakestFavorableMinusAdverse,
		CombinedWeakestCostBufferPct:         row.CombinedWeakestCostBufferPct,
		PeriodSplitsRequired:                 row.PeriodSplitsRequired,
		MinCandidatesPerSplit:                row.MinCandidatesPerSplit,
		SymbolSpecificMinCandidatesPerSplit:  row.SymbolSpecificMinCandidatesPerSplit,
		SymbolSpecificCostBufferMultiple:     row.SymbolSpecificCostBufferMultiple,
	}
}

func breakoutRetestAcceptanceCandidatesFromSelection(rows []FuturesRangeUniverseBreakoutRetestAcceptanceSelectionRow, cfg FuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig) []FuturesRangeUniverseBreakoutRetestAcceptanceCandidateConfig {
	out := make([]FuturesRangeUniverseBreakoutRetestAcceptanceCandidateConfig, 0, len(rows))
	for _, row := range rows {
		out = append(out, FuturesRangeUniverseBreakoutRetestAcceptanceCandidateConfig{
			CandidateID:                  row.CandidateID,
			SelectedOrder:                row.SelectedOrder,
			Timeframe:                    row.Timeframe,
			Side:                         row.Side,
			MaxHoldBars:                  row.HorizonBars,
			TargetRangeWidthMultiple:     cfg.TargetRangeWidthMultiple,
			StopBoundaryBufferRangeWidth: cfg.StopBoundaryBufferRangeWidth,
			SourceRank:                   row.SourceRank,
			RankScore:                    row.RankScore,
		})
	}
	return out
}

func SelectFuturesRangeUniverseBreakoutRetestAcceptanceCandidatesFromSignals(signals []FuturesRangeUniverseBreakoutRetestAcceptanceSignalRow) []FuturesRangeUniverseBreakoutRetestAcceptanceSelectionRow {
	seen := map[string]bool{}
	rows := []FuturesRangeUniverseBreakoutRetestAcceptanceSelectionRow{}
	for _, signal := range signals {
		if seen[signal.CandidateID] {
			continue
		}
		seen[signal.CandidateID] = true
		rows = append(rows, FuturesRangeUniverseBreakoutRetestAcceptanceSelectionRow{
			SelectedOrder: signal.SelectedOrder,
			CandidateID:   signal.CandidateID,
			SourceRank:    signal.SourceRank,
			Timeframe:     signal.Timeframe,
			Family:        signal.Family,
			Side:          RangeDiscoverySideAll,
			HorizonBars:   signal.MaxHoldBars,
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].SelectedOrder < rows[j].SelectedOrder
	})
	return rows
}

func breakoutRetestAcceptanceCandidatesFromSelectionFromSignals(signals []FuturesRangeUniverseBreakoutRetestAcceptanceSignalRow, trades []FuturesRangeUniverseBreakoutRetestAcceptanceTradeRow, cfg FuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig) []FuturesRangeUniverseBreakoutRetestAcceptanceCandidateConfig {
	seen := map[string]FuturesRangeUniverseBreakoutRetestAcceptanceCandidateConfig{}
	for _, signal := range signals {
		side := RangeDiscoverySideAll
		if signal.Side == Long {
			side = RangeDiscoverySideUp
		}
		if signal.Side == Short {
			side = RangeDiscoverySideDown
		}
		seen[signal.CandidateID] = FuturesRangeUniverseBreakoutRetestAcceptanceCandidateConfig{
			CandidateID:   signal.CandidateID,
			SelectedOrder: signal.SelectedOrder,
			Timeframe:     signal.Timeframe,
			Side:          side,
			MaxHoldBars:   signal.MaxHoldBars,
			SourceRank:    signal.SourceRank,
		}
	}
	for _, trade := range trades {
		if _, ok := seen[trade.CandidateID]; ok {
			continue
		}
		seen[trade.CandidateID] = FuturesRangeUniverseBreakoutRetestAcceptanceCandidateConfig{
			CandidateID:   trade.CandidateID,
			SelectedOrder: trade.SelectedOrder,
			Timeframe:     trade.Timeframe,
			Side:          RangeDiscoverySideAll,
			MaxHoldBars:   trade.HoldBars,
			SourceRank:    trade.SourceRank,
		}
	}
	out := make([]FuturesRangeUniverseBreakoutRetestAcceptanceCandidateConfig, 0, len(seen))
	for _, candidate := range seen {
		out = append(out, candidate)
	}
	sort.Slice(out, func(i, j int) bool {
		return candidateSelectedOrder(out[i]) < candidateSelectedOrder(out[j])
	})
	return out
}

func breakoutRetestAcceptanceCandidateID(timeframe, side string, horizonBars int) string {
	return fmt.Sprintf("breakout_retest_acceptance_%s_%s_h%d", strings.ToLower(timeframe), strings.ToLower(side), horizonBars)
}

func breakoutRetestAcceptanceDetectorConfig(cfg FuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig, barsPerDay int) RangeDetectorConfig {
	detectorCfg := DefaultCompressionRangeDetectorConfig()
	detectorCfg.BarsPerDay = barsPerDay
	detectorCfg.LookbackDays = cfg.DiscoveryConfig.Discovery.DetectorLookbackDays
	detectorCfg.Percentile = cfg.DiscoveryConfig.Discovery.DetectorPercentile
	detectorCfg.MinConsecutiveBars = cfg.DiscoveryConfig.Discovery.DetectorMinConsecutiveBars
	detectorCfg.UseBollinger = true
	detectorCfg.UseADX = false
	if cfg.DiscoveryConfig.Discovery.DetectorLookbackBarsOverride > 0 {
		detectorCfg.LookbackDays = 1
		detectorCfg.BarsPerDay = cfg.DiscoveryConfig.Discovery.DetectorLookbackBarsOverride
	}
	return detectorCfg
}

func breakoutRetestAcceptanceFrameDef(timeframe string) (rangeDiscoveryFrameDef, bool) {
	for _, frame := range rangeUniverseFrameDefs() {
		if frame.timeframe == timeframe {
			return frame, true
		}
	}
	return rangeDiscoveryFrameDef{}, false
}

func breakoutRetestAcceptanceCandidateAllowsSide(candidate FuturesRangeUniverseBreakoutRetestAcceptanceCandidateConfig, eventSide string) bool {
	switch eventSide {
	case RangeDiscoverySideUp:
		return candidate.Side == RangeDiscoverySideAll || candidate.Side == RangeDiscoverySideUp
	case RangeDiscoverySideDown:
		return candidate.Side == RangeDiscoverySideAll || candidate.Side == RangeDiscoverySideDown
	default:
		return false
	}
}

func filterBreakoutRetestAcceptanceSignals(signals []FuturesRangeUniverseBreakoutRetestAcceptanceSignalRow, candidateID string, symbol string, split Split, side string) []FuturesRangeUniverseBreakoutRetestAcceptanceSignalRow {
	out := make([]FuturesRangeUniverseBreakoutRetestAcceptanceSignalRow, 0, len(signals))
	for _, signal := range signals {
		if signal.CandidateID != candidateID {
			continue
		}
		if symbol != BreakoutRetestAcceptanceSummaryAggregateSymbol && signal.Symbol != symbol {
			continue
		}
		eventTime, err := parseTime(signal.RetestCloseTime)
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

func filterBreakoutRetestAcceptanceTrades(trades []FuturesRangeUniverseBreakoutRetestAcceptanceTradeRow, candidateID string, symbol string, split Split, side string) []FuturesRangeUniverseBreakoutRetestAcceptanceTradeRow {
	out := make([]FuturesRangeUniverseBreakoutRetestAcceptanceTradeRow, 0, len(trades))
	for _, trade := range trades {
		if trade.CandidateID != candidateID {
			continue
		}
		if symbol != BreakoutRetestAcceptanceSummaryAggregateSymbol && trade.Symbol != symbol {
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

func summarizeBreakoutRetestAcceptanceTrades(trades []FuturesRangeUniverseBreakoutRetestAcceptanceTradeRow, startBalance float64) FuturesRangeUniverseBreakoutRetestAcceptanceSummaryRow {
	row := FuturesRangeUniverseBreakoutRetestAcceptanceSummaryRow{TotalTrades: len(trades)}
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

func markBreakoutRetestAcceptanceReviewFlags(rows []FuturesRangeUniverseBreakoutRetestAcceptanceSummaryRow, cfg FuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig, startBalance float64, splits []Split) {
	candidates := breakoutRetestAcceptanceSummaryCandidates(rows)
	for _, candidate := range candidates {
		passes := breakoutRetestAcceptanceCandidatePasses(rows, candidate, cfg, startBalance, splits)
		for _, symbol := range breakoutRetestAcceptanceSummarySymbols(cfg) {
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
				if rows[i].Symbol == BreakoutRetestAcceptanceSummaryAggregateSymbol && rows[i].Side == "all" {
					rows[i].PassesReviewGate = passes
				}
			}
			if worstIndex >= 0 {
				rows[worstIndex].IsWorstPeriodSplit = true
			}
		}
	}
}

func breakoutRetestAcceptanceCandidatePasses(rows []FuturesRangeUniverseBreakoutRetestAcceptanceSummaryRow, candidate FuturesRangeUniverseBreakoutRetestAcceptanceCandidateConfig, cfg FuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig, startBalance float64, splits []Split) bool {
	byKey := breakoutRetestAcceptanceSummaryByKey(rows)
	aggregateFull := byKey[breakoutRetestAcceptanceSummaryKey(candidate.CandidateID, BreakoutRetestAcceptanceSummaryAggregateSymbol, fullSplitName, "all")]
	if aggregateFull.TotalTrades < cfg.MinFullTrades || aggregateFull.NetPnL <= 0 || aggregateFull.ProfitFactor < 1.2 {
		return false
	}
	if !breakoutRetestAcceptanceAggregateSplitsPass(byKey, candidate.CandidateID, cfg, startBalance, splits) {
		return false
	}
	symbolPositiveFull := 0
	for _, source := range cfg.DiscoveryConfig.Sources {
		symbol := strings.ToUpper(source.Symbol)
		full := byKey[breakoutRetestAcceptanceSummaryKey(candidate.CandidateID, symbol, fullSplitName, "all")]
		if full.TotalTrades >= cfg.MinKeySplitTrades && full.NetPnL > 0 && full.ProfitFactor >= 1.0 {
			symbolPositiveFull++
		}
	}
	return symbolPositiveFull >= 2
}

func breakoutRetestAcceptanceAggregateSplitsPass(byKey map[string]FuturesRangeUniverseBreakoutRetestAcceptanceSummaryRow, candidateID string, cfg FuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig, startBalance float64, splits []Split) bool {
	for _, splitName := range []string{"2023_2024_oos", "2025_2026_recent"} {
		row := byKey[breakoutRetestAcceptanceSummaryKey(candidateID, BreakoutRetestAcceptanceSummaryAggregateSymbol, splitName, "all")]
		if row.TotalTrades < cfg.MinKeySplitTrades || row.NetPnL <= 0 {
			return false
		}
	}
	for _, side := range []string{string(Long), string(Short)} {
		row := byKey[breakoutRetestAcceptanceSummaryKey(candidateID, BreakoutRetestAcceptanceSummaryAggregateSymbol, fullSplitName, side)]
		if row.TotalTrades >= cfg.MinKeySplitTrades && row.NetPnL < -startBalance*cfg.NearFlatNetPct*float64(len(cfg.DiscoveryConfig.Sources)) {
			return false
		}
	}
	return true
}

func breakoutRetestAcceptanceSummarySymbols(cfg FuturesRangeUniverseBreakoutRetestAcceptanceBaselineConfig) []string {
	symbols := []string{}
	for _, source := range cfg.DiscoveryConfig.Sources {
		symbols = append(symbols, strings.ToUpper(source.Symbol))
	}
	sort.Slice(symbols, func(i, j int) bool {
		return rangeUniverseSymbolSortKey(symbols[i]) < rangeUniverseSymbolSortKey(symbols[j])
	})
	symbols = append(symbols, BreakoutRetestAcceptanceSummaryAggregateSymbol)
	return symbols
}

func breakoutRetestAcceptanceSummaryCandidates(rows []FuturesRangeUniverseBreakoutRetestAcceptanceSummaryRow) []FuturesRangeUniverseBreakoutRetestAcceptanceCandidateConfig {
	seen := map[string]FuturesRangeUniverseBreakoutRetestAcceptanceCandidateConfig{}
	for _, row := range rows {
		if _, ok := seen[row.CandidateID]; ok {
			continue
		}
		seen[row.CandidateID] = FuturesRangeUniverseBreakoutRetestAcceptanceCandidateConfig{
			CandidateID:   row.CandidateID,
			SelectedOrder: 1,
			Timeframe:     row.Timeframe,
			Side:          RangeDiscoverySideAll,
			MaxHoldBars:   1,
		}
	}
	out := make([]FuturesRangeUniverseBreakoutRetestAcceptanceCandidateConfig, 0, len(seen))
	for _, candidate := range seen {
		out = append(out, candidate)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CandidateID < out[j].CandidateID
	})
	return out
}

func breakoutRetestAcceptanceSummaryKey(candidateID, symbol, split, side string) string {
	return candidateID + "|" + symbol + "|" + split + "|" + side
}

func breakoutRetestAcceptanceSummaryByKey(rows []FuturesRangeUniverseBreakoutRetestAcceptanceSummaryRow) map[string]FuturesRangeUniverseBreakoutRetestAcceptanceSummaryRow {
	byKey := map[string]FuturesRangeUniverseBreakoutRetestAcceptanceSummaryRow{}
	for _, row := range rows {
		byKey[breakoutRetestAcceptanceSummaryKey(row.CandidateID, row.Symbol, row.Split, row.Side)] = row
	}
	return byKey
}

func breakoutRetestAcceptanceTimeframeSortKey(timeframe string) int {
	switch timeframe {
	case RangeDiscoveryTimeframe5m:
		return 0
	case RangeDiscoveryTimeframe15m:
		return 1
	case RangeDiscoveryTimeframe1h:
		return 2
	case RangeDiscoveryTimeframe4h:
		return 3
	default:
		return 99
	}
}

func candidateSelectedOrder(candidate FuturesRangeUniverseBreakoutRetestAcceptanceCandidateConfig) int {
	if candidate.SelectedOrder > 0 {
		return candidate.SelectedOrder
	}
	switch candidate.CandidateID {
	case BreakoutRetestAcceptanceCandidate15MAllH12:
		return 1
	case BreakoutRetestAcceptanceCandidate1HAllH12:
		return 2
	default:
		return 99
	}
}
