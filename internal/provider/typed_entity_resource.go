package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type entityDefinition struct {
	endpoint string
	schema   map[string]*schema.Schema
	payload  func(*schema.ResourceData) map[string]interface{}
	flatten  func(*schema.ResourceData, map[string]interface{}) error
}

func typedEntity(definition entityDefinition) *schema.Resource {
	return &schema.Resource{
		CreateContext: typedEntityCreate(definition),
		ReadContext:   typedEntityRead(definition),
		UpdateContext: typedEntityUpdate(definition),
		DeleteContext: typedEntityDelete(definition),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: definition.schema,
	}
}

func typedEntityCreate(definition entityDefinition) schema.CreateContextFunc {
	return func(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
		client, err := ensureClient(meta)
		if err != nil {
			return diag.FromErr(err)
		}

		payload, err := json.Marshal(definition.payload(data))
		if err != nil {
			return diag.FromErr(err)
		}
		response, err := client.create(ctx, definition.endpoint, payload)
		if err != nil {
			return diag.FromErr(err)
		}

		id := extractID(response)
		if id == "" {
			return diag.Errorf("herodesk response did not include an id")
		}
		data.SetId(id)
		return typedEntityRead(definition)(ctx, data, meta)
	}
}

func typedEntityRead(definition entityDefinition) schema.ReadContextFunc {
	return func(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
		client, err := ensureClient(meta)
		if err != nil {
			return diag.FromErr(err)
		}

		response, err := client.read(ctx, definition.endpoint, data.Id())
		if err != nil {
			if strings.Contains(err.Error(), "status 404") {
				data.SetId("")
				return nil
			}
			return diag.FromErr(err)
		}

		object, err := responseObject(response)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := definition.flatten(data, object); err != nil {
			return diag.FromErr(err)
		}
		return nil
	}
}

func typedEntityUpdate(definition entityDefinition) schema.UpdateContextFunc {
	return func(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
		client, err := ensureClient(meta)
		if err != nil {
			return diag.FromErr(err)
		}

		payload, err := json.Marshal(definition.payload(data))
		if err != nil {
			return diag.FromErr(err)
		}
		if _, err := client.update(ctx, definition.endpoint, data.Id(), payload); err != nil {
			return diag.FromErr(err)
		}
		return typedEntityRead(definition)(ctx, data, meta)
	}
}

func typedEntityDelete(definition entityDefinition) schema.DeleteContextFunc {
	return func(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
		client, err := ensureClient(meta)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := client.delete(ctx, definition.endpoint, data.Id()); err != nil && !strings.Contains(err.Error(), "status 404") {
			return diag.FromErr(err)
		}
		data.SetId("")
		return nil
	}
}

func responseObject(response json.RawMessage) (map[string]interface{}, error) {
	var object map[string]interface{}
	if err := json.Unmarshal(response, &object); err != nil {
		return nil, fmt.Errorf("decode herodesk response: %w", err)
	}
	return object, nil
}

func schemaPayload(data *schema.ResourceData, fields ...string) map[string]interface{} {
	return schemaPayloadFromConfig(data, data.GetRawConfig(), fields...)
}

func schemaPayloadFromConfig(data *schema.ResourceData, rawConfig cty.Value, fields ...string) map[string]interface{} {
	payload := make(map[string]interface{}, len(fields))
	for _, field := range fields {
		if rawConfigHasField(rawConfig, field) {
			payload[field] = data.Get(field)
		}
	}
	return payload
}

func rawConfigHasField(rawConfig cty.Value, field string) bool {
	if !rawConfig.IsKnown() || rawConfig.IsNull() || !rawConfig.Type().IsObjectType() || !rawConfig.Type().HasAttribute(field) {
		return false
	}
	value := rawConfig.GetAttr(field)
	return value.IsKnown() && !value.IsNull()
}

func flattenSchema(resourceSchema map[string]*schema.Schema) func(*schema.ResourceData, map[string]interface{}) error {
	return func(data *schema.ResourceData, object map[string]interface{}) error {
		for field := range resourceSchema {
			if value, ok := object[field]; ok {
				if err := data.Set(field, value); err != nil {
					return fmt.Errorf("set %s: %w", field, err)
				}
			}
		}
		return nil
	}
}

func resourceHelpcenter() *schema.Resource {
	resourceSchema := map[string]*schema.Schema{
		"name":                requiredString("Name of the help center."),
		"language":            requiredString("Default language of the help center."),
		"description":         optionalComputedString("Description shown on the help center front page."),
		"available_languages": optionalComputedStringList("Additional translated languages."),
		"default_domain":      computedString("Default Herodesk help center subdomain."),
		"public_url":          computedString("Public URL of the help center."),
		"root_folder_id":      computedInt("ID of the automatically created root folder."),
		"created_at":          computedString("UTC timestamp when the help center was created."),
		"updated_at":          computedString("UTC timestamp when the help center was last updated."),
	}
	return typedEntity(entityDefinition{
		endpoint: "helpcenters",
		schema:   resourceSchema,
		payload: func(data *schema.ResourceData) map[string]interface{} {
			return schemaPayload(data, "name", "language", "description", "available_languages")
		},
		flatten: flattenSchema(resourceSchema),
	})
}

func resourceWebhook() *schema.Resource {
	resourceSchema := map[string]*schema.Schema{
		"name":            requiredString("Human-readable name for the webhook."),
		"url":             requiredString("URL that receives webhook events."),
		"secret":          optionalComputedSensitiveString("Secret used to sign webhook payloads."),
		"active":          optionalComputedBool("Whether the webhook receives events."),
		"events":          optionalComputedStringSet("Webhook events to subscribe to."),
		"failed_attempts": computedInt("Number of consecutive failed delivery attempts."),
		"last_failed_at":  computedString("UTC timestamp of the most recent failed delivery."),
		"last_error":      computedString("Error from the most recent failed delivery."),
		"created_at":      computedString("UTC timestamp when the webhook was created."),
		"updated_at":      computedString("UTC timestamp when the webhook was last updated."),
	}
	return typedEntity(entityDefinition{
		endpoint: "webhooks",
		schema:   resourceSchema,
		payload: func(data *schema.ResourceData) map[string]interface{} {
			return schemaPayload(data, "name", "url", "secret", "active", "events")
		},
		flatten: flattenSchema(resourceSchema),
	})
}

func resourceTag() *schema.Resource {
	resourceSchema := map[string]*schema.Schema{
		"name":       requiredString("Name of the tag."),
		"archived":   optionalComputedInt("Whether the tag is archived: 0 for active or 1 for archived."),
		"user_id":    computedInt("ID of the user who created the tag."),
		"created_at": computedString("UTC timestamp when the tag was created."),
		"updated_at": computedString("UTC timestamp when the tag was last updated."),
	}
	return typedEntity(entityDefinition{
		endpoint: "tags",
		schema:   resourceSchema,
		payload: func(data *schema.ResourceData) map[string]interface{} {
			return schemaPayload(data, "name", "archived")
		},
		flatten: flattenSchema(resourceSchema),
	})
}

func resourceTeam() *schema.Resource {
	resourceSchema := map[string]*schema.Schema{
		"name":            requiredString("Unique name of the team."),
		"manager_user_id": optionalComputedInt("ID of the active user who manages the team."),
		"users":           optionalComputedIntSet("IDs of users assigned to the team."),
		"inboxes":         optionalComputedIntSet("IDs of inboxes and Smart Folders assigned to the team."),
		"is_default":      optionalComputedInt("Whether the team is the default team: 0 or 1."),
		"created_at":      computedString("UTC timestamp when the team was created."),
		"updated_at":      computedString("UTC timestamp when the team was last updated."),
	}
	return typedEntity(entityDefinition{
		endpoint: "teams",
		schema:   resourceSchema,
		payload: func(data *schema.ResourceData) map[string]interface{} {
			return schemaPayload(data, "name", "manager_user_id", "users", "inboxes", "is_default")
		},
		flatten: flattenSchema(resourceSchema),
	})
}

func resourceSLAPolicy() *schema.Resource {
	resourceSchema := map[string]*schema.Schema{
		"name":               requiredString("Name of the SLA policy."),
		"enabled":            optionalComputedInt("Whether the SLA policy is enabled: 0 or 1."),
		"sort":               optionalComputedInt("Evaluation order; lower values run first."),
		"conditions":         optionalComputedConditions(),
		"first_reply_target": optionalComputedInt("First reply target in seconds."),
		"next_reply_target":  optionalComputedInt("Next reply target in seconds."),
		"resolution_target":  optionalComputedInt("Resolution target in seconds."),
		"hours_mode":         optionalComputedString("Target calculation mode: calendar or business."),
		"schedule_id":        optionalComputedInt("Business-hours schedule ID."),
		"warn_percent":       optionalComputedInt("Warning threshold percentage."),
		"backfill":           optionalComputedBool("Apply the policy to existing open conversations."),
		"created_at":         computedString("UTC timestamp when the policy was created."),
		"updated_at":         computedString("UTC timestamp when the policy was last updated."),
	}
	return typedEntity(entityDefinition{
		endpoint: "sla-policies",
		schema:   resourceSchema,
		payload:  slaPolicyPayload,
		flatten:  flattenSLAPolicy(resourceSchema),
	})
}

func requiredString(description string) *schema.Schema {
	return &schema.Schema{Type: schema.TypeString, Required: true, Description: description}
}

func optionalComputedString(description string) *schema.Schema {
	return &schema.Schema{Type: schema.TypeString, Optional: true, Computed: true, Description: description}
}

func optionalComputedSensitiveString(description string) *schema.Schema {
	return &schema.Schema{Type: schema.TypeString, Optional: true, Computed: true, Sensitive: true, Description: description}
}

func computedString(description string) *schema.Schema {
	return &schema.Schema{Type: schema.TypeString, Computed: true, Description: description}
}

func optionalComputedInt(description string) *schema.Schema {
	return &schema.Schema{Type: schema.TypeInt, Optional: true, Computed: true, Description: description}
}

func computedInt(description string) *schema.Schema {
	return &schema.Schema{Type: schema.TypeInt, Computed: true, Description: description}
}

func optionalComputedBool(description string) *schema.Schema {
	return &schema.Schema{Type: schema.TypeBool, Optional: true, Computed: true, Description: description}
}

func optionalComputedStringList(description string) *schema.Schema {
	return &schema.Schema{Type: schema.TypeList, Optional: true, Computed: true, Description: description, Elem: &schema.Schema{Type: schema.TypeString}}
}

func optionalComputedStringSet(description string) *schema.Schema {
	return &schema.Schema{Type: schema.TypeSet, Optional: true, Computed: true, Description: description, Elem: &schema.Schema{Type: schema.TypeString}}
}

func optionalComputedIntSet(description string) *schema.Schema {
	return &schema.Schema{Type: schema.TypeSet, Optional: true, Computed: true, Description: description, Elem: &schema.Schema{Type: schema.TypeInt}}
}

func optionalComputedConditions() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Computed:    true,
		Description: "Conditions that must all match for the policy to apply.",
		Elem: &schema.Resource{Schema: map[string]*schema.Schema{
			"attribute": {Type: schema.TypeString, Required: true},
			"values":    {Type: schema.TypeSet, Required: true, Elem: &schema.Schema{Type: schema.TypeString}},
		}},
	}
}

func slaPolicyPayload(data *schema.ResourceData) map[string]interface{} {
	payload := schemaPayload(data, "name", "enabled", "sort", "first_reply_target", "next_reply_target", "resolution_target", "hours_mode", "schedule_id", "warn_percent", "backfill")
	if rawConfigHasField(data.GetRawConfig(), "conditions") {
		payload["conditions"] = flattenSLAConditions(data.Get("conditions").([]interface{}))
	}
	return payload
}

func flattenSLAConditions(conditions []interface{}) []map[string]interface{} {
	payload := make([]map[string]interface{}, 0, len(conditions))
	for _, condition := range conditions {
		values := condition.(map[string]interface{})["values"].(*schema.Set).List()
		payload = append(payload, map[string]interface{}{
			"attribute": condition.(map[string]interface{})["attribute"],
			"value":     values,
		})
	}
	return payload
}

func flattenSLAPolicy(resourceSchema map[string]*schema.Schema) func(*schema.ResourceData, map[string]interface{}) error {
	return func(data *schema.ResourceData, object map[string]interface{}) error {
		if conditions, ok := object["conditions"].([]interface{}); ok {
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
			object["conditions"] = flattened
		}
		return flattenSchema(resourceSchema)(data, object)
	}
}
