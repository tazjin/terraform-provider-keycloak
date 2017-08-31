package provider

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tazjin/terraform-provider-keycloak/keycloak"
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
				Type:             schema.TypeMap,
				Optional:         true,
				DiffSuppressFunc: ignoreSmtpPasswordChange,
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
			"account_theme": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"admin_theme": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email_theme": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"login_theme": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

// Keycloak returns some asterisks instead of the plaintext password when querying for the realm configuration.
// Due to this Terraform will assume a change has happened and attempt to reset the password.
// This function will ignore the planned change in such a case, but it will also currently make it impossible to
// change the password (because it's not persisted in the state).
func ignoreSmtpPasswordChange(k, old, new string, d *schema.ResourceData) bool {
	if k == "smtp_server.password" && old == "**********" {
		// It would be nice to print a warning here, but it's unclear how/if providers can output things.
		return true
	}

	return false
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
	r := resourceDataToRealm(d)

	created, err := c.CreateRealm(r)
	if err != nil {
		return err
	}

	d.SetId(created.Id)

	return resourceRealmRead(d, m)
}

func resourceRealmUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*keycloak.KeycloakClient)
	r := resourceDataToRealm(d)
	return c.UpdateRealm(r)
}

func resourceRealmDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*keycloak.KeycloakClient)
	return c.DeleteRealm(d.Id())
}

// Type/struct conversion boilerplate (thanks, Go)

func resourceDataToRealm(d *schema.ResourceData) *keycloak.Realm {
	r := keycloak.Realm{
		Realm:   d.Get("realm").(string),
		Enabled: d.Get("enabled").(bool),

		SslRequired:      d.Get("ssl_required").(string),
		DisplayName:      d.Get("display_name").(string),
		SupportedLocales: getStringSlice(d, "supported_locales"),
		DefaultRoles:     getStringSlice(d, "default_roles"),

		AccountTheme: d.Get("account_theme").(string),
		AdminTheme:   d.Get("admin_theme").(string),
		EmailTheme:   d.Get("email_theme").(string),
		LoginTheme:   d.Get("login_theme").(string),

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

	if smtpMap, present := d.GetOk("smtp_server"); present {
		smtp := keycloak.SmtpServer(smtpMap.(map[string]interface{}))
		r.SmtpServer = &smtp
	}

	return &r
}

func realmToResourceData(r *keycloak.Realm, d *schema.ResourceData) {
	d.SetId(r.Id)
	d.Set("realm", r.Realm)
	d.Set("enabled", r.Enabled)

	d.Set("ssl_required", r.SslRequired)
	d.Set("display_name", r.DisplayName)
	d.Set("supported_locales", r.SupportedLocales)
	d.Set("default_roles", r.DefaultRoles)

	d.Set("account_theme", r.AccountTheme)
	d.Set("admin_theme", r.AdminTheme)
	d.Set("email_theme", r.EmailTheme)
	d.Set("login_theme", r.LoginTheme)

	if r.SmtpServer != nil {
		d.Set("smtp_server", *r.SmtpServer)
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
