package main

import (
	"fmt"
	"log"

	"github.com/ParthK7/GoStash/internal/wal"
)

func main() {
	logger, err := wal.NewWal("storage/active.log")
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 1000; i++ {
		msg := fmt.Sprintf("Some json data for the %dth time", i)
		err = logger.Write(msg)
		if err != nil {
			fmt.Printf("Error writing: %v \n", err)
		}
	}
}
