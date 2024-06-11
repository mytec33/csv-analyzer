package main

import (
	"log"
	"sync"
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

	var wg sync.WaitGroup

	// Don't delete! Have command line option to run single core or mutli
	// processCSVData(config, records)
	for r, record := range records {
		wg.Add(1)

		go processRecord(config, record, r, &wg)
	}

	wg.Wait()
}
