package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccServiceDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.knative_service.this", "status.address.url", "http://app.default.svc.cluster.local"),
					resource.TestCheckResourceAttr("data.knative_service.this", "status.latest_created_revision_name", "app-00001"),
					resource.TestCheckResourceAttr("data.knative_service.this", "status.latest_ready_revision_name", "app-00001"),
					resource.TestCheckResourceAttr("data.knative_service.this", "status.observed_generation", "1"),
					resource.TestCheckResourceAttr("data.knative_service.this", "status.url", "http://app.default.127.0.0.1.sslip.io"),
					resource.TestCheckResourceAttr("data.knative_service.this", "id", "default/app"),
				),
			},
		},
	})
}

const testAccServiceDataSourceConfig = `
data "knative_service" "this" {
  name = "app"
}
`
