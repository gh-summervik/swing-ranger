package service

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/summervik/swing-ranger/internal/model"

	"github.com/shopspring/decimal"
)

func GetEodPrices(db *sql.DB, symbol string) ([]model.EodPrice, error) {
	rows, err := db.Query(`
		SELECT symbol, date_eod, open, high, low, close, volume
		FROM public.eod_prices 
		WHERE symbol = $1 
		ORDER BY date_eod ASC
	`, symbol)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prices []model.EodPrice
	for rows.Next() {
		var dao model.EodPriceDao
		if err := rows.Scan(
			&dao.Symbol, &dao.DateEod, &dao.Open, &dao.High, &dao.Low,
			&dao.Close, &dao.Volume,
		); err != nil {
			return nil, err
		}

		open, _ := decimal.NewFromString(dao.Open)
		high, _ := decimal.NewFromString(dao.High)
		low, _ := decimal.NewFromString(dao.Low)
		closeP, _ := decimal.NewFromString(dao.Close)
		volume, _ := strconv.ParseFloat(dao.Volume, 64)

		prices = append(prices, model.EodPrice{
			Symbol:  dao.Symbol,
			DateEod: dao.DateEod,
			Open:    open,
			High:    high,
			Low:     low,
			Close:   closeP,
			Volume:  volume,
		})
	}
	return prices, rows.Err()
}

func UpsertEodPrices(db *sql.DB, prices []model.EodPrice, by string) error {
	now := time.Now().UTC()
	unixms := now.UnixMilli()

	for _, p := range prices {
		_, err := db.Exec(`
			INSERT INTO public.eod_prices 
			(symbol, date_eod, open, high, low, close, volume, 
			 created_by, updated_by, created_at, updated_at, 
			 created_at_unix_ms, updated_at_unix_ms)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
			ON CONFLICT (symbol, date_eod) DO UPDATE SET
				open = EXCLUDED.open,
				high = EXCLUDED.high,
				low = EXCLUDED.low,
				close = EXCLUDED.close,
				volume = EXCLUDED.volume,
				updated_by = EXCLUDED.updated_by,
				updated_at = EXCLUDED.updated_at,
				updated_at_unix_ms = EXCLUDED.updated_at_unix_ms
		`, p.Symbol, p.DateEod, p.Open, p.High, p.Low,
			p.Close, p.Volume, by, by, now, now, unixms, unixms)
		if err != nil {
			return err
		}
	}
	return nil
}