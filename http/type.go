package http

const DefaultServiceName = "local"

// for local context
type (
	localRouteRoles struct{}
)

// for global context
type (
	ContextClaimToken struct{} // a key to used to store claim token from access token
)
