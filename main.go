package main

import (
	"log"
	"runtime"
	"time"

	command "github.com/webhkp/godft/cmd"
)

func main() {
	// Print memory usage every 5 second
	go func() {
		for {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			log.Printf("\nAlloc = %v\tTotalAlloc = %v\tSys = %v\tNumGC = %v\n\n", m.Alloc/1024, m.TotalAlloc/1024, m.Sys/1024, m.NumGC)
			time.Sleep(5 * time.Second)
		}
	}()

	command.Execute()
}
