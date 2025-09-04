resource "cloudflare_api_token" "test_two_policies" {
  name = "%[1]s"

  policies = [
    {
      effect = "allow"
      permission_groups = [{
        id = "82e64a83756745bbbb1c9c2701bf816b" # DNS Read - fixed ID
      }]
      resources = {
        "com.cloudflare.api.account.zone.*" = "*"
      }
    },
    {
      effect = "allow" 
      permission_groups = [{
        id = "c8fed203ed3043cba015a93ad1616f1f" # Zone Read - fixed ID
      }]
      resources = {
        "com.cloudflare.api.account.*" = "*"
      }
    }
  ]
}