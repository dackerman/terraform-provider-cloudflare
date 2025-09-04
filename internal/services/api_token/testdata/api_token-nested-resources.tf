# This test attempts to create an API token with nested resources
# According to issue #5733, the API supports this but the provider doesn't

resource "cloudflare_api_token" "test_nested_resources" {
  name = "%[1]s"

  policies = [
    {
      effect = "allow"
      permission_groups = [{
        id = "c8fed203ed3043cba015a93ad1616f1f" # Zone Read
      }]
      # Try nested resource structure - this should fail if the bug exists
      resources = {
        "com.cloudflare.api.account.%[2]s" = jsonencode({
          "com.cloudflare.api.account.zone.*" = "*"
        })
      }
    }
  ]
}