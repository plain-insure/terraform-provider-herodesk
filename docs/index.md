---
page_title: "Provider: Herodesk"
description: |-
  The Herodesk provider manages Help Centers, webhooks, tags, teams, and SLA policies in Herodesk.
---

# Herodesk Provider

The Herodesk provider manages selected configuration objects through the
[Herodesk API](https://api.herodesk.io/v1/site/schema). It requires an API key
with permission to read and change the relevant objects.

## Example Usage

```terraform
terraform {
  required_providers {
    herodesk = {
      source = "plain-insure/herodesk"
    }
  }
}

provider "herodesk" {}
```

## Authentication

Set `HERODESK_TOKEN` to a Herodesk API key. The provider sends it as a Bearer
token on every request.

```terraform
provider "herodesk" {
  token = var.herodesk_token
}
```

For a custom API endpoint, such as a compatible test server, set
`HERODESK_BASE_URL` or configure `base_url`.

## Schema

### Optional

- `base_url` (String) API base URL. Defaults to `https://api.herodesk.io/v1` or
  `HERODESK_BASE_URL` when set.
- `token` (String, Sensitive) Herodesk API token. Defaults to `HERODESK_TOKEN`.