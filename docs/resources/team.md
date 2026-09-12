---
page_title: "herodesk_team Resource - Herodesk"
description: |-
  Manages a Herodesk team.
---

# herodesk_team

Manages a Herodesk team, including member and inbox assignments. The team name
must be unique. The API automatically includes the manager as a team member.
The default team cannot be deleted.

## Example Usage

```terraform
resource "herodesk_team" "support" {
  name            = "Support"
  manager_user_id = 42
  users           = [42, 43]
  inboxes         = [10, 11]
}
```

## Schema

### Required

- `name` (String) Unique team name.

### Optional

- `inboxes` (Set of Number) Standard inbox and Smart Folder IDs assigned to the
  team. Updating the value replaces the complete assignment list.
- `is_default` (Number) Set to `1` to move the default-team flag to this team.
- `manager_user_id` (Number) Active full user who manages the team.
- `users` (Set of Number) User IDs assigned to the team. Updating the value
  replaces the complete member list.

### Read-Only

- `created_at` (String) UTC creation timestamp.
- `updated_at` (String) UTC update timestamp.

## Import

```shell
terraform import herodesk_team.support 123
```