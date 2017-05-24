package keycloak

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	IdToken     string `json:"id_token"`
}

const (
	tokenEndpoint   = "%s/auth/realms/%s/protocol/openid-connect/token"
	formContentType = "application/x-www-form-urlencoded"
	loginBody       = "grant_type=client_credentials"
)

// Attempt to login to Keycloak with the provided information.
func Login(id string, secret string, baseUrl string, realm string) (*KeycloakClient, error) {
	url := fmt.Sprintf(tokenEndpoint, baseUrl, realm)

	req, _ := http.NewRequest("POST", url, bytes.NewBufferString(loginBody))
	req.Header.Set("Authorization", createBasicAuthorizationHeader(id, secret))
	req.Header.Set("Content-Type", formContentType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Keycloak login failed: %s (%d)", string(body), resp.StatusCode)
	}

	var t tokenResponse
	err = json.Unmarshal(body, &t)
	if err != nil {
		return nil, err
	}

	client := &KeycloakClient{
		token: t.AccessToken,
		url:   baseUrl,
		realm: realm,
	}
	return client, nil
}

func createBasicAuthorizationHeader(id string, secret string) string {
	input := fmt.Sprintf("%s:%s", id, secret)
	encoded := base64.StdEncoding.EncodeToString([]byte(input))
	return fmt.Sprintf("Basic %s", encoded)
}
