package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type dataSourceDefinition struct {
	endpoint   string
	itemSchema map[string]*schema.Schema
	flatten    func(map[string]interface{}) map[string]interface{}
}

func typedDataSource(definition dataSourceDefinition) *schema.Resource {
	return &schema.Resource{
		ReadContext: typedDataSourceRead(definition),
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Optional Herodesk object ID to retrieve.",
			},
			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Objects returned by Herodesk.",
				Elem: &schema.Resource{
					Schema: definition.itemSchema,
				},
			},
		},
	}
}

func typedDataSourceRead(definition dataSourceDefinition) schema.ReadContextFunc {
	return func(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
		client, err := ensureClient(meta)
		if err != nil {
			return diag.FromErr(err)
		}

		lookupID := data.Get("id").(int)
		var response json.RawMessage
		if lookupID != 0 {
			response, err = client.read(ctx, definition.endpoint, fmt.Sprint(lookupID))
		} else {
			response, err = client.list(ctx, definition.endpoint)
		}
		if err != nil {
			return diag.FromErr(err)
		}

		items, err := responseObjects(response)
		if err != nil {
			return diag.FromErr(err)
		}
		for index := range items {
			items[index] = definition.flatten(items[index])
		}
		if err := data.Set("items", items); err != nil {
			return diag.FromErr(err)
		}

		data.SetId(fmt.Sprintf("%s-%d", definition.endpoint, time.Now().UnixNano()))
		return nil
	}
}

func responseObjects(response json.RawMessage) ([]map[string]interface{}, error) {
	var items []map[string]interface{}
	if err := json.Unmarshal(response, &items); err == nil {
		return items, nil
	}

	var object map[string]interface{}
	if err := json.Unmarshal(response, &object); err != nil {
		return nil, fmt.Errorf("decode herodesk response: %w", err)
	}
	for _, field := range []string{"items", "data", "results"} {
		if nested, ok := object[field]; ok {
			encoded, err := json.Marshal(nested)
			if err != nil {
				return nil, fmt.Errorf("encode response %s: %w", field, err)
			}
			if err := json.Unmarshal(encoded, &items); err == nil {
				return items, nil
			}
		}
	}
	return []map[string]interface{}{object}, nil
}

func dataSourceTag() *schema.Resource {
	return typedDataSource(dataSourceDefinition{
		endpoint: "tags",
		itemSchema: map[string]*schema.Schema{
			"id":         computedInt("Tag ID."),
			"name":       computedString("Name of the tag."),
			"archived":   computedInt("Whether the tag is archived: 0 or 1."),
			"user_id":    computedInt("ID of the user who created the tag."),
			"created_at": computedString("UTC timestamp when the tag was created."),
			"updated_at": computedString("UTC timestamp when the tag was last updated."),
		},
		flatten: func(item map[string]interface{}) map[string]interface{} {
			return item
		},
	})
}

func dataSourceHelpcenter() *schema.Resource {
	return typedDataSource(dataSourceDefinition{
		endpoint: "helpcenters",
		itemSchema: map[string]*schema.Schema{
			"id":                  computedInt("Help center ID."),
			"name":                computedString("Name of the help center."),
			"language":            computedString("Default language of the help center."),
			"description":         computedString("Description shown on the help center front page."),
			"available_languages": computedStringList("Additional translated languages."),
			"custom_domain":       computedString("Custom domain for the help center."),
			"default_domain":      computedString("Default Herodesk help center subdomain."),
			"public_url":          computedString("Public URL of the help center."),
			"root_folder_id":      computedInt("Automatically created root folder ID."),
			"created_at":          computedString("UTC timestamp when the help center was created."),
			"updated_at":          computedString("UTC timestamp when the help center was last updated."),
		},
		flatten: identityItem,
	})
}

func dataSourceWebhook() *schema.Resource {
	return typedDataSource(dataSourceDefinition{
		endpoint: "webhooks",
		itemSchema: map[string]*schema.Schema{
			"id":              computedInt("Webhook ID."),
			"name":            computedString("Human-readable webhook name."),
			"url":             computedString("URL that receives webhook events."),
			"secret":          computedSensitiveString("Secret used to sign webhook payloads."),
			"active":          computedBool("Whether the webhook receives events."),
			"events":          computedStringSet("Webhook events."),
			"failed_attempts": computedInt("Number of consecutive failed delivery attempts."),
			"last_failed_at":  computedString("UTC timestamp of the most recent failed delivery."),
			"last_error":      computedString("Error from the most recent failed delivery."),
			"created_at":      computedString("UTC timestamp when the webhook was created."),
			"updated_at":      computedString("UTC timestamp when the webhook was last updated."),
		},
		flatten: identityItem,
	})
}

func dataSourceTeam() *schema.Resource {
	return typedDataSource(dataSourceDefinition{
		endpoint: "teams",
		itemSchema: map[string]*schema.Schema{
			"id":              computedInt("Team ID."),
			"name":            computedString("Team name."),
			"manager_user_id": computedInt("ID of the team manager."),
			"users":           computedIntSet("User IDs assigned to the team."),
			"inboxes":         computedIntSet("Inbox and Smart Folder IDs assigned to the team."),
			"is_default":      computedInt("Whether this is the default team: 0 or 1."),
			"created_at":      computedString("UTC timestamp when the team was created."),
			"updated_at":      computedString("UTC timestamp when the team was last updated."),
		},
		flatten: identityItem,
	})
}

func dataSourceSLAPolicy() *schema.Resource {
	return typedDataSource(dataSourceDefinition{
		endpoint: "sla-policies",
		itemSchema: map[string]*schema.Schema{
			"id":                 computedInt("SLA policy ID."),
			"name":               computedString("SLA policy name."),
			"enabled":            computedInt("Whether the SLA policy is enabled: 0 or 1."),
			"sort":               computedInt("Evaluation order; lower values run first."),
			"conditions":         computedConditions(),
			"first_reply_target": computedInt("First reply target in seconds."),
			"next_reply_target":  computedInt("Next reply target in seconds."),
			"resolution_target":  computedInt("Resolution target in seconds."),
			"hours_mode":         computedString("Target calculation mode."),
			"schedule_id":        computedInt("Business-hours schedule ID."),
			"warn_percent":       computedInt("Warning threshold percentage."),
			"created_at":         computedString("UTC timestamp when the policy was created."),
			"updated_at":         computedString("UTC timestamp when the policy was last updated."),
		},
		flatten: flattenSLAPolicyItem,
	})
}

func identityItem(item map[string]interface{}) map[string]interface{} {
	return item
}

func flattenSLAPolicyItem(item map[string]interface{}) map[string]interface{} {
	conditions, ok := item["conditions"].([]interface{})
	if !ok {
		return item
	}
	flattened := make([]map[string]interface{}, 0, len(conditions))
	for _, condition := range conditions {
		entry := condition.(map[string]interface{})
		values, _ := entry["value"].([]interface{})
		stringValues := make([]string, 0, len(values))
		for _, value := range values {
			stringValues = append(stringValues, fmt.Sprint(value))
		}
		flattened = append(flattened, map[string]interface{}{"attribute": entry["attribute"], "values": stringValues})
	}
	item["conditions"] = flattened
	return item
}

func computedSensitiveString(description string) *schema.Schema {
	return &schema.Schema{Type: schema.TypeString, Computed: true, Sensitive: true, Description: description}
}

func computedBool(description string) *schema.Schema {
	return &schema.Schema{Type: schema.TypeBool, Computed: true, Description: description}
}

func computedStringList(description string) *schema.Schema {
	return &schema.Schema{Type: schema.TypeList, Computed: true, Description: description, Elem: &schema.Schema{Type: schema.TypeString}}
}

func computedStringSet(description string) *schema.Schema {
	return &schema.Schema{Type: schema.TypeSet, Computed: true, Description: description, Elem: &schema.Schema{Type: schema.TypeString}}
}

func computedIntSet(description string) *schema.Schema {
	return &schema.Schema{Type: schema.TypeSet, Computed: true, Description: description, Elem: &schema.Schema{Type: schema.TypeInt}}
}

func computedConditions() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Conditions that must all match for the SLA policy to apply.",
		Elem: &schema.Resource{Schema: map[string]*schema.Schema{
			"attribute": {Type: schema.TypeString, Computed: true},
			"values":    {Type: schema.TypeSet, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
		}},
	}
}
