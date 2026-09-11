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
  json = jsonencode({
    name     = "Priority"
    archived = 0
  })
}
```

## Schema

### Required

- `json` (String) JSON request body. `name` is required for creation.
  `archived` is optional and uses `0` for active or `1` for archived.

### Read-Only

- `id` (String) Herodesk tag ID.

## Import

```shell
terraform import herodesk_tag.priority 123
```