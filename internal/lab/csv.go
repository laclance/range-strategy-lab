package lab

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func LoadCSV(path string) ([]Candle, error) {
	candles, _, err := LoadCSVWithHeader(path)
	return candles, err
}

func LoadCSVWithHeader(path string) ([]Candle, []string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	header, err := r.Read()
	if err != nil {
		return nil, nil, err
	}
	cols := columnMap(header)
	required := []string{"open_time", "open", "high", "low", "close", "volume"}
	for _, name := range required {
		if _, ok := cols[name]; !ok {
			return nil, header, fmt.Errorf("missing required column %q", name)
		}
	}

	var candles []Candle
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, header, err
		}
		c, err := parseCandle(rec, cols)
		if err != nil {
			return nil, header, err
		}
		candles = append(candles, c)
	}
	return candles, header, nil
}

func columnMap(header []string) map[string]int {
	out := make(map[string]int, len(header))
	for i, raw := range header {
		out[strings.ToLower(strings.TrimSpace(raw))] = i
	}
	return out
}

func parseCandle(rec []string, cols map[string]int) (Candle, error) {
	openTime, err := parseTime(rec[cols["open_time"]])
	if err != nil {
		return Candle{}, fmt.Errorf("open_time: %w", err)
	}
	closeTime := openTime.Add(5*time.Minute - time.Millisecond)
	if idx, ok := cols["close_time"]; ok && idx < len(rec) && strings.TrimSpace(rec[idx]) != "" {
		closeTime, err = parseTime(rec[idx])
		if err != nil {
			return Candle{}, fmt.Errorf("close_time: %w", err)
		}
	}
	open, err := parseFloat(rec[cols["open"]])
	if err != nil {
		return Candle{}, fmt.Errorf("open: %w", err)
	}
	high, err := parseFloat(rec[cols["high"]])
	if err != nil {
		return Candle{}, fmt.Errorf("high: %w", err)
	}
	low, err := parseFloat(rec[cols["low"]])
	if err != nil {
		return Candle{}, fmt.Errorf("low: %w", err)
	}
	closePrice, err := parseFloat(rec[cols["close"]])
	if err != nil {
		return Candle{}, fmt.Errorf("close: %w", err)
	}
	volume, err := parseFloat(rec[cols["volume"]])
	if err != nil {
		return Candle{}, fmt.Errorf("volume: %w", err)
	}
	return Candle{
		OpenTime:  openTime.UTC(),
		CloseTime: closeTime.UTC(),
		Open:      open,
		High:      high,
		Low:       low,
		Close:     closePrice,
		Volume:    volume,
	}, nil
}

func parseTime(raw string) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, fmt.Errorf("empty timestamp")
	}
	if n, err := strconv.ParseInt(raw, 10, 64); err == nil {
		if n > 1_000_000_000_000 {
			return time.UnixMilli(n).UTC(), nil
		}
		return time.Unix(n, 0).UTC(), nil
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return time.Time{}, err
	}
	return t.UTC(), nil
}

func parseFloat(raw string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(raw), 64)
}
