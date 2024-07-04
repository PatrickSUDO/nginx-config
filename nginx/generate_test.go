package nginx

import (
	"os"
	"testing"

	"github.com/PatrickSUDO/nginx-config/config"
)

func TestGenerateConfig(t *testing.T) {
	inputYaml, err := os.ReadFile("/mnt/data/input.yaml")
	if err != nil {
		t.Fatalf("Failed to read input.yaml: %v", err)
	}

	expectedNginxConf, err := os.ReadFile("/mnt/data/nginx.conf")
	if err != nil {
		t.Fatalf("Failed to read expected nginx.conf: %v", err)
	}

	cfg, err := config.LoadConfigFromString(string(inputYaml))
	if err != nil {
		t.Fatalf("Failed to load config from input.yaml: %v", err)
	}

	generatedNginxConf, err := GenerateConfig(cfg)
	if err != nil {
		t.Fatalf("Failed to generate nginx.conf: %v", err)
	}

	if string(generatedNginxConf) != string(expectedNginxConf) {
		t.Errorf("Generated nginx.conf does not match expected output.\nGenerated:\n%s\nExpected:\n%s", generatedNginxConf, expectedNginxConf)
	}
}
