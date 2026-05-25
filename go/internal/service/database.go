package service

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/summervik/swing-ranger/internal/config"
	"github.com/summervik/swing-ranger/internal/model"

	"github.com/shopspring/decimal"
)

type DbService struct {
	Command *sql.DB
	Query   *sql.DB
	Comms   *CommsService
}

func NewDbService(cfg config.Config, comms *CommsService) (*DbService, error) {
	cmdDb, err := sql.Open("postgres", cfg.Secrets.ConnectionStrings["Command"])
	if err != nil {
		return nil, err
	}

	qryDb, err := sql.Open("postgres", cfg.Secrets.ConnectionStrings["Query"])
	if err != nil {
		cmdDb.Close()
		return nil, err
	}

	return &DbService{
		Command: cmdDb,
		Query:   qryDb,
		Comms:   comms,
	}, nil
}

// GetEodCandlesticks returns all EOD data for a single symbol
func (s *DbService) GetEodCandlesticks(symbol string) ([]model.EodCandlestick, error) {
	rows, err := s.Query.Query(`
		SELECT symbol, date_eod, open, high, low, close, volume
		FROM public.eod_prices 
		WHERE symbol = $1 
		ORDER BY date_eod ASC
	`, symbol)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var candles []model.EodCandlestick
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

		candles = append(candles, model.NewEodCandlestick(
			dao.Symbol,
			dao.DateEod,
			open,
			high,
			low,
			closeP,
			volume,
		))
	}
	return candles, rows.Err()
}

// GetAllSymbols returns every symbol that has at least one EOD record
func (s *DbService) GetAllSymbols() ([]string, error) {
	rows, err := s.Query.Query(`
		SELECT DISTINCT symbol 
		FROM public.eod_prices 
		ORDER BY symbol ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var symbols []string
	for rows.Next() {
		var sym string
		if err := rows.Scan(&sym); err != nil {
			return nil, err
		}
		symbols = append(symbols, sym)
	}
	return symbols, rows.Err()
}

// UpsertEodPrices inserts or updates price records
func (s *DbService) UpsertEodPrices(candles []model.EodCandlestick, by string) error {
	now := time.Now().UTC()
	unixms := now.UnixMilli()

	for _, c := range candles {
		_, err := s.Command.Exec(`
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
		`, c.Symbol, c.DateEod, c.Open, c.High, c.Low,
			c.Close, c.Volume, by, by, now, now, unixms, unixms)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetWatchlistSymbols returns all symbols in a named watchlist
func (s *DbService) GetWatchlistSymbols(watchlistName string) ([]string, error) {
	rows, err := s.Query.Query(`
		SELECT symbol 
		FROM public.watchlists 
		WHERE watchlist_name = $1 
		ORDER BY symbol ASC
	`, watchlistName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var symbols []string
	for rows.Next() {
		var sym string
		if err := rows.Scan(&sym); err != nil {
			return nil, err
		}
		symbols = append(symbols, sym)
	}
	return symbols, rows.Err()
}

func (s *DbService) UpdateWatchlists(data map[string][]string) error {
	for watchlist, symbols := range data {
		_, err := s.Command.Exec(`DELETE FROM public.watchlists WHERE watchlist_name = $1`, watchlist)
		if err != nil {
			return err
		}

		for _, sym := range symbols {
			if !isValidTicker(sym) {
				continue
			}
			_, err := s.Command.Exec(`
			INSERT INTO public.watchlists (watchlist_name, symbol)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`, watchlist, sym)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isValidTicker(s string) bool {
	if len(s) < 1 || len(s) > 10 {
		return false
	}
	for _, c := range s {
		if !((c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-') {
			return false
		}
	}
	return true
}
