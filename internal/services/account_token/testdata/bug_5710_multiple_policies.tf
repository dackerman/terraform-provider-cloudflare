# Test case to reproduce issue #5710
# Multiple policies with different resource scopes

data "cloudflare_account_api_token_permission_groups_list" "workers_r2_storage" {
  account_id = "%[2]s"
  name       = "Workers R2 Storage Write"
  scope      = "com.cloudflare.api.account"
}

data "cloudflare_account_api_token_permission_groups_list" "dns_write" {
  account_id = "%[2]s"
  name       = "DNS Write"
  scope      = "com.cloudflare.api.account.zone"
}

data "cloudflare_account_api_token_permission_groups_list" "zone_settings" {
  account_id = "%[2]s"
  name       = "Zone Settings Write"
  scope      = "com.cloudflare.api.account.zone"
}

resource "cloudflare_account_token" "bug_5710_test" {
  name       = "%[1]s"
  account_id = "%[2]s"

  # First policy: Account-level permissions
  policies = [
    {
      effect = "allow"
      permission_groups = [
        { id = data.cloudflare_account_api_token_permission_groups_list.workers_r2_storage.result[0].id }
      ]
      resources = {
        "com.cloudflare.api.account.%[2]s" = "*"
      }
    },
    # Second policy: Zone-level permissions with different resource scope
    {
      effect = "allow"
      permission_groups = [
        { id = data.cloudflare_account_api_token_permission_groups_list.dns_write.result[0].id },
        { id = data.cloudflare_account_api_token_permission_groups_list.zone_settings.result[0].id }
      ]
      resources = {
        "com.cloudflare.api.account.zone.%[3]s" = "*"
      }
    }
  ]
}