package lab

import "sort"

type DetectorDutyCycleRow struct {
	Split                string  `json:"split"`
	ActiveBars           int     `json:"active_bars"`
	TotalBars            int     `json:"total_bars"`
	DutyCycle            float64 `json:"duty_cycle"`
	Episodes             int     `json:"episodes"`
	AvgEpisodeLength     float64 `json:"avg_episode_length"`
	MedianEpisodeLength  float64 `json:"median_episode_length"`
	LongestEpisodeLength int     `json:"longest_episode_length"`
}

type RangeEpisode struct {
	Split      string `json:"split"`
	StartIndex int    `json:"start_index"`
	EndIndex   int    `json:"end_index"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	LengthBars int    `json:"length_bars"`
}

func SummarizeDetectorSplits(candles []Candle, classifications []RangeClassification, splits []Split) ([]DetectorDutyCycleRow, []RangeEpisode) {
	if len(classifications) > len(candles) {
		classifications = classifications[:len(candles)]
	}
	rows := make([]DetectorDutyCycleRow, 0, len(splits))
	var episodes []RangeEpisode
	for _, split := range splits {
		row, splitEpisodes := summarizeDetectorSplit(candles, classifications, split)
		rows = append(rows, row)
		episodes = append(episodes, splitEpisodes...)
	}
	return rows, episodes
}

func summarizeDetectorSplit(candles []Candle, classifications []RangeClassification, split Split) (DetectorDutyCycleRow, []RangeEpisode) {
	row := DetectorDutyCycleRow{Split: split.Name}
	var episodes []RangeEpisode
	var episodeLengths []int

	inEpisode := false
	startIndex := 0
	endIndex := 0
	length := 0

	closeEpisode := func() {
		if !inEpisode {
			return
		}
		episode := RangeEpisode{
			Split:      split.Name,
			StartIndex: startIndex,
			EndIndex:   endIndex,
			StartTime:  candles[startIndex].CloseTime.Format(timeLayout),
			EndTime:    candles[endIndex].CloseTime.Format(timeLayout),
			LengthBars: length,
		}
		episodes = append(episodes, episode)
		episodeLengths = append(episodeLengths, length)
		inEpisode = false
		length = 0
	}

	for i, c := range candles {
		if !split.Contains(c.CloseTime) {
			closeEpisode()
			continue
		}
		row.TotalBars++
		active := i < len(classifications) && classifications[i].Active
		if !active {
			closeEpisode()
			continue
		}
		row.ActiveBars++
		if !inEpisode {
			inEpisode = true
			startIndex = i
			length = 0
		}
		endIndex = i
		length++
	}
	closeEpisode()

	row.Episodes = len(episodeLengths)
	if row.TotalBars > 0 {
		row.DutyCycle = float64(row.ActiveBars) / float64(row.TotalBars)
	}
	if len(episodeLengths) > 0 {
		sum := 0
		for _, episodeLength := range episodeLengths {
			sum += episodeLength
			if episodeLength > row.LongestEpisodeLength {
				row.LongestEpisodeLength = episodeLength
			}
		}
		row.AvgEpisodeLength = float64(sum) / float64(len(episodeLengths))
		row.MedianEpisodeLength = medianInt(episodeLengths)
	}
	return row, episodes
}

func medianInt(values []int) float64 {
	if len(values) == 0 {
		return 0
	}
	sortedValues := append([]int(nil), values...)
	sort.Ints(sortedValues)
	mid := len(sortedValues) / 2
	if len(sortedValues)%2 == 1 {
		return float64(sortedValues[mid])
	}
	return float64(sortedValues[mid-1]+sortedValues[mid]) / 2
}
