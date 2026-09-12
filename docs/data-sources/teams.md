---
page_title: "herodesk_teams Data Source - Herodesk"
description: |-
  Retrieves Herodesk teams.
---

# herodesk_teams

Retrieves all Herodesk teams, including their members and assigned inboxes, or
one team when `team_id` is set. API filters are not exposed by this provider.

## Example Usage

```terraform
data "herodesk_teams" "all" {}

data "herodesk_teams" "support" {
  team_id = 123
}
```

## Schema

### Optional

- `team_id` (Number) Team ID to retrieve. Omit to retrieve the collection.

### Read-Only

- `items` (List of Object) Returned teams with `id`, `name`, `manager_user_id`,
  `users`, `inboxes`, `is_default`, `created_at`, and `updated_at`.