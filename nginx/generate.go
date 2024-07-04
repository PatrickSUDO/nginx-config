package nginx

import (
	"bytes"
	"text/template"

	"github.com/PatrickSUDO/nginx-config/config"
)

func GenerateConfig(cfg *config.Config) (string, error) {
	tmpl, err := template.New("nginx.conf").Parse(`
{{- range $catchallName, $catchall := .Catchall }}
server {
    listen [::]:{{ $catchall.Port }} default_server ipv6only=on;
    listen 0.0.0.0:{{ $catchall.Port }} default_server;
    server_name _;
    root /var/www/;
    location / {
        return 503;
    }
}
{{- end }}

{{- range $appName, $app := .App }}
upstream {{ $appName }} {
    server 127.0.0.1:{{ $app.RuntimePort }};
}

server {
    listen [::]:{{ (index $.Catchall $app.Catchall).Port }};
    listen 0.0.0.0:{{ (index $.Catchall $app.Catchall).Port }};
    {{- range $app.FQDN }}
    server_name {{ . }};
    {{- end }}

    {{- range $path, $restriction := $app.PathBasedAccessRestriction }}
    location {{ $path }} {
        proxy_pass http://{{ $appName }};
        {{- range (index $.IPFilter $restriction.IPFilter) }}
        allow {{ . }};
        {{- end }}
        deny all;
    }
    {{- end }}

    {{- if eq (len $app.PathBasedAccessRestriction) 0 }}
    location / {
        proxy_pass http://{{ $appName }};
    }
    {{- end }}
}
{{- end }}
`)

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
