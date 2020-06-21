package backoff

import (
	"errors"
	"math"
	"math/rand"
	"time"
)

type Backoff struct {
	min               *int
	max               *int
	unit              time.Duration
	jitter            bool
	slotTime          float64
	attempts          uint
	maxAttempts       uint
	performedAttempts uint
}

var ErrMaxDurationMustBeGreater = errors.New("max duration must be greater than min duration")
var ErrInvalidSlotTime = errors.New("slot time cannot be negative")
var ErrInvalidMaxAttempts = errors.New("max attempt cannot be 0")

// New creates a fully parametrized Backoff
func New(min, max uint, unit time.Duration, jitter bool, slotTime float64, maxAttempts uint) (*Backoff, error) {
	if max < min {
		return nil, ErrMaxDurationMustBeGreater
	}

	if slotTime < 0 {
		return nil, ErrInvalidSlotTime
	}

	if maxAttempts == 0 {
		return nil, ErrInvalidMaxAttempts
	}

	nMin, nMax := int(min), int(max)
	return &Backoff{
		min:         &nMin,
		max:         &nMax,
		unit:        unit,
		jitter:      jitter,
		slotTime:    slotTime,
		maxAttempts: maxAttempts,
	}, nil
}

const (
	defaultMin         int     = 0
	defaultUnit                = time.Millisecond
	defaultJitter              = true
	defaultSlotTime    float64 = 100
	defaultMaxAttempts uint    = 10
)

// NewDefault creates a Backoff with default configuration
func NewDefault() Backoff {
	dMin := defaultMin
	return Backoff{
		min:         &dMin,
		unit:        defaultUnit,
		jitter:      defaultJitter,
		slotTime:    defaultSlotTime,
		maxAttempts: defaultMaxAttempts,
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
	nd := b.calcNextDuration()
	if nd == 0 {
		return time.Duration(nd)
	}

	rand.Seed(time.Now().UnixNano())

	x := nd - *b.min
	if x > 0 {
		x = rand.Intn(x)
	}

	return time.Duration(x + *b.min) * b.unit
}

func (b *Backoff) exp() time.Duration {
	return time.Duration(b.calcNextDuration()) * b.unit
}

func (b *Backoff) calcNextDuration() int {
	b.incAttempts()

	d := int((math.Pow(2, float64(b.attempts)) - 1) * b.slotTime)

	if b.min != nil && d < *b.min {
		return *b.min
	}

	if b.max != nil && d > *b.max {
		return *b.max
	}

	return d
}

func (b *Backoff) incAttempts() {
	if b.attempts < b.maxAttempts {
		b.attempts++
	}

	b.performedAttempts++
}

// HasReachedMaxAttempts returns whether the backoff has reached the max number of attempts
func (b *Backoff) HasReachedMaxAttempts() bool {
	return b.attempts == b.maxAttempts
}

// PerformedAttempts returns the current number of performed attempts
func (b *Backoff) PerformedAttempts() uint {
	return b.performedAttempts
}
