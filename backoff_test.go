package backoff_test

import (
	"testing"
	"time"

	"github.com/abelcervantes/backoff"
)

type newInput struct {
	min         time.Duration
	max         time.Duration
	jitter      bool
	slotTime    time.Duration
	maxAttempts uint
}

type newTestTableItem struct {
	inputs newInput
	err    error
}

func NewTestTable() []newTestTableItem {
	return []newTestTableItem{
		{
			inputs: newInput{
				min:         5 * time.Second,
				max:         2 * time.Second,
				jitter:      false,
				slotTime:    100 * time.Millisecond,
				maxAttempts: 10,
			},
			err: backoff.ErrMaxDurationMustBeGreater,
		},
		{
			inputs: newInput{
				min:         0 * time.Second,
				max:         10 * time.Second,
				jitter:      false,
				slotTime:    -10 * time.Millisecond,
				maxAttempts: 10,
			},
			err: backoff.ErrInvalidSlotTime,
		},
		{
			inputs: newInput{
				min:         0 * time.Second,
				max:         10 * time.Second,
				jitter:      false,
				slotTime:    100 * time.Millisecond,
				maxAttempts: 0,
			},
			err: backoff.ErrInvalidMaxAttempts,
		},
		{
			inputs: newInput{
				min:         0 * time.Second,
				max:         10 * time.Second,
				jitter:      false,
				slotTime:    100 * time.Millisecond,
				maxAttempts: 10,
			},
			err: nil,
		},
	}
}

func TestNew(t *testing.T) {
	for _, expected := range NewTestTable() {
		inp := expected.inputs
		b, err := backoff.New(inp.min, inp.max, inp.slotTime, inp.jitter, inp.maxAttempts)
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
				t.Fatalf("expected to not reach max attempts without calling next duration at least once, got: %v",
					b.HasReachedMaxAttempts())
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
		t.Fatalf("expected to not reach max attempts without calling next duration at least once, got: %v",
			b.HasReachedMaxAttempts())
	}
}

type jitterTestItem struct {
	min     time.Duration
	max     time.Duration
	attempt uint
	reachedMax bool
}

type noJitterTestItem struct {
	duration time.Duration
	attempt  uint
	reachedMax bool
}

func defaultTestTable() []jitterTestItem {
	return []jitterTestItem{
		{
			min:        0,
			max:        100 * time.Millisecond,
			attempt:    1,
			reachedMax: false,
		},
		{
			min:     0,
			max:     300 * time.Millisecond,
			attempt: 2,
			reachedMax: false,
		},
		{
			min:     0,
			max:     700 * time.Millisecond,
			attempt: 3,
			reachedMax: false,
		},
		{
			min:     0,
			max:     1500 * time.Millisecond,
			attempt: 4,
			reachedMax: false,
		},
		{
			min:     0,
			max:     3100 * time.Millisecond,
			attempt: 5,
			reachedMax: false,
		},
		{
			min:     0,
			max:     6300 * time.Millisecond,
			attempt: 6,
			reachedMax: false,
		},
		{
			min:     0,
			max:     12700 * time.Millisecond,
			attempt: 7,
			reachedMax: false,
		},
		{
			min:     0,
			max:     25500 * time.Millisecond,
			attempt: 8,
			reachedMax: false,
		},
		{
			min:     0,
			max:     51100 * time.Millisecond,
			attempt: 9,
			reachedMax: false,
		},
		{
			min:     0,
			max:     102300 * time.Millisecond,
			attempt: 10,
			reachedMax: true,
		},
		{
			min:     0,
			max:     102300 * time.Millisecond,
			attempt: 11,
			reachedMax: true,
		},
	}
}

type parametrizedTestItem struct {
	newInput        newInput
	noJitterOutputs []noJitterTestItem
	jitterOutputs   []jitterTestItem
}

func parametrizedTestTable() []parametrizedTestItem {
	return []parametrizedTestItem{
		{
			newInput: newInput{
				min:         0 * time.Millisecond,
				max:         0 * time.Millisecond,
				jitter:      false,
				slotTime:    100 * time.Millisecond,
				maxAttempts: 2,
			},
			noJitterOutputs: []noJitterTestItem{
				{
					duration: 0 * time.Millisecond,
					attempt:  1,
					reachedMax: false,
				},
				{
					duration: 0 * time.Millisecond,
					attempt:  2,
					reachedMax: true,
				},
				{
					duration: 0 * time.Millisecond,
					attempt:  3,
					reachedMax: true,
				},
			},
		},
		{
			newInput: newInput{
				min:         20 * time.Second,
				max:         60 * time.Second,
				jitter:      false,
				slotTime:    100 * time.Millisecond,
				maxAttempts: 2,
			},
			noJitterOutputs: []noJitterTestItem{
				{
					duration: 20 * time.Second,
					attempt:  1,
					reachedMax: false,
				},
				{
					duration: 20 * time.Second,
					attempt:  2,
					reachedMax: true,
				},
				{
					duration: 20 * time.Second,
					attempt:  3,
					reachedMax: true,
				},
			},
		},
		{
			newInput: newInput{
				min:         0 * time.Second,
				max:         20 * time.Second,
				jitter:      true,
				slotTime:    100 * time.Millisecond,
				maxAttempts: 2,
			},
			jitterOutputs: []jitterTestItem{
				{
					min:     0,
					max:     20 * time.Second,
					attempt: 1,
					reachedMax: false,
				},
				{
					min:     0,
					max:     20 * time.Second,
					attempt: 2,
					reachedMax: true,
				},
				{
					min:     0,
					max:     20 * time.Second,
					attempt: 3,
					reachedMax: true,
				},
				{
					min:     0,
					max:     20 * time.Second,
					attempt: 4,
					reachedMax: true,
				},
			},
		},
		{
			newInput: newInput{
				min:         0 * time.Second,
				max:         0 * time.Second,
				jitter:      true,
				slotTime:    100 * time.Millisecond,
				maxAttempts: 2,
			},
			jitterOutputs: []jitterTestItem{
				{
					min:     0,
					max:     0,
					attempt: 1,
					reachedMax: false,
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
			t.Fatalf("expected to not reach max attempts without calling next duration at least once, got: %v",
				b.HasReachedMaxAttempts())
		}

		for _, i := range defaultTestTable() {
			d := b.NextDuration()
			if b.PerformedAttempts() != i.attempt {
				t.Fatalf("expected attempt: %d , got: %d", i.attempt, b.PerformedAttempts())
			}

			if d > i.max || d < i.min {
				t.Fatalf("expected duration between: [%d, %d] got: %d", i.min, i.max, d)
			}

			if i.reachedMax != b.HasReachedMaxAttempts() {
				t.Fatalf("expected reached max attempts: %v , got: %v", i.reachedMax, b.HasReachedMaxAttempts())
			}
		}
	})

	t.Run("parametrized next duration", func(t *testing.T) {
		for _, i := range parametrizedTestTable() {
			input := i.newInput
			b, err := backoff.New(input.min, input.max, input.slotTime, input.jitter, input.maxAttempts)
			if err != nil {
				t.Fatalf("unexpected error, got: %s", err)
			}

			if b == nil {
				t.Fatalf("expected a non nil backoff")
			}

			if b.PerformedAttempts() != 0 {
				t.Fatalf("expected 0 attempts, got: %d", b.PerformedAttempts())
			}

			if b.HasReachedMaxAttempts() {
				t.Fatalf("expected to not reach max attempts without calling next duration at least once, got: %v",
					b.HasReachedMaxAttempts())
			}

			if input.jitter {
				for _, o := range i.jitterOutputs {
					d := b.NextDuration()
					if o.attempt != b.PerformedAttempts() {
						t.Fatalf("expected attempt: %d , got: %d", o.attempt, b.PerformedAttempts())
					}

					if d < o.min || d > o.max {
						t.Fatalf("expected duration between: [%d, %d] got: %d", o.min, o.max, d)
					}

					if o.reachedMax != b.HasReachedMaxAttempts() {
						t.Fatalf("expected reached max attempts: %v , got: %v", o.reachedMax, b.HasReachedMaxAttempts())
					}
				}
			} else {
				for _, o := range i.noJitterOutputs {
					d := b.NextDuration()
					if o.attempt != b.PerformedAttempts() {
						t.Fatalf("expected attempt: %d , got: %d", o.attempt, b.PerformedAttempts())
					}

					if o.duration != d {
						t.Fatalf("expected d: %d , got: %d", o.duration, d)
					}

					if o.reachedMax != b.HasReachedMaxAttempts() {
						t.Fatalf("expected reached max attempts: %v , got: %v", o.reachedMax, b.HasReachedMaxAttempts())
					}
				}
			}
		}
	})
}
