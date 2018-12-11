package keycloak

import "fmt"

type RoleRepresentation struct {
	Id             string  `json:"id"`
	Description    string  `json:"description"`
	Name           string  `json:"name"`
}

type CompositeRoleReference struct {
	Id             string  `json:"id"`
}

const (
	clientRolesUri  = "%s/auth/admin/realms/%s/clients/%s/roles"
	clientRoleUri   = "%s/auth/admin/realms/%s/clients/%s/roles/%s"
	clientRolesCompositesUri = "%s/auth/admin/realms/%s/clients/%s/roles/%s/composites"
)

func (c *KeycloakClient) GetClientRole(clientId string, realm string, representation *RoleRepresentation) (*RoleRepresentation, error) {
	var role RoleRepresentation
	roleUrl := fmt.Sprintf(clientRoleUri, c.url, realm, clientId, representation.Name)
	err := c.get(roleUrl, &role)
	return &role, err
}

func (c *KeycloakClient) CreateClientRole(clientId string, realm string, representation *RoleRepresentation) (*RoleRepresentation, error) {
	url := fmt.Sprintf(clientRolesUri, c.url, realm, clientId)


	_, err := c.post(url, representation)
	if err != nil {
		return nil, err
	}

	var createdRole RoleRepresentation
	roleUrl := fmt.Sprintf(clientRoleUri, c.url, realm, clientId, representation.Name)
	err = c.get(roleUrl, &createdRole)

	return &createdRole, err
}

func (c *KeycloakClient) UpdateClientRole(clientId string, realm string, representation *RoleRepresentation) error {
	url := fmt.Sprintf(clientRoleUri, c.url, realm, clientId, representation.Name)
	err := c.put(url, representation)
	return err
}

func (c *KeycloakClient) DeleteClientRole(clientId string, realm string, representation *RoleRepresentation) error {
	url := fmt.Sprintf(clientRoleUri, c.url, realm, clientId, representation.Name)
	err := c.delete(url, representation)
	return err
}

func (c *KeycloakClient) GetCompositeRoles(clientId string, realm string, representation *RoleRepresentation) ([]string, error) {
	var roles []CompositeRoleReference
	roleUrl := fmt.Sprintf(clientRolesCompositesUri, c.url, realm, clientId, representation.Name)
	err := c.get(roleUrl, &roles)
	
	var compositeRoleIds []string
	for _, value := range roles {
		compositeRoleIds = append(compositeRoleIds, value.Id)
	}
	
	return compositeRoleIds, err
}

func (c *KeycloakClient) AddRolesToCompositeRole(clientId string, realm string, representation *RoleRepresentation, roleIds []string) error {
	url := fmt.Sprintf(clientRolesCompositesUri, c.url, realm, clientId, representation.Name)
	_, err := c.post(url, toCompositeRoleRepresentation(roleIds))
	return err
}

func (c *KeycloakClient) RemoveRolesFromCompositeRole(clientId string, realm string, representation *RoleRepresentation, roleIds []string) error {
	url := fmt.Sprintf(clientRolesCompositesUri, c.url, realm, clientId, representation.Name)
	err := c.delete(url, toCompositeRoleRepresentation(roleIds))
	return err
}

func toCompositeRoleRepresentation(roleIds []string) []*CompositeRoleReference {
	var roles []*CompositeRoleReference
	for _, value := range roleIds {
		role := CompositeRoleReference{
			Id: value,
		}
		roles = append(roles, &role)
	}
	return roles
}