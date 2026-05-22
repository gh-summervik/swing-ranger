package main

import (
	"context"
	"fmt"
	"time"

	"github.com/summervik/swing-ranger/internal/config"
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
		prices, err := yahooSvc.GetHistorical(ctx, symbol, start, end)
		if err != nil {
			comms.Communicate(fmt.Sprintf("Failed to fetch data from Yahoo for %s: %v", symbol, err))
			continue
		}

		err = db.UpsertEodPrices(prices, "system")
		if err != nil {
			return err
		}

		comms.Communicate(fmt.Sprintf("Collected and saved %d records for %s", len(prices), symbol))
	}
	return nil
}
