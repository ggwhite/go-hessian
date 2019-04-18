package hessian

import (
	"log"
	"net/http"
	"reflect"
	"time"
)

type version int

// Versions
const (
	V1 version = iota
	V2
)

// Proxy hessian proxy
type Proxy struct {
	v       version
	url     string
	client  *http.Client
	o       Output
	i       Input
	typeMap map[string]reflect.Type
}

// Invoke input method name and arguments, it will send request to server, and parse response to interface
func (c *Proxy) Invoke(m string, args ...interface{}) ([]interface{}, error) {

	c.o.Flush()

	if err := c.o.Call(m, args...); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.url, c.o.Reader())
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "x-application/hessian")
	req.Header.Set("Accept-Encoding", "deflate")

	resp, err := c.client.Do(req)
	log.Println(resp, err)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("status error", resp.StatusCode)
	}

	c.i.SetReader(resp.Body)

	ans, err := c.i.Read()
	if err != nil {
		return nil, err
	}

	return ans, nil
}

// Input get InputStream from the proxy
func (c *Proxy) Input() Input {
	return c.i
}

// SetTypeMapping set type mapping
func (c *Proxy) SetTypeMapping(name string, t reflect.Type) {
	c.o.SetTypeMapping(name, t)
	c.i.SetTypeMapping(name, t)
}

// NewProxy create hessian proxy
func NewProxy(v version, url string, timeout time.Duration) *Proxy {
	switch v {
	default:
		return nil
	case V1:
		return &Proxy{
			v:   v,
			url: url,
			client: &http.Client{
				Timeout: timeout,
			},
			o: NewOutputV1(),
			i: NewInputV1(),
		}
	case V2:
		return nil
	}
}
