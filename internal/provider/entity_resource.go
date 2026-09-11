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

func resourceEntity(endpoint string) *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEntityCreate(endpoint),
		ReadContext:   resourceEntityRead(endpoint),
		UpdateContext: resourceEntityUpdate(endpoint),
		DeleteContext: resourceEntityDelete(endpoint),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Herodesk object ID.",
			},
			"json": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateJSON,
				Description:      "JSON payload sent to Herodesk.",
			},
		},
	}
}

func validateJSON(v interface{}, path cty.Path) diag.Diagnostics {
	value := strings.TrimSpace(v.(string))
	if !json.Valid([]byte(value)) {
		return diag.Diagnostics{diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Invalid JSON payload",
			Detail:        fmt.Sprintf("attribute at %v must be valid JSON", path),
			AttributePath: path,
		}}
	}
	return nil
}

func resourceEntityCreate(endpoint string) schema.CreateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		client, err := ensureClient(meta)
		if err != nil {
			return diag.FromErr(err)
		}

		payload := json.RawMessage(d.Get("json").(string))
		response, err := client.create(ctx, endpoint, payload)
		if err != nil {
			return diag.FromErr(err)
		}

		id := extractID(response)
		if id == "" {
			id = extractID(payload)
		}
		if id == "" {
			return diag.Errorf("herodesk response did not include an id")
		}

		d.SetId(id)
		if err := d.Set("json", normalizeJSON(response)); err != nil {
			return diag.FromErr(err)
		}

		return resourceEntityRead(endpoint)(ctx, d, meta)
	}
}

func resourceEntityRead(endpoint string) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		client, err := ensureClient(meta)
		if err != nil {
			return diag.FromErr(err)
		}

		response, err := client.read(ctx, endpoint, d.Id())
		if err != nil {
			if strings.Contains(err.Error(), "status 404") {
				d.SetId("")
				return nil
			}
			return diag.FromErr(err)
		}

		if err := d.Set("json", normalizeJSON(response)); err != nil {
			return diag.FromErr(err)
		}

		return nil
	}
}

func resourceEntityUpdate(endpoint string) schema.UpdateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		client, err := ensureClient(meta)
		if err != nil {
			return diag.FromErr(err)
		}

		payload := json.RawMessage(d.Get("json").(string))
		response, err := client.update(ctx, endpoint, d.Id(), payload)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("json", normalizeJSON(response)); err != nil {
			return diag.FromErr(err)
		}

		return resourceEntityRead(endpoint)(ctx, d, meta)
	}
}

func resourceEntityDelete(endpoint string) schema.DeleteContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		client, err := ensureClient(meta)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := client.delete(ctx, endpoint, d.Id()); err != nil {
			if strings.Contains(err.Error(), "status 404") {
				return nil
			}
			return diag.FromErr(err)
		}

		d.SetId("")
		return nil
	}
}

func normalizeJSON(v json.RawMessage) string {
	if len(v) == 0 {
		return "{}"
	}
	var decoded interface{}
	if err := json.Unmarshal(v, &decoded); err != nil {
		return strings.TrimSpace(string(v))
	}
	encoded, err := json.Marshal(decoded)
	if err != nil {
		return strings.TrimSpace(string(v))
	}
	return string(encoded)
}

func extractID(v json.RawMessage) string {
	var obj map[string]interface{}
	if err := json.Unmarshal(v, &obj); err != nil {
		return ""
	}
	if raw, ok := obj["id"]; ok {
		return fmt.Sprintf("%v", raw)
	}
	if raw, ok := obj["_id"]; ok {
		return fmt.Sprintf("%v", raw)
	}
	return ""
}
