
# 	resource "cloudflare_workers_script" "%[2]s" {
# 		account_id = "%[4]s"
# 		name       = "%[1]s"
# 		content    = "%[5]s"
# 	}

resource "cloudflare_workers_for_platforms_dispatch_namespace" "%[1]s" {
	account_id = "%[4]s"
	name       = "%[1]s"
}
	resource "cloudflare_workers_secret" "%[2]s" {
		account_id  = "%[4]s"
		script_name = "script_1"
		name 		= "%[2]s"
		dispatch_namespace = "%[1]s"
		type = "secret_text"
		depends_on = [cloudflare_workers_for_platforms_dispatch_namespace.%[1]s]
	}
