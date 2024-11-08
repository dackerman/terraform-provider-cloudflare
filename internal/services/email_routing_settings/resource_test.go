package email_routing_settings_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/cloudflare/terraform-provider-cloudflare/internal/acctest"
	"github.com/cloudflare/terraform-provider-cloudflare/internal/consts"
	"github.com/cloudflare/terraform-provider-cloudflare/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testEmailRoutingSettingsConfig(resourceID, zoneID string, subDomain string) string {
	return acctest.LoadTestCase("emailroutingsettingsconfig.tf", resourceID, zoneID, subDomain)
}

func TestAccTestEmailRoutingSettings(t *testing.T) {
	rnd := utils.GenerateRandomResourceName()
	name := "cloudflare_email_routing_settings." + rnd
	zoneID := os.Getenv("CLOUDFLARE_ZONE_ID")
	domain := os.Getenv("CLOUDFLARE_DOMAIN")
	subDomain := fmt.Sprintf("%s.%s", rnd, domain)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testEmailRoutingSettingsConfig(rnd, zoneID, subDomain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "enabled", "true"),
					resource.TestCheckResourceAttr(name, consts.ZoneIDSchemaKey, zoneID),
				),
			},
		},
	})
}
