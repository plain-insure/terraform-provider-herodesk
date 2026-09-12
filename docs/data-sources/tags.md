---
page_title: "herodesk_tags Data Source - Herodesk"
description: |-
  Retrieves Herodesk tags.
---

# herodesk_tags

Retrieves all Herodesk tags, or one tag when `tag_id` is set. API filter parameters
such as name and archived status are not exposed by this provider.

## Example Usage

```terraform
data "herodesk_tags" "all" {}

data "herodesk_tags" "priority" {
  tag_id = 123
}
```

## Schema

### Optional

- `tag_id` (Number) Tag ID to retrieve. Omit to retrieve the collection.

### Read-Only

- `items` (List of Object) Returned tags with `id`, `name`, `archived`,
  `user_id`, `created_at`, and `updated_at`.