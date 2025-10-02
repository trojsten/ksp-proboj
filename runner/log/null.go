package log

// NullLog provides a no-op logging implementation that discards all data.
// Useful for disabling logging entirely.
type NullLog struct{}

func NewNullLog() *NullLog {
	return &NullLog{}
}

func (NullLog) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (NullLog) Close() error {
	return nil
}
