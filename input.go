package hessian

import (
	"io"
	"reflect"
)

// Input stream for hessian response
type Input interface {
	Read() ([]interface{}, error)
	SetReader(io.Reader)
	SetTypeMapping(string, reflect.Type)
}
