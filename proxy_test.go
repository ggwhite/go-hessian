package hessian

import (
	"net/http"
	"reflect"
	"testing"
)

func TestProxyConfig_Validate(t *testing.T) {
	type fields struct {
		Version version
		URL     string
		TypeMap map[string]reflect.Type
		Client  *http.Client
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		want    *ProxyConfig
	}{
		{
			name: "Happy Pass",
			fields: fields{
				Version: V1,
				URL:     "http://localhost:8080/simple",
				TypeMap: make(map[string]reflect.Type),
				Client:  &http.Client{},
			},
			wantErr: false,
			want: &ProxyConfig{
				Version: V1,
				URL:     "http://localhost:8080/simple",
				TypeMap: make(map[string]reflect.Type),
				Client:  &http.Client{},
			},
		},
		{
			name: "Without URL",
			fields: fields{
				Version: V1,
				URL:     "",
				TypeMap: make(map[string]reflect.Type),
				Client:  &http.Client{},
			},
			wantErr: true,
			want: &ProxyConfig{
				Version: V1,
				URL:     "",
				TypeMap: make(map[string]reflect.Type),
				Client:  &http.Client{},
			},
		},
		{
			name: "Without TypeMap",
			fields: fields{
				Version: V1,
				URL:     "http://localhost:8080/simple",
				TypeMap: nil,
				Client:  &http.Client{},
			},
			wantErr: false,
			want: &ProxyConfig{
				Version: V1,
				URL:     "http://localhost:8080/simple",
				TypeMap: make(map[string]reflect.Type),
				Client:  &http.Client{},
			},
		},
		{
			name: "Without Client",
			fields: fields{
				Version: V1,
				URL:     "http://localhost:8080/simple",
				TypeMap: make(map[string]reflect.Type),
				Client:  nil,
			},
			wantErr: false,
			want: &ProxyConfig{
				Version: V1,
				URL:     "http://localhost:8080/simple",
				TypeMap: make(map[string]reflect.Type),
				Client:  &http.Client{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ProxyConfig{
				Version: tt.fields.Version,
				URL:     tt.fields.URL,
				TypeMap: tt.fields.TypeMap,
				Client:  tt.fields.Client,
			}
			if err := c.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ProxyConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.want, c) {
				t.Errorf("ProxyConfig.Validate() not right = %v, want %v", c, tt.want)
			}
		})
	}
}

func TestProxy_RegisterType(t *testing.T) {
	type User struct {
		Package `hessian:"path"`
		Name    string
	}
	type User2 struct {
		Package `hessian:"path2"`
		Name    string
	}
	type Account struct {
		Name string
	}
	type fields struct {
		conf         *ProxyConfig
		client       *http.Client
		serializer   Serializer
		deserializer Deserializer
	}
	type args struct {
		t reflect.Type
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		wantMap map[string]reflect.Type
	}{
		{
			name: "Happy Pass",
			fields: fields{
				conf: &ProxyConfig{
					Version: V1,
					URL:     "http://localhost:8080/simple",
					TypeMap: make(map[string]reflect.Type),
					Client:  &http.Client{},
				},
				client:       &http.Client{},
				serializer:   NewSerializerV1(),
				deserializer: NewDeserializerV1(),
			},
			args: args{
				t: reflect.TypeOf(&User{}),
			},
			wantErr: false,
			wantMap: map[string]reflect.Type{
				"path": reflect.TypeOf(&User{}),
			},
		},
		{
			name: "Happy Pass",
			fields: fields{
				conf: &ProxyConfig{
					Version: V1,
					URL:     "http://localhost:8080/simple",
					TypeMap: make(map[string]reflect.Type),
					Client:  &http.Client{},
				},
				client:       &http.Client{},
				serializer:   NewSerializerV1(),
				deserializer: NewDeserializerV1(),
			},
			args: args{
				t: reflect.TypeOf(User{}),
			},
			wantErr: false,
			wantMap: map[string]reflect.Type{
				"path": reflect.TypeOf(User{}),
			},
		},
		{
			name: "Exist Type",
			fields: fields{
				conf: &ProxyConfig{
					Version: V1,
					URL:     "http://localhost:8080/simple",
					TypeMap: map[string]reflect.Type{
						"path": reflect.TypeOf(User{}),
					},
					Client: &http.Client{},
				},
				client:       &http.Client{},
				serializer:   NewSerializerV1(),
				deserializer: NewDeserializerV1(),
			},
			args: args{
				t: reflect.TypeOf(User2{}),
			},
			wantErr: false,
			wantMap: map[string]reflect.Type{
				"path":  reflect.TypeOf(User{}),
				"path2": reflect.TypeOf(User2{}),
			},
		},
		{
			name: "Type withe Package",
			fields: fields{
				conf: &ProxyConfig{
					Version: V1,
					URL:     "http://localhost:8080/simple",
					TypeMap: make(map[string]reflect.Type),
					Client:  &http.Client{},
				},
				client:       &http.Client{},
				serializer:   NewSerializerV1(),
				deserializer: NewDeserializerV1(),
			},
			args: args{
				t: reflect.TypeOf(Account{}),
			},
			wantErr: true,
			wantMap: make(map[string]reflect.Type),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Proxy{
				conf:         tt.fields.conf,
				client:       tt.fields.client,
				serializer:   tt.fields.serializer,
				deserializer: tt.fields.deserializer,
			}
			if err := c.RegisterType(tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("Proxy.RegisterType() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.fields.conf.TypeMap, tt.wantMap) {
				t.Errorf("Proxy.RegisterType() map = %v, wantMap %v", tt.fields.conf.TypeMap, tt.wantMap)
			}
		})
	}
}

func TestNewProxy(t *testing.T) {
	type args struct {
		c *ProxyConfig
	}
	tests := []struct {
		name    string
		args    args
		want    *Proxy
		wantErr bool
	}{
		{
			name: "Happy Pass",
			args: args{
				c: &ProxyConfig{
					Version: V1,
					URL:     "http://localhost:8080/simple",
					TypeMap: make(map[string]reflect.Type),
					Client:  &http.Client{},
				},
			},
			want: &Proxy{
				conf: &ProxyConfig{
					Version: V1,
					URL:     "http://localhost:8080/simple",
					TypeMap: make(map[string]reflect.Type),
					Client:  &http.Client{},
				},
				client:       &http.Client{},
				serializer:   NewSerializerV1(),
				deserializer: NewDeserializerV1(),
			},
			wantErr: false,
		},
		{
			name: "Empty URL",
			args: args{
				c: &ProxyConfig{
					Version: V1,
					URL:     "",
					TypeMap: make(map[string]reflect.Type),
					Client:  &http.Client{},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Empty TypeMap",
			args: args{
				c: &ProxyConfig{
					Version: V1,
					URL:     "http://localhost:8080/simple",
					TypeMap: nil,
					Client:  &http.Client{},
				},
			},
			want: &Proxy{
				conf: &ProxyConfig{
					Version: V1,
					URL:     "http://localhost:8080/simple",
					TypeMap: make(map[string]reflect.Type),
					Client:  &http.Client{},
				},
				client:       &http.Client{},
				serializer:   NewSerializerV1(),
				deserializer: NewDeserializerV1(),
			},
			wantErr: false,
		},
		{
			name: "Empty Client",
			args: args{
				c: &ProxyConfig{
					Version: V1,
					URL:     "http://localhost:8080/simple",
					TypeMap: make(map[string]reflect.Type),
					Client:  nil,
				},
			},
			want: &Proxy{
				conf: &ProxyConfig{
					Version: V1,
					URL:     "http://localhost:8080/simple",
					TypeMap: make(map[string]reflect.Type),
					Client:  &http.Client{},
				},
				client:       &http.Client{},
				serializer:   NewSerializerV1(),
				deserializer: NewDeserializerV1(),
			},
			wantErr: false,
		},
		{
			name: "Wrong Version",
			args: args{
				c: &ProxyConfig{
					Version: version(123),
					URL:     "http://localhost:8080/simple",
					TypeMap: make(map[string]reflect.Type),
					Client:  nil,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Wrong Version V2",
			args: args{
				c: &ProxyConfig{
					Version: V2,
					URL:     "http://localhost:8080/simple",
					TypeMap: make(map[string]reflect.Type),
					Client:  nil,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProxy(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProxy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProxy() = %v, want %v", got, tt.want)
			}
		})
	}
}
