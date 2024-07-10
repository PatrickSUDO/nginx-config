package nginx

import (
	_ "embed"
	"regexp"
	"strings"
	"testing"

	"github.com/PatrickSUDO/nginx-config/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed fixtures/nginx.conf
var expectedNginxConf string

//go:embed fixtures/input.yaml
var inputYAML string

func TestGenerateConfig(t *testing.T) {
	// Load the configuration from the embedded YAML
	cfg, err := config.LoadConfigFromBytes([]byte(inputYAML))
	require.NoError(t, err)

	// Generate the NGINX configuration
	nginxConf, err := GenerateConfig(cfg)
	require.NoError(t, err)

	// Normalize configurations for comparison
	expectedNormalized := normalizeConfig(expectedNginxConf)
	actualNormalized := normalizeConfig(nginxConf)

	// Compare normalized configurations
	assert.Equal(t, expectedNormalized, actualNormalized, "Generated NGINX config doesn't match expected config")

	if expectedNormalized != actualNormalized {
		t.Logf("Expected:\n%s", expectedNormalized)
		t.Logf("Actual:\n%s", actualNormalized)
	}
}

func normalizeConfig(config string) string {
	// Remove all whitespace
	re := regexp.MustCompile(`\s+`)
	config = re.ReplaceAllString(config, "")

	// Remove all newlines
	config = strings.ReplaceAll(config, "\n", "")

	return config
}
