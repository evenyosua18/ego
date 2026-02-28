package http

import "time"

type RouteConfig struct {
	MaxConnection       int
	MainPrefix          string
	ShowRegisteredRoute bool
	HtmlPath            string
	CORS                CORSConfig
	Doc                 DocumentationConfig
	RateLimit           RateLimitConfig
	DisableAuthChecker  bool

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
}

type DocumentationConfig struct {
	Path string
}

// TODO: integrate with redis for support distributed rate limit
type RateLimitConfig struct {
	MaxLimit int
}
