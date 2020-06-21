package backoff

import (
	"errors"
	"math"
	"math/rand"
	"time"
)

type Backoff struct {
	min      int
	max      int
	unit     time.Duration
	jitter   bool
	attempts uint
}

var ErrMaxDurationMustBeGreater = errors.New("max duration must be greater than min duration")

// New creates a fully parametrized Backoff
func New(min, max uint, unit time.Duration, jitter bool) (*Backoff, error) {
	if max < min {
		return nil, ErrMaxDurationMustBeGreater
	}

	return &Backoff{
		min:      int(min),
		max:      int(max),
		unit:     unit,
		jitter:   jitter,
	}, nil
}

const (
	defaultMin int = 10
	defaultMax int = 120
	defaultUnit    = time.Second
	defaultJitter  = true
)

// NewDefault creates a Backoff with default configuration
func NewDefault() Backoff {
	return Backoff{
		min:      defaultMin,
		max:      defaultMax,
		unit:     defaultUnit,
		jitter:   defaultJitter,
	}
}

// NextDuration returns the next waiting time
func (b *Backoff) NextDuration() time.Duration {
	if b.jitter {
		return b.expJitter()
	}

	return b.exp()
}

func (b *Backoff) expJitter() time.Duration {
	nextD := b.calcNextDuration()

	rand.Seed(time.Now().UnixNano())

	return time.Duration(rand.Intn(nextD-b.min)+b.min) * b.unit
}

func (b *Backoff) exp() time.Duration {
	return time.Duration(b.calcNextDuration()) * b.unit
}

func (b *Backoff) calcNextDuration() int {
	d := int(math.Pow(2, float64(b.attempts)) * 100)

	if d < b.min {
		return b.min
	}

	if d > b.max {
		return b.max
	}

	b.attempts++

	return d
}

// Attempts returns the current number of performed attempts
func (b *Backoff) Attempts() uint {
	return b.attempts
}
