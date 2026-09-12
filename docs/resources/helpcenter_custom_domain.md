---
page_title: "herodesk_helpcenter_custom_domain Resource - Herodesk"
description: |-
  Manages the custom domain of an existing Herodesk Help Center.
---

# herodesk_helpcenter_custom_domain

Manages only the `custom_domain` property of an existing Herodesk Help Center.
It does not create, modify, or delete the Help Center itself. Removing this
resource clears the configured custom domain.

## Example Usage

```terraform
data "herodesk_helpcenters" "support" {
  id = 123
}

resource "herodesk_helpcenter_custom_domain" "support" {
  helpcenter_id = data.herodesk_helpcenters.support.items[0].id
  custom_domain = "help.example.com"
}
```

## Schema

### Required

- `helpcenter_id` (Number) ID of the existing Help Center. Changing this value
  replaces the resource.
- `custom_domain` (String) Custom domain served by the Help Center.

### Read-Only

- `cname_domain` (String) CNAME domain derived from the public help center URL.

## Import

Import using the Help Center ID:

```shell
terraform import herodesk_helpcenter_custom_domain.support 123
```