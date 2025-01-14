package main

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"golang.org/x/time/rate"
)

type HelloReq struct {
	g.Meta `path:"/hello" method:"get" sort:"1"`
	Name   string `v:"required" dc:"Your name"`
}

type HelloRes struct {
	Reply string `dc:"Reply content"`
}

type Hello struct{}

func (Hello) Say(ctx context.Context, req *HelloReq) (res *HelloRes, err error) {
	g.Log().Debugf(ctx, `receive say: %+v`, req)
	res = &HelloRes{
		Reply: fmt.Sprintf(`Hi %s`, req.Name),
	}
	return
}

var limiter = rate.NewLimiter(rate.Limit(10), 1) // 10 request per second

func Limiter(r *ghttp.Request) {
	if !limiter.Allow() {
		r.Response.WriteStatusExit(429)
		r.ExitAll()
	}
	r.Middleware.Next()
}

// curl "http://127.0.0.1:8080/hello?name=world"
func main() {
	s := g.Server()
	s.Use(Limiter, ghttp.MiddlewareHandlerResponse)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Bind(
			new(Hello),
		)
	})
	s.SetPort(8080)
	s.Run()
}
