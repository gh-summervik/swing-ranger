package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"github.com/summervik/swing-ranger/internal/model"
)

type YahooService struct {
	client *http.Client
}

func NewYahooService() *YahooService {
	return &YahooService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type yahooChartResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Symbol           string `json:"symbol"`
				Currency         string `json:"currency"`
				ExchangeTimezone string `json:"exchangeTimezone"`
			} `json:"meta"`
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Open   []float64 `json:"open"`
					High   []float64 `json:"high"`
					Low    []float64 `json:"low"`
					Close  []float64 `json:"close"`
					Volume []int64   `json:"volume"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
	} `json:"chart"`
}

func (s *YahooService) GetHistorical(ctx context.Context, symbol string, start, end time.Time) ([]model.EodCandlestick, error) {
	period1 := start.Unix()
	period2 := end.Unix()

	u := url.URL{
		Scheme: "https",
		Host:   "query1.finance.yahoo.com",
		Path:   "/v8/finance/chart/" + symbol,
	}
	q := u.Query()
	q.Set("period1", strconv.FormatInt(period1, 10))
	q.Set("period2", strconv.FormatInt(period2, 10))
	q.Set("interval", "1d")
	q.Set("includePrePost", "false")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("yahoo finance returned status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var chartResp yahooChartResponse
	if err := json.Unmarshal(body, &chartResp); err != nil {
		return nil, err
	}

	if len(chartResp.Chart.Result) == 0 {
		return nil, fmt.Errorf("no data returned for symbol %s", symbol)
	}

	result := chartResp.Chart.Result[0]
	if len(result.Timestamp) == 0 {
		return nil, fmt.Errorf("no timestamps in response")
	}

	eodCandles := make([]model.EodCandlestick, 0, len(result.Timestamp))
	quoteData := result.Indicators.Quote[0]

	for i, ts := range result.Timestamp {
		date := time.Unix(ts, 0).UTC()

		openVal := 0.0
		if i < len(quoteData.Open) {
			openVal = quoteData.Open[i]
		}
		highVal := 0.0
		if i < len(quoteData.High) {
			highVal = quoteData.High[i]
		}
		lowVal := 0.0
		if i < len(quoteData.Low) {
			lowVal = quoteData.Low[i]
		}
		closeVal := 0.0
		if i < len(quoteData.Close) {
			closeVal = quoteData.Close[i]
		}
		volumeVal := int64(0)
		if i < len(quoteData.Volume) {
			volumeVal = quoteData.Volume[i]
		}

		eodCandles = append(eodCandles, model.NewEodCandlestick(
			result.Meta.Symbol,
			date,
			decimal.NewFromFloat(openVal),
			decimal.NewFromFloat(highVal),
			decimal.NewFromFloat(lowVal),
			decimal.NewFromFloat(closeVal),
			float64(volumeVal),
		))
	}

	return eodCandles, nil
}