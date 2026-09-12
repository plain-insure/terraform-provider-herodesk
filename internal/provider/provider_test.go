package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/go-cty/cty"
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

func TestProviderResourcesExposeTypedSchemas(t *testing.T) {
	p := Provider(context.Background())

	for resourceName, expectedFields := range map[string][]string{
		"herodesk_helpcenter": {"name", "language", "description"},
		"herodesk_webhook":    {"name", "url", "events"},
		"herodesk_tag":        {"name", "archived"},
		"herodesk_team":       {"name", "users", "inboxes"},
		"herodesk_sla_policy": {"name", "conditions", "hours_mode"},
	} {
		resourceSchema := p.ResourcesMap[resourceName].Schema
		if _, ok := resourceSchema["json"]; ok {
			t.Errorf("resource %q must not expose a json argument", resourceName)
		}
		for _, field := range expectedFields {
			if _, ok := resourceSchema[field]; !ok {
				t.Errorf("resource %q is missing field %q", resourceName, field)
			}
		}
	}
}

func TestProviderDataSourcesExposeTypedItems(t *testing.T) {
	p := Provider(context.Background())

	for _, dataSourceName := range []string{
		"herodesk_helpcenters",
		"herodesk_webhooks",
		"herodesk_tags",
		"herodesk_teams",
		"herodesk_sla_policies",
	} {
		dataSourceSchema := p.DataSourcesMap[dataSourceName].Schema
		if _, ok := dataSourceSchema["raw_json"]; ok {
			t.Errorf("data source %q must not expose raw_json", dataSourceName)
		}
		items, ok := dataSourceSchema["items"]
		if !ok || items.Type != schema.TypeList {
			t.Errorf("data source %q must expose items as a list", dataSourceName)
			continue
		}
		if _, ok := items.Elem.(*schema.Resource); !ok {
			t.Errorf("data source %q items must contain typed objects", dataSourceName)
		}
	}
}

func TestSchemaPayloadPreservesConfiguredZeroValues(t *testing.T) {
	tagData := schema.TestResourceDataRaw(t, resourceTag().Schema, map[string]interface{}{
		"name":     "Priority",
		"archived": 0,
	})
	tagPayload := schemaPayloadFromConfig(tagData, cty.ObjectVal(map[string]cty.Value{
		"name":     cty.StringVal("Priority"),
		"archived": cty.NumberIntVal(0),
	}), "name", "archived")
	if value, ok := tagPayload["archived"]; !ok || value != 0 {
		t.Errorf("tag archived payload = %#v, want 0", value)
	}

	webhookData := schema.TestResourceDataRaw(t, resourceWebhook().Schema, map[string]interface{}{
		"name":   "Events",
		"url":    "https://example.com/hooks",
		"active": false,
	})
	webhookPayload := schemaPayloadFromConfig(webhookData, cty.ObjectVal(map[string]cty.Value{
		"name":   cty.StringVal("Events"),
		"url":    cty.StringVal("https://example.com/hooks"),
		"active": cty.False,
	}), "name", "url", "active")
	if value, ok := webhookPayload["active"]; !ok || value != false {
		t.Errorf("webhook active payload = %#v, want false", value)
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
