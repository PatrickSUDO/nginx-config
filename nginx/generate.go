package nginx

import (
	"bytes"
	"fmt"

	"github.com/PatrickSUDO/nginx-config/config"
	crossplane "github.com/nginxinc/nginx-go-crossplane"
)

func GenerateConfig(cfg *config.Config) (string, error) {
	payload := &crossplane.Payload{
		Status: "ok",
		Errors: []crossplane.PayloadError{},
		Config: []crossplane.Config{
			{
				File:   "nginx.conf",
				Parsed: crossplane.Directives{},
			},
		},
	}

	// Add catchall server
	for _, catchall := range cfg.Catchall {
		catchallServer := createCatchallServer(catchall.Port)
		payload.Config[0].Parsed = append(payload.Config[0].Parsed, &catchallServer)
	}

	// Add app servers
	for appName, app := range cfg.App {
		upstreamBlock := createUpstreamBlock(appName, app.RuntimePort)
		serverBlock := createServerBlock(appName, app, cfg)

		payload.Config[0].Parsed = append(payload.Config[0].Parsed, &upstreamBlock, &serverBlock)
	}

	// Convert the config to NGINX format
	var buf bytes.Buffer
	err := crossplane.Build(&buf, crossplane.Config{Parsed: payload.Config[0].Parsed}, &crossplane.BuildOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to build NGINX config: %w", err)
	}

	return buf.String(), nil
}

func createCatchallServer(port int) crossplane.Directive {
	return crossplane.Directive{
		Directive: "server",
		Block: crossplane.Directives{
			{Directive: "listen", Args: []string{fmt.Sprintf("[::]:%d", port), "default_server", "ipv6only=on"}},
			{Directive: "listen", Args: []string{fmt.Sprintf("0.0.0.0:%d", port), "default_server"}},
			{Directive: "server_name", Args: []string{"_"}},
			{Directive: "root", Args: []string{"/var/www/"}},
			{
				Directive: "location",
				Args:      []string{"/"},
				Block: crossplane.Directives{
					{Directive: "return", Args: []string{"503"}},
				},
			},
		},
	}
}

func createUpstreamBlock(appName string, runtimePort int) crossplane.Directive {
	return crossplane.Directive{
		Directive: "upstream",
		Args:      []string{appName},
		Block: crossplane.Directives{
			{Directive: "server", Args: []string{fmt.Sprintf("127.0.0.1:%d", runtimePort)}},
		},
	}
}

func createServerBlock(appName string, app config.AppConfig, cfg *config.Config) crossplane.Directive {
	serverBlock := crossplane.Directive{
		Directive: "server",
		Block:     crossplane.Directives{},
	}

	// Add listen directives
	port := cfg.Catchall[app.Catchall].Port
	serverBlock.Block = append(serverBlock.Block,
		&crossplane.Directive{Directive: "listen", Args: []string{fmt.Sprintf("[::]:%d", port)}},
		&crossplane.Directive{Directive: "listen", Args: []string{fmt.Sprintf("0.0.0.0:%d", port)}},
	)

	// Add server_name directives
	for _, fqdn := range app.FQDN {
		serverBlock.Block = append(serverBlock.Block,
			&crossplane.Directive{Directive: "server_name", Args: []string{fqdn}},
		)
	}

	// Add location blocks
	for path, restriction := range app.PathBasedAccessRestriction {
		locationBlock := createLocationBlock(appName, path, restriction.IPFilter, cfg.IPFilter)
		serverBlock.Block = append(serverBlock.Block, &locationBlock)
	}

	// Add default location block if not present
	if _, exists := app.PathBasedAccessRestriction["/"]; !exists {
		defaultLocationBlock := createDefaultLocationBlock(appName)
		serverBlock.Block = append(serverBlock.Block, &defaultLocationBlock)
	}

	return serverBlock
}

func createLocationBlock(appName, path, ipFilterName string, ipFilters map[string][]string) crossplane.Directive {
	locationBlock := crossplane.Directive{
		Directive: "location",
		Args:      []string{path},
		Block: crossplane.Directives{
			{Directive: "proxy_pass", Args: []string{fmt.Sprintf("http://%s", appName)}},
		},
	}

	if ips, ok := ipFilters[ipFilterName]; ok {
		for _, ip := range ips {
			locationBlock.Block = append(locationBlock.Block,
				&crossplane.Directive{Directive: "allow", Args: []string{ip}},
			)
		}
		locationBlock.Block = append(locationBlock.Block,
			&crossplane.Directive{Directive: "deny", Args: []string{"all"}},
		)
	}

	return locationBlock
}

func createDefaultLocationBlock(appName string) crossplane.Directive {
	return crossplane.Directive{
		Directive: "location",
		Args:      []string{"/"},
		Block: crossplane.Directives{
			{Directive: "proxy_pass", Args: []string{fmt.Sprintf("http://%s", appName)}},
		},
	}
}
