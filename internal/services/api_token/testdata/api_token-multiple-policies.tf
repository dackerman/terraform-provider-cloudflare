data "cloudflare_api_token_permission_groups_list" "dns_read" {
  name  = "DNS Read"
  scope = "com.cloudflare.api.account.zone"
}

data "cloudflare_api_token_permission_groups_list" "zone_read" {
  name  = "Zone Read"
  scope = "com.cloudflare.api.account.zone"
}

data "cloudflare_api_token_permission_groups_list" "analytics_read" {
  name  = "Analytics Read"
  scope = "com.cloudflare.api.account.zone"
}

resource "cloudflare_api_token" "test_multiple_policies" {
  name   = "%[1]s"
  status = "%[2]s"

  # First policy - allow DNS operations on all zones
  policies = [
    {
      effect = "allow"
      permission_groups = [{
        id = data.cloudflare_api_token_permission_groups_list.dns_read.result[0].id
      }]
      resources = {
        "com.cloudflare.api.account.zone.*" = "*"
      }
    },
    # Second policy - allow zone read on specific zones
    {
      effect = "allow"
      permission_groups = [{
        id = data.cloudflare_api_token_permission_groups_list.zone_read.result[0].id
      }]
      resources = {
        "com.cloudflare.api.account.zone.%[3]s" = "*"
      }
    },
    # Third policy - deny analytics on all zones (overrides allows)
    {
      effect = "%[4]s"
      permission_groups = [{
        id = data.cloudflare_api_token_permission_groups_list.analytics_read.result[0].id
      }]
      resources = {
        "com.cloudflare.api.account.zone.*" = "*"
      }
    }
  ]

  condition = {
    request_ip = {
      in = ["%[5]s"]
    }
  }
}