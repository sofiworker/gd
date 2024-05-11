# godog [![gd](gd.png)]()

[![GoDoc](https://pkg.go.dev/badge/github.com/chuck1024/gd?status.svg)](https://pkg.go.dev/github.com/chuck1024/gd@v1.7.17?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/chuck1024/gd)](https://goreportcard.com/report/github.com/chuck1024/gd)
[![Release](https://img.shields.io/github/v/release/chuck1024/gd.svg?style=flat-square)](https://github.com/chuck1024/gd/releases)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg)]()

"go" is the meaning of a dog in Chinese pronunciation, and dog's original intention is also a dog. So godog means "狗狗"
in Chinese, which is very cute.

## Quick start

```go
package main

import (
	"github.com/chuck1024/gd"
	"github.com/chuck1024/gd/net/dhttp"
	"github.com/chuck1024/gd/runtime/inject"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandlerHttp(c *gin.Context, req interface{}) (code int, message string, err error, ret string) {
	gd.Debug("httpServerTest req:%v", req)
	ret = "ok!!!"
	return http.StatusOK, "ok", nil, ret
}

func main() {
	d := gd.Default()
	inject.RegisterOrFail("httpServerInit", func(g *gin.Engine) error {
		r := g.Group("")
		r.Use(
			dhttp.GlFilter(),
			dhttp.GroupFilter(),
			dhttp.Logger("quick-start"),
		)

		d.HttpServer.GET(r, "test", HandlerHttp)

		if err := d.HttpServer.CheckHandle(); err != nil {
			return err
		}
		return nil
	})

	gd.SetConfig("Server", "httpPort", "10240")

	if err := d.Run(); err != nil {
		gd.Error("Error occurs, error = %s", err.Error())
		return
	}
}
```


## Docker Image

if you want to build docker images for your project, we recommand use [slim](https://github.com/slimtoolkit/slim) to reszie your images



`PS: More information can be obtained in the source code`


## License

gd is released under the [**MIT LICENSE**](http://opensource.org/licenses/mit-license.php).  
