---
page_title: "herodesk_webhook Resource - Herodesk"
description: |-
  Manages a Herodesk webhook.
---

# herodesk_webhook

Manages a webhook that receives Herodesk events. The API requires `name` and
`url` when creating a webhook. A `secret` is generated when it is omitted.

## Example Usage

```terraform
resource "herodesk_webhook" "events" {
  json = jsonencode({
    name   = "Operations events"
    url    = "https://example.com/hooks/herodesk"
    active = true
    events = ["conversation.created", "conversation.updated"]
  })
}
```

## Schema

### Required

- `json` (String) JSON request body. Creation supports `name`, `url`, `secret`,
  `active`, and `events`. `events` replaces the complete subscription list when
  included in an update.

  Supported events are `conversation.created`, `conversation.updated`,
  `conversation.deleted`, `message.sent`, `message.received`, `contact.created`,
  `contact.updated`, `quick_reply.sent`, `rule.triggered`, `company.created`,
  `company.updated`, `conversation.watcher_added`, and
  `conversation.watcher_removed`.

### Read-Only

- `id` (String) Herodesk webhook ID.

## Import

```shell
terraform import herodesk_webhook.events 123
```