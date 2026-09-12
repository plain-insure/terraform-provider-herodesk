---
page_title: "herodesk_helpcenter Resource - Herodesk"
description: |-
  Manages a Herodesk Help Center.
---

# herodesk_helpcenter

Manages a Herodesk Help Center. Creating a help center also creates its root
folder. Deleting it removes its folders, articles, and FAQs.

`name` and `language` are required for creation. `language` is the default
language, such as `en-US`; `available_languages` lists additional translated
languages.

## Example Usage

```terraform
resource "herodesk_helpcenter" "support" {
  name                = "Support"
  language            = "en-US"
  description         = "Help for our customers."
  available_languages = ["da-DK"]
}
```

## Schema

### Required

- `name` (String) Name of the help center.
- `language` (String) Default language of the help center.

### Optional

- `available_languages` (List of String) Additional translated languages.
- `custom_domain` (String) Custom domain for the help center.
- `description` (String) Description shown on the help center front page.

### Read-Only

- `created_at` (String) UTC creation timestamp.
- `default_domain` (String) Default Herodesk subdomain.
- `public_url` (String) Public help center URL.
- `root_folder_id` (Number) Automatically created root folder ID.
- `updated_at` (String) UTC update timestamp.

## Import

```shell
terraform import herodesk_helpcenter.support 123
```