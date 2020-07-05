package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/abelcervantes/backoff"
)

func example() {
	// http server
	go func() {
		maxTries, tries := 5, 0

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if tries == maxTries {
				fmt.Fprint(w, "Hello from server!")
				return
			}

			tries++
			http.Error(w, "internal server error", http.StatusInternalServerError)
		})

		_ = http.ListenAndServe(":8080", nil)
	}()

	time.Sleep(1 * time.Second)

	// http client
	b := backoff.NewDefault()
	for !b.HasReachedMaxAttempts() {
		response, err := http.Get("http://localhost:8080/")
		if err != nil {
			log.Fatal(err)
		}

		if response.StatusCode == http.StatusOK {
			log.Printf("external service answered with an OK response, performed attempts: %d", b.PerformedAttempts())

			body, _ := ioutil.ReadAll(response.Body)
			log.Printf("response: %s", string(body))
			response.Body.Close()
			break
		}

		log.Printf("external service answered with an error: %d", response.StatusCode)

		response.Body.Close()

		log.Printf("sleeping...")
		time.Sleep(b.NextDuration())
	}
}

func main() {
	example()
}
