package lab

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

const (
	FuturesRangeUniverseStructuredCompressionWalkForwardName = "futures_range_universe_structured_compression_walk_forward"

	StructuredCompressionWalkForwardStopStateSourceGap                  = "structured_compression_walk_forward_source_gap"
	StructuredCompressionWalkForwardStopStateCodegenOrTestBlocked       = "structured_compression_walk_forward_codegen_or_test_blocked"
	StructuredCompressionWalkForwardStopStateFailedNoPromotion          = "structured_compression_walk_forward_failed_no_promotion"
	StructuredCompressionWalkForwardStopStateFragileNeedsReview         = "structured_compression_walk_forward_fragile_needs_review"
	StructuredCompressionWalkForwardStopStatePassedCandidatePackage     = "structured_compression_walk_forward_passed_needs_candidate_strategy_package"
	StructuredCompressionWalkForwardStopStateReviewOnlyNoStrategy       = "structured_compression_walk_forward_review_only_no_strategy_change"
	StructuredCompressionWalkForwardFoldStressToOOS                     = "wf_2021_2022_train__2023_2024_test"
	StructuredCompressionWalkForwardFoldStressOOSToRecent               = "wf_2021_2024_train__2025_2026_test"
	StructuredCompressionWalkForwardFoldOOSToRecent                     = "wf_2023_2024_train__2025_2026_test"
	structuredCompressionWalkForwardReasonNoTrainingSelection           = "no_training_config_selected"
	structuredCompressionWalkForwardReasonSelectedTestNetNotPositive    = "selected_test_net_not_positive"
	structuredCompressionWalkForwardReasonFrozenTestNetNotPositive      = "frozen_test_net_not_positive"
	structuredCompressionWalkForwardReasonInadequateTestTrades          = "inadequate_test_trades"
	structuredCompressionWalkForwardReasonBTCRequired                   = "btcusdt_authority_required"
	structuredCompressionWalkForwardReasonFrozenMissing                 = "frozen_config_missing"
	structuredCompressionWalkForwardReasonSelectedWorseThanFrozen       = "selected_test_net_worse_than_frozen"
	structuredCompressionWalkForwardReasonSelectedNotFrozenOrEquivalent = "selected_not_frozen_or_equivalent"
)

type FuturesRangeUniverseStructuredCompressionWalkForwardConfig struct {
	OptimizationConfig       FuturesRangeUniverseStructuredCompressionOptimizationConfig
	FrozenReplayConfig       FuturesRangeUniverseStructuredCompressionStrategyReplayConfig
	FrozenConfigID           string
	MinMultiSplitTrainTrades int
	MinTrainSplitTrades      int
	MinTestTrades            int
	MinProfitFactor          float64
	NearFlatNetPct           float64
}

type FuturesRangeUniverseStructuredCompressionWalkForwardResult struct {
	SourceRows     []FuturesRangeUniverseSourceRow
	CoverageRows   []FuturesRangeUniverseCoverageRow
	GridRows       []FuturesRangeUniverseStructuredCompressionOptimizationGridRow
	FoldRows       []FuturesRangeUniverseStructuredCompressionWalkForwardFoldRow
	TradeRows      []FuturesRangeUniverseStructuredCompressionOptimizationTradeRow
	SummaryRows    []FuturesRangeUniverseStructuredCompressionOptimizationSummaryRow
	RankingRows    []FuturesRangeUniverseStructuredCompressionWalkForwardRankingRow
	Trades         []Trade
	FrozenConfigID string
}

type FuturesRangeUniverseStructuredCompressionWalkForwardFold struct {
	FoldID      string
	TrainSplits []string
	TestSplit   string
}

type FuturesRangeUniverseStructuredCompressionWalkForwardFoldRow struct {
	FoldID                               string  `json:"fold_id"`
	TrainSplits                          string  `json:"train_splits"`
	TestSplit                            string  `json:"test_split"`
	SelectedConfigID                     string  `json:"selected_config_id"`
	FrozenConfigID                       string  `json:"frozen_config_id"`
	SelectedRank                         int     `json:"selected_rank"`
	FrozenRank                           int     `json:"frozen_rank"`
	SelectedSymbolSet                    string  `json:"selected_symbol_set"`
	SelectedAuthoritySymbols             string  `json:"selected_authority_symbols"`
	SelectedDiagnosticSymbols            string  `json:"selected_diagnostic_symbols"`
	FrozenAuthoritySymbols               string  `json:"frozen_authority_symbols"`
	FrozenDiagnosticSymbols              string  `json:"frozen_diagnostic_symbols"`
	SelectedConfirmationWindowBars       int     `json:"selected_confirmation_window_bars"`
	SelectedMaxHoldBars                  int     `json:"selected_max_hold_bars"`
	SelectedTargetRangeWidthMultiple     float64 `json:"selected_target_range_width_multiple"`
	SelectedStopBoundaryBufferRangeWidth float64 `json:"selected_stop_boundary_buffer_range_width"`
	SelectedTrainTrades                  int     `json:"selected_train_trades"`
	SelectedTrainNetPnL                  float64 `json:"selected_train_net_pnl"`
	SelectedTrainProfitFactor            float64 `json:"selected_train_profit_factor"`
	SelectedTrainMaxDrawdown             float64 `json:"selected_train_max_drawdown"`
	SelectedTestTrades                   int     `json:"selected_test_trades"`
	SelectedTestNetPnL                   float64 `json:"selected_test_net_pnl"`
	SelectedTestProfitFactor             float64 `json:"selected_test_profit_factor"`
	SelectedTestMaxDrawdown              float64 `json:"selected_test_max_drawdown"`
	FrozenTrainTrades                    int     `json:"frozen_train_trades"`
	FrozenTrainNetPnL                    float64 `json:"frozen_train_net_pnl"`
	FrozenTrainProfitFactor              float64 `json:"frozen_train_profit_factor"`
	FrozenTrainMaxDrawdown               float64 `json:"frozen_train_max_drawdown"`
	FrozenTestTrades                     int     `json:"frozen_test_trades"`
	FrozenTestNetPnL                     float64 `json:"frozen_test_net_pnl"`
	FrozenTestProfitFactor               float64 `json:"frozen_test_profit_factor"`
	FrozenTestMaxDrawdown                float64 `json:"frozen_test_max_drawdown"`
	SelectedLongSideWeak                 bool    `json:"selected_long_side_weak"`
	SelectedShortSideWeak                bool    `json:"selected_short_side_weak"`
	SelectedETHUSDTWeak                  bool    `json:"selected_ethusdt_weak"`
	SelectedSOLUSDTWeak                  bool    `json:"selected_solusdt_weak"`
	BTCDiagnosticWeak                    bool    `json:"btc_diagnostic_weak"`
	RequiresBTCUSDTAuthority             bool    `json:"requires_btcusdt_authority"`
	SelectedFrozenOrEquivalent           bool    `json:"selected_frozen_or_equivalent"`
	SelectedExactFrozen                  bool    `json:"selected_exact_frozen"`
	PassesFoldGate                       bool    `json:"passes_fold_gate"`
	FailureReason                        string  `json:"failure_reason,omitempty"`
}

type FuturesRangeUniverseStructuredCompressionWalkForwardRankingRow struct {
	FoldID                       string  `json:"fold_id"`
	Rank                         int     `json:"rank"`
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
	TrainSplits                  string  `json:"train_splits"`
	TestSplit                    string  `json:"test_split"`
	TrainTrades                  int     `json:"train_trades"`
	TrainNetPnL                  float64 `json:"train_net_pnl"`
	TrainProfitFactor            float64 `json:"train_profit_factor"`
	TrainMaxDrawdown             float64 `json:"train_max_drawdown"`
	TestTrades                   int     `json:"test_trades"`
	TestNetPnL                   float64 `json:"test_net_pnl"`
	TestProfitFactor             float64 `json:"test_profit_factor"`
	TestMaxDrawdown              float64 `json:"test_max_drawdown"`
	TrainLongNetPnL              float64 `json:"train_long_net_pnl"`
	TrainShortNetPnL             float64 `json:"train_short_net_pnl"`
	TrainBTCUSDTNetPnL           float64 `json:"train_btcusdt_net_pnl"`
	TrainETHUSDTNetPnL           float64 `json:"train_ethusdt_net_pnl"`
	TrainSOLUSDTNetPnL           float64 `json:"train_solusdt_net_pnl"`
	TestBTCUSDTNetPnL            float64 `json:"test_btcusdt_net_pnl"`
	TestETHUSDTNetPnL            float64 `json:"test_ethusdt_net_pnl"`
	TestSOLUSDTNetPnL            float64 `json:"test_solusdt_net_pnl"`
	PassesTrainingGate           bool    `json:"passes_training_gate"`
	HistoricalComparisonOnly     bool    `json:"historical_comparison_only"`
	Selected                     bool    `json:"selected"`
	Frozen                       bool    `json:"frozen"`
	FrozenOrEquivalent           bool    `json:"frozen_or_equivalent"`
	RankScore                    float64 `json:"rank_score"`
	FailureReason                string  `json:"failure_reason,omitempty"`
}

type structuredCompressionWalkForwardMetrics struct {
	All   FuturesRangeUniverseStructuredCompressionSummaryRow
	Long  FuturesRangeUniverseStructuredCompressionSummaryRow
	Short FuturesRangeUniverseStructuredCompressionSummaryRow
	BTC   FuturesRangeUniverseStructuredCompressionSummaryRow
	ETH   FuturesRangeUniverseStructuredCompressionSummaryRow
	SOL   FuturesRangeUniverseStructuredCompressionSummaryRow
}

func DefaultFuturesRangeUniverseStructuredCompressionWalkForwardConfig() FuturesRangeUniverseStructuredCompressionWalkForwardConfig {
	optimizationCfg := DefaultFuturesRangeUniverseStructuredCompressionOptimizationConfig()
	replayCfg := DefaultFuturesRangeUniverseStructuredCompressionStrategyReplayConfig()
	return FuturesRangeUniverseStructuredCompressionWalkForwardConfig{
		OptimizationConfig:       optimizationCfg,
		FrozenReplayConfig:       replayCfg,
		FrozenConfigID:           StructuredCompressionStrategyReplayConfigID,
		MinMultiSplitTrainTrades: 100,
		MinTrainSplitTrades:      optimizationCfg.MinKeySplitTrades,
		MinTestTrades:            optimizationCfg.MinKeySplitTrades,
		MinProfitFactor:          1.2,
		NearFlatNetPct:           optimizationCfg.NearFlatNetPct,
	}
}

func RunFuturesRangeUniverseStructuredCompressionWalkForwardRobustness(cfg FuturesRangeUniverseStructuredCompressionWalkForwardConfig, btCfg BacktestConfig, splits []Split) (FuturesRangeUniverseStructuredCompressionWalkForwardResult, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return FuturesRangeUniverseStructuredCompressionWalkForwardResult{}, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}

	optimizationResult, err := RunFuturesRangeUniverseStructuredCompressionOptimization(cfg.OptimizationConfig, btCfg, splits)
	if err != nil {
		return FuturesRangeUniverseStructuredCompressionWalkForwardResult{}, err
	}
	replayResult, err := RunFuturesRangeUniverseStructuredCompressionStrategyReplay(cfg.FrozenReplayConfig, btCfg, splits)
	if err != nil {
		return FuturesRangeUniverseStructuredCompressionWalkForwardResult{}, err
	}

	result := FuturesRangeUniverseStructuredCompressionWalkForwardResult{
		SourceRows:     optimizationResult.SourceRows,
		CoverageRows:   optimizationResult.CoverageRows,
		GridRows:       optimizationResult.GridRows,
		TradeRows:      optimizationResult.TradeRows,
		SummaryRows:    optimizationResult.SummaryRows,
		Trades:         replayResult.Trades,
		FrozenConfigID: cfg.FrozenConfigID,
	}
	gridByID := structuredCompressionWalkForwardGridByID(structuredCompressionOptimizationGridConfigs(cfg.OptimizationConfig))
	for _, fold := range FuturesRangeUniverseStructuredCompressionWalkForwardFolds() {
		rankings := structuredCompressionWalkForwardRankingRows(fold, gridByID, optimizationResult.TradeRows, cfg, btCfg.StartBalance, splits)
		result.RankingRows = append(result.RankingRows, rankings...)
		result.FoldRows = append(result.FoldRows, structuredCompressionWalkForwardFoldRow(fold, rankings, cfg))
	}
	return result, nil
}

func FuturesRangeUniverseStructuredCompressionWalkForwardFolds() []FuturesRangeUniverseStructuredCompressionWalkForwardFold {
	return []FuturesRangeUniverseStructuredCompressionWalkForwardFold{
		{
			FoldID:      StructuredCompressionWalkForwardFoldStressToOOS,
			TrainSplits: []string{"2021_2022_stress"},
			TestSplit:   "2023_2024_oos",
		},
		{
			FoldID:      StructuredCompressionWalkForwardFoldStressOOSToRecent,
			TrainSplits: []string{"2021_2022_stress", "2023_2024_oos"},
			TestSplit:   "2025_2026_recent",
		},
		{
			FoldID:      StructuredCompressionWalkForwardFoldOOSToRecent,
			TrainSplits: []string{"2023_2024_oos"},
			TestSplit:   "2025_2026_recent",
		},
	}
}

func FuturesRangeUniverseStructuredCompressionWalkForwardStopState(folds []FuturesRangeUniverseStructuredCompressionWalkForwardFoldRow) string {
	if len(folds) == 0 {
		return StructuredCompressionWalkForwardStopStateFailedNoPromotion
	}
	passing := 0
	equivalent := 0
	exactFrozen := 0
	for _, fold := range folds {
		if fold.PassesFoldGate {
			passing++
		}
		if fold.SelectedFrozenOrEquivalent {
			equivalent++
		}
		if fold.SelectedExactFrozen {
			exactFrozen++
		}
	}
	if passing != len(folds) {
		if passing > 0 {
			return StructuredCompressionWalkForwardStopStateFragileNeedsReview
		}
		return StructuredCompressionWalkForwardStopStateFailedNoPromotion
	}
	if equivalent < 2 {
		return StructuredCompressionWalkForwardStopStateFragileNeedsReview
	}
	if exactFrozen < 2 {
		return StructuredCompressionWalkForwardStopStateReviewOnlyNoStrategy
	}
	return StructuredCompressionWalkForwardStopStatePassedCandidatePackage
}

func (cfg FuturesRangeUniverseStructuredCompressionWalkForwardConfig) withDefaults() FuturesRangeUniverseStructuredCompressionWalkForwardConfig {
	defaults := DefaultFuturesRangeUniverseStructuredCompressionWalkForwardConfig()
	cfg.OptimizationConfig = cfg.OptimizationConfig.withDefaults()
	cfg.FrozenReplayConfig = cfg.FrozenReplayConfig.withDefaults()
	if len(cfg.OptimizationConfig.Sources) == 0 {
		cfg.OptimizationConfig.Sources = defaults.OptimizationConfig.Sources
	}
	if len(cfg.FrozenReplayConfig.Sources) == 0 {
		cfg.FrozenReplayConfig.Sources = cfg.OptimizationConfig.Sources
	}
	if cfg.FrozenConfigID == "" {
		cfg.FrozenConfigID = defaults.FrozenConfigID
	}
	if cfg.MinMultiSplitTrainTrades == 0 {
		cfg.MinMultiSplitTrainTrades = defaults.MinMultiSplitTrainTrades
	}
	if cfg.MinTrainSplitTrades == 0 {
		cfg.MinTrainSplitTrades = defaults.MinTrainSplitTrades
	}
	if cfg.MinTestTrades == 0 {
		cfg.MinTestTrades = defaults.MinTestTrades
	}
	if cfg.MinProfitFactor == 0 {
		cfg.MinProfitFactor = defaults.MinProfitFactor
	}
	if cfg.NearFlatNetPct == 0 {
		cfg.NearFlatNetPct = defaults.NearFlatNetPct
	}
	cfg.FrozenReplayConfig.Sources = cfg.OptimizationConfig.Sources
	return cfg
}

func (cfg FuturesRangeUniverseStructuredCompressionWalkForwardConfig) validate() error {
	if err := cfg.OptimizationConfig.validate(); err != nil {
		return err
	}
	if err := cfg.FrozenReplayConfig.validate(); err != nil {
		return err
	}
	if cfg.FrozenConfigID != StructuredCompressionStrategyReplayConfigID {
		return fmt.Errorf("structured compression walk-forward frozen config must be %q", StructuredCompressionStrategyReplayConfigID)
	}
	if cfg.MinMultiSplitTrainTrades <= 0 || cfg.MinTrainSplitTrades <= 0 || cfg.MinTestTrades <= 0 {
		return fmt.Errorf("structured compression walk-forward trade gates must be positive")
	}
	if cfg.MinProfitFactor <= 0 || cfg.NearFlatNetPct <= 0 {
		return fmt.Errorf("structured compression walk-forward robustness gates must be positive")
	}
	seen := map[string]bool{}
	for _, source := range cfg.OptimizationConfig.Sources {
		seen[strings.ToUpper(strings.TrimSpace(source.Symbol))] = true
	}
	for _, symbol := range []string{RangeUniverseSymbolBTCUSDT, RangeUniverseSymbolETHUSDT, RangeUniverseSymbolSOLUSDT} {
		if !seen[symbol] {
			return fmt.Errorf("structured compression walk-forward requires approved %s source", symbol)
		}
	}
	return nil
}

func structuredCompressionWalkForwardRankingRows(fold FuturesRangeUniverseStructuredCompressionWalkForwardFold, gridByID map[string]FuturesRangeUniverseStructuredCompressionOptimizationGridConfig, trades []FuturesRangeUniverseStructuredCompressionOptimizationTradeRow, cfg FuturesRangeUniverseStructuredCompressionWalkForwardConfig, startBalance float64, splits []Split) []FuturesRangeUniverseStructuredCompressionWalkForwardRankingRow {
	rows := make([]FuturesRangeUniverseStructuredCompressionWalkForwardRankingRow, 0, len(gridByID))
	gridIDs := make([]string, 0, len(gridByID))
	for id := range gridByID {
		gridIDs = append(gridIDs, id)
	}
	sort.Strings(gridIDs)
	for _, id := range gridIDs {
		grid := gridByID[id]
		train := structuredCompressionWalkForwardMetricsForSplits(trades, grid.ConfigID, fold.TrainSplits, splits, startBalance)
		test := structuredCompressionWalkForwardMetricsForSplits(trades, grid.ConfigID, []string{fold.TestSplit}, splits, startBalance)
		passesTraining, historical, reason := structuredCompressionWalkForwardEvaluateTraining(grid, fold, train, trades, cfg, startBalance, splits)
		row := FuturesRangeUniverseStructuredCompressionWalkForwardRankingRow{
			FoldID:                       fold.FoldID,
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
			TrainSplits:                  strings.Join(fold.TrainSplits, "+"),
			TestSplit:                    fold.TestSplit,
			TrainTrades:                  train.All.TotalTrades,
			TrainNetPnL:                  train.All.NetPnL,
			TrainProfitFactor:            train.All.ProfitFactor,
			TrainMaxDrawdown:             train.All.MaxDrawdown,
			TestTrades:                   test.All.TotalTrades,
			TestNetPnL:                   test.All.NetPnL,
			TestProfitFactor:             test.All.ProfitFactor,
			TestMaxDrawdown:              test.All.MaxDrawdown,
			TrainLongNetPnL:              train.Long.NetPnL,
			TrainShortNetPnL:             train.Short.NetPnL,
			TrainBTCUSDTNetPnL:           train.BTC.NetPnL,
			TrainETHUSDTNetPnL:           train.ETH.NetPnL,
			TrainSOLUSDTNetPnL:           train.SOL.NetPnL,
			TestBTCUSDTNetPnL:            test.BTC.NetPnL,
			TestETHUSDTNetPnL:            test.ETH.NetPnL,
			TestSOLUSDTNetPnL:            test.SOL.NetPnL,
			PassesTrainingGate:           passesTraining,
			HistoricalComparisonOnly:     historical,
			Frozen:                       grid.ConfigID == cfg.FrozenConfigID,
			RankScore:                    structuredCompressionWalkForwardRankScore(grid, fold, train, trades, startBalance, splits),
			FailureReason:                reason,
		}
		rows = append(rows, row)
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].PassesTrainingGate != rows[j].PassesTrainingGate {
			return rows[i].PassesTrainingGate
		}
		if rows[i].HistoricalComparisonOnly != rows[j].HistoricalComparisonOnly {
			return !rows[i].HistoricalComparisonOnly
		}
		if rows[i].RankScore != rows[j].RankScore {
			return rows[i].RankScore > rows[j].RankScore
		}
		return rows[i].ConfigID < rows[j].ConfigID
	})
	selectedIndex := -1
	for i := range rows {
		rows[i].Rank = i + 1
		if selectedIndex == -1 && rows[i].PassesTrainingGate && !rows[i].HistoricalComparisonOnly {
			selectedIndex = i
		}
	}
	if selectedIndex >= 0 {
		selected := rows[selectedIndex].ConfigID
		selectedTestNet := rows[selectedIndex].TestNetPnL
		for i := range rows {
			rows[i].Selected = rows[i].ConfigID == selected
			rows[i].FrozenOrEquivalent = structuredCompressionWalkForwardFrozenOrEquivalent(rows[i], cfg.FrozenConfigID, selectedTestNet)
		}
	}
	return rows
}

func structuredCompressionWalkForwardFoldRow(fold FuturesRangeUniverseStructuredCompressionWalkForwardFold, rankings []FuturesRangeUniverseStructuredCompressionWalkForwardRankingRow, cfg FuturesRangeUniverseStructuredCompressionWalkForwardConfig) FuturesRangeUniverseStructuredCompressionWalkForwardFoldRow {
	var selected *FuturesRangeUniverseStructuredCompressionWalkForwardRankingRow
	var frozen *FuturesRangeUniverseStructuredCompressionWalkForwardRankingRow
	for i := range rankings {
		if rankings[i].Selected {
			selected = &rankings[i]
		}
		if rankings[i].Frozen {
			frozen = &rankings[i]
		}
	}
	row := FuturesRangeUniverseStructuredCompressionWalkForwardFoldRow{
		FoldID:         fold.FoldID,
		TrainSplits:    strings.Join(fold.TrainSplits, "+"),
		TestSplit:      fold.TestSplit,
		FrozenConfigID: cfg.FrozenConfigID,
	}
	reasons := []string{}
	if frozen == nil {
		reasons = append(reasons, structuredCompressionWalkForwardReasonFrozenMissing)
	} else {
		row.FrozenRank = frozen.Rank
		row.FrozenAuthoritySymbols = frozen.AuthoritySymbols
		row.FrozenDiagnosticSymbols = frozen.DiagnosticSymbols
		row.FrozenTrainTrades = frozen.TrainTrades
		row.FrozenTrainNetPnL = frozen.TrainNetPnL
		row.FrozenTrainProfitFactor = frozen.TrainProfitFactor
		row.FrozenTrainMaxDrawdown = frozen.TrainMaxDrawdown
		row.FrozenTestTrades = frozen.TestTrades
		row.FrozenTestNetPnL = frozen.TestNetPnL
		row.FrozenTestProfitFactor = frozen.TestProfitFactor
		row.FrozenTestMaxDrawdown = frozen.TestMaxDrawdown
	}
	if selected == nil {
		reasons = append(reasons, structuredCompressionWalkForwardReasonNoTrainingSelection)
	} else {
		row.SelectedConfigID = selected.ConfigID
		row.SelectedRank = selected.Rank
		row.SelectedSymbolSet = selected.SymbolSet
		row.SelectedAuthoritySymbols = selected.AuthoritySymbols
		row.SelectedDiagnosticSymbols = selected.DiagnosticSymbols
		row.SelectedConfirmationWindowBars = selected.ConfirmationWindowBars
		row.SelectedMaxHoldBars = selected.MaxHoldBars
		row.SelectedTargetRangeWidthMultiple = selected.TargetRangeWidthMultiple
		row.SelectedStopBoundaryBufferRangeWidth = selected.StopBoundaryBufferRangeWidth
		row.SelectedTrainTrades = selected.TrainTrades
		row.SelectedTrainNetPnL = selected.TrainNetPnL
		row.SelectedTrainProfitFactor = selected.TrainProfitFactor
		row.SelectedTrainMaxDrawdown = selected.TrainMaxDrawdown
		row.SelectedTestTrades = selected.TestTrades
		row.SelectedTestNetPnL = selected.TestNetPnL
		row.SelectedTestProfitFactor = selected.TestProfitFactor
		row.SelectedTestMaxDrawdown = selected.TestMaxDrawdown
		row.SelectedLongSideWeak = selected.TrainLongNetPnL < 0
		row.SelectedShortSideWeak = selected.TrainShortNetPnL < 0
		row.SelectedETHUSDTWeak = selected.TestETHUSDTNetPnL < 0
		row.SelectedSOLUSDTWeak = selected.TestSOLUSDTNetPnL < 0
		row.BTCDiagnosticWeak = strings.Contains(selected.DiagnosticSymbols, RangeUniverseSymbolBTCUSDT) && selected.TestBTCUSDTNetPnL < 0
		row.RequiresBTCUSDTAuthority = strings.Contains(selected.AuthoritySymbols, RangeUniverseSymbolBTCUSDT)
		row.SelectedExactFrozen = selected.ConfigID == cfg.FrozenConfigID
		row.SelectedFrozenOrEquivalent = structuredCompressionWalkForwardFrozenOrEquivalent(*selected, cfg.FrozenConfigID, selected.TestNetPnL)
		if selected.TestNetPnL <= 0 {
			reasons = append(reasons, structuredCompressionWalkForwardReasonSelectedTestNetNotPositive)
		}
		if row.RequiresBTCUSDTAuthority {
			reasons = append(reasons, structuredCompressionWalkForwardReasonBTCRequired)
		}
		if !row.SelectedFrozenOrEquivalent {
			reasons = append(reasons, structuredCompressionWalkForwardReasonSelectedNotFrozenOrEquivalent)
		}
	}
	if selected != nil && frozen != nil {
		if selected.TestTrades < cfg.MinTestTrades && frozen.TestTrades < cfg.MinTestTrades {
			reasons = append(reasons, structuredCompressionWalkForwardReasonInadequateTestTrades)
		}
		if frozen.TestTrades >= cfg.MinTestTrades && frozen.TestNetPnL <= 0 {
			reasons = append(reasons, structuredCompressionWalkForwardReasonFrozenTestNetNotPositive)
		}
		if selected.ConfigID != cfg.FrozenConfigID && selected.TestNetPnL < frozen.TestNetPnL {
			reasons = append(reasons, structuredCompressionWalkForwardReasonSelectedWorseThanFrozen)
		}
	}
	row.FailureReason = strings.Join(uniqueStrings(reasons), ";")
	row.PassesFoldGate = row.FailureReason == ""
	return row
}

func structuredCompressionWalkForwardEvaluateTraining(grid FuturesRangeUniverseStructuredCompressionOptimizationGridConfig, fold FuturesRangeUniverseStructuredCompressionWalkForwardFold, metrics structuredCompressionWalkForwardMetrics, trades []FuturesRangeUniverseStructuredCompressionOptimizationTradeRow, cfg FuturesRangeUniverseStructuredCompressionWalkForwardConfig, startBalance float64, splits []Split) (bool, bool, string) {
	reasons := []string{}
	historicalComparisonOnly := grid.SymbolSet == StructuredCompressionOptimizationSymbolSetAll
	if historicalComparisonOnly {
		reasons = append(reasons, "btcusdt_authority_comparison_only")
	}
	if len(fold.TrainSplits) > 1 && metrics.All.TotalTrades < cfg.MinMultiSplitTrainTrades {
		reasons = append(reasons, "inadequate_aggregate_train_trades")
	}
	for _, splitName := range fold.TrainSplits {
		splitMetrics := structuredCompressionWalkForwardMetricsForSplits(trades, grid.ConfigID, []string{splitName}, splits, startBalance)
		if splitMetrics.All.TotalTrades < cfg.MinTrainSplitTrades {
			reasons = append(reasons, "inadequate_"+splitName+"_train_trades")
		}
	}
	if metrics.All.NetPnL <= 0 {
		reasons = append(reasons, "train_net_not_positive")
	}
	if metrics.All.ProfitFactor < cfg.MinProfitFactor {
		reasons = append(reasons, "train_pf_below_1_2")
	}
	if metrics.Long.TotalTrades >= cfg.MinTrainSplitTrades && structuredCompressionWalkForwardLossDominates(metrics.Long.NetPnL, metrics.All.NetPnL) {
		reasons = append(reasons, "long_side_loss_dominates")
	}
	if metrics.Short.TotalTrades >= cfg.MinTrainSplitTrades && structuredCompressionWalkForwardLossDominates(metrics.Short.NetPnL, metrics.All.NetPnL) {
		reasons = append(reasons, "short_side_loss_dominates")
	}
	for _, symbol := range grid.AuthoritySymbols {
		if symbol == RangeUniverseSymbolBTCUSDT {
			reasons = append(reasons, "btcusdt_authority_required")
			continue
		}
		symbolMetrics := metrics.ETH
		if symbol == RangeUniverseSymbolSOLUSDT {
			symbolMetrics = metrics.SOL
		}
		if structuredCompressionWalkForwardLossDominates(symbolMetrics.NetPnL, metrics.All.NetPnL) {
			reasons = append(reasons, strings.ToLower(symbol)+"_authority_symbol_loss_dominates")
		}
	}
	return len(reasons) == 0, historicalComparisonOnly, strings.Join(uniqueStrings(reasons), ";")
}

func structuredCompressionWalkForwardLossDominates(partNetPnL float64, aggregateNetPnL float64) bool {
	return partNetPnL < 0 && aggregateNetPnL > 0 && -partNetPnL > aggregateNetPnL*0.5
}

func structuredCompressionWalkForwardRankScore(grid FuturesRangeUniverseStructuredCompressionOptimizationGridConfig, fold FuturesRangeUniverseStructuredCompressionWalkForwardFold, train structuredCompressionWalkForwardMetrics, trades []FuturesRangeUniverseStructuredCompressionOptimizationTradeRow, startBalance float64, splits []Split) float64 {
	row := FuturesRangeUniverseStructuredCompressionOptimizationGridRow{
		FullNetPnL:      train.All.NetPnL,
		FullMaxDrawdown: train.All.MaxDrawdown,
	}
	for _, splitName := range fold.TrainSplits {
		splitMetrics := structuredCompressionWalkForwardMetricsForSplits(trades, grid.ConfigID, []string{splitName}, splits, startBalance)
		switch splitName {
		case "2021_2022_stress":
			row.StressNetPnL = splitMetrics.All.NetPnL
		case "2023_2024_oos":
			row.OOSNetPnL = splitMetrics.All.NetPnL
		case "2025_2026_recent":
			row.RecentNetPnL = splitMetrics.All.NetPnL
		}
	}
	return structuredCompressionOptimizationRankScore(row, len(grid.AuthoritySymbols), startBalance)
}

func structuredCompressionWalkForwardFrozenOrEquivalent(row FuturesRangeUniverseStructuredCompressionWalkForwardRankingRow, frozenConfigID string, selectedTestNet float64) bool {
	if row.ConfigID == frozenConfigID {
		return true
	}
	if row.AuthoritySymbols != strings.Join([]string{RangeUniverseSymbolETHUSDT, RangeUniverseSymbolSOLUSDT}, ",") {
		return false
	}
	if row.ConfirmationWindowBars != 2 || row.MaxHoldBars != 12 {
		return false
	}
	return row.TestNetPnL >= selectedTestNet
}

func structuredCompressionWalkForwardMetricsForSplits(trades []FuturesRangeUniverseStructuredCompressionOptimizationTradeRow, configID string, splitNames []string, splits []Split, startBalance float64) structuredCompressionWalkForwardMetrics {
	return structuredCompressionWalkForwardMetrics{
		All:   structuredCompressionWalkForwardSummaryFor(trades, configID, StructuredCompressionSummaryAggregateSymbol, splitNames, "all", splits, startBalance),
		Long:  structuredCompressionWalkForwardSummaryFor(trades, configID, StructuredCompressionSummaryAggregateSymbol, splitNames, string(Long), splits, startBalance),
		Short: structuredCompressionWalkForwardSummaryFor(trades, configID, StructuredCompressionSummaryAggregateSymbol, splitNames, string(Short), splits, startBalance),
		BTC:   structuredCompressionWalkForwardSummaryFor(trades, configID, RangeUniverseSymbolBTCUSDT, splitNames, "all", splits, startBalance),
		ETH:   structuredCompressionWalkForwardSummaryFor(trades, configID, RangeUniverseSymbolETHUSDT, splitNames, "all", splits, startBalance),
		SOL:   structuredCompressionWalkForwardSummaryFor(trades, configID, RangeUniverseSymbolSOLUSDT, splitNames, "all", splits, startBalance),
	}
}

func structuredCompressionWalkForwardSummaryFor(trades []FuturesRangeUniverseStructuredCompressionOptimizationTradeRow, configID string, symbol string, splitNames []string, side string, splits []Split, startBalance float64) FuturesRangeUniverseStructuredCompressionSummaryRow {
	selected := []FuturesRangeUniverseStructuredCompressionTradeRow{}
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
		if side != "all" && string(trade.Side) != side {
			continue
		}
		exitTime, err := parseTime(trade.ExitTime)
		if err != nil || !structuredCompressionWalkForwardTimeInSplits(exitTime, splitNames, splits) {
			continue
		}
		selected = append(selected, trade.FuturesRangeUniverseStructuredCompressionTradeRow)
	}
	row := summarizeStructuredCompressionTrades(selected, startBalance)
	row.Symbol = symbol
	row.Split = strings.Join(splitNames, "+")
	row.Side = side
	return row
}

func structuredCompressionWalkForwardTimeInSplits(t time.Time, splitNames []string, splits []Split) bool {
	for _, splitName := range splitNames {
		for _, split := range splits {
			if split.Name == splitName && split.Contains(t) {
				return true
			}
		}
	}
	return false
}

func structuredCompressionWalkForwardGridByID(grids []FuturesRangeUniverseStructuredCompressionOptimizationGridConfig) map[string]FuturesRangeUniverseStructuredCompressionOptimizationGridConfig {
	byID := map[string]FuturesRangeUniverseStructuredCompressionOptimizationGridConfig{}
	for _, grid := range grids {
		byID[grid.ConfigID] = grid
	}
	return byID
}
