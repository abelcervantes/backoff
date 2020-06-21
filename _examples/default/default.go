package main

import (
	"log"
	"time"

	"github.com/abelcervantes/backoff"
)

func main() {
	b := backoff.NewDefault()
	for true {
		/*if b.HasReachedMaxAttempts() {
			log.Fatalf("reached max number of attempts")
		}*/

		duration := b.NextDuration()
		log.Printf("next duration: %.3f\tperformed: %d",
			float64(duration)/float64(time.Second), b.PerformedAttempts())
		log.Printf("sleeping...")
		time.Sleep(duration)
	}
}
