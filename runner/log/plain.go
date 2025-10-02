package log

import "io"

// PlainLog provides uncompressed logging to an underlying writer.
// It passes through all write operations without modification.
type PlainLog struct {
	w io.WriteCloser
}

func NewPlainLog(w io.WriteCloser) *PlainLog {
	return &PlainLog{w: w}
}

func (l *PlainLog) Write(p []byte) (n int, err error) {
	return l.w.Write(p)
}

func (l *PlainLog) Close() error {
	return l.w.Close()
}
