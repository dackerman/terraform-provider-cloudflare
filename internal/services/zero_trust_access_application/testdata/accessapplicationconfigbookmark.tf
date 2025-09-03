resource "cloudflare_zero_trust_access_application" "%[1]s" {
  %[3]s_id = "%[4]s"
  name     = "%[1]s"
  domain   = "https://example.com"
  type     = "bookmark"

  app_launcher_visible = true
  logo_url            = "https://www.cloudflare.com/img/logo-web-badges/cf-logo-on-white-bg.svg"
}