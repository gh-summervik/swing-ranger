package model

import "time"

type EodPriceDao struct {
	Symbol          string
	DateEod         time.Time
	Open            string
	High            string
	Low             string
	Close           string
	Volume          string
	CreatedBy       string
	UpdatedBy       string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	CreatedAtUnixMs int64
	UpdatedAtUnixMs int64
}

type WatchlistDao struct {
	WatchlistName string
	Symbol        string
	CreatedAt     time.Time
}