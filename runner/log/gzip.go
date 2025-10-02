package log

import (
	"compress/gzip"
	"io"
)

// GzipLog provides compressed logging using gzip compression.
// It wraps an underlying WriteCloser and compresses all written data
// on-the-fly, providing efficient storage for large log files.
type GzipLog struct {
	w  io.WriteCloser
	gz *gzip.Writer
}

func NewGzipLog(w io.WriteCloser) (*GzipLog, error) {
	var l GzipLog

	l.w = w
	l.gz = gzip.NewWriter(w)

	return &l, nil
}

func (c *GzipLog) Write(p []byte) (n int, err error) {
	n, err = c.gz.Write(p)
	if err != nil {
		return
	}
	err = c.gz.Flush()
	return
}

func (c *GzipLog) Close() error {
	err := c.gz.Close()
	if err != nil {
		return err
	}

	err = c.w.Close()
	return err
}
