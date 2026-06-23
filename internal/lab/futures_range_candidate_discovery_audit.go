package lab

import (
	"fmt"
	"math"
	"sort"
	"time"
)

const (
	RangeDiscoveryTimeframe5m  = "5m"
	RangeDiscoveryTimeframe15m = "15m"
	RangeDiscoveryTimeframe1h  = "1h"
	RangeDiscoveryTimeframe4h  = "4h"

	RangeDiscoveryFamilyBalancePersistence = "mature_balance_persistence"
	RangeDiscoveryFamilyBoundaryTouch      = "boundary_touch_rejection"
	RangeDiscoveryFamilyWickRejection      = "single_candle_wick_rejection"
	RangeDiscoveryFamilyFailedBreakReentry = "failed_breakout_reentry"
	RangeDiscoveryFamilyCleanBreakout      = "clean_breakout_continuation"

	RangeDiscoverySideAll        = "all"
	RangeDiscoverySideBalance    = "balance"
	RangeDiscoverySideSupport    = "support"
	RangeDiscoverySideResistance = "resistance"
	RangeDiscoverySideUp         = "up"
	RangeDiscoverySideDown       = "down"

	RangeDiscoveryClassFavorable = "favorable"
	RangeDiscoveryClassAdverse   = "adverse"
	RangeDiscoveryClassNeutral   = "neutral"
	RangeDiscoveryClassMissing   = "missing_future"

	RangeDiscoveryOutcomeInsidePersistence       = "inside_persistence"
	RangeDiscoveryOutcomeInternalRotation        = "internal_rotation"
	RangeDiscoveryOutcomeExpansionFailure        = "expansion_failure"
	RangeDiscoveryOutcomeRejectInward            = "reject_inward"
	RangeDiscoveryOutcomeAcceptOutside           = "accept_outside"
	RangeDiscoveryOutcomeStallNone               = "stall_none"
	RangeDiscoveryOutcomeInwardFollowThrough     = "inward_follow_through"
	RangeDiscoveryOutcomeBoundaryBreak           = "boundary_break"
	RangeDiscoveryOutcomeNoFollowThrough         = "no_follow_through"
	RangeDiscoveryOutcomeReentryContinuation     = "reentry_continuation"
	RangeDiscoveryOutcomeSecondBreak             = "second_break"
	RangeDiscoveryOutcomeContinuationBeyondBreak = "continuation_beyond_break"
	RangeDiscoveryOutcomeFailedBreakReentry      = "failed_break_reentry"
	RangeDiscoveryOutcomeNoExtension             = "no_extension"
	RangeDiscoveryOutcomeMissingFuture           = "missing_future"

	RangeDiscoveryStopStateAuditReady           = "range_discovery_audit_ready"
	RangeDiscoveryStopStateNoBacktestCandidate  = "range_discovery_no_backtest_candidate"
	RangeDiscoveryStopStateSourceOrResampleGap  = "range_discovery_source_or_resample_gap"
	RangeDiscoveryStopStateCodegenOrTestBlocked = "range_discovery_codegen_or_test_blocked"
	RangeDiscoveryStopStateReviewOnlyNoStrategy = "range_discovery_review_only_no_strategy_change"
)

type FuturesRangeCandidateDiscoveryAuditConfig struct {
	HorizonsBars                 []int
	QuickInvalidationBars        int
	MaxEventDelayBars            int
	ReentryWindowBars            int
	MinMoveRangeFraction         float64
	LargeWickFraction            float64
	RoundTripCostPct             float64
	MinCandidatesPerSplit        int
	DetectorLookbackDays         int
	DetectorLookbackBarsOverride int
	DetectorPercentile           float64
	DetectorMinConsecutiveBars   int
}

type FuturesRangeDiscoveryCoverageRow struct {
	Timeframe             string `json:"timeframe"`
	IntervalMinutes       int    `json:"interval_minutes"`
	ChildBars             int    `json:"child_bars"`
	BarsPerDay            int    `json:"bars_per_day"`
	RowCount              int    `json:"row_count"`
	FirstOpenTime         string `json:"first_open_time"`
	LastOpenTime          string `json:"last_open_time"`
	FirstCloseTime        string `json:"first_close_time"`
	LastCloseTime         string `json:"last_close_time"`
	ExpectedChildBars     int    `json:"expected_child_bars"`
	CompleteBucketCount   int    `json:"complete_bucket_count"`
	PartialFinalChildBars int    `json:"partial_final_child_bars"`
	PartialFinalDropped   bool   `json:"partial_final_dropped"`
	GapCount              int    `json:"gap_count"`
	DuplicateCount        int    `json:"duplicate_count"`
	MissingChildOpenCount int    `json:"missing_child_open_count"`
	Complete              bool   `json:"complete"`
	ValidationStatus      string `json:"validation_status"`
	ValidationError       string `json:"validation_error,omitempty"`
}

type FuturesRangeDiscoveryCandidateRow struct {
	EventID                     int     `json:"event_id"`
	EventIndex                  int     `json:"event_index"`
	EventOpenTime               string  `json:"event_open_time"`
	EventCloseTime              string  `json:"event_close_time"`
	Split                       string  `json:"split"`
	Timeframe                   string  `json:"timeframe"`
	Family                      string  `json:"family"`
	Side                        string  `json:"side"`
	HorizonBars                 int     `json:"horizon_bars"`
	EpisodeID                   int     `json:"episode_id"`
	EpisodeStartIndex           int     `json:"episode_start_index"`
	EpisodeEndIndex             int     `json:"episode_end_index"`
	EpisodeStartTime            string  `json:"episode_start_time"`
	EpisodeEndTime              string  `json:"episode_end_time"`
	EventDelayBars              int     `json:"event_delay_bars"`
	BreakoutDelayBars           int     `json:"breakout_delay_bars"`
	ReentryDelayBars            int     `json:"reentry_delay_bars"`
	RangeHigh                   float64 `json:"range_high"`
	RangeLow                    float64 `json:"range_low"`
	RangeMid                    float64 `json:"range_mid"`
	RangeWidthPct               float64 `json:"range_width_pct"`
	EpisodeRawLengthBars        int     `json:"episode_raw_length_bars"`
	EpisodeActiveLengthBars     int     `json:"episode_active_length_bars"`
	EpisodeRawLengthBucket      string  `json:"episode_raw_length_bucket"`
	EpisodeActiveLengthBucket   string  `json:"episode_active_length_bucket"`
	EpisodeRangeWidthBucket     string  `json:"episode_range_width_bucket"`
	EventOpen                   float64 `json:"event_open"`
	EventHigh                   float64 `json:"event_high"`
	EventLow                    float64 `json:"event_low"`
	EventClose                  float64 `json:"event_close"`
	EventVolume                 float64 `json:"event_volume"`
	EventClosePosition          float64 `json:"event_close_position"`
	EventWickFraction           float64 `json:"event_wick_fraction"`
	LabelWindowStartIndex       int     `json:"label_window_start_index"`
	LabelWindowEndIndex         int     `json:"label_window_end_index"`
	LabelWindowStartTime        string  `json:"label_window_start_time"`
	LabelWindowEndTime          string  `json:"label_window_end_time"`
	OutcomeLabel                string  `json:"outcome_label"`
	OutcomeClass                string  `json:"outcome_class"`
	FavorableMovePct            float64 `json:"favorable_move_pct"`
	AdverseMovePct              float64 `json:"adverse_move_pct"`
	FavorableMinusAdversePct    float64 `json:"favorable_minus_adverse_pct"`
	FavorableGreaterThanAdverse bool    `json:"favorable_greater_than_adverse"`
	QuickInvalidation           bool    `json:"quick_invalidation"`
	MissingFuture               bool    `json:"missing_future"`
	BarsToFavorable             int     `json:"bars_to_favorable"`
	BarsToAdverse               int     `json:"bars_to_adverse"`
}

type FuturesRangeDiscoverySummaryRow struct {
	Split                       string  `json:"split"`
	Timeframe                   string  `json:"timeframe"`
	Family                      string  `json:"family"`
	Side                        string  `json:"side"`
	HorizonBars                 int     `json:"horizon_bars"`
	CandidateCount              int     `json:"candidate_count"`
	LabeledCount                int     `json:"labeled_count"`
	MissingFutureCount          int     `json:"missing_future_count"`
	FavorableCount              int     `json:"favorable_count"`
	AdverseCount                int     `json:"adverse_count"`
	NeutralCount                int     `json:"neutral_count"`
	QuickInvalidationCount      int     `json:"quick_invalidation_count"`
	FavorableRate               float64 `json:"favorable_rate"`
	AdverseRate                 float64 `json:"adverse_rate"`
	NeutralRate                 float64 `json:"neutral_rate"`
	QuickInvalidationRate       float64 `json:"quick_invalidation_rate"`
	AvgFavorableMovePct         float64 `json:"avg_favorable_move_pct"`
	AvgAdverseMovePct           float64 `json:"avg_adverse_move_pct"`
	AvgFavorableMinusAdversePct float64 `json:"avg_favorable_minus_adverse_pct"`
	AvgRangeWidthPct            float64 `json:"avg_range_width_pct"`
	RoundTripCostPct            float64 `json:"round_trip_cost_pct"`
	CostBufferPct               float64 `json:"cost_buffer_pct"`
	CostBufferPass              bool    `json:"cost_buffer_pass"`
}

type FuturesRangeDiscoveryStabilityRow struct {
	Timeframe                     string  `json:"timeframe"`
	Family                        string  `json:"family"`
	Side                          string  `json:"side"`
	HorizonBars                   int     `json:"horizon_bars"`
	PeriodSplits                  int     `json:"period_splits"`
	CandidateCount                int     `json:"candidate_count"`
	CandidateCountMin             int     `json:"candidate_count_min"`
	CandidateCountMax             int     `json:"candidate_count_max"`
	CandidateCountDelta           int     `json:"candidate_count_delta"`
	FavorableRateMin              float64 `json:"favorable_rate_min"`
	FavorableRateMax              float64 `json:"favorable_rate_max"`
	FavorableRateDelta            float64 `json:"favorable_rate_delta"`
	AdverseRateMin                float64 `json:"adverse_rate_min"`
	AdverseRateMax                float64 `json:"adverse_rate_max"`
	AdverseRateDelta              float64 `json:"adverse_rate_delta"`
	QuickInvalidationRateMax      float64 `json:"quick_invalidation_rate_max"`
	FavorableMinusAdversePctMin   float64 `json:"favorable_minus_adverse_pct_min"`
	FavorableMinusAdversePctMax   float64 `json:"favorable_minus_adverse_pct_max"`
	FavorableMinusAdversePctDelta float64 `json:"favorable_minus_adverse_pct_delta"`
	CostBufferPctMin              float64 `json:"cost_buffer_pct_min"`
	CostBufferPassAllSplits       bool    `json:"cost_buffer_pass_all_splits"`
}

type FuturesRangeDiscoveryRankingRow struct {
	Rank                            int     `json:"rank"`
	Timeframe                       string  `json:"timeframe"`
	Family                          string  `json:"family"`
	Side                            string  `json:"side"`
	HorizonBars                     int     `json:"horizon_bars"`
	PassesGate                      bool    `json:"passes_gate"`
	RankScore                       float64 `json:"rank_score"`
	FullCandidateCount              int     `json:"full_candidate_count"`
	WeakestSplitCandidateCount      int     `json:"weakest_split_candidate_count"`
	FullFavorableRate               float64 `json:"full_favorable_rate"`
	WeakestSplitFavorableRate       float64 `json:"weakest_split_favorable_rate"`
	WorstSplitAdverseRate           float64 `json:"worst_split_adverse_rate"`
	WorstSplitQuickInvalidationRate float64 `json:"worst_split_quick_invalidation_rate"`
	FullFavorableMinusAdversePct    float64 `json:"full_favorable_minus_adverse_pct"`
	WeakestFavorableMinusAdversePct float64 `json:"weakest_favorable_minus_adverse_pct"`
	WeakestCostBufferPct            float64 `json:"weakest_cost_buffer_pct"`
	PeriodSplitsPassing             int     `json:"period_splits_passing"`
	PeriodSplitsRequired            int     `json:"period_splits_required"`
	FailureReason                   string  `json:"failure_reason,omitempty"`
}

type rangeDiscoveryFrameDef struct {
	timeframe  string
	interval   time.Duration
	childBars  int
	barsPerDay int
	families   []string
}

type rangeDiscoveryEvent struct {
	eventID           int
	timeframe         string
	family            string
	side              string
	eventIndex        int
	episode           rangeRegimeDurabilityEpisode
	eventDelayBars    int
	breakoutDelayBars int
	reentryDelayBars  int
	eventWickFraction float64
	favorableLabel    string
	adverseLabel      string
	neutralLabel      string
	expectedDirection string
	adverseBoundary   string
}

type rangeDiscoveryLabel struct {
	labelWindowStartIndex       int
	labelWindowEndIndex         int
	labelWindowStartTime        string
	labelWindowEndTime          string
	outcomeLabel                string
	outcomeClass                string
	favorableMovePct            float64
	adverseMovePct              float64
	favorableMinusAdversePct    float64
	favorableGreaterThanAdverse bool
	quickInvalidation           bool
	missingFuture               bool
	barsToFavorable             int
	barsToAdverse               int
}

type rangeDiscoverySummaryKey struct {
	split       string
	timeframe   string
	family      string
	side        string
	horizonBars int
}

type rangeDiscoverySummaryAccumulator struct {
	key               rangeDiscoverySummaryKey
	candidates        int
	labeled           int
	missingFuture     int
	favorable         int
	adverse           int
	neutral           int
	quickInvalidation int
	favorableMoveSum  float64
	adverseMoveSum    float64
	favorableMinusSum float64
	rangeWidthSum     float64
	roundTripCostPct  float64
}

type rangeDiscoveryStabilityKey struct {
	timeframe   string
	family      string
	side        string
	horizonBars int
}

func DefaultFuturesRangeCandidateDiscoveryAuditConfig() FuturesRangeCandidateDiscoveryAuditConfig {
	return FuturesRangeCandidateDiscoveryAuditConfig{
		HorizonsBars:               []int{3, 6, 12},
		QuickInvalidationBars:      3,
		MaxEventDelayBars:          24,
		ReentryWindowBars:          3,
		MinMoveRangeFraction:       0.25,
		LargeWickFraction:          0.40,
		RoundTripCostPct:           2 * (0.0004 + 0.000116),
		MinCandidatesPerSplit:      100,
		DetectorLookbackDays:       20,
		DetectorPercentile:         0.30,
		DetectorMinConsecutiveBars: 12,
	}
}

func RunFuturesRangeCandidateDiscoveryAudit(candles []Candle, cfg FuturesRangeCandidateDiscoveryAuditConfig, splits []Split) ([]FuturesRangeDiscoveryCoverageRow, []FuturesRangeDiscoveryCandidateRow, []FuturesRangeDiscoverySummaryRow, []FuturesRangeDiscoveryRankingRow, []FuturesRangeDiscoveryStabilityRow, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return nil, nil, nil, nil, nil, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}

	coverageRows := []FuturesRangeDiscoveryCoverageRow{}
	candidateRows := []FuturesRangeDiscoveryCandidateRow{}
	eventID := 0
	for _, frame := range rangeDiscoveryFrameDefs() {
		frameCandles, coverage, err := resampleRangeDiscoveryFrame(candles, frame)
		coverageRows = append(coverageRows, coverage)
		if err != nil {
			return coverageRows, nil, nil, nil, nil, err
		}
		if len(frameCandles) == 0 {
			continue
		}
		events, err := rangeDiscoveryEventsForFrame(frameCandles, frame, cfg, splits, &eventID)
		if err != nil {
			return coverageRows, nil, nil, nil, nil, err
		}
		for _, event := range events {
			for _, horizon := range cfg.HorizonsBars {
				candidateRows = append(candidateRows, newRangeDiscoveryCandidateRow(frameCandles, event, horizon, cfg))
			}
		}
	}
	sort.Slice(candidateRows, func(i, j int) bool {
		return lessRangeDiscoveryCandidate(candidateRows[i], candidateRows[j])
	})
	summaryRows := summarizeRangeDiscovery(candidateRows, cfg.RoundTripCostPct)
	stabilityRows := rangeDiscoveryStabilityRows(summaryRows, splits)
	rankingRows := rangeDiscoveryRankingRows(summaryRows, stabilityRows, cfg, splits)
	return coverageRows, candidateRows, summaryRows, rankingRows, stabilityRows, nil
}

func FuturesRangeDiscoveryReviewStopState(rankings []FuturesRangeDiscoveryRankingRow) string {
	for _, row := range rankings {
		if row.PassesGate {
			return RangeDiscoveryStopStateAuditReady
		}
	}
	return RangeDiscoveryStopStateNoBacktestCandidate
}

func (cfg FuturesRangeCandidateDiscoveryAuditConfig) withDefaults() FuturesRangeCandidateDiscoveryAuditConfig {
	defaults := DefaultFuturesRangeCandidateDiscoveryAuditConfig()
	if len(cfg.HorizonsBars) == 0 {
		cfg.HorizonsBars = append([]int(nil), defaults.HorizonsBars...)
	}
	if cfg.QuickInvalidationBars == 0 {
		cfg.QuickInvalidationBars = defaults.QuickInvalidationBars
	}
	if cfg.MaxEventDelayBars == 0 {
		cfg.MaxEventDelayBars = defaults.MaxEventDelayBars
	}
	if cfg.ReentryWindowBars == 0 {
		cfg.ReentryWindowBars = defaults.ReentryWindowBars
	}
	if cfg.MinMoveRangeFraction == 0 {
		cfg.MinMoveRangeFraction = defaults.MinMoveRangeFraction
	}
	if cfg.LargeWickFraction == 0 {
		cfg.LargeWickFraction = defaults.LargeWickFraction
	}
	if cfg.RoundTripCostPct == 0 {
		cfg.RoundTripCostPct = defaults.RoundTripCostPct
	}
	if cfg.MinCandidatesPerSplit == 0 {
		cfg.MinCandidatesPerSplit = defaults.MinCandidatesPerSplit
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
	return cfg
}

func (cfg FuturesRangeCandidateDiscoveryAuditConfig) validate() error {
	if cfg.QuickInvalidationBars <= 0 {
		return fmt.Errorf("range discovery quick invalidation bars must be positive")
	}
	if cfg.MaxEventDelayBars <= 0 {
		return fmt.Errorf("range discovery max event delay bars must be positive")
	}
	if cfg.ReentryWindowBars <= 0 {
		return fmt.Errorf("range discovery reentry window bars must be positive")
	}
	if cfg.MinMoveRangeFraction <= 0 {
		return fmt.Errorf("range discovery min move range fraction must be positive")
	}
	if cfg.LargeWickFraction <= 0 || cfg.LargeWickFraction > 1 {
		return fmt.Errorf("range discovery large wick fraction must be in (0,1]")
	}
	if cfg.RoundTripCostPct <= 0 {
		return fmt.Errorf("range discovery round trip cost pct must be positive")
	}
	if cfg.MinCandidatesPerSplit <= 0 {
		return fmt.Errorf("range discovery min candidates per split must be positive")
	}
	if cfg.DetectorLookbackDays <= 0 {
		return fmt.Errorf("range discovery detector lookback days must be positive")
	}
	if cfg.DetectorLookbackBarsOverride < 0 {
		return fmt.Errorf("range discovery detector lookback bars override cannot be negative")
	}
	if cfg.DetectorPercentile <= 0 || cfg.DetectorPercentile >= 1 {
		return fmt.Errorf("range discovery detector percentile must be between 0 and 1")
	}
	if cfg.DetectorMinConsecutiveBars <= 0 {
		return fmt.Errorf("range discovery detector min consecutive bars must be positive")
	}
	for _, horizon := range cfg.HorizonsBars {
		if horizon <= 0 {
			return fmt.Errorf("range discovery horizon bars must be positive")
		}
	}
	return nil
}

func rangeDiscoveryFrameDefs() []rangeDiscoveryFrameDef {
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
				RangeDiscoveryFamilyBoundaryTouch,
				RangeDiscoveryFamilyWickRejection,
				RangeDiscoveryFamilyFailedBreakReentry,
				RangeDiscoveryFamilyCleanBreakout,
			},
		},
		{
			timeframe:  RangeDiscoveryTimeframe1h,
			interval:   time.Hour,
			childBars:  12,
			barsPerDay: 24,
			families: []string{
				RangeDiscoveryFamilyBalancePersistence,
				RangeDiscoveryFamilyBoundaryTouch,
				RangeDiscoveryFamilyWickRejection,
				RangeDiscoveryFamilyFailedBreakReentry,
				RangeDiscoveryFamilyCleanBreakout,
			},
		},
		{
			timeframe:  RangeDiscoveryTimeframe4h,
			interval:   4 * time.Hour,
			childBars:  48,
			barsPerDay: 6,
			families: []string{
				RangeDiscoveryFamilyBalancePersistence,
				RangeDiscoveryFamilyCleanBreakout,
			},
		},
	}
}

func resampleRangeDiscoveryFrame(parent []Candle, frame rangeDiscoveryFrameDef) ([]Candle, FuturesRangeDiscoveryCoverageRow, error) {
	row := FuturesRangeDiscoveryCoverageRow{
		Timeframe:         frame.timeframe,
		IntervalMinutes:   int(frame.interval / time.Minute),
		ChildBars:         frame.childBars,
		BarsPerDay:        frame.barsPerDay,
		ExpectedChildBars: frame.childBars,
		ValidationStatus:  "accepted",
		Complete:          true,
	}
	reject := func(format string, args ...any) ([]Candle, FuturesRangeDiscoveryCoverageRow, error) {
		row.Complete = false
		row.ValidationStatus = "rejected"
		row.ValidationError = fmt.Sprintf(format, args...)
		return nil, row, fmt.Errorf("%s", row.ValidationError)
	}
	if frame.childBars <= 0 || frame.interval <= 0 {
		return reject("invalid range discovery timeframe definition for %s", frame.timeframe)
	}
	if len(parent) == 0 {
		row.RowCount = 0
		return nil, row, nil
	}
	if err := validateRangeDiscoveryParentCadence(parent, &row); err != nil {
		return reject(err.Error())
	}
	if frame.childBars == 1 {
		out := append([]Candle(nil), parent...)
		fillRangeDiscoveryCoverageTimes(&row, out)
		row.RowCount = len(out)
		row.CompleteBucketCount = len(out)
		return out, row, nil
	}

	out := make([]Candle, 0, len(parent)/frame.childBars)
	for start := 0; start < len(parent); {
		remaining := len(parent) - start
		if remaining < frame.childBars {
			row.PartialFinalChildBars = remaining
			row.PartialFinalDropped = true
			break
		}
		bucketStart := parent[start].OpenTime.UTC()
		if !bucketStart.Equal(bucketStart.Truncate(frame.interval)) {
			return reject("%s bucket start %s is not aligned to %s", frame.timeframe, bucketStart.Format(time.RFC3339), frame.interval)
		}
		high := parent[start].High
		low := parent[start].Low
		volume := 0.0
		for child := 0; child < frame.childBars; child++ {
			c := parent[start+child]
			expected := bucketStart.Add(time.Duration(child) * sourceInterval)
			if !c.OpenTime.UTC().Equal(expected) {
				row.MissingChildOpenCount++
				return reject("%s bucket %s missing child open %s", frame.timeframe, bucketStart.Format(time.RFC3339), expected.Format(time.RFC3339))
			}
			if c.High > high {
				high = c.High
			}
			if c.Low < low {
				low = c.Low
			}
			volume += c.Volume
		}
		last := parent[start+frame.childBars-1]
		out = append(out, Candle{
			OpenTime:  bucketStart,
			CloseTime: bucketStart.Add(frame.interval - time.Millisecond),
			Open:      parent[start].Open,
			High:      high,
			Low:       low,
			Close:     last.Close,
			Volume:    volume,
		})
		start += frame.childBars
	}
	row.RowCount = len(out)
	row.CompleteBucketCount = len(out)
	fillRangeDiscoveryCoverageTimes(&row, out)
	return out, row, nil
}

func validateRangeDiscoveryParentCadence(candles []Candle, row *FuturesRangeDiscoveryCoverageRow) error {
	seen := map[time.Time]struct{}{}
	for i, candle := range candles {
		open := candle.OpenTime.UTC()
		if _, ok := seen[open]; ok {
			row.DuplicateCount++
		}
		seen[open] = struct{}{}
		if i == 0 {
			continue
		}
		prev := candles[i-1].OpenTime.UTC()
		if !open.After(prev) {
			return fmt.Errorf("range discovery parent candles must be strictly increasing")
		}
		diff := open.Sub(prev)
		if diff != sourceInterval {
			if diff > sourceInterval && diff%sourceInterval == 0 {
				row.GapCount += int(diff/sourceInterval) - 1
				return fmt.Errorf("range discovery parent has %d missing 5m child candle(s)", row.GapCount)
			}
			return fmt.Errorf("range discovery parent has irregular child interval %s between %s and %s", diff, prev.Format(time.RFC3339), open.Format(time.RFC3339))
		}
	}
	if row.DuplicateCount > 0 {
		return fmt.Errorf("range discovery parent has %d duplicate child open(s)", row.DuplicateCount)
	}
	return nil
}

func fillRangeDiscoveryCoverageTimes(row *FuturesRangeDiscoveryCoverageRow, candles []Candle) {
	if len(candles) == 0 {
		return
	}
	row.FirstOpenTime = candles[0].OpenTime.UTC().Format(timeLayout)
	row.LastOpenTime = candles[len(candles)-1].OpenTime.UTC().Format(timeLayout)
	row.FirstCloseTime = candles[0].CloseTime.UTC().Format(timeLayout)
	row.LastCloseTime = candles[len(candles)-1].CloseTime.UTC().Format(timeLayout)
}

func rangeDiscoveryEventsForFrame(candles []Candle, frame rangeDiscoveryFrameDef, cfg FuturesRangeCandidateDiscoveryAuditConfig, splits []Split, nextEventID *int) ([]rangeDiscoveryEvent, error) {
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
			events = append(events, rangeDiscoveryEventsForFamily(candles, episode, frame, family, cfg, nextEventID)...)
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

func rangeDiscoveryEventsForFamily(candles []Candle, episode rangeRegimeDurabilityEpisode, frame rangeDiscoveryFrameDef, family string, cfg FuturesRangeCandidateDiscoveryAuditConfig, nextEventID *int) []rangeDiscoveryEvent {
	switch family {
	case RangeDiscoveryFamilyBalancePersistence:
		return rangeDiscoveryBalanceEvents(candles, episode, frame, nextEventID)
	case RangeDiscoveryFamilyBoundaryTouch:
		return rangeDiscoveryBoundaryTouchEvents(candles, episode, frame, cfg, nextEventID)
	case RangeDiscoveryFamilyWickRejection:
		return rangeDiscoveryWickRejectionEvents(candles, episode, frame, cfg, nextEventID)
	case RangeDiscoveryFamilyFailedBreakReentry:
		return rangeDiscoveryFailedBreakReentryEvents(candles, episode, frame, cfg, nextEventID)
	case RangeDiscoveryFamilyCleanBreakout:
		return rangeDiscoveryCleanBreakoutEvents(candles, episode, frame, cfg, nextEventID)
	default:
		return nil
	}
}

func rangeDiscoveryBalanceEvents(candles []Candle, episode rangeRegimeDurabilityEpisode, frame rangeDiscoveryFrameDef, nextEventID *int) []rangeDiscoveryEvent {
	if episode.EndIndex < 0 || episode.EndIndex >= len(candles) {
		return nil
	}
	(*nextEventID)++
	return []rangeDiscoveryEvent{{
		eventID:        *nextEventID,
		timeframe:      frame.timeframe,
		family:         RangeDiscoveryFamilyBalancePersistence,
		side:           RangeDiscoverySideBalance,
		eventIndex:     episode.EndIndex,
		episode:        episode,
		favorableLabel: RangeDiscoveryOutcomeInternalRotation,
		adverseLabel:   RangeDiscoveryOutcomeExpansionFailure,
		neutralLabel:   RangeDiscoveryOutcomeInsidePersistence,
	}}
}

func rangeDiscoveryBoundaryTouchEvents(candles []Candle, episode rangeRegimeDurabilityEpisode, frame rangeDiscoveryFrameDef, cfg FuturesRangeCandidateDiscoveryAuditConfig, nextEventID *int) []rangeDiscoveryEvent {
	events := []rangeDiscoveryEvent{}
	for i := episode.EndIndex + 1; i < len(candles) && i <= episode.EndIndex+cfg.MaxEventDelayBars; i++ {
		candle := candles[i]
		if candle.Low <= episode.Low && candle.Close >= episode.Low && candle.Close <= episode.High {
			(*nextEventID)++
			events = append(events, newRangeDiscoveryDirectionalEvent(*nextEventID, frame, episode, i, RangeDiscoveryFamilyBoundaryTouch, RangeDiscoverySideSupport, "up", i-episode.EndIndex, 0, 0, 0, RangeDiscoveryOutcomeRejectInward, RangeDiscoveryOutcomeAcceptOutside, RangeDiscoveryOutcomeStallNone))
		}
		if candle.High >= episode.High && candle.Close <= episode.High && candle.Close >= episode.Low {
			(*nextEventID)++
			events = append(events, newRangeDiscoveryDirectionalEvent(*nextEventID, frame, episode, i, RangeDiscoveryFamilyBoundaryTouch, RangeDiscoverySideResistance, "down", i-episode.EndIndex, 0, 0, 0, RangeDiscoveryOutcomeRejectInward, RangeDiscoveryOutcomeAcceptOutside, RangeDiscoveryOutcomeStallNone))
		}
	}
	return events
}

func rangeDiscoveryWickRejectionEvents(candles []Candle, episode rangeRegimeDurabilityEpisode, frame rangeDiscoveryFrameDef, cfg FuturesRangeCandidateDiscoveryAuditConfig, nextEventID *int) []rangeDiscoveryEvent {
	events := []rangeDiscoveryEvent{}
	for i := episode.EndIndex + 1; i < len(candles) && i <= episode.EndIndex+cfg.MaxEventDelayBars; i++ {
		candle := candles[i]
		rangeValue := candle.High - candle.Low
		if rangeValue <= 0 {
			continue
		}
		lowerWick := math.Min(candle.Open, candle.Close) - candle.Low
		upperWick := candle.High - math.Max(candle.Open, candle.Close)
		lowerFraction := lowerWick / rangeValue
		upperFraction := upperWick / rangeValue
		if candle.Low < episode.Low && candle.Close >= episode.Low && candle.Close <= episode.High && lowerFraction >= cfg.LargeWickFraction {
			(*nextEventID)++
			events = append(events, newRangeDiscoveryDirectionalEvent(*nextEventID, frame, episode, i, RangeDiscoveryFamilyWickRejection, RangeDiscoverySideSupport, "up", i-episode.EndIndex, 0, 0, lowerFraction, RangeDiscoveryOutcomeInwardFollowThrough, RangeDiscoveryOutcomeBoundaryBreak, RangeDiscoveryOutcomeNoFollowThrough))
		}
		if candle.High > episode.High && candle.Close <= episode.High && candle.Close >= episode.Low && upperFraction >= cfg.LargeWickFraction {
			(*nextEventID)++
			events = append(events, newRangeDiscoveryDirectionalEvent(*nextEventID, frame, episode, i, RangeDiscoveryFamilyWickRejection, RangeDiscoverySideResistance, "down", i-episode.EndIndex, 0, 0, upperFraction, RangeDiscoveryOutcomeInwardFollowThrough, RangeDiscoveryOutcomeBoundaryBreak, RangeDiscoveryOutcomeNoFollowThrough))
		}
	}
	return events
}

func rangeDiscoveryFailedBreakReentryEvents(candles []Candle, episode rangeRegimeDurabilityEpisode, frame rangeDiscoveryFrameDef, cfg FuturesRangeCandidateDiscoveryAuditConfig, nextEventID *int) []rangeDiscoveryEvent {
	events := []rangeDiscoveryEvent{}
	for breakIndex := episode.EndIndex + 1; breakIndex < len(candles) && breakIndex <= episode.EndIndex+cfg.MaxEventDelayBars; breakIndex++ {
		breakCandle := candles[breakIndex]
		if breakCandle.Close > episode.High {
			for reentryIndex := breakIndex + 1; reentryIndex < len(candles) && reentryIndex <= breakIndex+cfg.ReentryWindowBars; reentryIndex++ {
				reentry := candles[reentryIndex]
				if reentry.Close <= episode.High && reentry.Close >= episode.Low {
					(*nextEventID)++
					events = append(events, newRangeDiscoveryDirectionalEvent(*nextEventID, frame, episode, reentryIndex, RangeDiscoveryFamilyFailedBreakReentry, RangeDiscoverySideResistance, "down", reentryIndex-episode.EndIndex, breakIndex-episode.EndIndex, reentryIndex-breakIndex, 0, RangeDiscoveryOutcomeReentryContinuation, RangeDiscoveryOutcomeSecondBreak, RangeDiscoveryOutcomeNoFollowThrough))
					break
				}
			}
		}
		if breakCandle.Close < episode.Low {
			for reentryIndex := breakIndex + 1; reentryIndex < len(candles) && reentryIndex <= breakIndex+cfg.ReentryWindowBars; reentryIndex++ {
				reentry := candles[reentryIndex]
				if reentry.Close >= episode.Low && reentry.Close <= episode.High {
					(*nextEventID)++
					events = append(events, newRangeDiscoveryDirectionalEvent(*nextEventID, frame, episode, reentryIndex, RangeDiscoveryFamilyFailedBreakReentry, RangeDiscoverySideSupport, "up", reentryIndex-episode.EndIndex, breakIndex-episode.EndIndex, reentryIndex-breakIndex, 0, RangeDiscoveryOutcomeReentryContinuation, RangeDiscoveryOutcomeSecondBreak, RangeDiscoveryOutcomeNoFollowThrough))
					break
				}
			}
		}
	}
	return events
}

func rangeDiscoveryCleanBreakoutEvents(candles []Candle, episode rangeRegimeDurabilityEpisode, frame rangeDiscoveryFrameDef, cfg FuturesRangeCandidateDiscoveryAuditConfig, nextEventID *int) []rangeDiscoveryEvent {
	events := []rangeDiscoveryEvent{}
	for i := episode.EndIndex + 1; i < len(candles) && i <= episode.EndIndex+cfg.MaxEventDelayBars; i++ {
		candle := candles[i]
		if candle.Close > episode.High {
			(*nextEventID)++
			events = append(events, newRangeDiscoveryDirectionalEvent(*nextEventID, frame, episode, i, RangeDiscoveryFamilyCleanBreakout, RangeDiscoverySideUp, "up", i-episode.EndIndex, i-episode.EndIndex, 0, 0, RangeDiscoveryOutcomeContinuationBeyondBreak, RangeDiscoveryOutcomeFailedBreakReentry, RangeDiscoveryOutcomeNoExtension))
		}
		if candle.Close < episode.Low {
			(*nextEventID)++
			events = append(events, newRangeDiscoveryDirectionalEvent(*nextEventID, frame, episode, i, RangeDiscoveryFamilyCleanBreakout, RangeDiscoverySideDown, "down", i-episode.EndIndex, i-episode.EndIndex, 0, 0, RangeDiscoveryOutcomeContinuationBeyondBreak, RangeDiscoveryOutcomeFailedBreakReentry, RangeDiscoveryOutcomeNoExtension))
		}
	}
	return events
}

func newRangeDiscoveryDirectionalEvent(eventID int, frame rangeDiscoveryFrameDef, episode rangeRegimeDurabilityEpisode, eventIndex int, family string, side string, expectedDirection string, eventDelayBars int, breakoutDelayBars int, reentryDelayBars int, wickFraction float64, favorableLabel string, adverseLabel string, neutralLabel string) rangeDiscoveryEvent {
	return rangeDiscoveryEvent{
		eventID:           eventID,
		timeframe:         frame.timeframe,
		family:            family,
		side:              side,
		eventIndex:        eventIndex,
		episode:           episode,
		eventDelayBars:    eventDelayBars,
		breakoutDelayBars: breakoutDelayBars,
		reentryDelayBars:  reentryDelayBars,
		eventWickFraction: wickFraction,
		favorableLabel:    favorableLabel,
		adverseLabel:      adverseLabel,
		neutralLabel:      neutralLabel,
		expectedDirection: expectedDirection,
		adverseBoundary:   side,
	}
}

func newRangeDiscoveryCandidateRow(candles []Candle, event rangeDiscoveryEvent, horizon int, cfg FuturesRangeCandidateDiscoveryAuditConfig) FuturesRangeDiscoveryCandidateRow {
	eventCandle := candles[event.eventIndex]
	label := newRangeDiscoveryLabel(candles, event, horizon, cfg)
	width := event.episode.High - event.episode.Low
	eventClosePosition := 0.0
	if width > 0 {
		eventClosePosition = (eventCandle.Close - event.episode.Low) / width
	}
	return FuturesRangeDiscoveryCandidateRow{
		EventID:                     event.eventID,
		EventIndex:                  event.eventIndex,
		EventOpenTime:               eventCandle.OpenTime.UTC().Format(timeLayout),
		EventCloseTime:              eventCandle.CloseTime.UTC().Format(timeLayout),
		Split:                       splitNameForCloseTime(eventCandle.CloseTime, DefaultSplits()),
		Timeframe:                   event.timeframe,
		Family:                      event.family,
		Side:                        event.side,
		HorizonBars:                 horizon,
		EpisodeID:                   event.episode.EpisodeID,
		EpisodeStartIndex:           event.episode.StartIndex,
		EpisodeEndIndex:             event.episode.EndIndex,
		EpisodeStartTime:            candles[event.episode.StartIndex].CloseTime.UTC().Format(timeLayout),
		EpisodeEndTime:              candles[event.episode.EndIndex].CloseTime.UTC().Format(timeLayout),
		EventDelayBars:              event.eventDelayBars,
		BreakoutDelayBars:           event.breakoutDelayBars,
		ReentryDelayBars:            event.reentryDelayBars,
		RangeHigh:                   event.episode.High,
		RangeLow:                    event.episode.Low,
		RangeMid:                    (event.episode.High + event.episode.Low) / 2,
		RangeWidthPct:               event.episode.WidthPct,
		EpisodeRawLengthBars:        event.episode.RawLengthBars,
		EpisodeActiveLengthBars:     event.episode.ActiveLengthBars,
		EpisodeRawLengthBucket:      event.episode.RawLengthBucket,
		EpisodeActiveLengthBucket:   event.episode.ActiveLengthBucket,
		EpisodeRangeWidthBucket:     event.episode.WidthBucket,
		EventOpen:                   eventCandle.Open,
		EventHigh:                   eventCandle.High,
		EventLow:                    eventCandle.Low,
		EventClose:                  eventCandle.Close,
		EventVolume:                 eventCandle.Volume,
		EventClosePosition:          eventClosePosition,
		EventWickFraction:           event.eventWickFraction,
		LabelWindowStartIndex:       label.labelWindowStartIndex,
		LabelWindowEndIndex:         label.labelWindowEndIndex,
		LabelWindowStartTime:        label.labelWindowStartTime,
		LabelWindowEndTime:          label.labelWindowEndTime,
		OutcomeLabel:                label.outcomeLabel,
		OutcomeClass:                label.outcomeClass,
		FavorableMovePct:            label.favorableMovePct,
		AdverseMovePct:              label.adverseMovePct,
		FavorableMinusAdversePct:    label.favorableMinusAdversePct,
		FavorableGreaterThanAdverse: label.favorableGreaterThanAdverse,
		QuickInvalidation:           label.quickInvalidation,
		MissingFuture:               label.missingFuture,
		BarsToFavorable:             label.barsToFavorable,
		BarsToAdverse:               label.barsToAdverse,
	}
}

func newRangeDiscoveryLabel(candles []Candle, event rangeDiscoveryEvent, horizon int, cfg FuturesRangeCandidateDiscoveryAuditConfig) rangeDiscoveryLabel {
	start := event.eventIndex + 1
	end := event.eventIndex + horizon
	label := rangeDiscoveryLabel{
		labelWindowStartIndex: -1,
		labelWindowEndIndex:   -1,
		outcomeLabel:          RangeDiscoveryOutcomeMissingFuture,
		outcomeClass:          RangeDiscoveryClassMissing,
		missingFuture:         true,
	}
	if start < len(candles) {
		label.labelWindowStartIndex = start
		label.labelWindowStartTime = candles[start].CloseTime.UTC().Format(timeLayout)
		label.labelWindowEndIndex = minInt(end, len(candles)-1)
		label.labelWindowEndTime = candles[label.labelWindowEndIndex].CloseTime.UTC().Format(timeLayout)
	}
	if horizon <= 0 || end >= len(candles) {
		return label
	}
	label.missingFuture = false
	label.outcomeClass = RangeDiscoveryClassNeutral
	label.outcomeLabel = event.neutralLabel
	label.labelWindowEndIndex = end
	label.labelWindowEndTime = candles[end].CloseTime.UTC().Format(timeLayout)
	if event.family == RangeDiscoveryFamilyBalancePersistence {
		return newRangeDiscoveryBalanceLabel(candles, event, start, end, cfg, label)
	}
	return newRangeDiscoveryDirectionalLabel(candles, event, start, end, cfg, label)
}

func newRangeDiscoveryBalanceLabel(candles []Candle, event rangeDiscoveryEvent, start, end int, cfg FuturesRangeCandidateDiscoveryAuditConfig, label rangeDiscoveryLabel) rangeDiscoveryLabel {
	high := event.episode.High
	low := event.episode.Low
	width := high - low
	upperQuartile := low + width*0.75
	lowerQuartile := low + width*0.25
	touchedUpper := false
	touchedLower := false
	firstAdverse := 0
	for i := start; i <= end; i++ {
		candle := candles[i]
		delay := i - event.eventIndex
		if candle.High >= upperQuartile {
			touchedUpper = true
		}
		if candle.Low <= lowerQuartile {
			touchedLower = true
		}
		if firstAdverse == 0 && (candle.Close > high || candle.Close < low) {
			firstAdverse = delay
			label.barsToAdverse = delay
		}
		if delay <= cfg.QuickInvalidationBars && (candle.Close > high || candle.Close < low) {
			label.quickInvalidation = true
		}
		label.favorableMovePct = math.Max(label.favorableMovePct, movePct(math.Max(0, math.Max(candle.High-high, low-candle.Low)), event.episode.EndClose))
		label.adverseMovePct = math.Max(label.adverseMovePct, movePct(math.Max(0, math.Max(candle.Close-high, low-candle.Close)), event.episode.EndClose))
	}
	label.favorableMinusAdversePct = label.favorableMovePct - label.adverseMovePct
	label.favorableGreaterThanAdverse = label.favorableMovePct > label.adverseMovePct
	if firstAdverse > 0 {
		label.outcomeClass = RangeDiscoveryClassAdverse
		label.outcomeLabel = event.adverseLabel
		return label
	}
	label.outcomeClass = RangeDiscoveryClassFavorable
	if touchedUpper && touchedLower {
		label.outcomeLabel = RangeDiscoveryOutcomeInternalRotation
		label.barsToFavorable = minIntNonZero(firstTouchDelay(candles, start, end, event.eventIndex, upperQuartile, true), firstTouchDelay(candles, start, end, event.eventIndex, lowerQuartile, false))
	} else {
		label.outcomeLabel = RangeDiscoveryOutcomeInsidePersistence
		label.barsToFavorable = 1
	}
	return label
}

func newRangeDiscoveryDirectionalLabel(candles []Candle, event rangeDiscoveryEvent, start, end int, cfg FuturesRangeCandidateDiscoveryAuditConfig, label rangeDiscoveryLabel) rangeDiscoveryLabel {
	width := event.episode.High - event.episode.Low
	minMove := width * cfg.MinMoveRangeFraction
	firstFavorable := 0
	firstAdverse := 0
	for i := start; i <= end; i++ {
		candle := candles[i]
		delay := i - event.eventIndex
		favorableMove := directionalFavorableMove(candle, event)
		adverseMove := directionalAdverseMove(candle, event)
		if favorableMove > 0 {
			label.favorableMovePct = math.Max(label.favorableMovePct, movePct(favorableMove, event.episode.EndClose))
		}
		if adverseMove > 0 {
			label.adverseMovePct = math.Max(label.adverseMovePct, movePct(adverseMove, event.episode.EndClose))
		}
		if firstFavorable == 0 && favorableMove >= minMove && favorableMove > adverseMove {
			firstFavorable = delay
			label.barsToFavorable = delay
		}
		if firstAdverse == 0 && rangeDiscoveryAdverseClose(candle, event) {
			firstAdverse = delay
			label.barsToAdverse = delay
		}
		if delay <= cfg.QuickInvalidationBars && rangeDiscoveryAdverseClose(candle, event) {
			label.quickInvalidation = true
		}
	}
	label.favorableMinusAdversePct = label.favorableMovePct - label.adverseMovePct
	label.favorableGreaterThanAdverse = label.favorableMovePct > label.adverseMovePct
	switch {
	case firstFavorable > 0 && (firstAdverse == 0 || firstFavorable < firstAdverse || (firstFavorable == firstAdverse && label.favorableGreaterThanAdverse)):
		label.outcomeClass = RangeDiscoveryClassFavorable
		label.outcomeLabel = event.favorableLabel
	case firstAdverse > 0:
		label.outcomeClass = RangeDiscoveryClassAdverse
		label.outcomeLabel = event.adverseLabel
	default:
		label.outcomeClass = RangeDiscoveryClassNeutral
		label.outcomeLabel = event.neutralLabel
	}
	return label
}

func directionalFavorableMove(candle Candle, event rangeDiscoveryEvent) float64 {
	switch event.expectedDirection {
	case "up":
		switch event.family {
		case RangeDiscoveryFamilyCleanBreakout:
			return candle.High - event.episode.High
		default:
			return candle.High - event.episode.EndClose
		}
	case "down":
		switch event.family {
		case RangeDiscoveryFamilyCleanBreakout:
			return event.episode.Low - candle.Low
		default:
			return event.episode.EndClose - candle.Low
		}
	default:
		return 0
	}
}

func directionalAdverseMove(candle Candle, event rangeDiscoveryEvent) float64 {
	switch event.expectedDirection {
	case "up":
		return event.episode.EndClose - candle.Low
	case "down":
		return candle.High - event.episode.EndClose
	default:
		return 0
	}
}

func rangeDiscoveryAdverseClose(candle Candle, event rangeDiscoveryEvent) bool {
	switch event.family {
	case RangeDiscoveryFamilyCleanBreakout:
		if event.expectedDirection == "up" {
			return candle.Close <= event.episode.High
		}
		if event.expectedDirection == "down" {
			return candle.Close >= event.episode.Low
		}
	case RangeDiscoveryFamilyBoundaryTouch, RangeDiscoveryFamilyWickRejection, RangeDiscoveryFamilyFailedBreakReentry:
		if event.side == RangeDiscoverySideSupport {
			return candle.Close < event.episode.Low
		}
		if event.side == RangeDiscoverySideResistance {
			return candle.Close > event.episode.High
		}
	}
	return false
}

func summarizeRangeDiscovery(rows []FuturesRangeDiscoveryCandidateRow, roundTripCostPct float64) []FuturesRangeDiscoverySummaryRow {
	accumulators := map[rangeDiscoverySummaryKey]*rangeDiscoverySummaryAccumulator{}
	for _, row := range rows {
		for _, split := range rangeDiscoverySplitCombos(row.Split) {
			for _, side := range []string{row.Side, RangeDiscoverySideAll} {
				key := rangeDiscoverySummaryKey{
					split:       split,
					timeframe:   row.Timeframe,
					family:      row.Family,
					side:        side,
					horizonBars: row.HorizonBars,
				}
				acc := accumulators[key]
				if acc == nil {
					acc = &rangeDiscoverySummaryAccumulator{key: key, roundTripCostPct: roundTripCostPct}
					accumulators[key] = acc
				}
				acc.add(row)
			}
		}
	}
	out := make([]FuturesRangeDiscoverySummaryRow, 0, len(accumulators))
	for _, acc := range accumulators {
		out = append(out, acc.row())
	}
	sort.Slice(out, func(i, j int) bool {
		return lessRangeDiscoverySummary(out[i], out[j])
	})
	return out
}

func (acc *rangeDiscoverySummaryAccumulator) add(row FuturesRangeDiscoveryCandidateRow) {
	acc.candidates++
	acc.rangeWidthSum += row.RangeWidthPct
	if row.MissingFuture {
		acc.missingFuture++
		return
	}
	acc.labeled++
	switch row.OutcomeClass {
	case RangeDiscoveryClassFavorable:
		acc.favorable++
	case RangeDiscoveryClassAdverse:
		acc.adverse++
	case RangeDiscoveryClassNeutral:
		acc.neutral++
	}
	if row.QuickInvalidation {
		acc.quickInvalidation++
	}
	acc.favorableMoveSum += row.FavorableMovePct
	acc.adverseMoveSum += row.AdverseMovePct
	acc.favorableMinusSum += row.FavorableMinusAdversePct
}

func (acc *rangeDiscoverySummaryAccumulator) row() FuturesRangeDiscoverySummaryRow {
	row := FuturesRangeDiscoverySummaryRow{
		Split:                  acc.key.split,
		Timeframe:              acc.key.timeframe,
		Family:                 acc.key.family,
		Side:                   acc.key.side,
		HorizonBars:            acc.key.horizonBars,
		CandidateCount:         acc.candidates,
		LabeledCount:           acc.labeled,
		MissingFutureCount:     acc.missingFuture,
		FavorableCount:         acc.favorable,
		AdverseCount:           acc.adverse,
		NeutralCount:           acc.neutral,
		QuickInvalidationCount: acc.quickInvalidation,
		RoundTripCostPct:       acc.roundTripCostPct,
	}
	if acc.candidates > 0 {
		row.AvgRangeWidthPct = acc.rangeWidthSum / float64(acc.candidates)
	}
	if acc.labeled > 0 {
		denom := float64(acc.labeled)
		row.FavorableRate = float64(acc.favorable) / denom
		row.AdverseRate = float64(acc.adverse) / denom
		row.NeutralRate = float64(acc.neutral) / denom
		row.QuickInvalidationRate = float64(acc.quickInvalidation) / denom
		row.AvgFavorableMovePct = acc.favorableMoveSum / denom
		row.AvgAdverseMovePct = acc.adverseMoveSum / denom
		row.AvgFavorableMinusAdversePct = acc.favorableMinusSum / denom
		row.CostBufferPct = row.AvgFavorableMinusAdversePct - row.RoundTripCostPct
		row.CostBufferPass = row.CostBufferPct > 0
	}
	return row
}

func rangeDiscoveryStabilityRows(summaryRows []FuturesRangeDiscoverySummaryRow, splits []Split) []FuturesRangeDiscoveryStabilityRow {
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	periodSplitSet := map[string]bool{}
	for _, split := range splits {
		if split.Name != fullSplitName {
			periodSplitSet[split.Name] = true
		}
	}
	byKey := map[rangeDiscoveryStabilityKey][]FuturesRangeDiscoverySummaryRow{}
	for _, row := range summaryRows {
		if !periodSplitSet[row.Split] {
			continue
		}
		key := rangeDiscoveryStabilityKey{
			timeframe:   row.Timeframe,
			family:      row.Family,
			side:        row.Side,
			horizonBars: row.HorizonBars,
		}
		byKey[key] = append(byKey[key], row)
	}
	out := make([]FuturesRangeDiscoveryStabilityRow, 0, len(byKey))
	for key, rows := range byKey {
		out = append(out, newRangeDiscoveryStabilityRow(key, rows))
	}
	sort.Slice(out, func(i, j int) bool {
		return lessRangeDiscoveryStability(out[i], out[j])
	})
	return out
}

func newRangeDiscoveryStabilityRow(key rangeDiscoveryStabilityKey, rows []FuturesRangeDiscoverySummaryRow) FuturesRangeDiscoveryStabilityRow {
	row := FuturesRangeDiscoveryStabilityRow{
		Timeframe:                   key.timeframe,
		Family:                      key.family,
		Side:                        key.side,
		HorizonBars:                 key.horizonBars,
		PeriodSplits:                len(rows),
		CandidateCountMin:           math.MaxInt,
		FavorableRateMin:            math.Inf(1),
		AdverseRateMin:              math.Inf(1),
		FavorableMinusAdversePctMin: math.Inf(1),
		CostBufferPctMin:            math.Inf(1),
		CostBufferPassAllSplits:     true,
	}
	for _, summary := range rows {
		row.CandidateCount += summary.CandidateCount
		row.CandidateCountMin = minInt(row.CandidateCountMin, summary.CandidateCount)
		row.CandidateCountMax = maxInt(row.CandidateCountMax, summary.CandidateCount)
		row.FavorableRateMin = math.Min(row.FavorableRateMin, summary.FavorableRate)
		row.FavorableRateMax = math.Max(row.FavorableRateMax, summary.FavorableRate)
		row.AdverseRateMin = math.Min(row.AdverseRateMin, summary.AdverseRate)
		row.AdverseRateMax = math.Max(row.AdverseRateMax, summary.AdverseRate)
		row.QuickInvalidationRateMax = math.Max(row.QuickInvalidationRateMax, summary.QuickInvalidationRate)
		row.FavorableMinusAdversePctMin = math.Min(row.FavorableMinusAdversePctMin, summary.AvgFavorableMinusAdversePct)
		row.FavorableMinusAdversePctMax = math.Max(row.FavorableMinusAdversePctMax, summary.AvgFavorableMinusAdversePct)
		row.CostBufferPctMin = math.Min(row.CostBufferPctMin, summary.CostBufferPct)
		if !summary.CostBufferPass {
			row.CostBufferPassAllSplits = false
		}
	}
	if row.CandidateCountMin == math.MaxInt {
		row.CandidateCountMin = 0
	}
	if math.IsInf(row.FavorableRateMin, 1) {
		row.FavorableRateMin = 0
	}
	if math.IsInf(row.AdverseRateMin, 1) {
		row.AdverseRateMin = 0
	}
	if math.IsInf(row.FavorableMinusAdversePctMin, 1) {
		row.FavorableMinusAdversePctMin = 0
	}
	if math.IsInf(row.CostBufferPctMin, 1) {
		row.CostBufferPctMin = 0
	}
	row.CandidateCountDelta = row.CandidateCountMax - row.CandidateCountMin
	row.FavorableRateDelta = row.FavorableRateMax - row.FavorableRateMin
	row.AdverseRateDelta = row.AdverseRateMax - row.AdverseRateMin
	row.FavorableMinusAdversePctDelta = row.FavorableMinusAdversePctMax - row.FavorableMinusAdversePctMin
	return row
}

func rangeDiscoveryRankingRows(summaryRows []FuturesRangeDiscoverySummaryRow, stabilityRows []FuturesRangeDiscoveryStabilityRow, cfg FuturesRangeCandidateDiscoveryAuditConfig, splits []Split) []FuturesRangeDiscoveryRankingRow {
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	periodSplits := rangeDiscoveryPeriodSplits(splits)
	summaryByKey := map[rangeDiscoverySummaryKey]FuturesRangeDiscoverySummaryRow{}
	for _, row := range summaryRows {
		summaryByKey[rangeDiscoverySummaryKey{split: row.Split, timeframe: row.Timeframe, family: row.Family, side: row.Side, horizonBars: row.HorizonBars}] = row
	}
	stabilityByKey := map[rangeDiscoveryStabilityKey]FuturesRangeDiscoveryStabilityRow{}
	for _, row := range stabilityRows {
		stabilityByKey[rangeDiscoveryStabilityKey{timeframe: row.Timeframe, family: row.Family, side: row.Side, horizonBars: row.HorizonBars}] = row
	}
	keys := []rangeDiscoveryStabilityKey{}
	for key := range stabilityByKey {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return lessRangeDiscoveryStabilityKey(keys[i], keys[j])
	})
	rows := make([]FuturesRangeDiscoveryRankingRow, 0, len(keys))
	for _, key := range keys {
		full := summaryByKey[rangeDiscoverySummaryKey{split: fullSplitName, timeframe: key.timeframe, family: key.family, side: key.side, horizonBars: key.horizonBars}]
		stability := stabilityByKey[key]
		ranking := FuturesRangeDiscoveryRankingRow{
			Timeframe:                       key.timeframe,
			Family:                          key.family,
			Side:                            key.side,
			HorizonBars:                     key.horizonBars,
			FullCandidateCount:              full.CandidateCount,
			WeakestSplitCandidateCount:      stability.CandidateCountMin,
			FullFavorableRate:               full.FavorableRate,
			WeakestSplitFavorableRate:       stability.FavorableRateMin,
			WorstSplitAdverseRate:           stability.AdverseRateMax,
			WorstSplitQuickInvalidationRate: stability.QuickInvalidationRateMax,
			FullFavorableMinusAdversePct:    full.AvgFavorableMinusAdversePct,
			WeakestFavorableMinusAdversePct: stability.FavorableMinusAdversePctMin,
			WeakestCostBufferPct:            stability.CostBufferPctMin,
			PeriodSplitsRequired:            len(periodSplits),
		}
		passes := true
		reasons := []string{}
		for _, split := range periodSplits {
			summary, ok := summaryByKey[rangeDiscoverySummaryKey{split: split.Name, timeframe: key.timeframe, family: key.family, side: key.side, horizonBars: key.horizonBars}]
			if !ok {
				passes = false
				reasons = append(reasons, "missing_split")
				continue
			}
			if summary.CandidateCount < cfg.MinCandidatesPerSplit {
				passes = false
				reasons = append(reasons, "weak_count")
				continue
			}
			if summary.FavorableRate <= summary.AdverseRate {
				passes = false
				reasons = append(reasons, "favorable_not_above_adverse")
				continue
			}
			if summary.QuickInvalidationRate >= summary.FavorableRate {
				passes = false
				reasons = append(reasons, "quick_invalidation_dominates")
				continue
			}
			if summary.AvgFavorableMovePct <= summary.AvgAdverseMovePct {
				passes = false
				reasons = append(reasons, "adverse_excursion_dominates")
				continue
			}
			if !summary.CostBufferPass {
				passes = false
				reasons = append(reasons, "cost_buffer_fail")
				continue
			}
			ranking.PeriodSplitsPassing++
		}
		ranking.PassesGate = passes && ranking.PeriodSplitsPassing == ranking.PeriodSplitsRequired
		ranking.RankScore = rangeDiscoveryRankScore(ranking)
		if !ranking.PassesGate && len(reasons) > 0 {
			ranking.FailureReason = uniqueJoinedReasons(reasons)
		}
		rows = append(rows, ranking)
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].PassesGate != rows[j].PassesGate {
			return rows[i].PassesGate
		}
		if rows[i].RankScore != rows[j].RankScore {
			return rows[i].RankScore > rows[j].RankScore
		}
		return lessRangeDiscoveryRanking(rows[i], rows[j])
	})
	for i := range rows {
		rows[i].Rank = i + 1
	}
	return rows
}

func rangeDiscoveryRankScore(row FuturesRangeDiscoveryRankingRow) float64 {
	score := row.WeakestSplitFavorableRate - row.WorstSplitAdverseRate
	score += row.WeakestFavorableMinusAdversePct * 20
	score -= row.WorstSplitQuickInvalidationRate * 0.5
	score += math.Log1p(float64(row.WeakestSplitCandidateCount)) / 20
	return score
}

func rangeDiscoverySplitCombos(split string) []string {
	if split == "" || split == fullSplitName {
		return []string{fullSplitName}
	}
	return []string{split, fullSplitName}
}

func rangeDiscoveryPeriodSplits(splits []Split) []Split {
	out := []Split{}
	for _, split := range splits {
		if split.Name != fullSplitName {
			out = append(out, split)
		}
	}
	return out
}

func firstTouchDelay(candles []Candle, start, end, eventIndex int, level float64, above bool) int {
	for i := start; i <= end; i++ {
		if above && candles[i].High >= level {
			return i - eventIndex
		}
		if !above && candles[i].Low <= level {
			return i - eventIndex
		}
	}
	return 0
}

func minIntNonZero(a, b int) int {
	if a == 0 {
		return b
	}
	if b == 0 {
		return a
	}
	return minInt(a, b)
}

func uniqueJoinedReasons(reasons []string) string {
	seen := map[string]bool{}
	out := ""
	for _, reason := range reasons {
		if seen[reason] {
			continue
		}
		seen[reason] = true
		if out != "" {
			out += ","
		}
		out += reason
	}
	return out
}

func lessRangeDiscoveryCandidate(a, b FuturesRangeDiscoveryCandidateRow) bool {
	if a.EventIndex != b.EventIndex {
		return a.EventIndex < b.EventIndex
	}
	if rangeDiscoveryTimeframeSortKey(a.Timeframe) != rangeDiscoveryTimeframeSortKey(b.Timeframe) {
		return rangeDiscoveryTimeframeSortKey(a.Timeframe) < rangeDiscoveryTimeframeSortKey(b.Timeframe)
	}
	if rangeDiscoveryFamilySortKey(a.Family) != rangeDiscoveryFamilySortKey(b.Family) {
		return rangeDiscoveryFamilySortKey(a.Family) < rangeDiscoveryFamilySortKey(b.Family)
	}
	if rangeDiscoverySideSortKey(a.Side) != rangeDiscoverySideSortKey(b.Side) {
		return rangeDiscoverySideSortKey(a.Side) < rangeDiscoverySideSortKey(b.Side)
	}
	return a.HorizonBars < b.HorizonBars
}

func lessRangeDiscoverySummary(a, b FuturesRangeDiscoverySummaryRow) bool {
	if splitSortKey(a.Split) != splitSortKey(b.Split) {
		return splitSortKey(a.Split) < splitSortKey(b.Split)
	}
	if rangeDiscoveryTimeframeSortKey(a.Timeframe) != rangeDiscoveryTimeframeSortKey(b.Timeframe) {
		return rangeDiscoveryTimeframeSortKey(a.Timeframe) < rangeDiscoveryTimeframeSortKey(b.Timeframe)
	}
	if rangeDiscoveryFamilySortKey(a.Family) != rangeDiscoveryFamilySortKey(b.Family) {
		return rangeDiscoveryFamilySortKey(a.Family) < rangeDiscoveryFamilySortKey(b.Family)
	}
	if rangeDiscoverySideSortKey(a.Side) != rangeDiscoverySideSortKey(b.Side) {
		return rangeDiscoverySideSortKey(a.Side) < rangeDiscoverySideSortKey(b.Side)
	}
	return a.HorizonBars < b.HorizonBars
}

func lessRangeDiscoveryStability(a, b FuturesRangeDiscoveryStabilityRow) bool {
	return lessRangeDiscoveryStabilityKey(
		rangeDiscoveryStabilityKey{timeframe: a.Timeframe, family: a.Family, side: a.Side, horizonBars: a.HorizonBars},
		rangeDiscoveryStabilityKey{timeframe: b.Timeframe, family: b.Family, side: b.Side, horizonBars: b.HorizonBars},
	)
}

func lessRangeDiscoveryStabilityKey(a, b rangeDiscoveryStabilityKey) bool {
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

func lessRangeDiscoveryRanking(a, b FuturesRangeDiscoveryRankingRow) bool {
	if rangeDiscoveryTimeframeSortKey(a.Timeframe) != rangeDiscoveryTimeframeSortKey(b.Timeframe) {
		return rangeDiscoveryTimeframeSortKey(a.Timeframe) < rangeDiscoveryTimeframeSortKey(b.Timeframe)
	}
	if rangeDiscoveryFamilySortKey(a.Family) != rangeDiscoveryFamilySortKey(b.Family) {
		return rangeDiscoveryFamilySortKey(a.Family) < rangeDiscoveryFamilySortKey(b.Family)
	}
	if rangeDiscoverySideSortKey(a.Side) != rangeDiscoverySideSortKey(b.Side) {
		return rangeDiscoverySideSortKey(a.Side) < rangeDiscoverySideSortKey(b.Side)
	}
	return a.HorizonBars < b.HorizonBars
}

func rangeDiscoveryTimeframeSortKey(timeframe string) int {
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

func rangeDiscoveryFamilySortKey(family string) int {
	switch family {
	case RangeDiscoveryFamilyBalancePersistence:
		return 0
	case RangeDiscoveryFamilyBoundaryTouch:
		return 1
	case RangeDiscoveryFamilyWickRejection:
		return 2
	case RangeDiscoveryFamilyFailedBreakReentry:
		return 3
	case RangeDiscoveryFamilyCleanBreakout:
		return 4
	default:
		return 99
	}
}

func rangeDiscoverySideSortKey(side string) int {
	switch side {
	case RangeDiscoverySideAll:
		return 0
	case RangeDiscoverySideBalance:
		return 1
	case RangeDiscoverySideSupport:
		return 2
	case RangeDiscoverySideResistance:
		return 3
	case RangeDiscoverySideUp:
		return 4
	case RangeDiscoverySideDown:
		return 5
	default:
		return 99
	}
}
