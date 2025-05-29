package http

import (
	"reflect"
	"testing"
)

func TestRouterFuncOption(t *testing.T) {
	tests := []struct {
		name     string
		initial  RouterOptions
		option   RouterFuncOption
		expected RouterOptions
	}{
		{
			name:    "set full router option with roles",
			initial: RouterOptions{},
			option:  SetRouterFullOption(RouterOptions{Roles: []string{"admin", "user"}}),
			expected: RouterOptions{
				Roles: []string{"admin", "user"},
			},
		},
		{
			name:    "set router roles option",
			initial: RouterOptions{},
			option:  SetRouterRolesOption([]string{"editor", "viewer"}),
			expected: RouterOptions{
				Roles: []string{"editor", "viewer"},
			},
		},
		{
			name:    "append roles to existing",
			initial: RouterOptions{Roles: []string{"admin"}},
			option:  SetRouterRolesOption([]string{"user"}),
			expected: RouterOptions{
				Roles: []string{"admin", "user"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := tt.initial
			tt.option(&opt)

			if !reflect.DeepEqual(opt, tt.expected) {
				t.Errorf("expected %+v, got %+v", tt.expected, opt)
			}
		})
	}
}
