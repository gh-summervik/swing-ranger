package main

import (
	"github.com/shopspring/decimal"
	"github.com/summervik/swing-ranger/internal/config"
	"github.com/summervik/swing-ranger/internal/model"
	"github.com/summervik/swing-ranger/internal/service"

	_ "github.com/lib/pq"
)

func createSqueezeBreakoutDetector(cfg config.BacktestConfig) service.EventDetector {
	return func(chart *model.Chart) []model.Event {
		if chart == nil {
			return nil
		}
		var events []model.Event
		lookback := cfg.SqueezeLookback
		minBars := cfg.MinSqueezeBars
		minRSI := cfg.MinRSI

		bb := chart.BollingerBands
		if len(bb) == 0 {
			return nil
		}

		for i := lookback; i < len(chart.Candles); i++ {
			minWidth := decimal.NewFromInt(999999)
			minIdx := i
			for j := i - lookback; j <= i; j++ {
				width := bb[model.BBUpper1][j].Sub(bb[model.BBLower1][j])
				if width.LessThan(minWidth) {
					minWidth = width
					minIdx = j
				}
			}

			squeezeCount := 0
			for j := minIdx; j <= i; j++ {
				width := bb[model.BBUpper1][j].Sub(bb[model.BBLower1][j])
				if width.Equal(minWidth) || width.LessThan(minWidth.Mul(decimal.NewFromFloat(1.05))) {
					squeezeCount++
				}
			}
			if squeezeCount < minBars {
				continue
			}

			c := chart.Candles[i]
			if c.Close.GreaterThan(bb[model.BBUpper1][i]) &&
				c.IsBullish &&
				chart.RSI[model.RSIValue][i].GreaterThan(decimal.NewFromInt(int64(minRSI))) {

				events = append(events, model.Event{
					Symbol:    chart.Symbol,
					EventName: "squeeze_breakout",
					DateEod:   c.DateEod,
					Index:     i,
				})
			}
		}
		return events
	}
}
