package hessian

import (
	"io"
	"reflect"
)

// Deserializer input stream for hessian response
type Deserializer interface {
	// Read parse input (io.Reader) to return value
	Read() ([]interface{}, error)

	// Reset
	Reset(io.Reader)

	// SetTypeMap
	SetTypeMap(map[string]reflect.Type)
}
