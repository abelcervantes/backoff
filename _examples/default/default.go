package main

import (
	"log"
	"time"

	"github.com/abelcervantes/backoff"
)

func example1() {
	b := backoff.NewDefault()
	for i := 0; i < 15; i++ {
		if b.HasReachedMaxAttempts() {
			log.Printf("reached max number of attempts")
		}

		duration := b.NextDuration()
		log.Printf("next duration: %.3f\tattempts performed: %d",
			float64(duration)/float64(time.Second), b.PerformedAttempts())
		log.Printf("sleeping...")
		time.Sleep(duration)
	}
}

func example2() {
	b := backoff.NewDefault()
	for i := 0; i < 15; i++ {
		time.Sleep(b.NextDuration())
	}
}


func main() {
	example1()
	//example2()
}
