package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceHelpcenterCustomDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: helpcenterCustomDomainCreateOrUpdate,
		ReadContext:   helpcenterCustomDomainRead,
		UpdateContext: helpcenterCustomDomainCreateOrUpdate,
		DeleteContext: helpcenterCustomDomainDelete,
		Schema: map[string]*schema.Schema{
			"helpcenter_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the Help Center to configure.",
			},
			"custom_domain": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  "Custom domain served by the Help Center.",
			},
			"cname_domain": computedString("CNAME domain derived from the public help center URL."),
		},
	}
}

func helpcenterCustomDomainCreateOrUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if err := updateHelpcenterCustomDomain(ctx, data, meta, data.Get("custom_domain").(string)); err != nil {
		return diag.FromErr(err)
	}
	data.SetId(strconv.Itoa(data.Get("helpcenter_id").(int)))
	return helpcenterCustomDomainRead(ctx, data, meta)
}

func helpcenterCustomDomainRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := ensureClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	response, err := client.read(ctx, "helpcenters", strconv.Itoa(data.Get("helpcenter_id").(int)))
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
	customDomain, ok := object["custom_domain"].(string)
	if !ok || customDomain == "" {
		data.SetId("")
		return nil
	}
	if err := data.Set("custom_domain", customDomain); err != nil {
		return diag.FromErr(fmt.Errorf("set custom_domain: %w", err))
	}
	if err := data.Set("cname_domain", cnameDomain(object)); err != nil {
		return diag.FromErr(fmt.Errorf("set cname_domain: %w", err))
	}
	return nil
}

func helpcenterCustomDomainDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if err := updateHelpcenterCustomDomain(ctx, data, meta, ""); err != nil && !strings.Contains(err.Error(), "status 404") {
		return diag.FromErr(err)
	}
	data.SetId("")
	return nil
}

func updateHelpcenterCustomDomain(ctx context.Context, data *schema.ResourceData, meta interface{}, customDomain string) error {
	client, err := ensureClient(meta)
	if err != nil {
		return err
	}
	payload, err := json.Marshal(map[string]string{"custom_domain": customDomain})
	if err != nil {
		return err
	}
	_, err = client.update(ctx, "helpcenters", strconv.Itoa(data.Get("helpcenter_id").(int)), payload)
	return err
}
