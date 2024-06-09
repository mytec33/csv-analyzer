package main

import (
	"encoding/json"
	"os"
)

type CsvTest struct {
	Column        int    `json:"Column"`
	Test          string `json:"Test"`
	Value         int64  `json:"Value,omitempty"`
	Values        string `json:"Values,omitempty"`
	DateTimeValue string `json:"DateTimeValue,omitempty"`
	Length        int    `json:"Length,omitempty"`
}

type CsvConfiguration struct {
	Delimiter  string    `json:"Delimiter"`
	MaxColumns int       `json:"MaxColumns"`
	Tests      []CsvTest `json:"Tests"`
}

func ReadConfigFromFile(filename string) (CsvConfiguration, error) {
	var config CsvConfiguration

	file, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
