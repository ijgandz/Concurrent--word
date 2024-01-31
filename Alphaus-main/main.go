package main

import (
	"encoding/csv"
	"flag"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/alphauslabs/bluectl/pkg/logger"
)

var (
	file         = flag.String("file", "testcur.csv", "Sample file to process")
	wordToSearch = flag.String("word", "", "Specify the word to count occurrences")
)

func concurrent() int {
	fileContent, err := os.ReadFile(*file)
	if err != nil {
		logger.Error(err)
		return 0
	}

	reader := csv.NewReader(strings.NewReader(string(fileContent)))
	records, err := reader.ReadAll()
	if err != nil {
		logger.Error(err)
		return 0
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex
	var wordCount int

	processSlice := func(slice []string) {
		defer wg.Done()
		localCount := 0
		for _, word := range slice {
			// Use strings.Count for case-sensitive search
			localCount += strings.Count(word, *wordToSearch)
		}
		mutex.Lock()
		wordCount += localCount
		mutex.Unlock()
	}

	numCpu := len(records)
	for i := 0; i < numCpu; i++ {
		wg.Add(1)
		go processSlice(records[i])
	}

	wg.Wait()

	return wordCount
}

func main() {
	flag.Parse()
	if *file == "" || *wordToSearch == "" {
		logger.Error("Both -file and -word must be provided")
		return
	}

	// Concurrent Execution
	concurrentStart := time.Now()
	concurrentCount := concurrent()
	concurrentDuration := time.Since(concurrentStart)

	logger.Infof("Concurrent word count for \"%s\": %d, Duration: %v", *wordToSearch, concurrentCount, concurrentDuration)
}
