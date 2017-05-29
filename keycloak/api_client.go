package keycloak

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// An authenticated Keycloak API client
type KeycloakClient struct {
	token string
	url   string
}

// A function that mimics the default HTTP client 'Do' but authenticates all requests.
func (c *KeycloakClient) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.token)
	return http.DefaultClient.Do(req)
}

// Attempt to perform a GET request to the specified URL (with authentication).
// The result is decoded
// Go's type system is not able to type-check this function, so be careful - footguns ahead.
func (c *KeycloakClient) get(url string, v interface{}) error {
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := c.do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return fmt.Errorf("Could not get %s: %s (%d)", url, string(body), resp.StatusCode)
	}

	err = json.Unmarshal(body, v)

	if err != nil {
		return err
	}

	return nil
}

// Attempts to POST (create) a resource to Keycloak and returns the resource location.
func (c *KeycloakClient) post(url string, v interface{}) (string, error) {
	reqBody, _ := json.Marshal(v)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.do(req)

	if err != nil {
		return "", err
	}

	if resp.StatusCode != 201 && resp.StatusCode != 204 {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("Could not create resource: %s (%d)", string(body), resp.StatusCode)
	}

	return resp.Header.Get("Location"), nil
}

func (c *KeycloakClient) put(url string, v interface{}) error {
	reqBody, _ := json.Marshal(v)
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Could not update resource: %s (%d)", string(body), resp.StatusCode)
	}

	return nil
}

func (c *KeycloakClient) delete(url string, body interface{}) error {
	var req *http.Request
	if body != nil {
		bodyJson, _ := json.Marshal(body)
		req, _ = http.NewRequest("DELETE", url, bytes.NewBuffer(bodyJson))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest("DELETE", url, nil)
	}

	resp, err := c.do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Could not delete resource: %s (%d)", string(body), resp.StatusCode)
	}

	return nil
}
