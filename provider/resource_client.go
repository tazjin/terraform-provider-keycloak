// This file provides a Terraform resource for Keycloak clients
// The client resource is documented at http://www.keycloak.org/docs-api/3.1/rest-api/index.html#_clientrepresentation

package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tazjin/terraform-provider-keycloak/keycloak"
)

func resourceClient() *schema.Resource {
	return &schema.Resource{
		// API methods
		Read:   schema.ReadFunc(resourceClientRead),
		Create: schema.CreateFunc(resourceClientCreate),
		Update: schema.UpdateFunc(resourceClientUpdate),
		Delete: schema.DeleteFunc(resourceClientDelete),

		// Keycloak clients are importable by ID, so no import logic is required!
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"realm": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "master",
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"client_authenticator_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "client-secret",
			},
			"redirect_uris": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "openid-connect",
			},
			"public_client": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"bearer_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"service_accounts_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"web_origins": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			// Computed fields (i.e. things looked up in Keycloak after client creation)
			"client_secret": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_account_user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceClientRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*keycloak.KeycloakClient)

	client, err := c.GetClient(d.Id(), realm(d))
	if err != nil {
		return err
	}

	clientToResourceData(client, d)

	// Look up client secret in addition
	secret, err := c.GetClientSecret(d.Id(), realm(d))
	if err != nil {
		return err
	}
	d.Set("client_secret", secret.Value)

	// Look up service account user ID (if enabled)
	if client.ServiceAccountsEnabled {
		user, err := c.GetClientServiceAccountUser(d.Id(), realm(d))
		if err != nil {
			return err
		}

		d.Set("service_account_user_id", user.Id)
	}

	return nil
}

func resourceClientCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*keycloak.KeycloakClient)
	client := resourceDataToClient(d)
	created, err := apiClient.CreateClient(&client, realm(d))

	if err != nil {
		return err
	}

	d.SetId(created.Id)

	return resourceClientRead(d, m)
}

func resourceClientUpdate(d *schema.ResourceData, m interface{}) error {
	client := resourceDataToClient(d)
	apiClient := m.(*keycloak.KeycloakClient)
	return apiClient.UpdateClient(&client, realm(d))
}

func resourceClientDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*keycloak.KeycloakClient)
	return apiClient.DeleteClient(d.Id(), realm(d))
}

func resourceDataToClient(d *schema.ResourceData) keycloak.Client {
	redirectUris := []string{}
	webOrigins := []string{}

	for _, uri := range d.Get("redirect_uris").([]interface{}) {
		redirectUris = append(redirectUris, uri.(string))
	}

	rawOrigins, present := d.GetOk("web_origins")
	if present {
		for _, origin := range rawOrigins.([]interface{}) {
			webOrigins = append(webOrigins, origin.(string))
		}
	}

	c := keycloak.Client{
		ClientId:                d.Get("client_id").(string),
		Enabled:                 d.Get("enabled").(bool),
		ClientAuthenticatorType: d.Get("client_authenticator_type").(string),
		RedirectUris:            redirectUris,
		Protocol:                d.Get("protocol").(string),
		PublicClient:            d.Get("public_client").(bool),
		BearerOnly:              d.Get("bearer_only").(bool),
		ServiceAccountsEnabled:  d.Get("service_accounts_enabled").(bool),
		WebOrigins:              webOrigins,
	}

	if !d.IsNewResource() {
		c.Id = d.Id()
	}

	return c
}

func clientToResourceData(c *keycloak.Client, d *schema.ResourceData) {
	d.Set("client_id", c.ClientId)
	d.Set("enabled", c.Enabled)
	d.Set("client_authenticator_type", c.ClientAuthenticatorType)
	d.Set("redirect_uris", c.RedirectUris)
	d.Set("protocol", c.Protocol)
	d.Set("public_client", c.PublicClient)
	d.Set("bearer_only", c.BearerOnly)
	d.Set("service_accounts_enabled", c.ServiceAccountsEnabled)
	d.Set("web_origins", c.WebOrigins)
}
