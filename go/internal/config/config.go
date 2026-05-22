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
	MovingAverages []MovingAverageKey `json:"MovingAverages"`
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

func (m *MovingAverageKey) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	s = strings.TrimSpace(strings.ToUpper(s))
	if s == "" {
		return fmt.Errorf("empty moving average key")
	}

	hint := "try something like '20SC' (the 20 period SMA based on the Close)"

	i := 0
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		i++
	}
	if i == 0 {
		return fmt.Errorf("moving average must start with period number, got %q; %s", s, hint)
	}

	periodStr := s[:i]
	rest := s[i:]

	period, err := strconv.Atoi(periodStr)
	if err != nil {
		return fmt.Errorf("invalid period (%s) in %q; %s : %w", periodStr, s, hint, err)
	}
	if period < 1 || period > 1000 {
		return fmt.Errorf("A moving average period must be between 1 and 1000.")
	}
	m.Period = period

	if len(rest) != 2 {
		return fmt.Errorf("Invalid moving average key: %s. %s", s, hint)
	}

	// use the enum string values from the maps instead of hard-coding
	if t, ok := maTypeFromString[string(rest[0])]; ok {
		m.Type = t
	} else {
		return fmt.Errorf("unknown MA type '%c' in %q (use S or E)", rest[0], s)
	}

	if p, ok := pricePointFromString[string(rest[1])]; ok {
		m.PricePoint = p
	} else {
		return fmt.Errorf("unknown price point '%c' in %q (use O, H, L or C)", rest[1], s)
	}

	return nil
}

func (m MovingAverageKey) String() string {
	return fmt.Sprintf("%d%s%s", m.Period, m.Type.String(), m.PricePoint.String())
}