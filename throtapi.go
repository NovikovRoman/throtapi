package throtapi

import (
	"sync"
	"time"
)

const (
	Second = iota
	Minute
	Hour
	Day
	Month
)

type timeUnit int

type Config struct {
	PerSec   int
	PerMin   int
	PerHour  int
	PerDay   int
	PerMonth int
}

type Throtapi struct {
	sync.RWMutex
	limits map[timeUnit]TimeUnitParam
}

type TimeUnitParam struct {
	Timestamp   int64
	NumRequests int
	Limit       int
}

// New returns a new Throtapi instance.
func New(cfg Config) *Throtapi {
	return &Throtapi{
		limits: map[timeUnit]TimeUnitParam{
			Second: {Limit: cfg.PerSec},
			Minute: {Limit: cfg.PerMin},
			Hour:   {Limit: cfg.PerHour},
			Day:    {Limit: cfg.PerDay},
			Month:  {Limit: cfg.PerMonth},
		},
	}
}

// IsFree returns true if the API limits are not reached.
func (t *Throtapi) IsFree() bool {
	t.Lock()
	defer t.Unlock()

	if t.IsBusy() {
		return false
	}

	t.addRequest()
	return true
}

// IsBusy returns true if API limits are reached.
func (t *Throtapi) IsBusy() (busy bool) {
	for _, tu := range []timeUnit{Second, Minute, Hour, Day, Month} {
		if t.limits[tu].Limit < 1 {
			continue
		}

		if truncateTime(time.Now(), tu) == t.limits[tu].Timestamp &&
			t.limits[tu].NumRequests >= t.limits[tu].Limit {
			return true
		}
	}
	return false
}

// Limits returns the current API limits.
func (t *Throtapi) Limits() map[timeUnit]TimeUnitParam {
	t.RLock()
	defer t.RUnlock()

	mCopy := make(map[timeUnit]TimeUnitParam, len(t.limits))
	for k, v := range t.limits {
		mCopy[k] = v
	}
	return mCopy
}

// addRequest increments the number of requests for each time unit.
func (t *Throtapi) addRequest() {
	for _, tu := range []timeUnit{Second, Minute, Hour, Day, Month} {
		var tup TimeUnitParam
		var ok bool
		if tup, ok = t.limits[tu]; !ok || tup.Limit < 1 {
			continue
		}

		ts := truncateTime(time.Now(), tu)
		if ts != tup.Timestamp {
			tup.Timestamp = ts
			tup.NumRequests = 1
		} else {
			tup.NumRequests++
		}
		t.limits[tu] = tup
	}
}

// truncateTime returns the unix timestamp aligned for the given time unit.
func truncateTime(currTime time.Time, tu timeUnit) int64 {
	switch tu {
	case Second:
		return currTime.Unix()

	case Minute:
		return currTime.Truncate(time.Minute).Unix()

	case Hour:
		return currTime.Truncate(time.Hour).Unix()

	case Day:
		return currTime.Truncate(24 * time.Hour).Unix()

	default: //case Month:
		t := currTime
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()).Unix()
	}
}
