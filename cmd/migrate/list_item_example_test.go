package main

import "testing"

func TestListItemTransformation(t *testing.T) {
	tests := []TestCase{
		{
			Name: "",
			Config: `
resource "cloudflare_list" "known_openai_ips" {
  account_id  = local.account_id
  name        = "known_openai_ips"
  description = "The set of IPs that internal OpenAI requests may come from. Used for exemption from firewall and bot detection rules. Do not use for general access control!"
  kind        = "ip"

  dynamic "item" {
    for_each = local.ips_data.ips
    iterator = ip_entry
    content {
      comment = ip_entry.value.comment

      value {
        ip = ip_entry.value.ip
      }
    }
  }
}
`,
			Expected: []string{`
resource "cloudflare_list" "known_openai_ips" {
  account_id  = local.account_id
  description = "The set of IPs that internal OpenAI requests may come from. Used for exemption from firewall and bot detection rules. Do not use for general access control!"
  kind        = "ip"
  name        = "known_openai_ips"
}
`, `
resource "cloudflare_list_item" "known_openai_ips" {
  account_id = cloudflare_list_item.known_openai_ips.account_id
  comment    = each.value.content.comment
  for_each   = [for ip_entry in [for value in local.ips_data.ips : { key = value, value = value }] : { content = { comment = ip_entry.value.comment, value = { ip = ip_entry.value.ip } } }]
  ip         = each.value.content.ip.value
  list_id    = cloudflare_list_item.known_openai_ips.id
}
`,
			},
		},
		{

			Name: "",
			Config: `
resource "cloudflare_list" "persona_webhook_ips" {
  account_id  = local.account_id
  name        = "persona_webhook_ips"
  description = "Static IPs used by Persona to send webhook notifications"
  kind        = "ip"

  dynamic "item" {
    for_each = local.persona_webhook_ips
    content {
      value {
        ip = item.value
      }
    }
  }
}
`,
			Expected: []string{`
resource "cloudflare_list" "persona_webhook_ips" {
  account_id  = local.account_id
  description = "Static IPs used by Persona to send webhook notifications"
  kind        = "ip"
  name        = "persona_webhook_ips"
}
`, `
resource "cloudflare_list_item" "persona_webhook_ips" {
  account_id = cloudflare_list_item.persona_webhook_ips.account_id
  comment    = each.value.content.comment
  for_each   = [for item in [for value in local.persona_webhook_ips : { key = value, value = value }] : { content = { value = { ip = item.value } } }]
  ip         = each.value.content.ip.value
  list_id    = cloudflare_list_item.persona_webhook_ips.id
}
`,
			},
		},
	}

	RunTransformationTests(t, tests, transformListItem)
}
