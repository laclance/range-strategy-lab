package lab

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

const (
	FuturesRangeUniverseStructuredCompressionOptimizationName = "futures_range_universe_structured_compression_optimization"

	StructuredCompressionOptimizationStopStateSourceGap            = "structured_compression_optimization_source_gap"
	StructuredCompressionOptimizationStopStateCodegenOrTestBlocked = "structured_compression_optimization_codegen_or_test_blocked"
	StructuredCompressionOptimizationStopStateFailedNoPromotion    = "structured_compression_optimization_failed_no_promotion"
	StructuredCompressionOptimizationStopStatePassedStrategySpec   = "structured_compression_optimization_passed_needs_strategy_spec"
	StructuredCompressionOptimizationStopStateMixedPortfolioReview = "structured_compression_optimization_mixed_needs_portfolio_stream_review"
	StructuredCompressionOptimizationStopStateReviewOnlyNoStrategy = "structured_compression_optimization_review_only_no_strategy_change"

	StructuredCompressionOptimizationSymbolSetAll                 = "BTC_ETH_SOL"
	StructuredCompressionOptimizationSymbolSetETHSOL              = "ETH_SOL"
	StructuredCompressionOptimizationSymbolSetBTCDiagnosticETHSOL = "BTC_DIAGNOSTIC_ETH_SOL"
)

type FuturesRangeUniverseStructuredCompressionOptimizationConfig struct {
	Sources                       []FuturesRangeUniverseSourceConfig
	EventDelayBars                int
	ConfirmationWindowBars        []int
	MaxHoldBars                   []int
	TargetRangeWidthMultiples     []float64
	StopBoundaryBufferRangeWidths []float64
	SymbolSets                    []string
	DetectorLookbackDays          int
	DetectorPercentile            float64
	DetectorMinConsecutiveBars    int
	MinFullTrades                 int
	MinKeySplitTrades             int
	NearFlatNetPct                float64
	MaxDrawdownLimit              float64
}

type FuturesRangeUniverseStructuredCompressionOptimizationResult struct {
	SourceRows       []FuturesRangeUniverseSourceRow
	CoverageRows     []FuturesRangeUniverseCoverageRow
	GridRows         []FuturesRangeUniverseStructuredCompressionOptimizationGridRow
	TradeRows        []FuturesRangeUniverseStructuredCompressionOptimizationTradeRow
	SummaryRows      []FuturesRangeUniverseStructuredCompressionOptimizationSummaryRow
	RankingRows      []FuturesRangeUniverseStructuredCompressionOptimizationRankingRow
	Trades           []Trade
	SelectedConfigID string
}

type FuturesRangeUniverseStructuredCompressionOptimizationGridConfig struct {
	ConfigID                     string
	CandidateID                  string
	Timeframe                    string
	SymbolSet                    string
	AuthoritySymbols             []string
	DiagnosticSymbols            []string
	EvaluatedSymbols             []string
	ConfirmationWindowBars       int
	MaxHoldBars                  int
	TargetRangeWidthMultiple     float64
	StopBoundaryBufferRangeWidth float64
}

type FuturesRangeUniverseStructuredCompressionOptimizationGridRow struct {
	ConfigID                     string  `json:"config_id"`
	CandidateID                  string  `json:"candidate_id"`
	Timeframe                    string  `json:"timeframe"`
	SymbolSet                    string  `json:"symbol_set"`
	AuthoritySymbols             string  `json:"authority_symbols"`
	DiagnosticSymbols            string  `json:"diagnostic_symbols"`
	EvaluatedSymbols             string  `json:"evaluated_symbols"`
	ConfirmationWindowBars       int     `json:"confirmation_window_bars"`
	MaxHoldBars                  int     `json:"max_hold_bars"`
	TargetRangeWidthMultiple     float64 `json:"target_range_width_multiple"`
	StopBoundaryBufferRangeWidth float64 `json:"stop_boundary_buffer_range_width"`
	FullTrades                   int     `json:"full_trades"`
	KeyOOSTrades                 int     `json:"key_oos_trades"`
	KeyRecentTrades              int     `json:"key_recent_trades"`
	FullGrossPnL                 float64 `json:"full_gross_pnl"`
	FullNetPnL                   float64 `json:"full_net_pnl"`
	FullTotalCosts               float64 `json:"full_total_costs"`
	FullProfitFactor             float64 `json:"full_profit_factor"`
	FullMaxDrawdown              float64 `json:"full_max_drawdown"`
	FullAvgNetR                  float64 `json:"full_avg_net_r"`
	StressNetPnL                 float64 `json:"stress_net_pnl"`
	OOSNetPnL                    float64 `json:"oos_net_pnl"`
	RecentNetPnL                 float64 `json:"recent_net_pnl"`
	LongNetPnL                   float64 `json:"long_net_pnl"`
	ShortNetPnL                  float64 `json:"short_net_pnl"`
	BTCUSDTNetPnL                float64 `json:"btcusdt_net_pnl"`
	ETHUSDTNetPnL                float64 `json:"ethusdt_net_pnl"`
	SOLUSDTNetPnL                float64 `json:"solusdt_net_pnl"`
	ETHUSDTFullTrades            int     `json:"ethusdt_full_trades"`
	SOLUSDTFullTrades            int     `json:"solusdt_full_trades"`
	PassesGate                   bool    `json:"passes_gate"`
	NearViableForPortfolio       bool    `json:"near_viable_for_portfolio"`
	Selected                     bool    `json:"selected"`
	RankScore                    float64 `json:"rank_score"`
	FailureReason                string  `json:"failure_reason,omitempty"`
}

type FuturesRangeUniverseStructuredCompressionOptimizationTradeRow struct {
	ConfigID                     string  `json:"config_id"`
	SymbolSet                    string  `json:"symbol_set"`
	AuthoritySymbols             string  `json:"authority_symbols"`
	DiagnosticSymbols            string  `json:"diagnostic_symbols"`
	IsAuthority                  bool    `json:"is_authority"`
	IsDiagnostic                 bool    `json:"is_diagnostic"`
	ConfirmationWindowBars       int     `json:"confirmation_window_bars"`
	TargetRangeWidthMultiple     float64 `json:"target_range_width_multiple"`
	StopBoundaryBufferRangeWidth float64 `json:"stop_boundary_buffer_range_width"`
	FuturesRangeUniverseStructuredCompressionTradeRow
}

type FuturesRangeUniverseStructuredCompressionOptimizationSummaryRow struct {
	ConfigID                     string  `json:"config_id"`
	SymbolSet                    string  `json:"symbol_set"`
	AuthoritySymbols             string  `json:"authority_symbols"`
	DiagnosticSymbols            string  `json:"diagnostic_symbols"`
	IsAuthority                  bool    `json:"is_authority"`
	IsDiagnostic                 bool    `json:"is_diagnostic"`
	Selected                     bool    `json:"selected"`
	ConfirmationWindowBars       int     `json:"confirmation_window_bars"`
	TargetRangeWidthMultiple     float64 `json:"target_range_width_multiple"`
	StopBoundaryBufferRangeWidth float64 `json:"stop_boundary_buffer_range_width"`
	FuturesRangeUniverseStructuredCompressionSummaryRow
}

type FuturesRangeUniverseStructuredCompressionOptimizationRankingRow struct {
	Rank                         int     `json:"rank"`
	ConfigID                     string  `json:"config_id"`
	CandidateID                  string  `json:"candidate_id"`
	Timeframe                    string  `json:"timeframe"`
	SymbolSet                    string  `json:"symbol_set"`
	AuthoritySymbols             string  `json:"authority_symbols"`
	DiagnosticSymbols            string  `json:"diagnostic_symbols"`
	ConfirmationWindowBars       int     `json:"confirmation_window_bars"`
	MaxHoldBars                  int     `json:"max_hold_bars"`
	TargetRangeWidthMultiple     float64 `json:"target_range_width_multiple"`
	StopBoundaryBufferRangeWidth float64 `json:"stop_boundary_buffer_range_width"`
	FullTrades                   int     `json:"full_trades"`
	FullNetPnL                   float64 `json:"full_net_pnl"`
	FullProfitFactor             float64 `json:"full_profit_factor"`
	FullMaxDrawdown              float64 `json:"full_max_drawdown"`
	StressNetPnL                 float64 `json:"stress_net_pnl"`
	OOSNetPnL                    float64 `json:"oos_net_pnl"`
	RecentNetPnL                 float64 `json:"recent_net_pnl"`
	PassesGate                   bool    `json:"passes_gate"`
	NearViableForPortfolio       bool    `json:"near_viable_for_portfolio"`
	Selected                     bool    `json:"selected"`
	RankScore                    float64 `json:"rank_score"`
	FailureReason                string  `json:"failure_reason,omitempty"`
}

type futuresRangeUniverseStructuredCompressionOptimizationSignalRow struct {
	ConfigID                     string
	SymbolSet                    string
	IsAuthority                  bool
	IsDiagnostic                 bool
	ConfirmationWindowBars       int
	TargetRangeWidthMultiple     float64
	StopBoundaryBufferRangeWidth float64
	FuturesRangeUniverseStructuredCompressionSignalRow
}

func DefaultFuturesRangeUniverseStructuredCompressionOptimizationConfig() FuturesRangeUniverseStructuredCompressionOptimizationConfig {
	return FuturesRangeUniverseStructuredCompressionOptimizationConfig{
		Sources:                       DefaultFuturesRangeUniverseDiscoveryAuditConfig().Sources,
		EventDelayBars:                StructuredCompressionEventDelayDefaultBars,
		ConfirmationWindowBars:        []int{2, 3, 4},
		MaxHoldBars:                   []int{4, 6, 8, 12},
		TargetRangeWidthMultiples:     []float64{0.75, 1.0, 1.25},
		StopBoundaryBufferRangeWidths: []float64{0.0, 0.10},
		SymbolSets: []string{
			StructuredCompressionOptimizationSymbolSetAll,
			StructuredCompressionOptimizationSymbolSetETHSOL,
			StructuredCompressionOptimizationSymbolSetBTCDiagnosticETHSOL,
		},
		DetectorLookbackDays:       20,
		DetectorPercentile:         0.30,
		DetectorMinConsecutiveBars: 12,
		MinFullTrades:              100,
		MinKeySplitTrades:          25,
		NearFlatNetPct:             0.005,
		MaxDrawdownLimit:           0.20,
	}
}

func RunFuturesRangeUniverseStructuredCompressionOptimization(cfg FuturesRangeUniverseStructuredCompressionOptimizationConfig, btCfg BacktestConfig, splits []Split) (FuturesRangeUniverseStructuredCompressionOptimizationResult, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return FuturesRangeUniverseStructuredCompressionOptimizationResult{}, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}

	result := FuturesRangeUniverseStructuredCompressionOptimizationResult{}
	frame, ok := structuredCompressionFrameDef(RangeDiscoveryTimeframe4h)
	if !ok {
		return result, fmt.Errorf("structured compression optimization missing 4h frame definition")
	}
	candlesBySymbol := map[string][]Candle{}
	classificationsBySymbol := map[string][]RangeClassification{}
	for _, source := range cfg.Sources {
		candles, sourceRow, err := LoadFuturesRangeUniverseSource(source, splits)
		result.SourceRows = append(result.SourceRows, sourceRow)
		if err != nil {
			return result, err
		}
		resampled, coverage, err := resampleRangeDiscoveryFrame(candles, frame)
		uCoverage := FuturesRangeUniverseCoverageRow{Symbol: sourceRow.Symbol, FuturesRangeDiscoveryCoverageRow: coverage}
		result.CoverageRows = append(result.CoverageRows, uCoverage)
		if err != nil {
			return result, err
		}
		if !uCoverage.Complete || uCoverage.ValidationStatus != "accepted" {
			return result, fmt.Errorf("%s 4h structured compression optimization resample rejected: %s", sourceRow.Symbol, uCoverage.ValidationError)
		}
		detectorCfg := structuredCompressionOptimizationDetectorConfig(cfg, frame.barsPerDay)
		classifications, err := (CompressionRangeDetector{Config: detectorCfg}).Classify(resampled)
		if err != nil {
			return result, err
		}
		candlesBySymbol[sourceRow.Symbol] = resampled
		classificationsBySymbol[sourceRow.Symbol] = classifications
	}

	gridConfigs := structuredCompressionOptimizationGridConfigs(cfg)
	tradesByConfig := map[string][]Trade{}
	signalRows := []futuresRangeUniverseStructuredCompressionOptimizationSignalRow{}
	for _, grid := range gridConfigs {
		baseCfg := structuredCompressionOptimizationBaselineConfig(cfg, grid)
		candidate := FuturesRangeUniverseStructuredCompressionCandidateConfig{
			CandidateID:                  StructuredCompressionCandidate4HAllH6,
			Timeframe:                    RangeDiscoveryTimeframe4h,
			Side:                         RangeDiscoverySideAll,
			MaxHoldBars:                  grid.MaxHoldBars,
			TargetRangeWidthMultiple:     grid.TargetRangeWidthMultiple,
			StopBoundaryBufferRangeWidth: grid.StopBoundaryBufferRangeWidth,
		}
		for _, symbol := range grid.EvaluatedSymbols {
			frameCandles := candlesBySymbol[symbol]
			if len(frameCandles) == 0 {
				continue
			}
			strategy, err := newFuturesRangeUniverseStructuredCompressionStrategyFromClassifications(frameCandles, symbol, candidate, baseCfg, btCfg, classificationsBySymbol[symbol], splits)
			if err != nil {
				return result, err
			}
			isAuthority := stringInSlice(symbol, grid.AuthoritySymbols)
			isDiagnostic := stringInSlice(symbol, grid.DiagnosticSymbols)
			for _, signal := range strategy.SignalRows() {
				signalRows = append(signalRows, futuresRangeUniverseStructuredCompressionOptimizationSignalRow{
					ConfigID:                     grid.ConfigID,
					SymbolSet:                    grid.SymbolSet,
					IsAuthority:                  isAuthority,
					IsDiagnostic:                 isDiagnostic,
					ConfirmationWindowBars:       grid.ConfirmationWindowBars,
					TargetRangeWidthMultiple:     grid.TargetRangeWidthMultiple,
					StopBoundaryBufferRangeWidth: grid.StopBoundaryBufferRangeWidth,
					FuturesRangeUniverseStructuredCompressionSignalRow: signal,
				})
			}
			run := RunBacktest(frameCandles, strategy, btCfg)
			if isAuthority {
				tradesByConfig[grid.ConfigID] = append(tradesByConfig[grid.ConfigID], run.Trades...)
			}
			for _, trade := range strategy.TradeRows(run.Trades, splits) {
				result.TradeRows = append(result.TradeRows, FuturesRangeUniverseStructuredCompressionOptimizationTradeRow{
					ConfigID:                     grid.ConfigID,
					SymbolSet:                    grid.SymbolSet,
					AuthoritySymbols:             strings.Join(grid.AuthoritySymbols, ","),
					DiagnosticSymbols:            strings.Join(grid.DiagnosticSymbols, ","),
					IsAuthority:                  isAuthority,
					IsDiagnostic:                 isDiagnostic,
					ConfirmationWindowBars:       grid.ConfirmationWindowBars,
					TargetRangeWidthMultiple:     grid.TargetRangeWidthMultiple,
					StopBoundaryBufferRangeWidth: grid.StopBoundaryBufferRangeWidth,
					FuturesRangeUniverseStructuredCompressionTradeRow: trade,
				})
			}
		}
	}

	sortStructuredCompressionOptimizationRows(&result, signalRows)
	result.SummaryRows = SummarizeFuturesRangeUniverseStructuredCompressionOptimization(signalRows, result.TradeRows, gridConfigs, cfg, btCfg.StartBalance, splits)
	result.GridRows = structuredCompressionOptimizationGridRows(result.SummaryRows, gridConfigs, cfg, btCfg.StartBalance, splits)
	result.RankingRows = structuredCompressionOptimizationRankingRows(result.GridRows)
	for _, ranking := range result.RankingRows {
		if ranking.PassesGate {
			result.SelectedConfigID = ranking.ConfigID
			break
		}
	}
	if result.SelectedConfigID != "" {
		result.Trades = append([]Trade(nil), tradesByConfig[result.SelectedConfigID]...)
		sort.Slice(result.Trades, func(i, j int) bool {
			if result.Trades[i].EntryTime != result.Trades[j].EntryTime {
				return result.Trades[i].EntryTime < result.Trades[j].EntryTime
			}
			return result.Trades[i].Signal < result.Trades[j].Signal
		})
		markStructuredCompressionOptimizationSelected(result.SelectedConfigID, result.GridRows, result.SummaryRows, result.RankingRows)
	}
	return result, nil
}

func SummarizeFuturesRangeUniverseStructuredCompressionOptimization(signals []futuresRangeUniverseStructuredCompressionOptimizationSignalRow, trades []FuturesRangeUniverseStructuredCompressionOptimizationTradeRow, gridConfigs []FuturesRangeUniverseStructuredCompressionOptimizationGridConfig, cfg FuturesRangeUniverseStructuredCompressionOptimizationConfig, startBalance float64, splits []Split) []FuturesRangeUniverseStructuredCompressionOptimizationSummaryRow {
	cfg = cfg.withDefaults()
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	rows := []FuturesRangeUniverseStructuredCompressionOptimizationSummaryRow{}
	for _, grid := range gridConfigs {
		for _, symbol := range structuredCompressionOptimizationSummarySymbols(grid) {
			for _, split := range splits {
				for _, side := range []string{"all", string(Long), string(Short)} {
					filteredTrades := filterStructuredCompressionOptimizationTrades(trades, grid.ConfigID, symbol, split, side)
					filteredSignals := filterStructuredCompressionOptimizationSignals(signals, grid.ConfigID, symbol, split, side)
					row := summarizeStructuredCompressionTrades(filteredTrades, startBalance)
					row.CandidateID = grid.CandidateID
					row.Symbol = symbol
					row.Timeframe = grid.Timeframe
					row.Split = split.Name
					row.Side = side
					row.SignalCount = len(filteredSignals)
					for _, signal := range filteredSignals {
						if signal.SkippedReason != "" {
							row.SkippedSignalCount++
						}
					}
					rows = append(rows, FuturesRangeUniverseStructuredCompressionOptimizationSummaryRow{
						ConfigID:                     grid.ConfigID,
						SymbolSet:                    grid.SymbolSet,
						AuthoritySymbols:             strings.Join(grid.AuthoritySymbols, ","),
						DiagnosticSymbols:            strings.Join(grid.DiagnosticSymbols, ","),
						IsAuthority:                  symbol == StructuredCompressionSummaryAggregateSymbol || stringInSlice(symbol, grid.AuthoritySymbols),
						IsDiagnostic:                 stringInSlice(symbol, grid.DiagnosticSymbols),
						ConfirmationWindowBars:       grid.ConfirmationWindowBars,
						TargetRangeWidthMultiple:     grid.TargetRangeWidthMultiple,
						StopBoundaryBufferRangeWidth: grid.StopBoundaryBufferRangeWidth,
						FuturesRangeUniverseStructuredCompressionSummaryRow: row,
					})
				}
			}
		}
	}
	markStructuredCompressionOptimizationWorstSplits(rows)
	return rows
}

func FuturesRangeUniverseStructuredCompressionOptimizationStopState(rankings []FuturesRangeUniverseStructuredCompressionOptimizationRankingRow) string {
	near := 0
	for _, ranking := range rankings {
		if ranking.PassesGate {
			return StructuredCompressionOptimizationStopStatePassedStrategySpec
		}
		if ranking.NearViableForPortfolio {
			near++
		}
	}
	if near >= 2 {
		return StructuredCompressionOptimizationStopStateMixedPortfolioReview
	}
	return StructuredCompressionOptimizationStopStateFailedNoPromotion
}

func (cfg FuturesRangeUniverseStructuredCompressionOptimizationConfig) withDefaults() FuturesRangeUniverseStructuredCompressionOptimizationConfig {
	defaults := DefaultFuturesRangeUniverseStructuredCompressionOptimizationConfig()
	if len(cfg.Sources) == 0 {
		cfg.Sources = defaults.Sources
	}
	if cfg.EventDelayBars == 0 {
		cfg.EventDelayBars = defaults.EventDelayBars
	}
	if len(cfg.ConfirmationWindowBars) == 0 {
		cfg.ConfirmationWindowBars = append([]int(nil), defaults.ConfirmationWindowBars...)
	}
	if len(cfg.MaxHoldBars) == 0 {
		cfg.MaxHoldBars = append([]int(nil), defaults.MaxHoldBars...)
	}
	if len(cfg.TargetRangeWidthMultiples) == 0 {
		cfg.TargetRangeWidthMultiples = append([]float64(nil), defaults.TargetRangeWidthMultiples...)
	}
	if len(cfg.StopBoundaryBufferRangeWidths) == 0 {
		cfg.StopBoundaryBufferRangeWidths = append([]float64(nil), defaults.StopBoundaryBufferRangeWidths...)
	}
	if len(cfg.SymbolSets) == 0 {
		cfg.SymbolSets = append([]string(nil), defaults.SymbolSets...)
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
	if cfg.MaxDrawdownLimit == 0 {
		cfg.MaxDrawdownLimit = defaults.MaxDrawdownLimit
	}
	return cfg
}

func (cfg FuturesRangeUniverseStructuredCompressionOptimizationConfig) validate() error {
	if len(cfg.Sources) == 0 {
		return fmt.Errorf("structured compression optimization source list cannot be empty")
	}
	if cfg.EventDelayBars <= 0 {
		return fmt.Errorf("structured compression optimization event delay bars must be positive")
	}
	for _, value := range cfg.ConfirmationWindowBars {
		if value <= 0 {
			return fmt.Errorf("structured compression optimization confirmation window values must be positive")
		}
	}
	for _, value := range cfg.MaxHoldBars {
		if value <= 0 {
			return fmt.Errorf("structured compression optimization max hold values must be positive")
		}
	}
	for _, value := range cfg.TargetRangeWidthMultiples {
		if value <= 0 {
			return fmt.Errorf("structured compression optimization target multiples must be positive")
		}
	}
	for _, value := range cfg.StopBoundaryBufferRangeWidths {
		if value < 0 {
			return fmt.Errorf("structured compression optimization stop buffers must be non-negative")
		}
	}
	if cfg.DetectorLookbackDays <= 0 || cfg.DetectorMinConsecutiveBars <= 0 {
		return fmt.Errorf("structured compression optimization detector settings must be positive")
	}
	if cfg.DetectorPercentile <= 0 || cfg.DetectorPercentile >= 1 {
		return fmt.Errorf("structured compression optimization detector percentile must be between 0 and 1")
	}
	if cfg.MinFullTrades <= 0 || cfg.MinKeySplitTrades <= 0 {
		return fmt.Errorf("structured compression optimization trade gates must be positive")
	}
	if cfg.NearFlatNetPct <= 0 || cfg.MaxDrawdownLimit <= 0 {
		return fmt.Errorf("structured compression optimization robustness gates must be positive")
	}
	seen := map[string]bool{}
	for _, source := range cfg.Sources {
		symbol := strings.ToUpper(strings.TrimSpace(source.Symbol))
		if !rangeUniverseApprovedSymbol(symbol) {
			return fmt.Errorf("structured compression optimization source symbol %q is not approved", symbol)
		}
		if seen[symbol] {
			return fmt.Errorf("structured compression optimization source symbol %q is duplicated", symbol)
		}
		seen[symbol] = true
	}
	for _, symbolSet := range cfg.SymbolSets {
		if _, _, _, ok := structuredCompressionOptimizationSymbols(symbolSet); !ok {
			return fmt.Errorf("structured compression optimization unsupported symbol set %q", symbolSet)
		}
	}
	return nil
}

func structuredCompressionOptimizationDetectorConfig(cfg FuturesRangeUniverseStructuredCompressionOptimizationConfig, barsPerDay int) RangeDetectorConfig {
	detectorCfg := DefaultCompressionRangeDetectorConfig()
	detectorCfg.BarsPerDay = barsPerDay
	detectorCfg.LookbackDays = cfg.DetectorLookbackDays
	detectorCfg.Percentile = cfg.DetectorPercentile
	detectorCfg.MinConsecutiveBars = cfg.DetectorMinConsecutiveBars
	detectorCfg.UseBollinger = true
	detectorCfg.UseADX = false
	return detectorCfg
}

func structuredCompressionOptimizationBaselineConfig(cfg FuturesRangeUniverseStructuredCompressionOptimizationConfig, grid FuturesRangeUniverseStructuredCompressionOptimizationGridConfig) FuturesRangeUniverseStructuredCompressionBaselineConfig {
	return FuturesRangeUniverseStructuredCompressionBaselineConfig{
		Sources:                    cfg.Sources,
		EventDelayBars:             cfg.EventDelayBars,
		ConfirmationWindowBars:     grid.ConfirmationWindowBars,
		DetectorLookbackDays:       cfg.DetectorLookbackDays,
		DetectorPercentile:         cfg.DetectorPercentile,
		DetectorMinConsecutiveBars: cfg.DetectorMinConsecutiveBars,
		MinFullTrades:              cfg.MinFullTrades,
		MinKeySplitTrades:          cfg.MinKeySplitTrades,
		NearFlatNetPct:             cfg.NearFlatNetPct,
	}
}

func structuredCompressionOptimizationGridConfigs(cfg FuturesRangeUniverseStructuredCompressionOptimizationConfig) []FuturesRangeUniverseStructuredCompressionOptimizationGridConfig {
	cfg = cfg.withDefaults()
	rows := []FuturesRangeUniverseStructuredCompressionOptimizationGridConfig{}
	for _, symbolSet := range cfg.SymbolSets {
		evaluated, authority, diagnostic, _ := structuredCompressionOptimizationSymbols(symbolSet)
		for _, confirmation := range cfg.ConfirmationWindowBars {
			for _, maxHold := range cfg.MaxHoldBars {
				for _, targetMultiple := range cfg.TargetRangeWidthMultiples {
					for _, stopBuffer := range cfg.StopBoundaryBufferRangeWidths {
						rows = append(rows, FuturesRangeUniverseStructuredCompressionOptimizationGridConfig{
							ConfigID:                     structuredCompressionOptimizationConfigID(symbolSet, confirmation, maxHold, targetMultiple, stopBuffer),
							CandidateID:                  StructuredCompressionCandidate4HAllH6,
							Timeframe:                    RangeDiscoveryTimeframe4h,
							SymbolSet:                    symbolSet,
							AuthoritySymbols:             authority,
							DiagnosticSymbols:            diagnostic,
							EvaluatedSymbols:             evaluated,
							ConfirmationWindowBars:       confirmation,
							MaxHoldBars:                  maxHold,
							TargetRangeWidthMultiple:     targetMultiple,
							StopBoundaryBufferRangeWidth: stopBuffer,
						})
					}
				}
			}
		}
	}
	return rows
}

func structuredCompressionOptimizationConfigID(symbolSet string, confirmation int, maxHold int, targetMultiple float64, stopBuffer float64) string {
	return fmt.Sprintf("sc4h_%s_cw%d_h%d_t%s_sb%s", strings.ToLower(symbolSet), confirmation, maxHold, formatGridFloatID(targetMultiple), formatGridFloatID(stopBuffer))
}

func formatGridFloatID(value float64) string {
	return strings.ReplaceAll(fmt.Sprintf("%.2f", value), ".", "_")
}

func structuredCompressionOptimizationSymbols(symbolSet string) (evaluated []string, authority []string, diagnostic []string, ok bool) {
	switch symbolSet {
	case StructuredCompressionOptimizationSymbolSetAll:
		return []string{RangeUniverseSymbolBTCUSDT, RangeUniverseSymbolETHUSDT, RangeUniverseSymbolSOLUSDT}, []string{RangeUniverseSymbolBTCUSDT, RangeUniverseSymbolETHUSDT, RangeUniverseSymbolSOLUSDT}, nil, true
	case StructuredCompressionOptimizationSymbolSetETHSOL:
		return []string{RangeUniverseSymbolETHUSDT, RangeUniverseSymbolSOLUSDT}, []string{RangeUniverseSymbolETHUSDT, RangeUniverseSymbolSOLUSDT}, nil, true
	case StructuredCompressionOptimizationSymbolSetBTCDiagnosticETHSOL:
		return []string{RangeUniverseSymbolBTCUSDT, RangeUniverseSymbolETHUSDT, RangeUniverseSymbolSOLUSDT}, []string{RangeUniverseSymbolETHUSDT, RangeUniverseSymbolSOLUSDT}, []string{RangeUniverseSymbolBTCUSDT}, true
	default:
		return nil, nil, nil, false
	}
}

func structuredCompressionOptimizationSummarySymbols(grid FuturesRangeUniverseStructuredCompressionOptimizationGridConfig) []string {
	symbols := append([]string(nil), grid.EvaluatedSymbols...)
	symbols = append(symbols, StructuredCompressionSummaryAggregateSymbol)
	return symbols
}

func filterStructuredCompressionOptimizationSignals(signals []futuresRangeUniverseStructuredCompressionOptimizationSignalRow, configID string, symbol string, split Split, side string) []futuresRangeUniverseStructuredCompressionOptimizationSignalRow {
	out := []futuresRangeUniverseStructuredCompressionOptimizationSignalRow{}
	for _, signal := range signals {
		if signal.ConfigID != configID {
			continue
		}
		if symbol == StructuredCompressionSummaryAggregateSymbol {
			if !signal.IsAuthority {
				continue
			}
		} else if signal.Symbol != symbol {
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

func filterStructuredCompressionOptimizationTrades(trades []FuturesRangeUniverseStructuredCompressionOptimizationTradeRow, configID string, symbol string, split Split, side string) []FuturesRangeUniverseStructuredCompressionTradeRow {
	out := []FuturesRangeUniverseStructuredCompressionTradeRow{}
	for _, trade := range trades {
		if trade.ConfigID != configID {
			continue
		}
		if symbol == StructuredCompressionSummaryAggregateSymbol {
			if !trade.IsAuthority {
				continue
			}
		} else if trade.Symbol != symbol {
			continue
		}
		exitTime, err := parseTime(trade.ExitTime)
		if err != nil || !split.Contains(exitTime) {
			continue
		}
		if side != "all" && string(trade.Side) != side {
			continue
		}
		out = append(out, trade.FuturesRangeUniverseStructuredCompressionTradeRow)
	}
	return out
}

func structuredCompressionOptimizationGridRows(summary []FuturesRangeUniverseStructuredCompressionOptimizationSummaryRow, gridConfigs []FuturesRangeUniverseStructuredCompressionOptimizationGridConfig, cfg FuturesRangeUniverseStructuredCompressionOptimizationConfig, startBalance float64, splits []Split) []FuturesRangeUniverseStructuredCompressionOptimizationGridRow {
	byKey := structuredCompressionOptimizationSummaryByKey(summary)
	rows := []FuturesRangeUniverseStructuredCompressionOptimizationGridRow{}
	for _, grid := range gridConfigs {
		row := FuturesRangeUniverseStructuredCompressionOptimizationGridRow{
			ConfigID:                     grid.ConfigID,
			CandidateID:                  grid.CandidateID,
			Timeframe:                    grid.Timeframe,
			SymbolSet:                    grid.SymbolSet,
			AuthoritySymbols:             strings.Join(grid.AuthoritySymbols, ","),
			DiagnosticSymbols:            strings.Join(grid.DiagnosticSymbols, ","),
			EvaluatedSymbols:             strings.Join(grid.EvaluatedSymbols, ","),
			ConfirmationWindowBars:       grid.ConfirmationWindowBars,
			MaxHoldBars:                  grid.MaxHoldBars,
			TargetRangeWidthMultiple:     grid.TargetRangeWidthMultiple,
			StopBoundaryBufferRangeWidth: grid.StopBoundaryBufferRangeWidth,
		}
		full := byKey[structuredCompressionOptimizationSummaryKey(grid.ConfigID, StructuredCompressionSummaryAggregateSymbol, fullSplitName, "all")]
		stress := byKey[structuredCompressionOptimizationSummaryKey(grid.ConfigID, StructuredCompressionSummaryAggregateSymbol, "2021_2022_stress", "all")]
		oos := byKey[structuredCompressionOptimizationSummaryKey(grid.ConfigID, StructuredCompressionSummaryAggregateSymbol, "2023_2024_oos", "all")]
		recent := byKey[structuredCompressionOptimizationSummaryKey(grid.ConfigID, StructuredCompressionSummaryAggregateSymbol, "2025_2026_recent", "all")]
		longRow := byKey[structuredCompressionOptimizationSummaryKey(grid.ConfigID, StructuredCompressionSummaryAggregateSymbol, fullSplitName, string(Long))]
		shortRow := byKey[structuredCompressionOptimizationSummaryKey(grid.ConfigID, StructuredCompressionSummaryAggregateSymbol, fullSplitName, string(Short))]
		btc := byKey[structuredCompressionOptimizationSummaryKey(grid.ConfigID, RangeUniverseSymbolBTCUSDT, fullSplitName, "all")]
		eth := byKey[structuredCompressionOptimizationSummaryKey(grid.ConfigID, RangeUniverseSymbolETHUSDT, fullSplitName, "all")]
		sol := byKey[structuredCompressionOptimizationSummaryKey(grid.ConfigID, RangeUniverseSymbolSOLUSDT, fullSplitName, "all")]
		row.FullTrades = full.TotalTrades
		row.KeyOOSTrades = oos.TotalTrades
		row.KeyRecentTrades = recent.TotalTrades
		row.FullGrossPnL = full.GrossPnL
		row.FullNetPnL = full.NetPnL
		row.FullTotalCosts = full.TotalCosts
		row.FullProfitFactor = full.ProfitFactor
		row.FullMaxDrawdown = full.MaxDrawdown
		row.FullAvgNetR = full.AvgNetR
		row.StressNetPnL = stress.NetPnL
		row.OOSNetPnL = oos.NetPnL
		row.RecentNetPnL = recent.NetPnL
		row.LongNetPnL = longRow.NetPnL
		row.ShortNetPnL = shortRow.NetPnL
		row.BTCUSDTNetPnL = btc.NetPnL
		row.ETHUSDTNetPnL = eth.NetPnL
		row.SOLUSDTNetPnL = sol.NetPnL
		row.ETHUSDTFullTrades = eth.TotalTrades
		row.SOLUSDTFullTrades = sol.TotalTrades
		row.PassesGate, row.NearViableForPortfolio, row.FailureReason = structuredCompressionOptimizationEvaluateGrid(grid, byKey, cfg, startBalance, splits)
		row.RankScore = structuredCompressionOptimizationRankScore(row, len(grid.AuthoritySymbols), startBalance)
		rows = append(rows, row)
	}
	return rows
}

func structuredCompressionOptimizationEvaluateGrid(grid FuturesRangeUniverseStructuredCompressionOptimizationGridConfig, byKey map[string]FuturesRangeUniverseStructuredCompressionOptimizationSummaryRow, cfg FuturesRangeUniverseStructuredCompressionOptimizationConfig, startBalance float64, splits []Split) (bool, bool, string) {
	full := byKey[structuredCompressionOptimizationSummaryKey(grid.ConfigID, StructuredCompressionSummaryAggregateSymbol, fullSplitName, "all")]
	stress := byKey[structuredCompressionOptimizationSummaryKey(grid.ConfigID, StructuredCompressionSummaryAggregateSymbol, "2021_2022_stress", "all")]
	oos := byKey[structuredCompressionOptimizationSummaryKey(grid.ConfigID, StructuredCompressionSummaryAggregateSymbol, "2023_2024_oos", "all")]
	recent := byKey[structuredCompressionOptimizationSummaryKey(grid.ConfigID, StructuredCompressionSummaryAggregateSymbol, "2025_2026_recent", "all")]
	reasons := []string{}
	authorityBalance := startBalance * float64(len(grid.AuthoritySymbols))
	if full.TotalTrades < cfg.MinFullTrades {
		reasons = append(reasons, "inadequate_full_trades")
	}
	if oos.TotalTrades < cfg.MinKeySplitTrades || recent.TotalTrades < cfg.MinKeySplitTrades {
		reasons = append(reasons, "inadequate_key_split_trades")
	}
	if full.NetPnL <= 0 {
		reasons = append(reasons, "full_net_not_positive")
	}
	if full.ProfitFactor < 1.2 {
		reasons = append(reasons, "full_pf_below_1_2")
	}
	if oos.NetPnL <= 0 || recent.NetPnL <= 0 {
		reasons = append(reasons, "key_split_net_not_positive")
	}
	if !structuredCompressionOptimizationNearFlatOrBetter(stress, authorityBalance, cfg.NearFlatNetPct) {
		reasons = append(reasons, "stress_split_fragile")
	}
	if full.MaxDrawdown > cfg.MaxDrawdownLimit {
		reasons = append(reasons, "drawdown_above_limit")
	}
	for _, side := range []string{string(Long), string(Short)} {
		row := byKey[structuredCompressionOptimizationSummaryKey(grid.ConfigID, StructuredCompressionSummaryAggregateSymbol, fullSplitName, side)]
		if row.TotalTrades >= cfg.MinKeySplitTrades && row.NetPnL < -authorityBalance*cfg.NearFlatNetPct {
			reasons = append(reasons, side+"_side_loss_dominates")
		}
	}
	for _, symbol := range grid.AuthoritySymbols {
		row := byKey[structuredCompressionOptimizationSummaryKey(grid.ConfigID, symbol, fullSplitName, "all")]
		if symbol == RangeUniverseSymbolETHUSDT || symbol == RangeUniverseSymbolSOLUSDT {
			if row.TotalTrades < cfg.MinKeySplitTrades || row.NetPnL <= 0 || row.ProfitFactor < 1.0 {
				reasons = append(reasons, strings.ToLower(symbol)+"_transfer_evidence_weak")
			}
		}
		if row.NetPnL < 0 && full.NetPnL > 0 && math.Abs(row.NetPnL) > full.NetPnL*0.5 {
			reasons = append(reasons, strings.ToLower(symbol)+"_symbol_loss_dominates")
		}
	}
	passes := len(reasons) == 0
	near := structuredCompressionOptimizationNearViable(grid, byKey, cfg, authorityBalance)
	return passes, near, strings.Join(uniqueStrings(reasons), ";")
}

func structuredCompressionOptimizationNearViable(grid FuturesRangeUniverseStructuredCompressionOptimizationGridConfig, byKey map[string]FuturesRangeUniverseStructuredCompressionOptimizationSummaryRow, cfg FuturesRangeUniverseStructuredCompressionOptimizationConfig, authorityBalance float64) bool {
	full := byKey[structuredCompressionOptimizationSummaryKey(grid.ConfigID, StructuredCompressionSummaryAggregateSymbol, fullSplitName, "all")]
	if full.TotalTrades < cfg.MinFullTrades || full.GrossPnL <= 0 || full.ProfitFactor < 1.0 {
		return false
	}
	if full.NetPnL < -authorityBalance*cfg.NearFlatNetPct {
		return false
	}
	for _, splitName := range []string{"2023_2024_oos", "2025_2026_recent"} {
		row := byKey[structuredCompressionOptimizationSummaryKey(grid.ConfigID, StructuredCompressionSummaryAggregateSymbol, splitName, "all")]
		if row.TotalTrades < cfg.MinKeySplitTrades || !structuredCompressionOptimizationNearFlatOrBetter(row, authorityBalance, cfg.NearFlatNetPct) {
			return false
		}
	}
	return true
}

func structuredCompressionOptimizationNearFlatOrBetter(row FuturesRangeUniverseStructuredCompressionOptimizationSummaryRow, balance float64, nearFlatPct float64) bool {
	if row.TotalTrades == 0 {
		return false
	}
	if row.NetPnL >= 0 {
		return true
	}
	return row.NetPnL >= -balance*nearFlatPct && row.GrossPnL > 0
}

func structuredCompressionOptimizationRankScore(row FuturesRangeUniverseStructuredCompressionOptimizationGridRow, authoritySymbolCount int, startBalance float64) float64 {
	drawdownPenalty := row.FullMaxDrawdown * startBalance * float64(authoritySymbolCount)
	return row.FullNetPnL + row.OOSNetPnL + row.RecentNetPnL + row.StressNetPnL - drawdownPenalty
}

func structuredCompressionOptimizationRankingRows(gridRows []FuturesRangeUniverseStructuredCompressionOptimizationGridRow) []FuturesRangeUniverseStructuredCompressionOptimizationRankingRow {
	rows := make([]FuturesRangeUniverseStructuredCompressionOptimizationRankingRow, 0, len(gridRows))
	for _, grid := range gridRows {
		rows = append(rows, FuturesRangeUniverseStructuredCompressionOptimizationRankingRow{
			ConfigID:                     grid.ConfigID,
			CandidateID:                  grid.CandidateID,
			Timeframe:                    grid.Timeframe,
			SymbolSet:                    grid.SymbolSet,
			AuthoritySymbols:             grid.AuthoritySymbols,
			DiagnosticSymbols:            grid.DiagnosticSymbols,
			ConfirmationWindowBars:       grid.ConfirmationWindowBars,
			MaxHoldBars:                  grid.MaxHoldBars,
			TargetRangeWidthMultiple:     grid.TargetRangeWidthMultiple,
			StopBoundaryBufferRangeWidth: grid.StopBoundaryBufferRangeWidth,
			FullTrades:                   grid.FullTrades,
			FullNetPnL:                   grid.FullNetPnL,
			FullProfitFactor:             grid.FullProfitFactor,
			FullMaxDrawdown:              grid.FullMaxDrawdown,
			StressNetPnL:                 grid.StressNetPnL,
			OOSNetPnL:                    grid.OOSNetPnL,
			RecentNetPnL:                 grid.RecentNetPnL,
			PassesGate:                   grid.PassesGate,
			NearViableForPortfolio:       grid.NearViableForPortfolio,
			RankScore:                    grid.RankScore,
			FailureReason:                grid.FailureReason,
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].PassesGate != rows[j].PassesGate {
			return rows[i].PassesGate
		}
		if rows[i].RankScore != rows[j].RankScore {
			return rows[i].RankScore > rows[j].RankScore
		}
		return rows[i].ConfigID < rows[j].ConfigID
	})
	for i := range rows {
		rows[i].Rank = i + 1
	}
	return rows
}

func markStructuredCompressionOptimizationSelected(selected string, gridRows []FuturesRangeUniverseStructuredCompressionOptimizationGridRow, summaryRows []FuturesRangeUniverseStructuredCompressionOptimizationSummaryRow, rankingRows []FuturesRangeUniverseStructuredCompressionOptimizationRankingRow) {
	for i := range gridRows {
		gridRows[i].Selected = gridRows[i].ConfigID == selected
	}
	for i := range summaryRows {
		summaryRows[i].Selected = summaryRows[i].ConfigID == selected
	}
	for i := range rankingRows {
		rankingRows[i].Selected = rankingRows[i].ConfigID == selected
	}
}

func markStructuredCompressionOptimizationWorstSplits(rows []FuturesRangeUniverseStructuredCompressionOptimizationSummaryRow) {
	worst := map[string]int{}
	for i := range rows {
		if rows[i].Split == fullSplitName || rows[i].Side != "all" {
			continue
		}
		key := rows[i].ConfigID + "|" + rows[i].Symbol
		if existing, ok := worst[key]; !ok || rows[i].NetPnL < rows[existing].NetPnL {
			worst[key] = i
		}
	}
	for _, index := range worst {
		rows[index].IsWorstPeriodSplit = true
	}
}

func structuredCompressionOptimizationSummaryKey(configID, symbol, split, side string) string {
	return configID + "|" + symbol + "|" + split + "|" + side
}

func structuredCompressionOptimizationSummaryByKey(rows []FuturesRangeUniverseStructuredCompressionOptimizationSummaryRow) map[string]FuturesRangeUniverseStructuredCompressionOptimizationSummaryRow {
	byKey := map[string]FuturesRangeUniverseStructuredCompressionOptimizationSummaryRow{}
	for _, row := range rows {
		byKey[structuredCompressionOptimizationSummaryKey(row.ConfigID, row.Symbol, row.Split, row.Side)] = row
	}
	return byKey
}

func sortStructuredCompressionOptimizationRows(result *FuturesRangeUniverseStructuredCompressionOptimizationResult, signals []futuresRangeUniverseStructuredCompressionOptimizationSignalRow) {
	sort.Slice(result.CoverageRows, func(i, j int) bool {
		return rangeUniverseSymbolSortKey(result.CoverageRows[i].Symbol) < rangeUniverseSymbolSortKey(result.CoverageRows[j].Symbol)
	})
	sort.Slice(result.TradeRows, func(i, j int) bool {
		if result.TradeRows[i].ConfigID != result.TradeRows[j].ConfigID {
			return result.TradeRows[i].ConfigID < result.TradeRows[j].ConfigID
		}
		if result.TradeRows[i].EntryTime != result.TradeRows[j].EntryTime {
			return result.TradeRows[i].EntryTime < result.TradeRows[j].EntryTime
		}
		return result.TradeRows[i].SignalID < result.TradeRows[j].SignalID
	})
	sort.Slice(signals, func(i, j int) bool {
		if signals[i].ConfigID != signals[j].ConfigID {
			return signals[i].ConfigID < signals[j].ConfigID
		}
		if signals[i].Symbol != signals[j].Symbol {
			return rangeUniverseSymbolSortKey(signals[i].Symbol) < rangeUniverseSymbolSortKey(signals[j].Symbol)
		}
		return signals[i].ConfirmationIndex < signals[j].ConfirmationIndex
	})
}

func uniqueStrings(values []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, value := range values {
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
	}
	return out
}

func stringInSlice(value string, values []string) bool {
	for _, candidate := range values {
		if value == candidate {
			return true
		}
	}
	return false
}
