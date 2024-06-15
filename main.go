package main

import (
	"flag"
	"log"
	"sync"
)

func main() {
	configFile := flag.String("cfile", "", "json configuration file containing CSV tests")
	csvFile := flag.String("csv", "", "csv file to analyze")
	multiCore := flag.Bool("multicore", false, "run using multiple cores")
	flag.Parse()

	config, err := ReadConfigFromFile(*configFile)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	records := readCsvFile(*csvFile, config.MaxColumns)

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
