package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tazjin/terraform-provider-keycloak/keycloak"
	"sort"
	"io/ioutil"
)

func resourceRealm() *schema.Resource {
	return &schema.Resource{
		// API methods
		Read:   schema.ReadFunc(resourceRealmRead),
		Create: schema.CreateFunc(resourceRealmCreate),
		Update: schema.UpdateFunc(resourceRealmUpdate),
		Delete: schema.DeleteFunc(resourceRealmDelete),

		// Realms are importable by ID
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"realm": {
				Description: "Realm name and ID",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"ssl_required": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "external",
				ValidateFunc: validateSslRequired,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"supported_locales": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Due to Terraform being unable to deal with list default values this field *must* be set.
			// The default values from Keycloak are: ["offline_access", "uma_authorization"]
			"default_roles": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"smtp_server": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     smtpMapSchema(),
				Set:      schema.SchemaSetFunc(smtpSettingSetHash),
			},
			"internationalization_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"registration_allowed": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"registration_email_as_username": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"remember_me": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"verify_email": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"reset_password_allowed": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"edit_username_allowed": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"brute_force_protected": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"access_token_lifespan": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
			},
			"access_token_lifespan_for_implicit_flow": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  900,
			},
			"sso_session_idle_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1800,
			},
			"sso_session_max_lifespan": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  36000,
			},
			"offline_session_idle_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2592000,
			},
			"access_code_lifespan": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  60,
			},
			"access_code_lifespan_user_action": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
			},
			"access_code_lifespan_login": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1800,
			},
			"max_failure_wait_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  900,
			},
			"minimum_quick_login_wait_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  60,
			},
			"wait_increment_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  60,
			},
			"quick_login_check_milli_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1000,
			},
			"max_delta_time_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  43200,
			},
			"failure_factor": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
		},
	}
}

func smtpMapSchema() *schema.Resource {
	return &schema.Resource{
		// Every type in the map schema is 'string' because this is sent to the server as a Map<String, String>
		// and parsed there as a Java property object.
		Schema: map[string]*schema.Schema{
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"starttls": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auth": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type: schema.TypeString,
				Sensitive: true,
				Optional: true,
				DiffSuppressFunc: func(_, old, _ string, d *schema.ResourceData) bool {
					if old == "**********" {
						return true
					}
					return false
				},
			},
			"from": {
				Type:     schema.TypeString,
				Required: true,
			},
			"fromDisplayName": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"replyTo": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"replyToDisplayName": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"envelopeFrom": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func validateSslRequired(v interface{}, _ string) (w []string, err []error) {
	switch v.(string) {
	case
		"all",
		"external",
		"none":
		return
	}
	err = []error{
		fmt.Errorf("Invalid value for ssl_required. Valid are ALL, EXTERNAL or NONE"),
	}
	return
}

func resourceRealmRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*keycloak.KeycloakClient)

	r, err := c.GetRealm(d.Id())
	if err != nil {
		return err
	}

	realmToResourceData(r, d)
	return nil
}

func resourceRealmCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*keycloak.KeycloakClient)
	r, err := resourceDataToRealm(d)
	if err != nil {
		return err
	}

	created, err := c.CreateRealm(r)
	if err != nil {
		return err
	}

	d.SetId(created.Id)

	return resourceRealmRead(d, m)
}

func resourceRealmUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*keycloak.KeycloakClient)
	r, err := resourceDataToRealm(d)
	if err != nil {
		return err
	}
	return c.UpdateRealm(r)
}

func resourceRealmDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*keycloak.KeycloakClient)
	return c.DeleteRealm(d.Id())
}

// Type/struct conversion boilerplate (thanks, Go)

func resourceDataToRealm(d *schema.ResourceData) (*keycloak.Realm, error) {
	r := keycloak.Realm{
		Realm:   d.Get("realm").(string),
		Enabled: d.Get("enabled").(bool),

		SslRequired:      d.Get("ssl_required").(string),
		DisplayName:      d.Get("display_name").(string),
		SupportedLocales: getStringSlice(d, "supported_locales"),
		DefaultRoles:     getStringSlice(d, "default_roles"),

		InternationalizationEnabled: getOptionalBool(d, "internationalization_enabled"),
		RegistrationAllowed:         getOptionalBool(d, "registration_allowed"),
		RegistrationEmailAsUsername: getOptionalBool(d, "registration_email_as_username"),
		RememberMe:                  getOptionalBool(d, "remember_me"),
		VerifyEmail:                 getOptionalBool(d, "verify_email"),
		ResetPasswordAllowed:        getOptionalBool(d, "reset_password_allowed"),
		EditUsernameAllowed:         getOptionalBool(d, "edit_username_allowed"),
		BruteForceProtected:         getOptionalBool(d, "brute_force_protected"),

		AccessTokenLifespan:                getOptionalInt(d, "access_token_lifespan"),
		AccessTokenLifespanForImplicitFlow: getOptionalInt(d, "access_token_lifespan_for_implicit_flow"),
		SsoSessionIdleTimeout:              getOptionalInt(d, "sso_session_idle_timeout"),
		SsoSessionMaxLifespan:              getOptionalInt(d, "sso_session_max_lifespan"),
		OfflineSessionIdleTimeout:          getOptionalInt(d, "offline_session_idle_timeout"),
		AccessCodeLifespan:                 getOptionalInt(d, "access_code_lifespan"),
		AccessCodeLifespanUserAction:       getOptionalInt(d, "access_code_lifespan_user_action"),
		AccessCodeLifespanLogin:            getOptionalInt(d, "access_code_lifespan_login"),
		MaxFailureWaitSeconds:              getOptionalInt(d, "max_failure_wait_seconds"),
		MinimumQuickLoginWaitSeconds:       getOptionalInt(d, "minimum_quick_login_wait_seconds"),
		WaitIncrementSeconds:               getOptionalInt(d, "wait_increment_seconds"),
		QuickLoginCheckMilliSeconds:        getOptionalInt(d, "quick_login_check_milli_seconds"),
		MaxDeltaTimeSeconds:                getOptionalInt(d, "max_delta_time_seconds"),
		FailureFactor:                      getOptionalInt(d, "failure_factor"),
	}

	if !d.IsNewResource() {
		r.Id = d.Id()
	} else {
		r.Id = r.Realm
	}

	if smtpSet, present := d.GetOk("smtp_server"); present {
		settings, err := setToSmtpSettings(smtpSet.(*schema.Set))

		if err != nil {
			return &r, err
		}
		r.SmtpServer = settings
	}

	return &r, nil
}

func realmToResourceData(r *keycloak.Realm, d *schema.ResourceData) {
	d.SetId(r.Id)
	d.Set("realm", r.Realm)
	d.Set("enabled", r.Enabled)

	d.Set("ssl_required", r.SslRequired)
	d.Set("display_name", r.DisplayName)
	d.Set("supported_locales", r.SupportedLocales)
	d.Set("default_roles", r.DefaultRoles)

	if r.SmtpServer != nil {
		d.Set("smtp_server", smtpSettingsToSet(r.SmtpServer))
	}

	setOptionalBool(d, "internationalization_enabled", r.InternationalizationEnabled)
	setOptionalBool(d, "registration_allowed", r.RegistrationAllowed)
	setOptionalBool(d, "registration_email_as_username", r.RegistrationEmailAsUsername)
	setOptionalBool(d, "remember_me", r.RememberMe)
	setOptionalBool(d, "verify_email", r.VerifyEmail)
	setOptionalBool(d, "reset_password_allowed", r.ResetPasswordAllowed)
	setOptionalBool(d, "edit_username_allowed", r.EditUsernameAllowed)
	setOptionalBool(d, "brute_force_protected", r.BruteForceProtected)

	setOptionalInt(d, "access_token_lifespan", r.AccessTokenLifespan)
	setOptionalInt(d, "access_token_lifespan_for_implicit_flow", r.AccessTokenLifespanForImplicitFlow)
	setOptionalInt(d, "sso_session_idle_timeout", r.SsoSessionIdleTimeout)
	setOptionalInt(d, "sso_session_max_lifespan", r.SsoSessionMaxLifespan)
	setOptionalInt(d, "offline_session_idle_timeout", r.OfflineSessionIdleTimeout)
	setOptionalInt(d, "access_code_lifespan", r.AccessCodeLifespan)
	setOptionalInt(d, "access_code_lifespan_user_action", r.AccessCodeLifespanUserAction)
	setOptionalInt(d, "access_code_lifespan_login", r.AccessCodeLifespanLogin)
	setOptionalInt(d, "max_failure_wait_seconds", r.MaxFailureWaitSeconds)
	setOptionalInt(d, "minimum_quick_login_wait_seconds", r.MinimumQuickLoginWaitSeconds)
	setOptionalInt(d, "wait_increment_seconds", r.WaitIncrementSeconds)
	setOptionalInt(d, "quick_login_check_milli_seconds", r.QuickLoginCheckMilliSeconds)
	setOptionalInt(d, "max_delta_time_seconds", r.MaxDeltaTimeSeconds)
	setOptionalInt(d, "failure_factor", r.FailureFactor)
}

func setToSmtpSettings(set *schema.Set) (*map[string]interface{}, error) {
	settings := set.List()

	if len(settings) == 0 {
		return nil, nil
	}

	if len(settings) > 1 {
		return nil, fmt.Errorf("Only one SMTP server can be defined per realm")
	}

	smtpServer := settings[0].(map[string]interface{})

	return &smtpServer, nil
}

func smtpSettingsToSet(settings *map[string]interface{}) *schema.Set {
	return schema.NewSet(schema.HashResource(smtpMapSchema()), []interface{}{*settings})
}

// Perform a consistent hash of an smtp setting set (keys are normally unordered, causing
// differences in hashes).
func smtpSettingSetHash(v interface{}) int {
	ioutil.WriteFile("/tmp/debug", []byte(fmt.Sprint("thing: ", v)), 0666)
	settingSet := v.(*schema.Set)
	// Can't do anything with the error here v0v
	settings, _ := setToSmtpSettings(settingSet)

	keys := make([]string, len(*settings))

	i := 0
	for key, _ := range *settings {
		keys[i] = key
		i++
	}

	sort.Strings(keys)

	var hashInput string
	for _, key := range keys {
		hashInput += fmt.Sprint(key, ":", (*settings)[key])
	}

	return hashcode.String(hashInput)
}
