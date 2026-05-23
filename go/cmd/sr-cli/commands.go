package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/summervik/swing-ranger/internal/config"
	"github.com/summervik/swing-ranger/internal/model"
	"github.com/summervik/swing-ranger/internal/service"

	_ "github.com/lib/pq"
)

func initApp() []string {
	result := make([]string, 0, 10)
	secrets, err := config.LoadSecrets()
	if err != nil {
		result = append(result, "A valid secrets.json file is required at the same level as the application being executed.")
	} else if len(secrets.ConnectionStrings) < 2 {
		result = append(result, "The secrets.json file should contain two connections strings, one with a key of 'Command' and the other with a key of 'Query'. They can be the same connection string.")
	} else if secrets.ConnectionStrings["Command"] == "" {
		result = append(result, "A 'Command' connection was not found in the secrets.json file.")
	} else if secrets.ConnectionStrings["Query"] == "" {
		result = append(result, "A 'Query' connection was not found in the secrets.json file.")
	} else {
		result = append(result, "No issues found; app is initialized.")
	}

	return result
}

func collectSymbols(cfg config.Config, comms *service.CommsService, db *service.DbService) error {
	comms.Communicate(fmt.Sprintf("Collecting historical data for %d symbol(s)", len(cfg.Data)))
	yahooSvc := service.NewYahooService()

	ctx := context.Background()
	start := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Now().UTC()

	for _, symbol := range cfg.Data {
		comms.Communicate(fmt.Sprintf("Fetching historical data for %s", symbol))
		candles, err := yahooSvc.GetHistorical(ctx, symbol, start, end)
		if err != nil {
			comms.Communicate(fmt.Sprintf("Failed to fetch data from Yahoo for %s: %v", symbol, err))
			continue
		}

		err = db.UpsertEodPrices(candles, "system")
		if err != nil {
			return err
		}

		comms.Communicate(fmt.Sprintf("Collected and saved %d records for %s", len(candles), symbol))
	}
	return nil
}

func doBacktest(cfg config.Config, db *service.DbService) error {
	backtestName := cfg.Data[0]
	btCfg, ok := cfg.AppConfig.Chart.Backtests[backtestName]
	if !ok {
		return fmt.Errorf("unknown backtest: %s", backtestName)
	}

	symbols, err := db.GetAllSymbols()
	if err != nil {
		return fmt.Errorf("failed to load symbols: %w", err)
	}
	if len(symbols) == 0 {
		return fmt.Errorf("no symbols found in the database")
	}

	var charts []*model.Chart
	for _, sym := range symbols {
		candles, err := db.GetEodCandlesticks(sym)
		if err != nil {
			fmt.Printf("warning: skipping %s: %v\n", sym, err)
			continue
		}
		if len(candles) < 200 {
			continue
		}

		chart, err := model.NewChart(sym, candles, cfg.AppConfig.Chart)
		if err != nil {
			fmt.Printf("warning: skipping %s (chart creation failed): %v\n", sym, err)
			continue
		}
		charts = append(charts, chart)
	}

	if len(charts) == 0 {
		return fmt.Errorf("no valid charts could be built")
	}

	detector := createSqueezeBreakoutDetector(btCfg)

	backtestSvc := service.NewBacktestService()
	results := backtestSvc.Run(detector, charts, []int{1, 3, 5, 10, 15, 20, 30})

	reportSvc := service.NewReportService()
	reportSvc.PrintResults(results)

	return nil
}

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

func doTest(cfg config.Config, db *service.DbService) error {
	data, _ := json.MarshalIndent(cfg, "", "  ")
	fmt.Println(string(data))

	prices, err := db.GetEodCandlesticks("SPY")
	if err != nil {
		return err
	}

	if len(prices) == 0 {
		return fmt.Errorf("No prices found for SPY")
	}

	chart, err := model.NewChart("SPY", prices, cfg.AppConfig.Chart)
	if err != nil {
		return err
	}

	for i := 0; i < len(chart.Candles); i++ {
		fmt.Printf("%s\t%s\t%s\t%s\t%s\n", chart.Candles[i].DateEod.Format("2006-01-02"), chart.Candles[i].Close.Round(4).String(), chart.MovingAverages["fast"][i].Round(4).String(), chart.MovingAverages["mid"][i].Round(4).String(), chart.MovingAverages["slow"][i].Round(4).String())
	}

	return nil
}
