package http

type RouteConfig struct {
	MaxLimit            int
	MaxConnection       int
	MainPrefix          string
	ShowRegisteredRoute bool
	HtmlPath            string
	CORS                CORSConfig
	Doc                 DocumentationConfig
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
