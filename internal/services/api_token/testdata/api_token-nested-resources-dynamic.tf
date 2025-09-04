# Test nested resources using the dynamic attribute
resource "cloudflare_api_token" "test_nested_resources_dynamic" {
  name = "%[1]s"

  policies = [
    {
      effect = "allow"
      permission_groups = [{
        id = "c8fed203ed3043cba015a93ad1616f1f" # Zone Read
      }]
      # Nested resource structure - should work with dynamic attribute
      resources = {
        "com.cloudflare.api.account.%[2]s" = {
          "com.cloudflare.api.account.zone.*" = "*"
        }
      }
    }
  ]
}