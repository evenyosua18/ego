package http

import (
	"context"
	"time"
)

type RouteHandler func(ctx Context) error

type IHttpRouter interface {
	Get(path string, handler RouteHandler, opts ...RouterFuncOption)
	Post(path string, handler RouteHandler, opts ...RouterFuncOption)
	Put(path string, handler RouteHandler, opts ...RouterFuncOption)
	Delete(path string, handler RouteHandler, opts ...RouterFuncOption)
	Patch(path string, handler RouteHandler, opts ...RouterFuncOption)

	Group(prefix string) IHttpRouter
	Listen(port string) error
}

type Context interface {
	Context() context.Context
	Param(string) string
	Query(string) string
	BindQuery(destination any) error
	Body() []byte
	RequestBody(res any) error

	Cookie(name string, defaultValue ...string) string
	SetCookie(name, value string, expiredAt time.Time)

	ResponseRedirect(to string) error

	Send(int, []byte) error
	JSON(int, any) error
	ResponseError(err error) error
	ResponseSuccess(data any) error

	HttpData() (map[string]any, string)

	Render(name string, data map[string]any) error
	FormValue(key string) string
}
