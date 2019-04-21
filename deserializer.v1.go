package hessian

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"reflect"
)

// DeserializerV1 input stream for hessian 1.0 response
type DeserializerV1 struct {
	version int
	r       io.Reader
	typeMap map[string]reflect.Type
}

// ReadAt Read object from given bytes and begin index.
func (i *DeserializerV1) ReadAt(p []byte, begin int) ([]interface{}, int, error) {
	var ans []interface{}
	var j = begin
	var err error

	for ; j < len(p); j++ {
		ch := p[j]
		switch ch {
		default:
		case 'c':
		case 'r':
		case 'f':
			var args []interface{}
			args, j, err = i.ReadAt(p, j+1)
			if err != nil {
				return nil, j, err
			}
			return nil, j, fmt.Errorf("Call Service Error, Exception: %s, Message: %s", args[1].(string), args[3].(string))
		case 'z':
			return ans, j, nil
		case 'T':
			ans = append(ans, true)
		case 'F':
			ans = append(ans, false)
		case 'N':
			ans = append(ans, nil)
		case 'B':
			// B b16 b8 byte-value
			var val []byte
			val, j, err = i.ReadBytesAt(p, j+1)
			if err != nil {
				return nil, j, err
			}
			ans = append(ans, val)
		case 'S':
			// S b16 b8 string-value
			var val string
			val, j, err = i.ReadStringAt(p, j+1)
			if err != nil {
				return nil, j, err
			}
			ans = append(ans, val)
		case 'I':
			// I b32 b24 b16 b8
			var val int32
			val, j, err = i.ReadInt32At(p, j+1)
			if err != nil {
				return nil, j, err
			}
			ans = append(ans, val)
		case 'L':
			// L b64 b56 b48 b40 b32 b24 b16 b8
			var val int64
			val, j, err = i.ReadInt64At(p, j+1)
			if err != nil {
				return nil, j, err
			}
			ans = append(ans, val)
		case 'D':
			// D b64 b56 b48 b40 b32 b24 b16 b8
			var val float64
			val, j, err = i.ReadFloat64At(p, j+1)
			if err != nil {
				return nil, j, err
			}
			ans = append(ans, val)
		case 'M':
			if p[j+1] == 't' {
				j++
				var pkg string
				var m = make(map[interface{}]interface{})
				pkg, j, err = i.ReadStringAt(p, j+1)
				if err != nil {
					return nil, j, err
				}

				// parse 'Mt' arguments
				m, j, err = i.ReadMapAt(p, j+1)
				if err != nil {
					return nil, j, err
				}

				if len(pkg) > 0 {
					obj, err := i.BuildObject(pkg, m)
					if err != nil {
						ans = append(ans, nil)
					}
					ans = append(ans, obj)
				} else {
					ans = append(ans, m)
				}
			}
		case 'V':
			if p[j+1] == 't' {
				j++
				var arr []interface{}

				arr, j, err = i.ReadArrayAt(p, j+1)
				if err != nil {
					return nil, j, err
				}

				ans = append(ans, arr)
			}
		}
	}

	return ans, j, nil
}

// ReadBytesAt Read bytes from given bytes and begin index.
//
// After 'B' chart, find b16 b8 <bytes-value>
func (i *DeserializerV1) ReadBytesAt(p []byte, begin int) ([]byte, int, error) {
	var ans []byte

	var idx = begin

	// find bytes length
	// b16 b8 <bytes-value>
	l := int(p[idx])<<8 + int(p[idx+1])
	idx++

	ans = make([]byte, l)

	for k := 0; k < l; k++ {
		idx++
		ans[k] = p[idx]
	}

	return ans, idx, nil
}

// ReadStringAt Read string from given bytes and begin index.
//
// After 'S' chart, find b16 b8 <string-value>
func (i *DeserializerV1) ReadStringAt(p []byte, begin int) (string, int, error) {
	var ans []byte

	var idx = begin

	// find string length
	// b16 b8 <string-value>
	l := int(p[idx])<<8 + int(p[idx+1])
	idx++

	ans = make([]byte, l)

	for k := 0; k < l; k++ {
		idx++
		ans[k] = p[idx]
	}

	return string(ans), idx, nil
}

// ReadInt32At Read string from given bytes and begin index.
//
// After 'I' chart, find b32 b24 b16 b8
func (i *DeserializerV1) ReadInt32At(p []byte, begin int) (int32, int, error) {
	var ans int32
	var idx = begin

	// b32 b24 b16 b8
	ans = int32(p[idx])<<24 + int32(p[idx+1])<<16 + int32(p[idx+2])<<8 + int32(p[idx+3])

	return ans, idx + 3, nil
}

// ReadInt64At Read string from given bytes and begin index.
//
// After 'L' chart, find b64 b56 b48 b40 b32 b24 b16 b8
func (i *DeserializerV1) ReadInt64At(p []byte, begin int) (int64, int, error) {
	var ans int64
	var idx = begin

	// b64 b56 b48 b40 b32 b24 b16 b8
	ans = int64(p[idx])<<56 + int64(p[idx+1])<<48 + int64(p[idx+2])<<40 + int64(p[idx+3])<<32 + int64(p[idx+4])<<24 + int64(p[idx+5])<<16 + int64(p[idx+6])<<8 + int64(p[idx+7])

	return ans, idx + 7, nil
}

// ReadFloat64At Read string from given bytes and begin index.
//
// After 'D' chart, find b64 b56 b48 b40 b32 b24 b16 b8
func (i *DeserializerV1) ReadFloat64At(p []byte, begin int) (float64, int, error) {
	var ans uint64
	var idx = begin

	// b64 b56 b48 b40 b32 b24 b16 b8
	ans = uint64(p[idx])<<56 + uint64(p[idx+1])<<48 + uint64(p[idx+2])<<40 + uint64(p[idx+3])<<32 + uint64(p[idx+4])<<24 + uint64(p[idx+5])<<16 + uint64(p[idx+6])<<8 + uint64(p[idx+7])

	return math.Float64frombits(ans), idx + 7, nil
}

// ReadMapAt Read string from given bytes and begin index.
//
// After 'Mt' chart, find <key> + <value> ...
func (i *DeserializerV1) ReadMapAt(p []byte, begin int) (map[interface{}]interface{}, int, error) {
	var ans = make(map[interface{}]interface{})
	var params []interface{}
	var idx = begin
	var err error

	params, idx, err = i.ReadAt(p, idx)
	if err != nil {
		return nil, idx, err
	}

	// put params to map
	for j := 0; j < len(params); j++ {
		ans[params[j]] = params[j+1]
		j++
	}

	return ans, idx, nil
}

// ReadArrayAt Read string from given bytes and begin index.
//
// After 'Vt <type> <size>' chart, find <value> + <value> ...
func (i *DeserializerV1) ReadArrayAt(p []byte, begin int) ([]interface{}, int, error) {
	// find array type
	var arrlen int32
	var arr, args []interface{}
	var idx = begin
	var err error

	_, idx, err = i.ReadStringAt(p, idx)
	if err != nil {
		return nil, idx, err
	}

	// find array length
	if p[idx+1] != 'l' {
		return nil, idx, fmt.Errorf("parse unknowen array length")
	}
	idx++

	arrlen, idx, err = i.ReadInt32At(p, idx+1)
	if err != nil {
		return nil, idx, err
	}

	arr = make([]interface{}, int(arrlen))
	// parse array val
	args, idx, err = i.ReadAt(p, idx+1)
	if err != nil {
		return nil, idx, err
	}

	if len(args) != len(arr) {
		return nil, idx, fmt.Errorf("array length and arguments length not equal %d != %d", len(arr), len(args))
	}

	for k := 0; k < len(arr); k++ {
		arr[k] = args[k]
	}

	return arr, idx, nil
}

// BuildObject build a struct or ptr(of struct) from given package and data
func (i *DeserializerV1) BuildObject(pkg string, data map[interface{}]interface{}) (interface{}, error) {
	var tt, t reflect.Type
	var exist bool

	// find pkg from type map
	if tt, exist = i.typeMap[pkg]; !exist {
		return nil, fmt.Errorf("can not find package %s in type map", pkg)
	}

	// check
	if tt.Kind() == reflect.Ptr {
		t = tt.Elem()
	} else {
		t = tt
	}

	// create new object (ptr)
	tv := reflect.New(t)

	// write value to fields
	for j, l := 0, t.NumField(); j < l; j++ {
		if val, exist := data[t.Field(j).Tag.Get(tagName)]; exist {
			if val == nil {
				continue
			}
			tv.Elem().Field(j).Set(reflect.ValueOf(val))
		}
	}

	// return ptr if mapping type is ptr
	if tt.Kind() == reflect.Ptr {
		return tv.Interface(), nil
	}

	// return struct if mapping type is struct
	return tv.Elem().Interface(), nil
}

// Read parse input (io.Reader) to return value
func (i *DeserializerV1) Read() ([]interface{}, error) {
	if i.r == nil {
		return nil, fmt.Errorf("io.Reader is not set yet")
	}

	p, err := ioutil.ReadAll(i.r)
	if err != nil {
		return nil, err
	}

	ans, _, err := i.ReadAt(p, 0)
	return ans, err
}

// Reset the io.Reader
func (i *DeserializerV1) Reset(r io.Reader) {
	i.r = r
}

// SetTypeMap set type map
func (i *DeserializerV1) SetTypeMap(typeMap map[string]reflect.Type) {
	i.typeMap = typeMap
}

// NewDeserializerV1 create DeserializerV1
func NewDeserializerV1() *DeserializerV1 {
	return &DeserializerV1{
		version: 1,
		typeMap: make(map[string]reflect.Type),
	}
}
