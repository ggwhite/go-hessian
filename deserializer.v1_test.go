package hessian

import (
	"bytes"
	"io"
	"math"
	"reflect"
	"testing"
)

func TestDeserializerV1_ReadAt(t *testing.T) {
	type User struct {
		Package `hessian:"package.User"`
		Name    string `hessian:"name"`
	}
	type fields struct {
		version int
		r       io.Reader
		typeMap map[string]reflect.Type
	}
	type args struct {
		p     []byte
		begin int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []interface{}
		want1   int
		wantErr bool
	}{
		{
			name: "bool & nil",
			fields: fields{
				version: 1,
				r:       nil,
				typeMap: nil,
			},
			args: args{
				p:     []byte{'r', 0x01, 0x00, 'T', 'F', 'N', 'z'},
				begin: 0,
			},
			want: []interface{}{
				true,
				false,
				nil,
			},
			want1:   6,
			wantErr: false,
		},
		{
			name: "bytes",
			fields: fields{
				version: 1,
				r:       nil,
				typeMap: nil,
			},
			args: args{
				p:     []byte{'r', 0x01, 0x00, 'B', 0x00, 0x02, 'a', 'b', 'z'},
				begin: 0,
			},
			want: []interface{}{
				[]byte{'a', 'b'},
			},
			want1:   8,
			wantErr: false,
		},
		{
			name: "string",
			fields: fields{
				version: 1,
				r:       nil,
				typeMap: nil,
			},
			args: args{
				p:     []byte{'r', 0x01, 0x00, 'S', 0x00, 0x02, 'a', 'b', 'z'},
				begin: 0,
			},
			want: []interface{}{
				"ab",
			},
			want1:   8,
			wantErr: false,
		},
		{
			name: "int",
			fields: fields{
				version: 1,
				r:       nil,
				typeMap: nil,
			},
			args: args{
				p:     []byte{'r', 0x01, 0x00, 'I', 0x00, 0x00, 0x00, 0x02, 'z'},
				begin: 0,
			},
			want:    []interface{}{int32(2)},
			want1:   8,
			wantErr: false,
		},
		{
			name: "long",
			fields: fields{
				version: 1,
				r:       nil,
				typeMap: nil,
			},
			args: args{
				p:     []byte{'r', 0x01, 0x00, 'L', 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 'z'},
				begin: 0,
			},
			want:    []interface{}{int64(2)},
			want1:   12,
			wantErr: false,
		},
		{
			name: "double",
			fields: fields{
				version: 1,
				r:       nil,
				typeMap: nil,
			},
			args: args{
				p:     []byte{'r', 0x01, 0x00, 'D', 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 'z'},
				begin: 0,
			},
			want:    []interface{}{math.SmallestNonzeroFloat64},
			want1:   12,
			wantErr: false,
		},
		{
			name: "map",
			fields: fields{
				version: 1,
				r:       nil,
				typeMap: nil,
			},
			args: args{
				p:     []byte{'r', 0x01, 0x00, 'M', 't', 0x00, 0x00, 'S', 0x00, 0x01, 'a', 'S', 0x00, 0x02, 'b', 'c', 'z', 'z'},
				begin: 0,
			},
			want: []interface{}{
				map[interface{}]interface{}{
					"a": "bc",
				},
			},
			want1:   17,
			wantErr: false,
		},
		{
			name: "struct",
			fields: fields{
				version: 1,
				r:       nil,
				typeMap: map[string]reflect.Type{
					"package.User": reflect.TypeOf(User{}),
				},
			},
			args: args{
				p: []byte{
					'r', 0x01, 0x00,
					'M', 't',
					0x00, 0x0c, 'p', 'a', 'c', 'k', 'a', 'g', 'e', '.', 'U', 's', 'e', 'r',
					'S', 0x00, 0x04, 'n', 'a', 'm', 'e',
					'S', 0x00, 0x03, 'w', 't', 'f',
					'z',
					'z',
				},
				begin: 0,
			},
			want: []interface{}{
				User{
					Name: "wtf",
				},
			},
			want1:   33,
			wantErr: false,
		},
		{
			name: "array",
			fields: fields{
				version: 1,
				r:       nil,
				typeMap: map[string]reflect.Type{
					"package.User": reflect.TypeOf(User{}),
				},
			},
			args: args{
				p: []byte{
					'r', 0x01, 0x00,
					'V', 't',
					0x00, 0x07, '[', 'o', 'b', 'j', 'e', 'c', 't',
					'l', 0x00, 0x00, 0x00, 0x02,
					'S', 0x00, 0x05, 'b', 'n', 'a', 'm', 'e',
					'B', 0x00, 0x07, 'g', 'g', 'w', 'h', 'i', 't', 'e',
					'z',
					'z',
				},
				begin: 0,
			},
			want: []interface{}{
				[]interface{}{
					"bname",
					[]byte("ggwhite"),
				},
			},
			want1:   38,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &DeserializerV1{
				version: tt.fields.version,
				r:       tt.fields.r,
				typeMap: tt.fields.typeMap,
			}
			got, got1, err := i.ReadAt(tt.args.p, tt.args.begin)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeserializerV1.ReadAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeserializerV1.ReadAt() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("DeserializerV1.ReadAt() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDeserializerV1_ReadBytesAt(t *testing.T) {
	type fields struct {
		version int
		r       io.Reader
		typeMap map[string]reflect.Type
	}
	type args struct {
		p     []byte
		begin int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		want1   int
		wantErr bool
	}{
		{
			name: "Read begin 0",
			fields: fields{
				version: 1,
				r:       nil,
				typeMap: nil,
			},
			args: args{
				p:     []byte{0x00, 0x02, 'a', 'b'},
				begin: 0,
			},
			want:    []byte{'a', 'b'},
			want1:   3,
			wantErr: false,
		},
		{
			name: "Read begin not 0",
			fields: fields{
				version: 1,
				r:       nil,
				typeMap: nil,
			},
			args: args{
				p:     []byte{'a', 'b', 'c', 0x00, 0x02, 'a', 'b'},
				begin: 3,
			},
			want:    []byte{'a', 'b'},
			want1:   6,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &DeserializerV1{
				version: tt.fields.version,
				r:       tt.fields.r,
				typeMap: tt.fields.typeMap,
			}
			got, got1, err := i.ReadBytesAt(tt.args.p, tt.args.begin)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeserializerV1.ReadBytesAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeserializerV1.ReadBytesAt() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("DeserializerV1.ReadBytesAt() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDeserializerV1_ReadStringAt(t *testing.T) {
	type fields struct {
		version int
		r       io.Reader
		typeMap map[string]reflect.Type
	}
	type args struct {
		p     []byte
		begin int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		want1   int
		wantErr bool
	}{
		{
			name: "Read begin 0",
			fields: fields{
				version: 1,
				r:       nil,
				typeMap: nil,
			},
			args: args{
				p:     []byte{0x00, 0x03, 'a', 'b', 'c'},
				begin: 0,
			},
			want:    "abc",
			want1:   4,
			wantErr: false,
		},
		{
			name: "Read begin not 0",
			fields: fields{
				version: 1,
				r:       nil,
				typeMap: nil,
			},
			args: args{
				p:     []byte{'a', 'b', 'c', 0x00, 0x05, 'a', 'b', 'c', 'd', 'e'},
				begin: 3,
			},
			want:    "abcde",
			want1:   9,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &DeserializerV1{
				version: tt.fields.version,
				r:       tt.fields.r,
				typeMap: tt.fields.typeMap,
			}
			got, got1, err := i.ReadStringAt(tt.args.p, tt.args.begin)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeserializerV1.ReadStringAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeserializerV1.ReadStringAt() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("DeserializerV1.ReadStringAt() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDeserializerV1_ReadInt32At(t *testing.T) {
	type fields struct {
		version int
		r       io.Reader
		typeMap map[string]reflect.Type
	}
	type args struct {
		p     []byte
		begin int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int32
		want1   int
		wantErr bool
	}{
		{
			name: "Read begin 0",
			fields: fields{
				version: 1,
				r:       nil,
				typeMap: nil,
			},
			args: args{
				p:     []byte{0x7f, 0xff, 0xff, 0xff},
				begin: 0,
			},
			want:    math.MaxInt32,
			want1:   3,
			wantErr: false,
		},
		{
			name: "Read begin not 0",
			fields: fields{
				version: 1,
				r:       nil,
				typeMap: nil,
			},
			args: args{
				p:     []byte{'b', 'c', 0x80, 0x00, 0x00, 0x00},
				begin: 2,
			},
			want:    math.MinInt32,
			want1:   5,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &DeserializerV1{
				version: tt.fields.version,
				r:       tt.fields.r,
				typeMap: tt.fields.typeMap,
			}
			got, got1, err := i.ReadInt32At(tt.args.p, tt.args.begin)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeserializerV1.ReadInt32At() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeserializerV1.ReadInt32At() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("DeserializerV1.ReadInt32At() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDeserializerV1_ReadInt64At(t *testing.T) {
	type fields struct {
		version int
		r       io.Reader
		typeMap map[string]reflect.Type
	}
	type args struct {
		p     []byte
		begin int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		want1   int
		wantErr bool
	}{
		{
			name: "Read begin 0",
			fields: fields{
				version: 1,
				r:       nil,
				typeMap: nil,
			},
			args: args{
				p:     []byte{0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
				begin: 0,
			},
			want:    math.MaxInt64,
			want1:   7,
			wantErr: false,
		},
		{
			name: "Read begin not 0",
			fields: fields{
				version: 1,
				r:       nil,
				typeMap: nil,
			},
			args: args{
				p:     []byte{'b', 'c', 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
				begin: 2,
			},
			want:    math.MinInt64,
			want1:   9,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &DeserializerV1{
				version: tt.fields.version,
				r:       tt.fields.r,
				typeMap: tt.fields.typeMap,
			}
			got, got1, err := i.ReadInt64At(tt.args.p, tt.args.begin)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeserializerV1.ReadInt64At() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeserializerV1.ReadInt64At() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("DeserializerV1.ReadInt64At() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDeserializerV1_ReadFloat64At(t *testing.T) {
	type fields struct {
		version int
		r       io.Reader
		typeMap map[string]reflect.Type
	}
	type args struct {
		p     []byte
		begin int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		want1   int
		wantErr bool
	}{
		{
			name: "Read begin 0",
			fields: fields{
				version: 1,
				r:       nil,
				typeMap: nil,
			},
			args: args{
				p:     []byte{0x7f, 0xef, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
				begin: 0,
			},
			want:    math.MaxFloat64,
			want1:   7,
			wantErr: false,
		},
		{
			name: "Read begin not 0",
			fields: fields{
				version: 1,
				r:       nil,
				typeMap: nil,
			},
			args: args{
				p:     []byte{'b', 'c', 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01},
				begin: 2,
			},
			want:    math.SmallestNonzeroFloat64,
			want1:   9,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &DeserializerV1{
				version: tt.fields.version,
				r:       tt.fields.r,
				typeMap: tt.fields.typeMap,
			}
			got, got1, err := i.ReadFloat64At(tt.args.p, tt.args.begin)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeserializerV1.ReadFloat64At() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DeserializerV1.ReadFloat64At() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("DeserializerV1.ReadFloat64At() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDeserializerV1_ReadMapAt(t *testing.T) {
	type fields struct {
		version int
		r       io.Reader
		typeMap map[string]reflect.Type
	}
	type args struct {
		p     []byte
		begin int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[interface{}]interface{}
		want1   int
		wantErr bool
	}{
		{
			name: "Read begin 0",
			fields: fields{
				version: 1,
				r:       nil,
				typeMap: nil,
			},
			args: args{
				p: []byte{
					'S', 0x00, 0x05, 'b', 'n', 'a', 'm', 'e',
					'B', 0x00, 0x07, 'g', 'g', 'w', 'h', 'i', 't', 'e',
				},
				begin: 0,
			},
			want: map[interface{}]interface{}{
				"bname": []byte("ggwhite"),
			},
			want1:   18,
			wantErr: false,
		},
		{
			name: "Read begin not 0",
			fields: fields{
				version: 1,
				r:       nil,
				typeMap: nil,
			},
			args: args{
				p: []byte{
					'a', 'b', 'c', 'd',
					'S', 0x00, 0x04, 'n', 'a', 'm', 'e',
					'S', 0x00, 0x07, 'g', 'g', 'w', 'h', 'i', 't', 'e',
				},
				begin: 4,
			},
			want: map[interface{}]interface{}{
				"name": "ggwhite",
			},
			want1:   21,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &DeserializerV1{
				version: tt.fields.version,
				r:       tt.fields.r,
				typeMap: tt.fields.typeMap,
			}
			got, got1, err := i.ReadMapAt(tt.args.p, tt.args.begin)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeserializerV1.ReadMapAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeserializerV1.ReadMapAt() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("DeserializerV1.ReadMapAt() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDeserializerV1_ReadArrayAt(t *testing.T) {
	type fields struct {
		version int
		r       io.Reader
		typeMap map[string]reflect.Type
	}
	type args struct {
		p     []byte
		begin int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []interface{}
		want1   int
		wantErr bool
	}{
		{
			name: "Read begin 0",
			fields: fields{
				version: 1,
				r:       nil,
				typeMap: nil,
			},
			args: args{
				p: []byte{
					0x00, 0x07, '[', 'o', 'b', 'j', 'e', 'c', 't',
					'l', 0x00, 0x00, 0x00, 0x02,
					'S', 0x00, 0x05, 'b', 'n', 'a', 'm', 'e',
					'B', 0x00, 0x07, 'g', 'g', 'w', 'h', 'i', 't', 'e',
				},
				begin: 0,
			},
			want: []interface{}{
				"bname",
				[]byte("ggwhite"),
			},
			want1:   32,
			wantErr: false,
		},
		{
			name: "Read begin not 0",
			fields: fields{
				version: 1,
				r:       nil,
				typeMap: nil,
			},
			args: args{
				p: []byte{
					'a', 'b', 'c', 'd',
					0x00, 0x07, '[', 'o', 'b', 'j', 'e', 'c', 't',
					'l', 0x00, 0x00, 0x00, 0x02,
					'S', 0x00, 0x05, 'b', 'n', 'a', 'm', 'e',
					'B', 0x00, 0x07, 'g', 'g', 'w', 'h', 'i', 't', 'e',
				},
				begin: 4,
			},
			want: []interface{}{
				"bname",
				[]byte("ggwhite"),
			},
			want1:   36,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &DeserializerV1{
				version: tt.fields.version,
				r:       tt.fields.r,
				typeMap: tt.fields.typeMap,
			}
			got, got1, err := i.ReadArrayAt(tt.args.p, tt.args.begin)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeserializerV1.ReadArrayAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeserializerV1.ReadArrayAt() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("DeserializerV1.ReadArrayAt() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDeserializerV1_BuildObject(t *testing.T) {
	type User struct {
		Package `hessian:"package.User"`
		Name    string `hessian:"name"`
	}
	type Account struct {
		Package `hessian:"package.Account"`
		Name    string `hessian:"name"`
	}
	type fields struct {
		version int
		r       io.Reader
		typeMap map[string]reflect.Type
	}
	type args struct {
		pkg  string
		data map[interface{}]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "struct",
			fields: fields{
				version: 1,
				r:       nil,
				typeMap: map[string]reflect.Type{
					"package.User":    reflect.TypeOf(User{}),
					"package.Account": reflect.TypeOf(Account{}),
				},
			},
			args: args{
				pkg: "package.User",
				data: map[interface{}]interface{}{
					"name": "ggwhite",
				},
			},
			want: User{
				Name: "ggwhite",
			},
			wantErr: false,
		},
		{
			name: "ptr",
			fields: fields{
				version: 1,
				r:       nil,
				typeMap: map[string]reflect.Type{
					"package.User":    reflect.TypeOf(User{}),
					"package.Account": reflect.TypeOf(&Account{}),
				},
			},
			args: args{
				pkg: "package.Account",
				data: map[interface{}]interface{}{
					"name": "ggwhite",
				},
			},
			want: &Account{
				Name: "ggwhite",
			},
			wantErr: false,
		},
		{
			name: "no mapping",
			fields: fields{
				version: 1,
				r:       nil,
				typeMap: map[string]reflect.Type{
					"package.User": reflect.TypeOf(User{}),
				},
			},
			args: args{
				pkg: "package.Account",
				data: map[interface{}]interface{}{
					"name": "ggwhite",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &DeserializerV1{
				version: tt.fields.version,
				r:       tt.fields.r,
				typeMap: tt.fields.typeMap,
			}
			got, err := i.BuildObject(tt.args.pkg, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeserializerV1.BuildObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeserializerV1.BuildObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeserializerV1_Read(t *testing.T) {
	type fields struct {
		version int
		r       io.Reader
		typeMap map[string]reflect.Type
	}
	tests := []struct {
		name    string
		fields  fields
		want    []interface{}
		wantErr bool
	}{
		{
			name: "reader not set",
			fields: fields{
				version: 1,
				r:       nil,
				typeMap: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "string",
			fields: fields{
				version: 1,
				r: bytes.NewReader([]byte{
					'r', 0x01, 0x00,
					'S', 0x00, 0x07, 'g', 'g', 'w', 'h', 'i', 't', 'e',
					'z',
				}),
				typeMap: nil,
			},
			want:    []interface{}{"ggwhite"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &DeserializerV1{
				version: tt.fields.version,
				r:       tt.fields.r,
				typeMap: tt.fields.typeMap,
			}
			got, err := i.Read()
			if (err != nil) != tt.wantErr {
				t.Errorf("DeserializerV1.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeserializerV1.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeserializerV1_Reset(t *testing.T) {
	type fields struct {
		version int
		r       io.Reader
		typeMap map[string]reflect.Type
	}
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "ok",
			fields: fields{
				version: 1,
				r:       nil,
				typeMap: nil,
			},
			args: args{
				r: bytes.NewReader(nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &DeserializerV1{
				version: tt.fields.version,
				r:       tt.fields.r,
				typeMap: tt.fields.typeMap,
			}
			i.Reset(tt.args.r)
		})
	}
}
