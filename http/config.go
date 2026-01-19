package http

type RouteConfig struct {
	MaxConnection       int
	MainPrefix          string
	ShowRegisteredRoute bool
	HtmlPath            string
	CORS                CORSConfig
	Doc                 DocumentationConfig
	RateLimit           RateLimitConfig
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
