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

// TestAccCloudflareRuleset_ModifyPlanBug tests the specific ModifyPlan bug where
// rule IDs and refs are not properly preserved between state and plan, causing
// phantom changes on every apply.
func TestAccCloudflareRuleset_ModifyPlanBug(t *testing.T) {
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
				Config: testAccCheckCloudflareRulesetWithRefs(rnd, zoneID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-ruleset-with-refs"),
					resource.TestCheckResourceAttr(resourceName, "kind", "zone"),
					resource.TestCheckResourceAttr(resourceName, "phase", "http_request_firewall_managed"),
					resource.TestCheckResourceAttr(resourceName, "rules.#", "1"),
					
					// Check that the rule has both ID and ref
					resource.TestCheckResourceAttrSet(resourceName, "rules.0.id"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.ref", "my-custom-ref"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.action", "execute"),
					resource.TestCheckResourceAttr(resourceName, "rules.0.action_parameters.id", "efb7b8c949ac4650a09736fc376e9aee"),
				),
			},
			{
				// Apply the same configuration again - should have no changes
				Config: testAccCheckCloudflareRulesetWithRefs(rnd, zoneID),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// This should pass if ModifyPlan is working correctly
						// but will fail with the current bug
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				// Make a small change to ensure updates work
				Config: testAccCheckCloudflareRulesetWithRefsUpdated(rnd, zoneID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "rules.0.description", "Updated description"),
				),
			},
			{
				// Apply again to verify no phantom changes after real update
				Config: testAccCheckCloudflareRulesetWithRefsUpdated(rnd, zoneID),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func testAccCheckCloudflareRulesetWithRefs(rnd, zoneID string) string {
	return fmt.Sprintf(`
resource "cloudflare_ruleset" "%[1]s" {
  zone_id     = "%[2]s"
  name        = "test-ruleset-with-refs"
  kind        = "zone"
  phase       = "http_request_firewall_managed"
  description = "Testing ModifyPlan with refs"

  rules = [{
    expression  = "true"
    action      = "execute"
    ref         = "my-custom-ref"
    description = "Test rule with ref"
    action_parameters = {
      id = "efb7b8c949ac4650a09736fc376e9aee"
      overrides = {
        enabled = true
      }
    }
  }]
}
`, rnd, zoneID)
}

func testAccCheckCloudflareRulesetWithRefsUpdated(rnd, zoneID string) string {
	return fmt.Sprintf(`
resource "cloudflare_ruleset" "%[1]s" {
  zone_id     = "%[2]s"
  name        = "test-ruleset-with-refs"
  kind        = "zone"
  phase       = "http_request_firewall_managed"
  description = "Testing ModifyPlan with refs"

  rules = [{
    expression  = "true"
    action      = "execute"
    ref         = "my-custom-ref"
    description = "Updated description"
    action_parameters = {
      id = "efb7b8c949ac4650a09736fc376e9aee"
      overrides = {
        enabled = true
      }
    }
  }]
}
`, rnd, zoneID)
}