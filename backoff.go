package main

import (
	"errors"
	"math"
	"math/rand"
	"time"
)

type Backoff struct {
	min      uint
	max      uint
	unit     time.Duration
	jitter   bool
	attempts uint
}

var ErrMaxDurationMustBeGreater = errors.New("max duration must be greater than min duration")

// NewBackoff creates a fully parametrized Backoff
func NewBackoff(min, max uint, unit time.Duration, jitter bool, attempts uint) (*Backoff, error) {
	if max < min {
		return nil, ErrMaxDurationMustBeGreater
	}

	return &Backoff{
		min:      min,
		max:      max,
		unit:     unit,
		jitter:   jitter,
		attempts: attempts,
	}, nil
}

// NewDefaultBackoff creates a Backoff with default configuration
func NewDefaultBackoff() Backoff {
	return Backoff{
		min:      10,
		max:      120,
		unit:     time.Second,
		jitter:   true,
		attempts: 20,
	}
}

// NextDuration returns the next waiting time
func (b *Backoff) NextDuration() time.Duration {
	if b.jitter {
		return b.expJitter()
	}

	return b.exp()
}

func (b Backoff) expJitter() time.Duration {
	nextD := b.calcNextDuration()

	rand.Seed(time.Now().UnixNano())

	return time.Duration(rand.Intn(nextD-b.min)+b.min) * b.unit
}

func (b Backoff) exp() time.Duration {
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
func (b *Backoff) Attempts() int {
	return b.attempts
}