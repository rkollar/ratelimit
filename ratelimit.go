package ratelimit

import (
	"math"
	"time"
)

const resolution = 1 * time.Millisecond

type Bucket struct {
	ts       time.Time
	rate     int64
	capacity int64
	avail    int64
}

func NewBucket(rate int64, capacity int64) *Bucket {
	b := &Bucket{
		ts:    time.Now().Round(resolution),
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

	// convert rate from tokens per second to tokens per resolution unit
	conv := float64(time.Second) / float64(resolution)
	self.rate = int64(math.Ceil(float64(rate) / conv))

	self.capacity = capacity
	if self.avail > self.capacity {
		self.avail = self.capacity
	}
}

func (self *Bucket) Wait(count int64) {
	if count <= self.avail {
		self.avail -= count
		return
	}

	// refill the bucket
	since := time.Since(self.ts)
	diff := int64(since / resolution)
	self.avail += diff * self.rate
	if self.avail > self.capacity {
		self.avail = self.capacity
	}
	self.ts = self.ts.Add(time.Duration(diff) * resolution)

	// re-check
	if count <= self.avail {
		self.avail -= count
		return
	}

	// we need to wait
	count -= self.avail
	self.avail = 0

	dur := time.Duration(math.Ceil(float64(count)/float64(self.rate))) * resolution
	if dur < resolution {
		dur = resolution
	}
	time.Sleep(dur)
	self.Wait(count)
}

func (self *Bucket) Fill() {
	self.avail = self.capacity
}
