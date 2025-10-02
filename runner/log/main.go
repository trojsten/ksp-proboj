// Package log provides logging infrastructure for the Proboj runner.
// It supports multiple log formats including plain text, gzip compression,
// and null logging (disabled). The package abstracts log writing to allow
// flexible logging strategies for different use cases.
package log

import (
	"io"
)

// Log defines the interface for log writers in the Proboj system.
type Log io.WriteCloser
