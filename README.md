# Golang Hessian

[![Build Status](https://travis-ci.org/ggwhite/go-hessian.svg?branch=master)](https://travis-ci.org/ggwhite/go-hessian)
[![codecov](https://codecov.io/gh/ggwhite/go-hessian/branch/master/graph/badge.svg)](https://codecov.io/gh/ggwhite/go-hessian)
[![Go Report Card](https://goreportcard.com/badge/github.com/ggwhite/go-hessian)](https://goreportcard.com/report/github.com/ggwhite/go-hessian)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/ggwhite/go-hessian/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/ggwhite/go-hessian?status.svg)](https://godoc.org/github.com/ggwhite/go-hessian)

Golang Hessian can use hessian proxy to connect to hessian service.

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
    proxy := hessian.NewProxy(hessian.V1, addr, 5*time.Second)
    
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
    proxy := hessian.NewProxy(hessian.V1, addr, 5*time.Second)
    
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
    proxy := hessian.NewProxy(hessian.V1, addr, 5*time.Second)
    proxy.SetTypeMapping("lab.ggw.shs.service.User", reflect.TypeOf(&User{})) // or User{}
    
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

> Give `hessian.Package` to your struct and add tag `hessian` to let encoder know what package(ClassName) of your POJO,
> set `TypeMapping` to proxy let decoder know what type that you want convert to.
> 
> Mapping type can be a type of struct or a pointer of the struct.