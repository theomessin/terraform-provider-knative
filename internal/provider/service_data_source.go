package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	client "knative.dev/serving/pkg/client/clientset/versioned/typed/serving/v1"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ServiceDataSource{}

func NewServiceDataSource() datasource.DataSource {
	return &ServiceDataSource{}
}

// ServiceDataSource defines the data source implementation.
type ServiceDataSource struct {
	client *client.ServingV1Client
}

// ServiceDataSourceModel describes the data source data model.
type ServiceDataSourceModel struct {
	Namespace types.String `tfsdk:"namespace"`
	Name      types.String `tfsdk:"name"`
	Status    types.Object `tfsdk:"status"`
	Id        types.String `tfsdk:"id"`
}

func (d *ServiceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service"
}

func (d *ServiceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Knative Service data source",

		Attributes: map[string]schema.Attribute{
			"namespace": schema.StringAttribute{
				MarkdownDescription: "The namespace where the Knative Service resource is located. Defaults to `default`.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the Knative Service resource.",
				Required:            true,
			},
			"status": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"address": schema.ObjectAttribute{
						Computed:            true,
						MarkdownDescription: "Address holds the information needed for a Route to be the target of an event.",
						AttributeTypes: map[string]attr.Type{
							"url": types.StringType,
						},
					},
					"latest_created_revision_name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "LatestCreatedRevisionName is the last revision that was created from this Configuration. It might not be ready yet, for that use LatestReadyRevisionName.",
					},
					"latest_ready_revision_name": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "LatestReadyRevisionName holds the name of the latest Revision stamped out from this Configuration that has had its \"Ready\" condition become \"True\".",
					},
					"observed_generation": schema.Int64Attribute{
						Computed:            true,
						MarkdownDescription: "ObservedGeneration is the 'Generation' of the Service that was last processed by the controller.",
					},
					"url": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: "URL holds the url that will distribute traffic over the provided traffic targets. It generally has the form http[s]://{route-name}.{route-namespace}.{cluster-level-suffix}",
					},
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the Knative Service resource. In the form of `namespace/name`.",
				Computed:            true,
			},
		},
	}
}

func (d *ServiceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.ServingV1Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.ServingV1Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *ServiceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ServiceDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Namespace.IsNull() {
		data.Namespace = types.StringValue("default")
	}

	client := *d.client
	service, err := client.
		Services(data.Namespace.ValueString()).
		Get(ctx, data.Name.ValueString(), metav1.GetOptions{})

	if err != nil {
		resp.Diagnostics.AddError("Knative Client Error", err.Error())
		return
	}

	addressAttrType := map[string]attr.Type{"url": types.StringType}
	address, _ := types.ObjectValue(addressAttrType, map[string]attr.Value{
		"url": types.StringValue(service.Status.Address.URL.String()),
	})

	statusAttrType := map[string]attr.Type{
		"address":                      types.ObjectType{AttrTypes: addressAttrType},
		"latest_created_revision_name": types.StringType,
		"latest_ready_revision_name":   types.StringType,
		"observed_generation":          types.Int64Type,
		"url":                          types.StringType,
	}

	data.Status, _ = types.ObjectValue(statusAttrType, map[string]attr.Value{
		"address":                      address,
		"latest_created_revision_name": types.StringValue(service.Status.LatestCreatedRevisionName),
		"latest_ready_revision_name":   types.StringValue(service.Status.LatestReadyRevisionName),
		"observed_generation":          types.Int64Value(service.Status.ObservedGeneration),
		"url":                          types.StringValue(service.Status.URL.String()),
	})

	data.Id = types.StringValue(service.Namespace + "/" + service.Name)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
