package throtapi

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testCfg = Config{
		PerSec:   1,
		PerMin:   2,
		PerHour:  3,
		PerDay:   4,
		PerMonth: 5,
	}
)

func TestNew(t *testing.T) {
	tapi := New(testCfg)
	assert.Equal(t, 1, tapi.limits[Second].Limit)
	assert.Equal(t, 0, tapi.limits[Second].NumRequests)
	assert.Equal(t, int64(0), tapi.limits[Second].Timestamp)

	assert.Equal(t, 2, tapi.limits[Minute].Limit)
	assert.Equal(t, 0, tapi.limits[Minute].NumRequests)
	assert.Equal(t, int64(0), tapi.limits[Minute].Timestamp)

	assert.Equal(t, 3, tapi.limits[Hour].Limit)
	assert.Equal(t, 0, tapi.limits[Hour].NumRequests)
	assert.Equal(t, int64(0), tapi.limits[Hour].Timestamp)

	assert.Equal(t, 4, tapi.limits[Day].Limit)
	assert.Equal(t, 0, tapi.limits[Day].NumRequests)
	assert.Equal(t, int64(0), tapi.limits[Day].Timestamp)

	assert.Equal(t, 5, tapi.limits[Month].Limit)
	assert.Equal(t, 0, tapi.limits[Month].NumRequests)
	assert.Equal(t, int64(0), tapi.limits[Month].Timestamp)

	assert.Equal(t, tapi.Limits(), map[timeUnit]TimeUnitParam{
		Second: {
			Limit:       1,
			NumRequests: 0,
			Timestamp:   0,
		},
		Minute: {
			Limit:       2,
			NumRequests: 0,
			Timestamp:   0,
		},
		Hour: {
			Limit:       3,
			NumRequests: 0,
			Timestamp:   0,
		},
		Day: {
			Limit:       4,
			NumRequests: 0,
			Timestamp:   0,
		},
		Month: {
			Limit:       5,
			NumRequests: 0,
			Timestamp:   0,
		},
	})
}

func TestStatus(t *testing.T) {
	cfg := Config{
		PerSec:   3,
		PerMin:   4,
		PerHour:  5,
		PerDay:   6,
		PerMonth: 7,
	}

	tapi := New(cfg)
	for i := 0; i < 3; i++ {
		assert.True(t, tapi.IsFree())
	}
	assert.False(t, tapi.IsFree())
	assert.Equal(t, 3, tapi.limits[Second].NumRequests)
	assert.Equal(t, 3, tapi.limits[Minute].NumRequests)
	assert.Equal(t, 3, tapi.limits[Hour].NumRequests)
	assert.Equal(t, 3, tapi.limits[Day].NumRequests)
	assert.Equal(t, 3, tapi.limits[Month].NumRequests)
	tapi.limits[Second] = TimeUnitParam{
		Limit:       3,
		NumRequests: 0,
		Timestamp:   0,
	}

	assert.True(t, tapi.IsFree())
	assert.False(t, tapi.IsFree())

	assert.Equal(t, 1, tapi.limits[Second].NumRequests)
	assert.Equal(t, 4, tapi.limits[Minute].NumRequests)
	assert.Equal(t, 4, tapi.limits[Hour].NumRequests)
	assert.Equal(t, 4, tapi.limits[Day].NumRequests)
	assert.Equal(t, 4, tapi.limits[Month].NumRequests)

	tapi.limits[Second] = TimeUnitParam{
		Limit:       3,
		NumRequests: 0,
		Timestamp:   0,
	}
	tapi.limits[Minute] = TimeUnitParam{
		Limit:       4,
		NumRequests: 0,
		Timestamp:   0,
	}

	assert.True(t, tapi.IsFree())
	assert.False(t, tapi.IsFree())
	assert.Equal(t, 1, tapi.limits[Second].NumRequests)
	assert.Equal(t, 1, tapi.limits[Minute].NumRequests)
	assert.Equal(t, 5, tapi.limits[Hour].NumRequests)
	assert.Equal(t, 5, tapi.limits[Day].NumRequests)
	assert.Equal(t, 5, tapi.limits[Month].NumRequests)

	tapi.limits[Second] = TimeUnitParam{
		Limit:       3,
		NumRequests: 0,
		Timestamp:   0,
	}
	tapi.limits[Minute] = TimeUnitParam{
		Limit:       4,
		NumRequests: 0,
		Timestamp:   0,
	}
	tapi.limits[Hour] = TimeUnitParam{
		Limit:       5,
		NumRequests: 0,
		Timestamp:   0,
	}

	assert.True(t, tapi.IsFree())
	assert.False(t, tapi.IsFree())
	assert.Equal(t, 1, tapi.limits[Second].NumRequests)
	assert.Equal(t, 1, tapi.limits[Minute].NumRequests)
	assert.Equal(t, 1, tapi.limits[Hour].NumRequests)
	assert.Equal(t, 6, tapi.limits[Day].NumRequests)
	assert.Equal(t, 6, tapi.limits[Month].NumRequests)

	tapi.limits[Second] = TimeUnitParam{
		Limit:       3,
		NumRequests: 0,
		Timestamp:   0,
	}
	tapi.limits[Minute] = TimeUnitParam{
		Limit:       4,
		NumRequests: 0,
		Timestamp:   0,
	}
	tapi.limits[Hour] = TimeUnitParam{
		Limit:       5,
		NumRequests: 0,
		Timestamp:   0,
	}
	tapi.limits[Day] = TimeUnitParam{
		Limit:       6,
		NumRequests: 0,
		Timestamp:   0,
	}

	assert.True(t, tapi.IsFree())
	assert.False(t, tapi.IsFree())
	assert.Equal(t, 1, tapi.limits[Second].NumRequests)
	assert.Equal(t, 1, tapi.limits[Minute].NumRequests)
	assert.Equal(t, 1, tapi.limits[Hour].NumRequests)
	assert.Equal(t, 1, tapi.limits[Day].NumRequests)
	assert.Equal(t, 7, tapi.limits[Month].NumRequests)

	tapi.limits[Month] = TimeUnitParam{
		Limit:       tapi.limits[Month].Limit,
		NumRequests: tapi.limits[Month].NumRequests,
		Timestamp:   tapi.limits[Month].Timestamp - 1,
	}

	assert.True(t, tapi.IsFree())
	assert.True(t, tapi.IsFree())
	assert.False(t, tapi.IsFree())
	assert.Equal(t, 3, tapi.limits[Second].NumRequests)
	assert.Equal(t, 3, tapi.limits[Minute].NumRequests)
	assert.Equal(t, 3, tapi.limits[Hour].NumRequests)
	assert.Equal(t, 3, tapi.limits[Day].NumRequests)
	assert.Equal(t, 2, tapi.limits[Month].NumRequests)

	cfg = Config{
		PerSec:   2,
		PerMin:   0,
		PerHour:  0,
		PerDay:   0,
		PerMonth: 0,
	}
	tapi = New(cfg)
	assert.True(t, tapi.IsFree())
	assert.True(t, tapi.IsFree())
	assert.False(t, tapi.IsFree())
	assert.Equal(t, 2, tapi.limits[Second].NumRequests)
	assert.Equal(t, 0, tapi.limits[Minute].NumRequests)
	assert.Equal(t, 0, tapi.limits[Hour].NumRequests)
	assert.Equal(t, 0, tapi.limits[Day].NumRequests)
	assert.Equal(t, 0, tapi.limits[Month].NumRequests)

	tapi.limits[Second] = TimeUnitParam{
		Limit:       tapi.limits[Second].Limit,
		NumRequests: tapi.limits[Second].NumRequests,
		Timestamp:   tapi.limits[Second].Timestamp - 1,
	}
	assert.True(t, tapi.IsFree())
	assert.True(t, tapi.IsFree())
	assert.False(t, tapi.IsFree())
	assert.Equal(t, 2, tapi.limits[Second].NumRequests)
	assert.Equal(t, 0, tapi.limits[Minute].NumRequests)
	assert.Equal(t, 0, tapi.limits[Hour].NumRequests)
	assert.Equal(t, 0, tapi.limits[Day].NumRequests)
	assert.Equal(t, 0, tapi.limits[Month].NumRequests)
}

func Test_truncateTime(t *testing.T) {
	tests := []struct {
		time time.Time
		want map[timeUnit]time.Time
	}{
		{
			time: time.Date(2005, 5, 20, 18, 5, 4, 444, time.UTC),
			want: map[timeUnit]time.Time{
				Second: time.Date(2005, 5, 20, 18, 5, 4, 0, time.UTC),
				Minute: time.Date(2005, 5, 20, 18, 5, 0, 0, time.UTC),
				Hour:   time.Date(2005, 5, 20, 18, 0, 0, 0, time.UTC),
				Day:    time.Date(2005, 5, 20, 0, 0, 0, 0, time.UTC),
				Month:  time.Date(2005, 5, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			time: time.Date(2012, 5, 27, 2, 45, 44, 444, time.UTC),
			want: map[timeUnit]time.Time{
				Second: time.Date(2012, 5, 27, 2, 45, 44, 0, time.UTC),
				Minute: time.Date(2012, 5, 27, 2, 45, 0, 0, time.UTC),
				Hour:   time.Date(2012, 5, 27, 2, 0, 0, 0, time.UTC),
				Day:    time.Date(2012, 5, 27, 0, 0, 0, 0, time.UTC),
				Month:  time.Date(2012, 5, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.time.String(), func(t *testing.T) {
			for tu, tt := range test.want {
				res := truncateTime(test.time, tu)
				if tt.Unix() != res {
					tuStr := "month"
					switch tu {
					case Second:
						tuStr = "second"
					case Minute:
						tuStr = "minute"
					case Hour:
						tuStr = "hour"
					case Day:
						tuStr = "day"
					}
					t.Errorf("Not equal %v %s %v", test.time, tuStr, res)
				}
			}
		})
	}
}
