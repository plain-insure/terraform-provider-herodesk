---
page_title: "herodesk_helpcenter Resource - Herodesk"
description: |-
  Manages a Herodesk Help Center.
---

# herodesk_helpcenter

Manages a Herodesk Help Center. Creating a help center also creates its root
folder. Deleting it removes its folders, articles, and FAQs.

The `json` payload maps to the Herodesk Help Centers API. `name` and `language`
are required for creation. `language` is the default language, such as `en-US`;
`available_languages` lists additional translated languages.

## Example Usage

```terraform
resource "herodesk_helpcenter" "support" {
  json = jsonencode({
    name                = "Support"
    language            = "en-US"
    description         = "Help for our customers."
    available_languages = ["da-DK"]
  })
}
```

## Schema

### Required

- `json` (String) JSON request body. Creation supports `name`, `language`,
  `description`, and `available_languages`. Updates may also set `custom_domain`.

### Read-Only

- `id` (String) Herodesk help center ID.

## Import

```shell
terraform import herodesk_helpcenter.support 123
```