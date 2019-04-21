package hessian

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"math"
	"reflect"
	"testing"
)

func TestSerializerV1_printString(t *testing.T) {
	type fields struct {
		version int
		buf     *bytes.Buffer
		typeMap map[string]reflect.Type
	}
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		wantBuf []byte
	}{
		{
			name: "Hello",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				s: "Hello",
			},
			wantErr: false,
			wantBuf: []byte{0, 5, 'H', 'e', 'l', 'l', 'o'},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &SerializerV1{
				version: tt.fields.version,
				buf:     tt.fields.buf,
				typeMap: tt.fields.typeMap,
			}
			if err := o.printString(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("SerializerV1.printString() error = %v, wantErr %v", err, tt.wantErr)
			}
			b, _ := ioutil.ReadAll(tt.fields.buf)
			if !reflect.DeepEqual(tt.wantBuf, b) {
				t.Errorf("SerializerV1.printString() content is %v, wantBuf %v", b, tt.wantBuf)
			}
		})
	}
}

func TestSerializerV1_printInt32(t *testing.T) {
	type fields struct {
		version int
		buf     *bytes.Buffer
		typeMap map[string]reflect.Type
	}
	type args struct {
		i int32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		wantBuf []byte
	}{
		{
			name: "Max of Int32",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				i: math.MaxInt32,
			},
			wantErr: false,
			wantBuf: []byte{0x7f, 0xff, 0xff, 0xff},
		},
		{
			name: "Min of Int32",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				i: math.MinInt32,
			},
			wantErr: false,
			wantBuf: []byte{0x80, 0x00, 0x00, 0x00},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &SerializerV1{
				version: tt.fields.version,
				buf:     tt.fields.buf,
				typeMap: tt.fields.typeMap,
			}
			if err := o.printInt32(tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("SerializerV1.printInt32() error = %v, wantErr %v", err, tt.wantErr)
			}
			b, _ := ioutil.ReadAll(tt.fields.buf)
			if !reflect.DeepEqual(tt.wantBuf, b) {
				t.Errorf("SerializerV1.printInt32() content is %v, wantBuf %v", b, tt.wantBuf)
			}
		})
	}
}

func TestSerializerV1_printInt64(t *testing.T) {
	type fields struct {
		version int
		buf     *bytes.Buffer
		typeMap map[string]reflect.Type
	}
	type args struct {
		i int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		wantBuf []byte
	}{
		{
			name: "Max of Int64",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				i: math.MaxInt64,
			},
			wantErr: false,
			wantBuf: []byte{0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		},
		{
			name: "Min of Int64",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				i: math.MinInt64,
			},
			wantErr: false,
			wantBuf: []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &SerializerV1{
				version: tt.fields.version,
				buf:     tt.fields.buf,
				typeMap: tt.fields.typeMap,
			}
			if err := o.printInt64(tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("SerializerV1.printInt64() error = %v, wantErr %v", err, tt.wantErr)
			}
			b, _ := ioutil.ReadAll(tt.fields.buf)
			if !reflect.DeepEqual(tt.wantBuf, b) {
				t.Errorf("SerializerV1.printInt64() content is %v, wantBuf %v", b, tt.wantBuf)
			}
		})
	}
}

func TestSerializerV1_Call(t *testing.T) {
	type fields struct {
		version int
		buf     *bytes.Buffer
		typeMap map[string]reflect.Type
	}
	type args struct {
		m    string
		args []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		wantBuf []byte
	}{
		{
			name: "call without arguments",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				m:    "str",
				args: []interface{}{},
			},
			wantErr: false,
			wantBuf: []byte{
				'c', 0x01, 0x00,
				'm', 0x00, 0x03, 's', 't', 'r',
				'z',
			},
		},
		{
			name: "call with arguments",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				m:    "strI",
				args: []interface{}{"hello"},
			},
			wantErr: false,
			wantBuf: []byte{
				'c', 0x01, 0x00,
				'm', 0x00, 0x04, 's', 't', 'r', 'I',
				'S', 0x00, 0x05, 'h', 'e', 'l', 'l', 'o',
				'z',
			},
		},
		{
			name: "call with arguments",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				m:    "intI",
				args: []interface{}{1},
			},
			wantErr: false,
			wantBuf: []byte{
				'c', 0x01, 0x00,
				'm', 0x00, 0x04, 'i', 'n', 't', 'I',
				'I', 0x00, 0x00, 0x00, 0x01,
				'z',
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &SerializerV1{
				version: tt.fields.version,
				buf:     tt.fields.buf,
				typeMap: tt.fields.typeMap,
			}
			if err := o.Call(tt.args.m, tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("SerializerV1.Call() error = %v, wantErr %v", err, tt.wantErr)
			}
			b, _ := ioutil.ReadAll(tt.fields.buf)
			if !reflect.DeepEqual(tt.wantBuf, b) {
				t.Errorf("SerializerV1.Call() content is %v, wantBuf %v", b, tt.wantBuf)
			}
		})
	}
}

func TestSerializerV1_StartCall(t *testing.T) {
	type fields struct {
		version int
		buf     *bytes.Buffer
		typeMap map[string]reflect.Type
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		wantBuf []byte
	}{
		{
			name: "StartCall",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			wantErr: false,
			wantBuf: []byte{'c', 0x01, 0x00},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &SerializerV1{
				version: tt.fields.version,
				buf:     tt.fields.buf,
				typeMap: tt.fields.typeMap,
			}
			if err := o.StartCall(); (err != nil) != tt.wantErr {
				t.Errorf("SerializerV1.StartCall() error = %v, wantErr %v", err, tt.wantErr)
			}
			b, _ := ioutil.ReadAll(tt.fields.buf)
			if !reflect.DeepEqual(tt.wantBuf, b) {
				t.Errorf("SerializerV1.StartCall() content is %v, wantBuf %v", b, tt.wantBuf)
			}
		})
	}
}

func TestSerializerV1_CompleteCall(t *testing.T) {
	type fields struct {
		version int
		buf     *bytes.Buffer
		typeMap map[string]reflect.Type
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		wantBuf []byte
	}{
		{
			name: "CompleteCall",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			wantErr: false,
			wantBuf: []byte{'z'},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &SerializerV1{
				version: tt.fields.version,
				buf:     tt.fields.buf,
				typeMap: tt.fields.typeMap,
			}
			if err := o.CompleteCall(); (err != nil) != tt.wantErr {
				t.Errorf("SerializerV1.CompleteCall() error = %v, wantErr %v", err, tt.wantErr)
			}
			b, _ := ioutil.ReadAll(tt.fields.buf)
			if !reflect.DeepEqual(tt.wantBuf, b) {
				t.Errorf("SerializerV1.CompleteCall() content is %v, wantBuf %v", b, tt.wantBuf)
			}
		})
	}
}

func TestSerializerV1_WriteMethod(t *testing.T) {
	type fields struct {
		version int
		buf     *bytes.Buffer
		typeMap map[string]reflect.Type
	}
	type args struct {
		m string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		wantBuf []byte
	}{
		{
			name: "method",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				m: "str",
			},
			wantErr: false,
			wantBuf: []byte{'m', 0x00, 0x03, 's', 't', 'r'},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &SerializerV1{
				version: tt.fields.version,
				buf:     tt.fields.buf,
				typeMap: tt.fields.typeMap,
			}
			if err := o.WriteMethod(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("SerializerV1.WriteMethod() error = %v, wantErr %v", err, tt.wantErr)
			}
			b, _ := ioutil.ReadAll(tt.fields.buf)
			if !reflect.DeepEqual(tt.wantBuf, b) {
				t.Errorf("SerializerV1.WriteMethod() content is %v, wantBuf %v", b, tt.wantBuf)
			}
		})
	}
}

func TestSerializerV1_WriteObject(t *testing.T) {
	type User struct {
		Package `hessian:"lab.ggw.shs.User"`
		Name    string `hessian:"name"`
	}
	type Account struct {
		Package  `hessian:"lab.ggw.shs.Account"`
		Password string
	}
	type fields struct {
		version int
		buf     *bytes.Buffer
		typeMap map[string]reflect.Type
	}
	type args struct {
		arg interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		wantBuf []byte
	}{
		{
			name: "nil",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				arg: nil,
			},
			wantErr: false,
			wantBuf: []byte{'N'},
		},
		{
			name: "bool true",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				arg: true,
			},
			wantErr: false,
			wantBuf: []byte{'T'},
		},
		{
			name: "bool false",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				arg: false,
			},
			wantErr: false,
			wantBuf: []byte{'F'},
		},
		{
			name: "int",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				arg: int(1),
			},
			wantErr: false,
			wantBuf: []byte{'I', 0x00, 0x00, 0x00, 0x01},
		},
		{
			name: "int8",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				arg: int8(1),
			},
			wantErr: false,
			wantBuf: []byte{'I', 0x00, 0x00, 0x00, 0x01},
		},
		{
			name: "int16",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				arg: int16(1),
			},
			wantErr: false,
			wantBuf: []byte{'I', 0x00, 0x00, 0x00, 0x01},
		},
		{
			name: "int32",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				arg: int32(1),
			},
			wantErr: false,
			wantBuf: []byte{'I', 0x00, 0x00, 0x00, 0x01},
		},
		{
			name: "int64",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				arg: int64(1),
			},
			wantErr: false,
			wantBuf: []byte{'L', 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		},
		{
			name: "float32",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				arg: float32(1),
			},
			wantErr: false,
			wantBuf: []byte{'D', 0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			name: "float64",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				arg: float64(1),
			},
			wantErr: false,
			wantBuf: []byte{'D', 0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			name: "map",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				arg: map[interface{}]interface{}{
					"KeyA": "ValueA",
				},
			},
			wantErr: false,
			wantBuf: []byte{
				'M', 't', 0x00, 0x00,
				'S', 0x00, 0x04, 'K', 'e', 'y', 'A',
				'S', 0x00, 0x06, 'V', 'a', 'l', 'u', 'e', 'A',
				'z',
			},
		},
		{
			name: "slice",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				arg: []int{1, 2, 3},
			},
			wantErr: false,
			wantBuf: []byte{
				'V', 't', 0x00, 0x04,
				'[', 'i', 'n', 't',
				'l', 0x00, 0x00, 0x00, 0x03,
				'I', 0x00, 0x00, 0x00, 0x01,
				'I', 0x00, 0x00, 0x00, 0x02,
				'I', 0x00, 0x00, 0x00, 0x03,
				'z',
			},
		},
		{
			name: "struct",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				arg: User{
					Name: "ggwhite",
				},
			},
			wantErr: false,
			wantBuf: []byte{
				'M', 't',
				0x00, 0x10, 'l', 'a', 'b', '.', 'g', 'g', 'w', '.', 's', 'h', 's', '.', 'U', 's', 'e', 'r',
				'S', 0x00, 0x04, 'n', 'a', 'm', 'e',
				'S', 0x00, 0x07, 'g', 'g', 'w', 'h', 'i', 't', 'e',
				'z',
			},
		},
		{
			name: "ptr",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				arg: &User{
					Name: "ggwhite",
				},
			},
			wantErr: false,
			wantBuf: []byte{
				'M', 't',
				0x00, 0x10, 'l', 'a', 'b', '.', 'g', 'g', 'w', '.', 's', 'h', 's', '.', 'U', 's', 'e', 'r',
				'S', 0x00, 0x04, 'n', 'a', 'm', 'e',
				'S', 0x00, 0x07, 'g', 'g', 'w', 'h', 'i', 't', 'e',
				'z',
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &SerializerV1{
				version: tt.fields.version,
				buf:     tt.fields.buf,
				typeMap: tt.fields.typeMap,
			}
			if err := o.WriteObject(tt.args.arg); (err != nil) != tt.wantErr {
				t.Errorf("SerializerV1.WriteObject() error = %v, wantErr %v", err, tt.wantErr)
			}
			b, _ := ioutil.ReadAll(tt.fields.buf)
			if !reflect.DeepEqual(tt.wantBuf, b) {
				t.Errorf("SerializerV1.WriteObject() content is %v (%s), wantBuf %v (%s)", b, string(b), tt.wantBuf, string(tt.wantBuf))
			}
		})
	}
}

func TestSerializerV1_WriteNull(t *testing.T) {
	type fields struct {
		version int
		buf     *bytes.Buffer
		typeMap map[string]reflect.Type
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		wantBuf []byte
	}{
		{
			name: "null",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			wantErr: false,
			wantBuf: []byte{'N'},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &SerializerV1{
				version: tt.fields.version,
				buf:     tt.fields.buf,
				typeMap: tt.fields.typeMap,
			}
			if err := o.WriteNull(); (err != nil) != tt.wantErr {
				t.Errorf("SerializerV1.WriteNull() error = %v, wantErr %v", err, tt.wantErr)
			}
			b, _ := ioutil.ReadAll(tt.fields.buf)
			if !reflect.DeepEqual(tt.wantBuf, b) {
				t.Errorf("SerializerV1.WriteNull() content is %v, wantBuf %v", b, tt.wantBuf)
			}
		})
	}
}

func TestSerializerV1_WriteBytes(t *testing.T) {
	type fields struct {
		version int
		buf     *bytes.Buffer
		typeMap map[string]reflect.Type
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		wantBuf []byte
	}{
		{
			name: "byte",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				b: []byte("qq"),
			},
			wantErr: false,
			wantBuf: []byte{'B', 0x00, 0x02, 'q', 'q'},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &SerializerV1{
				version: tt.fields.version,
				buf:     tt.fields.buf,
				typeMap: tt.fields.typeMap,
			}
			if err := o.WriteBytes(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("SerializerV1.WriteBytes() error = %v, wantErr %v", err, tt.wantErr)
			}
			b, _ := ioutil.ReadAll(tt.fields.buf)
			if !reflect.DeepEqual(tt.wantBuf, b) {
				t.Errorf("SerializerV1.WriteBytes() content is %v, wantBuf %v", b, tt.wantBuf)
			}
		})
	}
}

func TestSerializerV1_WriteString(t *testing.T) {
	type fields struct {
		version int
		buf     *bytes.Buffer
		typeMap map[string]reflect.Type
	}
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		wantBuf []byte
	}{
		{
			name: "string",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				s: "qq",
			},
			wantErr: false,
			wantBuf: []byte{'S', 0x00, 0x02, 'q', 'q'},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &SerializerV1{
				version: tt.fields.version,
				buf:     tt.fields.buf,
				typeMap: tt.fields.typeMap,
			}
			if err := o.WriteString(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("SerializerV1.WriteString() error = %v, wantErr %v", err, tt.wantErr)
			}
			b, _ := ioutil.ReadAll(tt.fields.buf)
			if !reflect.DeepEqual(tt.wantBuf, b) {
				t.Errorf("SerializerV1.WriteString() content is %v, wantBuf %v", b, tt.wantBuf)
			}
		})
	}
}

func TestSerializerV1_WriteBool(t *testing.T) {
	type fields struct {
		version int
		buf     *bytes.Buffer
		typeMap map[string]reflect.Type
	}
	type args struct {
		b bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		wantBuf []byte
	}{
		{
			name: "bool true",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				b: true,
			},
			wantErr: false,
			wantBuf: []byte{'T'},
		},
		{
			name: "bool false",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				b: false,
			},
			wantErr: false,
			wantBuf: []byte{'F'},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &SerializerV1{
				version: tt.fields.version,
				buf:     tt.fields.buf,
				typeMap: tt.fields.typeMap,
			}
			if err := o.WriteBool(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("SerializerV1.WriteBool() error = %v, wantErr %v", err, tt.wantErr)
			}
			b, _ := ioutil.ReadAll(tt.fields.buf)
			if !reflect.DeepEqual(tt.wantBuf, b) {
				t.Errorf("SerializerV1.WriteBool() content is %v, wantBuf %v", b, tt.wantBuf)
			}
		})
	}
}

func TestSerializerV1_WriteInt(t *testing.T) {
	type fields struct {
		version int
		buf     *bytes.Buffer
		typeMap map[string]reflect.Type
	}
	type args struct {
		i int32
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		wantBuf []byte
	}{
		{
			name: "max of int32",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				i: math.MaxInt32,
			},
			wantErr: false,
			wantBuf: []byte{'I', 0x7f, 0xff, 0xff, 0xff},
		},
		{
			name: "min of int32",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				i: math.MinInt32,
			},
			wantErr: false,
			wantBuf: []byte{'I', 0x80, 0x00, 0x00, 0x00},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &SerializerV1{
				version: tt.fields.version,
				buf:     tt.fields.buf,
				typeMap: tt.fields.typeMap,
			}
			if err := o.WriteInt(tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("SerializerV1.WriteInt() error = %v, wantErr %v", err, tt.wantErr)
			}
			b, _ := ioutil.ReadAll(tt.fields.buf)
			if !reflect.DeepEqual(tt.wantBuf, b) {
				t.Errorf("SerializerV1.WriteInt() content is %v, wantBuf %v", b, tt.wantBuf)
			}
		})
	}
}

func TestSerializerV1_WriteLong(t *testing.T) {
	type fields struct {
		version int
		buf     *bytes.Buffer
		typeMap map[string]reflect.Type
	}
	type args struct {
		i int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		wantBuf []byte
	}{
		{
			name: "max of int64",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				i: math.MaxInt64,
			},
			wantErr: false,
			wantBuf: []byte{'L', 0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		},
		{
			name: "min of int64",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				i: math.MinInt64,
			},
			wantErr: false,
			wantBuf: []byte{'L', 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &SerializerV1{
				version: tt.fields.version,
				buf:     tt.fields.buf,
				typeMap: tt.fields.typeMap,
			}
			if err := o.WriteLong(tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("SerializerV1.WriteLong() error = %v, wantErr %v", err, tt.wantErr)
			}
			b, _ := ioutil.ReadAll(tt.fields.buf)
			if !reflect.DeepEqual(tt.wantBuf, b) {
				t.Errorf("SerializerV1.WriteLong() content is %v, wantBuf %v", b, tt.wantBuf)
			}
		})
	}
}

func TestSerializerV1_WriteDouble(t *testing.T) {
	log.Println(math.MaxFloat64)
	type fields struct {
		version int
		buf     *bytes.Buffer
		typeMap map[string]reflect.Type
	}
	type args struct {
		i float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		wantBuf []byte
	}{
		{
			name: "max of float64",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				i: math.MaxFloat64,
			},
			wantErr: false,
			wantBuf: []byte{'D', 0x7f, 0xef, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		},
		{
			name: "min of float64",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				i: math.SmallestNonzeroFloat64,
			},
			wantErr: false,
			wantBuf: []byte{'D', 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &SerializerV1{
				version: tt.fields.version,
				buf:     tt.fields.buf,
				typeMap: tt.fields.typeMap,
			}
			if err := o.WriteDouble(tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("SerializerV1.WriteDouble() error = %v, wantErr %v", err, tt.wantErr)
			}
			b, _ := ioutil.ReadAll(tt.fields.buf)
			if !reflect.DeepEqual(tt.wantBuf, b) {
				t.Errorf("SerializerV1.WriteDouble() content is %v, wantBuf %v", b, tt.wantBuf)
			}
		})
	}
}

func TestSerializerV1_WriteMap(t *testing.T) {
	type fields struct {
		version int
		buf     *bytes.Buffer
		typeMap map[string]reflect.Type
	}
	type args struct {
		m interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		wantBuf []byte
	}{
		{
			name: "map",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				m: map[interface{}]interface{}{
					"KeyA": "ValueA",
				},
			},
			wantErr: false,
			// M t 0 0 S 0 4 K e y A S 0 6 V a l u e A I 0 0 0 2 S 0 6 V a l u e B S 0 4 V a l u e C z
			wantBuf: []byte{
				'M', 't', 0x00, 0x00,
				'S', 0x00, 0x04, 'K', 'e', 'y', 'A',
				'S', 0x00, 0x06, 'V', 'a', 'l', 'u', 'e', 'A',
				'z',
			},
		},
		{
			name: "map",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				m: map[interface{}]interface{}{
					2: "ValueB",
				},
			},
			wantErr: false,
			// M t 0 0 S 0 4 K e y A S 0 6 V a l u e A I 0 0 0 2 S 0 6 V a l u e B S 0 4 V a l u e C z
			wantBuf: []byte{
				'M', 't', 0x00, 0x00,
				'I', 0x00, 0x00, 0x00, 0x02,
				'S', 0x00, 0x06, 'V', 'a', 'l', 'u', 'e', 'B',
				'z',
			},
		},
		{
			name: "map",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				m: map[interface{}]interface{}{
					"KeyC": "ValueC",
				},
			},
			wantErr: false,
			// M t 0 0 S 0 4 K e y A S 0 6 V a l u e A I 0 0 0 2 S 0 6 V a l u e B S 0 4 V a l u e C z
			wantBuf: []byte{
				'M', 't', 0x00, 0x00,
				'S', 0x00, 0x04, 'K', 'e', 'y', 'C',
				'S', 0x00, 0x06, 'V', 'a', 'l', 'u', 'e', 'C',
				'z',
			},
		},
		{
			name: "not map type",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				m: "abc",
			},
			wantErr: true,
			wantBuf: []byte{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &SerializerV1{
				version: tt.fields.version,
				buf:     tt.fields.buf,
				typeMap: tt.fields.typeMap,
			}
			if err := o.WriteMap(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("SerializerV1.WriteMap() error = %v, wantErr %v", err, tt.wantErr)
			}
			b, _ := ioutil.ReadAll(tt.fields.buf)
			if !reflect.DeepEqual(tt.wantBuf, b) {
				t.Errorf("SerializerV1.WriteMap() content is %v (%s), wantBuf %v (%s)", b, string(b), tt.wantBuf, string(tt.wantBuf))
			}
		})
	}
}

func TestSerializerV1_WriteArray(t *testing.T) {
	type fields struct {
		version int
		buf     *bytes.Buffer
		typeMap map[string]reflect.Type
	}
	type args struct {
		arr interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		wantBuf []byte
	}{
		{
			name: "interface slice",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				arr: []interface{}{
					"ValueA", int32(2), "ValueB", int64(4),
				},
			},
			wantErr: false,
			// V t 0 7 [object l 0 0 0 4 S 0 6 ValueA I 0 0 0 2 S 0 6 ValueB L 0 0 0 0 0 0 0 4 z
			wantBuf: []byte{
				'V', 't', 0x00, 0x07,
				'[', 'o', 'b', 'j', 'e', 'c', 't',
				'l', 0x00, 0x00, 0x00, 0x04,
				'S', 0x00, 0x06, 'V', 'a', 'l', 'u', 'e', 'A',
				'I', 0x00, 0x00, 0x00, 0x02,
				'S', 0x00, 0x06, 'V', 'a', 'l', 'u', 'e', 'B',
				'L', 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04,
				'z',
			},
		},
		{
			name: "string slice",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				arr: []string{
					"ValueA", "ValueB", "ValueC",
				},
			},
			wantErr: false,
			// V t 0 7 [string l 0 0 0 3 S 0 6 ValueA S 0 6 ValueB S 0 6 ValueC z
			wantBuf: []byte{
				'V', 't', 0x00, 0x07,
				'[', 's', 't', 'r', 'i', 'n', 'g',
				'l', 0x00, 0x00, 0x00, 0x03,
				'S', 0x00, 0x06, 'V', 'a', 'l', 'u', 'e', 'A',
				'S', 0x00, 0x06, 'V', 'a', 'l', 'u', 'e', 'B',
				'S', 0x00, 0x06, 'V', 'a', 'l', 'u', 'e', 'C',
				'z',
			},
		},
		{
			name: "int slice",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				arr: []int{1, 2, 3},
			},
			wantErr: false,
			// V t 0 4 [int l 0 0 0 3 I 0 0 0 1 I 0 0 0 2 I 0 0 0 3 z
			wantBuf: []byte{
				'V', 't', 0x00, 0x04,
				'[', 'i', 'n', 't',
				'l', 0x00, 0x00, 0x00, 0x03,
				'I', 0x00, 0x00, 0x00, 0x01,
				'I', 0x00, 0x00, 0x00, 0x02,
				'I', 0x00, 0x00, 0x00, 0x03,
				'z',
			},
		},
		{
			name: "int8 slice",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				arr: []int8{int8(1), int8(2), int8(3)},
			},
			wantErr: false,
			// V t 0 4 [int l 0 0 0 3 I 0 0 0 1 I 0 0 0 2 I 0 0 0 3 z
			wantBuf: []byte{
				'V', 't', 0x00, 0x04,
				'[', 'i', 'n', 't',
				'l', 0x00, 0x00, 0x00, 0x03,
				'I', 0x00, 0x00, 0x00, 0x01,
				'I', 0x00, 0x00, 0x00, 0x02,
				'I', 0x00, 0x00, 0x00, 0x03,
				'z',
			},
		},
		{
			name: "int16 slice",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				arr: []int16{int16(1), int16(2), int16(3)},
			},
			wantErr: false,
			// V t 0 4 [int l 0 0 0 3 I 0 0 0 1 I 0 0 0 2 I 0 0 0 3 z
			wantBuf: []byte{
				'V', 't', 0x00, 0x04,
				'[', 'i', 'n', 't',
				'l', 0x00, 0x00, 0x00, 0x03,
				'I', 0x00, 0x00, 0x00, 0x01,
				'I', 0x00, 0x00, 0x00, 0x02,
				'I', 0x00, 0x00, 0x00, 0x03,
				'z',
			},
		},
		{
			name: "int32 slice",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				arr: []int32{int32(1), int32(2), int32(3)},
			},
			wantErr: false,
			// V t 0 4 [int l 0 0 0 3 I 0 0 0 1 I 0 0 0 2 I 0 0 0 3 z
			wantBuf: []byte{
				'V', 't', 0x00, 0x04,
				'[', 'i', 'n', 't',
				'l', 0x00, 0x00, 0x00, 0x03,
				'I', 0x00, 0x00, 0x00, 0x01,
				'I', 0x00, 0x00, 0x00, 0x02,
				'I', 0x00, 0x00, 0x00, 0x03,
				'z',
			},
		},
		{
			name: "int64 slice",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				arr: []int64{int64(1), int64(2), int64(3)},
			},
			wantErr: false,
			// V t 0 5 [long l 0 0 0 3 L 0 0 0 0 0 0 0 1 L 0 0 0 0 0 0 0 2 L 0 0 0 0 0 0 0 3 z
			wantBuf: []byte{
				'V', 't', 0x00, 0x05,
				'[', 'l', 'o', 'n', 'g',
				'l', 0x00, 0x00, 0x00, 0x03,
				'L', 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
				'L', 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
				'L', 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03,
				'z',
			},
		},
		{
			name: "float32 slice",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				arr: []float32{float32(1), float32(2), float32(3)},
			},
			wantErr: false,
			// V t 0 7 [double l 0 0 0 3 D 0 0 0 0 0 0 0 1 D 0 0 0 0 0 0 0 2 D 0 0 0 0 0 0 0 3 z
			wantBuf: []byte{
				'V', 't', 0x00, 0x07,
				'[', 'd', 'o', 'u', 'b', 'l', 'e',
				'l', 0x00, 0x00, 0x00, 0x03,
				'D', 0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				'D', 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				'D', 0x40, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				'z',
			},
		},
		{
			name: "float64 slice",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				arr: []float64{float64(1), float64(2), float64(3)},
			},
			wantErr: false,
			// V t 0 7 [double l 0 0 0 3 D 0 0 0 0 0 0 0 1 D 0 0 0 0 0 0 0 2 D 0 0 0 0 0 0 0 3 z
			wantBuf: []byte{
				'V', 't', 0x00, 0x07,
				'[', 'd', 'o', 'u', 'b', 'l', 'e',
				'l', 0x00, 0x00, 0x00, 0x03,
				'D', 0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				'D', 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				'D', 0x40, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				'z',
			},
		},
		{
			name: "not array type",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				arr: "abc",
			},
			wantErr: true,
			wantBuf: []byte{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &SerializerV1{
				version: tt.fields.version,
				buf:     tt.fields.buf,
				typeMap: tt.fields.typeMap,
			}
			if err := o.WriteArray(tt.args.arr); (err != nil) != tt.wantErr {
				t.Errorf("SerializerV1.WriteArray() error = %v, wantErr %v", err, tt.wantErr)
			}
			b, _ := ioutil.ReadAll(tt.fields.buf)
			if !reflect.DeepEqual(tt.wantBuf, b) {
				t.Errorf("SerializerV1.WriteArray() content is %v, wantBuf %v", b, tt.wantBuf)
			}
		})
	}
}

func TestSerializerV1_WriteStruct(t *testing.T) {
	type User struct {
		Package `hessian:"lab.ggw.shs.User"`
		Name    string `hessian:"name"`
	}
	type Account struct {
		Package  `hessian:"lab.ggw.shs.Account"`
		Password string
	}
	type fields struct {
		version int
		buf     *bytes.Buffer
		typeMap map[string]reflect.Type
	}
	type args struct {
		s interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		wantBuf []byte
	}{
		{
			name: "struct",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				s: User{
					Name: "ggwhite",
				},
			},
			wantErr: false,
			wantBuf: []byte{
				'M', 't',
				0x00, 0x10, 'l', 'a', 'b', '.', 'g', 'g', 'w', '.', 's', 'h', 's', '.', 'U', 's', 'e', 'r',
				'S', 0x00, 0x04, 'n', 'a', 'm', 'e',
				'S', 0x00, 0x07, 'g', 'g', 'w', 'h', 'i', 't', 'e',
				'z',
			},
		},
		{
			name: "struct without tag",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				s: Account{
					Password: "ggw.chang@gmail.com",
				},
			},
			wantErr: false,
			wantBuf: []byte{
				'M', 't',
				0x00, 0x13, 'l', 'a', 'b', '.', 'g', 'g', 'w', '.', 's', 'h', 's', '.', 'A', 'c', 'c', 'o', 'u', 'n', 't',
				'z',
			},
		},
		{
			name: "not struct type",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				s: "abc",
			},
			wantErr: true,
			wantBuf: []byte{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &SerializerV1{
				version: tt.fields.version,
				buf:     tt.fields.buf,
				typeMap: tt.fields.typeMap,
			}
			if err := o.WriteStruct(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("SerializerV1.WriteStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
			b, _ := ioutil.ReadAll(tt.fields.buf)
			if !reflect.DeepEqual(tt.wantBuf, b) {
				t.Errorf("SerializerV1.WriteStruct() content is %v (%s), wantBuf %v (%s)", b, string(b), tt.wantBuf, string(tt.wantBuf))
			}
		})
	}
}

func TestSerializerV1_WritePtr(t *testing.T) {
	type User struct {
		Package `hessian:"lab.ggw.shs.User"`
		Name    string `hessian:"name"`
	}
	type Account struct {
		Package  `hessian:"lab.ggw.shs.Account"`
		Password string
	}
	var str string
	str = "abc"
	type fields struct {
		version int
		buf     *bytes.Buffer
		typeMap map[string]reflect.Type
	}
	type args struct {
		p interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		wantBuf []byte
	}{
		{
			name: "ptr struct",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				p: &User{
					Name: "ggwhite",
				},
			},
			wantErr: false,
			wantBuf: []byte{
				'M', 't',
				0x00, 0x10, 'l', 'a', 'b', '.', 'g', 'g', 'w', '.', 's', 'h', 's', '.', 'U', 's', 'e', 'r',
				'S', 0x00, 0x04, 'n', 'a', 'm', 'e',
				'S', 0x00, 0x07, 'g', 'g', 'w', 'h', 'i', 't', 'e',
				'z',
			},
		},
		{
			name: "ptr struct without tag",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				p: &Account{
					Password: "ggw.chang@gmail.com",
				},
			},
			wantErr: false,
			wantBuf: []byte{
				'M', 't',
				0x00, 0x13, 'l', 'a', 'b', '.', 'g', 'g', 'w', '.', 's', 'h', 's', '.', 'A', 'c', 'c', 'o', 'u', 'n', 't',
				'z',
			},
		},
		{
			name: "ptr string",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				p: &str,
			},
			wantErr: true,
			wantBuf: []byte{},
		},
		{
			name: "not ptr type",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			args: args{
				p: "abc",
			},
			wantErr: true,
			wantBuf: []byte{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &SerializerV1{
				version: tt.fields.version,
				buf:     tt.fields.buf,
				typeMap: tt.fields.typeMap,
			}
			if err := o.WritePtr(tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("SerializerV1.WritePtr() error = %v, wantErr %v", err, tt.wantErr)
			}
			b, _ := ioutil.ReadAll(tt.fields.buf)
			if !reflect.DeepEqual(tt.wantBuf, b) {
				t.Errorf("SerializerV1.WritePtr() content is %v (%s), wantBuf %v (%s)", b, string(b), tt.wantBuf, string(tt.wantBuf))
			}
		})
	}
}

func TestSerializerV1_Reader(t *testing.T) {
	type fields struct {
		version int
		buf     *bytes.Buffer
		typeMap map[string]reflect.Type
	}
	tests := []struct {
		name   string
		fields fields
		want   io.Reader
	}{
		{
			name: "get reader",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			want: bytes.NewBuffer(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &SerializerV1{
				version: tt.fields.version,
				buf:     tt.fields.buf,
				typeMap: tt.fields.typeMap,
			}
			if got := o.Reader(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SerializerV1.Reader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializerV1_Writer(t *testing.T) {
	type fields struct {
		version int
		buf     *bytes.Buffer
		typeMap map[string]reflect.Type
	}
	tests := []struct {
		name   string
		fields fields
		want   io.Writer
	}{
		{
			name: "get writer",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
			want: bytes.NewBuffer(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &SerializerV1{
				version: tt.fields.version,
				buf:     tt.fields.buf,
				typeMap: tt.fields.typeMap,
			}
			if got := o.Writer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SerializerV1.Writer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSerializerV1_Flush(t *testing.T) {
	type fields struct {
		version int
		buf     *bytes.Buffer
		typeMap map[string]reflect.Type
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "flush",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &SerializerV1{
				version: tt.fields.version,
				buf:     tt.fields.buf,
				typeMap: tt.fields.typeMap,
			}
			o.Flush()
		})
	}
}

func TestSerializerV1_SetTypeMap(t *testing.T) {
	type fields struct {
		version int
		buf     *bytes.Buffer
		typeMap map[string]reflect.Type
	}
	type args struct {
		typeMap map[string]reflect.Type
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "set type map",
			fields: fields{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: make(map[string]reflect.Type),
			},
			args: args{
				typeMap: make(map[string]reflect.Type),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &SerializerV1{
				version: tt.fields.version,
				buf:     tt.fields.buf,
				typeMap: tt.fields.typeMap,
			}
			o.SetTypeMap(tt.args.typeMap)
		})
	}
}

func TestNewSerializerV1(t *testing.T) {
	tests := []struct {
		name string
		want *SerializerV1
	}{
		{
			name: "new one",
			want: &SerializerV1{
				version: 1,
				buf:     bytes.NewBuffer(nil),
				typeMap: make(map[string]reflect.Type),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSerializerV1(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSerializerV1() = %v, want %v", got, tt.want)
			}
		})
	}
}
