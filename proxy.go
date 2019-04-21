package hessian

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
)

type version int

// Versions
const (
	V1 version = iota
	V2
)

// ProxyConfig config for NewProxy
type ProxyConfig struct {
	Version version
	URL     string
	TypeMap map[string]reflect.Type
	Client  *http.Client
}

// Validate Proxy Config
func (c *ProxyConfig) Validate() error {
	if len(c.URL) == 0 {
		return fmt.Errorf("Proxy Config: URL is required")
	}

	if c.TypeMap == nil {
		c.TypeMap = make(map[string]reflect.Type)
	}

	if c.Client == nil {
		c.Client = &http.Client{}
	}

	return nil
}

// Proxy hessian proxy
type Proxy struct {
	conf         *ProxyConfig
	client       *http.Client
	serializer   Serializer
	deserializer Deserializer
}

// Invoke input method name and arguments, it will send request to server, and parse response to interface
func (c *Proxy) Invoke(m string, args ...interface{}) ([]interface{}, error) {

	c.serializer.Flush()

	if err := c.serializer.Call(m, args...); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.conf.URL, c.serializer.Reader())
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "x-application/hessian")
	req.Header.Set("Accept-Encoding", "deflate")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("status error", resp.StatusCode)
	}

	c.deserializer.Reset(resp.Body)

	ans, err := c.deserializer.Read()
	if err != nil {
		return nil, err
	}

	return ans, nil
}

// RegisterType input a type with Package field, register type mapping
func (c *Proxy) RegisterType(t reflect.Type) error {
	var tmp reflect.Type
	var pkg string

	if t.Kind() == reflect.Ptr {
		tmp = t.Elem()
	} else {
		tmp = t
	}

	for i, l := 0, tmp.NumField(); i < l; i++ {
		if tmp.Field(i).Type == reflect.TypeOf(Package("")) {
			pkg = tmp.Field(i).Tag.Get(tagName)
			break
		}
	}

	if len(pkg) == 0 {
		return fmt.Errorf("input type is without Package field")
	}

	c.conf.TypeMap[pkg] = t

	c.serializer.SetTypeMap(c.conf.TypeMap)
	c.deserializer.SetTypeMap(c.conf.TypeMap)

	return nil
}

// NewProxy create a proxy of the hessian service
func NewProxy(c *ProxyConfig) (*Proxy, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	switch c.Version {
	default:
		return nil, fmt.Errorf("Please set proxy version is V1 or V2")
	case V1:
		return &Proxy{
			conf:         c,
			client:       c.Client,
			serializer:   NewSerializerV1(),
			deserializer: NewDeserializerV1(),
		}, nil
	case V2:
		return nil, fmt.Errorf("Hessian V2.0 is unsupported")
	}
}
