# Golang Hessian

[![Build Status](https://travis-ci.org/ggwhite/go-hessian.svg?branch=master)](https://travis-ci.org/ggwhite/go-hessian)
[![codecov](https://codecov.io/gh/ggwhite/go-hessian/branch/master/graph/badge.svg)](https://codecov.io/gh/ggwhite/go-hessian)
[![Go Report Card](https://goreportcard.com/badge/github.com/ggwhite/go-hessian)](https://goreportcard.com/report/github.com/ggwhite/go-hessian)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/ggwhite/go-hessian/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/ggwhite/go-hessian?status.svg)](https://godoc.org/github.com/ggwhite/go-hessian)
[![Release](https://img.shields.io/github/release/ggwhite/go-hessian.svg?style=flat-square)](https://github.com/ggwhite/go-hessian/releases/latest)

Golang Hessian can use hessian proxy to connect to hessian service.

> Reference: 
> * Hessian Official Site: http://hessian.caucho.com
> * Hessian Java source: https://github.com/ebourg/hessian
> * Project test server: [simple-hessian-server](https://github.com/ggwhite/simple-hessian-server)

# How to use ?

### No argument method call

``` golang
package main

import (
    "log"
    "time"

    hessian "github.com/ggwhite/go-hessian"
)

func main() {
    var addr = "http://localhost:8080/simple"
    proxy, err := hessian.NewProxy(&hessian.ProxyConfig{
        Version: hessian.V1,
        URL:     addr,
    })
    if err != nil {
        panic(err)
    }
    
    args, err := proxy.Invoke("str")
    log.Println(args, err)
}
```

Result:
```
2019/04/18 16:20:52 [Hello] <nil>
```

> Create a proxy to your hessian service, invoke from given `method` name, it return a slice of interface and error.

### With argument method call

``` golang
package main

import (
    "log"
    "time"

    hessian "github.com/ggwhite/go-hessian"
)

func main() {
    var addr = "http://localhost:8080/simple"
    proxy, err := hessian.NewProxy(&hessian.ProxyConfig{
        Version: hessian.V1,
        URL:     addr,
    })
    if err != nil {
        panic(err)
    }
    
    log.Println(proxy.Invoke("strI2", "ggwhite", "this is message"))
    log.Println(args, err)
}
```

Result:
```
2019/04/18 16:43:18 [Hello[ggwhite], this is message] <nil>
```

### With struct argument method call

``` golang
package main

import (
    "log"
    "time"

    hessian "github.com/ggwhite/go-hessian"
)

type User struct {
    hessian.Package `hessian:"lab.ggw.shs.service.User"`
    Name            string      `hessian:"name"`
    Email           interface{} `hessian:"email"`
    Father          *User       `hessian:"father"`
}

func main() {
    var addr = "http://localhost:8080/simple"
    proxy, err := hessian.NewProxy(&hessian.ProxyConfig{
        Version: hessian.V1,
        URL:     addr,
    })
    if err != nil {
        panic(err)
    }
    proxy.RegisterType(reflect.TypeOf(User{}))
    
    ans, err := proxy.Invoke("obj")
    log.Println(ans[0], err)
    
    log.Println(proxy.Invoke("objI", &User{
        Name:  "ggwhite",
        Email: "ggw.chang@gmail.com",
    }))
}
```

Result:
```
2019/04/18 16:46:13 &{ ggwhite ggw.chang@gmail.com <nil>} <nil>
2019/04/18 16:46:13 [ggwhite] <nil>
```

> Give `hessian.Package` to your struct and add tag `hessian` to let proxy know what package(ClassName) of your POJO.
> 
> Mapping type can be a type of struct or a pointer of the struct.


## Supported

* Hessian v1.0 specification.
  > Ref: [hessian 1.0 spec](http://hessian.caucho.com/doc/hessian-1.0-spec.xtp)

|Type    |Serialize |Deserialize |Syntax                                      |
|:------:|:--------:|:----------:|--------------------------------------------|
|null    |Y         |Y           |N                                           |
|boolean |Y         |Y           |T/F                                         |
|int     |Y         |Y           |I b32 b24 b16 b8                            |
|long    |Y         |Y           |L b64 b56 b48 b40 b32 b24 b16 b8            |
|double  |Y         |Y           |D b64 b56 b48 b40 b32 b24 b16 b8            |
|date    |Y         |Y           |d b64 b56 b48 b40 b32 b24 b16 b8            |
|string  |Y         |Y           |S b16 b8 utf-8-data                         |
|xml     |N         |N           |                                            |
|binary  |Y         |Y           |B b16 b8 binary-data                        |
|list    |Y         |Y           |V type? length? object* z                   |
|map     |Y         |Y           |M t b16 b8 type-string (object, object)* z  |
|ref     |N         |N           |                                            |
|remote  |N         |N           |                                            |