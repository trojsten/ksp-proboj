package log

import "io"

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
