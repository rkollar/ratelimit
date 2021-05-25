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
	var ns int
	var s []byte

	for pos := 0; pos < len(p); pos += self.maxChunkSize {
		s = p[pos:min(len(p), pos+self.maxChunkSize)]

		self.bucket.Wait(int64(len(s)))

		ns, err = self.w.Write(s)
		n += ns
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
	if s <= 0 {
		panic("max chunk size <= 0")
	}
	self.maxChunkSize = s
}

func (self *Writer) BucketFillMax() {
	self.bucket.FillMax()
}
func (self *Writer) BucketFill(val int64) {
	self.bucket.Fill(val)
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
