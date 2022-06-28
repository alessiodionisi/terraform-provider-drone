package provider

import (
	"fmt"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mavimo/terraform-provider-drone/internal/provider/utils"
)

func TestAccResourceCron(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccResourceCronCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCron(testDroneUser, "hook-test", "cron_job"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"drone_cron.cron", "repository", fmt.Sprintf("%s/hook-test", testDroneUser),
					),
					resource.TestCheckResourceAttr(
						"drone_cron.cron", "name", "cron_job",
					),
				),
			},
		},
	})
}

func testAccResourceCron(user, repo, name string) string {
	return fmt.Sprintf(`
resource "drone_repo" "repo" {
	repository = "%s/%s"
}

resource "drone_cron" "cron" {
	repository = "${drone_repo.repo.repository}"
	name = "%s"
	expr = "@monthly"
	event = "push"
}
`, user, repo, name)
}

func testAccResourceCronCheckDestroy(state *terraform.State) error {
	provider, err := providerFactories["drone"]()
	if err != nil {
		return err
	}

	client := provider.Meta().(drone.Client)

	for _, resource := range state.RootModule().Resources {
		if resource.Type != "drone_cron" {
			continue
		}

		owner, repo, err := utils.ParseRepo(resource.Primary.Attributes["repository"])

		if err != nil {
			return err
		}

		err = client.CronDelete(owner, repo, resource.Primary.Attributes["name"])

		if err == nil {
			return fmt.Errorf(
				"Cron job still exists: %s/%s:%s",
				owner,
				repo,
				resource.Primary.Attributes["name"],
			)
		}
	}

	return nil
}
