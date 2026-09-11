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
  json = jsonencode({
    name            = "Support"
    manager_user_id = 42
    users           = [42, 43]
    inboxes         = [10, 11]
  })
}
```

## Schema

### Required

- `json` (String) JSON request body. `name` is required for creation.
  `manager_user_id` selects an active full user. `users` is an array of member
  IDs and `inboxes` is an array of standard inbox or Smart Folder IDs. Sending
  either array during an update replaces the complete assignment list. Set
  `is_default` to `1` to move the default-team flag to this team.

### Read-Only

- `id` (String) Herodesk team ID.

## Import

```shell
terraform import herodesk_team.support 123
```