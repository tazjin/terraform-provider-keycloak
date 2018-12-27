package keycloak

import (
	"fmt"
)

type UserGroupMap struct {
	GroupId string   `json:"groupId"`
	UserIds []string `json:"userIds"`
	Realm   string   `json:"realm"`
}

const (
	userGroupsUri = "%s/auth/admin/realms/%s/users/%s/groups/%s"
	getUsersUri   = "%s/auth/admin/realms/%s/groups/%s/members"
)

func (c *KeycloakClient) GetUsersInGroup(groupId string, realm string) (*UserGroupMap, error) {
	var users []User
	url := fmt.Sprintf(getUsersUri, c.url, realm, groupId)
	err := c.get(url, &users)

	if err != nil {
		return nil, err
	}

	ug := UserGroupMap{}
	ug.GroupId = groupId
	ug.Realm = realm
	ug.UserIds = []string{}
	for index := 0; index < len(users); index++ {
		ug.UserIds = append(ug.UserIds, users[index].Id)
	}
	return &ug, err
}

// This attempts to add a Keycloak role to a user based after looking up the role ID from the available rolesUri.
func (c *KeycloakClient) AddUsersToGroup(userIds []string, groupId string, realm string) error {
	for index := 0; index < len(userIds); index++ {
		url := fmt.Sprintf(userGroupsUri, c.url, realm, userIds[index], groupId)
		//body := []UserGroupMap{}
		err := c.put(url, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *KeycloakClient) RemoveUsersFromGroup(userIds []string, groupId string, realm string) error {
	for index := 0; index < len(userIds); index++ {
		url := fmt.Sprintf(userGroupsUri, c.url, realm, userIds[index], groupId)
		//body := []UserGroupMap{}
		err := c.delete(url, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
