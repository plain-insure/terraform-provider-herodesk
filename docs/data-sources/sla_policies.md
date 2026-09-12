---
page_title: "herodesk_sla_policies Data Source - Herodesk"
description: |-
  Retrieves Herodesk SLA policies.
---

# herodesk_sla_policies

Retrieves all Herodesk SLA policies in their API evaluation order, or one policy
when `sla_policy_id` is set. This data source requires a Herodesk plan with SLA
support.

## Example Usage

```terraform
data "herodesk_sla_policies" "all" {}

data "herodesk_sla_policies" "priority" {
  sla_policy_id = 123
}
```

## Schema

### Optional

- `sla_policy_id` (Number) SLA policy ID to retrieve. Omit to retrieve the
  collection.

### Read-Only

- `items` (List of Object) Returned policies with `id`, `name`, `enabled`,
  `sort`, `conditions`, reply and resolution targets, `hours_mode`,
  `schedule_id`, `warn_percent`, and timestamps. Each condition exposes an
  `attribute` and a set of string `values`.