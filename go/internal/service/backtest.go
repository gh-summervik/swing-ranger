package service

import (
	"sort"

	"github.com/summervik/swing-ranger/internal/model"
	"github.com/shopspring/decimal"
)

type BacktestService struct{}

func NewBacktestService() *BacktestService {
	return &BacktestService{}
}

type EventDetector func(*model.Chart) []model.Event

func (s *BacktestService) Run(detector EventDetector, charts []*model.Chart, horizons []int) model.BacktestResults {
	if len(horizons) == 0 {
		horizons = []int{1, 3, 5, 10, 15, 20, 30}
	}

	results := model.BacktestResults{
		EventName: "",
		Overall: model.PerformanceStats{
			Horizons: make([]model.HorizonStats, len(horizons)),
		},
		BySymbol: make(map[string]model.PerformanceStats),
		BySector: make(map[string]model.PerformanceStats),
	}

	allReturns := make(map[int][]decimal.Decimal)

	for _, chart := range charts {
		if chart == nil || len(chart.Candles) == 0 {
			continue
		}

		events := detector(chart)
		if len(events) == 0 {
			continue
		}

		if results.EventName == "" && len(events) > 0 {
			results.EventName = events[0].EventName
		}

		symbolStats := results.BySymbol[chart.Symbol]
		if symbolStats.Horizons == nil {
			symbolStats = model.PerformanceStats{
				EventName: results.EventName,
				Horizons:  make([]model.HorizonStats, len(horizons)),
			}
		}

		symbolReturns := make(map[int][]decimal.Decimal) // per-symbol returns

		for _, e := range events {
			if e.Index < 0 || e.Index >= len(chart.Candles) {
				continue
			}
			if !chart.Candles[e.Index].DateEod.Equal(e.DateEod) {
				continue
			}

			basePrice := chart.Candles[e.Index].Close

			for hIdx, days := range horizons {
				futureIdx := e.Index + days
				if futureIdx >= len(chart.Candles) {
					continue
				}

				futurePrice := chart.Candles[futureIdx].Close
				if basePrice.IsZero() {
					continue
				}

				ret := futurePrice.Sub(basePrice).Div(basePrice)

				// overall
				allReturns[days] = append(allReturns[days], ret)

				// per-symbol
				symbolReturns[days] = append(symbolReturns[days], ret)

				symbolStats.Horizons[hIdx].DaysForward = days
				symbolStats.Horizons[hIdx].Count++
				symbolStats.Horizons[hIdx].MeanReturn = symbolStats.Horizons[hIdx].MeanReturn.Add(ret)
			}
		}

		// finalize this symbol's stats
		for hIdx, days := range horizons {
			returns := symbolReturns[days]
			if len(returns) == 0 {
				continue
			}

			sum := decimal.Zero
			for _, r := range returns {
				sum = sum.Add(r)
			}
			mean := sum.Div(decimal.NewFromInt(int64(len(returns))))

			sorted := make([]decimal.Decimal, len(returns))
			copy(sorted, returns)
			sort.Slice(sorted, func(i, j int) bool { return sorted[i].LessThan(sorted[j]) })
			median := sorted[len(sorted)/2]

			var wins, losses int
			var totalWin, totalLoss, maxGain, maxLoss decimal.Decimal

			for _, r := range returns {
				if r.GreaterThan(decimal.Zero) {
					wins++
					totalWin = totalWin.Add(r)
					if r.GreaterThan(maxGain) {
						maxGain = r
					}
				} else if r.LessThan(decimal.Zero) {
					losses++
					totalLoss = totalLoss.Add(r.Abs())
					if r.LessThan(maxLoss) {
						maxLoss = r
					}
				}
			}

			winRate := 0.0
			if wins+losses > 0 {
				winRate = float64(wins) / float64(wins+losses)
			}

			avgWin := decimal.Zero
			if wins > 0 {
				avgWin = totalWin.Div(decimal.NewFromInt(int64(wins)))
			}
			avgLoss := decimal.Zero
			if losses > 0 {
				avgLoss = totalLoss.Div(decimal.NewFromInt(int64(losses)))
			}

			// Expected Value
			expectedValue := decimal.Zero
			if wins+losses > 0 {
				lossRate := decimal.NewFromFloat(1 - winRate)
				expectedValue = avgWin.Mul(decimal.NewFromFloat(winRate)).Add(avgLoss.Neg().Mul(lossRate))
			}

			symbolStats.Horizons[hIdx] = model.HorizonStats{
				DaysForward:   days,
				Count:         len(returns),
				MeanReturn:    mean,
				MedianReturn:  median,
				WinRate:       winRate,
				AvgWin:        avgWin,
				AvgLoss:       avgLoss,
				MaxGain:       maxGain,
				MaxLoss:       maxLoss,
				ExpectedValue: expectedValue,
			}
		}

		results.BySymbol[chart.Symbol] = symbolStats
	}

	// Overall (unchanged)
	results.Overall.EventName = results.EventName
	if len(allReturns) > 0 && len(horizons) > 0 {
		results.Overall.TotalEvents = len(allReturns[horizons[0]])
	}

	for hIdx, days := range horizons {
		returns := allReturns[days]
		if len(returns) == 0 {
			continue
		}

		sum := decimal.Zero
		for _, r := range returns {
			sum = sum.Add(r)
		}
		mean := sum.Div(decimal.NewFromInt(int64(len(returns))))

		sorted := make([]decimal.Decimal, len(returns))
		copy(sorted, returns)
		sort.Slice(sorted, func(i, j int) bool { return sorted[i].LessThan(sorted[j]) })
		median := sorted[len(sorted)/2]

		var wins, losses int
		var totalWin, totalLoss, maxGain, maxLoss decimal.Decimal

		for _, r := range returns {
			if r.GreaterThan(decimal.Zero) {
				wins++
				totalWin = totalWin.Add(r)
				if r.GreaterThan(maxGain) {
					maxGain = r
				}
			} else if r.LessThan(decimal.Zero) {
				losses++
				totalLoss = totalLoss.Add(r.Abs())
				if r.LessThan(maxLoss) {
					maxLoss = r
				}
			}
		}

		winRate := 0.0
		if wins+losses > 0 {
			winRate = float64(wins) / float64(wins+losses)
		}

		avgWin := decimal.Zero
		if wins > 0 {
			avgWin = totalWin.Div(decimal.NewFromInt(int64(wins)))
		}
		avgLoss := decimal.Zero
		if losses > 0 {
			avgLoss = totalLoss.Div(decimal.NewFromInt(int64(losses)))
		}

		expectedValue := decimal.Zero
		if wins+losses > 0 {
			lossRate := decimal.NewFromFloat(1 - winRate)
			expectedValue = avgWin.Mul(decimal.NewFromFloat(winRate)).Add(avgLoss.Neg().Mul(lossRate))
		}

		results.Overall.Horizons[hIdx] = model.HorizonStats{
			DaysForward:   days,
			Count:         len(returns),
			MeanReturn:    mean,
			MedianReturn:  median,
			WinRate:       winRate,
			AvgWin:        avgWin,
			AvgLoss:       avgLoss,
			MaxGain:       maxGain,
			MaxLoss:       maxLoss,
			ExpectedValue: expectedValue,
		}
	}

	return results
}