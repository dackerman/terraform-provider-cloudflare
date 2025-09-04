package account_token_test

import (
	"os"
	"testing"

	"github.com/cloudflare/terraform-provider-cloudflare/internal/acctest"
	"github.com/cloudflare/terraform-provider-cloudflare/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

// TestAccAccountToken_Bug5710_MultiplePolicies tests the issue reported in
// https://github.com/cloudflare/terraform-provider-cloudflare/issues/5710
// where multiple policies with different resource scopes cause 
// "resources does not correlate with any element in actual" error
func TestAccAccountToken_Bug5710_MultiplePolicies(t *testing.T) {
	rnd := utils.GenerateRandomResourceName()
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	zoneID := os.Getenv("CLOUDFLARE_ZONE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: acctest.LoadTestCase("bug_5710_multiple_policies.tf", rnd, accountID, zoneID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudflare_account_token.bug_5710_test", "name", rnd),
					resource.TestCheckResourceAttr("cloudflare_account_token.bug_5710_test", "account_id", accountID),
					// Check first policy (account-level)
					resource.TestCheckResourceAttr("cloudflare_account_token.bug_5710_test", "policies.0.effect", "allow"),
					resource.TestCheckResourceAttrSet("cloudflare_account_token.bug_5710_test", "policies.0.permission_groups.0.id"),
					resource.TestCheckResourceAttr("cloudflare_account_token.bug_5710_test", "policies.0.resources.com.cloudflare.api.account."+accountID, "*"),
					// Check second policy (zone-level)
					resource.TestCheckResourceAttr("cloudflare_account_token.bug_5710_test", "policies.1.effect", "allow"),
					resource.TestCheckResourceAttrSet("cloudflare_account_token.bug_5710_test", "policies.1.permission_groups.0.id"),
					resource.TestCheckResourceAttr("cloudflare_account_token.bug_5710_test", "policies.1.resources.com.cloudflare.api.account.zone."+zoneID, "*"),
				),
			},
			{
				// Re-apply to check for drift
				Config: acctest.LoadTestCase("bug_5710_multiple_policies.tf", rnd, accountID, zoneID),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				// Update the name to trigger an update
				Config: acctest.LoadTestCase("bug_5710_multiple_policies.tf", rnd+"-updated", accountID, zoneID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudflare_account_token.bug_5710_test", "name", rnd+"-updated"),
					// Verify policies remain intact after update
					resource.TestCheckResourceAttr("cloudflare_account_token.bug_5710_test", "policies.0.resources.com.cloudflare.api.account."+accountID, "*"),
					resource.TestCheckResourceAttr("cloudflare_account_token.bug_5710_test", "policies.1.resources.com.cloudflare.api.account.zone."+zoneID, "*"),
				),
			},
			{
				// Re-apply again to check for drift after update
				Config: acctest.LoadTestCase("bug_5710_multiple_policies.tf", rnd+"-updated", accountID, zoneID),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

// TestAccAccountToken_Bug5710_UpdatePolicyResources tests modifying resources in existing policies
func TestAccAccountToken_Bug5710_UpdatePolicyResources(t *testing.T) {
	rnd := utils.GenerateRandomResourceName()
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	zoneID := os.Getenv("CLOUDFLARE_ZONE_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Start with a single policy
				Config: acctest.LoadTestCase("account_token-without-condition.tf", rnd, accountID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudflare_account_token.test_account_token", "name", rnd),
					resource.TestCheckResourceAttr("cloudflare_account_token.test_account_token", "policies.#", "1"),
				),
			},
			{
				// Update to multiple policies with different resource scopes
				Config: acctest.LoadTestCase("bug_5710_multiple_policies.tf", rnd, accountID, zoneID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudflare_account_token.bug_5710_test", "name", rnd),
					resource.TestCheckResourceAttr("cloudflare_account_token.bug_5710_test", "policies.#", "2"),
					// Verify both policies are present
					resource.TestCheckResourceAttr("cloudflare_account_token.bug_5710_test", "policies.0.resources.com.cloudflare.api.account."+accountID, "*"),
					resource.TestCheckResourceAttr("cloudflare_account_token.bug_5710_test", "policies.1.resources.com.cloudflare.api.account.zone."+zoneID, "*"),
				),
			},
		},
	})
}