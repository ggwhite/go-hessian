package hessian

import (
	"io"
	"reflect"
)

// Serializer output stream for hessian requests
type Serializer interface {
	// Call Writes a complete method call.
	Call(string, ...interface{}) error

	// Writer get writer
	Writer() io.Writer

	// Reader get reader
	Reader() io.Reader

	// Flush clean writer
	Flush()

	// SetTypeMap
	SetTypeMap(map[string]reflect.Type)
}
