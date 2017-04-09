package ratelimit

import (
	"io"
)

type Writer struct {
	w            io.Writer
	maxChunkSize int
	bucket       *Bucket
}

func NewWriter(w io.Writer, maxChunkSize int, rate int64, capacity int64) *Writer {
	return &Writer{
		w:            w,
		maxChunkSize: maxChunkSize,
		bucket:       NewBucket(rate, capacity),
	}
}

func (self *Writer) Write(p []byte) (n int, err error) {
	var n_ int
	var chunk []byte

	for pos := 0; pos < len(p); pos += self.maxChunkSize {
		chunk = p[pos:min(len(p), pos+self.maxChunkSize)]

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

func (self *Writer) SetLimits(rate int64, capacity int64) {
	self.bucket.Set(rate, capacity)
}

func (self *Writer) SetMaxChunkSize(s int) {
	self.maxChunkSize = s
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
