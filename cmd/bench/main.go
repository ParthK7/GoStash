// I have not written this file myself. Used Gemini to give me the tentative load test file.

package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	const (
		url         = "http://localhost:8080/ingest"
		concurrency = 50    // Number of concurrent "users"
		totalReqs   = 10000 // Total logs to send
	)

	payload := []byte("Benchmark log entry: system status nominal")
	var wg sync.WaitGroup
	start := time.Now()

	reqChan := make(chan struct{}, totalReqs)
	for i := 0; i < totalReqs; i++ {
		reqChan <- struct{}{}
	}
	close(reqChan)

	fmt.Printf("Starting benchmark: %d requests with %d workers...\n", totalReqs, concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Add(-1)
			for range reqChan {
				resp, err := http.Post(url, "text/plain", bytes.NewBuffer(payload))
				if err == nil {
					resp.Body.Close()
				}
			}
		}()
	}

	wg.Wait()
	duration := time.Since(start)
	rps := float64(totalReqs) / duration.Seconds()

	fmt.Println("--- Benchmark Results ---")
	fmt.Printf("Total Time: %v\n", duration)
	fmt.Printf("Throughput: %.2f requests/sec\n", rps)
	fmt.Printf("Avg Latency: %v\n", duration/time.Duration(totalReqs))
}
