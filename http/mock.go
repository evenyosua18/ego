package http

import "fmt"

type MockRouter struct {
	Calls      *[]string
	RouteGroup string
}

func (m *MockRouter) Get(path string, _ RouteHandler, opts ...RouterFuncOption) {
	*m.Calls = append(*m.Calls, fmt.Sprintf("GET %s", m.fullPath(path)))
}

func (m *MockRouter) Post(path string, _ RouteHandler, opts ...RouterFuncOption) {
	*m.Calls = append(*m.Calls, fmt.Sprintf("POST %s", m.fullPath(path)))
}

func (m *MockRouter) Put(path string, _ RouteHandler, opts ...RouterFuncOption) {
	*m.Calls = append(*m.Calls, fmt.Sprintf("PUT %s", m.fullPath(path)))
}

func (m *MockRouter) Delete(path string, _ RouteHandler, opts ...RouterFuncOption) {
	*m.Calls = append(*m.Calls, fmt.Sprintf("DELETE %s", m.fullPath(path)))
}

func (m *MockRouter) Patch(path string, _ RouteHandler, opts ...RouterFuncOption) {
	*m.Calls = append(*m.Calls, fmt.Sprintf("PATCH %s", m.fullPath(path)))
}

func (m *MockRouter) Group(prefix string) IHttpRouter {
	newPrefix := m.RouteGroup + prefix
	if len(newPrefix) > 1 && newPrefix[len(newPrefix)-1] == '/' {
		newPrefix = newPrefix[:len(newPrefix)-1]
	}
	return &MockRouter{Calls: m.Calls, RouteGroup: newPrefix}
}

func (m *MockRouter) Listen(port string) error {
	*m.Calls = append(*m.Calls, fmt.Sprintf("LISTEN %s", port))
	return nil
}

func (m *MockRouter) fullPath(path string) string {
	if m.RouteGroup == "/" {
		return path
	}

	return m.RouteGroup + path
}
