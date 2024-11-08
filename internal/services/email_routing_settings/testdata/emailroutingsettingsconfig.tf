
resource "cloudflare_email_routing_settings" "%[1]s" {
	zone_id = "%[2]s"
	name = "%[3]s"
}
		