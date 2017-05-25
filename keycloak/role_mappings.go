package keycloak

import (
	"fmt"
)

type Role struct {
	Id                 string `json:"id"`
	Name               string `json:"name,omitempty"`
	ScopeParamRequired bool   `json:"scopeParamRequired"`
}

const (
	rolesUri          = "%s/auth/admin/realms/%s/users/%s/role-mappings/realm"
	availableRolesUri = "%s/auth/admin/realms/%s/users/%s/role-mappings/realm/available"
	compositeRolesUri = "%s/auth/admin/realms/%s/users/%s/role-mappings/realm/composite"
)

// Attempt to look up available roles for a given user ID
func (c *KeycloakClient) GetAvailableRolesForUser(userId string, realm string) ([]Role, error) {
	url := fmt.Sprintf(availableRolesUri, c.url, realm, userId)

	var roles []Role
	err := c.get(url, &roles)

	return roles, err
}

// Attempt to look up copmosite (effective) roles for a given user ID
func (c *KeycloakClient) GetCompositeRolesForUser(userId string, realm string) ([]Role, error) {
	url := fmt.Sprintf(compositeRolesUri, c.url, realm, userId)

	var roles []Role
	err := c.get(url, &roles)

	return roles, err
}

// Find a role for a given user based on the role ID.
// The idea is to hide the complexity of the randomly generated role IDs from the user.
// TODO: Evaluate whether this is the most sensible approach vs. some sort of data provider
func (c *KeycloakClient) FindRoleForUser(roles []Role, roleIdentifier string) (*Role, error) {
	var role Role
	for _, value := range roles {
		if value.Name == roleIdentifier || value.Id == roleIdentifier {
			role = value
		}
	}

	if role.Id == "" {
		return nil, fmt.Errorf("Role %s not found", roleIdentifier)
	}

	return &role, nil
}

// This attempts to add a Keycloak role to a user based after looking up the role ID from the available rolesUri.
func (c *KeycloakClient) AddRoleToUser(userId string, roleName string, realm string) (*Role, error) {
	roles, err := c.GetAvailableRolesForUser(userId, realm)
	if err != nil {
		return nil, err
	}

	role, err := c.FindRoleForUser(roles, roleName)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(rolesUri, c.url, realm, userId)
	body := []Role{*role}

	_, err = c.post(url, &body)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (c *KeycloakClient) RemoveRoleFromUser(userId string, role *Role, realm string) error {
	url := fmt.Sprintf(rolesUri, c.url, realm, userId)
	body := []Role{*role}

	err := c.delete(url, body)
	if err != nil {
		return err
	}

	return nil
}
