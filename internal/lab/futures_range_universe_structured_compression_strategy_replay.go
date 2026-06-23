package lab

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

const (
	FuturesRangeUniverseStructuredCompressionStrategyReplayName = "futures_range_universe_structured_compression_strategy_replay"

	StructuredCompressionStrategyReplayConfigID = "sc4h_btc_diagnostic_eth_sol_cw2_h12_t1_00_sb0_00"

	StructuredCompressionStrategyReplayStopStateSourceGap            = "structured_compression_strategy_replay_source_gap"
	StructuredCompressionStrategyReplayStopStateCodegenOrTestBlocked = "structured_compression_strategy_replay_codegen_or_test_blocked"
	StructuredCompressionStrategyReplayStopStateRegressionOrMismatch = "structured_compression_strategy_replay_regression_or_mismatch"
	StructuredCompressionStrategyReplayStopStateFailedNoPromotion    = "structured_compression_strategy_replay_failed_no_promotion"
	StructuredCompressionStrategyReplayStopStatePassedWalkForward    = "structured_compression_strategy_replay_passed_needs_walk_forward_robustness_brief"
	StructuredCompressionStrategyReplayStopStateReviewOnlyNoStrategy = "structured_compression_strategy_replay_review_only_no_strategy_change"
)

type FuturesRangeUniverseStructuredCompressionStrategyReplayConfig struct {
	Sources                       []FuturesRangeUniverseSourceConfig
	ConfigID                      string
	SymbolSet                     string
	EventDelayBars                int
	ConfirmationWindowBars        int
	MaxHoldBars                   int
	TargetRangeWidthMultiple      float64
	StopBoundaryBufferRangeWidth  float64
	DetectorLookbackDays          int
	DetectorPercentile            float64
	DetectorMinConsecutiveBars    int
	MinFullTrades                 int
	MinKeySplitTrades             int
	NearFlatNetPct                float64
	MaxDrawdownLimit              float64
	ExpectedFullTrades            int
	ExpectedOOSTrades             int
	ExpectedRecentTrades          int
	ExpectedFullNetPnL            float64
	ExpectedFullProfitFactor      float64
	ExpectedFullMaxDrawdown       float64
	ExpectedMetricTolerance       float64
	ExpectedProfitFactorTolerance float64
	ExpectedDrawdownTolerance     float64
}

type FuturesRangeUniverseStructuredCompressionStrategyReplayResult struct {
	SourceRows   []FuturesRangeUniverseSourceRow
	CoverageRows []FuturesRangeUniverseCoverageRow
	SignalRows   []FuturesRangeUniverseStructuredCompressionStrategyReplaySignalRow
	TradeRows    []FuturesRangeUniverseStructuredCompressionStrategyReplayTradeRow
	SummaryRows  []FuturesRangeUniverseStructuredCompressionStrategyReplaySummaryRow
	Trades       []Trade
}

type FuturesRangeUniverseStructuredCompressionStrategyReplaySignalRow struct {
	ConfigID                     string  `json:"config_id"`
	SymbolSet                    string  `json:"symbol_set"`
	AuthoritySymbols             string  `json:"authority_symbols"`
	DiagnosticSymbols            string  `json:"diagnostic_symbols"`
	IsAuthority                  bool    `json:"is_authority"`
	IsDiagnostic                 bool    `json:"is_diagnostic"`
	ConfirmationWindowBars       int     `json:"confirmation_window_bars"`
	TargetRangeWidthMultiple     float64 `json:"target_range_width_multiple"`
	StopBoundaryBufferRangeWidth float64 `json:"stop_boundary_buffer_range_width"`
	FuturesRangeUniverseStructuredCompressionSignalRow
}

type FuturesRangeUniverseStructuredCompressionStrategyReplayTradeRow struct {
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

type FuturesRangeUniverseStructuredCompressionStrategyReplaySummaryRow struct {
	ConfigID                     string  `json:"config_id"`
	SymbolSet                    string  `json:"symbol_set"`
	AuthoritySymbols             string  `json:"authority_symbols"`
	DiagnosticSymbols            string  `json:"diagnostic_symbols"`
	IsAuthority                  bool    `json:"is_authority"`
	IsDiagnostic                 bool    `json:"is_diagnostic"`
	ConfirmationWindowBars       int     `json:"confirmation_window_bars"`
	TargetRangeWidthMultiple     float64 `json:"target_range_width_multiple"`
	StopBoundaryBufferRangeWidth float64 `json:"stop_boundary_buffer_range_width"`
	PassesReplayGate             bool    `json:"passes_replay_gate"`
	ReplayMismatch               bool    `json:"replay_mismatch"`
	FailureReason                string  `json:"failure_reason,omitempty"`
	FuturesRangeUniverseStructuredCompressionSummaryRow
}

func DefaultFuturesRangeUniverseStructuredCompressionStrategyReplayConfig() FuturesRangeUniverseStructuredCompressionStrategyReplayConfig {
	return FuturesRangeUniverseStructuredCompressionStrategyReplayConfig{
		Sources:                       DefaultFuturesRangeUniverseDiscoveryAuditConfig().Sources,
		ConfigID:                      StructuredCompressionStrategyReplayConfigID,
		SymbolSet:                     StructuredCompressionOptimizationSymbolSetBTCDiagnosticETHSOL,
		EventDelayBars:                StructuredCompressionEventDelayDefaultBars,
		ConfirmationWindowBars:        2,
		MaxHoldBars:                   12,
		TargetRangeWidthMultiple:      1.0,
		StopBoundaryBufferRangeWidth:  0.0,
		DetectorLookbackDays:          20,
		DetectorPercentile:            0.30,
		DetectorMinConsecutiveBars:    12,
		MinFullTrades:                 100,
		MinKeySplitTrades:             25,
		NearFlatNetPct:                0.005,
		MaxDrawdownLimit:              0.20,
		ExpectedFullTrades:            129,
		ExpectedOOSTrades:             43,
		ExpectedRecentTrades:          32,
		ExpectedFullNetPnL:            573.87,
		ExpectedFullProfitFactor:      1.8089,
		ExpectedFullMaxDrawdown:       0.0982,
		ExpectedMetricTolerance:       0.25,
		ExpectedProfitFactorTolerance: 0.002,
		ExpectedDrawdownTolerance:     0.002,
	}
}

func RunFuturesRangeUniverseStructuredCompressionStrategyReplay(cfg FuturesRangeUniverseStructuredCompressionStrategyReplayConfig, btCfg BacktestConfig, splits []Split) (FuturesRangeUniverseStructuredCompressionStrategyReplayResult, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return FuturesRangeUniverseStructuredCompressionStrategyReplayResult{}, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}

	result := FuturesRangeUniverseStructuredCompressionStrategyReplayResult{}
	frame, ok := structuredCompressionFrameDef(RangeDiscoveryTimeframe4h)
	if !ok {
		return result, fmt.Errorf("structured compression strategy replay missing 4h frame definition")
	}
	evaluated, authority, diagnostic, _ := structuredCompressionOptimizationSymbols(cfg.SymbolSet)
	candidate := FuturesRangeUniverseStructuredCompressionCandidateConfig{
		CandidateID:                  StructuredCompressionCandidate4HAllH12,
		Timeframe:                    RangeDiscoveryTimeframe4h,
		Side:                         RangeDiscoverySideAll,
		MaxHoldBars:                  cfg.MaxHoldBars,
		TargetRangeWidthMultiple:     cfg.TargetRangeWidthMultiple,
		StopBoundaryBufferRangeWidth: cfg.StopBoundaryBufferRangeWidth,
	}
	baseCfg := cfg.baselineConfig()

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
			return result, fmt.Errorf("%s 4h structured compression strategy replay resample rejected: %s", sourceRow.Symbol, uCoverage.ValidationError)
		}
		classifications, err := (CompressionRangeDetector{Config: structuredCompressionDetectorConfig(baseCfg, frame.barsPerDay)}).Classify(resampled)
		if err != nil {
			return result, err
		}
		candlesBySymbol[sourceRow.Symbol] = resampled
		classificationsBySymbol[sourceRow.Symbol] = classifications
	}

	authorityText := strings.Join(authority, ",")
	diagnosticText := strings.Join(diagnostic, ",")
	for _, symbol := range evaluated {
		frameCandles := candlesBySymbol[symbol]
		if len(frameCandles) == 0 {
			continue
		}
		strategy, err := newFuturesRangeUniverseStructuredCompressionStrategyFromClassifications(frameCandles, symbol, candidate, baseCfg, btCfg, classificationsBySymbol[symbol], splits)
		if err != nil {
			return result, err
		}
		isAuthority := stringInSlice(symbol, authority)
		isDiagnostic := stringInSlice(symbol, diagnostic)
		for _, signal := range strategy.SignalRows() {
			result.SignalRows = append(result.SignalRows, FuturesRangeUniverseStructuredCompressionStrategyReplaySignalRow{
				ConfigID:                     cfg.ConfigID,
				SymbolSet:                    cfg.SymbolSet,
				AuthoritySymbols:             authorityText,
				DiagnosticSymbols:            diagnosticText,
				IsAuthority:                  isAuthority,
				IsDiagnostic:                 isDiagnostic,
				ConfirmationWindowBars:       cfg.ConfirmationWindowBars,
				TargetRangeWidthMultiple:     cfg.TargetRangeWidthMultiple,
				StopBoundaryBufferRangeWidth: cfg.StopBoundaryBufferRangeWidth,
				FuturesRangeUniverseStructuredCompressionSignalRow: signal,
			})
		}
		run := RunBacktest(frameCandles, strategy, btCfg)
		if isAuthority {
			result.Trades = append(result.Trades, run.Trades...)
		}
		for _, trade := range strategy.TradeRows(run.Trades, splits) {
			result.TradeRows = append(result.TradeRows, FuturesRangeUniverseStructuredCompressionStrategyReplayTradeRow{
				ConfigID:                     cfg.ConfigID,
				SymbolSet:                    cfg.SymbolSet,
				AuthoritySymbols:             authorityText,
				DiagnosticSymbols:            diagnosticText,
				IsAuthority:                  isAuthority,
				IsDiagnostic:                 isDiagnostic,
				ConfirmationWindowBars:       cfg.ConfirmationWindowBars,
				TargetRangeWidthMultiple:     cfg.TargetRangeWidthMultiple,
				StopBoundaryBufferRangeWidth: cfg.StopBoundaryBufferRangeWidth,
				FuturesRangeUniverseStructuredCompressionTradeRow: trade,
			})
		}
	}

	sortStructuredCompressionStrategyReplayRows(&result)
	result.SummaryRows = SummarizeFuturesRangeUniverseStructuredCompressionStrategyReplay(result.SignalRows, result.TradeRows, cfg, btCfg.StartBalance, splits)
	return result, nil
}

func SummarizeFuturesRangeUniverseStructuredCompressionStrategyReplay(signals []FuturesRangeUniverseStructuredCompressionStrategyReplaySignalRow, trades []FuturesRangeUniverseStructuredCompressionStrategyReplayTradeRow, cfg FuturesRangeUniverseStructuredCompressionStrategyReplayConfig, startBalance float64, splits []Split) []FuturesRangeUniverseStructuredCompressionStrategyReplaySummaryRow {
	cfg = cfg.withDefaults()
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	evaluated, authority, diagnostic, _ := structuredCompressionOptimizationSymbols(cfg.SymbolSet)
	symbols := append([]string(nil), evaluated...)
	symbols = append(symbols, StructuredCompressionSummaryAggregateSymbol)
	rows := make([]FuturesRangeUniverseStructuredCompressionStrategyReplaySummaryRow, 0, len(symbols)*len(splits)*3)
	for _, symbol := range symbols {
		for _, split := range splits {
			for _, side := range []string{"all", string(Long), string(Short)} {
				filteredTrades := filterStructuredCompressionStrategyReplayTrades(trades, cfg.ConfigID, symbol, split, side)
				filteredSignals := filterStructuredCompressionStrategyReplaySignals(signals, cfg.ConfigID, symbol, split, side)
				base := summarizeStructuredCompressionTrades(filteredTrades, startBalance)
				base.CandidateID = StructuredCompressionCandidate4HAllH12
				base.Symbol = symbol
				base.Timeframe = RangeDiscoveryTimeframe4h
				base.Split = split.Name
				base.Side = side
				base.SignalCount = len(filteredSignals)
				for _, signal := range filteredSignals {
					if signal.SkippedReason != "" {
						base.SkippedSignalCount++
					}
				}
				rows = append(rows, FuturesRangeUniverseStructuredCompressionStrategyReplaySummaryRow{
					ConfigID:                     cfg.ConfigID,
					SymbolSet:                    cfg.SymbolSet,
					AuthoritySymbols:             strings.Join(authority, ","),
					DiagnosticSymbols:            strings.Join(diagnostic, ","),
					IsAuthority:                  symbol == StructuredCompressionSummaryAggregateSymbol || stringInSlice(symbol, authority),
					IsDiagnostic:                 stringInSlice(symbol, diagnostic),
					ConfirmationWindowBars:       cfg.ConfirmationWindowBars,
					TargetRangeWidthMultiple:     cfg.TargetRangeWidthMultiple,
					StopBoundaryBufferRangeWidth: cfg.StopBoundaryBufferRangeWidth,
					FuturesRangeUniverseStructuredCompressionSummaryRow: base,
				})
			}
		}
	}
	pass, mismatch, reason := structuredCompressionStrategyReplayEvaluate(rows, cfg, startBalance)
	markStructuredCompressionStrategyReplayRows(rows, pass, mismatch, reason)
	return rows
}

func FuturesRangeUniverseStructuredCompressionStrategyReplayStopState(summary []FuturesRangeUniverseStructuredCompressionStrategyReplaySummaryRow, cfg FuturesRangeUniverseStructuredCompressionStrategyReplayConfig, startBalance float64, splits []Split) string {
	cfg = cfg.withDefaults()
	if len(summary) == 0 {
		return StructuredCompressionStrategyReplayStopStateFailedNoPromotion
	}
	pass, mismatch, _ := structuredCompressionStrategyReplayEvaluate(summary, cfg, startBalance)
	if mismatch {
		return StructuredCompressionStrategyReplayStopStateRegressionOrMismatch
	}
	if pass {
		return StructuredCompressionStrategyReplayStopStatePassedWalkForward
	}
	return StructuredCompressionStrategyReplayStopStateFailedNoPromotion
}

func (cfg FuturesRangeUniverseStructuredCompressionStrategyReplayConfig) withDefaults() FuturesRangeUniverseStructuredCompressionStrategyReplayConfig {
	defaults := DefaultFuturesRangeUniverseStructuredCompressionStrategyReplayConfig()
	if len(cfg.Sources) == 0 {
		cfg.Sources = defaults.Sources
	}
	if cfg.ConfigID == "" {
		cfg.ConfigID = defaults.ConfigID
	}
	if cfg.SymbolSet == "" {
		cfg.SymbolSet = defaults.SymbolSet
	}
	if cfg.EventDelayBars == 0 {
		cfg.EventDelayBars = defaults.EventDelayBars
	}
	if cfg.ConfirmationWindowBars == 0 {
		cfg.ConfirmationWindowBars = defaults.ConfirmationWindowBars
	}
	if cfg.MaxHoldBars == 0 {
		cfg.MaxHoldBars = defaults.MaxHoldBars
	}
	if cfg.TargetRangeWidthMultiple == 0 {
		cfg.TargetRangeWidthMultiple = defaults.TargetRangeWidthMultiple
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
	if cfg.ExpectedFullTrades == 0 {
		cfg.ExpectedFullTrades = defaults.ExpectedFullTrades
	}
	if cfg.ExpectedOOSTrades == 0 {
		cfg.ExpectedOOSTrades = defaults.ExpectedOOSTrades
	}
	if cfg.ExpectedRecentTrades == 0 {
		cfg.ExpectedRecentTrades = defaults.ExpectedRecentTrades
	}
	if cfg.ExpectedFullNetPnL == 0 {
		cfg.ExpectedFullNetPnL = defaults.ExpectedFullNetPnL
	}
	if cfg.ExpectedFullProfitFactor == 0 {
		cfg.ExpectedFullProfitFactor = defaults.ExpectedFullProfitFactor
	}
	if cfg.ExpectedFullMaxDrawdown == 0 {
		cfg.ExpectedFullMaxDrawdown = defaults.ExpectedFullMaxDrawdown
	}
	if cfg.ExpectedMetricTolerance == 0 {
		cfg.ExpectedMetricTolerance = defaults.ExpectedMetricTolerance
	}
	if cfg.ExpectedProfitFactorTolerance == 0 {
		cfg.ExpectedProfitFactorTolerance = defaults.ExpectedProfitFactorTolerance
	}
	if cfg.ExpectedDrawdownTolerance == 0 {
		cfg.ExpectedDrawdownTolerance = defaults.ExpectedDrawdownTolerance
	}
	return cfg
}

func (cfg FuturesRangeUniverseStructuredCompressionStrategyReplayConfig) validate() error {
	if len(cfg.Sources) == 0 {
		return fmt.Errorf("structured compression strategy replay source list cannot be empty")
	}
	if cfg.ConfigID != StructuredCompressionStrategyReplayConfigID {
		return fmt.Errorf("structured compression strategy replay config_id must be %q", StructuredCompressionStrategyReplayConfigID)
	}
	if cfg.SymbolSet != StructuredCompressionOptimizationSymbolSetBTCDiagnosticETHSOL {
		return fmt.Errorf("structured compression strategy replay symbol set must be %q", StructuredCompressionOptimizationSymbolSetBTCDiagnosticETHSOL)
	}
	if cfg.EventDelayBars <= 0 || cfg.ConfirmationWindowBars <= 0 || cfg.MaxHoldBars <= 0 {
		return fmt.Errorf("structured compression strategy replay event, confirmation, and max-hold bars must be positive")
	}
	if cfg.TargetRangeWidthMultiple <= 0 || cfg.StopBoundaryBufferRangeWidth < 0 {
		return fmt.Errorf("structured compression strategy replay target multiple must be positive and stop buffer non-negative")
	}
	if cfg.DetectorLookbackDays <= 0 || cfg.DetectorMinConsecutiveBars <= 0 {
		return fmt.Errorf("structured compression strategy replay detector settings must be positive")
	}
	if cfg.DetectorPercentile <= 0 || cfg.DetectorPercentile >= 1 {
		return fmt.Errorf("structured compression strategy replay detector percentile must be between 0 and 1")
	}
	if cfg.MinFullTrades <= 0 || cfg.MinKeySplitTrades <= 0 {
		return fmt.Errorf("structured compression strategy replay trade gates must be positive")
	}
	if cfg.NearFlatNetPct <= 0 || cfg.MaxDrawdownLimit <= 0 {
		return fmt.Errorf("structured compression strategy replay robustness gates must be positive")
	}
	seen := map[string]bool{}
	for _, source := range cfg.Sources {
		symbol := strings.ToUpper(strings.TrimSpace(source.Symbol))
		if !rangeUniverseApprovedSymbol(symbol) {
			return fmt.Errorf("structured compression strategy replay source symbol %q is not approved", symbol)
		}
		if seen[symbol] {
			return fmt.Errorf("structured compression strategy replay source symbol %q is duplicated", symbol)
		}
		seen[symbol] = true
	}
	for _, symbol := range []string{RangeUniverseSymbolBTCUSDT, RangeUniverseSymbolETHUSDT, RangeUniverseSymbolSOLUSDT} {
		if !seen[symbol] {
			return fmt.Errorf("structured compression strategy replay requires approved %s source", symbol)
		}
	}
	return nil
}

func (cfg FuturesRangeUniverseStructuredCompressionStrategyReplayConfig) baselineConfig() FuturesRangeUniverseStructuredCompressionBaselineConfig {
	return FuturesRangeUniverseStructuredCompressionBaselineConfig{
		Sources:                    cfg.Sources,
		EventDelayBars:             cfg.EventDelayBars,
		ConfirmationWindowBars:     cfg.ConfirmationWindowBars,
		DetectorLookbackDays:       cfg.DetectorLookbackDays,
		DetectorPercentile:         cfg.DetectorPercentile,
		DetectorMinConsecutiveBars: cfg.DetectorMinConsecutiveBars,
		MinFullTrades:              cfg.MinFullTrades,
		MinKeySplitTrades:          cfg.MinKeySplitTrades,
		NearFlatNetPct:             cfg.NearFlatNetPct,
		Candidates: []FuturesRangeUniverseStructuredCompressionCandidateConfig{{
			CandidateID:                  StructuredCompressionCandidate4HAllH12,
			Timeframe:                    RangeDiscoveryTimeframe4h,
			Side:                         RangeDiscoverySideAll,
			MaxHoldBars:                  cfg.MaxHoldBars,
			TargetRangeWidthMultiple:     cfg.TargetRangeWidthMultiple,
			StopBoundaryBufferRangeWidth: cfg.StopBoundaryBufferRangeWidth,
		}},
	}
}

func filterStructuredCompressionStrategyReplaySignals(signals []FuturesRangeUniverseStructuredCompressionStrategyReplaySignalRow, configID string, symbol string, split Split, side string) []FuturesRangeUniverseStructuredCompressionStrategyReplaySignalRow {
	out := []FuturesRangeUniverseStructuredCompressionStrategyReplaySignalRow{}
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

func filterStructuredCompressionStrategyReplayTrades(trades []FuturesRangeUniverseStructuredCompressionStrategyReplayTradeRow, configID string, symbol string, split Split, side string) []FuturesRangeUniverseStructuredCompressionTradeRow {
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

func structuredCompressionStrategyReplayEvaluate(rows []FuturesRangeUniverseStructuredCompressionStrategyReplaySummaryRow, cfg FuturesRangeUniverseStructuredCompressionStrategyReplayConfig, startBalance float64) (bool, bool, string) {
	byKey := structuredCompressionStrategyReplaySummaryByKey(rows)
	full := byKey[structuredCompressionStrategyReplaySummaryKey(cfg.ConfigID, StructuredCompressionSummaryAggregateSymbol, fullSplitName, "all")]
	stress := byKey[structuredCompressionStrategyReplaySummaryKey(cfg.ConfigID, StructuredCompressionSummaryAggregateSymbol, "2021_2022_stress", "all")]
	oos := byKey[structuredCompressionStrategyReplaySummaryKey(cfg.ConfigID, StructuredCompressionSummaryAggregateSymbol, "2023_2024_oos", "all")]
	recent := byKey[structuredCompressionStrategyReplaySummaryKey(cfg.ConfigID, StructuredCompressionSummaryAggregateSymbol, "2025_2026_recent", "all")]
	reasons := []string{}
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
	if stress.NetPnL <= 0 || oos.NetPnL <= 0 || recent.NetPnL <= 0 {
		reasons = append(reasons, "period_split_net_not_positive")
	}
	for _, side := range []string{string(Long), string(Short)} {
		row := byKey[structuredCompressionStrategyReplaySummaryKey(cfg.ConfigID, StructuredCompressionSummaryAggregateSymbol, fullSplitName, side)]
		if row.TotalTrades >= cfg.MinKeySplitTrades && row.NetPnL < -startBalance*2*cfg.NearFlatNetPct {
			reasons = append(reasons, side+"_side_loss_dominates")
		}
	}
	for _, symbol := range []string{RangeUniverseSymbolETHUSDT, RangeUniverseSymbolSOLUSDT} {
		row := byKey[structuredCompressionStrategyReplaySummaryKey(cfg.ConfigID, symbol, fullSplitName, "all")]
		if row.TotalTrades < cfg.MinKeySplitTrades || row.NetPnL <= 0 || row.ProfitFactor < 1.0 {
			reasons = append(reasons, strings.ToLower(symbol)+"_authority_symbol_weak")
		}
	}
	btc := byKey[structuredCompressionStrategyReplaySummaryKey(cfg.ConfigID, RangeUniverseSymbolBTCUSDT, fullSplitName, "all")]
	if btc.IsAuthority {
		reasons = append(reasons, "btcusdt_marked_authority")
	}

	mismatchReasons := []string{}
	if full.TotalTrades != cfg.ExpectedFullTrades {
		mismatchReasons = append(mismatchReasons, "full_trade_count_mismatch")
	}
	if oos.TotalTrades != cfg.ExpectedOOSTrades || recent.TotalTrades != cfg.ExpectedRecentTrades {
		mismatchReasons = append(mismatchReasons, "key_split_trade_count_mismatch")
	}
	if math.Abs(full.NetPnL-cfg.ExpectedFullNetPnL) > cfg.ExpectedMetricTolerance {
		mismatchReasons = append(mismatchReasons, "full_net_pnl_mismatch")
	}
	if math.Abs(full.ProfitFactor-cfg.ExpectedFullProfitFactor) > cfg.ExpectedProfitFactorTolerance {
		mismatchReasons = append(mismatchReasons, "full_pf_mismatch")
	}
	if math.Abs(full.MaxDrawdown-cfg.ExpectedFullMaxDrawdown) > cfg.ExpectedDrawdownTolerance {
		mismatchReasons = append(mismatchReasons, "max_drawdown_mismatch")
	}
	if len(mismatchReasons) > 0 {
		reasons = append(reasons, mismatchReasons...)
	}

	pass := len(reasons) == 0
	return pass, len(mismatchReasons) > 0, strings.Join(uniqueStrings(reasons), ";")
}

func markStructuredCompressionStrategyReplayRows(rows []FuturesRangeUniverseStructuredCompressionStrategyReplaySummaryRow, pass bool, mismatch bool, reason string) {
	worst := map[string]int{}
	for i := range rows {
		rows[i].PassesReplayGate = pass
		rows[i].ReplayMismatch = mismatch
		rows[i].FailureReason = reason
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

func structuredCompressionStrategyReplaySummaryKey(configID, symbol, split, side string) string {
	return configID + "|" + symbol + "|" + split + "|" + side
}

func structuredCompressionStrategyReplaySummaryByKey(rows []FuturesRangeUniverseStructuredCompressionStrategyReplaySummaryRow) map[string]FuturesRangeUniverseStructuredCompressionStrategyReplaySummaryRow {
	byKey := map[string]FuturesRangeUniverseStructuredCompressionStrategyReplaySummaryRow{}
	for _, row := range rows {
		byKey[structuredCompressionStrategyReplaySummaryKey(row.ConfigID, row.Symbol, row.Split, row.Side)] = row
	}
	return byKey
}

func sortStructuredCompressionStrategyReplayRows(result *FuturesRangeUniverseStructuredCompressionStrategyReplayResult) {
	sort.Slice(result.CoverageRows, func(i, j int) bool {
		return rangeUniverseSymbolSortKey(result.CoverageRows[i].Symbol) < rangeUniverseSymbolSortKey(result.CoverageRows[j].Symbol)
	})
	sort.Slice(result.SignalRows, func(i, j int) bool {
		if result.SignalRows[i].Symbol != result.SignalRows[j].Symbol {
			return rangeUniverseSymbolSortKey(result.SignalRows[i].Symbol) < rangeUniverseSymbolSortKey(result.SignalRows[j].Symbol)
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
}
