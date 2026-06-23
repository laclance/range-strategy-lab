package lab

import (
	"fmt"
	"math"
	"sort"
)

const HoldInsideMidlineTouchPrototypeName = "hold_inside_midline_touch_prototype"

type HoldInsideMidlineTouchPrototypeConfig struct {
	MaxMidlineEventDelayBars int
	MaxHoldBars              int
	Profile                  DetectorSweepProfile
	ContextRule              DetectorContextRefinementRule
}

type HoldInsideMidlineTouchPrototypeStrategy struct {
	signalsByEventIndex map[int]Signal
	signals             []HoldInsideMidlineTouchPrototypeSignalRow
	signalsByID         map[string]HoldInsideMidlineTouchPrototypeSignalRow
}

type HoldInsideMidlineTouchPrototypeSignalRow struct {
	SignalID                 string    `json:"signal_id"`
	SkippedReason            string    `json:"skipped_reason"`
	ProfileID                string    `json:"profile_id"`
	ContextRule              string    `json:"context_rule"`
	HoldBars                 int       `json:"hold_bars"`
	EventType                string    `json:"event_type"`
	Split                    string    `json:"split"`
	SourceEpisodeID          int       `json:"source_episode_id"`
	EpisodeStartIndex        int       `json:"episode_start_index"`
	EpisodeEndIndex          int       `json:"episode_end_index"`
	EpisodeStartTime         string    `json:"episode_start_time"`
	EpisodeEndTime           string    `json:"episode_end_time"`
	EpisodeHigh              float64   `json:"episode_high"`
	EpisodeLow               float64   `json:"episode_low"`
	EpisodeMid               float64   `json:"episode_mid"`
	HoldDecisionIndex        int       `json:"hold_decision_index"`
	HoldDecisionTime         string    `json:"hold_decision_time"`
	HoldDecisionClose        float64   `json:"hold_decision_close"`
	HoldDecisionCloseBucket  string    `json:"hold_decision_close_position_bucket"`
	EventIndex               int       `json:"event_index"`
	EventTime                string    `json:"event_time"`
	EventDelayBars           int       `json:"event_delay_bars"`
	EventClose               float64   `json:"event_close"`
	EventClosePosition       float64   `json:"event_close_position"`
	EventClosePositionBucket string    `json:"event_close_position_bucket"`
	EventMidSide             string    `json:"event_mid_side"`
	Side                     Direction `json:"side"`
	Stop                     float64   `json:"stop"`
	Target                   float64   `json:"target"`
	MaxHoldBars              int       `json:"max_hold_bars"`
	EntryIndex               int       `json:"entry_index"`
	EntryTime                string    `json:"entry_time"`
	EntryOpen                float64   `json:"entry_open"`
	EntryGeometryValid       bool      `json:"entry_geometry_valid"`
}

type HoldInsideMidlineTouchPrototypeTradeRow struct {
	SignalID                 string    `json:"signal_id"`
	ProfileID                string    `json:"profile_id"`
	ContextRule              string    `json:"context_rule"`
	EventType                string    `json:"event_type"`
	SourceEpisodeID          int       `json:"source_episode_id"`
	EpisodeHigh              float64   `json:"episode_high"`
	EpisodeLow               float64   `json:"episode_low"`
	EpisodeMid               float64   `json:"episode_mid"`
	HoldDecisionIndex        int       `json:"hold_decision_index"`
	HoldDecisionTime         string    `json:"hold_decision_time"`
	HoldDecisionClose        float64   `json:"hold_decision_close"`
	EventIndex               int       `json:"event_index"`
	EventTime                string    `json:"event_time"`
	EventDelayBars           int       `json:"event_delay_bars"`
	EventClose               float64   `json:"event_close"`
	EventClosePositionBucket string    `json:"event_close_position_bucket"`
	EventMidSide             string    `json:"event_mid_side"`
	EntrySplit               string    `json:"entry_split"`
	CloseSplit               string    `json:"close_split"`
	Side                     Direction `json:"side"`
	EntryTime                string    `json:"entry_time"`
	ExitTime                 string    `json:"exit_time"`
	OpenIndex                int       `json:"open_index"`
	CloseIndex               int       `json:"close_index"`
	EntryPrice               float64   `json:"entry_price"`
	ExitPrice                float64   `json:"exit_price"`
	Stop                     float64   `json:"stop"`
	Target                   float64   `json:"target"`
	Size                     float64   `json:"size"`
	InitialRisk              float64   `json:"initial_risk"`
	GrossPnL                 float64   `json:"gross_pnl"`
	NetPnL                   float64   `json:"net_pnl"`
	Fees                     float64   `json:"fees"`
	Slippage                 float64   `json:"slippage"`
	GrossR                   float64   `json:"gross_r"`
	NetR                     float64   `json:"net_r"`
	ExitReason               string    `json:"exit_reason"`
	HoldBars                 int       `json:"hold_bars"`
}

type HoldInsideMidlineTouchPrototypeSummaryRow struct {
	Split              string  `json:"split"`
	Side               string  `json:"side"`
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
	IsWorstPeriodSplit bool    `json:"is_worst_period_split"`
}

func NewHoldInsideMidlineTouchPrototypeStrategy(candles []Candle, detectorCfg RangeDetectorConfig, cfg HoldInsideMidlineTouchPrototypeConfig, splits []Split) (HoldInsideMidlineTouchPrototypeStrategy, error) {
	detectorCfg = detectorCfg.withDefaults()
	if err := detectorCfg.validate(); err != nil {
		return HoldInsideMidlineTouchPrototypeStrategy{}, err
	}
	cfg = cfg.withDefaults(detectorCfg.LookbackDays)
	if err := cfg.validate(); err != nil {
		return HoldInsideMidlineTouchPrototypeStrategy{}, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}

	lookbackBars := detectorCfg.LookbackDays * detectorCfg.BarsPerDay
	atr := NormalizedATR(candles, detectorCfg.ATRPeriod)
	donchian := DonchianWidth(candles, detectorCfg.DonchianPeriod)
	bollinger := BollingerWidth(candles, detectorCfg.BollingerPeriod)
	adx := ADX(candles, detectorCfg.ADXPeriod)
	normalizedATR := NormalizedATR(candles, detectorCfg.ATRPeriod)

	atrThresholds := rollingPriorPercentile(atr, lookbackBars, cfg.Profile.Percentile)
	donchianThresholds := rollingPriorPercentile(donchian, lookbackBars, cfg.Profile.Percentile)
	bollingerThresholds := rollingPriorPercentile(bollinger, lookbackBars, cfg.Profile.Percentile)
	adxThresholds := rollingPriorPercentile(adx, lookbackBars, cfg.Profile.Percentile)
	classifications := classifyDetectorSweepProfile(
		candles,
		cfg.Profile,
		atr,
		donchian,
		bollinger,
		adx,
		atrThresholds,
		donchianThresholds,
		bollingerThresholds,
		adxThresholds,
	)

	return newHoldInsideMidlineTouchPrototypeStrategyFromClassifications(candles, classifications, normalizedATR, cfg, splits)
}

func newHoldInsideMidlineTouchPrototypeStrategyFromClassifications(candles []Candle, classifications []RangeClassification, normalizedATR []float64, cfg HoldInsideMidlineTouchPrototypeConfig, splits []Split) (HoldInsideMidlineTouchPrototypeStrategy, error) {
	cfg = cfg.withDefaults(cfg.Profile.LookbackDays)
	if err := cfg.validate(); err != nil {
		return HoldInsideMidlineTouchPrototypeStrategy{}, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	if len(classifications) > len(candles) {
		classifications = classifications[:len(candles)]
	}

	strategy := HoldInsideMidlineTouchPrototypeStrategy{
		signalsByEventIndex: map[int]Signal{},
		signalsByID:         map[string]HoldInsideMidlineTouchPrototypeSignalRow{},
	}
	episodes := rangeRegimeDurabilityEpisodes(candles, classifications, normalizedATR, splits, cfg.Profile.ProfileID)
	for _, episode := range episodes {
		row, ok := newHoldInsideMidlineTouchPrototypeSignalRow(candles, episode, cfg, splits, len(strategy.signals)+1)
		if !ok {
			continue
		}
		strategy.signals = append(strategy.signals, row)
		if row.SignalID == "" {
			continue
		}
		strategy.signalsByEventIndex[row.EventIndex] = Signal{
			Side:        row.Side,
			Stop:        row.Stop,
			Target:      row.Target,
			MaxHoldBars: row.MaxHoldBars,
			Reason:      row.SignalID,
		}
		strategy.signalsByID[row.SignalID] = row
	}
	sort.Slice(strategy.signals, func(i, j int) bool {
		return strategy.signals[i].EventIndex < strategy.signals[j].EventIndex
	})
	return strategy, nil
}

func (HoldInsideMidlineTouchPrototypeStrategy) Name() string {
	return HoldInsideMidlineTouchPrototypeName
}

func (s HoldInsideMidlineTouchPrototypeStrategy) OnCandle(ctx StrategyContext) (Signal, bool) {
	sig, ok := s.signalsByEventIndex[ctx.Index]
	return sig, ok
}

func (s HoldInsideMidlineTouchPrototypeStrategy) SignalRows() []HoldInsideMidlineTouchPrototypeSignalRow {
	return append([]HoldInsideMidlineTouchPrototypeSignalRow(nil), s.signals...)
}

func (s HoldInsideMidlineTouchPrototypeStrategy) TradeRows(trades []Trade, splits []Split) []HoldInsideMidlineTouchPrototypeTradeRow {
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	rows := make([]HoldInsideMidlineTouchPrototypeTradeRow, 0, len(trades))
	for _, trade := range trades {
		signal, ok := s.signalsByID[trade.Signal]
		if !ok {
			continue
		}
		entryTime, _ := parseTime(trade.EntryTime)
		exitTime, _ := parseTime(trade.ExitTime)
		initialRisk := math.Abs(trade.EntryPrice-trade.Stop) * trade.Size
		row := HoldInsideMidlineTouchPrototypeTradeRow{
			SignalID:                 signal.SignalID,
			ProfileID:                signal.ProfileID,
			ContextRule:              signal.ContextRule,
			EventType:                signal.EventType,
			SourceEpisodeID:          signal.SourceEpisodeID,
			EpisodeHigh:              signal.EpisodeHigh,
			EpisodeLow:               signal.EpisodeLow,
			EpisodeMid:               signal.EpisodeMid,
			HoldDecisionIndex:        signal.HoldDecisionIndex,
			HoldDecisionTime:         signal.HoldDecisionTime,
			HoldDecisionClose:        signal.HoldDecisionClose,
			EventIndex:               signal.EventIndex,
			EventTime:                signal.EventTime,
			EventDelayBars:           signal.EventDelayBars,
			EventClose:               signal.EventClose,
			EventClosePositionBucket: signal.EventClosePositionBucket,
			EventMidSide:             signal.EventMidSide,
			EntrySplit:               splitNameForCloseTime(entryTime, splits),
			CloseSplit:               splitNameForCloseTime(exitTime, splits),
			Side:                     trade.Side,
			EntryTime:                trade.EntryTime,
			ExitTime:                 trade.ExitTime,
			OpenIndex:                trade.OpenIndex,
			CloseIndex:               trade.CloseIndex,
			EntryPrice:               trade.EntryPrice,
			ExitPrice:                trade.ExitPrice,
			Stop:                     trade.Stop,
			Target:                   trade.Target,
			Size:                     trade.Size,
			InitialRisk:              initialRisk,
			GrossPnL:                 trade.GrossPnL,
			NetPnL:                   trade.NetPnL,
			Fees:                     trade.Fees,
			Slippage:                 trade.Slippage,
			ExitReason:               trade.Reason,
			HoldBars:                 trade.HoldBars,
		}
		if initialRisk > 0 {
			row.GrossR = trade.GrossPnL / initialRisk
			row.NetR = trade.NetPnL / initialRisk
		}
		rows = append(rows, row)
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].OpenIndex < rows[j].OpenIndex
	})
	return rows
}

func SummarizeHoldInsideMidlineTouchPrototype(trades []HoldInsideMidlineTouchPrototypeTradeRow, startBalance float64, splits []Split) []HoldInsideMidlineTouchPrototypeSummaryRow {
	if len(splits) == 0 {
		splits = DefaultSplits()
	}
	rows := make([]HoldInsideMidlineTouchPrototypeSummaryRow, 0, len(splits)*3)
	for _, split := range splits {
		for _, side := range []string{"all", "long", "short"} {
			filtered := filterHoldInsideMidlineTouchPrototypeTrades(trades, split, side)
			row := summarizeHoldInsideMidlineTouchPrototypeTrades(filtered, startBalance)
			row.Split = split.Name
			row.Side = side
			rows = append(rows, row)
		}
	}
	worstIndex := -1
	for i, row := range rows {
		if row.Split == fullSplitName || row.Side != "all" {
			continue
		}
		if worstIndex == -1 || row.NetPnL < rows[worstIndex].NetPnL {
			worstIndex = i
		}
	}
	if worstIndex >= 0 {
		rows[worstIndex].IsWorstPeriodSplit = true
	}
	return rows
}

func (cfg HoldInsideMidlineTouchPrototypeConfig) withDefaults(lookbackDays int) HoldInsideMidlineTouchPrototypeConfig {
	if cfg.MaxMidlineEventDelayBars == 0 {
		cfg.MaxMidlineEventDelayBars = 12
	}
	if cfg.MaxHoldBars == 0 {
		cfg.MaxHoldBars = 6
	}
	if cfg.Profile == (DetectorSweepProfile{}) {
		cfg.Profile = defaultHoldInsideDirectionalEdgeProfiles(lookbackDays)[0]
	}
	if cfg.ContextRule == (DetectorContextRefinementRule{}) {
		cfg.ContextRule = DetectorContextRefinementRule{RuleID: DetectorContextRuleHold3Inside, HoldBars: 3}
	}
	return cfg
}

func (cfg HoldInsideMidlineTouchPrototypeConfig) validate() error {
	if cfg.MaxMidlineEventDelayBars <= 0 {
		return fmt.Errorf("hold-inside midline touch prototype max event delay bars must be positive")
	}
	if cfg.MaxHoldBars <= 0 {
		return fmt.Errorf("hold-inside midline touch prototype max hold bars must be positive")
	}
	if cfg.Profile.ProfileID == "" {
		return fmt.Errorf("hold-inside midline touch prototype profile id must not be empty")
	}
	if cfg.ContextRule.RuleID != DetectorContextRuleHold3Inside || cfg.ContextRule.HoldBars != 3 || cfg.ContextRule.RequireMid50 {
		return fmt.Errorf("hold-inside midline touch prototype only supports hold_3_inside without hold-decision mid_50 filtering")
	}
	return nil
}

func newHoldInsideMidlineTouchPrototypeSignalRow(candles []Candle, episode rangeRegimeDurabilityEpisode, cfg HoldInsideMidlineTouchPrototypeConfig, splits []Split, sequence int) (HoldInsideMidlineTouchPrototypeSignalRow, bool) {
	holdDecisionIndex := episode.EndIndex + cfg.ContextRule.HoldBars
	if holdDecisionIndex < 0 || holdDecisionIndex >= len(candles) {
		return HoldInsideMidlineTouchPrototypeSignalRow{}, false
	}
	if !detectorContextRulePasses(candles, episode, cfg.ContextRule, holdDecisionIndex) {
		return HoldInsideMidlineTouchPrototypeSignalRow{}, false
	}
	episodeMid := (episode.High + episode.Low) / 2
	eventIndex, found := firstHoldInsideMidlineReactionEventIndex(candles, holdDecisionIndex, episodeMid, HoldInsideMidlineReactionEventTouch, cfg.MaxMidlineEventDelayBars)
	if !found {
		return HoldInsideMidlineTouchPrototypeSignalRow{}, false
	}
	event := candles[eventIndex]
	eventPosition := decisionClosePosition(event.Close, episode.Low, episode.High)
	if decisionClosePositionBucket(eventPosition) != decisionClosePositionBucketMid50 {
		return HoldInsideMidlineTouchPrototypeSignalRow{}, false
	}
	holdDecision := candles[holdDecisionIndex]
	holdPosition := decisionClosePosition(holdDecision.Close, episode.Low, episode.High)
	entryIndex := eventIndex + 1
	entryTime := ""
	entryOpen := 0.0
	if entryIndex < len(candles) {
		entryTime = candles[entryIndex].OpenTime.Format(timeLayout)
		entryOpen = candles[entryIndex].Open
	}

	row := HoldInsideMidlineTouchPrototypeSignalRow{
		ProfileID:                cfg.Profile.ProfileID,
		ContextRule:              cfg.ContextRule.RuleID,
		HoldBars:                 cfg.ContextRule.HoldBars,
		EventType:                HoldInsideMidlineReactionEventTouch,
		Split:                    splitNameForCloseTime(event.CloseTime, splits),
		SourceEpisodeID:          episode.EpisodeID,
		EpisodeStartIndex:        episode.StartIndex,
		EpisodeEndIndex:          episode.EndIndex,
		EpisodeStartTime:         candles[episode.StartIndex].CloseTime.Format(timeLayout),
		EpisodeEndTime:           candles[episode.EndIndex].CloseTime.Format(timeLayout),
		EpisodeHigh:              episode.High,
		EpisodeLow:               episode.Low,
		EpisodeMid:               episodeMid,
		HoldDecisionIndex:        holdDecisionIndex,
		HoldDecisionTime:         holdDecision.CloseTime.Format(timeLayout),
		HoldDecisionClose:        holdDecision.Close,
		HoldDecisionCloseBucket:  decisionClosePositionBucket(holdPosition),
		EventIndex:               eventIndex,
		EventTime:                event.CloseTime.Format(timeLayout),
		EventDelayBars:           eventIndex - holdDecisionIndex,
		EventClose:               event.Close,
		EventClosePosition:       normalizedPositionValue(eventPosition),
		EventClosePositionBucket: decisionClosePositionBucket(eventPosition),
		EventMidSide:             holdInsideDecisionMidSide(event.Close, episodeMid),
		MaxHoldBars:              cfg.MaxHoldBars,
		EntryIndex:               entryIndex,
		EntryTime:                entryTime,
		EntryOpen:                entryOpen,
	}

	switch row.EventMidSide {
	case holdInsideDecisionMidSideBelow:
		row.Side = Long
		row.Stop = episode.Low
		row.Target = episode.High
	case holdInsideDecisionMidSideAbove:
		row.Side = Short
		row.Stop = episode.High
		row.Target = episode.Low
	default:
		row.SkippedReason = "event_close_at_mid"
		return row, true
	}
	if entryIndex >= len(candles) {
		row.SkippedReason = "missing_entry_candle"
		return row, true
	}
	row.EntryGeometryValid = validEntryGeometry(Signal{Side: row.Side, Stop: row.Stop, Target: row.Target}, entryOpen)
	row.SignalID = fmt.Sprintf("hold_inside_midline_touch_%06d", sequence)
	return row, true
}

func normalizedPositionValue(position float64) float64 {
	if !validNumber(position) {
		return 0
	}
	return position
}

func filterHoldInsideMidlineTouchPrototypeTrades(trades []HoldInsideMidlineTouchPrototypeTradeRow, split Split, side string) []HoldInsideMidlineTouchPrototypeTradeRow {
	filtered := make([]HoldInsideMidlineTouchPrototypeTradeRow, 0, len(trades))
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
	return filtered
}

func summarizeHoldInsideMidlineTouchPrototypeTrades(trades []HoldInsideMidlineTouchPrototypeTradeRow, startBalance float64) HoldInsideMidlineTouchPrototypeSummaryRow {
	row := HoldInsideMidlineTouchPrototypeSummaryRow{TotalTrades: len(trades)}
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

func HoldInsideMidlineTouchPrototypeStopState(summary []HoldInsideMidlineTouchPrototypeSummaryRow) string {
	byKey := map[string]HoldInsideMidlineTouchPrototypeSummaryRow{}
	for _, row := range summary {
		byKey[row.Split+"|"+row.Side] = row
	}
	full := byKey[fullSplitName+"|all"]
	if full.TotalTrades == 0 || full.GrossPnL <= 0 || full.NetPnL <= 0 || full.ProfitFactor < 1.2 {
		return "prototype_failed_no_promotion"
	}
	for _, split := range []string{"2023_2024_oos", "2025_2026_recent"} {
		row := byKey[split+"|all"]
		if row.TotalTrades == 0 || row.NetPnL <= 0 {
			return "prototype_failed_no_promotion"
		}
	}
	for _, side := range []string{"long", "short"} {
		row := byKey[fullSplitName+"|"+side]
		if row.TotalTrades == 0 || row.NetPnL <= 0 {
			return "prototype_failed_no_promotion"
		}
	}
	return "prototype_review_passed_needs_stricter_oos_review"
}
