package provider

import (
	"keycloak"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUserRoleMapping() *schema.Resource {
	return &schema.Resource{
		// API methods
		Read:   schema.ReadFunc(resourceUserRoleMappingRead),
		Create: schema.CreateFunc(resourceUserRoleMappingCreate),
		Delete: schema.DeleteFunc(resourceUserRoleMappingDelete),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scope_param_required": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default: "",
				ForceNew: true,
			},
			"realm": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "master",
				ForceNew: true,
			},
		},
	}
}

func resourceUserRoleMappingRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*keycloak.KeycloakClient)
	userId := d.Get("user_id").(string)
	clientId := d.Get("client_id").(string)

	roles, err := c.GetCompositeRolesForUser(userId, realm(d), clientId)
	if err != nil {
		return err
	}

	role, err := c.FindRoleForUser(roles, d.Id())
	if err != nil {
		return err
	}

	userRoleMappingToResourceData(userId, role, d)
	return nil
}

func resourceUserRoleMappingCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*keycloak.KeycloakClient)

	role, err := c.AddRoleToUser(
		d.Get("user_id").(string),
		d.Get("name").(string),
		realm(d),
		d.Get("client_id").(string),
	)

	if err != nil {
		return err
	}

	d.SetId(role.Id)

	return resourceUserRoleMappingRead(d, m)
}

func resourceUserRoleMappingDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*keycloak.KeycloakClient)
	role := resourceDataToUserRoleMapping(d)
	return c.RemoveRoleFromUser(d.Get("user_id").(string), &role, realm(d), d.Get("client_id").(string))
}

func userRoleMappingToResourceData(userId string, r *keycloak.Role, d *schema.ResourceData) {
	d.SetId(r.Id)
	d.Set("user_id", userId)
	d.Set("name", r.Name)
	d.Set("scope_param_required", r.ScopeParamRequired)
}

func resourceDataToUserRoleMapping(d *schema.ResourceData) keycloak.Role {
	return keycloak.Role{
		Id:                 d.Id(),
		Name:               d.Get("name").(string),
		ScopeParamRequired: d.Get("scope_param_required").(bool),
	}
}
