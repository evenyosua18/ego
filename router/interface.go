package router

type RouteHandler func(ctx Context) error

type IRouter interface {
	Get(path string, handler RouteHandler)
	Post(path string, handler RouteHandler)
	Put(path string, handler RouteHandler)
	Delete(path string, handler RouteHandler)
	Patch(path string, handler RouteHandler)

	Group(prefix string) IRouter
}

type Context interface {
	Param(string) string
	Query(string) string
	Body() []byte
	Send(int, []byte) error
	JSON(int, any) error
}
