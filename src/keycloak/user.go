package keycloak

import "fmt"

type User struct {
	Id string `json:"id"`
	Username string `json:"username"`
	Enabled bool `json:"enabled"`
	FirstName string `json:"firstName,omitempty"`
	LastName string `json:"lastName,omitempty"`
	Email string `json:"email"`
	RequiredActions []string `json:"requiredActions,omitempty"`
}

const (
	userUri          = "%s/auth/admin/realms/%s/users/%s"
	userList          = "%s/auth/admin/realms/%s/users"
)

func (c *KeycloakClient) AddUser(user *User, realm string) (*User, error) {
	url := fmt.Sprintf(userList, c.url, realm)
	userLocation, err := c.post(url, *user)
	if err != nil {
		return nil, err
	}

	var createdUser User
	err = c.get(userLocation, &createdUser)

	return &createdUser, err
}


// Attempt to look up user by given user ID
func (c *KeycloakClient) GetUser(userId string, realm string) (*User, error) {
	url := fmt.Sprintf(userUri, c.url, realm, userId)

	var user User
	err := c.get(url, &user)

	return &user, err
}


// Attempt to update user
func (c *KeycloakClient) UpdateUser(user *User, realm string) (error) {
	url := fmt.Sprintf(userUri, c.url, realm, user.Id)
	err := c.put(url, *user)

	if err != nil {
		return err
	}

	return nil
}


func (c *KeycloakClient) DeleteUser(id string, realm string) error {
	url := fmt.Sprintf(userUri, c.url, realm, id)
	return c.delete(url, nil)
}