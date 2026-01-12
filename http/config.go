package http

type RouteConfig struct {
	MaxLimit            int
	MainPrefix          string
	ShowRegisteredRoute bool
	HtmlPath            string
	CORS                CORSConfig
}

type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
}
