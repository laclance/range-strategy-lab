package lab

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

const (
	FuturesRangeFirstOccupancyRotationV1OptimizationName = "futures_range_first_occupancy_rotation_v1_optimization"

	RangeFirstOccupancyRotationV1BaselineConfigID = "range_occupancy_rotation_v1_1h_l48_w035_ow12_occ060_rec33_t66_h12_sb005"

	RangeFirstStrategyV1StopStateSourceGap                  = "range_first_strategy_v1_source_gap"
	RangeFirstStrategyV1StopStateNoValidSignals             = "range_first_strategy_v1_no_valid_signals"
	RangeFirstStrategyV1StopStateBaselineFailedGridReviewed = "range_first_strategy_v1_baseline_failed_grid_still_reviewed"
	RangeFirstStrategyV1StopStateOptimizerFailedNoReplay    = "range_first_strategy_v1_optimizer_failed_no_replay"
	RangeFirstStrategyV1StopStatePassedNeedsFixedReplaySpec = "range_first_strategy_v1_passed_needs_fixed_replay_spec"
	RangeFirstStrategyV1StopStateRejectedClosedFamily       = "range_first_strategy_v1_rejected_closed_family_reslice"

	rangeFirstOccupancyRotationSourcePath      = "../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv"
	rangeFirstOccupancyRotationExpectedRows    = 573984
	rangeFirstOccupancyRotationExpectedFirst   = "2021-01-01T00:00:00Z"
	rangeFirstOccupancyRotationExpectedLast    = "2026-06-16T23:55:00Z"
	rangeFirstOccupancyRotationExpectedZeroVol = 66
	rangeFirstOccupancyRotationExpected15MRows = 191328
	rangeFirstOccupancyRotationExpected15MLast = "2026-06-16T23:45:00Z"
	rangeFirstOccupancyRotationExpected1HRows  = 47832
	rangeFirstOccupancyRotationExpected1HLast  = "2026-06-16T23:00:00Z"

	rangeFirstOccupancyRotationFamily = "range_occupancy_rotation_v1"
)

type FuturesRangeFirstOccupancyRotationV1OptimizationConfig struct {
	SourcePath               string
	ApprovedSourcePath       string
	ExpectedSourceRows       int
	ExpectedFirstOpenTime    string
	ExpectedLastOpenTime     string
	ExpectedGapCount         int
	ExpectedDuplicateCount   int
	ExpectedZeroVolumeCount  int
	SkipSourceFactCheck      bool
	Timeframes               []string
	Expected15MRows          int
	Expected15MLastOpenTime  string
	Expected1HRows           int
	Expected1HLastOpenTime   string
	SkipCoverageCountCheck   bool
	LookbackHours            []int
	MaxWidthPcts             []float64
	OccupancyWindows         []int
	OccupancyZoneLevels      []float64
	OccupancyMinFractions    []float64
	RecaptureLevels          []float64
	TargetLevels             []float64
	MaxHoldBars              []int
	StopBufferWidths         []float64
	SideModes                []string
	MinTrainTrades           int
	MinOOSTrades             int
	MinRecentTrades          int
	MinFullTrades            int
	MinTrainProfitFactor     float64
	MinOOSProfitFactor       float64
	MinRecentProfitFactor    float64
	MinFullProfitFactor      float64
	MinNetToDrawdown         float64
	MinSideTradesForWeakness int
	MaxSideNetConcentration  float64
	SideConcentrationPenalty float64
	ThinSideTradeThreshold   int
	ThinSidePenalty          float64
}

type FuturesRangeFirstOccupancyRotationV1OptimizationResult struct {
	SourceRows       []FuturesRangeFirstOccupancyRotationV1SourceRow    `json:"source_rows"`
	CoverageRows     []FuturesRangeFirstOccupancyRotationV1CoverageRow  `json:"coverage_rows"`
	GridRows         []FuturesRangeFirstOccupancyRotationV1GridRow      `json:"grid_rows"`
	BaselineRows     []FuturesRangeFirstOccupancyRotationV1BaselineRow  `json:"baseline_rows"`
	SignalRows       []FuturesRangeFirstOccupancyRotationV1SignalRow    `json:"signal_rows"`
	TradeRows        []FuturesRangeFirstOccupancyRotationV1TradeRow     `json:"trade_rows"`
	SummaryRows      []FuturesRangeFirstOccupancyRotationV1SummaryRow   `json:"summary_rows"`
	RankingRows      []FuturesRangeFirstOccupancyRotationV1RankingRow   `json:"ranking_rows"`
	SelectionRows    []FuturesRangeFirstOccupancyRotationV1SelectionRow `json:"selection_rows"`
	SkipRows         []FuturesRangeFirstOccupancyRotationV1SkipRow      `json:"skip_rows"`
	Trades           []Trade                                            `json:"trades"`
	BaselineConfigID string                                             `json:"baseline_config_id"`
	SelectedConfigID string                                             `json:"selected_config_id,omitempty"`
	StopState        string                                             `json:"stop_state"`
}

type FuturesRangeFirstOccupancyRotationV1GridConfig struct {
	ConfigID             string
	Timeframe            string
	LookbackHours        int
	MaxWidthPct          float64
	OccupancyWindow      int
	OccupancyZoneLevel   float64
	OccupancyMinFraction float64
	RecaptureLevel       float64
	TargetLevel          float64
	MaxHoldBars          int
	StopBufferWidth      float64
	SideMode             string
}

type FuturesRangeFirstOccupancyRotationV1SourceRow struct {
	Path                    string `json:"path"`
	ApprovedPath            string `json:"approved_path"`
	Venue                   string `json:"venue"`
	Product                 string `json:"product"`
	Symbol                  string `json:"symbol"`
	Interval                string `json:"interval"`
	RowCount                int    `json:"row_count"`
	ExpectedRowCount        int    `json:"expected_row_count"`
	FirstOpenTime           string `json:"first_open_time"`
	ExpectedFirstOpenTime   string `json:"expected_first_open_time"`
	LastOpenTime            string `json:"last_open_time"`
	ExpectedLastOpenTime    string `json:"expected_last_open_time"`
	GapCount                int    `json:"gap_count"`
	ExpectedGapCount        int    `json:"expected_gap_count"`
	DuplicateCount          int    `json:"duplicate_count"`
	ExpectedDuplicateCount  int    `json:"expected_duplicate_count"`
	ZeroVolumeCount         int    `json:"zero_volume_count"`
	ExpectedZeroVolumeCount int    `json:"expected_zero_volume_count"`
	ComparisonOnly          bool   `json:"comparison_only"`
	SourceFactsPass         bool   `json:"source_facts_pass"`
	ValidationStatus        string `json:"validation_status"`
	ValidationError         string `json:"validation_error,omitempty"`
}

type FuturesRangeFirstOccupancyRotationV1CoverageRow struct {
	ExpectedRowCount      int    `json:"expected_row_count"`
	ExpectedFirstOpenTime string `json:"expected_first_open_time"`
	ExpectedLastOpenTime  string `json:"expected_last_open_time"`
	CoverageFactsPass     bool   `json:"coverage_facts_pass"`
	FuturesRangeDiscoveryCoverageRow
}

type FuturesRangeFirstOccupancyRotationV1SignalRow struct {
	ConfigID               string    `json:"config_id"`
	Family                 string    `json:"family"`
	IsBaseline             bool      `json:"is_baseline"`
	Selected               bool      `json:"selected"`
	SignalID               string    `json:"signal_id"`
	Timeframe              string    `json:"timeframe"`
	Split                  string    `json:"split"`
	SignalIndex            int       `json:"signal_index"`
	SignalOpenTime         string    `json:"signal_open_time"`
	SignalCloseTime        string    `json:"signal_close_time"`
	SignalOpen             float64   `json:"signal_open"`
	SignalHigh             float64   `json:"signal_high"`
	SignalLow              float64   `json:"signal_low"`
	SignalClose            float64   `json:"signal_close"`
	LookbackHours          int       `json:"lookback_hours"`
	MaxWidthPct            float64   `json:"max_width_pct"`
	OccupancyWindow        int       `json:"occupancy_window"`
	OccupancyZoneLevel     float64   `json:"occupancy_zone_level"`
	OccupancyMinFraction   float64   `json:"occupancy_min_fraction"`
	RequiredOccupancyCount int       `json:"required_occupancy_count"`
	RecaptureLevel         float64   `json:"recapture_level"`
	TargetLevel            float64   `json:"target_level"`
	StopBufferWidth        float64   `json:"stop_buffer_width"`
	SideMode               string    `json:"side_mode"`
	RangeHigh              float64   `json:"range_high"`
	RangeLow               float64   `json:"range_low"`
	RangeMid               float64   `json:"range_mid"`
	RangeQ1                float64   `json:"range_q1"`
	RangeQ3                float64   `json:"range_q3"`
	RangeWidth             float64   `json:"range_width"`
	RangeWidthPct          float64   `json:"range_width_pct"`
	LowerRecapture         float64   `json:"lower_recapture"`
	UpperRecapture         float64   `json:"upper_recapture"`
	InsideOccupancyCount   int       `json:"inside_occupancy_count"`
	LowerOccupancyCount    int       `json:"lower_occupancy_count"`
	UpperOccupancyCount    int       `json:"upper_occupancy_count"`
	Side                   Direction `json:"side"`
	EntryIndex             int       `json:"entry_index"`
	EntryOpenTime          string    `json:"entry_open_time"`
	EntryOpen              float64   `json:"entry_open"`
	ExpectedEntryPrice     float64   `json:"expected_entry_price"`
	Stop                   float64   `json:"stop"`
	Target                 float64   `json:"target"`
	MaxHoldBars            int       `json:"max_hold_bars"`
	EntryGeometryValid     bool      `json:"entry_geometry_valid"`
	Executed               bool      `json:"executed"`
	SkippedReason          string    `json:"skipped_reason,omitempty"`
}

type FuturesRangeFirstOccupancyRotationV1TradeRow struct {
	ConfigID        string    `json:"config_id"`
	Family          string    `json:"family"`
	IsBaseline      bool      `json:"is_baseline"`
	Selected        bool      `json:"selected"`
	SignalID        string    `json:"signal_id"`
	Timeframe       string    `json:"timeframe"`
	SignalIndex     int       `json:"signal_index"`
	SignalCloseTime string    `json:"signal_close_time"`
	LookbackHours   int       `json:"lookback_hours"`
	MaxWidthPct     float64   `json:"max_width_pct"`
	OccupancyWindow int       `json:"occupancy_window"`
	RecaptureLevel  float64   `json:"recapture_level"`
	TargetLevel     float64   `json:"target_level"`
	StopBufferWidth float64   `json:"stop_buffer_width"`
	RangeHigh       float64   `json:"range_high"`
	RangeLow        float64   `json:"range_low"`
	RangeWidth      float64   `json:"range_width"`
	RangeMid        float64   `json:"range_mid"`
	EntrySplit      string    `json:"entry_split"`
	CloseSplit      string    `json:"close_split"`
	Side            Direction `json:"side"`
	EntryTime       string    `json:"entry_time"`
	ExitTime        string    `json:"exit_time"`
	OpenIndex       int       `json:"open_index"`
	CloseIndex      int       `json:"close_index"`
	EntryPrice      float64   `json:"entry_price"`
	ExitPrice       float64   `json:"exit_price"`
	Stop            float64   `json:"stop"`
	Target          float64   `json:"target"`
	Size            float64   `json:"size"`
	InitialRisk     float64   `json:"initial_risk"`
	GrossPnL        float64   `json:"gross_pnl"`
	NetPnL          float64   `json:"net_pnl"`
	Fees            float64   `json:"fees"`
	Slippage        float64   `json:"slippage"`
	GrossR          float64   `json:"gross_r"`
	NetR            float64   `json:"net_r"`
	ExitReason      string    `json:"exit_reason"`
	HoldBars        int       `json:"hold_bars"`
}

type FuturesRangeFirstOccupancyRotationV1SummaryRow struct {
	ConfigID             string  `json:"config_id"`
	Family               string  `json:"family"`
	IsBaseline           bool    `json:"is_baseline"`
	Selected             bool    `json:"selected"`
	Timeframe            string  `json:"timeframe"`
	LookbackHours        int     `json:"lookback_hours"`
	MaxWidthPct          float64 `json:"max_width_pct"`
	OccupancyWindow      int     `json:"occupancy_window"`
	OccupancyZoneLevel   float64 `json:"occupancy_zone_level"`
	OccupancyMinFraction float64 `json:"occupancy_min_fraction"`
	RecaptureLevel       float64 `json:"recapture_level"`
	TargetLevel          float64 `json:"target_level"`
	MaxHoldBars          int     `json:"max_hold_bars"`
	StopBufferWidth      float64 `json:"stop_buffer_width"`
	SideMode             string  `json:"side_mode"`
	Split                string  `json:"split"`
	Side                 string  `json:"side"`
	SignalCount          int     `json:"signal_count"`
	SkippedSignalCount   int     `json:"skipped_signal_count"`
	TotalTrades          int     `json:"total_trades"`
	Wins                 int     `json:"wins"`
	Losses               int     `json:"losses"`
	WinRate              float64 `json:"win_rate"`
	GrossPnL             float64 `json:"gross_pnl"`
	NetPnL               float64 `json:"net_pnl"`
	TotalCosts           float64 `json:"total_costs"`
	ProfitFactor         float64 `json:"profit_factor"`
	GrossProfitFactor    float64 `json:"gross_profit_factor"`
	MaxDrawdown          float64 `json:"max_drawdown"`
	AvgGrossR            float64 `json:"avg_gross_r"`
	AvgNetR              float64 `json:"avg_net_r"`
	AvgInitialRisk       float64 `json:"avg_initial_risk"`
	AvgHoldBars          float64 `json:"avg_hold_bars"`
	IsWorstPeriodSplit   bool    `json:"is_worst_period_split"`
}

type FuturesRangeFirstOccupancyRotationV1GridRow struct {
	ConfigID                string  `json:"config_id"`
	Family                  string  `json:"family"`
	IsBaseline              bool    `json:"is_baseline"`
	Selected                bool    `json:"selected"`
	Timeframe               string  `json:"timeframe"`
	LookbackHours           int     `json:"lookback_hours"`
	MaxWidthPct             float64 `json:"max_width_pct"`
	OccupancyWindow         int     `json:"occupancy_window"`
	OccupancyZoneLevel      float64 `json:"occupancy_zone_level"`
	OccupancyMinFraction    float64 `json:"occupancy_min_fraction"`
	RecaptureLevel          float64 `json:"recapture_level"`
	TargetLevel             float64 `json:"target_level"`
	MaxHoldBars             int     `json:"max_hold_bars"`
	StopBufferWidth         float64 `json:"stop_buffer_width"`
	SideMode                string  `json:"side_mode"`
	SignalCount             int     `json:"signal_count"`
	SkippedSignalCount      int     `json:"skipped_signal_count"`
	TrainTrades             int     `json:"train_trades"`
	OOSTrades               int     `json:"oos_trades"`
	RecentTrades            int     `json:"recent_trades"`
	FullTrades              int     `json:"full_trades"`
	TrainNetPnL             float64 `json:"train_net_pnl"`
	OOSNetPnL               float64 `json:"oos_net_pnl"`
	RecentNetPnL            float64 `json:"recent_net_pnl"`
	FullNetPnL              float64 `json:"full_net_pnl"`
	TrainProfitFactor       float64 `json:"train_profit_factor"`
	OOSProfitFactor         float64 `json:"oos_profit_factor"`
	RecentProfitFactor      float64 `json:"recent_profit_factor"`
	FullProfitFactor        float64 `json:"full_profit_factor"`
	TrainMaxDrawdown        float64 `json:"train_max_drawdown"`
	FullMaxDrawdown         float64 `json:"full_max_drawdown"`
	TrainNetToDrawdown      float64 `json:"train_net_to_drawdown"`
	FullNetToDrawdown       float64 `json:"full_net_to_drawdown"`
	FullLongTrades          int     `json:"full_long_trades"`
	FullShortTrades         int     `json:"full_short_trades"`
	FullLongNetPnL          float64 `json:"full_long_net_pnl"`
	FullShortNetPnL         float64 `json:"full_short_net_pnl"`
	OOSLongNetPnL           float64 `json:"oos_long_net_pnl"`
	OOSShortNetPnL          float64 `json:"oos_short_net_pnl"`
	RecentLongNetPnL        float64 `json:"recent_long_net_pnl"`
	RecentShortNetPnL       float64 `json:"recent_short_net_pnl"`
	SideConcentration       float64 `json:"side_concentration"`
	SideWeaknessCaveat      bool    `json:"side_weakness_caveat"`
	SideConcentrationCaveat bool    `json:"side_concentration_caveat"`
	ThinSideCaveat          bool    `json:"thin_side_caveat"`
	CaveatPenalty           float64 `json:"caveat_penalty"`
	RankScore               float64 `json:"rank_score"`
	PassesGate              bool    `json:"passes_gate"`
	FailureReason           string  `json:"failure_reason,omitempty"`
}

type FuturesRangeFirstOccupancyRotationV1BaselineRow struct {
	ConfigID           string  `json:"config_id"`
	PassesGate         bool    `json:"passes_gate"`
	Rank               int     `json:"rank"`
	FullTrades         int     `json:"full_trades"`
	FullNetPnL         float64 `json:"full_net_pnl"`
	FullProfitFactor   float64 `json:"full_profit_factor"`
	FullMaxDrawdown    float64 `json:"full_max_drawdown"`
	TrainTrades        int     `json:"train_trades"`
	TrainNetPnL        float64 `json:"train_net_pnl"`
	TrainProfitFactor  float64 `json:"train_profit_factor"`
	OOSTrades          int     `json:"oos_trades"`
	OOSNetPnL          float64 `json:"oos_net_pnl"`
	OOSProfitFactor    float64 `json:"oos_profit_factor"`
	RecentTrades       int     `json:"recent_trades"`
	RecentNetPnL       float64 `json:"recent_net_pnl"`
	RecentProfitFactor float64 `json:"recent_profit_factor"`
	RankScore          float64 `json:"rank_score"`
	FailureReason      string  `json:"failure_reason,omitempty"`
}

type FuturesRangeFirstOccupancyRotationV1RankingRow struct {
	Rank                 int     `json:"rank"`
	ConfigID             string  `json:"config_id"`
	Family               string  `json:"family"`
	IsBaseline           bool    `json:"is_baseline"`
	Selected             bool    `json:"selected"`
	Timeframe            string  `json:"timeframe"`
	LookbackHours        int     `json:"lookback_hours"`
	MaxWidthPct          float64 `json:"max_width_pct"`
	OccupancyWindow      int     `json:"occupancy_window"`
	OccupancyMinFraction float64 `json:"occupancy_min_fraction"`
	RecaptureLevel       float64 `json:"recapture_level"`
	TargetLevel          float64 `json:"target_level"`
	MaxHoldBars          int     `json:"max_hold_bars"`
	StopBufferWidth      float64 `json:"stop_buffer_width"`
	TrainTrades          int     `json:"train_trades"`
	TrainNetPnL          float64 `json:"train_net_pnl"`
	TrainProfitFactor    float64 `json:"train_profit_factor"`
	TrainMaxDrawdown     float64 `json:"train_max_drawdown"`
	OOSTrades            int     `json:"oos_trades"`
	OOSNetPnL            float64 `json:"oos_net_pnl"`
	OOSProfitFactor      float64 `json:"oos_profit_factor"`
	RecentTrades         int     `json:"recent_trades"`
	RecentNetPnL         float64 `json:"recent_net_pnl"`
	RecentProfitFactor   float64 `json:"recent_profit_factor"`
	FullTrades           int     `json:"full_trades"`
	FullNetPnL           float64 `json:"full_net_pnl"`
	FullProfitFactor     float64 `json:"full_profit_factor"`
	FullMaxDrawdown      float64 `json:"full_max_drawdown"`
	CaveatPenalty        float64 `json:"caveat_penalty"`
	RankScore            float64 `json:"rank_score"`
	PassesGate           bool    `json:"passes_gate"`
	FailureReason        string  `json:"failure_reason,omitempty"`
}

type FuturesRangeFirstOccupancyRotationV1SelectionRow struct {
	Role               string  `json:"role"`
	ConfigID           string  `json:"config_id"`
	Rank               int     `json:"rank"`
	PassesGate         bool    `json:"passes_gate"`
	FullTrades         int     `json:"full_trades"`
	FullNetPnL         float64 `json:"full_net_pnl"`
	FullProfitFactor   float64 `json:"full_profit_factor"`
	FullMaxDrawdown    float64 `json:"full_max_drawdown"`
	TrainTrades        int     `json:"train_trades"`
	TrainNetPnL        float64 `json:"train_net_pnl"`
	TrainProfitFactor  float64 `json:"train_profit_factor"`
	OOSTrades          int     `json:"oos_trades"`
	OOSNetPnL          float64 `json:"oos_net_pnl"`
	OOSProfitFactor    float64 `json:"oos_profit_factor"`
	RecentTrades       int     `json:"recent_trades"`
	RecentNetPnL       float64 `json:"recent_net_pnl"`
	RecentProfitFactor float64 `json:"recent_profit_factor"`
	RankScore          float64 `json:"rank_score"`
	FailureReason      string  `json:"failure_reason,omitempty"`
	StopState          string  `json:"stop_state"`
}

type FuturesRangeFirstOccupancyRotationV1SkipRow struct {
	ConfigID   string `json:"config_id"`
	Family     string `json:"family"`
	IsBaseline bool   `json:"is_baseline"`
	Selected   bool   `json:"selected"`
	Timeframe  string `json:"timeframe"`
	Split      string `json:"split"`
	Side       string `json:"side"`
	Reason     string `json:"reason"`
	Count      int    `json:"count"`
}

type rangeFirstOccupancyFrameCache struct {
	candles         []Candle
	coverage        FuturesRangeFirstOccupancyRotationV1CoverageRow
	splitsByIndex   []string
	envelopes       map[int][]rangeFirstOccupancyEnvelope
	occupancyCounts map[string][]rangeFirstOccupancyCounts
}

type rangeFirstOccupancyEnvelope struct {
	valid      bool
	reason     string
	high       float64
	low        float64
	mid        float64
	q1         float64
	q3         float64
	width      float64
	widthPct   float64
	closePrice float64
}

type rangeFirstOccupancyCounts struct {
	inside int
	lower  int
	upper  int
}

type rangeFirstOccupancyStrategy struct {
	config      FuturesRangeFirstOccupancyRotationV1GridConfig
	signals     []FuturesRangeFirstOccupancyRotationV1SignalRow
	signalsByID map[string]FuturesRangeFirstOccupancyRotationV1SignalRow
	byIndex     map[int]Signal
}

type rangeFirstOccupancyConfigRun struct {
	signals []FuturesRangeFirstOccupancyRotationV1SignalRow
	skips   []FuturesRangeFirstOccupancyRotationV1SkipRow
	trades  []Trade
	rows    []FuturesRangeFirstOccupancyRotationV1TradeRow
	summary []FuturesRangeFirstOccupancyRotationV1SummaryRow
}

type rangeFirstOccupancySkipKey struct {
	configID  string
	timeframe string
	split     string
	side      string
	reason    string
}

func DefaultFuturesRangeFirstOccupancyRotationV1OptimizationConfig() FuturesRangeFirstOccupancyRotationV1OptimizationConfig {
	return FuturesRangeFirstOccupancyRotationV1OptimizationConfig{
		SourcePath:               rangeFirstOccupancyRotationSourcePath,
		ApprovedSourcePath:       rangeFirstOccupancyRotationSourcePath,
		ExpectedSourceRows:       rangeFirstOccupancyRotationExpectedRows,
		ExpectedFirstOpenTime:    rangeFirstOccupancyRotationExpectedFirst,
		ExpectedLastOpenTime:     rangeFirstOccupancyRotationExpectedLast,
		ExpectedGapCount:         0,
		ExpectedDuplicateCount:   0,
		ExpectedZeroVolumeCount:  rangeFirstOccupancyRotationExpectedZeroVol,
		Timeframes:               []string{RangeDiscoveryTimeframe15m, RangeDiscoveryTimeframe1h},
		Expected15MRows:          rangeFirstOccupancyRotationExpected15MRows,
		Expected15MLastOpenTime:  rangeFirstOccupancyRotationExpected15MLast,
		Expected1HRows:           rangeFirstOccupancyRotationExpected1HRows,
		Expected1HLastOpenTime:   rangeFirstOccupancyRotationExpected1HLast,
		LookbackHours:            []int{24, 48, 72},
		MaxWidthPcts:             []float64{0.020, 0.035},
		OccupancyWindows:         []int{8, 12},
		OccupancyZoneLevels:      []float64{0.25},
		OccupancyMinFractions:    []float64{0.60, 0.70},
		RecaptureLevels:          []float64{0.25, 0.33},
		TargetLevels:             []float64{0.50, 0.66},
		MaxHoldBars:              []int{8, 12, 24},
		StopBufferWidths:         []float64{0.00, 0.05},
		SideModes:                []string{RangeDiscoverySideAll},
		MinTrainTrades:           100,
		MinOOSTrades:             50,
		MinRecentTrades:          25,
		MinFullTrades:            200,
		MinTrainProfitFactor:     1.20,
		MinOOSProfitFactor:       1.05,
		MinRecentProfitFactor:    1.05,
		MinFullProfitFactor:      1.15,
		MinNetToDrawdown:         1.00,
		MinSideTradesForWeakness: 25,
		MaxSideNetConcentration:  0.75,
		SideConcentrationPenalty: 0.50,
		ThinSideTradeThreshold:   50,
		ThinSidePenalty:          0.50,
	}
}

func RunFuturesRangeFirstOccupancyRotationV1Optimization(parent []Candle, manifest SourceManifest, cfg FuturesRangeFirstOccupancyRotationV1OptimizationConfig, btCfg BacktestConfig, splits []Split) (FuturesRangeFirstOccupancyRotationV1OptimizationResult, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return FuturesRangeFirstOccupancyRotationV1OptimizationResult{}, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}

	result := FuturesRangeFirstOccupancyRotationV1OptimizationResult{
		BaselineConfigID: RangeFirstOccupancyRotationV1BaselineConfigID,
	}
	sourceRow, sourceErr := rangeFirstOccupancyRotationSourceRow(parent, manifest, cfg)
	result.SourceRows = append(result.SourceRows, sourceRow)
	if sourceErr != nil {
		result.StopState = RangeFirstStrategyV1StopStateSourceGap
		return result, sourceErr
	}

	frames, err := rangeFirstOccupancyRotationBuildFrames(parent, cfg, splits)
	if err != nil {
		result.CoverageRows = rangeFirstOccupancyRotationCoverageRows(frames)
		result.StopState = RangeFirstStrategyV1StopStateSourceGap
		return result, err
	}
	result.CoverageRows = rangeFirstOccupancyRotationCoverageRows(frames)

	gridConfigs := rangeFirstOccupancyRotationGridConfigs(cfg)
	baselineConfig := rangeFirstOccupancyRotationBaselineGridConfig()
	baselineRun, err := rangeFirstOccupancyRotationRunConfig(baselineConfig, frames[baselineConfig.Timeframe], btCfg, splits, true)
	if err != nil {
		return result, err
	}
	result.Trades = append([]Trade(nil), baselineRun.trades...)

	for _, grid := range gridConfigs {
		run, err := rangeFirstOccupancyRotationRunConfig(grid, frames[grid.Timeframe], btCfg, splits, false)
		if err != nil {
			return result, err
		}
		result.SummaryRows = append(result.SummaryRows, run.summary...)
		result.SkipRows = append(result.SkipRows, run.skips...)
	}
	rangeFirstOccupancyRotationMarkWorstSplits(result.SummaryRows)
	result.GridRows = rangeFirstOccupancyRotationGridRows(result.SummaryRows, result.SkipRows, gridConfigs, cfg)
	result.RankingRows = rangeFirstOccupancyRotationRankingRows(result.GridRows)
	for _, ranking := range result.RankingRows {
		if ranking.PassesGate {
			result.SelectedConfigID = ranking.ConfigID
			break
		}
	}
	rangeFirstOccupancyRotationMarkSelected(result.SelectedConfigID, result.GridRows, result.SummaryRows, result.RankingRows, result.SkipRows)

	selectedRun := rangeFirstOccupancyConfigRun{}
	if result.SelectedConfigID != "" {
		selectedGrid, _ := rangeFirstOccupancyRotationGridConfigByID(gridConfigs, result.SelectedConfigID)
		selectedRun, err = rangeFirstOccupancyRotationRunConfig(selectedGrid, frames[selectedGrid.Timeframe], btCfg, splits, true)
		if err != nil {
			return result, err
		}
	}
	baselineSelected := result.SelectedConfigID == baselineConfig.ConfigID
	for i := range baselineRun.signals {
		baselineRun.signals[i].IsBaseline = true
		baselineRun.signals[i].Selected = baselineSelected
	}
	for i := range baselineRun.rows {
		baselineRun.rows[i].IsBaseline = true
		baselineRun.rows[i].Selected = baselineSelected
	}
	result.SignalRows = append(result.SignalRows, baselineRun.signals...)
	result.TradeRows = append(result.TradeRows, baselineRun.rows...)
	if result.SelectedConfigID != "" && !baselineSelected {
		for i := range selectedRun.signals {
			selectedRun.signals[i].Selected = true
		}
		for i := range selectedRun.rows {
			selectedRun.rows[i].Selected = true
		}
		result.SignalRows = append(result.SignalRows, selectedRun.signals...)
		result.TradeRows = append(result.TradeRows, selectedRun.rows...)
	}
	sort.Slice(result.SignalRows, func(i, j int) bool {
		if result.SignalRows[i].ConfigID != result.SignalRows[j].ConfigID {
			return result.SignalRows[i].ConfigID < result.SignalRows[j].ConfigID
		}
		return result.SignalRows[i].SignalIndex < result.SignalRows[j].SignalIndex
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

	result.StopState = FuturesRangeFirstOccupancyRotationV1OptimizationStopState(result.SourceRows, result.CoverageRows, result.GridRows, result.RankingRows)
	result.BaselineRows = rangeFirstOccupancyRotationBaselineRows(result.GridRows, result.RankingRows)
	result.SelectionRows = rangeFirstOccupancyRotationSelectionRows(result.GridRows, result.RankingRows, result.BaselineConfigID, result.SelectedConfigID, result.StopState)
	return result, nil
}

func FuturesRangeFirstOccupancyRotationV1OptimizationStopState(source []FuturesRangeFirstOccupancyRotationV1SourceRow, coverage []FuturesRangeFirstOccupancyRotationV1CoverageRow, grid []FuturesRangeFirstOccupancyRotationV1GridRow, rankings []FuturesRangeFirstOccupancyRotationV1RankingRow) string {
	if len(source) == 0 || !source[0].SourceFactsPass || source[0].ValidationStatus != "accepted" {
		return RangeFirstStrategyV1StopStateSourceGap
	}
	for _, row := range coverage {
		if !row.CoverageFactsPass || !row.Complete || row.ValidationStatus != "accepted" {
			return RangeFirstStrategyV1StopStateSourceGap
		}
	}
	validSignals := 0
	for _, row := range grid {
		validSignals += row.SignalCount
	}
	if validSignals == 0 {
		return RangeFirstStrategyV1StopStateNoValidSignals
	}
	for _, row := range rankings {
		if row.PassesGate {
			return RangeFirstStrategyV1StopStatePassedNeedsFixedReplaySpec
		}
	}
	return RangeFirstStrategyV1StopStateOptimizerFailedNoReplay
}

func (cfg FuturesRangeFirstOccupancyRotationV1OptimizationConfig) withDefaults() FuturesRangeFirstOccupancyRotationV1OptimizationConfig {
	defaults := DefaultFuturesRangeFirstOccupancyRotationV1OptimizationConfig()
	if cfg.SourcePath == "" {
		cfg.SourcePath = defaults.SourcePath
	}
	if cfg.ApprovedSourcePath == "" {
		cfg.ApprovedSourcePath = defaults.ApprovedSourcePath
	}
	if cfg.ExpectedSourceRows == 0 {
		cfg.ExpectedSourceRows = defaults.ExpectedSourceRows
	}
	if cfg.ExpectedFirstOpenTime == "" {
		cfg.ExpectedFirstOpenTime = defaults.ExpectedFirstOpenTime
	}
	if cfg.ExpectedLastOpenTime == "" {
		cfg.ExpectedLastOpenTime = defaults.ExpectedLastOpenTime
	}
	if cfg.ExpectedZeroVolumeCount == 0 {
		cfg.ExpectedZeroVolumeCount = defaults.ExpectedZeroVolumeCount
	}
	if len(cfg.Timeframes) == 0 {
		cfg.Timeframes = append([]string(nil), defaults.Timeframes...)
	}
	if cfg.Expected15MRows == 0 {
		cfg.Expected15MRows = defaults.Expected15MRows
	}
	if cfg.Expected15MLastOpenTime == "" {
		cfg.Expected15MLastOpenTime = defaults.Expected15MLastOpenTime
	}
	if cfg.Expected1HRows == 0 {
		cfg.Expected1HRows = defaults.Expected1HRows
	}
	if cfg.Expected1HLastOpenTime == "" {
		cfg.Expected1HLastOpenTime = defaults.Expected1HLastOpenTime
	}
	if len(cfg.LookbackHours) == 0 {
		cfg.LookbackHours = append([]int(nil), defaults.LookbackHours...)
	}
	if len(cfg.MaxWidthPcts) == 0 {
		cfg.MaxWidthPcts = append([]float64(nil), defaults.MaxWidthPcts...)
	}
	if len(cfg.OccupancyWindows) == 0 {
		cfg.OccupancyWindows = append([]int(nil), defaults.OccupancyWindows...)
	}
	if len(cfg.OccupancyZoneLevels) == 0 {
		cfg.OccupancyZoneLevels = append([]float64(nil), defaults.OccupancyZoneLevels...)
	}
	if len(cfg.OccupancyMinFractions) == 0 {
		cfg.OccupancyMinFractions = append([]float64(nil), defaults.OccupancyMinFractions...)
	}
	if len(cfg.RecaptureLevels) == 0 {
		cfg.RecaptureLevels = append([]float64(nil), defaults.RecaptureLevels...)
	}
	if len(cfg.TargetLevels) == 0 {
		cfg.TargetLevels = append([]float64(nil), defaults.TargetLevels...)
	}
	if len(cfg.MaxHoldBars) == 0 {
		cfg.MaxHoldBars = append([]int(nil), defaults.MaxHoldBars...)
	}
	if len(cfg.StopBufferWidths) == 0 {
		cfg.StopBufferWidths = append([]float64(nil), defaults.StopBufferWidths...)
	}
	if len(cfg.SideModes) == 0 {
		cfg.SideModes = append([]string(nil), defaults.SideModes...)
	}
	if cfg.MinTrainTrades == 0 {
		cfg.MinTrainTrades = defaults.MinTrainTrades
	}
	if cfg.MinOOSTrades == 0 {
		cfg.MinOOSTrades = defaults.MinOOSTrades
	}
	if cfg.MinRecentTrades == 0 {
		cfg.MinRecentTrades = defaults.MinRecentTrades
	}
	if cfg.MinFullTrades == 0 {
		cfg.MinFullTrades = defaults.MinFullTrades
	}
	if cfg.MinTrainProfitFactor == 0 {
		cfg.MinTrainProfitFactor = defaults.MinTrainProfitFactor
	}
	if cfg.MinOOSProfitFactor == 0 {
		cfg.MinOOSProfitFactor = defaults.MinOOSProfitFactor
	}
	if cfg.MinRecentProfitFactor == 0 {
		cfg.MinRecentProfitFactor = defaults.MinRecentProfitFactor
	}
	if cfg.MinFullProfitFactor == 0 {
		cfg.MinFullProfitFactor = defaults.MinFullProfitFactor
	}
	if cfg.MinNetToDrawdown == 0 {
		cfg.MinNetToDrawdown = defaults.MinNetToDrawdown
	}
	if cfg.MinSideTradesForWeakness == 0 {
		cfg.MinSideTradesForWeakness = defaults.MinSideTradesForWeakness
	}
	if cfg.MaxSideNetConcentration == 0 {
		cfg.MaxSideNetConcentration = defaults.MaxSideNetConcentration
	}
	if cfg.SideConcentrationPenalty == 0 {
		cfg.SideConcentrationPenalty = defaults.SideConcentrationPenalty
	}
	if cfg.ThinSideTradeThreshold == 0 {
		cfg.ThinSideTradeThreshold = defaults.ThinSideTradeThreshold
	}
	if cfg.ThinSidePenalty == 0 {
		cfg.ThinSidePenalty = defaults.ThinSidePenalty
	}
	return cfg
}

func (cfg FuturesRangeFirstOccupancyRotationV1OptimizationConfig) validate() error {
	if cfg.SourcePath == "" || cfg.ApprovedSourcePath == "" {
		return fmt.Errorf("range occupancy rotation source paths must not be empty")
	}
	for _, timeframe := range cfg.Timeframes {
		if timeframe != RangeDiscoveryTimeframe15m && timeframe != RangeDiscoveryTimeframe1h {
			return fmt.Errorf("range occupancy rotation unsupported timeframe %q", timeframe)
		}
	}
	for _, value := range cfg.LookbackHours {
		if value <= 0 {
			return fmt.Errorf("range occupancy rotation lookback hours must be positive")
		}
	}
	for _, value := range cfg.MaxWidthPcts {
		if value <= 0 {
			return fmt.Errorf("range occupancy rotation max width pct must be positive")
		}
	}
	for _, value := range cfg.OccupancyWindows {
		if value <= 0 {
			return fmt.Errorf("range occupancy rotation occupancy windows must be positive")
		}
	}
	for _, value := range cfg.OccupancyZoneLevels {
		if value <= 0 || value >= 0.5 {
			return fmt.Errorf("range occupancy rotation occupancy zone levels must be in (0,0.5)")
		}
	}
	for _, value := range cfg.OccupancyMinFractions {
		if value <= 0 || value > 1 {
			return fmt.Errorf("range occupancy rotation occupancy min fractions must be in (0,1]")
		}
	}
	for _, value := range cfg.RecaptureLevels {
		if value <= 0 || value >= 0.5 {
			return fmt.Errorf("range occupancy rotation recapture levels must be in (0,0.5)")
		}
	}
	for _, value := range cfg.TargetLevels {
		if value <= 0 || value >= 1 {
			return fmt.Errorf("range occupancy rotation target levels must be in (0,1)")
		}
	}
	for _, value := range cfg.MaxHoldBars {
		if value <= 0 {
			return fmt.Errorf("range occupancy rotation max hold bars must be positive")
		}
	}
	for _, value := range cfg.StopBufferWidths {
		if value < 0 {
			return fmt.Errorf("range occupancy rotation stop buffers must be non-negative")
		}
	}
	for _, sideMode := range cfg.SideModes {
		if sideMode != RangeDiscoverySideAll {
			return fmt.Errorf("range occupancy rotation unsupported side mode %q", sideMode)
		}
	}
	if cfg.MinTrainTrades <= 0 || cfg.MinOOSTrades <= 0 || cfg.MinRecentTrades <= 0 || cfg.MinFullTrades <= 0 {
		return fmt.Errorf("range occupancy rotation trade gates must be positive")
	}
	if cfg.MinTrainProfitFactor <= 0 || cfg.MinOOSProfitFactor <= 0 || cfg.MinRecentProfitFactor <= 0 || cfg.MinFullProfitFactor <= 0 {
		return fmt.Errorf("range occupancy rotation profit factor gates must be positive")
	}
	if cfg.MinNetToDrawdown <= 0 || cfg.MinSideTradesForWeakness <= 0 || cfg.MaxSideNetConcentration <= 0 || cfg.SideConcentrationPenalty <= 0 || cfg.ThinSideTradeThreshold <= 0 || cfg.ThinSidePenalty <= 0 {
		return fmt.Errorf("range occupancy rotation review gates must be positive")
	}
	return nil
}

func rangeFirstOccupancyRotationSourceRow(parent []Candle, manifest SourceManifest, cfg FuturesRangeFirstOccupancyRotationV1OptimizationConfig) (FuturesRangeFirstOccupancyRotationV1SourceRow, error) {
	row := FuturesRangeFirstOccupancyRotationV1SourceRow{
		Path:                    manifest.Path,
		ApprovedPath:            cfg.ApprovedSourcePath,
		Venue:                   manifest.Venue,
		Product:                 manifest.Product,
		Symbol:                  manifest.Symbol,
		Interval:                manifest.Interval,
		RowCount:                manifest.RowCount,
		ExpectedRowCount:        cfg.ExpectedSourceRows,
		FirstOpenTime:           manifest.FirstOpenTime,
		ExpectedFirstOpenTime:   cfg.ExpectedFirstOpenTime,
		LastOpenTime:            manifest.LastOpenTime,
		ExpectedLastOpenTime:    cfg.ExpectedLastOpenTime,
		GapCount:                manifest.GapCount,
		ExpectedGapCount:        cfg.ExpectedGapCount,
		DuplicateCount:          manifest.DuplicateCount,
		ExpectedDuplicateCount:  cfg.ExpectedDuplicateCount,
		ZeroVolumeCount:         manifest.ZeroVolumeCount,
		ExpectedZeroVolumeCount: cfg.ExpectedZeroVolumeCount,
		ComparisonOnly:          manifest.ComparisonOnly,
		ValidationStatus:        manifest.ValidationStatus,
		SourceFactsPass:         true,
	}
	reject := func(format string, args ...any) (FuturesRangeFirstOccupancyRotationV1SourceRow, error) {
		row.SourceFactsPass = false
		row.ValidationStatus = "rejected"
		row.ValidationError = fmt.Sprintf(format, args...)
		return row, fmt.Errorf("%s", row.ValidationError)
	}
	if manifest.Product != "Binance USDT-M futures" || manifest.ComparisonOnly {
		return reject("range occupancy rotation requires Binance USDT-M futures source; got product=%q comparison_only=%t", manifest.Product, manifest.ComparisonOnly)
	}
	if manifest.Symbol != RangeUniverseSymbolBTCUSDT || manifest.Interval != "5m" {
		return reject("range occupancy rotation requires BTCUSDT 5m source; got symbol=%q interval=%q", manifest.Symbol, manifest.Interval)
	}
	if manifest.ValidationStatus != "accepted" {
		return reject("range occupancy rotation source manifest status is %q", manifest.ValidationStatus)
	}
	if !cfg.SkipSourceFactCheck {
		if manifest.Path != cfg.ApprovedSourcePath {
			return reject("range occupancy rotation source path %q does not match approved path %q", manifest.Path, cfg.ApprovedSourcePath)
		}
		if len(parent) != cfg.ExpectedSourceRows || manifest.RowCount != cfg.ExpectedSourceRows {
			return reject("range occupancy rotation source rows=%d manifest_rows=%d expected=%d", len(parent), manifest.RowCount, cfg.ExpectedSourceRows)
		}
		if manifest.FirstOpenTime != cfg.ExpectedFirstOpenTime || manifest.LastOpenTime != cfg.ExpectedLastOpenTime {
			return reject("range occupancy rotation source coverage %s..%s expected %s..%s", manifest.FirstOpenTime, manifest.LastOpenTime, cfg.ExpectedFirstOpenTime, cfg.ExpectedLastOpenTime)
		}
		if manifest.GapCount != cfg.ExpectedGapCount || manifest.DuplicateCount != cfg.ExpectedDuplicateCount || manifest.ZeroVolumeCount != cfg.ExpectedZeroVolumeCount {
			return reject("range occupancy rotation source facts gap=%d duplicate=%d zero_volume=%d expected gap=%d duplicate=%d zero_volume=%d", manifest.GapCount, manifest.DuplicateCount, manifest.ZeroVolumeCount, cfg.ExpectedGapCount, cfg.ExpectedDuplicateCount, cfg.ExpectedZeroVolumeCount)
		}
	}
	return row, nil
}

func rangeFirstOccupancyRotationBuildFrames(parent []Candle, cfg FuturesRangeFirstOccupancyRotationV1OptimizationConfig, splits []Split) (map[string]rangeFirstOccupancyFrameCache, error) {
	out := map[string]rangeFirstOccupancyFrameCache{}
	for _, timeframe := range cfg.Timeframes {
		frame, ok := rangeFirstOccupancyRotationFrameDef(timeframe)
		if !ok {
			return out, fmt.Errorf("range occupancy rotation missing frame definition for %s", timeframe)
		}
		candles, coverage, err := resampleRangeDiscoveryFrame(parent, frame)
		row := rangeFirstOccupancyRotationCoverageRow(coverage, cfg)
		cache := rangeFirstOccupancyFrameCache{
			candles:         candles,
			coverage:        row,
			splitsByIndex:   make([]string, len(candles)),
			envelopes:       map[int][]rangeFirstOccupancyEnvelope{},
			occupancyCounts: map[string][]rangeFirstOccupancyCounts{},
		}
		for i, candle := range candles {
			cache.splitsByIndex[i] = splitNameForCloseTime(candle.CloseTime, splits)
		}
		out[timeframe] = cache
		if err != nil {
			return out, err
		}
		if !row.CoverageFactsPass || !row.Complete || row.ValidationStatus != "accepted" {
			return out, fmt.Errorf("range occupancy rotation %s resample rejected: %s", timeframe, row.ValidationError)
		}
	}
	for timeframe, cache := range out {
		for _, lookback := range cfg.LookbackHours {
			envelopes, err := rangeFirstOccupancyRotationPrecomputeEnvelopes(cache.candles, timeframe, lookback)
			if err != nil {
				return out, err
			}
			cache.envelopes[lookback] = envelopes
			for _, window := range cfg.OccupancyWindows {
				cache.occupancyCounts[rangeFirstOccupancyCountsKey(lookback, window)] = rangeFirstOccupancyRotationPrecomputeOccupancy(cache.candles, envelopes, window, cfg.OccupancyZoneLevels[0])
			}
		}
		out[timeframe] = cache
	}
	return out, nil
}

func rangeFirstOccupancyRotationCoverageRows(frames map[string]rangeFirstOccupancyFrameCache) []FuturesRangeFirstOccupancyRotationV1CoverageRow {
	rows := make([]FuturesRangeFirstOccupancyRotationV1CoverageRow, 0, len(frames))
	for _, frame := range frames {
		rows = append(rows, frame.coverage)
	}
	sort.Slice(rows, func(i, j int) bool {
		return rangeFirstOccupancyRotationTimeframeSortKey(rows[i].Timeframe) < rangeFirstOccupancyRotationTimeframeSortKey(rows[j].Timeframe)
	})
	return rows
}

func rangeFirstOccupancyRotationCoverageRow(base FuturesRangeDiscoveryCoverageRow, cfg FuturesRangeFirstOccupancyRotationV1OptimizationConfig) FuturesRangeFirstOccupancyRotationV1CoverageRow {
	row := FuturesRangeFirstOccupancyRotationV1CoverageRow{
		FuturesRangeDiscoveryCoverageRow: base,
		ExpectedFirstOpenTime:            rangeFirstOccupancyRotationExpectedFirst,
		CoverageFactsPass:                true,
	}
	switch base.Timeframe {
	case RangeDiscoveryTimeframe15m:
		row.ExpectedRowCount = cfg.Expected15MRows
		row.ExpectedLastOpenTime = cfg.Expected15MLastOpenTime
	case RangeDiscoveryTimeframe1h:
		row.ExpectedRowCount = cfg.Expected1HRows
		row.ExpectedLastOpenTime = cfg.Expected1HLastOpenTime
	}
	if base.ValidationStatus != "accepted" || !base.Complete {
		row.CoverageFactsPass = false
		return row
	}
	if !cfg.SkipCoverageCountCheck {
		if row.RowCount != row.ExpectedRowCount || row.FirstOpenTime != row.ExpectedFirstOpenTime || row.LastOpenTime != row.ExpectedLastOpenTime {
			row.CoverageFactsPass = false
			row.ValidationStatus = "rejected"
			row.ValidationError = fmt.Sprintf("%s coverage row_count=%d first=%s last=%s expected row_count=%d first=%s last=%s", base.Timeframe, row.RowCount, row.FirstOpenTime, row.LastOpenTime, row.ExpectedRowCount, row.ExpectedFirstOpenTime, row.ExpectedLastOpenTime)
		}
	}
	if row.GapCount != 0 || row.DuplicateCount != 0 || row.MissingChildOpenCount != 0 {
		row.CoverageFactsPass = false
		if row.ValidationError == "" {
			row.ValidationStatus = "rejected"
			row.ValidationError = fmt.Sprintf("%s coverage gap=%d duplicate=%d missing_child_open=%d", base.Timeframe, row.GapCount, row.DuplicateCount, row.MissingChildOpenCount)
		}
	}
	return row
}

func rangeFirstOccupancyRotationFrameDef(timeframe string) (rangeDiscoveryFrameDef, bool) {
	for _, frame := range rangeDiscoveryFrameDefs() {
		if frame.timeframe == timeframe {
			return frame, true
		}
	}
	return rangeDiscoveryFrameDef{}, false
}

func rangeFirstOccupancyRotationPrecomputeEnvelopes(candles []Candle, timeframe string, lookbackHours int) ([]rangeFirstOccupancyEnvelope, error) {
	bars, err := rangeFirstOccupancyRotationLookbackBars(timeframe, lookbackHours)
	if err != nil {
		return nil, err
	}
	out := make([]rangeFirstOccupancyEnvelope, len(candles))
	for i := range candles {
		if i < bars {
			out[i] = rangeFirstOccupancyEnvelope{reason: "insufficient_history"}
			continue
		}
		closePrice := candles[i].Close
		if closePrice <= 0 {
			out[i] = rangeFirstOccupancyEnvelope{reason: "non_positive_price", closePrice: closePrice}
			continue
		}
		high := candles[i-bars].High
		low := candles[i-bars].Low
		if high <= 0 || low <= 0 {
			out[i] = rangeFirstOccupancyEnvelope{reason: "non_positive_price", closePrice: closePrice}
			continue
		}
		for j := i - bars; j < i; j++ {
			if candles[j].High <= 0 || candles[j].Low <= 0 || candles[j].Close <= 0 {
				out[i] = rangeFirstOccupancyEnvelope{reason: "non_positive_price", closePrice: closePrice}
				break
			}
			if candles[j].High > high {
				high = candles[j].High
			}
			if candles[j].Low < low {
				low = candles[j].Low
			}
		}
		if out[i].reason != "" {
			continue
		}
		width := high - low
		if width <= 0 {
			out[i] = rangeFirstOccupancyEnvelope{reason: "non_positive_range_width", high: high, low: low, closePrice: closePrice}
			continue
		}
		out[i] = rangeFirstOccupancyEnvelope{
			valid:      true,
			high:       high,
			low:        low,
			mid:        low + 0.50*width,
			q1:         low + 0.25*width,
			q3:         low + 0.75*width,
			width:      width,
			widthPct:   width / closePrice,
			closePrice: closePrice,
		}
	}
	return out, nil
}

func rangeFirstOccupancyRotationPrecomputeOccupancy(candles []Candle, envelopes []rangeFirstOccupancyEnvelope, window int, zoneLevel float64) []rangeFirstOccupancyCounts {
	out := make([]rangeFirstOccupancyCounts, len(candles))
	for i := range candles {
		if i+1 < window || i >= len(envelopes) || !envelopes[i].valid {
			continue
		}
		env := envelopes[i]
		lowerLimit := env.low + zoneLevel*env.width
		upperLimit := env.high - zoneLevel*env.width
		for j := i - window + 1; j <= i; j++ {
			closePrice := candles[j].Close
			if closePrice >= env.low && closePrice <= env.high {
				out[i].inside++
			}
			if closePrice <= lowerLimit {
				out[i].lower++
			}
			if closePrice >= upperLimit {
				out[i].upper++
			}
		}
	}
	return out
}

func rangeFirstOccupancyRotationGridConfigs(cfg FuturesRangeFirstOccupancyRotationV1OptimizationConfig) []FuturesRangeFirstOccupancyRotationV1GridConfig {
	cfg = cfg.withDefaults()
	rows := []FuturesRangeFirstOccupancyRotationV1GridConfig{}
	for _, timeframe := range cfg.Timeframes {
		for _, lookback := range cfg.LookbackHours {
			for _, maxWidth := range cfg.MaxWidthPcts {
				for _, occupancyWindow := range cfg.OccupancyWindows {
					for _, occupancyZone := range cfg.OccupancyZoneLevels {
						for _, occupancyMin := range cfg.OccupancyMinFractions {
							for _, recapture := range cfg.RecaptureLevels {
								for _, target := range cfg.TargetLevels {
									for _, maxHold := range cfg.MaxHoldBars {
										for _, stopBuffer := range cfg.StopBufferWidths {
											for _, sideMode := range cfg.SideModes {
												grid := FuturesRangeFirstOccupancyRotationV1GridConfig{
													Timeframe:            timeframe,
													LookbackHours:        lookback,
													MaxWidthPct:          maxWidth,
													OccupancyWindow:      occupancyWindow,
													OccupancyZoneLevel:   occupancyZone,
													OccupancyMinFraction: occupancyMin,
													RecaptureLevel:       recapture,
													TargetLevel:          target,
													MaxHoldBars:          maxHold,
													StopBufferWidth:      stopBuffer,
													SideMode:             sideMode,
												}
												grid.ConfigID = rangeFirstOccupancyRotationConfigID(grid)
												rows = append(rows, grid)
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return rows
}

func rangeFirstOccupancyRotationBaselineGridConfig() FuturesRangeFirstOccupancyRotationV1GridConfig {
	return FuturesRangeFirstOccupancyRotationV1GridConfig{
		ConfigID:             RangeFirstOccupancyRotationV1BaselineConfigID,
		Timeframe:            RangeDiscoveryTimeframe1h,
		LookbackHours:        48,
		MaxWidthPct:          0.035,
		OccupancyWindow:      12,
		OccupancyZoneLevel:   0.25,
		OccupancyMinFraction: 0.60,
		RecaptureLevel:       0.33,
		TargetLevel:          0.66,
		MaxHoldBars:          12,
		StopBufferWidth:      0.05,
		SideMode:             RangeDiscoverySideAll,
	}
}

func rangeFirstOccupancyRotationConfigID(grid FuturesRangeFirstOccupancyRotationV1GridConfig) string {
	return fmt.Sprintf("range_occupancy_rotation_v1_%s_l%d_w%s_ow%d_occ%s_rec%s_t%s_h%d_sb%s",
		grid.Timeframe,
		grid.LookbackHours,
		rangeFirstOccupancyRotationMilliID(grid.MaxWidthPct),
		grid.OccupancyWindow,
		rangeFirstOccupancyRotationPercent3ID(grid.OccupancyMinFraction),
		rangeFirstOccupancyRotationPercentID(grid.RecaptureLevel),
		rangeFirstOccupancyRotationPercentID(grid.TargetLevel),
		grid.MaxHoldBars,
		rangeFirstOccupancyRotationPercent3ID(grid.StopBufferWidth),
	)
}

func rangeFirstOccupancyRotationMilliID(value float64) string {
	return fmt.Sprintf("%03d", int(math.Round(value*1000)))
}

func rangeFirstOccupancyRotationPercentID(value float64) string {
	return fmt.Sprintf("%02d", int(math.Round(value*100)))
}

func rangeFirstOccupancyRotationPercent3ID(value float64) string {
	return fmt.Sprintf("%03d", int(math.Round(value*100)))
}

func rangeFirstOccupancyRotationRunConfig(grid FuturesRangeFirstOccupancyRotationV1GridConfig, frame rangeFirstOccupancyFrameCache, btCfg BacktestConfig, splits []Split, keepDetails bool) (rangeFirstOccupancyConfigRun, error) {
	strategy, skips, err := rangeFirstOccupancyRotationBuildStrategy(grid, frame, btCfg, splits)
	if err != nil {
		return rangeFirstOccupancyConfigRun{}, err
	}
	run := RunBacktest(frame.candles, strategy, btCfg)
	signals := strategy.SignalRows()
	rangeFirstOccupancyRotationMarkExecutedSignals(signals, run.Trades, &skips, grid, frame)
	tradeRows := []FuturesRangeFirstOccupancyRotationV1TradeRow{}
	if keepDetails {
		tradeRows = strategy.TradeRows(run.Trades, splits)
	}
	summary := rangeFirstOccupancyRotationSummarizeConfig(grid, signals, skips, run.Trades, btCfg.StartBalance, splits)
	if !keepDetails {
		signals = nil
	}
	return rangeFirstOccupancyConfigRun{
		signals: signals,
		skips:   skips,
		trades:  run.Trades,
		rows:    tradeRows,
		summary: summary,
	}, nil
}

func rangeFirstOccupancyRotationBuildStrategy(grid FuturesRangeFirstOccupancyRotationV1GridConfig, frame rangeFirstOccupancyFrameCache, btCfg BacktestConfig, splits []Split) (rangeFirstOccupancyStrategy, []FuturesRangeFirstOccupancyRotationV1SkipRow, error) {
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	envelopes, ok := frame.envelopes[grid.LookbackHours]
	if !ok {
		return rangeFirstOccupancyStrategy{}, nil, fmt.Errorf("range occupancy rotation missing envelopes for %s lookback %d", grid.Timeframe, grid.LookbackHours)
	}
	occupancy, ok := frame.occupancyCounts[rangeFirstOccupancyCountsKey(grid.LookbackHours, grid.OccupancyWindow)]
	if !ok {
		return rangeFirstOccupancyStrategy{}, nil, fmt.Errorf("range occupancy rotation missing occupancy counts for %s lookback %d window %d", grid.Timeframe, grid.LookbackHours, grid.OccupancyWindow)
	}
	strategy := rangeFirstOccupancyStrategy{
		config:      grid,
		signalsByID: map[string]FuturesRangeFirstOccupancyRotationV1SignalRow{},
		byIndex:     map[int]Signal{},
	}
	skipCounts := map[rangeFirstOccupancySkipKey]int{}
	addSkip := func(index int, side string, reason string) {
		split := fullSplitName
		if index >= 0 && index < len(frame.splitsByIndex) {
			split = frame.splitsByIndex[index]
		}
		for _, splitName := range rangeFirstOccupancyRotationSplitNames(split) {
			skipCounts[rangeFirstOccupancySkipKey{configID: grid.ConfigID, timeframe: grid.Timeframe, split: splitName, side: side, reason: reason}]++
		}
	}
	requiredCount := int(math.Ceil(grid.OccupancyMinFraction * float64(grid.OccupancyWindow)))
	for i := range frame.candles {
		env := envelopes[i]
		if !env.valid {
			reason := env.reason
			if reason == "" {
				reason = "insufficient_history"
			}
			addSkip(i, "all", reason)
			continue
		}
		if env.widthPct > grid.MaxWidthPct {
			addSkip(i, "all", "range_too_wide")
			continue
		}
		if i+1 >= len(frame.candles) {
			addSkip(i, "all", "missing_entry_bar")
			continue
		}
		counts := occupancy[i]
		if counts.inside < grid.OccupancyWindow {
			addSkip(i, "all", "recent_close_outside_range")
			continue
		}
		signalCandle := frame.candles[i]
		lowerRecapture := env.low + grid.RecaptureLevel*env.width
		upperRecapture := env.high - grid.RecaptureLevel*env.width
		longPrelim := true
		shortPrelim := true
		if counts.lower < requiredCount {
			addSkip(i, string(Long), "occupancy_threshold_not_met")
			longPrelim = false
		}
		if counts.upper < requiredCount {
			addSkip(i, string(Short), "occupancy_threshold_not_met")
			shortPrelim = false
		}
		if longPrelim && signalCandle.Close < lowerRecapture {
			addSkip(i, string(Long), "recapture_not_confirmed")
			longPrelim = false
		}
		if shortPrelim && signalCandle.Close > upperRecapture {
			addSkip(i, string(Short), "recapture_not_confirmed")
			shortPrelim = false
		}
		if longPrelim && shortPrelim {
			addSkip(i, "all", "ambiguous_dual_side_signal")
			continue
		}
		if longPrelim && signalCandle.Close >= env.mid {
			addSkip(i, string(Long), "signal_crossed_midpoint")
			longPrelim = false
		}
		if shortPrelim && signalCandle.Close <= env.mid {
			addSkip(i, string(Short), "signal_crossed_midpoint")
			shortPrelim = false
		}
		if !longPrelim && !shortPrelim {
			continue
		}
		side := Long
		if shortPrelim {
			side = Short
		}
		row := rangeFirstOccupancyRotationSignalRow(frame, grid, env, counts, i, side, requiredCount, lowerRecapture, upperRecapture, btCfg)
		if row.SkippedReason != "" {
			addSkip(i, string(side), row.SkippedReason)
			continue
		}
		strategy.byIndex[i] = Signal{Side: row.Side, Stop: row.Stop, Target: row.Target, MaxHoldBars: row.MaxHoldBars, Reason: row.SignalID}
		strategy.signalsByID[row.SignalID] = row
		strategy.signals = append(strategy.signals, row)
	}
	return strategy, rangeFirstOccupancyRotationSkipRowsFromCounts(skipCounts, grid, false, false), nil
}

func rangeFirstOccupancyRotationSignalRow(frame rangeFirstOccupancyFrameCache, grid FuturesRangeFirstOccupancyRotationV1GridConfig, env rangeFirstOccupancyEnvelope, counts rangeFirstOccupancyCounts, index int, side Direction, requiredCount int, lowerRecapture float64, upperRecapture float64, btCfg BacktestConfig) FuturesRangeFirstOccupancyRotationV1SignalRow {
	signal := frame.candles[index]
	entryIndex := index + 1
	signalID := fmt.Sprintf("%s_%06d", grid.ConfigID, len(frame.candles)+index)
	row := FuturesRangeFirstOccupancyRotationV1SignalRow{
		ConfigID:               grid.ConfigID,
		Family:                 rangeFirstOccupancyRotationFamily,
		IsBaseline:             grid.ConfigID == RangeFirstOccupancyRotationV1BaselineConfigID,
		SignalID:               signalID,
		Timeframe:              grid.Timeframe,
		Split:                  frame.splitsByIndex[index],
		SignalIndex:            index,
		SignalOpenTime:         signal.OpenTime.UTC().Format(timeLayout),
		SignalCloseTime:        signal.CloseTime.UTC().Format(timeLayout),
		SignalOpen:             signal.Open,
		SignalHigh:             signal.High,
		SignalLow:              signal.Low,
		SignalClose:            signal.Close,
		LookbackHours:          grid.LookbackHours,
		MaxWidthPct:            grid.MaxWidthPct,
		OccupancyWindow:        grid.OccupancyWindow,
		OccupancyZoneLevel:     grid.OccupancyZoneLevel,
		OccupancyMinFraction:   grid.OccupancyMinFraction,
		RequiredOccupancyCount: requiredCount,
		RecaptureLevel:         grid.RecaptureLevel,
		TargetLevel:            grid.TargetLevel,
		StopBufferWidth:        grid.StopBufferWidth,
		SideMode:               grid.SideMode,
		RangeHigh:              env.high,
		RangeLow:               env.low,
		RangeMid:               env.mid,
		RangeQ1:                env.q1,
		RangeQ3:                env.q3,
		RangeWidth:             env.width,
		RangeWidthPct:          env.widthPct,
		LowerRecapture:         lowerRecapture,
		UpperRecapture:         upperRecapture,
		InsideOccupancyCount:   counts.inside,
		LowerOccupancyCount:    counts.lower,
		UpperOccupancyCount:    counts.upper,
		Side:                   side,
		EntryIndex:             entryIndex,
		MaxHoldBars:            grid.MaxHoldBars,
	}
	if entryIndex >= len(frame.candles) {
		row.SkippedReason = "missing_entry_bar"
		return row
	}
	entry := frame.candles[entryIndex]
	row.EntryOpenTime = entry.OpenTime.UTC().Format(timeLayout)
	row.EntryOpen = entry.Open
	row.ExpectedEntryPrice = applySlippage(entry.Open, btCfg.SlippagePct, side, true)
	if side == Long {
		row.Stop = env.low - grid.StopBufferWidth*env.width
		row.Target = env.low + grid.TargetLevel*env.width
	} else {
		row.Stop = env.high + grid.StopBufferWidth*env.width
		row.Target = env.high - grid.TargetLevel*env.width
	}
	if row.Stop <= 0 || row.Target <= 0 || row.ExpectedEntryPrice <= 0 {
		row.SkippedReason = "non_positive_price"
		return row
	}
	row.EntryGeometryValid = validEntryGeometry(Signal{Side: row.Side, Stop: row.Stop, Target: row.Target}, row.ExpectedEntryPrice)
	if !row.EntryGeometryValid {
		row.SkippedReason = "entry_stop_target_invalid"
	}
	return row
}

func (s rangeFirstOccupancyStrategy) Name() string {
	return s.config.ConfigID
}

func (s rangeFirstOccupancyStrategy) OnCandle(ctx StrategyContext) (Signal, bool) {
	sig, ok := s.byIndex[ctx.Index]
	return sig, ok
}

func (s rangeFirstOccupancyStrategy) SignalRows() []FuturesRangeFirstOccupancyRotationV1SignalRow {
	return append([]FuturesRangeFirstOccupancyRotationV1SignalRow(nil), s.signals...)
}

func (s rangeFirstOccupancyStrategy) TradeRows(trades []Trade, splits []Split) []FuturesRangeFirstOccupancyRotationV1TradeRow {
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	rows := make([]FuturesRangeFirstOccupancyRotationV1TradeRow, 0, len(trades))
	for _, trade := range trades {
		signal, ok := s.signalsByID[trade.Signal]
		if !ok {
			continue
		}
		entryTime, _ := parseTime(trade.EntryTime)
		exitTime, _ := parseTime(trade.ExitTime)
		initialRisk := math.Abs(trade.EntryPrice-trade.Stop) * trade.Size
		row := FuturesRangeFirstOccupancyRotationV1TradeRow{
			ConfigID:        signal.ConfigID,
			Family:          signal.Family,
			IsBaseline:      signal.IsBaseline,
			Selected:        signal.Selected,
			SignalID:        signal.SignalID,
			Timeframe:       signal.Timeframe,
			SignalIndex:     signal.SignalIndex,
			SignalCloseTime: signal.SignalCloseTime,
			LookbackHours:   signal.LookbackHours,
			MaxWidthPct:     signal.MaxWidthPct,
			OccupancyWindow: signal.OccupancyWindow,
			RecaptureLevel:  signal.RecaptureLevel,
			TargetLevel:     signal.TargetLevel,
			StopBufferWidth: signal.StopBufferWidth,
			RangeHigh:       signal.RangeHigh,
			RangeLow:        signal.RangeLow,
			RangeWidth:      signal.RangeWidth,
			RangeMid:        signal.RangeMid,
			EntrySplit:      splitNameForCloseTime(entryTime, splits),
			CloseSplit:      splitNameForCloseTime(exitTime, splits),
			Side:            trade.Side,
			EntryTime:       trade.EntryTime,
			ExitTime:        trade.ExitTime,
			OpenIndex:       trade.OpenIndex,
			CloseIndex:      trade.CloseIndex,
			EntryPrice:      trade.EntryPrice,
			ExitPrice:       trade.ExitPrice,
			Stop:            trade.Stop,
			Target:          trade.Target,
			Size:            trade.Size,
			InitialRisk:     initialRisk,
			GrossPnL:        trade.GrossPnL,
			NetPnL:          trade.NetPnL,
			Fees:            trade.Fees,
			Slippage:        trade.Slippage,
			ExitReason:      trade.Reason,
			HoldBars:        trade.HoldBars,
		}
		if initialRisk > 0 {
			row.GrossR = trade.GrossPnL / initialRisk
			row.NetR = trade.NetPnL / initialRisk
		}
		rows = append(rows, row)
	}
	return rows
}

func rangeFirstOccupancyRotationMarkExecutedSignals(signals []FuturesRangeFirstOccupancyRotationV1SignalRow, trades []Trade, skips *[]FuturesRangeFirstOccupancyRotationV1SkipRow, grid FuturesRangeFirstOccupancyRotationV1GridConfig, frame rangeFirstOccupancyFrameCache) {
	executed := map[string]bool{}
	for _, trade := range trades {
		executed[trade.Signal] = true
	}
	skipCounts := map[rangeFirstOccupancySkipKey]int{}
	for i := range signals {
		if executed[signals[i].SignalID] {
			signals[i].Executed = true
			continue
		}
		signals[i].SkippedReason = "already_in_position"
		split := signals[i].Split
		for _, splitName := range rangeFirstOccupancyRotationSplitNames(split) {
			skipCounts[rangeFirstOccupancySkipKey{configID: grid.ConfigID, timeframe: grid.Timeframe, split: splitName, side: string(signals[i].Side), reason: "already_in_position"}]++
		}
	}
	*skips = append(*skips, rangeFirstOccupancyRotationSkipRowsFromCounts(skipCounts, grid, false, false)...)
	_ = frame
}

func rangeFirstOccupancyRotationSummarizeConfig(grid FuturesRangeFirstOccupancyRotationV1GridConfig, signals []FuturesRangeFirstOccupancyRotationV1SignalRow, skips []FuturesRangeFirstOccupancyRotationV1SkipRow, trades []Trade, startBalance float64, splits []Split) []FuturesRangeFirstOccupancyRotationV1SummaryRow {
	rows := make([]FuturesRangeFirstOccupancyRotationV1SummaryRow, 0, len(splits)*3)
	signalCounts := rangeFirstOccupancyRotationSignalCounts(signals)
	skipCounts := rangeFirstOccupancyRotationSkipCounts(skips)
	for _, split := range splits {
		for _, side := range []string{"all", string(Long), string(Short)} {
			filtered := []Trade{}
			for _, trade := range trades {
				exitTime, err := parseTime(trade.ExitTime)
				if err != nil || !split.Contains(exitTime) {
					continue
				}
				if side != "all" && string(trade.Side) != side {
					continue
				}
				filtered = append(filtered, trade)
			}
			summary := rangeFirstOccupancyRotationSummarizeTrades(filtered, startBalance)
			row := FuturesRangeFirstOccupancyRotationV1SummaryRow{
				ConfigID:             grid.ConfigID,
				Family:               rangeFirstOccupancyRotationFamily,
				IsBaseline:           grid.ConfigID == RangeFirstOccupancyRotationV1BaselineConfigID,
				Timeframe:            grid.Timeframe,
				LookbackHours:        grid.LookbackHours,
				MaxWidthPct:          grid.MaxWidthPct,
				OccupancyWindow:      grid.OccupancyWindow,
				OccupancyZoneLevel:   grid.OccupancyZoneLevel,
				OccupancyMinFraction: grid.OccupancyMinFraction,
				RecaptureLevel:       grid.RecaptureLevel,
				TargetLevel:          grid.TargetLevel,
				MaxHoldBars:          grid.MaxHoldBars,
				StopBufferWidth:      grid.StopBufferWidth,
				SideMode:             grid.SideMode,
				Split:                split.Name,
				Side:                 side,
				SignalCount:          signalCounts[split.Name+"|"+side] + skipCounts[split.Name+"|"+side],
				SkippedSignalCount:   skipCounts[split.Name+"|"+side],
				TotalTrades:          summary.TotalTrades,
				Wins:                 summary.Wins,
				Losses:               summary.Losses,
				WinRate:              summary.WinRate,
				GrossPnL:             summary.GrossPnL,
				NetPnL:               summary.NetPnL,
				TotalCosts:           summary.TotalCosts,
				ProfitFactor:         summary.ProfitFactor,
				GrossProfitFactor:    summary.GrossProfitFactor,
				MaxDrawdown:          summary.MaxDrawdown,
				AvgGrossR:            summary.AvgGrossR,
				AvgNetR:              summary.AvgNetR,
				AvgInitialRisk:       summary.AvgInitialRisk,
				AvgHoldBars:          summary.AvgHoldBars,
			}
			rows = append(rows, row)
		}
	}
	return rows
}

type rangeFirstOccupancyTradeSummary struct {
	TotalTrades       int
	Wins              int
	Losses            int
	WinRate           float64
	GrossPnL          float64
	NetPnL            float64
	TotalCosts        float64
	ProfitFactor      float64
	GrossProfitFactor float64
	MaxDrawdown       float64
	AvgGrossR         float64
	AvgNetR           float64
	AvgInitialRisk    float64
	AvgHoldBars       float64
}

func rangeFirstOccupancyRotationSummarizeTrades(trades []Trade, startBalance float64) rangeFirstOccupancyTradeSummary {
	row := rangeFirstOccupancyTradeSummary{TotalTrades: len(trades)}
	balance := startBalance
	equity := []float64{startBalance}
	netProfit, netLoss := 0.0, 0.0
	grossProfit, grossLoss := 0.0, 0.0
	holdBars := 0
	for _, trade := range trades {
		initialRisk := math.Abs(trade.EntryPrice-trade.Stop) * trade.Size
		row.GrossPnL += trade.GrossPnL
		row.NetPnL += trade.NetPnL
		row.TotalCosts += trade.Fees + trade.Slippage
		row.AvgInitialRisk += initialRisk
		if initialRisk > 0 {
			row.AvgGrossR += trade.GrossPnL / initialRisk
			row.AvgNetR += trade.NetPnL / initialRisk
		}
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

func rangeFirstOccupancyRotationGridRows(summary []FuturesRangeFirstOccupancyRotationV1SummaryRow, skips []FuturesRangeFirstOccupancyRotationV1SkipRow, gridConfigs []FuturesRangeFirstOccupancyRotationV1GridConfig, cfg FuturesRangeFirstOccupancyRotationV1OptimizationConfig) []FuturesRangeFirstOccupancyRotationV1GridRow {
	byKey := rangeFirstOccupancyRotationSummaryByKey(summary)
	rows := make([]FuturesRangeFirstOccupancyRotationV1GridRow, 0, len(gridConfigs))
	for _, grid := range gridConfigs {
		train := byKey[rangeFirstOccupancyRotationSummaryKey(grid.ConfigID, "2021_2022_stress", "all")]
		oos := byKey[rangeFirstOccupancyRotationSummaryKey(grid.ConfigID, "2023_2024_oos", "all")]
		recent := byKey[rangeFirstOccupancyRotationSummaryKey(grid.ConfigID, "2025_2026_recent", "all")]
		full := byKey[rangeFirstOccupancyRotationSummaryKey(grid.ConfigID, fullSplitName, "all")]
		fullLong := byKey[rangeFirstOccupancyRotationSummaryKey(grid.ConfigID, fullSplitName, string(Long))]
		fullShort := byKey[rangeFirstOccupancyRotationSummaryKey(grid.ConfigID, fullSplitName, string(Short))]
		oosLong := byKey[rangeFirstOccupancyRotationSummaryKey(grid.ConfigID, "2023_2024_oos", string(Long))]
		oosShort := byKey[rangeFirstOccupancyRotationSummaryKey(grid.ConfigID, "2023_2024_oos", string(Short))]
		recentLong := byKey[rangeFirstOccupancyRotationSummaryKey(grid.ConfigID, "2025_2026_recent", string(Long))]
		recentShort := byKey[rangeFirstOccupancyRotationSummaryKey(grid.ConfigID, "2025_2026_recent", string(Short))]
		row := FuturesRangeFirstOccupancyRotationV1GridRow{
			ConfigID:                grid.ConfigID,
			Family:                  rangeFirstOccupancyRotationFamily,
			IsBaseline:              grid.ConfigID == RangeFirstOccupancyRotationV1BaselineConfigID,
			Timeframe:               grid.Timeframe,
			LookbackHours:           grid.LookbackHours,
			MaxWidthPct:             grid.MaxWidthPct,
			OccupancyWindow:         grid.OccupancyWindow,
			OccupancyZoneLevel:      grid.OccupancyZoneLevel,
			OccupancyMinFraction:    grid.OccupancyMinFraction,
			RecaptureLevel:          grid.RecaptureLevel,
			TargetLevel:             grid.TargetLevel,
			MaxHoldBars:             grid.MaxHoldBars,
			StopBufferWidth:         grid.StopBufferWidth,
			SideMode:                grid.SideMode,
			SignalCount:             full.SignalCount - full.SkippedSignalCount,
			SkippedSignalCount:      full.SkippedSignalCount,
			TrainTrades:             train.TotalTrades,
			OOSTrades:               oos.TotalTrades,
			RecentTrades:            recent.TotalTrades,
			FullTrades:              full.TotalTrades,
			TrainNetPnL:             train.NetPnL,
			OOSNetPnL:               oos.NetPnL,
			RecentNetPnL:            recent.NetPnL,
			FullNetPnL:              full.NetPnL,
			TrainProfitFactor:       train.ProfitFactor,
			OOSProfitFactor:         oos.ProfitFactor,
			RecentProfitFactor:      recent.ProfitFactor,
			FullProfitFactor:        full.ProfitFactor,
			TrainMaxDrawdown:        train.MaxDrawdown,
			FullMaxDrawdown:         full.MaxDrawdown,
			TrainNetToDrawdown:      rangeFirstOccupancyRotationNetToDrawdown(train.NetPnL, train.MaxDrawdown),
			FullNetToDrawdown:       rangeFirstOccupancyRotationNetToDrawdown(full.NetPnL, full.MaxDrawdown),
			FullLongTrades:          fullLong.TotalTrades,
			FullShortTrades:         fullShort.TotalTrades,
			FullLongNetPnL:          fullLong.NetPnL,
			FullShortNetPnL:         fullShort.NetPnL,
			OOSLongNetPnL:           oosLong.NetPnL,
			OOSShortNetPnL:          oosShort.NetPnL,
			RecentLongNetPnL:        recentLong.NetPnL,
			RecentShortNetPnL:       recentShort.NetPnL,
			SideConcentration:       rangeFirstOccupancyRotationSideConcentration(fullLong.NetPnL, fullShort.NetPnL, full.NetPnL),
			SideWeaknessCaveat:      rangeFirstOccupancyRotationSideWeakness(fullLong, fullShort, oosLong, oosShort, recentLong, recentShort, cfg),
			SideConcentrationCaveat: rangeFirstOccupancyRotationSideConcentrationCaveat(fullLong.NetPnL, fullShort.NetPnL, full.NetPnL, cfg),
			ThinSideCaveat:          fullLong.TotalTrades < cfg.ThinSideTradeThreshold || fullShort.TotalTrades < cfg.ThinSideTradeThreshold,
		}
		row.CaveatPenalty = rangeFirstOccupancyRotationCaveatPenalty(row, cfg)
		row.PassesGate, row.FailureReason = rangeFirstOccupancyRotationEvaluateGrid(row, cfg)
		row.RankScore = rangeFirstOccupancyRotationRankScore(row)
		rows = append(rows, row)
	}
	return rows
}

func rangeFirstOccupancyRotationEvaluateGrid(row FuturesRangeFirstOccupancyRotationV1GridRow, cfg FuturesRangeFirstOccupancyRotationV1OptimizationConfig) (bool, string) {
	reasons := []string{}
	if row.TrainTrades < cfg.MinTrainTrades {
		reasons = append(reasons, "train_trades_below_100")
	}
	if row.OOSTrades < cfg.MinOOSTrades {
		reasons = append(reasons, "oos_trades_below_50")
	}
	if row.RecentTrades < cfg.MinRecentTrades {
		reasons = append(reasons, "recent_trades_below_25")
	}
	if row.FullTrades < cfg.MinFullTrades {
		reasons = append(reasons, "full_trades_below_200")
	}
	if row.TrainNetPnL <= 0 {
		reasons = append(reasons, "train_net_not_positive")
	}
	if row.OOSNetPnL <= 0 {
		reasons = append(reasons, "oos_net_not_positive")
	}
	if row.RecentNetPnL <= 0 {
		reasons = append(reasons, "recent_net_not_positive")
	}
	if row.FullNetPnL <= 0 {
		reasons = append(reasons, "full_net_not_positive")
	}
	if row.TrainProfitFactor < cfg.MinTrainProfitFactor {
		reasons = append(reasons, "train_pf_below_1_20")
	}
	if row.OOSProfitFactor < cfg.MinOOSProfitFactor {
		reasons = append(reasons, "oos_pf_below_1_05")
	}
	if row.RecentProfitFactor < cfg.MinRecentProfitFactor {
		reasons = append(reasons, "recent_pf_below_1_05")
	}
	if row.FullProfitFactor < cfg.MinFullProfitFactor {
		reasons = append(reasons, "full_pf_below_1_15")
	}
	if row.TrainMaxDrawdown > 0 && row.TrainNetToDrawdown < cfg.MinNetToDrawdown {
		reasons = append(reasons, "train_net_to_drawdown_below_1")
	}
	if row.FullMaxDrawdown > 0 && row.FullNetToDrawdown < cfg.MinNetToDrawdown {
		reasons = append(reasons, "full_net_to_drawdown_below_1")
	}
	if row.SideWeaknessCaveat {
		reasons = append(reasons, "side_loses_in_oos_and_recent")
	}
	if row.SideConcentrationCaveat {
		reasons = append(reasons, "side_net_concentration_above_75pct")
	}
	return len(reasons) == 0, strings.Join(uniqueStrings(reasons), ";")
}

func rangeFirstOccupancyRotationRankingRows(gridRows []FuturesRangeFirstOccupancyRotationV1GridRow) []FuturesRangeFirstOccupancyRotationV1RankingRow {
	rows := make([]FuturesRangeFirstOccupancyRotationV1RankingRow, 0, len(gridRows))
	for _, grid := range gridRows {
		rows = append(rows, FuturesRangeFirstOccupancyRotationV1RankingRow{
			ConfigID:             grid.ConfigID,
			Family:               grid.Family,
			IsBaseline:           grid.IsBaseline,
			Selected:             grid.Selected,
			Timeframe:            grid.Timeframe,
			LookbackHours:        grid.LookbackHours,
			MaxWidthPct:          grid.MaxWidthPct,
			OccupancyWindow:      grid.OccupancyWindow,
			OccupancyMinFraction: grid.OccupancyMinFraction,
			RecaptureLevel:       grid.RecaptureLevel,
			TargetLevel:          grid.TargetLevel,
			MaxHoldBars:          grid.MaxHoldBars,
			StopBufferWidth:      grid.StopBufferWidth,
			TrainTrades:          grid.TrainTrades,
			TrainNetPnL:          grid.TrainNetPnL,
			TrainProfitFactor:    grid.TrainProfitFactor,
			TrainMaxDrawdown:     grid.TrainMaxDrawdown,
			OOSTrades:            grid.OOSTrades,
			OOSNetPnL:            grid.OOSNetPnL,
			OOSProfitFactor:      grid.OOSProfitFactor,
			RecentTrades:         grid.RecentTrades,
			RecentNetPnL:         grid.RecentNetPnL,
			RecentProfitFactor:   grid.RecentProfitFactor,
			FullTrades:           grid.FullTrades,
			FullNetPnL:           grid.FullNetPnL,
			FullProfitFactor:     grid.FullProfitFactor,
			FullMaxDrawdown:      grid.FullMaxDrawdown,
			CaveatPenalty:        grid.CaveatPenalty,
			RankScore:            grid.RankScore,
			PassesGate:           grid.PassesGate,
			FailureReason:        grid.FailureReason,
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].PassesGate != rows[j].PassesGate {
			return rows[i].PassesGate
		}
		if rows[i].RankScore != rows[j].RankScore {
			return rows[i].RankScore > rows[j].RankScore
		}
		if rows[i].TrainNetPnL != rows[j].TrainNetPnL {
			return rows[i].TrainNetPnL > rows[j].TrainNetPnL
		}
		if rows[i].TrainProfitFactor != rows[j].TrainProfitFactor {
			return rows[i].TrainProfitFactor > rows[j].TrainProfitFactor
		}
		if rows[i].TrainMaxDrawdown != rows[j].TrainMaxDrawdown {
			return rows[i].TrainMaxDrawdown < rows[j].TrainMaxDrawdown
		}
		if rows[i].FullTrades != rows[j].FullTrades {
			return rows[i].FullTrades > rows[j].FullTrades
		}
		if rangeFirstOccupancyRotationTimeframeSortKey(rows[i].Timeframe) != rangeFirstOccupancyRotationTimeframeSortKey(rows[j].Timeframe) {
			return rangeFirstOccupancyRotationTimeframeSortKey(rows[i].Timeframe) < rangeFirstOccupancyRotationTimeframeSortKey(rows[j].Timeframe)
		}
		if rows[i].MaxHoldBars != rows[j].MaxHoldBars {
			return rows[i].MaxHoldBars < rows[j].MaxHoldBars
		}
		return rows[i].ConfigID < rows[j].ConfigID
	})
	for i := range rows {
		rows[i].Rank = i + 1
	}
	return rows
}

func rangeFirstOccupancyRotationBaselineRows(gridRows []FuturesRangeFirstOccupancyRotationV1GridRow, rankingRows []FuturesRangeFirstOccupancyRotationV1RankingRow) []FuturesRangeFirstOccupancyRotationV1BaselineRow {
	ranks := map[string]int{}
	for _, row := range rankingRows {
		ranks[row.ConfigID] = row.Rank
	}
	for _, row := range gridRows {
		if row.ConfigID != RangeFirstOccupancyRotationV1BaselineConfigID {
			continue
		}
		return []FuturesRangeFirstOccupancyRotationV1BaselineRow{{
			ConfigID:           row.ConfigID,
			PassesGate:         row.PassesGate,
			Rank:               ranks[row.ConfigID],
			FullTrades:         row.FullTrades,
			FullNetPnL:         row.FullNetPnL,
			FullProfitFactor:   row.FullProfitFactor,
			FullMaxDrawdown:    row.FullMaxDrawdown,
			TrainTrades:        row.TrainTrades,
			TrainNetPnL:        row.TrainNetPnL,
			TrainProfitFactor:  row.TrainProfitFactor,
			OOSTrades:          row.OOSTrades,
			OOSNetPnL:          row.OOSNetPnL,
			OOSProfitFactor:    row.OOSProfitFactor,
			RecentTrades:       row.RecentTrades,
			RecentNetPnL:       row.RecentNetPnL,
			RecentProfitFactor: row.RecentProfitFactor,
			RankScore:          row.RankScore,
			FailureReason:      row.FailureReason,
		}}
	}
	return nil
}

func rangeFirstOccupancyRotationSelectionRows(gridRows []FuturesRangeFirstOccupancyRotationV1GridRow, rankingRows []FuturesRangeFirstOccupancyRotationV1RankingRow, baselineID string, selectedID string, stopState string) []FuturesRangeFirstOccupancyRotationV1SelectionRow {
	gridByID := map[string]FuturesRangeFirstOccupancyRotationV1GridRow{}
	rankByID := map[string]int{}
	for _, row := range gridRows {
		gridByID[row.ConfigID] = row
	}
	for _, row := range rankingRows {
		rankByID[row.ConfigID] = row.Rank
	}
	rows := []FuturesRangeFirstOccupancyRotationV1SelectionRow{}
	rows = append(rows, rangeFirstOccupancyRotationSelectionRow("baseline", gridByID[baselineID], rankByID[baselineID], stopState))
	if selectedID != "" {
		rows = append(rows, rangeFirstOccupancyRotationSelectionRow("selected", gridByID[selectedID], rankByID[selectedID], stopState))
	} else {
		rows = append(rows, FuturesRangeFirstOccupancyRotationV1SelectionRow{Role: "selected", StopState: stopState, FailureReason: "no_selectable_config"})
	}
	return rows
}

func rangeFirstOccupancyRotationSelectionRow(role string, row FuturesRangeFirstOccupancyRotationV1GridRow, rank int, stopState string) FuturesRangeFirstOccupancyRotationV1SelectionRow {
	return FuturesRangeFirstOccupancyRotationV1SelectionRow{
		Role:               role,
		ConfigID:           row.ConfigID,
		Rank:               rank,
		PassesGate:         row.PassesGate,
		FullTrades:         row.FullTrades,
		FullNetPnL:         row.FullNetPnL,
		FullProfitFactor:   row.FullProfitFactor,
		FullMaxDrawdown:    row.FullMaxDrawdown,
		TrainTrades:        row.TrainTrades,
		TrainNetPnL:        row.TrainNetPnL,
		TrainProfitFactor:  row.TrainProfitFactor,
		OOSTrades:          row.OOSTrades,
		OOSNetPnL:          row.OOSNetPnL,
		OOSProfitFactor:    row.OOSProfitFactor,
		RecentTrades:       row.RecentTrades,
		RecentNetPnL:       row.RecentNetPnL,
		RecentProfitFactor: row.RecentProfitFactor,
		RankScore:          row.RankScore,
		FailureReason:      row.FailureReason,
		StopState:          stopState,
	}
}

func rangeFirstOccupancyRotationMarkSelected(selected string, gridRows []FuturesRangeFirstOccupancyRotationV1GridRow, summaryRows []FuturesRangeFirstOccupancyRotationV1SummaryRow, rankingRows []FuturesRangeFirstOccupancyRotationV1RankingRow, skipRows []FuturesRangeFirstOccupancyRotationV1SkipRow) {
	if selected == "" {
		return
	}
	for i := range gridRows {
		gridRows[i].Selected = gridRows[i].ConfigID == selected
	}
	for i := range summaryRows {
		summaryRows[i].Selected = summaryRows[i].ConfigID == selected
	}
	for i := range rankingRows {
		rankingRows[i].Selected = rankingRows[i].ConfigID == selected
	}
	for i := range skipRows {
		skipRows[i].Selected = skipRows[i].ConfigID == selected
	}
}

func rangeFirstOccupancyRotationMarkWorstSplits(rows []FuturesRangeFirstOccupancyRotationV1SummaryRow) {
	worst := map[string]int{}
	for i := range rows {
		if rows[i].Split == fullSplitName || rows[i].Side != "all" {
			continue
		}
		key := rows[i].ConfigID
		if existing, ok := worst[key]; !ok || rows[i].NetPnL < rows[existing].NetPnL {
			worst[key] = i
		}
	}
	for _, index := range worst {
		rows[index].IsWorstPeriodSplit = true
	}
}

func rangeFirstOccupancyRotationSignalCounts(signals []FuturesRangeFirstOccupancyRotationV1SignalRow) map[string]int {
	out := map[string]int{}
	for _, signal := range signals {
		for _, split := range rangeFirstOccupancyRotationSplitNames(signal.Split) {
			out[split+"|all"]++
			out[split+"|"+string(signal.Side)]++
		}
	}
	return out
}

func rangeFirstOccupancyRotationSkipCounts(skips []FuturesRangeFirstOccupancyRotationV1SkipRow) map[string]int {
	out := map[string]int{}
	for _, row := range skips {
		out[row.Split+"|"+row.Side] += row.Count
		if row.Side != "all" {
			out[row.Split+"|all"] += row.Count
		}
	}
	return out
}

func rangeFirstOccupancyRotationSkipRowsFromCounts(counts map[rangeFirstOccupancySkipKey]int, grid FuturesRangeFirstOccupancyRotationV1GridConfig, isBaseline bool, selected bool) []FuturesRangeFirstOccupancyRotationV1SkipRow {
	rows := make([]FuturesRangeFirstOccupancyRotationV1SkipRow, 0, len(counts))
	for key, count := range counts {
		if count == 0 {
			continue
		}
		rows = append(rows, FuturesRangeFirstOccupancyRotationV1SkipRow{
			ConfigID:   key.configID,
			Family:     rangeFirstOccupancyRotationFamily,
			IsBaseline: isBaseline || key.configID == RangeFirstOccupancyRotationV1BaselineConfigID,
			Selected:   selected,
			Timeframe:  key.timeframe,
			Split:      key.split,
			Side:       key.side,
			Reason:     key.reason,
			Count:      count,
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].ConfigID != rows[j].ConfigID {
			return rows[i].ConfigID < rows[j].ConfigID
		}
		if rows[i].Split != rows[j].Split {
			return rows[i].Split < rows[j].Split
		}
		if rows[i].Side != rows[j].Side {
			return rows[i].Side < rows[j].Side
		}
		return rows[i].Reason < rows[j].Reason
	})
	return rows
}

func rangeFirstOccupancyRotationSummaryByKey(rows []FuturesRangeFirstOccupancyRotationV1SummaryRow) map[string]FuturesRangeFirstOccupancyRotationV1SummaryRow {
	out := map[string]FuturesRangeFirstOccupancyRotationV1SummaryRow{}
	for _, row := range rows {
		out[rangeFirstOccupancyRotationSummaryKey(row.ConfigID, row.Split, row.Side)] = row
	}
	return out
}

func rangeFirstOccupancyRotationSummaryKey(configID, split, side string) string {
	return configID + "|" + split + "|" + side
}

func rangeFirstOccupancyCountsKey(lookback int, window int) string {
	return fmt.Sprintf("%d|%d", lookback, window)
}

func rangeFirstOccupancyRotationSplitNames(split string) []string {
	if split == "" || split == fullSplitName {
		return []string{fullSplitName}
	}
	return []string{split, fullSplitName}
}

func rangeFirstOccupancyRotationLookbackBars(timeframe string, lookbackHours int) (int, error) {
	switch timeframe {
	case RangeDiscoveryTimeframe15m:
		return lookbackHours * 4, nil
	case RangeDiscoveryTimeframe1h:
		return lookbackHours, nil
	default:
		return 0, fmt.Errorf("range occupancy rotation unsupported timeframe %q", timeframe)
	}
}

func rangeFirstOccupancyRotationGridConfigByID(rows []FuturesRangeFirstOccupancyRotationV1GridConfig, id string) (FuturesRangeFirstOccupancyRotationV1GridConfig, bool) {
	for _, row := range rows {
		if row.ConfigID == id {
			return row, true
		}
	}
	return FuturesRangeFirstOccupancyRotationV1GridConfig{}, false
}

func rangeFirstOccupancyRotationRankScore(row FuturesRangeFirstOccupancyRotationV1GridRow) float64 {
	pfComponent := math.Min(row.TrainProfitFactor, 3.0) - 1.0
	tradeComponent := math.Min(float64(row.TrainTrades), 400.0) / 400.0
	drawdownComponent := row.TrainNetPnL / math.Max(row.TrainMaxDrawdown, 1.0)
	return drawdownComponent + pfComponent + tradeComponent - row.CaveatPenalty
}

func rangeFirstOccupancyRotationCaveatPenalty(row FuturesRangeFirstOccupancyRotationV1GridRow, cfg FuturesRangeFirstOccupancyRotationV1OptimizationConfig) float64 {
	penalty := 0.0
	if row.SideConcentration > 0.60 {
		penalty += cfg.SideConcentrationPenalty
	}
	if row.ThinSideCaveat {
		penalty += cfg.ThinSidePenalty
	}
	return penalty
}

func rangeFirstOccupancyRotationNetToDrawdown(net float64, drawdown float64) float64 {
	if drawdown <= 0 {
		if net > 0 {
			return 999.99
		}
		return 0
	}
	return net / drawdown
}

func rangeFirstOccupancyRotationSideConcentration(longNet float64, shortNet float64, totalNet float64) float64 {
	if totalNet <= 0 {
		return 0
	}
	longContribution := math.Max(longNet, 0)
	shortContribution := math.Max(shortNet, 0)
	totalContribution := longContribution + shortContribution
	if totalContribution <= 0 {
		return 0
	}
	return math.Max(longContribution, shortContribution) / totalContribution
}

func rangeFirstOccupancyRotationSideWeakness(fullLong, fullShort, oosLong, oosShort, recentLong, recentShort FuturesRangeFirstOccupancyRotationV1SummaryRow, cfg FuturesRangeFirstOccupancyRotationV1OptimizationConfig) bool {
	if fullLong.TotalTrades >= cfg.MinSideTradesForWeakness && oosLong.NetPnL < 0 && recentLong.NetPnL < 0 {
		return true
	}
	if fullShort.TotalTrades >= cfg.MinSideTradesForWeakness && oosShort.NetPnL < 0 && recentShort.NetPnL < 0 {
		return true
	}
	return false
}

func rangeFirstOccupancyRotationSideConcentrationCaveat(longNet float64, shortNet float64, totalNet float64, cfg FuturesRangeFirstOccupancyRotationV1OptimizationConfig) bool {
	concentration := rangeFirstOccupancyRotationSideConcentration(longNet, shortNet, totalNet)
	if concentration <= cfg.MaxSideNetConcentration {
		return false
	}
	if longNet < shortNet {
		return longNet <= 0
	}
	return shortNet <= 0
}

func rangeFirstOccupancyRotationTimeframeSortKey(timeframe string) int {
	switch timeframe {
	case RangeDiscoveryTimeframe1h:
		return 1
	case RangeDiscoveryTimeframe15m:
		return 2
	default:
		return 99
	}
}
