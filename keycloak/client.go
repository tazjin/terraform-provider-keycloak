package keycloak

import (
	"fmt"
)

// Client resource as documented in the Keycloak REST API docs.
// Some fields are not mapped here.
// http://www.keycloak.org/docs-api/3.1/rest-api/index.html#_clientrepresentation

type Client struct {
	Id                      string   `json:"id,omitempty"`
	ClientId                string   `json:"clientId"`
	Enabled                 bool     `json:"enabled"`
	ClientAuthenticatorType string   `json:"clientAuthenticatorType,omitempty"`
	RedirectUris            []string `json:"redirectUris"`
	Protocol                string   `json:"protocol,omitempty"`
	PublicClient            bool     `json:"publicClient"`
	BearerOnly              bool     `json:"bearerOnly"`
	ServiceAccountsEnabled  bool     `json:"serviceAccountsEnabled"`
}

type ClientSecret struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

const (
	clientUri       = "%s/auth/admin/realms/%s/clients/%s"
	clientList      = "%s/auth/admin/realms/%s/clients"
	clientSecretUri = "%s/auth/admin/realms/%s/clients/%s/client-secret"
)

func (c *KeycloakClient) GetClient(id string) (*Client, error) {
	url := fmt.Sprintf(clientUri, c.url, c.realm, id)

	var client Client
	err := c.get(url, &client)

	if err != nil {
		return nil, err
	}

	return &client, nil
}

func (c *KeycloakClient) GetClientSecret(id string) (*ClientSecret, error) {
	url := fmt.Sprintf(clientSecretUri, c.url, c.realm, id)

	var secret ClientSecret
	err := c.get(url, &secret)

	if err != nil {
		return nil, err
	}

	return &secret, nil
}

// Attempt to create a Keycloak client and return the created client.
func (c *KeycloakClient) CreateClient(client *Client) (*Client, error) {
	url := fmt.Sprintf(clientList, c.url, c.realm)
	clientLocation, err := c.post(url, *client)
	if err != nil {
		return nil, err
	}

	var createdClient Client
	err = c.get(clientLocation, &createdClient)

	if err != nil {
		return nil, err
	}

	return &createdClient, nil
}

func (c *KeycloakClient) UpdateClient(client *Client) (*Client, error) {
	url := fmt.Sprintf(clientUri, c.url, c.realm, client.Id)
	err := c.put(url, *client)

	if err != nil {
		return nil, err
	}

	updated, err := c.GetClient(client.Id)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (c *KeycloakClient) DeleteClient(id string) error {
	url := fmt.Sprintf(clientUri, c.url, c.realm, id)
	return c.delete(url)
}
