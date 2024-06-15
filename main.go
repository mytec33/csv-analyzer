package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
	"sync"
)

const batchSize = 1000 // Adjust batch size as needed

func main() {
	configFile := flag.String("cfile", "", "json configuration file containing CSV tests")
	csvFile := flag.String("csv", "", "csv file to analyze")
	multiCore := flag.Bool("multicore", false, "run using multiple cores")
	cpuProfile := flag.String("cpuprofile", "", "write cpu profile to `file`")
	flag.Parse()

	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	config, err := ReadConfigFromFile(*configFile)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	testsByColumn := PreprocessTests(config.Tests)

	records := readCsvFile(*csvFile, config.MaxColumns)

	hasHeader := true
	if hasHeader {
		records = records[1:]
	}

	if *multiCore {
		var wg sync.WaitGroup

		recordPool := sync.Pool{
			New: func() interface{} {
				slice := make([][]string, 0, batchSize)
				return &slice
			},
		}

		for i := 0; i < len(records); i += batchSize {
			end := i + batchSize
			if end > len(records) {
				end = len(records)
			}

			pooledBatch := recordPool.Get().(*[][]string)
			*pooledBatch = append((*pooledBatch)[:0], records[i:end]...)

			wg.Add(1)
			go func(batch *[][]string, startIdx int) {
				defer wg.Done()
				for j, record := range *batch {
					processRecord(record, startIdx+j, testsByColumn)
				}
				recordPool.Put(batch)
			}(pooledBatch, i)
		}

		wg.Wait()
	} else {
		processCSVData(config, records)
	}

}
