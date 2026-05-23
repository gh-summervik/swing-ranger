package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type Event struct {
	Symbol    string
	EventName string
	DateEod   time.Time
	Index     int
}

type HorizonStats struct {
	DaysForward    int
	Count          int
	MeanReturn     decimal.Decimal
	MedianReturn   decimal.Decimal
	WinRate        float64          // 0.0 to 1.0
	AvgWin         decimal.Decimal
	AvgLoss        decimal.Decimal
	MaxGain        decimal.Decimal
	MaxLoss        decimal.Decimal
	ExpectedValue  decimal.Decimal  
}

type PerformanceStats struct {
	EventName   string
	TotalEvents int
	Horizons    []HorizonStats // sorted by DaysForward
}

type BacktestResults struct {
	EventName string
	Overall   PerformanceStats
	BySymbol  map[string]PerformanceStats
	BySector  map[string]PerformanceStats
}