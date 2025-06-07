package http

import "context"

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
	Param(string) string
	Query(string) string
	Body() []byte
	Send(int, []byte) error
	JSON(int, any) error
	Context() context.Context
	ResponseError(err error) error
	ResponseSuccess(data any) error
	HttpData() map[string]any
}
