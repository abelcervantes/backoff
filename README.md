# :stopwatch: ️backoff :stopwatch:

![Go](https://github.com/abelcervantes/backoff/workflows/Go/badge.svg)
[![Build Status](https://travis-ci.com/abelcervantes/backoff.svg?branch=master)](https://travis-ci.com/abelcervantes/backoff)
[![codecov](https://codecov.io/gh/abelcervantes/backoff/branch/master/graph/badge.svg)](https://codecov.io/gh/abelcervantes/backoff)

A simple truncated backoff algorithm implementation made with go!

## Installation
```shell
go get github.com/abelcervantes/backoff
```

## Formula
To calculate the next duration:
((2^attempts) - 1) * slotTime 

## Usage
default: 
``` go
package main

import (
	"time"

	"github.com/abelcervantes/backoff"
)

func main() {
    b := backoff.NewDefault()
    time.Sleep(b.NextDuration())
}
```
custom: 
``` go
package main

import (
	"time"

	"github.com/abelcervantes/backoff"
)

func main() {
    b, err := backoff.New(2*time.Second, 20*time.Second, 400*time.Millisecond, false, 10)
    if err != nil {
        log.Fatal(err)
    }	
    time.Sleep(b.NextDuration())
}
```

- [Examples][]

[Examples]: ./_examples

## References
- https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
- https://en.wikipedia.org/wiki/Exponential_backoff
