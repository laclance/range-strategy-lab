package lab

import (
	"fmt"
	"math"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	RangeUniverseSymbolBTCUSDT = "BTCUSDT"
	RangeUniverseSymbolETHUSDT = "ETHUSDT"
	RangeUniverseSymbolSOLUSDT = "SOLUSDT"

	RangeUniverseFamilyBreakoutRetestAcceptance   = "breakout_retest_acceptance"
	RangeUniverseFamilyStructuredCompressionBreak = "structured_compression_expansion"

	RangeUniverseOutcomeRetestContinuation       = "retest_continuation"
	RangeUniverseOutcomeRetestReentry            = "retest_reentry"
	RangeUniverseOutcomeRetestStall              = "retest_stall"
	RangeUniverseOutcomeStructuredContinuation   = "structured_continuation"
	RangeUniverseOutcomeStructuredFailedReentry  = "structured_failed_reentry"
	RangeUniverseOutcomeStructuredNoContinuation = "structured_no_continuation"

	RangeUniverseStopStateSourceGap            = "range_universe_source_gap"
	RangeUniverseStopStateNoBacktestCandidate  = "range_universe_no_backtest_candidate"
	RangeUniverseStopStateAuditReady           = "range_universe_audit_ready_for_baseline_backtest"
	RangeUniverseStopStateCodegenOrTestBlocked = "range_universe_codegen_or_test_blocked"
	RangeUniverseStopStateReviewOnlyNoStrategy = "range_universe_review_only_no_strategy_change"
)

type FuturesRangeUniverseSourceConfig struct {
	Symbol                    string
	Path                      string
	ApprovedPath              string
	SkipSplitEligibilityCheck bool
}

type FuturesRangeUniverseDiscoveryAuditConfig struct {
	Sources                             []FuturesRangeUniverseSourceConfig
	Discovery                           FuturesRangeCandidateDiscoveryAuditConfig
	SymbolSpecificMinCandidatesPerSplit int
	SymbolSpecificCostBufferMultiple    float64
}

type FuturesRangeUniverseDiscoveryResult struct {
	SourceRows    []FuturesRangeUniverseSourceRow
	CoverageRows  []FuturesRangeUniverseCoverageRow
	CandidateRows []FuturesRangeUniverseCandidateRow
	SummaryRows   []FuturesRangeUniverseSummaryRow
	RankingRows   []FuturesRangeUniverseRankingRow
	StabilityRows []FuturesRangeUniverseStabilityRow
}

type FuturesRangeUniverseSourceRow struct {
	Symbol                    string `json:"symbol"`
	Path                      string `json:"path"`
	ApprovedPath              string `json:"approved_path"`
	Venue                     string `json:"venue"`
	Product                   string `json:"product"`
	Interval                  string `json:"interval"`
	RowCount                  int    `json:"row_count"`
	FirstOpenTime             string `json:"first_open_time"`
	LastOpenTime              string `json:"last_open_time"`
	Schema                    string `json:"schema"`
	TimestampSemantics        string `json:"timestamp_semantics"`
	FinalityRule              string `json:"finality_rule"`
	PhysicalNonMonotonicCount int    `json:"physical_non_monotonic_count"`
	SortedForValidation       bool   `json:"sorted_for_validation"`
	AcceptedMonotonic         bool   `json:"accepted_monotonic"`
	GapCount                  int    `json:"gap_count"`
	DuplicateCount            int    `json:"duplicate_count"`
	ZeroVolumeCount           int    `json:"zero_volume_count"`
	StressSplitRows           int    `json:"stress_split_rows"`
	OOSSplitRows              int    `json:"oos_split_rows"`
	RecentSplitRows           int    `json:"recent_split_rows"`
	AllPeriodSplitsEligible   bool   `json:"all_period_splits_eligible"`
	ComparisonOnly            bool   `json:"comparison_only"`
	ValidationStatus          string `json:"validation_status"`
	ValidationError           string `json:"validation_error,omitempty"`
}

type FuturesRangeUniverseCoverageRow struct {
	Symbol string `json:"symbol"`
	FuturesRangeDiscoveryCoverageRow
}

type FuturesRangeUniverseCandidateRow struct {
	Symbol string `json:"symbol"`
	FuturesRangeDiscoveryCandidateRow
}

type FuturesRangeUniverseSummaryRow struct {
	Symbol string `json:"symbol"`
	FuturesRangeDiscoverySummaryRow
}

type FuturesRangeUniverseStabilityRow struct {
	Symbol string `json:"symbol"`
	FuturesRangeDiscoveryStabilityRow
}

type FuturesRangeUniverseRankingRow struct {
	Rank                                 int     `json:"rank"`
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
	SymbolSpecificWeakestSplitCount      int     `json:"symbol_specific_weakest_split_count"`
	SymbolSpecificWeakestCostBufferPct   float64 `json:"symbol_specific_weakest_cost_buffer_pct"`
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
	FailureReason                        string  `json:"failure_reason,omitempty"`
}

type rangeUniverseRankingKey struct {
	timeframe   string
	family      string
	side        string
	horizonBars int
}

type rangeUniverseSymbolGate struct {
	symbol          string
	full            FuturesRangeUniverseSummaryRow
	stability       FuturesRangeUniverseStabilityRow
	passes          bool
	symbolException bool
	reasons         []string
	score           float64
}

func DefaultFuturesRangeUniverseDiscoveryAuditConfig() FuturesRangeUniverseDiscoveryAuditConfig {
	discovery := DefaultFuturesRangeCandidateDiscoveryAuditConfig()
	return FuturesRangeUniverseDiscoveryAuditConfig{
		Sources: []FuturesRangeUniverseSourceConfig{
			{Symbol: RangeUniverseSymbolBTCUSDT, Path: "../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv", ApprovedPath: "../binance-bot/data/btcusdt_futures_um_5m_2021_2026.csv"},
			{Symbol: RangeUniverseSymbolETHUSDT, Path: "../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv", ApprovedPath: "../binance-bot/data/ethusdt_futures_um_5m_2021_2026.csv"},
			{Symbol: RangeUniverseSymbolSOLUSDT, Path: "../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv", ApprovedPath: "../binance-bot/data/solusdt_futures_um_5m_2021_2026.csv"},
		},
		Discovery:                           discovery,
		SymbolSpecificMinCandidatesPerSplit: 200,
		SymbolSpecificCostBufferMultiple:    2,
	}
}

func RunFuturesRangeUniverseDiscoveryAudit(cfg FuturesRangeUniverseDiscoveryAuditConfig, splits []Split) (FuturesRangeUniverseDiscoveryResult, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return FuturesRangeUniverseDiscoveryResult{}, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}

	result := FuturesRangeUniverseDiscoveryResult{}
	eventOffset := 0
	for _, source := range cfg.Sources {
		candles, sourceRow, err := LoadFuturesRangeUniverseSource(source, splits)
		result.SourceRows = append(result.SourceRows, sourceRow)
		if err != nil {
			return result, err
		}

		coverage, candidates, err := runRangeUniverseSymbolCandidates(source.Symbol, candles, cfg.Discovery, splits, &eventOffset)
		result.CoverageRows = append(result.CoverageRows, coverage...)
		if err != nil {
			return result, err
		}
		result.CandidateRows = append(result.CandidateRows, candidates...)
	}

	sort.Slice(result.CandidateRows, func(i, j int) bool {
		if result.CandidateRows[i].Symbol != result.CandidateRows[j].Symbol {
			return rangeUniverseSymbolSortKey(result.CandidateRows[i].Symbol) < rangeUniverseSymbolSortKey(result.CandidateRows[j].Symbol)
		}
		return lessRangeDiscoveryCandidate(result.CandidateRows[i].FuturesRangeDiscoveryCandidateRow, result.CandidateRows[j].FuturesRangeDiscoveryCandidateRow)
	})
	result.SummaryRows = summarizeRangeUniverseDiscovery(result.CandidateRows, cfg.Discovery.RoundTripCostPct)
	result.StabilityRows = rangeUniverseStabilityRows(result.SummaryRows, splits)
	result.RankingRows = rangeUniverseRankingRows(result.SummaryRows, result.StabilityRows, cfg, splits)
	return result, nil
}

func LoadFuturesRangeUniverseSource(source FuturesRangeUniverseSourceConfig, splits []Split) ([]Candle, FuturesRangeUniverseSourceRow, error) {
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	row := FuturesRangeUniverseSourceRow{
		Symbol:             strings.ToUpper(strings.TrimSpace(source.Symbol)),
		Path:               source.Path,
		ApprovedPath:       source.ApprovedPath,
		Venue:              "Binance",
		Product:            "Binance USDT-M futures",
		Interval:           "5m",
		TimestampSemantics: "open_time",
		FinalityRule:       "closed 5m candles only; close_time must equal open_time plus 5m minus 1ms",
		AcceptedMonotonic:  true,
		ValidationStatus:   "accepted",
	}
	reject := func(format string, args ...any) ([]Candle, FuturesRangeUniverseSourceRow, error) {
		row.ValidationStatus = "rejected"
		row.ValidationError = fmt.Sprintf(format, args...)
		return nil, row, fmt.Errorf("%s", row.ValidationError)
	}
	if row.Symbol == "" {
		return reject("range universe source symbol is required")
	}
	if source.Path == "" {
		return reject("range universe source path is required for %s", row.Symbol)
	}
	approved := source.ApprovedPath
	if approved == "" {
		approved = source.Path
		row.ApprovedPath = approved
	}
	if !sameCleanPath(source.Path, approved) {
		return reject("%s source path %q is not the approved local path %q", row.Symbol, source.Path, approved)
	}
	if err := validateRangeUniverseSourcePath(row.Symbol, source.Path); err != nil {
		return reject(err.Error())
	}

	candles, header, err := LoadCSVWithHeader(source.Path)
	if err != nil {
		return reject("load %s source: %v", row.Symbol, err)
	}
	row.Schema = strings.Join(normalizeSchema(header), ",")
	if len(candles) == 0 {
		return reject("%s source CSV contains no candles", row.Symbol)
	}
	row.PhysicalNonMonotonicCount = physicalNonMonotonicCount(candles)
	if row.PhysicalNonMonotonicCount > 0 {
		row.SortedForValidation = true
		sort.SliceStable(candles, func(i, j int) bool {
			return candles[i].OpenTime.Before(candles[j].OpenTime)
		})
	}
	if err := validateRangeUniverseSortedCandles(row.Symbol, candles, &row, splits); err != nil {
		return reject(err.Error())
	}
	row.RowCount = len(candles)
	row.FirstOpenTime = candles[0].OpenTime.UTC().Format(timeLayout)
	row.LastOpenTime = candles[len(candles)-1].OpenTime.UTC().Format(timeLayout)
	row.AllPeriodSplitsEligible = row.StressSplitRows > 0 && row.OOSSplitRows > 0 && row.RecentSplitRows > 0
	if !source.SkipSplitEligibilityCheck && !row.AllPeriodSplitsEligible {
		return reject("%s source does not cover every required period split", row.Symbol)
	}
	return candles, row, nil
}

func runRangeUniverseSymbolCandidates(symbol string, candles []Candle, cfg FuturesRangeCandidateDiscoveryAuditConfig, splits []Split, eventOffset *int) ([]FuturesRangeUniverseCoverageRow, []FuturesRangeUniverseCandidateRow, error) {
	coverageRows := []FuturesRangeUniverseCoverageRow{}
	candidateRows := []FuturesRangeUniverseCandidateRow{}
	for _, frame := range rangeUniverseFrameDefs() {
		frameCandles, coverage, err := resampleRangeDiscoveryFrame(candles, frame)
		coverageRows = append(coverageRows, FuturesRangeUniverseCoverageRow{Symbol: symbol, FuturesRangeDiscoveryCoverageRow: coverage})
		if err != nil {
			return coverageRows, nil, err
		}
		if len(frameCandles) == 0 {
			continue
		}
		events, err := rangeUniverseEventsForFrame(frameCandles, frame, cfg, splits, eventOffset)
		if err != nil {
			return coverageRows, nil, err
		}
		for _, event := range events {
			for _, horizon := range cfg.HorizonsBars {
				candidateRows = append(candidateRows, FuturesRangeUniverseCandidateRow{
					Symbol:                            symbol,
					FuturesRangeDiscoveryCandidateRow: newRangeDiscoveryCandidateRow(frameCandles, event, horizon, cfg),
				})
			}
		}
	}
	return coverageRows, candidateRows, nil
}

func rangeUniverseFrameDefs() []rangeDiscoveryFrameDef {
	return []rangeDiscoveryFrameDef{
		{
			timeframe:  RangeDiscoveryTimeframe5m,
			interval:   5 * time.Minute,
			childBars:  1,
			barsPerDay: 24 * 12,
			families: []string{
				RangeDiscoveryFamilyBoundaryTouch,
				RangeDiscoveryFamilyWickRejection,
				RangeDiscoveryFamilyFailedBreakReentry,
			},
		},
		{
			timeframe:  RangeDiscoveryTimeframe15m,
			interval:   15 * time.Minute,
			childBars:  3,
			barsPerDay: 24 * 4,
			families: []string{
				RangeUniverseFamilyBreakoutRetestAcceptance,
				RangeDiscoveryFamilyBoundaryTouch,
				RangeDiscoveryFamilyWickRejection,
				RangeDiscoveryFamilyFailedBreakReentry,
				RangeUniverseFamilyStructuredCompressionBreak,
			},
		},
		{
			timeframe:  RangeDiscoveryTimeframe1h,
			interval:   time.Hour,
			childBars:  12,
			barsPerDay: 24,
			families: []string{
				RangeUniverseFamilyBreakoutRetestAcceptance,
				RangeDiscoveryFamilyBalancePersistence,
				RangeDiscoveryFamilyBoundaryTouch,
				RangeDiscoveryFamilyWickRejection,
				RangeDiscoveryFamilyFailedBreakReentry,
				RangeUniverseFamilyStructuredCompressionBreak,
			},
		},
		{
			timeframe:  RangeDiscoveryTimeframe4h,
			interval:   4 * time.Hour,
			childBars:  48,
			barsPerDay: 6,
			families: []string{
				RangeUniverseFamilyBreakoutRetestAcceptance,
				RangeDiscoveryFamilyBalancePersistence,
				RangeUniverseFamilyStructuredCompressionBreak,
			},
		},
	}
}

func rangeUniverseEventsForFrame(candles []Candle, frame rangeDiscoveryFrameDef, cfg FuturesRangeCandidateDiscoveryAuditConfig, splits []Split, nextEventID *int) ([]rangeDiscoveryEvent, error) {
	detectorCfg := DefaultCompressionRangeDetectorConfig()
	detectorCfg.BarsPerDay = frame.barsPerDay
	detectorCfg.LookbackDays = cfg.DetectorLookbackDays
	detectorCfg.Percentile = cfg.DetectorPercentile
	detectorCfg.MinConsecutiveBars = cfg.DetectorMinConsecutiveBars
	if cfg.DetectorLookbackBarsOverride > 0 {
		detectorCfg.LookbackDays = 1
		detectorCfg.BarsPerDay = cfg.DetectorLookbackBarsOverride
	}
	classifications, err := (CompressionRangeDetector{Config: detectorCfg}).Classify(candles)
	if err != nil {
		return nil, err
	}
	normalizedATR := NormalizedATR(candles, detectorCfg.ATRPeriod)
	episodes := rangeRegimeDurabilityEpisodes(candles, classifications, normalizedATR, splits, BalancedDetectorProfileID)
	events := []rangeDiscoveryEvent{}
	for _, episode := range episodes {
		for _, family := range frame.families {
			events = append(events, rangeUniverseEventsForFamily(candles, episode, frame, family, cfg, nextEventID)...)
		}
	}
	sort.Slice(events, func(i, j int) bool {
		if events[i].eventIndex != events[j].eventIndex {
			return events[i].eventIndex < events[j].eventIndex
		}
		if events[i].family != events[j].family {
			return rangeDiscoveryFamilySortKey(events[i].family) < rangeDiscoveryFamilySortKey(events[j].family)
		}
		return rangeDiscoverySideSortKey(events[i].side) < rangeDiscoverySideSortKey(events[j].side)
	})
	return events, nil
}

func rangeUniverseEventsForFamily(candles []Candle, episode rangeRegimeDurabilityEpisode, frame rangeDiscoveryFrameDef, family string, cfg FuturesRangeCandidateDiscoveryAuditConfig, nextEventID *int) []rangeDiscoveryEvent {
	switch family {
	case RangeUniverseFamilyBreakoutRetestAcceptance:
		return rangeUniverseBreakoutRetestEvents(candles, episode, frame, cfg, nextEventID)
	case RangeUniverseFamilyStructuredCompressionBreak:
		return rangeUniverseStructuredCompressionEvents(candles, episode, frame, cfg, nextEventID)
	default:
		return rangeDiscoveryEventsForFamily(candles, episode, frame, family, cfg, nextEventID)
	}
}

func rangeUniverseBreakoutRetestEvents(candles []Candle, episode rangeRegimeDurabilityEpisode, frame rangeDiscoveryFrameDef, cfg FuturesRangeCandidateDiscoveryAuditConfig, nextEventID *int) []rangeDiscoveryEvent {
	events := []rangeDiscoveryEvent{}
	for breakIndex := episode.EndIndex + 1; breakIndex < len(candles) && breakIndex <= episode.EndIndex+cfg.MaxEventDelayBars; breakIndex++ {
		breakCandle := candles[breakIndex]
		if breakCandle.Close > episode.High {
			for retestIndex := breakIndex + 1; retestIndex < len(candles) && retestIndex <= breakIndex+cfg.ReentryWindowBars; retestIndex++ {
				retest := candles[retestIndex]
				if retest.Low <= episode.High && retest.Close > episode.High {
					(*nextEventID)++
					events = append(events, newRangeDiscoveryDirectionalEvent(*nextEventID, frame, episode, retestIndex, RangeUniverseFamilyBreakoutRetestAcceptance, RangeDiscoverySideUp, "up", retestIndex-episode.EndIndex, breakIndex-episode.EndIndex, retestIndex-breakIndex, 0, RangeUniverseOutcomeRetestContinuation, RangeUniverseOutcomeRetestReentry, RangeUniverseOutcomeRetestStall))
					break
				}
			}
		}
		if breakCandle.Close < episode.Low {
			for retestIndex := breakIndex + 1; retestIndex < len(candles) && retestIndex <= breakIndex+cfg.ReentryWindowBars; retestIndex++ {
				retest := candles[retestIndex]
				if retest.High >= episode.Low && retest.Close < episode.Low {
					(*nextEventID)++
					events = append(events, newRangeDiscoveryDirectionalEvent(*nextEventID, frame, episode, retestIndex, RangeUniverseFamilyBreakoutRetestAcceptance, RangeDiscoverySideDown, "down", retestIndex-episode.EndIndex, breakIndex-episode.EndIndex, retestIndex-breakIndex, 0, RangeUniverseOutcomeRetestContinuation, RangeUniverseOutcomeRetestReentry, RangeUniverseOutcomeRetestStall))
					break
				}
			}
		}
	}
	return events
}

func rangeUniverseStructuredCompressionEvents(candles []Candle, episode rangeRegimeDurabilityEpisode, frame rangeDiscoveryFrameDef, cfg FuturesRangeCandidateDiscoveryAuditConfig, nextEventID *int) []rangeDiscoveryEvent {
	events := []rangeDiscoveryEvent{}
	for breakIndex := episode.EndIndex + 1; breakIndex < len(candles) && breakIndex <= episode.EndIndex+cfg.MaxEventDelayBars; breakIndex++ {
		breakCandle := candles[breakIndex]
		if breakCandle.Close > episode.High {
			for confirmIndex := breakIndex + 1; confirmIndex < len(candles) && confirmIndex <= breakIndex+cfg.ReentryWindowBars; confirmIndex++ {
				confirm := candles[confirmIndex]
				if confirm.Close > episode.High && confirm.High > breakCandle.High {
					(*nextEventID)++
					events = append(events, newRangeDiscoveryDirectionalEvent(*nextEventID, frame, episode, confirmIndex, RangeUniverseFamilyStructuredCompressionBreak, RangeDiscoverySideUp, "up", confirmIndex-episode.EndIndex, breakIndex-episode.EndIndex, confirmIndex-breakIndex, 0, RangeUniverseOutcomeStructuredContinuation, RangeUniverseOutcomeStructuredFailedReentry, RangeUniverseOutcomeStructuredNoContinuation))
					break
				}
			}
		}
		if breakCandle.Close < episode.Low {
			for confirmIndex := breakIndex + 1; confirmIndex < len(candles) && confirmIndex <= breakIndex+cfg.ReentryWindowBars; confirmIndex++ {
				confirm := candles[confirmIndex]
				if confirm.Close < episode.Low && confirm.Low < breakCandle.Low {
					(*nextEventID)++
					events = append(events, newRangeDiscoveryDirectionalEvent(*nextEventID, frame, episode, confirmIndex, RangeUniverseFamilyStructuredCompressionBreak, RangeDiscoverySideDown, "down", confirmIndex-episode.EndIndex, breakIndex-episode.EndIndex, confirmIndex-breakIndex, 0, RangeUniverseOutcomeStructuredContinuation, RangeUniverseOutcomeStructuredFailedReentry, RangeUniverseOutcomeStructuredNoContinuation))
					break
				}
			}
		}
	}
	return events
}

func summarizeRangeUniverseDiscovery(rows []FuturesRangeUniverseCandidateRow, roundTripCostPct float64) []FuturesRangeUniverseSummaryRow {
	bySymbol := map[string][]FuturesRangeDiscoveryCandidateRow{}
	for _, row := range rows {
		bySymbol[row.Symbol] = append(bySymbol[row.Symbol], row.FuturesRangeDiscoveryCandidateRow)
	}
	out := []FuturesRangeUniverseSummaryRow{}
	for symbol, symbolRows := range bySymbol {
		for _, row := range summarizeRangeDiscovery(symbolRows, roundTripCostPct) {
			out = append(out, FuturesRangeUniverseSummaryRow{Symbol: symbol, FuturesRangeDiscoverySummaryRow: row})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Symbol != out[j].Symbol {
			return rangeUniverseSymbolSortKey(out[i].Symbol) < rangeUniverseSymbolSortKey(out[j].Symbol)
		}
		return lessRangeDiscoverySummary(out[i].FuturesRangeDiscoverySummaryRow, out[j].FuturesRangeDiscoverySummaryRow)
	})
	return out
}

func rangeUniverseStabilityRows(summaryRows []FuturesRangeUniverseSummaryRow, splits []Split) []FuturesRangeUniverseStabilityRow {
	bySymbol := map[string][]FuturesRangeDiscoverySummaryRow{}
	for _, row := range summaryRows {
		bySymbol[row.Symbol] = append(bySymbol[row.Symbol], row.FuturesRangeDiscoverySummaryRow)
	}
	out := []FuturesRangeUniverseStabilityRow{}
	for symbol, rows := range bySymbol {
		for _, row := range rangeDiscoveryStabilityRows(rows, splits) {
			out = append(out, FuturesRangeUniverseStabilityRow{Symbol: symbol, FuturesRangeDiscoveryStabilityRow: row})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Symbol != out[j].Symbol {
			return rangeUniverseSymbolSortKey(out[i].Symbol) < rangeUniverseSymbolSortKey(out[j].Symbol)
		}
		return lessRangeDiscoveryStability(out[i].FuturesRangeDiscoveryStabilityRow, out[j].FuturesRangeDiscoveryStabilityRow)
	})
	return out
}

func rangeUniverseRankingRows(summaryRows []FuturesRangeUniverseSummaryRow, stabilityRows []FuturesRangeUniverseStabilityRow, cfg FuturesRangeUniverseDiscoveryAuditConfig, splits []Split) []FuturesRangeUniverseRankingRow {
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	periodSplits := rangeDiscoveryPeriodSplits(splits)
	summaryByKey := map[string]map[rangeDiscoverySummaryKey]FuturesRangeUniverseSummaryRow{}
	for _, row := range summaryRows {
		if summaryByKey[row.Symbol] == nil {
			summaryByKey[row.Symbol] = map[rangeDiscoverySummaryKey]FuturesRangeUniverseSummaryRow{}
		}
		key := rangeDiscoverySummaryKey{split: row.Split, timeframe: row.Timeframe, family: row.Family, side: row.Side, horizonBars: row.HorizonBars}
		summaryByKey[row.Symbol][key] = row
	}
	stabilityByKey := map[string]map[rangeDiscoveryStabilityKey]FuturesRangeUniverseStabilityRow{}
	rankingKeys := map[rangeUniverseRankingKey]bool{}
	for _, row := range stabilityRows {
		if stabilityByKey[row.Symbol] == nil {
			stabilityByKey[row.Symbol] = map[rangeDiscoveryStabilityKey]FuturesRangeUniverseStabilityRow{}
		}
		stabilityKey := rangeDiscoveryStabilityKey{timeframe: row.Timeframe, family: row.Family, side: row.Side, horizonBars: row.HorizonBars}
		stabilityByKey[row.Symbol][stabilityKey] = row
		rankingKeys[rangeUniverseRankingKey{timeframe: row.Timeframe, family: row.Family, side: row.Side, horizonBars: row.HorizonBars}] = true
	}
	keys := make([]rangeUniverseRankingKey, 0, len(rankingKeys))
	for key := range rankingKeys {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return lessRangeUniverseRankingKey(keys[i], keys[j])
	})

	rows := []FuturesRangeUniverseRankingRow{}
	for _, key := range keys {
		row := FuturesRangeUniverseRankingRow{
			Timeframe:                            key.timeframe,
			Family:                               key.family,
			Side:                                 key.side,
			HorizonBars:                          key.horizonBars,
			PeriodSplitsRequired:                 len(periodSplits),
			MinCandidatesPerSplit:                cfg.Discovery.MinCandidatesPerSplit,
			SymbolSpecificMinCandidatesPerSplit:  cfg.SymbolSpecificMinCandidatesPerSplit,
			SymbolSpecificCostBufferMultiple:     cfg.SymbolSpecificCostBufferMultiple,
			CombinedWeakestSplitCandidateCount:   math.MaxInt,
			CombinedWeakestFavorableRate:         math.Inf(1),
			CombinedWeakestFavorableMinusAdverse: math.Inf(1),
			CombinedWeakestCostBufferPct:         math.Inf(1),
		}
		gates := []rangeUniverseSymbolGate{}
		for _, source := range cfg.Sources {
			symbol := strings.ToUpper(source.Symbol)
			gate := rangeUniverseSymbolGateForKey(symbol, key, summaryByKey[symbol], stabilityByKey[symbol], cfg, splits)
			gates = append(gates, gate)
			if gate.passes {
				row.SymbolGatePassCount++
				if symbol == RangeUniverseSymbolBTCUSDT {
					row.BTCUSDTGatePass = true
				} else {
					row.TransferSymbolGatePassCount++
				}
			}
			if gate.symbolException && (!row.SymbolSpecificException || gate.score > rangeUniverseSymbolExceptionScore(row)) {
				row.SymbolSpecificException = true
				row.SymbolSpecificBestSymbol = symbol
				row.SymbolSpecificWeakestSplitCount = gate.stability.CandidateCountMin
				row.SymbolSpecificWeakestCostBufferPct = gate.stability.CostBufferPctMin
			}
			if symbol == RangeUniverseSymbolBTCUSDT {
				fillRangeUniverseBTCFields(&row, gate)
			} else if gate.score > rangeUniverseTransferScore(row) {
				fillRangeUniverseTransferFields(&row, gate)
			}
			row.CombinedFullCandidateCount += gate.full.CandidateCount
			if gate.stability.CandidateCountMin < row.CombinedWeakestSplitCandidateCount {
				row.CombinedWeakestSplitCandidateCount = gate.stability.CandidateCountMin
			}
			row.CombinedWeakestFavorableRate = math.Min(row.CombinedWeakestFavorableRate, gate.stability.FavorableRateMin)
			row.CombinedWorstAdverseRate = math.Max(row.CombinedWorstAdverseRate, gate.stability.AdverseRateMax)
			row.CombinedWorstQuickInvalidationRate = math.Max(row.CombinedWorstQuickInvalidationRate, gate.stability.QuickInvalidationRateMax)
			row.CombinedWeakestFavorableMinusAdverse = math.Min(row.CombinedWeakestFavorableMinusAdverse, gate.stability.FavorableMinusAdversePctMin)
			row.CombinedWeakestCostBufferPct = math.Min(row.CombinedWeakestCostBufferPct, gate.stability.CostBufferPctMin)
		}
		if row.CombinedWeakestSplitCandidateCount == math.MaxInt {
			row.CombinedWeakestSplitCandidateCount = 0
		}
		if math.IsInf(row.CombinedWeakestFavorableRate, 1) {
			row.CombinedWeakestFavorableRate = 0
		}
		if math.IsInf(row.CombinedWeakestFavorableMinusAdverse, 1) {
			row.CombinedWeakestFavorableMinusAdverse = 0
		}
		if math.IsInf(row.CombinedWeakestCostBufferPct, 1) {
			row.CombinedWeakestCostBufferPct = 0
		}
		row.PassesGate = (row.BTCUSDTGatePass && row.TransferSymbolGatePassCount >= 1) || row.SymbolSpecificException
		row.RankScore = rangeUniverseRankScore(row)
		if !row.PassesGate {
			row.FailureReason = rangeUniverseFailureReason(row, gates)
		}
		rows = append(rows, row)
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].PassesGate != rows[j].PassesGate {
			return rows[i].PassesGate
		}
		if rows[i].RankScore != rows[j].RankScore {
			return rows[i].RankScore > rows[j].RankScore
		}
		return lessRangeUniverseRankingRows(rows[i], rows[j])
	})
	for i := range rows {
		rows[i].Rank = i + 1
	}
	return rows
}

func FuturesRangeUniverseReviewStopState(rankings []FuturesRangeUniverseRankingRow) string {
	for _, row := range rankings {
		if row.PassesGate {
			return RangeUniverseStopStateAuditReady
		}
	}
	return RangeUniverseStopStateNoBacktestCandidate
}

func (cfg FuturesRangeUniverseDiscoveryAuditConfig) withDefaults() FuturesRangeUniverseDiscoveryAuditConfig {
	defaults := DefaultFuturesRangeUniverseDiscoveryAuditConfig()
	if len(cfg.Sources) == 0 {
		cfg.Sources = defaults.Sources
	}
	cfg.Discovery = cfg.Discovery.withDefaults()
	if cfg.SymbolSpecificMinCandidatesPerSplit == 0 {
		cfg.SymbolSpecificMinCandidatesPerSplit = defaults.SymbolSpecificMinCandidatesPerSplit
	}
	if cfg.SymbolSpecificCostBufferMultiple == 0 {
		cfg.SymbolSpecificCostBufferMultiple = defaults.SymbolSpecificCostBufferMultiple
	}
	return cfg
}

func (cfg FuturesRangeUniverseDiscoveryAuditConfig) validate() error {
	if len(cfg.Sources) == 0 {
		return fmt.Errorf("range universe source list cannot be empty")
	}
	if err := cfg.Discovery.validate(); err != nil {
		return err
	}
	if cfg.SymbolSpecificMinCandidatesPerSplit <= 0 {
		return fmt.Errorf("range universe symbol-specific min candidates must be positive")
	}
	if cfg.SymbolSpecificCostBufferMultiple <= 0 {
		return fmt.Errorf("range universe symbol-specific cost buffer multiple must be positive")
	}
	seen := map[string]bool{}
	for _, source := range cfg.Sources {
		symbol := strings.ToUpper(strings.TrimSpace(source.Symbol))
		if symbol == "" {
			return fmt.Errorf("range universe source symbol is required")
		}
		if !rangeUniverseApprovedSymbol(symbol) {
			return fmt.Errorf("range universe symbol %q is not approved", symbol)
		}
		if seen[symbol] {
			return fmt.Errorf("range universe source symbol %q is duplicated", symbol)
		}
		seen[symbol] = true
	}
	return nil
}

func validateRangeUniverseSourcePath(symbol string, path string) error {
	base := strings.ToLower(filepath.Base(path))
	if !strings.Contains(base, strings.ToLower(symbol)) {
		return fmt.Errorf("%s source CSV path must identify %s", symbol, symbol)
	}
	if !strings.Contains(base, "5m") {
		return fmt.Errorf("%s source CSV path must identify 5m candles", symbol)
	}
	if strings.Contains(base, "spot") {
		return fmt.Errorf("%s spot-looking CSV path %q cannot be used for range universe futures research", symbol, path)
	}
	if !(strings.Contains(base, "futures") || strings.Contains(base, "usdm") || strings.Contains(base, "usd-m") || strings.Contains(base, "_um_")) {
		return fmt.Errorf("%s source CSV path must identify Binance USDT-M futures", symbol)
	}
	return nil
}

func validateRangeUniverseSortedCandles(symbol string, candles []Candle, row *FuturesRangeUniverseSourceRow, splits []Split) error {
	seen := map[time.Time]struct{}{}
	for i, candle := range candles {
		if err := validateCandleValues(i, candle); err != nil {
			return fmt.Errorf("%s %s", symbol, err.Error())
		}
		if candle.Volume == 0 {
			row.ZeroVolumeCount++
		}
		openTime := candle.OpenTime.UTC()
		if _, ok := seen[openTime]; ok {
			row.DuplicateCount++
			return fmt.Errorf("%s source CSV contains %d duplicate open_time value(s)", symbol, row.DuplicateCount)
		}
		seen[openTime] = struct{}{}

		expectedClose := openTime.Add(sourceInterval - time.Millisecond)
		if !candle.CloseTime.UTC().Equal(expectedClose) {
			return fmt.Errorf("%s row %d close_time=%s does not match expected closed 5m candle close_time=%s", symbol, i, candle.CloseTime.UTC().Format(time.RFC3339Nano), expectedClose.Format(time.RFC3339Nano))
		}
		for _, split := range splits {
			if split.Name == fullSplitName || !split.Contains(candle.CloseTime) {
				continue
			}
			switch split.Name {
			case "2021_2022_stress":
				row.StressSplitRows++
			case "2023_2024_oos":
				row.OOSSplitRows++
			case "2025_2026_recent":
				row.RecentSplitRows++
			}
		}
		if i == 0 {
			continue
		}
		prevOpen := candles[i-1].OpenTime.UTC()
		if !openTime.After(prevOpen) {
			row.AcceptedMonotonic = false
			return fmt.Errorf("%s accepted source open_time values must be strictly increasing after validation sort", symbol)
		}
		diff := openTime.Sub(prevOpen)
		if diff != sourceInterval {
			if diff > sourceInterval && diff%sourceInterval == 0 {
				row.GapCount += int(diff/sourceInterval) - 1
				return fmt.Errorf("%s source CSV has %d missing 5m candle(s)", symbol, row.GapCount)
			}
			return fmt.Errorf("%s irregular 5m interval between %s and %s: got %s", symbol, prevOpen.Format(time.RFC3339), openTime.Format(time.RFC3339), diff)
		}
	}
	return nil
}

func physicalNonMonotonicCount(candles []Candle) int {
	count := 0
	for i := 1; i < len(candles); i++ {
		if !candles[i].OpenTime.UTC().After(candles[i-1].OpenTime.UTC()) {
			count++
		}
	}
	return count
}

func sameCleanPath(a, b string) bool {
	return filepath.Clean(a) == filepath.Clean(b)
}

func rangeUniverseApprovedSymbol(symbol string) bool {
	switch symbol {
	case RangeUniverseSymbolBTCUSDT, RangeUniverseSymbolETHUSDT, RangeUniverseSymbolSOLUSDT:
		return true
	default:
		return false
	}
}

func rangeUniverseSymbolGateForKey(symbol string, key rangeUniverseRankingKey, summaryByKey map[rangeDiscoverySummaryKey]FuturesRangeUniverseSummaryRow, stabilityByKey map[rangeDiscoveryStabilityKey]FuturesRangeUniverseStabilityRow, cfg FuturesRangeUniverseDiscoveryAuditConfig, splits []Split) rangeUniverseSymbolGate {
	gate := rangeUniverseSymbolGate{symbol: symbol}
	if summaryByKey == nil || stabilityByKey == nil {
		gate.reasons = append(gate.reasons, symbol+"_missing_surface")
		return gate
	}
	fullKey := rangeDiscoverySummaryKey{split: fullSplitName, timeframe: key.timeframe, family: key.family, side: key.side, horizonBars: key.horizonBars}
	stabilityKey := rangeDiscoveryStabilityKey{timeframe: key.timeframe, family: key.family, side: key.side, horizonBars: key.horizonBars}
	gate.full = summaryByKey[fullKey]
	gate.stability = stabilityByKey[stabilityKey]
	if gate.stability.PeriodSplits == 0 {
		gate.reasons = append(gate.reasons, symbol+"_missing_stability")
		return gate
	}
	passingSplits := 0
	symbolExceptionSplits := 0
	for _, split := range rangeDiscoveryPeriodSplits(splits) {
		summary, ok := summaryByKey[rangeDiscoverySummaryKey{split: split.Name, timeframe: key.timeframe, family: key.family, side: key.side, horizonBars: key.horizonBars}]
		if !ok {
			gate.reasons = append(gate.reasons, symbol+"_missing_split")
			continue
		}
		if rangeUniverseSummaryPasses(summary.FuturesRangeDiscoverySummaryRow, cfg.Discovery.MinCandidatesPerSplit, cfg.Discovery.RoundTripCostPct) {
			passingSplits++
		} else {
			gate.reasons = append(gate.reasons, symbol+"_split_gate_fail")
		}
		if rangeUniverseSummaryPasses(summary.FuturesRangeDiscoverySummaryRow, cfg.SymbolSpecificMinCandidatesPerSplit, cfg.Discovery.RoundTripCostPct*cfg.SymbolSpecificCostBufferMultiple) {
			symbolExceptionSplits++
		}
	}
	required := len(rangeDiscoveryPeriodSplits(splits))
	gate.passes = passingSplits == required
	gate.symbolException = symbolExceptionSplits == required
	gate.score = rangeUniverseSymbolGateScore(gate)
	return gate
}

func rangeUniverseSummaryPasses(summary FuturesRangeDiscoverySummaryRow, minCandidates int, minCostBuffer float64) bool {
	if summary.CandidateCount < minCandidates {
		return false
	}
	if summary.FavorableRate <= summary.AdverseRate {
		return false
	}
	if summary.QuickInvalidationRate >= summary.FavorableRate {
		return false
	}
	if summary.AvgFavorableMovePct <= summary.AvgAdverseMovePct {
		return false
	}
	if summary.CostBufferPct <= minCostBuffer {
		return false
	}
	return true
}

func fillRangeUniverseBTCFields(row *FuturesRangeUniverseRankingRow, gate rangeUniverseSymbolGate) {
	row.BTCUSDTWeakestSplitCandidateCount = gate.stability.CandidateCountMin
	row.BTCUSDTWeakestFavorableRate = gate.stability.FavorableRateMin
	row.BTCUSDTWorstAdverseRate = gate.stability.AdverseRateMax
	row.BTCUSDTWorstQuickInvalidationRate = gate.stability.QuickInvalidationRateMax
	row.BTCUSDTWeakestFavorableMinusAdverse = gate.stability.FavorableMinusAdversePctMin
	row.BTCUSDTWeakestCostBufferPct = gate.stability.CostBufferPctMin
}

func fillRangeUniverseTransferFields(row *FuturesRangeUniverseRankingRow, gate rangeUniverseSymbolGate) {
	row.BestTransferSymbol = gate.symbol
	row.BestTransferWeakestSplitCount = gate.stability.CandidateCountMin
	row.BestTransferWeakestFavorableRate = gate.stability.FavorableRateMin
	row.BestTransferWorstAdverseRate = gate.stability.AdverseRateMax
	row.BestTransferWorstQuickInvalidation = gate.stability.QuickInvalidationRateMax
	row.BestTransferWeakestFavorableMinusAdv = gate.stability.FavorableMinusAdversePctMin
	row.BestTransferWeakestCostBufferPct = gate.stability.CostBufferPctMin
}

func rangeUniverseSymbolGateScore(gate rangeUniverseSymbolGate) float64 {
	if gate.stability.PeriodSplits == 0 {
		return math.Inf(-1)
	}
	score := gate.stability.FavorableRateMin - gate.stability.AdverseRateMax
	score += gate.stability.FavorableMinusAdversePctMin * 20
	score -= gate.stability.QuickInvalidationRateMax * 0.5
	score += math.Log1p(float64(gate.stability.CandidateCountMin)) / 20
	if gate.passes {
		score += 1
	}
	return score
}

func rangeUniverseTransferScore(row FuturesRangeUniverseRankingRow) float64 {
	if row.BestTransferSymbol == "" {
		return math.Inf(-1)
	}
	score := row.BestTransferWeakestFavorableRate - row.BestTransferWorstAdverseRate
	score += row.BestTransferWeakestFavorableMinusAdv * 20
	score -= row.BestTransferWorstQuickInvalidation * 0.5
	score += math.Log1p(float64(row.BestTransferWeakestSplitCount)) / 20
	return score
}

func rangeUniverseSymbolExceptionScore(row FuturesRangeUniverseRankingRow) float64 {
	if !row.SymbolSpecificException {
		return math.Inf(-1)
	}
	return row.SymbolSpecificWeakestCostBufferPct + math.Log1p(float64(row.SymbolSpecificWeakestSplitCount))/20
}

func rangeUniverseRankScore(row FuturesRangeUniverseRankingRow) float64 {
	score := row.CombinedWeakestFavorableRate - row.CombinedWorstAdverseRate
	score += row.CombinedWeakestFavorableMinusAdverse * 20
	score -= row.CombinedWorstQuickInvalidationRate * 0.5
	score += math.Log1p(float64(row.CombinedWeakestSplitCandidateCount)) / 20
	if row.BTCUSDTGatePass {
		score += 2
	}
	score += float64(row.TransferSymbolGatePassCount) * 1.5
	if row.SymbolSpecificException {
		score += 0.5
	}
	return score
}

func rangeUniverseFailureReason(row FuturesRangeUniverseRankingRow, gates []rangeUniverseSymbolGate) string {
	reasons := []string{}
	if !row.BTCUSDTGatePass {
		reasons = append(reasons, "btcusdt_gate_fail")
	}
	if row.TransferSymbolGatePassCount == 0 {
		reasons = append(reasons, "no_transfer_symbol_confirmation")
	}
	if !row.SymbolSpecificException {
		reasons = append(reasons, "no_symbol_specific_exception")
	}
	for _, gate := range gates {
		reasons = append(reasons, gate.reasons...)
	}
	return uniqueJoinedReasons(reasons)
}

func lessRangeUniverseRankingKey(a, b rangeUniverseRankingKey) bool {
	if rangeDiscoveryTimeframeSortKey(a.timeframe) != rangeDiscoveryTimeframeSortKey(b.timeframe) {
		return rangeDiscoveryTimeframeSortKey(a.timeframe) < rangeDiscoveryTimeframeSortKey(b.timeframe)
	}
	if rangeDiscoveryFamilySortKey(a.family) != rangeDiscoveryFamilySortKey(b.family) {
		return rangeDiscoveryFamilySortKey(a.family) < rangeDiscoveryFamilySortKey(b.family)
	}
	if rangeDiscoverySideSortKey(a.side) != rangeDiscoverySideSortKey(b.side) {
		return rangeDiscoverySideSortKey(a.side) < rangeDiscoverySideSortKey(b.side)
	}
	return a.horizonBars < b.horizonBars
}

func lessRangeUniverseRankingRows(a, b FuturesRangeUniverseRankingRow) bool {
	return lessRangeUniverseRankingKey(
		rangeUniverseRankingKey{timeframe: a.Timeframe, family: a.Family, side: a.Side, horizonBars: a.HorizonBars},
		rangeUniverseRankingKey{timeframe: b.Timeframe, family: b.Family, side: b.Side, horizonBars: b.HorizonBars},
	)
}

func rangeUniverseSymbolSortKey(symbol string) int {
	switch symbol {
	case RangeUniverseSymbolBTCUSDT:
		return 0
	case RangeUniverseSymbolETHUSDT:
		return 1
	case RangeUniverseSymbolSOLUSDT:
		return 2
	default:
		return 99
	}
}
