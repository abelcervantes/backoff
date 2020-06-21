package backoff_test

import (
	"testing"
	"time"

	"github.com/abelcervantes/backoff"
)

type newInput struct {
	min         uint
	max         uint
	unit        time.Duration
	jitter      bool
	slotTime    float64
	maxAttempts uint
}

type newTestTableItem struct {
	inputs newInput
	err error
}

func NewTestTable() []newTestTableItem {
	return []newTestTableItem{
		{
			inputs: newInput{
				min:         5,
				max:         2,
				unit:        time.Second,
				jitter:      false,
				slotTime:    100,
				maxAttempts: 10,
			},
			err:    backoff.ErrMaxDurationMustBeGreater,
		},
		{
			inputs: newInput{
				min:         0,
				max:         10,
				unit:        time.Second,
				jitter:      false,
				slotTime:    -10,
				maxAttempts: 10,
			},
			err:    backoff.ErrInvalidSlotTime,
		},
		{
			inputs: newInput{
				min:         0,
				max:         10,
				unit:        time.Second,
				jitter:      false,
				slotTime:    100,
				maxAttempts: 0,
			},
			err:    backoff.ErrInvalidMaxAttempts,
		},
		{
			inputs: newInput{
				min:         0,
				max:         10,
				unit:        time.Second,
				jitter:      false,
				slotTime:    100,
				maxAttempts: 10,
			},
			err:    nil,
		},
	}
}

func TestNew(t *testing.T) {
	for _, expected := range NewTestTable() {
		inp := expected.inputs
		b, err := backoff.New(inp.min, inp.max, inp.unit, inp.jitter, inp.slotTime, inp.maxAttempts)
		if err != expected.err {
			t.Fatalf("expected: %s got: %s", expected.err, err)
		}

		if err == nil {
			if b == nil {
				t.Fatalf("expected a non nil backoff")
			}

			if b.PerformedAttempts() != 0 {
				t.Fatalf("expected 0 attempts, got: %d", b.PerformedAttempts())
			}

			if b.HasReachedMaxAttempts() {
				t.Fatalf("expected to not reach max attempts without calling next duration at least once, got: %v", b.HasReachedMaxAttempts())
			}
		}
	}
}

func TestNewDefault(t *testing.T) {
	b := backoff.NewDefault()

	if b.PerformedAttempts() != 0 {
		t.Fatalf("expected 0 attempts, got: %d", b.PerformedAttempts())
	}

	if b.HasReachedMaxAttempts() {
		t.Fatalf("expected to not reach max attempts without calling next duration at least once, got: %v", b.HasReachedMaxAttempts())
	}
}

type nextDurationDefaultTestItem struct {
	min     time.Duration
	max     time.Duration
	attempt uint
}

func nextDurationDefaultTable() []nextDurationDefaultTestItem {
	return []nextDurationDefaultTestItem{
		{
			min:     0,
			max:     100 * time.Millisecond,
			attempt: 1,
		},
		{
			min:     0,
			max:     300 * time.Millisecond,
			attempt: 2,
		},
		{
			min:     0,
			max:     700 * time.Millisecond,
			attempt: 3,
		},
		{
			min:     0,
			max:     1500 * time.Millisecond,
			attempt: 4,
		},
		{
			min:     0,
			max:     3100 * time.Millisecond,
			attempt: 5,
		},
		{
			min:     0,
			max:     6300 * time.Millisecond,
			attempt: 6,
		},
		{
			min:     0,
			max:     12700 * time.Millisecond,
			attempt: 7,
		},
		{
			min:     0,
			max:     25500 * time.Millisecond,
			attempt: 8,
		},
		{
			min:     0,
			max:     51100 * time.Millisecond,
			attempt: 9,
		},
		{
			min:     0,
			max:     102300 * time.Millisecond,
			attempt: 10,
		},
		{
			min:     0,
			max:     102300 * time.Millisecond,
			attempt: 11,
		},
	}
}

type nextDurationParametrizedTestItem struct {
	newInput            newInput
	nextDurationOutputs []nextDurationDefaultTestItem
}

func nextDurationParametrizedTable() []nextDurationParametrizedTestItem {
	return []nextDurationParametrizedTestItem{
		{
			newInput: newInput{
				min:         0,
				max:         0,
				unit:        time.Millisecond,
				jitter:      false,
				slotTime:    100,
				maxAttempts: 2,
			},
			nextDurationOutputs: []nextDurationDefaultTestItem{
				{
					min:     0,
					max:     0,
					attempt: 1,
				},
				{
					min:     0,
					max:     0,
					attempt: 2,
				},
				{
					min:     0,
					max:     0,
					attempt: 3,
				},
			},
		},
		{
			newInput: newInput{
				min:         0,
				max:         20,
				unit:        time.Second,
				jitter:      true,
				slotTime:    100,
				maxAttempts: 2,
			},
			nextDurationOutputs: []nextDurationDefaultTestItem{
				{
					min:     0,
					max:     0,
					attempt: 1,
				},
				{
					min:     0,
					max:     0,
					attempt: 2,
				},
				{
					min:     0,
					max:     0,
					attempt: 3,
				},
			},
		},
	}
}

func TestBackoff_NextDuration(t *testing.T) {
	t.Parallel()

	t.Run("default next duration", func(t *testing.T) {
		b := backoff.NewDefault()

		if b.PerformedAttempts() != 0 {
			t.Fatalf("expected 0 attempts, got: %d", b.PerformedAttempts())
		}

		if b.HasReachedMaxAttempts() {
			t.Fatalf("expected to not reach max attempts without calling next duration at least once, got: %v", b.HasReachedMaxAttempts())
		}
		
		for _, i := range nextDurationDefaultTable() {
			d := b.NextDuration()
			if b.PerformedAttempts() !=  i.attempt {
				t.Fatalf("expected attempt: %d , got: %d", i.attempt, b.PerformedAttempts())
			}

			if d > i.max || d < i.min {
				t.Fatalf("expected duration between: [%d, %d] got: %d", i.min, i.max, d)
			}
		}
	})

	t.Run("parametrized next duration (no jitter)", func(t *testing.T) {
		for _, i := range nextDurationParametrizedTable() {
			input := i.newInput
			b, err := backoff.New(input.min, input.max, input.unit, input.jitter, input.slotTime, input.maxAttempts)
			if err != nil {
				t.Fatalf("expected no error, got: %s", err)
			}

			if b == nil {
				t.Fatalf("expected a non nil backoff")
			}

			if b.PerformedAttempts() != 0 {
				t.Fatalf("expected 0 attempts, got: %d", b.PerformedAttempts())
			}

			if b.HasReachedMaxAttempts() {
				t.Fatalf("expected to not reach max attempts without calling next duration at least once, got: %v", b.HasReachedMaxAttempts())
			}

			for _, o := range i.nextDurationOutputs {
				_ = b.NextDuration()
				if o.attempt != b.PerformedAttempts() {
					t.Fatalf("expected attempt: %d , got: %d", o.attempt, b.PerformedAttempts())
				}
			}
		}
	})
}