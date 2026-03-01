package http

import (
	"context"
	"fmt"
	"net/http"
)

type HttpRouter struct {
	Calls      *[]string
	RouteGroup string
}

func (m *HttpRouter) Get(path string, _ RouteHandler, opts ...RouterFuncOption) {
	*m.Calls = append(*m.Calls, fmt.Sprintf("GET %s", m.fullPath(path)))
}

func (m *HttpRouter) Post(path string, _ RouteHandler, opts ...RouterFuncOption) {
	*m.Calls = append(*m.Calls, fmt.Sprintf("POST %s", m.fullPath(path)))
}

func (m *HttpRouter) Put(path string, _ RouteHandler, opts ...RouterFuncOption) {
	*m.Calls = append(*m.Calls, fmt.Sprintf("PUT %s", m.fullPath(path)))
}

func (m *HttpRouter) Delete(path string, _ RouteHandler, opts ...RouterFuncOption) {
	*m.Calls = append(*m.Calls, fmt.Sprintf("DELETE %s", m.fullPath(path)))
}

func (m *HttpRouter) Patch(path string, _ RouteHandler, opts ...RouterFuncOption) {
	*m.Calls = append(*m.Calls, fmt.Sprintf("PATCH %s", m.fullPath(path)))
}

func (m *HttpRouter) Use(args ...any) IHttpRouter {
	return m
}

func (m *HttpRouter) Group(prefix string, handlers ...any) IHttpRouter {
	newPrefix := m.RouteGroup + prefix
	if len(newPrefix) > 1 && newPrefix[len(newPrefix)-1] == '/' {
		newPrefix = newPrefix[:len(newPrefix)-1]
	}

	for _, handler := range handlers {
		if name, ok := handler.(string); ok {
			*m.Calls = append(*m.Calls, fmt.Sprintf("GROUP %s AS %s", newPrefix, name))
		}
	}

	return &HttpRouter{Calls: m.Calls, RouteGroup: newPrefix}
}

func (m *HttpRouter) Listen(port string) error {
	*m.Calls = append(*m.Calls, fmt.Sprintf("LISTEN %s", port))
	return nil
}

func (m *HttpRouter) Shutdown() error {
	*m.Calls = append(*m.Calls, "SHUTDOWN")
	return nil
}

func (m *HttpRouter) ShutdownWithContext(ctx context.Context) error {
	*m.Calls = append(*m.Calls, "SHUTDOWN WITH CONTEXT")
	return nil
}

func (m *HttpRouter) ActiveConnections() int {
	return 0
}

func (m *HttpRouter) Test(req *http.Request) (*http.Response, error) {
	*m.Calls = append(*m.Calls, "TEST")
	return nil, nil
}

func (m *HttpRouter) fullPath(path string) string {
	if m.RouteGroup == "/" {
		return path
	}

	return m.RouteGroup + path
}
