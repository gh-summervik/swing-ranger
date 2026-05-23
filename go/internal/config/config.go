package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Secrets struct {
	ConnectionStrings map[string]string `json:"ConnectionStrings"`
}

type AppConfig struct {
	Chart ChartConfig `json:"Chart"`
}

type ChartConfig struct {
	MovingAverages             []MovingAverageKey `json:"MovingAverages"`
	BollingBandsMovingAverage  MovingAverageKey   `json:"BollingerBandsMovingAverage"`
	MACD                       MACDConfig         `json:"MACD"`
	RSI                        RSIConfig          `json:"RSI"`
	Backtests                  map[string]BacktestConfig `json:"Backtests"`
}

type BacktestConfig struct {
	Type            string `json:"type"`
	SqueezeLookback int    `json:"squeezeLookback"`
	MinSqueezeBars  int    `json:"minSqueezeBars"`
	MinRSI          int    `json:"minRSI"`
}

type MACDConfig struct {
	FastPeriod   int
	SlowPeriod   int
	SignalPeriod int
	PricePoint   PricePoint
}

type RSIConfig struct {
	Period     int
	PricePoint PricePoint
}

type Config struct {
	ShowHelp  bool
	Verbose   bool
	Command   string
	Data      []string
	Secrets   *Secrets
	AppConfig *AppConfig
}

func LoadSecrets() (*Secrets, error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("could not locate executable: %w", err)
	}
	jsonPath := filepath.Join(filepath.Dir(exe), "secrets.json")
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("secrets.json not found next to binary (%s): %w", jsonPath, err)
	}
	var s Secrets
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("invalid secrets.json: %w", err)
	}
	return &s, nil
}

func LoadAppConfig() (*AppConfig, error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("could not locate executable: %w", err)
	}

	defaultCfg := AppConfig{
		Chart: ChartConfig{
			MovingAverages: []MovingAverageKey{
				{Type: Sma, Period: 20, PricePoint: Close},
				{Type: Sma, Period: 50, PricePoint: Close},
				{Type: Sma, Period: 200, PricePoint: Close},
			},
			BollingBandsMovingAverage: MovingAverageKey{
				Type:       Sma,
				Period:     20,
				PricePoint: Close,
			},
			MACD: MACDConfig{
				FastPeriod:   12,
				SlowPeriod:   26,
				SignalPeriod: 9,
				PricePoint:   Close,
			},
			RSI: RSIConfig{
				Period:     14,
				PricePoint: Close,
			},
			Backtests: map[string]BacktestConfig{
				"squeeze_breakout": {
					Type:            "squeeze_breakout",
					SqueezeLookback: 50,
					MinSqueezeBars:  3,
					MinRSI:          50,
				},
			},
		},
	}

	jsonPath := filepath.Join(filepath.Dir(exe), "config.json")
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return &defaultCfg, nil
	}

	var ac AppConfig
	if err := json.Unmarshal(data, &ac); err != nil {
		return nil, fmt.Errorf("Invalid config.json: %w", err)
	}
	return &ac, nil
}

func ParseMovingAverageKey(s string) (MovingAverageKey, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	if s == "" {
		return MovingAverageKey{}, fmt.Errorf("empty moving average key")
	}

	hint := "try something like '20SC' (the 20 period SMA based on the Close)"

	i := 0
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		i++
	}
	if i == 0 {
		return MovingAverageKey{}, fmt.Errorf("moving average must start with period number, got %q; %s", s, hint)
	}

	periodStr := s[:i]
	rest := s[i:]

	period, err := strconv.Atoi(periodStr)
	if err != nil {
		return MovingAverageKey{}, fmt.Errorf("invalid period (%s) in %q; %s : %w", periodStr, s, hint, err)
	}
	if period < 1 || period > 1000 {
		return MovingAverageKey{}, fmt.Errorf("A moving average period must be between 1 and 1000.")
	}

	if len(rest) != 2 {
		return MovingAverageKey{}, fmt.Errorf("Invalid moving average key: %s. %s", s, hint)
	}

	var key MovingAverageKey
	key.Period = period

	if t, ok := maTypeFromString[string(rest[0])]; ok {
		key.Type = t
	} else {
		return MovingAverageKey{}, fmt.Errorf("unknown MA type '%c' in %q (use S or E)", rest[0], s)
	}

	if p, ok := pricePointFromString[string(rest[1])]; ok {
		key.PricePoint = p
	} else {
		return MovingAverageKey{}, fmt.Errorf("unknown price point '%c' in %q (use O, H, L or C)", rest[1], s)
	}

	return key, nil
}

func (m *MovingAverageKey) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	key, err := ParseMovingAverageKey(s)
	if err != nil {
		return err
	}
	*m = key
	return nil
}

func (m MovingAverageKey) String() string {
	return fmt.Sprintf("%d%s%s", m.Period, m.Type.String(), m.PricePoint.String())
}

func ParseMACDConfig(s string) (MACDConfig, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	if s == "" {
		return MACDConfig{}, fmt.Errorf("empty MACD config")
	}

	parts := strings.Split(s, ",")
	if len(parts) != 3 {
		return MACDConfig{}, fmt.Errorf("MACD must be in format 'fast,slow,signalC' e.g. '12,26,9C'")
	}

	fast, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil || fast < 1 {
		return MACDConfig{}, fmt.Errorf("invalid fast period in MACD")
	}

	slow, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil || slow < 1 {
		return MACDConfig{}, fmt.Errorf("invalid slow period in MACD")
	}

	signalStr := strings.TrimSpace(parts[2])
	if len(signalStr) < 2 {
		return MACDConfig{}, fmt.Errorf("invalid signal part in MACD")
	}

	signal, err := strconv.Atoi(signalStr[:len(signalStr)-1])
	if err != nil || signal < 1 {
		return MACDConfig{}, fmt.Errorf("invalid signal period in MACD")
	}

	ppChar := signalStr[len(signalStr)-1:]
	pp, ok := pricePointFromString[ppChar]
	if !ok {
		return MACDConfig{}, fmt.Errorf("unknown price point '%s' in MACD (use O,H,L,C)", ppChar)
	}

	return MACDConfig{
		FastPeriod:   fast,
		SlowPeriod:   slow,
		SignalPeriod: signal,
		PricePoint:   pp,
	}, nil
}

func (m *MACDConfig) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	cfg, err := ParseMACDConfig(s)
	if err != nil {
		return err
	}
	*m = cfg
	return nil
}

func ParseRSIConfig(s string) (RSIConfig, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	if s == "" {
		return RSIConfig{}, fmt.Errorf("empty RSI config")
	}

	hint := "try something like '14C' (14 period RSI on Close)"

	i := 0
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		i++
	}
	if i == 0 {
		return RSIConfig{}, fmt.Errorf("RSI must start with period number, got %q; %s", s, hint)
	}

	periodStr := s[:i]
	rest := s[i:]

	period, err := strconv.Atoi(periodStr)
	if err != nil {
		return RSIConfig{}, fmt.Errorf("invalid period (%s) in %q; %s : %w", periodStr, s, hint, err)
	}
	if period < 1 || period > 1000 {
		return RSIConfig{}, fmt.Errorf("RSI period must be between 1 and 1000.")
	}

	if len(rest) != 1 {
		return RSIConfig{}, fmt.Errorf("RSI must end with price point (O,H,L,C); got %q; %s", s, hint)
	}

	ppChar := rest
	pp, ok := pricePointFromString[ppChar]
	if !ok {
		return RSIConfig{}, fmt.Errorf("unknown price point '%c' in %q (use O, H, L or C)", rest[0], s)
	}

	return RSIConfig{
		Period:     period,
		PricePoint: pp,
	}, nil
}

func (r *RSIConfig) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	cfg, err := ParseRSIConfig(s)
	if err != nil {
		return err
	}
	*r = cfg
	return nil
}