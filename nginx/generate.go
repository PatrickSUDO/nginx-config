package nginx

import (
	"bytes"
	"text/template"

	"github.com/PatrickSUDO/nginx-config/config"
)

func GenerateConfig(cfg *config.Config) (string, error) {
	tmpl, err := template.ParseFiles("nginx/template.conf")
	if err != nil {
		return "", err
	}
	var result bytes.Buffer
	err = tmpl.Execute(&result, cfg)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}
