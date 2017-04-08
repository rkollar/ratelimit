package ratelimit

import (
	"io"
)

const chunkSize = 32 * 1024

type Writer struct {
	w      io.Writer
	bucket *Bucket
}

func NewWriter(w io.Writer, rate int64) *Writer {
	return &Writer{
		w:      w,
		bucket: NewBucketWithRate(rate, capacityFromRate(rate)),
	}
}

func (self *Writer) Write(p []byte) (n int, err error) {
	var n_ int
	var chunk []byte

	for pos := 0; pos < len(p); pos += chunkSize {
		chunk = p[pos:min(len(p), pos+chunkSize)]

		self.bucket.Wait(int64(len(chunk)))

		n_, err = self.w.Write(chunk)
		n += n_
		if err != nil {
			return
		}
	}
	return
}

func (self *Writer) Close() error {
	closer, ok := self.w.(io.Closer)
	if ok {
		return closer.Close()
	}
	return nil
}

func (self *Writer) SetLimit(rate int64) {
	self.bucket.Set(rate, capacityFromRate(rate))
}
func (self *Writer) FillBucket() {
	self.bucket.Fill()
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func capacityFromRate(rate int64) int64 {
	return rate * 3
}
