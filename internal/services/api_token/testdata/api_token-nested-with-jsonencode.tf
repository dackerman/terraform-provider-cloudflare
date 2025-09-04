# Test nested resources using jsonencode to pass complex structure as JSON string
resource "cloudflare_api_token" "test_nested_jsonencode" {
  name = "%[1]s"

  policies = [
    {
      effect = "allow"
      permission_groups = [{
        id = "c8fed203ed3043cba015a93ad1616f1f" # Zone Read
      }]
      # Using jsonencode to pass nested structure as a JSON string value
      resources = {
        "com.cloudflare.api.account.%[2]s" = jsonencode({
          "com.cloudflare.api.account.zone.*" = "*"
        })
      }
    }
  ]
}