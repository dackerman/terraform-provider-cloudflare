package api_token_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cloudflare/terraform-provider-cloudflare/internal/acctest"
	"github.com/cloudflare/terraform-provider-cloudflare/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccAPIToken_Basic(t *testing.T) {
	rnd := utils.GenerateRandomResourceName()
	resourceName := "cloudflare_api_token.test_account_token"

	var policyId string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudflareAPITokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.LoadTestCase("api_token-without-condition.tf", rnd),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rnd),
					resource.TestCheckResourceAttr(resourceName, "policies.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "policies.0.id"),
					resource.TestCheckResourceAttrWith(resourceName, "policies.0.id", func(value string) error {
						policyId = value
						return nil
					}),
					// conditions by default should not be set
					resource.TestCheckNoResourceAttr(resourceName, "condition.request_ip.0.in"),
					resource.TestCheckNoResourceAttr(resourceName, "condition.request_ip.0.not_in"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("name"), knownvalue.StringExact(rnd)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("policies"), knownvalue.ListSizeExact(1)),
				},
			},
			{
				Config: acctest.LoadTestCase("api_token-without-condition.tf", rnd),
				// re-plan should not detect drift
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"value"},
			},
			{
				Config: acctest.LoadTestCase("api_token-without-condition.tf", rnd+"-updated"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rnd+"-updated"),
					resource.TestCheckResourceAttr(resourceName, "policies.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "policies.0.id"),
					resource.TestCheckResourceAttrWith(resourceName, "policies.0.id", func(value string) error {
						if value != policyId {
							return fmt.Errorf("policy ID changed from %s to %s", policyId, value)
						}
						return nil
					}),
					// conditions still not set
					resource.TestCheckNoResourceAttr(resourceName, "condition.request_ip.0.in"),
					resource.TestCheckNoResourceAttr(resourceName, "condition.request_ip.0.not_in"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("name"), knownvalue.StringExact(rnd+"-updated")),
				},
			},
			{
				Config: acctest.LoadTestCase("api_token-without-condition.tf", rnd+"-updated"),
				// re-plan should not detect drift
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccAPIToken_SetIndividualCondition(t *testing.T) {
	rnd := utils.GenerateRandomResourceName()
	resourceName := "cloudflare_api_token.test_account_token"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudflareAPITokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.LoadTestCase("api_token-with-individual-condition.tf", rnd),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rnd),
					resource.TestCheckResourceAttr(resourceName, "condition.request_ip.in.0", "192.0.2.1/32"),
					resource.TestCheckNoResourceAttr(resourceName, "condition.request_ip.not_in"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("name"), knownvalue.StringExact(rnd)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("condition").AtMapKey("request_ip").AtMapKey("in"), knownvalue.ListSizeExact(1)),
				},
			},
			{
				Config: acctest.LoadTestCase("api_token-with-individual-condition.tf", rnd),
				// re-plan should not detect drift
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"value"},
			},
		},
	})
}

func TestAccAPIToken_SetAllCondition(t *testing.T) {
	rnd := utils.GenerateRandomResourceName()
	resourceName := "cloudflare_api_token.test_account_token"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudflareAPITokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.LoadTestCase("api_token-with-all-condition.tf", rnd),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rnd),
					resource.TestCheckResourceAttr(resourceName, "condition.request_ip.in.0", "192.0.2.1/32"),
					resource.TestCheckResourceAttr(resourceName, "condition.request_ip.not_in.0", "198.51.100.1/32"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("name"), knownvalue.StringExact(rnd)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("condition").AtMapKey("request_ip").AtMapKey("in"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("condition").AtMapKey("request_ip").AtMapKey("not_in"), knownvalue.ListSizeExact(1)),
				},
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"value"},
			},
		},
	})
}

func TestAccAPIToken_TokenTTL(t *testing.T) {
	rnd := utils.GenerateRandomResourceName()
	resourceName := "cloudflare_api_token.test_account_token"

	oneDaysFromNow := time.Now().UTC().AddDate(0, 0, 1)
	expireTime := oneDaysFromNow.Format(time.RFC3339)
	twoDaysFromNow := time.Now().UTC().AddDate(0, 0, 2)
	updatedExpireTime := twoDaysFromNow.Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudflareAPITokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.LoadTestCase("api_token-with-ttl.tf", rnd, expireTime),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rnd),
					resource.TestCheckResourceAttr(resourceName, "not_before", "2018-07-01T05:20:00Z"),
					resource.TestCheckResourceAttr(resourceName, "expires_on", expireTime),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("name"), knownvalue.StringExact(rnd)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("not_before"), knownvalue.StringExact("2018-07-01T05:20:00Z")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("expires_on"), knownvalue.StringExact(expireTime)),
				},
			},
			{
				Config: acctest.LoadTestCase("api_token-with-ttl.tf", rnd, expireTime),
				// re-plan should not detect drift
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"value"},
			},
			{
				Config: acctest.LoadTestCase("api_token-with-ttl.tf", rnd, updatedExpireTime),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rnd),
					resource.TestCheckResourceAttr(resourceName, "not_before", "2018-07-01T05:20:00Z"),
					resource.TestCheckResourceAttr(resourceName, "expires_on", updatedExpireTime),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("expires_on"), knownvalue.StringExact(updatedExpireTime)),
				},
			},
		},
	})
}

func TestAccAPIToken_PermissionGroupOrder(t *testing.T) {
	rnd := utils.GenerateRandomResourceName()
	name := "cloudflare_api_token." + rnd
	permissionID1 := "82e64a83756745bbbb1c9c2701bf816b" // DNS read
	permissionID2 := "e199d584e69344eba202452019deafe3" // Disable ESC read

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudflareAPITokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.LoadTestCase("api_token-permissiongroup-order.tf", rnd, permissionID1, permissionID2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rnd),
					resource.TestCheckResourceAttr(name, "policies.0.permission_groups.0.id", permissionID1),
					resource.TestCheckResourceAttr(name, "policies.0.permission_groups.1.id", permissionID2),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(name, tfjsonpath.New("name"), knownvalue.StringExact(rnd)),
				},
			},
			{
				Config: acctest.LoadTestCase("api_token-permissiongroup-order.tf", rnd, permissionID2, permissionID1),
				// changing the order of permission groups should not affect plan
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ResourceName:            name,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"value"},
			},
		},
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudflareAPITokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: acctest.LoadTestCase("api_token-permissiongroup-order.tf", rnd, permissionID2, permissionID1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "name", rnd),
					resource.TestCheckResourceAttr(name, "policies.0.permission_groups.0.id", permissionID1),
					resource.TestCheckResourceAttr(name, "policies.0.permission_groups.1.id", permissionID2),
				),
			},
			{
				Config: acctest.LoadTestCase("api_token-permissiongroup-order.tf", rnd, permissionID2, permissionID1),
				// re-applying same change does not produce drift
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: acctest.LoadTestCase("api_token-permissiongroup-order.tf", rnd, permissionID1, permissionID2),
				// changing the order of permission groups should not affect plan
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})

	// Dynamic permission retrieval test - commented out because it requires data sources
	// that need to be configured properly with actual permission IDs
	/*
	rnd2 := utils.GenerateRandomResourceName()
	permissionID0 := ""
	permissionID1Dynamic := ""

	var policyId string

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudflareAPITokenDestroy,
		Steps: []resource.TestStep{
			{
				// not setting permission IDs first, retrieving them from API by name
				Config: acctest.LoadTestCase("api_token-permissiongroup-order.tf", rnd2, "", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudflare_api_token.test_account_token", "name", rnd2),
					resource.TestCheckResourceAttr("cloudflare_api_token.test_account_token", "policies.#", "1"),
					resource.TestCheckResourceAttrSet("cloudflare_api_token.test_account_token", "policies.0.id"),
					resource.TestCheckResourceAttrWith("cloudflare_api_token.test_account_token", "policies.0.id", func(value string) error {
						policyId = value
						return nil
					}),
					resource.TestCheckResourceAttr("cloudflare_api_token.test_account_token", "policies.0.permission_groups.#", "2"),
					resource.TestCheckResourceAttrWith("cloudflare_api_token.test_account_token", "policies.0.permission_groups.0.id", func(value string) error {
						permissionID0 = value
						return nil
					}),
					resource.TestCheckResourceAttrWith("cloudflare_api_token.test_account_token", "policies.0.permission_groups.1.id", func(value string) error {
						permissionID1Dynamic = value
						return nil
					}),
				),
			},
			// below we try changing the order of the permission group IDs and
			// verify there are no plan changes
			{
				Config: acctest.LoadTestCase("api_token-permissiongroup-order.tf", rnd2, permissionID0, permissionID1Dynamic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: acctest.LoadTestCase("api_token-permissiongroup-order.tf", rnd2, permissionID1Dynamic, permissionID0),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// try updating the token and ensure policy information hasn't
			// changed
			{
				Config: acctest.LoadTestCase("api_token-permissiongroup-order.tf", rnd2+"updated", permissionID1Dynamic, permissionID0),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudflare_api_token.test_account_token", "name", rnd2+"updated"),
					resource.TestCheckResourceAttr("cloudflare_api_token.test_account_token", "policies.#", "1"),
					resource.TestCheckResourceAttrWith("cloudflare_api_token.test_account_token", "policies.0.id", func(value string) error {
						if value != policyId {
							return fmt.Errorf("policy ID changed from %s to %s", policyId, value)
						}
						return nil
					}),
					resource.TestCheckResourceAttr("cloudflare_api_token.test_account_token", "policies.0.permission_groups.#", "2"),
					resource.TestCheckResourceAttrWith("cloudflare_api_token.test_account_token", "policies.0.permission_groups.0.id", func(value string) error {
						if value != permissionID0 {
							return fmt.Errorf("permission ID 0 changed from %s to %s", permissionID0, value)
						}
						return nil
					}),
					resource.TestCheckResourceAttrWith("cloudflare_api_token.test_account_token", "policies.0.permission_groups.1.id", func(value string) error {
						if value != permissionID1Dynamic {
							return fmt.Errorf("permission ID 1 changed from %s to %s", permissionID1Dynamic, value)
						}
						return nil
					}),
				),
			},
			{
				Config: acctest.LoadTestCase("api_token-permissiongroup-order.tf", rnd2+"updated2", permissionID0, permissionID1Dynamic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("cloudflare_api_token.test_account_token", "name", rnd2+"updated2"),
					resource.TestCheckResourceAttr("cloudflare_api_token.test_account_token", "policies.#", "1"),
					resource.TestCheckResourceAttrWith("cloudflare_api_token.test_account_token", "policies.0.id", func(value string) error {
						if value != policyId {
							return fmt.Errorf("policy ID changed from %s to %s", policyId, value)
						}
						return nil
					}),
					resource.TestCheckResourceAttr("cloudflare_api_token.test_account_token", "policies.0.permission_groups.#", "2"),
					resource.TestCheckResourceAttrWith("cloudflare_api_token.test_account_token", "policies.0.permission_groups.0.id", func(value string) error {
						if value != permissionID0 {
							return fmt.Errorf("permission ID 0 changed from %s to %s", permissionID0, value)
						}
						return nil
					}),
					resource.TestCheckResourceAttrWith("cloudflare_api_token.test_account_token", "policies.0.permission_groups.1.id", func(value string) error {
						if value != permissionID1Dynamic {
							return fmt.Errorf("permission ID 1 changed from %s to %s", permissionID1Dynamic, value)
						}
						return nil
					}),
				),
			},
		},
	})
	*/
}

func TestAccAPIToken_MultiplePolicies(t *testing.T) {
	rnd := utils.GenerateRandomResourceName()
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	resourceName := "cloudflare_api_token.test_multiple_policies"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudflareAPITokenDestroy,
		Steps: []resource.TestStep{
			{
				// Create with multiple policies - mix of allow and deny
				Config: acctest.LoadTestCase("api_token-multiple-policies-simple.tf", rnd, "active", "allow", accountID, "192.0.2.0/24"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rnd),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "policies.#", "3"),
					// Check that we have the expected effects
					resource.TestCheckResourceAttr(resourceName, "policies.0.effect", "allow"),
					resource.TestCheckResourceAttr(resourceName, "policies.1.effect", "allow"),
					resource.TestCheckResourceAttr(resourceName, "policies.2.effect", "deny"),
					// Check permission groups are set
					resource.TestCheckResourceAttr(resourceName, "policies.0.permission_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policies.1.permission_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policies.2.permission_groups.#", "1"),
					// Check condition
					resource.TestCheckResourceAttr(resourceName, "condition.request_ip.in.0", "192.0.2.0/24"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("name"), knownvalue.StringExact(rnd)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("status"), knownvalue.StringExact("active")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("policies"), knownvalue.ListSizeExact(3)),
					// Verify each policy has an ID assigned
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("policies").AtSliceIndex(0).AtMapKey("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("policies").AtSliceIndex(1).AtMapKey("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("policies").AtSliceIndex(2).AtMapKey("id"), knownvalue.NotNull()),
					// Check computed fields are populated
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("issued_on"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("modified_on"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("value"), knownvalue.NotNull()),
				},
			},
			{
				// Update second policy to deny effect
				Config: acctest.LoadTestCase("api_token-multiple-policies-simple.tf", rnd, "active", "deny", accountID, "192.0.2.0/24"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rnd),
					resource.TestCheckResourceAttr(resourceName, "policies.#", "3"),
					// Second policy should now be deny
					resource.TestCheckResourceAttr(resourceName, "policies.1.effect", "deny"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("policies").AtSliceIndex(1).AtMapKey("effect"), knownvalue.StringExact("deny")),
				},
			},
			{
				// Test no drift on re-apply
				Config: acctest.LoadTestCase("api_token-multiple-policies-simple.tf", rnd, "active", "deny", accountID, "192.0.2.0/24"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				// Import test
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"value"},
			},
			{
				// Update status to disabled and change IP condition
				Config: acctest.LoadTestCase("api_token-multiple-policies-simple.tf", rnd, "disabled", "deny", accountID, "203.0.113.0/24"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "condition.request_ip.in.0", "203.0.113.0/24"),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("status"), knownvalue.StringExact("disabled")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("condition").AtMapKey("request_ip").AtMapKey("in").AtSliceIndex(0), knownvalue.StringExact("203.0.113.0/24")),
				},
			},
		},
	})
}

func testAccCheckCloudflareAPITokenDestroy(s *terraform.State) error {
	client := acctest.SharedClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudflare_api_token" {
			continue
		}

		// Try to get the token - it should not exist
		_, err := client.User.Tokens.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("api token still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}