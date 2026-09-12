---
page_title: "herodesk_tag Resource - Herodesk"
description: |-
  Manages a Herodesk tag.
---

# herodesk_tag

Manages a Herodesk tag. Tags used by one or more Smart Folders cannot be deleted
by the Herodesk API.

## Example Usage

```terraform
resource "herodesk_tag" "priority" {
  name     = "Priority"
  archived = 0
}
```

## Schema

### Required

- `name` (String) Name of the tag.

### Optional

- `archived` (Number) `0` for an active tag or `1` for an archived tag.

### Read-Only

- `created_at` (String) UTC creation timestamp.
- `updated_at` (String) UTC update timestamp.
- `user_id` (Number) ID of the user who created the tag.

## Import

```shell
terraform import herodesk_tag.priority 123
```