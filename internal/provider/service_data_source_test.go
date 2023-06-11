package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test Data Source Read works with an existing service.
			{
				Config: exampleAppConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.knative_service.this", "status.address.url", "http://app.example.svc.cluster.local"),
					resource.TestCheckResourceAttr("data.knative_service.this", "status.latest_created_revision_name", "app-00001"),
					resource.TestCheckResourceAttr("data.knative_service.this", "status.latest_ready_revision_name", "app-00001"),
					resource.TestCheckResourceAttr("data.knative_service.this", "status.observed_generation", "1"),
					resource.TestCheckResourceAttr("data.knative_service.this", "status.url", "http://app.example.127.0.0.1.sslip.io"),
					resource.TestCheckResourceAttr("data.knative_service.this", "id", "example/app"),
				),
			},
			// Test Data Source Read errors when the service does not exist.
			{
				Config:      defaultAppConfig,
				ExpectError: regexp.MustCompile("services.serving.knative.dev \"app\" not found"),
			},
		},
	})
}

const exampleAppConfig = `
data "knative_service" "this" {
  namespace = "example"
  name = "app"
}
`

const defaultAppConfig = `
data "knative_service" "this" {
  name = "app"
}
`
