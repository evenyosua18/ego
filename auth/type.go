package auth

type ClaimToken struct {
	AccessType      string   `json:"access_type"`
	AllowedServices []string `json:"allowed_services"`
	ClientID        string   `json:"client_id"`
	ExpiredAt       string   `json:"expired_at"`
	GrantType       string   `json:"grant_type"`
	Roles           []string `json:"roles"`
	Scopes          []string `json:"scopes"`
	UserID          int      `json:"user_id"`
}
