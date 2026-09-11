---
page_title: "herodesk_helpcenters Data Source - Herodesk"
description: |-
  Retrieves Herodesk Help Centers.
---

# herodesk_helpcenters

Retrieves all Herodesk Help Centers, or one help center when `id` is set. The
Herodesk API can return localized fields; this data source uses the default API
response because it does not expose API query filters.

## Example Usage

```terraform
data "herodesk_helpcenters" "all" {}

data "herodesk_helpcenters" "support" {
  id = "123"
}
```

## Schema

### Optional

- `id` (String) Help center ID to retrieve. Omit to retrieve the collection.

### Read-Only

- `items` (List of String) One normalized JSON string for each returned help
  center.
- `raw_json` (String) Normalized JSON response returned by the API.