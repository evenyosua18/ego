package http

import "fmt"

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

func (m *HttpRouter) Group(prefix string) IHttpRouter {
	newPrefix := m.RouteGroup + prefix
	if len(newPrefix) > 1 && newPrefix[len(newPrefix)-1] == '/' {
		newPrefix = newPrefix[:len(newPrefix)-1]
	}
	return &HttpRouter{Calls: m.Calls, RouteGroup: newPrefix}
}

func (m *HttpRouter) Listen(port string) error {
	*m.Calls = append(*m.Calls, fmt.Sprintf("LISTEN %s", port))
	return nil
}

func (m *HttpRouter) fullPath(path string) string {
	if m.RouteGroup == "/" {
		return path
	}

	return m.RouteGroup + path
}
