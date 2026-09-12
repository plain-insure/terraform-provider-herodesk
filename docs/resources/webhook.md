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
  name   = "Operations events"
  url    = "https://example.com/hooks/herodesk"
  active = true
  events = ["conversation.created", "conversation.updated"]
}
```

## Schema

### Required

- `name` (String) Human-readable webhook name.
- `url` (String) URL that receives webhook events.

### Optional

- `active` (Boolean) Whether the webhook receives events.
- `events` (Set of String) Subscribed event names. When changed, replaces the
  complete event subscription list.
- `secret` (String, Sensitive) Secret used to sign webhook payloads.

  Supported events are `conversation.created`, `conversation.updated`,
  `conversation.deleted`, `message.sent`, `message.received`, `contact.created`,
  `contact.updated`, `quick_reply.sent`, `rule.triggered`, `company.created`,
  `company.updated`, `conversation.watcher_added`, and
  `conversation.watcher_removed`.

### Read-Only

- `created_at` (String) UTC creation timestamp.
- `failed_attempts` (Number) Consecutive failed deliveries.
- `last_error` (String) Error from the most recent failed delivery.
- `last_failed_at` (String) UTC timestamp of the most recent failed delivery.
- `updated_at` (String) UTC update timestamp.

## Import

```shell
terraform import herodesk_webhook.events 123
```