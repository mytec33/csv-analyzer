package main

import (
	"flag"
	"log"
	"sync"
)

func main() {
	multiCore := flag.Bool("multicore", false, "run using multiple cores")
	flag.Parse()

	config, err := ReadConfigFromFile("./Examples/motor_vehicle_collisions_config.json")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	records := readCsvFile("./TestFiles/Motor_Vehicle_Collisions_-_Crashes.csv", config.MaxColumns)

	hasHeader := true
	if hasHeader {
		records = records[1:]
	}

	if *multiCore {
		var wg sync.WaitGroup

		for r, record := range records {
			wg.Add(1)

			go processRecord(config, record, r, &wg)
		}

		wg.Wait()
	} else {
		processCSVData(config, records)
	}

}
