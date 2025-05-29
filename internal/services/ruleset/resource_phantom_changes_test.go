package ruleset_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/cloudflare/terraform-provider-cloudflare/internal/acctest"
	"github.com/cloudflare/terraform-provider-cloudflare/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

// TestAccCloudflareRuleset_PhantomChanges tests the bug where the ruleset resource
// shows constant changes despite no fields being updated. This is caused by the
// ModifyPlan function incorrectly extracting state elements twice instead of
// extracting plan elements.
func TestAccCloudflareRuleset_PhantomChanges(t *testing.T) {
	// Temporarily unset CLOUDFLARE_API_TOKEN if it is set as the WAF
	// service does not yet support the API tokens and it results in
	// misleading state error messages.
	if os.Getenv("CLOUDFLARE_API_TOKEN") != "" {
		t.Setenv("CLOUDFLARE_API_TOKEN", "")
	}

	rnd := utils.GenerateRandomResourceName()
	zoneID := os.Getenv("CLOUDFLARE_ZONE_ID")
	resourceName := "cloudflare_ruleset." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCloudflareRulesetPhantomChanges(rnd, zoneID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "default"),
					resource.TestCheckResourceAttr(resourceName, "kind", "zone"),
					resource.TestCheckResourceAttr(resourceName, "phase", "http_request_firewall_managed"),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "2"),
					
					// Check first rule (OWASP)
					resource.TestCheckResourceAttr(resourceName, "rules.0.action", "execute"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.expression", "true"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.action_parameters.id", "4814384a9e5d4991b9815dcfc25d2f1f"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.action_parameters.overrides.categories.#", "3"),
					
					// Check second rule (Cloudflare Managed Ruleset)
					resource.TestCheckResourceAttr(resourceName, "rules.1.action", "execute"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.expression", "true"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.action_parameters.id", "efb7b8c949ac4650a09736fc376e9aee"),
					resource.TestCheckResourceAttr(resourceName, "rules.1.action_parameters.overrides.categories.#", "3"),
				),
			},
			{
				// This is the critical test - applying the same configuration again
				// should result in no changes, but due to the bug it will show changes
				Config: testAccCheckCloudflareRulesetPhantomChanges(rnd, zoneID),
				// This plan check should pass in a correctly functioning provider
				// but will fail with the current bug, demonstrating the phantom changes
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func testAccCheckCloudflareRulesetPhantomChanges(rnd, zoneID string) string {
	return fmt.Sprintf(`
resource "cloudflare_ruleset" "%[1]s" {
  zone_id     = "%[2]s"
  name        = "default"
  kind        = "zone"
  phase       = "http_request_firewall_managed"
  description = "Test ruleset for phantom changes bug"

  rules = [
    {
      expression = "true"
      action     = "execute"
      action_parameters = {
        id  = "4814384a9e5d4991b9815dcfc25d2f1f" # OWASP Rules
        ref = "4814384a9e5d4991b9815dcfc25d2f1f"
        overrides = {
          categories = [
            {
              category = "paranoia-level-2"
              enabled  = false
            },
            {
              category = "paranoia-level-3"
              enabled  = false
            },
            {
              category = "paranoia-level-4"
              enabled  = false
            }
          ]
        }
      }
    },
    {
      expression = "true"
      action     = "execute"
      action_parameters = {
        id  = "efb7b8c949ac4650a09736fc376e9aee" # Cloudflare Managed Ruleset
        ref = "efb7b8c949ac4650a09736fc376e9aee"
        overrides = {
          categories = [
            {
              category = "wordpress"
              enabled  = false
            },
            {
              category = "drupal"
              enabled  = false
            },
            {
              category = "joomla"
              enabled  = false
            }
          ]
        }
      }
    }
  ]
}
`, rnd, zoneID)
}