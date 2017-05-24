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
	clientUserUri   = "%s/auth/admin/realms/%s/clients/%s/service-account-user"
)

func (c *KeycloakClient) GetClient(id string, realm string) (*Client, error) {
	url := fmt.Sprintf(clientUri, c.url, realm, id)

	var client Client
	err := c.get(url, &client)

	if err != nil {
		return nil, err
	}

	return &client, nil
}

func (c *KeycloakClient) GetClientSecret(id string, realm string) (*ClientSecret, error) {
	url := fmt.Sprintf(clientSecretUri, c.url, realm, id)

	var secret ClientSecret
	err := c.get(url, &secret)

	if err != nil {
		return nil, err
	}

	return &secret, nil
}

// Attempt to create a Keycloak client and return the created client.
func (c *KeycloakClient) CreateClient(client *Client, realm string) (*Client, error) {
	url := fmt.Sprintf(clientList, c.url, realm)
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

func (c *KeycloakClient) UpdateClient(client *Client, realm string) error {
	url := fmt.Sprintf(clientUri, c.url, realm, client.Id)
	err := c.put(url, *client)

	if err != nil {
		return err
	}

	return nil
}

func (c *KeycloakClient) DeleteClient(id string, realm string) error {
	url := fmt.Sprintf(clientUri, c.url, realm, id)
	return c.delete(url)
}

func (c *KeycloakClient) GetClientServiceAccountUser(id string, realm string) (*User, error) {
	url := fmt.Sprintf(clientUserUri, c.url, realm, id)

	var user User
	err := c.get(url, &user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
