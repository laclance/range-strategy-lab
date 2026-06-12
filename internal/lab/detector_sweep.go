package lab

import "fmt"

type DetectorSweepProfile struct {
	ProfileID          string  `json:"profile_id"`
	IsBalancedBaseline bool    `json:"is_balanced_baseline"`
	IsADXComparison    bool    `json:"is_adx_comparison"`
	Percentile         float64 `json:"percentile"`
	MinConsecutiveBars int     `json:"min_consecutive_bars"`
	UseBollinger       bool    `json:"use_bollinger"`
	UseADX             bool    `json:"use_adx"`
	LookbackDays       int     `json:"lookback_days"`
}

type DetectorSweepRow struct {
	ProfileID            string  `json:"profile_id"`
	IsBalancedBaseline   bool    `json:"is_balanced_baseline"`
	IsADXComparison      bool    `json:"is_adx_comparison"`
	Percentile           float64 `json:"percentile"`
	MinConsecutiveBars   int     `json:"min_consecutive_bars"`
	UseBollinger         bool    `json:"use_bollinger"`
	UseADX               bool    `json:"use_adx"`
	LookbackDays         int     `json:"lookback_days"`
	Split                string  `json:"split"`
	ActiveBars           int     `json:"active_bars"`
	TotalBars            int     `json:"total_bars"`
	DutyCycle            float64 `json:"duty_cycle"`
	Episodes             int     `json:"episodes"`
	AvgEpisodeLength     float64 `json:"avg_episode_length"`
	MedianEpisodeLength  float64 `json:"median_episode_length"`
	LongestEpisodeLength int     `json:"longest_episode_length"`
}

func DefaultDetectorSweepProfiles(lookbackDays int) []DetectorSweepProfile {
	if lookbackDays <= 0 {
		lookbackDays = DefaultCompressionRangeDetectorConfig().LookbackDays
	}

	percentiles := []float64{0.20, 0.30, 0.40}
	minConsecutiveBars := []int{6, 12, 24}
	useBollingerOptions := []bool{true, false}

	var profiles []DetectorSweepProfile
	for _, percentile := range percentiles {
		for _, minConsecutive := range minConsecutiveBars {
			for _, useBollinger := range useBollingerOptions {
				profiles = append(profiles, newDetectorSweepProfile(
					percentile,
					minConsecutive,
					useBollinger,
					false,
					lookbackDays,
					false,
				))
			}
		}
	}

	profiles = append(profiles, newDetectorSweepProfile(0.30, 12, true, true, lookbackDays, true))
	return profiles
}

func RunDetectorSweep(candles []Candle, cfg RangeDetectorConfig, splits []Split) ([]DetectorSweepRow, error) {
	cfg = cfg.withDefaults()
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	if len(splits) == 0 {
		splits = DefaultSplits()
	}

	profiles := DefaultDetectorSweepProfiles(cfg.LookbackDays)
	lookbackBars := cfg.LookbackDays * cfg.BarsPerDay

	atr := NormalizedATR(candles, cfg.ATRPeriod)
	donchian := DonchianWidth(candles, cfg.DonchianPeriod)
	bollinger := BollingerWidth(candles, cfg.BollingerPeriod)
	adx := ADX(candles, cfg.ADXPeriod)

	percentiles := []float64{0.20, 0.30, 0.40}
	atrThresholds := make(map[float64][]float64, len(percentiles))
	donchianThresholds := make(map[float64][]float64, len(percentiles))
	bollingerThresholds := make(map[float64][]float64, len(percentiles))
	adxThresholds := make(map[float64][]float64, len(percentiles))
	for _, percentile := range percentiles {
		atrThresholds[percentile] = rollingPriorPercentile(atr, lookbackBars, percentile)
		donchianThresholds[percentile] = rollingPriorPercentile(donchian, lookbackBars, percentile)
		bollingerThresholds[percentile] = rollingPriorPercentile(bollinger, lookbackBars, percentile)
		adxThresholds[percentile] = rollingPriorPercentile(adx, lookbackBars, percentile)
	}

	rows := make([]DetectorSweepRow, 0, len(profiles)*len(splits))
	for _, profile := range profiles {
		classifications := classifyDetectorSweepProfile(
			candles,
			profile,
			atr,
			donchian,
			bollinger,
			adx,
			atrThresholds[profile.Percentile],
			donchianThresholds[profile.Percentile],
			bollingerThresholds[profile.Percentile],
			adxThresholds[profile.Percentile],
		)
		dutyRows, _ := SummarizeDetectorSplits(candles, classifications, splits)
		rows = append(rows, detectorSweepRows(profile, dutyRows)...)
	}
	return rows, nil
}

func newDetectorSweepProfile(percentile float64, minConsecutiveBars int, useBollinger, useADX bool, lookbackDays int, isADXComparison bool) DetectorSweepProfile {
	isBalancedBaseline := percentile == 0.30 && minConsecutiveBars == 12 && useBollinger && !useADX
	return DetectorSweepProfile{
		ProfileID:          detectorSweepProfileID(percentile, minConsecutiveBars, useBollinger, useADX),
		IsBalancedBaseline: isBalancedBaseline,
		IsADXComparison:    isADXComparison,
		Percentile:         percentile,
		MinConsecutiveBars: minConsecutiveBars,
		UseBollinger:       useBollinger,
		UseADX:             useADX,
		LookbackDays:       lookbackDays,
	}
}

func detectorSweepProfileID(percentile float64, minConsecutiveBars int, useBollinger, useADX bool) string {
	return fmt.Sprintf(
		"p%02d_c%02d_bollinger_%s_adx_%s",
		int(percentile*100+0.5),
		minConsecutiveBars,
		boolToken(useBollinger),
		boolToken(useADX),
	)
}

func boolToken(value bool) string {
	if value {
		return "on"
	}
	return "off"
}

func classifyDetectorSweepProfile(
	candles []Candle,
	profile DetectorSweepProfile,
	atr []float64,
	donchian []float64,
	bollinger []float64,
	adx []float64,
	atrThresholds []float64,
	donchianThresholds []float64,
	bollingerThresholds []float64,
	adxThresholds []float64,
) []RangeClassification {
	raw := make([]bool, len(candles))
	for i := range candles {
		raw[i] = metricLow(atr[i], atrThresholds[i]) &&
			metricLow(donchian[i], donchianThresholds[i])
		if raw[i] && profile.UseBollinger {
			raw[i] = metricLow(bollinger[i], bollingerThresholds[i])
		}
		if raw[i] && profile.UseADX {
			raw[i] = metricLow(adx[i], adxThresholds[i])
		}
	}
	active := markConsecutiveActive(raw, profile.MinConsecutiveBars)

	out := make([]RangeClassification, len(candles))
	for i, c := range candles {
		out[i] = RangeClassification{
			Index:     i,
			CloseTime: c.CloseTime.Format(timeLayout),
			RawActive: raw[i],
			Active:    active[i],
		}
	}
	return out
}

func detectorSweepRows(profile DetectorSweepProfile, dutyRows []DetectorDutyCycleRow) []DetectorSweepRow {
	rows := make([]DetectorSweepRow, 0, len(dutyRows))
	for _, dutyRow := range dutyRows {
		rows = append(rows, DetectorSweepRow{
			ProfileID:            profile.ProfileID,
			IsBalancedBaseline:   profile.IsBalancedBaseline,
			IsADXComparison:      profile.IsADXComparison,
			Percentile:           profile.Percentile,
			MinConsecutiveBars:   profile.MinConsecutiveBars,
			UseBollinger:         profile.UseBollinger,
			UseADX:               profile.UseADX,
			LookbackDays:         profile.LookbackDays,
			Split:                dutyRow.Split,
			ActiveBars:           dutyRow.ActiveBars,
			TotalBars:            dutyRow.TotalBars,
			DutyCycle:            dutyRow.DutyCycle,
			Episodes:             dutyRow.Episodes,
			AvgEpisodeLength:     dutyRow.AvgEpisodeLength,
			MedianEpisodeLength:  dutyRow.MedianEpisodeLength,
			LongestEpisodeLength: dutyRow.LongestEpisodeLength,
		})
	}
	return rows
}
