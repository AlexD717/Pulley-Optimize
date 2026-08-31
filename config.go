package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type AppConfig struct {
	UnitSelected      int    `json:"unit-selected"`
	Ratio             string `json:"ratio"`
	UseAvailableBelts bool   `json:"use-available-belts"`
	MaxSlack          string `json:"max-slack"`
	MinSlack          string `json:"min-slack"`
	MaxPulley         string `json:"max-pulley"`
	MinPulley         string `json:"min-pulley"`
}

func defaultConfig() AppConfig {
	return AppConfig{
		UnitSelected:      0,
		Ratio:             "1.0",
		UseAvailableBelts: true,
		MaxSlack:          "0.5",
		MinSlack:          "-0.2",
		MaxPulley:         "100",
		MinPulley:         "8",
	}
}

func getConfigPath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	exeDir := filepath.Dir(exePath)
	return filepath.Join(exeDir, "config.json"), nil
}

func loadConfig() AppConfig {
	cfg := defaultConfig()

	path, err := getConfigPath()
	if err == nil {
		data, err := os.ReadFile(path)
		if err == nil {
			json.Unmarshal(data, &cfg)
		}
	}
	return cfg
}

func saveConfig(m Model) {
	cfg := AppConfig{
		UnitSelected:      m.Inputs[1].Selected,
		Ratio:             m.Inputs[2].Input.Value(),
		UseAvailableBelts: m.Inputs[3].Checked,
		MaxSlack:          m.Inputs[4].Input.Value(),
		MinSlack:          m.Inputs[5].Input.Value(),
		MaxPulley:         m.Inputs[6].Input.Value(),
		MinPulley:         m.Inputs[7].Input.Value(),
	}

	path, err := getConfigPath()
	if err == nil {
		data, _ := json.MarshalIndent(cfg, "", "  ")
		os.WriteFile(path, data, 0644)
	}
}
