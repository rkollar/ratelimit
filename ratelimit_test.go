package ratelimit

import (
	"fmt"
	"testing"
	"time"
)

func ASSERT(tb testing.TB, condition bool, msg string, v ...interface{}) {
	tb.Helper()

	if !condition {
		tb.Fatalf(fmt.Sprintf(msg, v...))
	}
}

func TestFiles(t *testing.T) {
	t.Run("refill rate", func(t *testing.T) {
		tf := func(tb testing.TB, tps int64, dur time.Duration) {
			tb.Helper()

			var now = time.Now()
			var bucket = NewBucket(tps, tps*100000)
			bucket.avail = 0

			bucket.refill(now.Add(dur))
			var actualRate = bucket.avail / int64(dur.Seconds())
			ASSERT(tb, actualRate == tps, "actualRate(%d) != rate(%d)", actualRate, tps)
		}
		tf(t, 10000, 5*time.Second)
		tf(t, 10, 10*time.Second)
		tf(t, 1, 30*time.Second)
	})
}
