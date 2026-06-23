package lab

import "math"

type BacktestConfig struct {
	StartBalance   float64
	RiskPct        float64
	MaxNotionalPct float64
	FeePct         float64
	SlippagePct    float64
	MaxHoldBars    int
}

type Position struct {
	Side        Direction `json:"side"`
	EntryPrice  float64   `json:"entry_price"`
	Stop        float64   `json:"stop"`
	Target      float64   `json:"target"`
	Size        float64   `json:"size"`
	OpenIndex   int       `json:"open_index"`
	OpenTime    string    `json:"open_time"`
	EntryFee    float64   `json:"entry_fee"`
	EntrySlip   float64   `json:"entry_slippage"`
	MaxHoldBars int       `json:"max_hold_bars"`
	Reason      string    `json:"reason"`
}

type Trade struct {
	Side       Direction `json:"side"`
	EntryTime  string    `json:"entry_time"`
	ExitTime   string    `json:"exit_time"`
	OpenIndex  int       `json:"open_index"`
	CloseIndex int       `json:"close_index"`
	EntryPrice float64   `json:"entry_price"`
	ExitPrice  float64   `json:"exit_price"`
	Stop       float64   `json:"stop"`
	Target     float64   `json:"target"`
	Size       float64   `json:"size"`
	GrossPnL   float64   `json:"gross_pnl"`
	NetPnL     float64   `json:"net_pnl"`
	Fees       float64   `json:"fees"`
	Slippage   float64   `json:"slippage"`
	Reason     string    `json:"reason"`
	Signal     string    `json:"signal"`
	HoldBars   int       `json:"hold_bars"`
}

type BacktestResult struct {
	Trades      []Trade   `json:"trades"`
	EquityCurve []float64 `json:"equity_curve"`
}

func RunBacktest(candles []Candle, strategy Strategy, cfg BacktestConfig) BacktestResult {
	balance := cfg.StartBalance
	equity := []float64{balance}
	var trades []Trade
	var pos *Position
	var pending *Position

	for i, bar := range candles {
		if pending != nil && pending.OpenIndex == i {
			pos = pending
			pending = nil
		}
		if pos != nil && i >= pos.OpenIndex {
			if tr, ok := maybeExit(*pos, bar, i, cfg); ok {
				trades = append(trades, tr)
				balance += tr.NetPnL
				pos = nil
			}
		}
		equity = append(equity, markToMarket(balance, pos, bar.Close))
		if pos != nil || pending != nil || i+1 >= len(candles) {
			continue
		}
		sig, ok := strategy.OnCandle(StrategyContext{Candles: candles, Index: i})
		if !ok || !validSignal(sig) {
			continue
		}
		entryBar := candles[i+1]
		entryRaw := entryBar.Open
		entry := applySlippage(entryRaw, cfg.SlippagePct, sig.Side, true)
		if !validEntryGeometry(sig, entry) {
			continue
		}
		stopDist := math.Abs(entry - sig.Stop)
		size := positionSize(balance, entry, stopDist, cfg)
		if size <= 0 {
			continue
		}
		maxHold := sig.MaxHoldBars
		if maxHold <= 0 {
			maxHold = cfg.MaxHoldBars
		}
		entryFee := entry * size * cfg.FeePct
		pending = &Position{
			Side:        sig.Side,
			EntryPrice:  entry,
			Stop:        sig.Stop,
			Target:      sig.Target,
			Size:        size,
			OpenIndex:   i + 1,
			OpenTime:    entryBar.OpenTime.Format(timeLayout),
			EntryFee:    entryFee,
			EntrySlip:   math.Abs(entry-entryRaw) * size,
			MaxHoldBars: maxHold,
			Reason:      sig.Reason,
		}
	}
	if pos != nil && len(candles) > 0 {
		lastIdx := len(candles) - 1
		last := candles[lastIdx]
		tr := closePosition(*pos, last.Close, last.Close, "force_close", lastIdx, last.CloseTime.Format(timeLayout), cfg)
		trades = append(trades, tr)
		balance += tr.NetPnL
		equity = append(equity, balance)
	}
	return BacktestResult{Trades: trades, EquityCurve: equity}
}

func validSignal(sig Signal) bool {
	if sig.Side != Long && sig.Side != Short {
		return false
	}
	if sig.Stop <= 0 || sig.Target <= 0 {
		return false
	}
	if sig.Side == Long {
		return sig.Target > sig.Stop
	}
	return sig.Stop > sig.Target
}

func validEntryGeometry(sig Signal, entry float64) bool {
	if entry <= 0 {
		return false
	}
	if sig.Side == Long {
		return sig.Stop < entry && sig.Target > entry
	}
	return sig.Stop > entry && sig.Target < entry
}

func maybeExit(pos Position, bar Candle, idx int, cfg BacktestConfig) (Trade, bool) {
	stopHit, targetHit := false, false
	if pos.Side == Long {
		stopHit = bar.Low <= pos.Stop
		targetHit = bar.High >= pos.Target
	} else {
		stopHit = bar.High >= pos.Stop
		targetHit = bar.Low <= pos.Target
	}
	if stopHit {
		raw := pos.Stop
		exit := applySlippage(raw, cfg.SlippagePct, pos.Side, false)
		return closePosition(pos, exit, raw, "stop_loss", idx, bar.CloseTime.Format(timeLayout), cfg), true
	}
	if targetHit {
		raw := pos.Target
		exit := applySlippage(raw, cfg.SlippagePct, pos.Side, false)
		return closePosition(pos, exit, raw, "take_profit", idx, bar.CloseTime.Format(timeLayout), cfg), true
	}
	if pos.MaxHoldBars > 0 && idx-pos.OpenIndex >= pos.MaxHoldBars {
		raw := bar.Close
		exit := applySlippage(raw, cfg.SlippagePct, pos.Side, false)
		return closePosition(pos, exit, raw, "time_stop", idx, bar.CloseTime.Format(timeLayout), cfg), true
	}
	return Trade{}, false
}

func closePosition(pos Position, exit, rawExit float64, reason string, closeIndex int, exitTime string, cfg BacktestConfig) Trade {
	gross := 0.0
	if pos.Side == Long {
		gross = (exit - pos.EntryPrice) * pos.Size
	} else {
		gross = (pos.EntryPrice - exit) * pos.Size
	}
	exitFee := exit * pos.Size * cfg.FeePct
	exitSlip := math.Abs(exit-rawExit) * pos.Size
	fees := pos.EntryFee + exitFee
	slippage := pos.EntrySlip + exitSlip
	net := gross - fees
	return Trade{
		Side:       pos.Side,
		EntryTime:  pos.OpenTime,
		ExitTime:   exitTime,
		OpenIndex:  pos.OpenIndex,
		CloseIndex: closeIndex,
		EntryPrice: pos.EntryPrice,
		ExitPrice:  exit,
		Stop:       pos.Stop,
		Target:     pos.Target,
		Size:       pos.Size,
		GrossPnL:   gross,
		NetPnL:     net,
		Fees:       fees,
		Slippage:   slippage,
		Reason:     reason,
		Signal:     pos.Reason,
		HoldBars:   closeIndex - pos.OpenIndex,
	}
}

func applySlippage(price, slippagePct float64, side Direction, isEntry bool) float64 {
	if slippagePct <= 0 {
		return price
	}
	if side == Long {
		if isEntry {
			return price * (1 + slippagePct)
		}
		return price * (1 - slippagePct)
	}
	if isEntry {
		return price * (1 - slippagePct)
	}
	return price * (1 + slippagePct)
}

func positionSize(balance, entry, stopDist float64, cfg BacktestConfig) float64 {
	if balance <= 0 || entry <= 0 || stopDist <= 0 || cfg.RiskPct <= 0 {
		return 0
	}
	size := balance * cfg.RiskPct / stopDist
	if cfg.MaxNotionalPct > 0 {
		maxSize := balance * cfg.MaxNotionalPct / entry
		if size > maxSize {
			size = maxSize
		}
	}
	return size
}

func markToMarket(balance float64, pos *Position, close float64) float64 {
	if pos == nil {
		return balance
	}
	if pos.Side == Long {
		return balance + (close-pos.EntryPrice)*pos.Size - pos.EntryFee
	}
	return balance + (pos.EntryPrice-close)*pos.Size - pos.EntryFee
}

const timeLayout = "2006-01-02T15:04:05Z"
