package keycloak

// An authenticated Keycloak API client
type Client struct {
	token string
	url   string
}
