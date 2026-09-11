---
page_title: "herodesk_teams Data Source - Herodesk"
description: |-
  Retrieves Herodesk teams.
---

# herodesk_teams

Retrieves all Herodesk teams, including their members and assigned inboxes, or
one team when `id` is set. API filters are not exposed by this provider.

## Example Usage

```terraform
data "herodesk_teams" "all" {}

data "herodesk_teams" "support" {
  id = "123"
}
```

## Schema

### Optional

- `id` (String) Team ID to retrieve. Omit to retrieve the collection.

### Read-Only

- `items` (List of String) One normalized JSON string for each returned team.
- `raw_json` (String) Normalized JSON response returned by the API.