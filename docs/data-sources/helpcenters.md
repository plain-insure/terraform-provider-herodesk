---
page_title: "herodesk_helpcenters Data Source - Herodesk"
description: |-
  Retrieves Herodesk Help Centers.
---

# herodesk_helpcenters

Retrieves all Herodesk Help Centers, or one help center when `helpcenter_id` is set. The
Herodesk API can return localized fields; this data source uses the default API
response because it does not expose API query filters.

## Example Usage

```terraform
data "herodesk_helpcenters" "all" {}

data "herodesk_helpcenters" "support" {
  helpcenter_id = 123
}
```

## Schema

### Optional

- `helpcenter_id` (Number) Help center ID to retrieve. Omit to retrieve the
  collection.

### Read-Only

- `items` (List of Object) Returned help centers with `id`, `name`, `language`,
  `description`, `available_languages`, `custom_domain`, `default_domain`,
  `public_url`, `root_folder_id`, `created_at`, and `updated_at`.