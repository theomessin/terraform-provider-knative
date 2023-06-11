package provider

import (
	"context"
	"os/user"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"k8s.io/client-go/tools/clientcmd"
	client "knative.dev/serving/pkg/client/clientset/versioned/typed/serving/v1"
)

// Ensure KnativeProvider satisfies various provider interfaces.
var _ provider.Provider = &KnativeProvider{}

// KnativeProvider defines the provider implementation.
type KnativeProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// KnativeProviderModel describes the provider data model.
type KnativeProviderModel struct {
	//
}

func (p *KnativeProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "knative"
	resp.Version = p.version
}

func (p *KnativeProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			//
		},
	}
}

func (p *KnativeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data KnativeProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	usr, _ := user.Current()
	dir := usr.HomeDir
	kubeconfigPath := filepath.Join(dir, ".kube", "config")
	config, _ := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	client, _ := client.NewForConfig(config)
	resp.DataSourceData = client
}

func (p *KnativeProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		//
	}
}

func (p *KnativeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewServiceDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &KnativeProvider{
			version: version,
		}
	}
}
