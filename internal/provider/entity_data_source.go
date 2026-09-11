package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEntity(endpoint string) *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEntityRead(endpoint),
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional object ID to fetch a single entry.",
			},
			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "JSON payloads returned from Herodesk.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"raw_json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Raw JSON response from Herodesk.",
			},
		},
	}
}

func dataSourceEntityRead(endpoint string) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		client, err := ensureClient(meta)
		if err != nil {
			return diag.FromErr(err)
		}

		lookupID := d.Get("id").(string)
		var response json.RawMessage
		if lookupID != "" {
			response, err = client.read(ctx, endpoint, lookupID)
		} else {
			response, err = client.list(ctx, endpoint)
		}
		if err != nil {
			return diag.FromErr(err)
		}

		items := flattenItems(response)
		if err := d.Set("items", items); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("raw_json", normalizeJSON(response)); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(fmt.Sprintf("%s-%d", endpoint, time.Now().UnixNano()))
		return nil
	}
}

func flattenItems(response json.RawMessage) []string {
	var array []json.RawMessage
	if err := json.Unmarshal(response, &array); err == nil {
		items := make([]string, 0, len(array))
		for _, item := range array {
			items = append(items, normalizeJSON(item))
		}
		return items
	}

	var object map[string]json.RawMessage
	if err := json.Unmarshal(response, &object); err == nil {
		for _, field := range []string{"items", "data", "results"} {
			if nested, ok := object[field]; ok {
				if err := json.Unmarshal(nested, &array); err == nil {
					items := make([]string, 0, len(array))
					for _, item := range array {
						items = append(items, normalizeJSON(item))
					}
					return items
				}
			}
		}
	}

	return []string{normalizeJSON(response)}
}
