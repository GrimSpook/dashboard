package data

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func getDataDir() (string, string, error) {

	home, err := os.UserHomeDir()
	if err != nil {
		return "", "", err
	}

	dataDir := filepath.Join(home, ".config", "dashboard", "tracker-data")

	file := filepath.Join(dataDir, "tracker-data.json")

	return file, dataDir, nil
}

func SaveCompanies(companies []Company) error {
	file, dir, err := getDataDir()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(companies, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(file, data, 0644)
}

func LoadCompanies() ([]Company, error) {

	file, _, err := getDataDir()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var companies []Company
	err = json.Unmarshal(data, &companies)
	if err != nil {
		return nil, err
	}

	return companies, nil
}
