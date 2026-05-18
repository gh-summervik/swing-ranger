package model

import (
	"time"
	"github.com/shopspring/decimal"
)

type EodPrice struct {
	Symbol  string
	DateEod time.Time
	Open    decimal.Decimal
	High    decimal.Decimal
	Low     decimal.Decimal
	Close   decimal.Decimal
	Volume  float64
}