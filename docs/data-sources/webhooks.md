---
page_title: "herodesk_webhooks Data Source - Herodesk"
description: |-
  Retrieves Herodesk webhooks.
---

# herodesk_webhooks

Retrieves all Herodesk webhooks, or one webhook when `id` is set. API filter
parameters such as name, URL, and active state are not exposed by this provider.

## Example Usage

```terraform
data "herodesk_webhooks" "all" {}

data "herodesk_webhooks" "events" {
  id = "123"
}
```

## Schema

### Optional

- `id` (Number) Webhook ID to retrieve. Omit to retrieve the collection.

### Read-Only

- `items` (List of Object) Returned webhooks with `id`, `name`, `url`, `secret`,
  `active`, `events`, delivery status fields, and timestamps. `secret` is
  sensitive.