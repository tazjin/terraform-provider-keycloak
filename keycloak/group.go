package keycloak

import (
	"fmt"
)

type Group struct {
	Id          string            `json:"id,omitempty"`
	Name        string            `json:"name"`
	Attributes  map[string]string `json:"attributes,omitempty"`
	RealmRoles  []string          `json:"realmRoles,omitempty"`
	ClientRoles map[string]string `json:"clientRoles,omitempty"`
	SubGroups   []Group           `json:"subgroups,omitempty"`
}

const (
	groupUri  = "%s/auth/admin/realms/%s/groups/%s"
	groupList = "%s/auth/admin/realms/%s/groups"
)

func (c *KeycloakClient) AddGroup(group *Group, realm string) (*Group, error) {
	url := fmt.Sprintf(groupList, c.url, realm)

	groupLocation, err := c.post(url, *group)
	if err != nil {
		return nil, err
	}

	var createdGroup Group
	err = c.get(groupLocation, &createdGroup)

	return &createdGroup, err
}

// Attempt to look up group by given group ID
func (c *KeycloakClient) GetGroup(groupId string, realm string) (*Group, error) {
	url := fmt.Sprintf(groupUri, c.url, realm, groupId)

	var group Group
	err := c.get(url, &group)

	return &group, err
}

// Attempt to update group
func (c *KeycloakClient) UpdateGroup(group *Group, realm string) error {
	url := fmt.Sprintf(groupUri, c.url, realm, group.Id)
	err := c.put(url, *group)

	if err != nil {
		return err
	}

	return nil
}

func (c *KeycloakClient) DeleteGroup(id string, realm string) error {
	url := fmt.Sprintf(groupUri, c.url, realm, id)
	return c.delete(url, nil)
}
