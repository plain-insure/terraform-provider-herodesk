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
  name               = "Priority conversations"
  enabled            = 1
  sort               = 10
  conditions {
    attribute = "tag_id"
    values    = ["123"]
  }
  first_reply_target = 900
  resolution_target  = 14400
  hours_mode         = "calendar"
  warn_percent       = 80
}
```

## Schema

### Required

- `name` (String) Name of the SLA policy. The Herodesk API also requires at
  least one reply or resolution target for creation.

### Optional

- `backfill` (Boolean) Apply the policy to existing open conversations.
- `conditions` (Block List) Conditions that must all match for the policy to
  apply. Each block has an `attribute` of `inbox_id`, `channel_type`, `tag_id`,
  or `company_id`, and a set of string `values`.
- `enabled` (Number) `0` to disable the policy or `1` to enable it.
- `first_reply_target` (Number) First reply target in seconds.
- `hours_mode` (String) Target calculation mode: `calendar` or `business`.
- `next_reply_target` (Number) Next reply target in seconds.
- `resolution_target` (Number) Resolution target in seconds.
- `schedule_id` (Number) Business-hours schedule ID. Required by the API when
  `hours_mode` is `business`.
- `sort` (Number) Evaluation order. Lower values run first.
- `warn_percent` (Number) Warning threshold percentage from 1 through 99.

### Read-Only

- `created_at` (String) UTC creation timestamp.
- `updated_at` (String) UTC update timestamp.

## Import

```shell
terraform import herodesk_sla_policy.priority 123
```