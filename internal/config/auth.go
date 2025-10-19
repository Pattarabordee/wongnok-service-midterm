package config

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type Keycloak struct {
	ClientID     string `env:"KEYCLOCK_CLIENT_ID"`
	ClientSecret string `env:"KEYCLOCK_CLIENT_SECRET"`
	RedirectURL  string `env:"KEYCLOCK_REDIRECT_URL"`
	Realm        string `env:"KEYCLOAK_REALM"`
	URL          string `env:"KEYCLOAK_URL"`
}

func (kc Keycloak) RealmURL() string {
	return fmt.Sprintf("%s/realms/%s", kc.URL, kc.Realm)
}

func (kc Keycloak) LogoutURL() string {
	return fmt.Sprintf("%s/protocol/openid-connect/logout", kc.RealmURL())
}

type IOAuth2Config interface {
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
}

type IOIDCTokenVerifier interface {
	Verify(ctx context.Context, rawIDToken string) (*oidc.IDToken, error)
}

type IOIDCIDToken interface {
	Claims(v any) error
}
