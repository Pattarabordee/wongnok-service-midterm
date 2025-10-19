package config_test

import (
	"testing"
	"wongnok/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestKeycloakRealmURL(t *testing.T) {
	t.Run("ShouldReturnRealmURL", func(t *testing.T) {
		kc := config.Keycloak{
			Realm: "demo",
			URL:   "http://example.com",
		}

		assert.Equal(t, "http://example.com/realms/demo", kc.RealmURL())
	})
}

func TestKeycloakLogoutURL(t *testing.T) {
	t.Run("ShouldReturnLogoutURL", func(t *testing.T) {
		kc := config.Keycloak{
			Realm: "demo",
			URL:   "http://example.com",
		}

		assert.Equal(t, "http://example.com/realms/demo/protocol/openid-connect/logout", kc.LogoutURL())
	})
}
