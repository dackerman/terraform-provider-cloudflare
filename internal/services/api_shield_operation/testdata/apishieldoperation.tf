
	resource "cloudflare_api_shield_operation" "%[1]s" {
		zone_id = "%[2]s"
		operation_id = "494da477-9db2-427d-ba4b-beb454c7e698"
		state = "ignored"
	}
