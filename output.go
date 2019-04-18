package hessian

import (
	"io"
	"reflect"
)

// Output stream for hessian requests
type Output interface {
	// Call Writes a complete method call.
	Call(string, ...interface{}) error

	// Writer get writer
	Writer() io.Writer

	// Reader get reader
	Reader() io.Reader

	// Flush clean writer
	Flush()

	// SetTypeMapping set type mapping
	SetTypeMapping(string, reflect.Type)
}
