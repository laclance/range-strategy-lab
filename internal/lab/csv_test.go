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
