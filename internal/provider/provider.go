package provider

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const defaultBaseURL = "https://api.herodesk.io/v1"

func Provider(_ context.Context) *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Herodesk API base URL.",
				DefaultFunc: schema.EnvDefaultFunc("HERODESK_BASE_URL", defaultBaseURL),
			},
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Herodesk API token.",
				DefaultFunc: schema.EnvDefaultFunc("HERODESK_TOKEN", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"herodesk_helpcenter": resourceEntity("helpcenters"),
			"herodesk_webhook":    resourceEntity("webhooks"),
			"herodesk_tag":        resourceEntity("tags"),
			"herodesk_team":       resourceEntity("teams"),
			"herodesk_sla_policy": resourceEntity("sla-policies"),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"herodesk_helpcenters":  dataSourceEntity("helpcenters"),
			"herodesk_webhooks":     dataSourceEntity("webhooks"),
			"herodesk_tags":         dataSourceEntity("tags"),
			"herodesk_teams":        dataSourceEntity("teams"),
			"herodesk_sla_policies": dataSourceEntity("sla-policies"),
		},
		ConfigureContextFunc: configure,
	}

	return p
}

func configure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	token := strings.TrimSpace(d.Get("token").(string))
	if token == "" {
		token = strings.TrimSpace(os.Getenv("HERODESK_TOKEN"))
	}
	if token == "" {
		return nil, diag.Errorf("provider token is required: set token or HERODESK_TOKEN")
	}

	baseURL := strings.TrimSpace(d.Get("base_url").(string))
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	client := &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	return client, nil
}

func ensureClient(meta interface{}) (*Client, error) {
	client, ok := meta.(*Client)
	if !ok || client == nil {
		return nil, fmt.Errorf("provider client not configured")
	}
	return client, nil
}
