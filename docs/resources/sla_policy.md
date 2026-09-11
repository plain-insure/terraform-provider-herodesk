---
page_title: "herodesk_sla_policy Resource - Herodesk"
description: |-
  Manages a Herodesk SLA policy.
---

# herodesk_sla_policy

Manages a Herodesk SLA policy. Policies are evaluated in ascending `sort`
order, and the first matching policy wins. This resource requires a Herodesk
plan that includes SLA functionality.

## Example Usage

```terraform
resource "herodesk_sla_policy" "priority" {
  json = jsonencode({
    name               = "Priority conversations"
    enabled            = 1
    sort               = 10
    conditions         = [{ attribute = "tag_id", value = [123] }]
    first_reply_target = 900
    resolution_target  = 14400
    hours_mode         = "calendar"
    warn_percent       = 80
  })
}
```

## Schema

### Required

- `json` (String) JSON request body. `name` and at least one of
  `first_reply_target`, `next_reply_target`, or `resolution_target` are required
  for creation. Targets are measured in seconds.

  `conditions` is an array of objects with `attribute` set to `inbox_id`,
  `channel_type`, `tag_id`, or `company_id`, and a `value` array. Other fields
  are `enabled` (`0` or `1`), `sort`, `hours_mode` (`calendar` or `business`),
  `schedule_id` (required for business hours), `warn_percent` (1-99), and
  `backfill`.

### Read-Only

- `id` (String) Herodesk SLA policy ID.

## Import

```shell
terraform import herodesk_sla_policy.priority 123
```