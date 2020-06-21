package backoff_test

import (
	"testing"
	"time"

	"github.com/abelcervantes/backoff"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("returns max duration error", func(t *testing.T) {
		b, err := backoff.New(5, 2, time.Second, true)
		if err != backoff.ErrMaxDurationMustBeGreater {
			t.Fatal("should return 'ErrMaxDurationMustBeGreater' error")
		}

		if b != nil {
			t.Fatal("should return a nil Backoff")
		}
	})

	t.Run("initial number of attempts must be 0", func(t *testing.T) {
		b, err := backoff.New(5, 10, time.Second, true)
		if err != nil {
			t.Fatal("unexpected error")
		}

		if b == nil {
			t.Fatal("should return a non nil Backoff")
		}

		if b.Attempts() != 0 {
			t.Fatalf("expected 0 attempts, got: %d", b.Attempts())
		}
	})
}
