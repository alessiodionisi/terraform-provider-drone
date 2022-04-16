---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "drone_secret Resource - terraform-provider-drone"
subcategory: ""
description: |-
  
---

# drone_secret (Resource)

Manage a repository secret.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String)
- `repository` (String)
- `value` (String, Sensitive)

### Optional

- `allow_on_pull_request` (Boolean)
- `id` (String) The ID of this resource.

