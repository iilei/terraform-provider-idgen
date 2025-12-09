package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &TemplatedDataSource{}

func NewTemplatedDataSource() datasource.DataSource {
	return &TemplatedDataSource{}
}

// TemplatedDataSource defines the data source implementation.
type TemplatedDataSource struct{}

// TemplatedDataSourceModel describes the data source data model.
type TemplatedDataSourceModel struct {
	ID types.String `tfsdk:"id"`
}

func (d *TemplatedDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_templated"
}

func (d *TemplatedDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Generates a templated identifier combining multiple ID types.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The generated templated ID.",
				Computed:    true,
			},
		},
	}
}

func (d *TemplatedDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Provider configuration is not needed for this simple implementation
}

func (d *TemplatedDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TemplatedDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// For now, just return a fixed string
	data.ID = types.StringValue("fixed-templated-thing-abc-123")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
