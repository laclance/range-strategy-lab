package lab

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadCSVParsesBinanceShape(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "candles.csv")
	data := "open_time,open,high,low,close,volume,close_time,ignore\n" +
		"1609459200000,100,110,90,105,12.5,1609459499999,0\n"
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
	candles, err := LoadCSV(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(candles) != 1 {
		t.Fatalf("len=%d, want 1", len(candles))
	}
	wantOpen := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	if !candles[0].OpenTime.Equal(wantOpen) || candles[0].Close != 105 {
		t.Fatalf("unexpected candle: %+v", candles[0])
	}
}

func TestLoadCSVParsesNormalizedShapeWithDefaultCloseTime(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "candles.csv")
	data := " Open_Time , Open , HIGH , low , close , volume \n" +
		"2021-01-01T00:00:00Z,100,110,90,105,12.5\n"
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	candles, err := LoadCSV(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(candles) != 1 {
		t.Fatalf("len=%d, want 1", len(candles))
	}
	wantOpen := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	wantClose := wantOpen.Add(5*time.Minute - time.Millisecond)
	if !candles[0].OpenTime.Equal(wantOpen) || !candles[0].CloseTime.Equal(wantClose) {
		t.Fatalf("unexpected timestamps: %+v", candles[0])
	}
}

func TestLoadCSVRejectsMissingRequiredColumnAndEmptyFile(t *testing.T) {
	dir := t.TempDir()
	missingPath := filepath.Join(dir, "missing.csv")
	if err := os.WriteFile(missingPath, []byte("open_time,open,high,low,close\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadCSV(missingPath); err == nil {
		t.Fatalf("expected missing volume error")
	}

	emptyPath := filepath.Join(dir, "empty.csv")
	if err := os.WriteFile(emptyPath, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadCSV(emptyPath); err == nil {
		t.Fatalf("expected empty CSV error")
	}
}

func TestLoadCSVPropagatesOpenReadAndCandleErrors(t *testing.T) {
	if _, err := LoadCSV(filepath.Join(t.TempDir(), "missing.csv")); err == nil {
		t.Fatalf("expected open error")
	}

	dir := t.TempDir()
	badReadPath := filepath.Join(dir, "bad-read.csv")
	badRead := "open_time,open,high,low,close,volume\n" +
		"2021-01-01T00:00:00Z,\"100,110,90,105,1\n"
	if err := os.WriteFile(badReadPath, []byte(badRead), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadCSV(badReadPath); err == nil {
		t.Fatalf("expected CSV read error")
	}

	badCandlePath := filepath.Join(dir, "bad-candle.csv")
	badCandle := "open_time,open,high,low,close,volume\n" +
		"2021-01-01T00:00:00Z,bad,110,90,105,1\n"
	if err := os.WriteFile(badCandlePath, []byte(badCandle), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadCSV(badCandlePath); err == nil {
		t.Fatalf("expected candle parse error")
	}
}

func TestParseCandleRejectsBadFields(t *testing.T) {
	cols := columnMap([]string{"open_time", "open", "high", "low", "close", "volume", "close_time"})
	good := []string{"1609459200", "100", "110", "90", "105", "12.5", "1609459499999"}
	cases := []struct {
		name string
		idx  int
	}{
		{name: "open_time", idx: cols["open_time"]},
		{name: "open", idx: cols["open"]},
		{name: "high", idx: cols["high"]},
		{name: "low", idx: cols["low"]},
		{name: "close", idx: cols["close"]},
		{name: "volume", idx: cols["volume"]},
		{name: "close_time", idx: cols["close_time"]},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := append([]string(nil), good...)
			rec[tc.idx] = "bad"
			if _, err := parseCandle(rec, cols); err == nil {
				t.Fatalf("expected %s parse error", tc.name)
			}
		})
	}
}

func TestParseTimeRejectsEmptyTimestamp(t *testing.T) {
	if _, err := parseTime(" "); err == nil {
		t.Fatalf("expected empty timestamp error")
	}
}
