# Terraform Provider for Herodesk

Terraform provider for managing Help Centers, webhooks, tags, teams, and SLA
policies in [Herodesk](https://herodesk.io/).

## Requirements

- Terraform 1.0 or later
- Go 1.24 or later to build the provider
- A Herodesk API key

## Using the provider

Configure the provider with `HERODESK_TOKEN`, or supply `token` directly. The
token is sent to Herodesk using Bearer authentication.

```hcl
terraform {
	required_providers {
		herodesk = {
			source = "plain-insure/herodesk"
		}
	}
}

provider "herodesk" {}
```

```powershell
$env:HERODESK_TOKEN = "your-api-key"
terraform init
terraform apply
```

The provider defaults to `https://api.herodesk.io/v1`. Set
`HERODESK_BASE_URL`, or set the provider's `base_url` argument, when using a
compatible endpoint.

## JSON object model

Herodesk objects are represented by a required `json` argument so the provider
can support the complete request payload defined by the Herodesk API. The API
response is stored back in `json` after creation and refresh. Values returned by
the API, including generated IDs and timestamps, are therefore visible in the
resource's `json` value.

For example, this creates a tag:

```hcl
resource "herodesk_tag" "priority" {
	json = jsonencode({
		name     = "Priority"
		archived = 0
	})
}
```

Collection data sources accept an optional `id`. They return the normalized API
response in `raw_json` and an `items` list containing one normalized JSON object
per returned item.

```hcl
data "herodesk_tags" "all" {}

locals {
	tags = [for item in data.herodesk_tags.all.items : jsondecode(item)]
}
```

## Resources

- [`herodesk_helpcenter`](docs/resources/helpcenter.md): create and manage help centers.
- [`herodesk_webhook`](docs/resources/webhook.md): create and manage event webhooks.
- [`herodesk_tag`](docs/resources/tag.md): create and manage conversation tags.
- [`herodesk_team`](docs/resources/team.md): create and manage teams, members, and inbox assignments.
- [`herodesk_sla_policy`](docs/resources/sla_policy.md): create and manage SLA policies.

## Data sources

- [`herodesk_helpcenters`](docs/data-sources/helpcenters.md): list help centers or get one by ID.
- [`herodesk_webhooks`](docs/data-sources/webhooks.md): list webhooks or get one by ID.
- [`herodesk_tags`](docs/data-sources/tags.md): list tags or get one by ID.
- [`herodesk_teams`](docs/data-sources/teams.md): list teams or get one by ID.
- [`herodesk_sla_policies`](docs/data-sources/sla_policies.md): list SLA policies or get one by ID.

## Releases

Pushing a `v*` tag creates a signed GitHub release through GoReleaser. Before
the first release, configure these repository Actions secrets:

- `GPG_PRIVATE_KEY`: ASCII-armored private key used to sign checksums.
- `GPG_PASSPHRASE`: passphrase for that key.
