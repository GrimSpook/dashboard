package data

import (
	"encoding/json"
	"os"
)

func SaveCompanies(companies []Company) error {
	data, err := json.MarshalIndent(companies, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("tracker-data.json", data, 0644)
}

func LoadCompanies() ([]Company, error) {
	data, err := os.ReadFile("tracker-data.json")
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
