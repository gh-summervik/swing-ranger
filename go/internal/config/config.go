package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Secrets struct {
	ConnectionStrings map[string]string `json:"ConnectionStrings"`
}

type Config struct {
	ShowHelp bool
	Verbose  bool
	Command  string
	Data     []string
	Secrets  *Secrets
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