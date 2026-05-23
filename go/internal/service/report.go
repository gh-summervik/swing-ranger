package service

import (
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/summervik/swing-ranger/internal/model"
)

type ReportService struct{}

func NewReportService() *ReportService {
	return &ReportService{}
}

func (s *ReportService) PrintResults(r model.BacktestResults) {
	fmt.Printf("\n=== BACKTEST RESULTS: %s ===\n", r.EventName)
	fmt.Printf("Total Events: %d\n\n", r.Overall.TotalEvents)

	fmt.Println("OVERALL")
	fmt.Println("Days | Count | Mean % | Median % | Win Rate | Avg Win | Avg Loss | Max Gain | Max Loss | EV %")
	fmt.Println("-----|-------|--------|----------|----------|---------|----------|----------|----------|-----")
	for _, h := range r.Overall.Horizons {
		if h.Count == 0 {
			continue
		}
		fmt.Printf("%4d | %5d | %6.2f%% | %6.2f%% | %7.1f%% | %7.2f%% | %7.2f%% | %7.2f%% | %7.2f%% | %6.2f%%\n",
			h.DaysForward, h.Count,
			h.MeanReturn.Mul(decimal.NewFromInt(100)).Round(2).InexactFloat64(),
			h.MedianReturn.Mul(decimal.NewFromInt(100)).Round(2).InexactFloat64(),
			h.WinRate*100,
			h.AvgWin.Mul(decimal.NewFromInt(100)).Round(2).InexactFloat64(),
			h.AvgLoss.Mul(decimal.NewFromInt(100)).Round(2).InexactFloat64(),
			h.MaxGain.Mul(decimal.NewFromInt(100)).Round(2).InexactFloat64(),
			h.MaxLoss.Mul(decimal.NewFromInt(100)).Round(2).InexactFloat64(),
			h.ExpectedValue.Mul(decimal.NewFromInt(100)).Round(2).InexactFloat64())
	}

	fmt.Println("\nBY SYMBOL")
	for sym, stats := range r.BySymbol {
		fmt.Printf("\n%s (%d events)\n", sym, stats.Horizons[0].Count)
		fmt.Println("Days | Mean % | EV %")
		fmt.Println("-----|--------|-----")
		for _, h := range stats.Horizons {
			if h.Count == 0 {
				continue
			}
			fmt.Printf("%4d | %6.2f%% | %6.2f%%\n",
				h.DaysForward,
				h.MeanReturn.Mul(decimal.NewFromInt(100)).Round(2).InexactFloat64(),
				h.ExpectedValue.Mul(decimal.NewFromInt(100)).Round(2).InexactFloat64())
		}
	}
}
