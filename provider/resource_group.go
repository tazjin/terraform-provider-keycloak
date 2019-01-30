// This file provides a Terraform resource for Keycloak clients
// The client resource is documented at http://www.keycloak.org/docs-api/3.1/rest-api/index.html#_clientrepresentation

package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tazjin/terraform-provider-keycloak/keycloak"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		// API methods
		Read:   schema.ReadFunc(resourceGroupRead),
		Create: schema.CreateFunc(resourceGroupCreate),
		Update: schema.UpdateFunc(resourceGroupUpdate),
		Delete: schema.DeleteFunc(resourceGroupDelete),

		// Keycloak clients are importable by ID, but the realm must also be provided by the user.
		Importer: &schema.ResourceImporter{
			State: importGroupHelper,
		},

		Schema: map[string]*schema.Schema{
			"realm": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"attributes": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"realmroles": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"clientroles": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"subgroups": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func importGroupHelper(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	realm, id, err := splitRealmId(d.Id())
	if err != nil {
		return nil, err
	}

	d.SetId(id)
	d.Set("realm", realm)

	resourceGroupRead(d, m)

	return []*schema.ResourceData{d}, nil
}

func resourceGroupRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*keycloak.KeycloakClient)

	group, err := c.GetGroup(d.Id(), realm(d))
	if err != nil {
		return err
	}

	groupToResourceData(group, d)

	return nil
}

func resourceGroupCreate(d *schema.ResourceData, m interface{}) error {

	apiGroup := m.(*keycloak.KeycloakClient)
	group := resourceDataToGroup(d)
	created, err := apiGroup.AddGroup(&group, realm(d))

	if err != nil {
		return err
	}

	d.SetId(created.Id)

	return resourceGroupRead(d, m)
}

func resourceGroupUpdate(d *schema.ResourceData, m interface{}) error {
	user := resourceDataToGroup(d)
	apiGroup := m.(*keycloak.KeycloakClient)
	return apiGroup.UpdateGroup(&user, realm(d))
}

func resourceGroupDelete(d *schema.ResourceData, m interface{}) error {
	apiGroup := m.(*keycloak.KeycloakClient)
	return apiGroup.DeleteGroup(d.Id(), realm(d))
}

//TODO: Support subgroups, might have to make API calls to read based on name
func resourceDataToGroup(d *schema.ResourceData) keycloak.Group {
	u := keycloak.Group{
		Name:        d.Get("name").(string),
		Attributes:  getOptionalStringMap(d, "attributes"),
		RealmRoles:  getOptionalStringList(d, "realmroles"),
		ClientRoles: getOptionalStringMap(d, "clientroles"),
	}

	if !d.IsNewResource() {
		u.Id = d.Id()
	}

	return u
}

//TODO: Support subgroups
func groupToResourceData(g *keycloak.Group, d *schema.ResourceData) {
	d.SetId(g.Id)
	d.Set("id", g.Id)
	d.Set("name", g.Name)
	d.Set("attributes", g.Attributes)
	d.Set("realmroles", g.RealmRoles)
	d.Set("clientroles", g.ClientRoles)
	//d.Set("subgroups", u.SubGroups)
}
