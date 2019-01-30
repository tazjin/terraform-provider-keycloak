package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tazjin/terraform-provider-keycloak/keycloak"
	"log"
	"strings"
)

func resourceClientRole() *schema.Resource {
	return &schema.Resource{
		// API methods
		Read:   schema.ReadFunc(resourceClientRoleRead),
		Create: schema.CreateFunc(resourceClientRoleCreate),
		Update: schema.UpdateFunc(resourceClientRoleUpdate),
		Delete: schema.DeleteFunc(resourceClientRoleDelete),

		// Keycloak clients are importable by ID, but the realm must also be provided by the user.
		Importer: &schema.ResourceImporter{
			State: importClientRoleHelper,
		},

		Schema: map[string]*schema.Schema{
			"realm": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"composite_role_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func importClientRoleHelper(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	split := strings.Split(d.Id(), ".")

	if len(split) != 3 {
		return nil, fmt.Errorf("Import ID must be specified as '${realm}.${client_id}.${resource_id}'")
	}

	realm := split[0]
	client_id := split[1]
	role_name := split[2]

	d.Partial(true)
	d.Set("realm", realm)
	d.Set("client_id", client_id)
	d.Set("name", role_name)

	apiClient := m.(*keycloak.KeycloakClient)
	readRole, err := apiClient.GetClientRole(client_id, realm, role_name)
	if err != nil {
		return nil, err
	}

	d.Set("description", readRole.Description)
	d.SetId(readRole.Id)

	d.Partial(false)
	return []*schema.ResourceData{d}, nil
}

func resourceDataToRoleRepresentation(d *schema.ResourceData) *keycloak.RoleRepresentation {
	c := keycloak.RoleRepresentation{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	if !d.IsNewResource() {
		c.Id = d.Id()
	}

	return &c
}

func resourceClientRoleRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*keycloak.KeycloakClient)
	d.Partial(true)
	readRole, err := apiClient.GetClientRole(clientId(d), realm(d), d.Get("name").(string))
	if err != nil {
		return err
	}

	d.Set("name", readRole.Name)
	d.Set("description", readRole.Description)
	d.SetId(readRole.Id)

	roleIds := getCompositeRoleIds(d)
	if len(roleIds) > 0 {
		compositeRoleIds, err := apiClient.GetCompositeRoles(clientId(d), realm(d), resourceDataToRoleRepresentation(d))
		if err != nil {
			return err
		}

		setCompositeRoleIds(compositeRoleIds, d)
	}
	d.Partial(false)
	return nil
}

func resourceClientRoleCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*keycloak.KeycloakClient)
	d.Partial(true)
	createdRole, err := apiClient.CreateClientRole(clientId(d), realm(d), resourceDataToRoleRepresentation(d))
	if err != nil {
		log.Printf("[WARN] Error when creating client role: %s", err.Error())
		return err
	}

	d.Set("name", createdRole.Name)
	d.Set("description", createdRole.Description)
	d.SetId(createdRole.Id)

	roleIds := getCompositeRoleIds(d)
	if len(roleIds) > 0 {
		err = apiClient.AddRolesToCompositeRole(
			clientId(d), realm(d), resourceDataToRoleRepresentation(d), roleIds)
		if err != nil {
			log.Printf("[WARN] Error when adding composite roles: %s", err.Error())
			return err
		}

		compositeRoleIds, err := apiClient.GetCompositeRoles(clientId(d), realm(d), resourceDataToRoleRepresentation(d))
		if err != nil {
			log.Printf("[WARN] Error when fetching composite roles: %s", err.Error())
			return err
		}

		setCompositeRoleIds(compositeRoleIds, d)
	}

	d.Partial(false)
	return nil
}

func resourceClientRoleUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*keycloak.KeycloakClient)
	log.Printf("[WARN] Updating keycloak client role")
	d.Partial(true)
	err := apiClient.UpdateClientRole(clientId(d), realm(d), resourceDataToRoleRepresentation(d))
	if err != nil {
		return err
	}

	currentRoles, err := apiClient.GetCompositeRoles(clientId(d), realm(d), resourceDataToRoleRepresentation(d))
	if err != nil {
		return err
	}

	desiredRoles := getCompositeRoleIds(d)

	var rolesToAdd []string
	for _, desiredRole := range desiredRoles {
		if !contains(currentRoles, desiredRole) {
			rolesToAdd = append(rolesToAdd, desiredRole)
		}
	}

	var rolesToRemove []string
	for _, currentRole := range currentRoles {
		if !contains(desiredRoles, currentRole) {
			rolesToRemove = append(rolesToRemove, currentRole)
		}
	}

	if len(rolesToAdd) > 0 {
		err = apiClient.AddRolesToCompositeRole(clientId(d), realm(d), resourceDataToRoleRepresentation(d), rolesToAdd)
		if err != nil {
			return err
		}
	}

	if len(rolesToRemove) > 0 {
		err = apiClient.RemoveRolesFromCompositeRole(clientId(d), realm(d), resourceDataToRoleRepresentation(d), rolesToRemove)
		if err != nil {
			return err
		}
	}

	currentRoles, err = apiClient.GetCompositeRoles(clientId(d), realm(d), resourceDataToRoleRepresentation(d))
	if err != nil {
		return err
	}

	setCompositeRoleIds(currentRoles, d)
	d.Partial(false)
	return nil
}

func resourceClientRoleDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*keycloak.KeycloakClient)
	return apiClient.DeleteClientRole(clientId(d), realm(d), resourceDataToRoleRepresentation(d))
}

func getCompositeRoleIds(d *schema.ResourceData) []string {
	return getStringSlice(d, "composite_role_ids")
}

// This method avoids setting the role_ids if they contain the same elements
// This is to avoid reordering based on sorting / retrieval method after updating / changes
func setCompositeRoleIds(currentRoleIds []string, d *schema.ResourceData) {
	storedRoleIds := getStringSlice(d, "composite_role_ids")
	if !containsSameElements(currentRoleIds, storedRoleIds) {
		d.Set("composite_role_ids", currentRoleIds)
	}
}
