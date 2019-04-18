package hessian

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"reflect"
)

// InputV1 input stream for hessian 1.0 response
type InputV1 struct {
	version int
	peek    int
	r       io.Reader
	typeMap map[string]reflect.Type
}

func (i *InputV1) read(p []byte, begin int) ([]interface{}, int, error) {

	var ans []interface{}

	var j = begin

	for ; j < len(p); j++ {
		ch := p[j]
		switch ch {
		default:
		case 'c':
		case 'r':
		case 'f':
			params, n, err := i.read(p, j+1)
			log.Println(params)
			j = n
			if err != nil {
				return nil, j, err
			}
			return nil, j, fmt.Errorf("Call Service Error, Exception: %s, Message: %s", params[1].(string), params[3].(string))
		case 'z':
			break
		case 'T':
			ans = append(ans, true)
		case 'F':
			ans = append(ans, false)
		case 'N':
			ans = append(ans, nil)
		case 'S':
			// S b16 b8
			val1 := int(p[j+1])
			val2 := int(p[j+2])
			j += 2
			sl := val1<<8 + val2
			byt := make([]byte, sl)
			for k := 0; k < sl; k++ {
				j++
				byt[k] = p[j]
			}
			ans = append(ans, string(byt))
		case 'I':
			// I b32 b24 b16 b8
			val1 := int32(p[j+1])
			val2 := int32(p[j+2])
			val3 := int32(p[j+3])
			val4 := int32(p[j+4])
			val := val1<<24 + val2<<16 + val3<<8 + val4
			ans = append(ans, val)
			j += 4
		case 'L':
			// L b64 b56 b48 b40 b32 b24 b16 b8
			val1 := int64(p[j+1])
			val2 := int64(p[j+2])
			val3 := int64(p[j+3])
			val4 := int64(p[j+4])
			val5 := int64(p[j+5])
			val6 := int64(p[j+6])
			val7 := int64(p[j+7])
			val8 := int64(p[j+8])
			val := val1<<56 + val2<<48 + val3<<40 + val4<<32 + val5<<24 + val6<<16 + val7<<8 + val8
			ans = append(ans, val)
			j += 8
		case 'D':
			// D b64 b56 b48 b40 b32 b24 b16 b8
			val1 := uint64(p[j+1])
			val2 := uint64(p[j+2])
			val3 := uint64(p[j+3])
			val4 := uint64(p[j+4])
			val5 := uint64(p[j+5])
			val6 := uint64(p[j+6])
			val7 := uint64(p[j+7])
			val8 := uint64(p[j+8])
			val := val1<<56 + val2<<48 + val3<<40 + val4<<32 + val5<<24 + val6<<16 + val7<<8 + val8
			ans = append(ans, math.Float64frombits(val))
			j += 8
		case 'M':
			if p[j+1] == 't' {
				j++
				val1 := int(p[j+1])
				val2 := int(p[j+2])
				j += 2
				sl := val1<<8 + val2
				byt := make([]byte, sl)
				for k := 0; k < sl; k++ {
					j++
					byt[k] = p[j]
				}
				pkg := string(byt)
				if len(pkg) > 0 {
					m := make(map[string]interface{})
					params, n, err := i.read(p, j+1)
					j = n
					if err != nil {
						return nil, j, err
					}

					for k := 0; k < len(params); k++ {
						m[params[k].(string)] = params[k+1]
						k++
					}

					if subt, exist := i.typeMap[pkg]; exist {
						var t reflect.Type
						if subt.Kind() == reflect.Ptr {
							t = subt.Elem()
						} else {
							t = subt
						}
						subv := reflect.New(t)

						for k, l := 0, t.NumField(); k < l; k++ {
							if val, exist := m[t.Field(k).Tag.Get(tagName)]; exist {
								if val == nil {
									continue
								}
								subv.Elem().Field(k).Set(reflect.ValueOf(val))
							}
						}
						if subt.Kind() == reflect.Ptr {
							ans = append(ans, subv.Interface())
							continue
						}
						ans = append(ans, subv.Elem().Interface())
					}
				} else {
					m := make(map[interface{}]interface{})
					params, n, err := i.read(p, j+1)
					j = n
					if err != nil {
						return nil, j, err
					}
					for k := 0; k < len(params); k++ {
						m[params[k]] = params[k+1]
						k++
					}
					ans = append(ans, m)
				}
			}
		}
	}

	return ans, j, nil
}

//
func (i *InputV1) Read() ([]interface{}, error) {
	if i.r == nil {
		return nil, fmt.Errorf("io.Reader is not set yet")
	}

	p, err := ioutil.ReadAll(i.r)
	if err != nil {
		return nil, err
	}

	log.Println(p)

	ans, _, err := i.read(p, 0)
	return ans, err
}

// SetReader give reader to Input
func (i *InputV1) SetReader(r io.Reader) {
	i.r = r
}

// SetTypeMapping set type
func (i *InputV1) SetTypeMapping(name string, t reflect.Type) {
	i.typeMap[name] = t
}

// NewInputV1 create InputV1
func NewInputV1() *InputV1 {
	return &InputV1{
		version: 1,
		peek:    -1,
		typeMap: make(map[string]reflect.Type),
	}
}
