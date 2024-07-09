package nginx

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/PatrickSUDO/nginx-config/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed fixtures/input.yaml
var inputYAML []byte

//go:embed fixtures/nginx.conf
var expectedNginxConf []byte

func TestGenerateConfig(t *testing.T) {
	// Load config from embedded YAML
	cfg, err := config.LoadConfigFromBytes(inputYAML)
	require.NoError(t, err, "Failed to load config from input.yaml")

	// Generate NGINX config
	generatedNginxConf, err := GenerateConfig(cfg)
	require.NoError(t, err, "Failed to generate nginx.conf")

	// Normalize both configs
	normalizedExpected := normalizeConfig(string(expectedNginxConf))
	normalizedGenerated := normalizeConfig(generatedNginxConf)

	// Compare normalized configs
	if !assert.Equal(t, normalizedExpected, normalizedGenerated, "Generated nginx.conf does not match expected output") {
		// If configs don't match, print both for easier debugging
		t.Logf("Expected config:\n%s", normalizedExpected)
		t.Logf("Generated config:\n%s", normalizedGenerated)
	}

	// Additional specific checks
	assert.Contains(t, normalizedGenerated, "upstream acceptance", "Generated config should contain upstream for acceptance")
	assert.Contains(t, normalizedGenerated, "upstream production", "Generated config should contain upstream for production")
	assert.Contains(t, normalizedGenerated, "location / {", "Generated config should contain default location block")
	assert.Contains(t, normalizedGenerated, "location /public {", "Generated config should contain location block for /public")
	assert.Contains(t, normalizedGenerated, "location /secret {", "Generated config should contain location block for /secret")
	assert.Contains(t, normalizedGenerated, "allow 82.94.188.0/25;", "Generated config should contain correct IP allow rule")
	assert.Contains(t, normalizedGenerated, "allow 2001:888:2177::/48;", "Generated config should contain correct IPv6 allow rule")
}

// normalizeConfig removes whitespace and newline differences for comparison
func normalizeConfig(config string) string {
	lines := strings.Split(config, "\n")
	var normalizedLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			normalizedLines = append(normalizedLines, trimmed)
		}
	}
	return strings.Join(normalizedLines, "\n")
}
