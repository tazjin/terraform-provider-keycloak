package provider

import (
	"keycloak"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema:        keycloakProviderSchema(),
		ConfigureFunc: schema.ConfigureFunc(keycloakProviderSetup),
		ResourcesMap: map[string]*schema.Resource{
			"keycloak_client":             resourceClient(),
			"keycloak_user_role_mapping":  resourceUserRoleMapping(),
			"keycloak_realm":              resourceRealm(),
			"keycloak_user":               resourceUser(),
			"keycloak_group":              resourceGroup(),
			"keycloak_user_group_mapping": resourceUserGroupMapping(),
		},
	}
}

// This method provides the provider schema, that is the configuration parameters accepted in the
// provider{} block in Terraform configuration.
func keycloakProviderSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"client_id": {
			Required:    true,
			Type:        schema.TypeString,
			DefaultFunc: schema.EnvDefaultFunc("KEYCLOAK_CLIENT_ID", nil),
		},
		"client_secret": {
			Required:    true,
			Type:        schema.TypeString,
			DefaultFunc: schema.EnvDefaultFunc("KEYCLOAK_CLIENT_SECRET", nil),
		},
		"api_base": {
			Required:    true,
			Type:        schema.TypeString,
			DefaultFunc: schema.EnvDefaultFunc("KEYCLOAK_API_BASE", nil),
		},
		"realm": {
			Optional:    true,
			Type:        schema.TypeString,
			DefaultFunc: schema.EnvDefaultFunc("KEYCLOAK_REALM", "master"),
		},
	}
}

// This method attempts to log in to Keycloak with the provided client credentials
// and returns a configured Keycloak client.
func keycloakProviderSetup(data *schema.ResourceData) (interface{}, error) {
	return keycloak.Login(
		data.Get("client_id").(string),
		data.Get("client_secret").(string),
		data.Get("api_base").(string),
		data.Get("realm").(string),
	)
}
