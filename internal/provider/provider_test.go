package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestProvider_RegistersRequestedResourcesAndDataSources(t *testing.T) {
	p := Provider(context.Background())

	resourceKeys := []string{
		"herodesk_helpcenter",
		"herodesk_webhook",
		"herodesk_tag",
		"herodesk_team",
		"herodesk_sla_policy",
	}
	for _, key := range resourceKeys {
		if _, ok := p.ResourcesMap[key]; !ok {
			t.Fatalf("expected resource %q to be registered", key)
		}
	}

	dataSourceKeys := []string{
		"herodesk_helpcenters",
		"herodesk_webhooks",
		"herodesk_tags",
		"herodesk_teams",
		"herodesk_sla_policies",
	}
	for _, key := range dataSourceKeys {
		if _, ok := p.DataSourcesMap[key]; !ok {
			t.Fatalf("expected data source %q to be registered", key)
		}
	}
}

func TestProviderConfigure_RequiresToken(t *testing.T) {
	t.Setenv("HERODESK_TOKEN", "")

	p := Provider(context.Background())
	d := schema.TestResourceDataRaw(t, p.Schema, map[string]interface{}{
		"token":    "",
		"base_url": defaultBaseURL,
	})

	_, diags := p.ConfigureContextFunc(context.Background(), d)
	if !diags.HasError() {
		t.Fatal("expected configuration to fail without token")
	}
}
