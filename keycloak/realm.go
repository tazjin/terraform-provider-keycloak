package keycloak

import "fmt"

// Representation of top-level realm keys. According to the Keycloak documentation other keys than top-level keys will
// be ignored on realm updates, which is why they are not included here.
// http://www.keycloak.org/docs-api/3.1/rest-api/index.html#_realmrepresentation
type Realm struct {
	// General realm settings
	Id      string `json:"id"`
	Realm   string `json:"realm"`
	Enabled bool   `json:"enabled"`

	// Optional realm settings
	SslRequired      string      `json:"sslRequired,omitempty"` // valid values are ALL, NONE or EXTERNAL
	DisplayName      string      `json:"displayName,omitempty"`
	SupportedLocales []string    `json:"supportedLocales,omitempty"`
	DefaultRoles     []string    `json:"defaultRoles,omitempty"`

	// The available keys of the SMTP server map are not documented in Keycloak's API docs, but they can be found in the
	// source code at:
	// https://github.com/keycloak/keycloak/blob/master/services/src/main/java/org/keycloak/email/DefaultEmailSenderProvider.java
	SmtpServer       *map[string]interface{} `json:"smtpServer,omitempty"`

	InternationalizationEnabled *bool `json:"internationalizationEnabled,omitempty"`
	RegistrationAllowed         *bool `json:"registrationAllowed,omitempty"`
	RegistrationEmailAsUsername *bool `json:"registrationEmailAsUsername,omitempty"`
	RememberMe                  *bool `json:"rememberMe,omitempty"`
	VerifyEmail                 *bool `json:"verifyEmail,omitempty"`
	ResetPasswordAllowed        *bool `json:"resetPasswordAllowed,omitempty"`
	EditUsernameAllowed         *bool `json:"editUsernameAllowed,omitempty"`
	BruteForceProtected         *bool `json:"bruteForceProtected,omitempty"`

	// Token & session settings
	AccessTokenLifespan                *int `json:"accessTokenLifespan,omitempty"`
	AccessTokenLifespanForImplicitFlow *int `json:"accessTokenLifespanForImplicitFlow,omitempty"`
	SsoSessionIdleTimeout              *int `json:"ssoSessionIdleTimeout,omitempty"`
	SsoSessionMaxLifespan              *int `json:"ssoSessionMaxLifespan,omitempty"`
	OfflineSessionIdleTimeout          *int `json:"offlineSessionIdleTimeout,omitempty"`
	AccessCodeLifespan                 *int `json:"accessCodeLifespan,omitempty"`
	AccessCodeLifespanUserAction       *int `json:"accessCodeLifespanUserAction,omitempty"`
	AccessCodeLifespanLogin            *int `json:"accessCodeLifespanLogin,omitempty"`
	MaxFailureWaitSeconds              *int `json:"maxFailureWaitSeconds,omitempty"`
	MinimumQuickLoginWaitSeconds       *int `json:"minimumQuickLoginWaitSeconds,omitempty"`
	WaitIncrementSeconds               *int `json:"waitIncrementSeconds,omitempty"`
	QuickLoginCheckMilliSeconds        *int `json:"quickLoginCheckMilliSeconds,omitempty"`
	MaxDeltaTimeSeconds                *int `json:"maxDeltaTimeSeconds,omitempty"`
	FailureFactor                      *int `json:"failureFactor,omitempty"`
}

const (
	realmsUri = "%s/auth/admin/realms"
	realmUri  = "%s/auth/admin/realms/%s"
)

func (c *KeycloakClient) GetRealm(id string) (*Realm, error) {
	url := fmt.Sprintf(realmUri, c.url, id)

	var r Realm
	err := c.get(url, &r)

	return &r, err
}

// This "imports" (i.e. creates) a realm from a realm representation.
func (c *KeycloakClient) CreateRealm(r *Realm) (*Realm, error) {
	url := fmt.Sprintf(realmsUri, c.url)

	realmLocation, err := c.post(url, *r)
	if err != nil {
		return nil, err
	}

	var createdRealm Realm
	err = c.get(realmLocation, &createdRealm)

	return &createdRealm, err
}

func (c *KeycloakClient) UpdateRealm(r *Realm) error {
	url := fmt.Sprintf(realmUri, c.url, r.Id)
	return c.put(url, *r)
}

func (c *KeycloakClient) DeleteRealm(id string) error {
	url := fmt.Sprintf(realmUri, c.url, id)
	return c.delete(url, nil)
}
