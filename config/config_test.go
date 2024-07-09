package config

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed fixtures/input.yaml
var inputYAML []byte

func TestLoadConfig(t *testing.T) {
	// Load config from embedded YAML
	cfg, err := LoadConfigFromBytes(inputYAML)
	assert.NoError(t, err, "Failed to load config from input.yaml")

	// Check expected values
	expectedIPFilter := map[string][]string{
		"myfilter": {"82.94.188.0/25", "2001:888:2177::/48"},
		"allowall": {"0.0.0.0/0", "::/0"},
	}
	assert.Equal(t, expectedIPFilter, cfg.IPFilter, "IPFilter mismatch")

	assert.Equal(t, 7000, cfg.Catchall["default"].Port, "Catchall default port mismatch")

	// Acceptance app checks
	assert.Equal(t, 8000, cfg.App["acceptance"].RuntimePort, "App acceptance runtime port mismatch")

	expectedFQDN := []string{"myapp-accp.mendixcloud.com", "accp.myapp.com"}
	assert.Equal(t, expectedFQDN, cfg.App["acceptance"].FQDN, "App acceptance FQDN mismatch")

	// Check path-based access restrictions for acceptance
	acceptanceRestrictions := cfg.App["acceptance"].PathBasedAccessRestriction
	assert.Equal(t, "myfilter", acceptanceRestrictions["/"].IPFilter, "Acceptance '/' path restriction mismatch")
	assert.Equal(t, "allowall", acceptanceRestrictions["/public"].IPFilter, "Acceptance '/public' path restriction mismatch")

	// Production app checks
	assert.Equal(t, 8001, cfg.App["production"].RuntimePort, "App production runtime port mismatch")

	expectedProductionFQDN := []string{"myapp.mendixcloud.com", "myapp.com"}
	assert.Equal(t, expectedProductionFQDN, cfg.App["production"].FQDN, "App production FQDN mismatch")

	// Check path-based access restrictions for production
	productionRestrictions := cfg.App["production"].PathBasedAccessRestriction
	assert.Equal(t, "myfilter", productionRestrictions["/secret"].IPFilter, "Production '/secret' path restriction mismatch")
}
