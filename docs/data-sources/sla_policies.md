---
page_title: "herodesk_sla_policies Data Source - Herodesk"
description: |-
  Retrieves Herodesk SLA policies.
---

# herodesk_sla_policies

Retrieves all Herodesk SLA policies in their API evaluation order, or one policy
when `id` is set. This data source requires a Herodesk plan with SLA support.

## Example Usage

```terraform
data "herodesk_sla_policies" "all" {}

data "herodesk_sla_policies" "priority" {
  id = "123"
}
```

## Schema

### Optional

- `id` (String) SLA policy ID to retrieve. Omit to retrieve the collection.

### Read-Only

- `items` (List of String) One normalized JSON string for each returned SLA
  policy.
- `raw_json` (String) Normalized JSON response returned by the API.