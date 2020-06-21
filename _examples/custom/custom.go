package main

import (
	"log"
	"time"

	"github.com/abelcervantes/backoff"
)

func main() {
	b, err := backoff.New(2*time.Second, 20*time.Second, 400*time.Millisecond, false, 10)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		d := b.NextDuration()
		log.Printf("next duration: %d", d/time.Second)
		log.Printf("sleeping...")
		time.Sleep(d)
	}
}
