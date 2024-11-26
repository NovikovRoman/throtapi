# ThrotAPI

> `Throttle API` API limit self-monitoring library.

## Getting Started

Download the `Throttle API` library:

```shell
go get github.com/NovikovRoman/throtapi
```

## Config

```go
type Config struct {
    PerSec   int // limit of requests per second
    PerMin   int // limit of requests per minute
    PerHour  int // limit of requests per hour
    PerDay   int // limit of requests per day
    PerMonth int // limit of requests per month
}
```

Limit <1 is ignored.

## Methods

- `New(cfg throtapi.Config) *Throtapi`
- `IsFree() bool` - returns true if the API limits are not reached.
- `IsBusy() bool` - returns true if API limits are reached.
- `Limits() map[timeUnit]TimeUnitParam` - returns the current API limits.

## Usage Example

```go
package main

import (
    "context"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/NovikovRoman/throtapi"
)

func main() {
    tapi := throtapi.New(throtapi.Config{
        PerSec: 3, // limit of 3 requests per second
        PerDay: 3000, // limit of 3000 requests per day
        PerMonth: 30000, // limit of 30000 requests per month
    })

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    interrupt := make(chan os.Signal, 1)
    defer close(interrupt)
    signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

    go worker(ctx, tapi)
    <-interrupt
}

func worker(ctx context.Context, tapi *throtapi.Throtapi) {
    for {
        select {
        case <-ctx.Done():
            return

        default:
            if tapi.IsFree() {
                // accessing the API
                // …
                continue
            }

        // waiting
        time.Sleep(time.Second / 100)
        }
    }
}
```
