package model

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/summervik/swing-ranger/internal/config"
	"github.com/shopspring/decimal"
)

type EodCandlestick struct {
	Symbol    string
	DateEod   time.Time
	Open      decimal.Decimal
	High      decimal.Decimal
	Low       decimal.Decimal
	Close     decimal.Decimal
	Volume    float64
	Body      decimal.Decimal
	UpperWick decimal.Decimal
	LowerWick decimal.Decimal
	IsBullish bool
	Range     decimal.Decimal
}

type MA string

const (
	Fast MA = "fast"
	Mid  MA = "mid"
	Slow MA = "slow"
)

type BB string

const (
	BBMiddle BB = "middle"
	BBUpper1 BB = "upper1"
	BBUpper2 BB = "upper2"
	BBUpper3 BB = "upper3"
	BBLower1 BB = "lower1"
	BBLower2 BB = "lower2"
	BBLower3 BB = "lower3"
)

type MACD string

const (
	MACDLine   MACD = "line"
	MACDSignal MACD = "signal"
	MACDHist   MACD = "hist"
)

type RSI string

const (
	RSIValue RSI = "rsi"
)

type Chart struct {
	Symbol         string
	Candles        []EodCandlestick
	MovingAverages map[MA][]decimal.Decimal
	BollingerBands map[BB][]decimal.Decimal
	MACD           map[MACD][]decimal.Decimal
	RSI            map[RSI][]decimal.Decimal
}

func NewChart(symbol string, candles []EodCandlestick, chartConfig config.ChartConfig) (*Chart, error) {
	if len(chartConfig.MovingAverages) != 3 {
		return nil, fmt.Errorf("exactly 3 moving averages are required, got %d", len(chartConfig.MovingAverages))
	}

	mas := make([]config.MovingAverageKey, len(chartConfig.MovingAverages))
	copy(mas, chartConfig.MovingAverages)
	sort.Slice(mas, func(i, j int) bool {
		return mas[i].Period < mas[j].Period
	})

	c := &Chart{
		Symbol:         symbol,
		Candles:        make([]EodCandlestick, len(candles)),
		MovingAverages: make(map[MA][]decimal.Decimal),
		BollingerBands: make(map[BB][]decimal.Decimal),
		MACD:           make(map[MACD][]decimal.Decimal),
		RSI:            make(map[RSI][]decimal.Decimal),
	}
	copy(c.Candles, candles)

	for i := range c.Candles {
		computeCandleMetrics(&c.Candles[i])
	}

	c.MovingAverages[Fast] = calculateMA(c.Candles, mas[0])
	c.MovingAverages[Mid] = calculateMA(c.Candles, mas[1])
	c.MovingAverages[Slow] = calculateMA(c.Candles, mas[2])

	c.BollingerBands = calculateBollingerBands(c.Candles, mas[0])
	c.MACD = calculateMACD(c.Candles, chartConfig.MACD)
	c.RSI = calculateRSI(c.Candles, chartConfig.RSI)

	return c, nil
}

func NewEodCandlestick(symbol string, dateEod time.Time, open, high, low, close decimal.Decimal, volume float64) EodCandlestick {
	c := EodCandlestick{
		Symbol:  symbol,
		DateEod: dateEod,
		Open:    open,
		High:    high,
		Low:     low,
		Close:   close,
		Volume:  volume,
	}
	computeCandleMetrics(&c)
	return c
}

func computeCandleMetrics(c *EodCandlestick) {
	if c == nil {
		return
	}
	c.Body = c.Close.Sub(c.Open).Abs()
	c.Range = c.High.Sub(c.Low)
	if c.Close.GreaterThan(c.Open) {
		c.IsBullish = true
		c.UpperWick = c.High.Sub(c.Close)
		c.LowerWick = c.Open.Sub(c.Low)
	} else {
		c.IsBullish = false
		c.UpperWick = c.High.Sub(c.Open)
		c.LowerWick = c.Close.Sub(c.Low)
	}
}

func calculateMA(candles []EodCandlestick, key config.MovingAverageKey) []decimal.Decimal {
	switch key.Type {
	case config.Sma:
		return calculateSMA(candles, key.Period, key.PricePoint)
	case config.Ema:
		return calculateEMA(candles, key.Period, key.PricePoint)
	default:
		return calculateSMA(candles, key.Period, key.PricePoint)
	}
}

func getPrice(p EodCandlestick, pp config.PricePoint) decimal.Decimal {
	switch pp {
	case config.Open:
		return p.Open
	case config.High:
		return p.High
	case config.Low:
		return p.Low
	case config.Close:
		return p.Close
	}
	return p.Close
}

func calculateSMA(candles []EodCandlestick, period int, pp config.PricePoint) []decimal.Decimal {
	n := len(candles)
	ma := make([]decimal.Decimal, n)
	if n == 0 || period < 1 {
		return ma
	}
	if period > n {
		return ma
	}
	p := period - 1
	for i := 0; i < p; i++ {
		ma[i] = decimal.Zero
	}

	sum := decimal.Zero
	for i := 0; i < period; i++ {
		sum = sum.Add(getPrice(candles[i], pp))
	}
	ma[p] = sum.Div(decimal.NewFromInt(int64(period)))

	for i := period; i < n; i++ {
		sum = sum.Sub(getPrice(candles[i-period], pp)).Add(getPrice(candles[i], pp))
		ma[i] = sum.Div(decimal.NewFromInt(int64(period)))
	}
	return ma
}

func calculateEMA(candles []EodCandlestick, period int, pp config.PricePoint) []decimal.Decimal {
	n := len(candles)
	ma := make([]decimal.Decimal, n)
	if n == 0 || period < 1 {
		return ma
	}
	if period > n {
		return ma
	}

	p := period - 1
	for i := 0; i < p; i++ {
		ma[i] = decimal.Zero
	}

	sum := decimal.Zero
	for i := 0; i < period; i++ {
		sum = sum.Add(getPrice(candles[i], pp))
	}
	ema := sum.Div(decimal.NewFromInt(int64(period)))
	ma[p] = ema

	multiplier := decimal.NewFromInt(2).Div(decimal.NewFromInt(int64(period + 1)))
	oneMinusMult := decimal.NewFromInt(1).Sub(multiplier)

	for i := period; i < n; i++ {
		current := getPrice(candles[i], pp)
		ema = current.Mul(multiplier).Add(ema.Mul(oneMinusMult))
		ma[i] = ema
	}
	return ma
}

func calculateBollingerBands(candles []EodCandlestick, baseKey config.MovingAverageKey) map[BB][]decimal.Decimal {
	n := len(candles)
	bb := make(map[BB][]decimal.Decimal)
	if n == 0 {
		return bb
	}

	middle := calculateMA(candles, baseKey)
	stdDevs := calculateRollingStdDev(candles, baseKey.Period, baseKey.PricePoint)

	bb[BBMiddle] = middle
	bb[BBUpper1] = shiftByStdDev(middle, stdDevs, 1)
	bb[BBUpper2] = shiftByStdDev(middle, stdDevs, 2)
	bb[BBUpper3] = shiftByStdDev(middle, stdDevs, 3)
	bb[BBLower1] = shiftByStdDev(middle, stdDevs, -1)
	bb[BBLower2] = shiftByStdDev(middle, stdDevs, -2)
	bb[BBLower3] = shiftByStdDev(middle, stdDevs, -3)

	return bb
}

func calculateRollingStdDev(candles []EodCandlestick, period int, pp config.PricePoint) []decimal.Decimal {
	n := len(candles)
	stds := make([]decimal.Decimal, n)
	if n == 0 || period < 2 {
		return stds
	}
	p := period - 1
	for i := 0; i < p; i++ {
		stds[i] = decimal.Zero
	}

	for i := p; i < n; i++ {
		sum := 0.0
		sumSq := 0.0
		count := 0
		for j := i - period + 1; j <= i; j++ {
			val, _ := getPrice(candles[j], pp).Float64()
			sum += val
			sumSq += val * val
			count++
		}
		if count == 0 {
			stds[i] = decimal.Zero
			continue
		}
		mean := sum / float64(count)
		variance := (sumSq / float64(count)) - (mean * mean)
		if variance > 0 {
			stds[i] = decimal.NewFromFloat(math.Sqrt(variance))
		} else {
			stds[i] = decimal.Zero
		}
	}
	return stds
}

func shiftByStdDev(base, stddevs []decimal.Decimal, k int) []decimal.Decimal {
	n := len(base)
	result := make([]decimal.Decimal, n)
	mult := decimal.NewFromInt(int64(k))
	for i := range base {
		if stddevs[i].IsZero() {
			result[i] = decimal.Zero
		} else {
			dev := stddevs[i].Mul(mult)
			if k >= 0 {
				result[i] = base[i].Add(dev)
			} else {
				result[i] = base[i].Sub(dev.Abs())
			}
		}
	}
	return result
}

func calculateMACD(candles []EodCandlestick, cfg config.MACDConfig) map[MACD][]decimal.Decimal {
	n := len(candles)
	m := make(map[MACD][]decimal.Decimal)
	if n == 0 {
		return m
	}

	fastKey := config.MovingAverageKey{Type: config.Ema, Period: cfg.FastPeriod, PricePoint: cfg.PricePoint}
	slowKey := config.MovingAverageKey{Type: config.Ema, Period: cfg.SlowPeriod, PricePoint: cfg.PricePoint}

	fast := calculateMA(candles, fastKey)
	slow := calculateMA(candles, slowKey)

	macdLine := make([]decimal.Decimal, n)
	for i := 0; i < n; i++ {
		macdLine[i] = fast[i].Sub(slow[i])
	}

	signal := calculateEMAOnSeries(macdLine, cfg.SignalPeriod)

	hist := make([]decimal.Decimal, n)
	for i := 0; i < n; i++ {
		hist[i] = macdLine[i].Sub(signal[i])
	}

	m[MACDLine] = macdLine
	m[MACDSignal] = signal
	m[MACDHist] = hist

	return m
}

func calculateEMAOnSeries(series []decimal.Decimal, period int) []decimal.Decimal {
	n := len(series)
	ma := make([]decimal.Decimal, n)
	if n == 0 || period < 1 {
		return ma
	}
	if period > n {
		return ma
	}

	p := period - 1
	for i := 0; i < p; i++ {
		ma[i] = decimal.Zero
	}

	sum := decimal.Zero
	for i := 0; i < period; i++ {
		sum = sum.Add(series[i])
	}
	ema := sum.Div(decimal.NewFromInt(int64(period)))
	ma[p] = ema

	multiplier := decimal.NewFromInt(2).Div(decimal.NewFromInt(int64(period + 1)))
	oneMinusMult := decimal.NewFromInt(1).Sub(multiplier)

	for i := period; i < n; i++ {
		current := series[i]
		ema = current.Mul(multiplier).Add(ema.Mul(oneMinusMult))
		ma[i] = ema
	}
	return ma
}

func calculateRSI(candles []EodCandlestick, cfg config.RSIConfig) map[RSI][]decimal.Decimal {
	n := len(candles)
	rsiMap := make(map[RSI][]decimal.Decimal)
	if n == 0 {
		return rsiMap
	}

	rsi := make([]decimal.Decimal, n)
	period := cfg.Period
	if period < 2 || period > n {
		for i := range rsi {
			rsi[i] = decimal.Zero
		}
		rsiMap[RSIValue] = rsi
		return rsiMap
	}

	// leading zeros
	p := period - 1
	for i := 0; i < p; i++ {
		rsi[i] = decimal.Zero
	}

	// first RSI value uses simple average
	sumGain := decimal.Zero
	sumLoss := decimal.Zero
	for i := 1; i <= period; i++ {
		change := getPrice(candles[i], cfg.PricePoint).Sub(getPrice(candles[i-1], cfg.PricePoint))
		if change.GreaterThan(decimal.Zero) {
			sumGain = sumGain.Add(change)
		} else {
			sumLoss = sumLoss.Add(change.Abs())
		}
	}
	avgGain := sumGain.Div(decimal.NewFromInt(int64(period)))
	avgLoss := sumLoss.Div(decimal.NewFromInt(int64(period)))

	if avgLoss.IsZero() {
		rsi[p] = decimal.NewFromInt(100)
	} else {
		rs := avgGain.Div(avgLoss)
		rsi[p] = decimal.NewFromInt(100).Sub(decimal.NewFromInt(100).Div(decimal.NewFromInt(1).Add(rs)))
	}

	// subsequent values use Wilder's smoothing
	for i := period + 1; i < n; i++ {
		change := getPrice(candles[i], cfg.PricePoint).Sub(getPrice(candles[i-1], cfg.PricePoint))
		var curGain, curLoss decimal.Decimal
		if change.GreaterThan(decimal.Zero) {
			curGain = change
			curLoss = decimal.Zero
		} else {
			curGain = decimal.Zero
			curLoss = change.Abs()
		}

		avgGain = avgGain.Mul(decimal.NewFromInt(int64(period-1))).Add(curGain).Div(decimal.NewFromInt(int64(period)))
		avgLoss = avgLoss.Mul(decimal.NewFromInt(int64(period-1))).Add(curLoss).Div(decimal.NewFromInt(int64(period)))

		if avgLoss.IsZero() {
			rsi[i-1] = decimal.NewFromInt(100)
		} else {
			rs := avgGain.Div(avgLoss)
			rsi[i-1] = decimal.NewFromInt(100).Sub(decimal.NewFromInt(100).Div(decimal.NewFromInt(1).Add(rs)))
		}
	}

	rsiMap[RSIValue] = rsi
	return rsiMap
}