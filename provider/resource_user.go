// This file provides a Terraform resource for Keycloak clients
// The client resource is documented at http://www.keycloak.org/docs-api/3.1/rest-api/index.html#_clientrepresentation

package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tazjin/terraform-provider-keycloak/keycloak"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		// API methods
		Read:   schema.ReadFunc(resourceUserRead),
		Create: schema.CreateFunc(resourceUserCreate),
		Update: schema.UpdateFunc(resourceUserUpdate),
		Delete: schema.DeleteFunc(resourceUserDelete),

		// Keycloak clients are importable by ID, but the realm must also be provided by the user.
		Importer: &schema.ResourceImporter{
			State: importUserHelper,
		},

		Schema: map[string]*schema.Schema{
			"realm": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"firstname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lastname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Valid actions are: "CONFIGURE_TOTP", "UPDATE_PASSWORD", "UPDATE_PROFILE", "VERIFY_EMAIL"
			"initial_required_actions": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				// Suppress planned changes to required actions if the user has an ID,
				// meaning the user account has already been created
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Id() != ""
				},
			},
		},
	}
}

func importUserHelper(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	realm, id, err := splitRealmId(d.Id())
	if err != nil {
		return nil, err
	}

	d.SetId(id)
	d.Set("realm", realm)

	resourceUserRead(d, m)

	return []*schema.ResourceData{d}, nil
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*keycloak.KeycloakClient)

	user, err := c.GetUser(d.Id(), realm(d))
	if err != nil {
		return err
	}

	userToResourceData(user, d)

	return nil
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {

	apiUser := m.(*keycloak.KeycloakClient)
	user := resourceDataToUser(d)
	created, err := apiUser.AddUser(&user, realm(d))

	if err != nil {
		return err
	}

	d.SetId(created.Id)

	return resourceUserRead(d, m)
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	user := resourceDataToUser(d)
	apiUser := m.(*keycloak.KeycloakClient)
	return apiUser.UpdateUser(&user, realm(d))
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	apiUser := m.(*keycloak.KeycloakClient)
	return apiUser.DeleteUser(d.Id(), realm(d))
}

func resourceDataToUser(d *schema.ResourceData) keycloak.User {
	u := keycloak.User{
		Username:  d.Get("username").(string),
		Enabled:   d.Get("enabled").(bool),
		FirstName: d.Get("firstname").(string),
		LastName:  d.Get("lastname").(string),
		Email:     d.Get("email").(string),
	}

	if !d.IsNewResource() {
		u.Id = d.Id()
	} else {
		u.Id = u.Username
		u.RequiredActions = getOptionalStringSet(d, "initial_required_actions")
	}

	return u
}

func userToResourceData(u *keycloak.User, d *schema.ResourceData) {
	d.SetId(u.Id)
	d.Set("username", u.Username)
	d.Set("enabled", u.Enabled)
	d.Set("firstname", u.FirstName)
	d.Set("lastname", u.LastName)
	d.Set("email", u.Email)
}

// Custom helper function needed to handle optional schema.TypeSet fields because
// getOptionalStringList() from the terraform helper only supports schema.TypeList
func getOptionalStringSet(d *schema.ResourceData, key string) []string {
	stringList := []string{}
	rawSet, present := d.GetOk(key)
	if present {
		for _, stringVal := range (rawSet.(*schema.Set)).List() {
			stringList = append(stringList, stringVal.(string))
		}
	}
	return stringList
}
