package main

import (
	"log"
)

func main() {
	config, err := ReadConfigFromFile("./Examples/motor_vehicle_collisions_config.json")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	records := readCsvFile("./TestFiles/Motor_Vehicle_Collisions_-_Crashes.csv", config.MaxColumns)

	hasHeader := true
	if hasHeader {
		records = records[1:]
	}

	processCSVData(config, records)
}
