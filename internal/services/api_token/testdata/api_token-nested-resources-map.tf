# Test nested resources using the new resources_map attribute
resource "cloudflare_api_token" "test_nested_resources_map" {
  name = "%[1]s"

  policies = [
    {
      effect = "allow"
      permission_groups = [{
        id = "c8fed203ed3043cba015a93ad1616f1f" # Zone Read
      }]
      # Using resources_map for nested structure - all zones in specific account
      resources_map = {
        "com.cloudflare.api.account.%[2]s" = {
          "com.cloudflare.api.account.zone.*" = "*"
        }
      }
    }
  ]
}