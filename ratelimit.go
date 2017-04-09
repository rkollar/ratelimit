package ratelimit

import (
	"time"
)

type Bucket struct {
	ts       time.Time
	rate     int64
	capacity int64
	avail    int64
}

func NewBucket(rate int64, capacity int64) *Bucket {
	b := &Bucket{
		ts:    time.Now().Round(1 * time.Millisecond),
		avail: capacity,
	}
	b.Set(rate, capacity)

	return b
}

func (self *Bucket) Set(rate int64, capacity int64) {
	if rate <= 0 {
		panic("rate <= 0")
	}
	if capacity <= 0 {
		panic("capacity <= 0")
	}
	if rate < 1000 {
		rate = 1000
	}
	self.rate = int64(float64(rate) / 1000)
	self.capacity = capacity
}

func (self *Bucket) Wait(count int64) {
	if count <= self.avail {
		self.avail -= count
		return
	}

	// refill the bucket
	diff := int64(time.Since(self.ts).Nanoseconds() / 1000000)
	self.avail += diff * self.rate
	if self.avail > self.capacity {
		self.avail = self.capacity
	}
	self.ts = self.ts.Add(time.Duration(diff) * time.Millisecond)

	// re-check
	if count <= self.avail {
		self.avail -= count
		return
	}

	// we need to wait
	count -= self.avail
	self.avail = 0
	dur := time.Duration((float64(count) / float64(self.rate))) * time.Millisecond
	if dur < 1*time.Millisecond {
		dur = 1 * time.Millisecond
	}
	time.Sleep(dur)
	self.Wait(count)
}

func (self *Bucket) Fill() {
	self.avail = self.capacity
}
