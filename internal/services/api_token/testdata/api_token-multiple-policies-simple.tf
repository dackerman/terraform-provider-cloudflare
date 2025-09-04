resource "cloudflare_api_token" "test_multiple_policies" {
  name   = "%[1]s"
  status = "%[2]s"

  # First policy - allow DNS operations on all zones
  policies = [
    {
      effect = "allow"
      permission_groups = [{
        id = "82e64a83756745bbbb1c9c2701bf816b" # DNS Read
      }]
      resources = {
        "com.cloudflare.api.account.zone.*" = "*"
      }
    },
    # Second policy - allow zone read on specific account
    {
      effect = "%[3]s"
      permission_groups = [{
        id = "c8fed203ed3043cba015a93ad1616f1f" # Zone Read
      }]
      resources = {
        "com.cloudflare.api.account.%[4]s" = "*"
      }
    },
    # Third policy - deny analytics on all zones
    {
      effect = "deny"
      permission_groups = [{
        id = "9c88f9c5bce24ce7af9a958ba9c504db" # Analytics Read
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