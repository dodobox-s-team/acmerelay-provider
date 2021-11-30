package acmerelay

import (
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

// Provider wraps the provider implementation as a Caddy module.
type ProviderWrapper struct{ *Provider }

func init() {
	caddy.RegisterModule(ProviderWrapper{})
}

// CaddyModule returns the Caddy module information.
func (ProviderWrapper) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.acmerelay",
		New: func() caddy.Module { return &ProviderWrapper{new(Provider)} },
	}
}

// Before using the provider config, resolve placeholders in the API token.
// Implements caddy.Provisioner.
func (p *ProviderWrapper) Provision(ctx caddy.Context) error {
	p.Provider.APIKey = caddy.NewReplacer().ReplaceAll(p.Provider.APIKey, "")
	p.Provider.APIEndpoint = caddy.NewReplacer().ReplaceAll(p.Provider.APIEndpoint, "")
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
// acmerelay {
//     api_key <api_key>
//     endpoint <url>
// }
//
// Expansion of placeholders is left to the JSON config caddy.Provisioner (above).
func (p *ProviderWrapper) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "api_key":
				if p.Provider.APIKey != "" {
					return d.Err("API key already set")
				}
				if !d.NextArg() {
					return d.ArgErr()
				}
				p.Provider.APIKey = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}
			case "endpoint":
				if p.Provider.APIEndpoint != "" {
					return d.Err("API endpoint already set")
				}
				if !d.NextArg() {
					return d.ArgErr()
				}
				p.Provider.APIEndpoint = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	if p.Provider.APIKey == "" {
		return d.Err("missing API key")
	}
	if p.Provider.APIEndpoint == "" {
		// Set the default url
		p.Provider.APIEndpoint = "https://api.dodobox.site/acmerelay"
	} else {
		p.Provider.APIEndpoint = strings.TrimSuffix(p.Provider.APIEndpoint, "/")
	}
	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*ProviderWrapper)(nil)
	_ caddy.Provisioner     = (*ProviderWrapper)(nil)
)
