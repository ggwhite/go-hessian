package hessian

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"reflect"
)

const tagName = "hessian"

// SerializerV1 output stream for hessian 1.0 requests
type SerializerV1 struct {
	version int
	buf     *bytes.Buffer
	typeMap map[string]reflect.Type
}

// Writes a string value to the stream using UTF-8 encoding.
// b16 b8 string-value
func (o *SerializerV1) printString(s string) error {
	l := len(s)
	if err := o.buf.WriteByte(byte(l >> 8)); err != nil {
		return err
	}
	if err := o.buf.WriteByte(byte(l)); err != nil {
		return err
	}
	if _, err := o.buf.WriteString(s); err != nil {
		return err
	}
	return nil
}

// Writes an integer value to the stream.  The integer will be written with the following syntax:
// b32 b24 b16 b8
func (o *SerializerV1) printInt32(i int32) error {
	if err := o.buf.WriteByte(byte(i >> 24)); err != nil {
		return err
	}
	if err := o.buf.WriteByte(byte(i >> 16)); err != nil {
		return err
	}
	if err := o.buf.WriteByte(byte(i >> 8)); err != nil {
		return err
	}
	if err := o.buf.WriteByte(byte(i)); err != nil {
		return err
	}
	return nil
}

// Writes an long value to the stream.  The long will be written with the following syntax:
// b64 b56 b48 b40 b32 b24 b16 b8
func (o *SerializerV1) printInt64(i int64) error {
	if err := o.buf.WriteByte(byte(i >> 56)); err != nil {
		return err
	}
	if err := o.buf.WriteByte(byte(i >> 48)); err != nil {
		return err
	}
	if err := o.buf.WriteByte(byte(i >> 40)); err != nil {
		return err
	}
	if err := o.buf.WriteByte(byte(i >> 32)); err != nil {
		return err
	}
	if err := o.buf.WriteByte(byte(i >> 24)); err != nil {
		return err
	}
	if err := o.buf.WriteByte(byte(i >> 16)); err != nil {
		return err
	}
	if err := o.buf.WriteByte(byte(i >> 8)); err != nil {
		return err
	}
	if err := o.buf.WriteByte(byte(i)); err != nil {
		return err
	}
	return nil
}

// Call Writes a complete method call.
func (o *SerializerV1) Call(m string, args ...interface{}) error {
	if err := o.StartCall(); err != nil {
		return err
	}
	if err := o.WriteMethod(m); err != nil {
		return err
	}
	for _, arg := range args {
		if err := o.WriteObject(arg); err != nil {
			return err
		}
	}
	if err := o.CompleteCall(); err != nil {
		return err
	}
	return nil
}

// StartCall Starts the method call.
func (o *SerializerV1) StartCall() error {
	if err := o.buf.WriteByte('c'); err != nil {
		return err
	}
	if err := o.buf.WriteByte(byte(o.version)); err != nil {
		return err
	}
	if err := o.buf.WriteByte(byte(0)); err != nil {
		return err
	}
	return nil
}

// CompleteCall Completes
func (o *SerializerV1) CompleteCall() error {
	if err := o.buf.WriteByte('z'); err != nil {
		return err
	}
	return nil
}

// WriteMethod Writes the method tag.
//
// m b16 b8 method-name
func (o *SerializerV1) WriteMethod(m string) error {
	if err := o.buf.WriteByte('m'); err != nil {
		return err
	}
	o.printString(m)
	return nil
}

// WriteObject Writes any object to the output stream.
func (o *SerializerV1) WriteObject(arg interface{}) error {
	t := reflect.TypeOf(arg)
	if arg == nil {
		if err := o.WriteNull(); err != nil {
			return err
		}
		return nil
	}
	switch t.Kind() {
	default:
	case reflect.String:
		if err := o.WriteString(arg.(string)); err != nil {
			return err
		}
	case reflect.Bool:
		if err := o.WriteBool(arg.(bool)); err != nil {
			return err
		}
	case reflect.Int:
		if err := o.WriteInt(int32(arg.(int))); err != nil {
			return err
		}
	case reflect.Int8:
		if err := o.WriteInt(int32(arg.(int8))); err != nil {
			return err
		}
	case reflect.Int16:
		if err := o.WriteInt(int32(arg.(int16))); err != nil {
			return err
		}
	case reflect.Int32:
		if err := o.WriteInt(int32(arg.(int32))); err != nil {
			return err
		}
	case reflect.Int64:
		if err := o.WriteLong(arg.(int64)); err != nil {
			return err
		}
	case reflect.Float32:
		if err := o.WriteDouble(float64(arg.(float32))); err != nil {
			return err
		}
	case reflect.Float64:
		if err := o.WriteDouble(arg.(float64)); err != nil {
			return err
		}
	case reflect.Map:
		if err := o.WriteMap(arg); err != nil {
			return err
		}
	case reflect.Struct:
		if err := o.WriteStruct(arg); err != nil {
			return err
		}
	case reflect.Ptr:
		if err := o.WritePtr(arg); err != nil {
			return err
		}
	case reflect.Array:
		if err := o.WriteArray(arg); err != nil {
			return err
		}
	case reflect.Slice:
		if err := o.WriteArray(arg); err != nil {
			return err
		}
	}

	return nil
}

// WriteNull Writes a null value to the stream.
func (o *SerializerV1) WriteNull() error {
	if err := o.buf.WriteByte('N'); err != nil {
		return err
	}
	return nil
}

// WriteBytes Writes a bytes value to the stream
func (o *SerializerV1) WriteBytes(b []byte) error {
	if err := o.buf.WriteByte('B'); err != nil {
		return err
	}
	o.printString(string(b))
	return nil
}

// WriteString Writes a string value to the stream using UTF-8 encoding.
//
// The string will be written with the following syntax:
//
// S b16 b8 string-value
func (o *SerializerV1) WriteString(s string) error {
	if err := o.buf.WriteByte('S'); err != nil {
		return err
	}
	o.printString(s)
	return nil
}

// WriteBool Writes a boolean value to the stream.  The boolean will be written with the following syntax:
//
// T or F
func (o *SerializerV1) WriteBool(b bool) error {
	if b {
		if err := o.buf.WriteByte('T'); err != nil {
			return err
		}
		return nil
	}

	if err := o.buf.WriteByte('F'); err != nil {
		return err
	}
	return nil
}

// WriteInt Writes an integer value to the stream.  The integer will be written with the following syntax:
//
// I b32 b24 b16 b8
func (o *SerializerV1) WriteInt(i int32) error {
	if err := o.buf.WriteByte('I'); err != nil {
		return err
	}
	o.printInt32(i)
	return nil
}

// WriteLong Writes an long value to the stream.  The long will be written with the following syntax:
//
// L b64 b56 b48 b40 b32 b24 b16 b8
func (o *SerializerV1) WriteLong(i int64) error {
	if err := o.buf.WriteByte('L'); err != nil {
		return err
	}
	o.printInt64(i)
	return nil
}

// WriteDouble Writes an double value to the stream.  The double will be written with the following syntax:
//
// D b64 b56 b48 b40 b32 b24 b16 b8
func (o *SerializerV1) WriteDouble(i float64) error {
	n := math.Float64bits(i)

	if err := o.buf.WriteByte('D'); err != nil {
		return err
	}
	if err := o.buf.WriteByte(byte(n >> 56)); err != nil {
		return err
	}
	if err := o.buf.WriteByte(byte(n >> 48)); err != nil {
		return err
	}
	if err := o.buf.WriteByte(byte(n >> 40)); err != nil {
		return err
	}
	if err := o.buf.WriteByte(byte(n >> 32)); err != nil {
		return err
	}
	if err := o.buf.WriteByte(byte(n >> 24)); err != nil {
		return err
	}
	if err := o.buf.WriteByte(byte(n >> 16)); err != nil {
		return err
	}
	if err := o.buf.WriteByte(byte(n >> 8)); err != nil {
		return err
	}
	if err := o.buf.WriteByte(byte(n)); err != nil {
		return err
	}
	return nil
}

// WriteMap Write an map value to the stream. The map will be written with the following syntax:
//
// Mt b16 b8 (<key> <value>)z
func (o *SerializerV1) WriteMap(m interface{}) error {
	t := reflect.TypeOf(m)
	if t.Kind() != reflect.Map {
		return fmt.Errorf("WriteMap input is not a map")
	}
	if err := o.buf.WriteByte('M'); err != nil {
		return err
	}
	if err := o.buf.WriteByte('t'); err != nil {
		return err
	}
	if err := o.buf.WriteByte(byte(0)); err != nil {
		return err
	}
	if err := o.buf.WriteByte(byte(0)); err != nil {
		return err
	}

	v := reflect.ValueOf(m)
	for _, key := range v.MapKeys() {
		if err := o.WriteObject(key.Interface()); err != nil {
			return err
		}
		if err := o.WriteObject(v.MapIndex(key).Interface()); err != nil {
			return err
		}
	}

	if err := o.buf.WriteByte('z'); err != nil {
		return err
	}

	return nil
}

// WriteArray Write an map value to the stream. The map will be written with the following syntax:
//
// Vt b16 b8 <array-type> l b32 b24 b16 b8 <object> ... ... z
func (o *SerializerV1) WriteArray(arr interface{}) error {
	// check input
	t := reflect.TypeOf(arr)
	if t.Kind() != reflect.Array && t.Kind() != reflect.Slice {
		return fmt.Errorf("WriteArray input is not a array or slice")
	}
	if t.Elem().Kind() == reflect.Uint8 {
		return o.WriteBytes(arr.([]byte))
	}

	// write begin
	if err := o.buf.WriteByte('V'); err != nil {
		return err
	}
	if err := o.buf.WriteByte('t'); err != nil {
		return err
	}

	v := reflect.ValueOf(arr)

	var arrtype string
	// define array type
	switch t.Elem().Kind() {
	default:
		arrtype = "[object"
	case reflect.Interface:
		arrtype = "[object"
	case reflect.String:
		arrtype = "[string"
	case reflect.Int:
		arrtype = "[int"
	case reflect.Int8:
		arrtype = "[int"
	case reflect.Int16:
		arrtype = "[int"
	case reflect.Int32:
		arrtype = "[int"
	case reflect.Int64:
		arrtype = "[long"
	case reflect.Float32:
		arrtype = "[double"
	case reflect.Float64:
		arrtype = "[double"
	}
	o.printString(arrtype)

	if err := o.buf.WriteByte('l'); err != nil {
		return err
	}
	o.printInt32(int32(v.Len()))

	for i := 0; i < v.Len(); i++ {
		if err := o.WriteObject(v.Index(i).Interface()); err != nil {
			return err
		}
	}

	if err := o.buf.WriteByte('z'); err != nil {
		return err
	}

	return nil
}

// WriteStruct Writes an object value to the stream.
func (o *SerializerV1) WriteStruct(s interface{}) error {
	t := reflect.TypeOf(s)
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("WriteStruct input is not a struct")
	}
	if err := o.buf.WriteByte('M'); err != nil {
		return err
	}
	if err := o.buf.WriteByte('t'); err != nil {
		return err
	}

	v := reflect.ValueOf(s)

	var pkg string
	m := make(map[string]interface{})
	for i, l := 0, v.NumField(); i < l; i++ {
		if t.Field(i).Type == reflect.TypeOf(Package("")) {
			pkg = t.Field(i).Tag.Get(tagName)
			continue
		}
		if t.Field(i).Type.Kind() == reflect.Ptr && v.Field(i).IsNil() {
			m[t.Field(i).Tag.Get(tagName)] = nil
			continue
		}
		m[t.Field(i).Tag.Get(tagName)] = v.Field(i).Interface()
	}

	o.printString(pkg)

	for k, v := range m {
		if err := o.WriteObject(k); err != nil {
			return err
		}
		if err := o.WriteObject(v); err != nil {
			return err
		}
	}

	if err := o.buf.WriteByte('z'); err != nil {
		return err
	}
	return nil
}

// WritePtr Writes an object value to the stream.
func (o *SerializerV1) WritePtr(p interface{}) error {
	t := reflect.TypeOf(p)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("WritePtr input is not a pointer")
	}
	if err := o.buf.WriteByte('M'); err != nil {
		return err
	}
	if err := o.buf.WriteByte('t'); err != nil {
		return err
	}

	v := reflect.ValueOf(p).Elem()
	t = v.Type()

	var pkg string
	m := make(map[string]interface{})
	for i, l := 0, v.NumField(); i < l; i++ {
		if t.Field(i).Type == reflect.TypeOf(Package("")) {
			pkg = t.Field(i).Tag.Get(tagName)
			continue
		}
		if t.Field(i).Type.Kind() == reflect.Ptr && v.Field(i).IsNil() {
			m[t.Field(i).Tag.Get(tagName)] = nil
			continue
		}
		m[t.Field(i).Tag.Get(tagName)] = v.Field(i).Interface()
	}

	o.printString(pkg)

	for k, v := range m {
		if err := o.WriteObject(k); err != nil {
			return err
		}
		if err := o.WriteObject(v); err != nil {
			return err
		}
	}

	if err := o.buf.WriteByte('z'); err != nil {
		return err
	}
	return nil
}

// Reader get reader
func (o *SerializerV1) Reader() io.Reader {
	return o.buf
}

// Writer get writer
func (o *SerializerV1) Writer() io.Writer {
	return o.buf
}

// Flush clean writer
func (o *SerializerV1) Flush() {
	o.buf.Reset()
}

// SetTypeMap set type map
func (o *SerializerV1) SetTypeMap(typeMap map[string]reflect.Type) {
	o.typeMap = typeMap
}

// NewSerializerV1 create SerializerV1
func NewSerializerV1() *SerializerV1 {
	return &SerializerV1{
		buf:     bytes.NewBuffer(nil),
		version: 1,
		typeMap: make(map[string]reflect.Type),
	}
}
