package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/iilei/terraform-provider-idgen/internal/idgen"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ProquintDataSource{}

func NewProquintDataSource() datasource.DataSource {
	return &ProquintDataSource{}
}

// ProquintDataSource defines the data source implementation.
type ProquintDataSource struct{}

// ProquintDataSourceModel describes the data source data model.
type ProquintDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	Length    types.Int64  `tfsdk:"length"`
	GroupSize types.Int64  `tfsdk:"group_size"`
	Seed      types.Int64  `tfsdk:"seed"`
}

func (d *ProquintDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_proquint"
}

func (d *ProquintDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Generates a Proquint identifier.\n\n" +
			"**Security Notice:** When using `seed`, IDs become deterministic and predictable. " +
			"Never use seeded IDs for security tokens, passwords, or cryptographic purposes.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The generated Proquint.",
				Computed:    true,
			},
			"length": schema.Int64Attribute{
				Description: "The length of the generated ID in characters. Defaults to 11 (2 proquint words).",
				Optional:    true,
			},
			"group_size": schema.Int64Attribute{
				Description: "Number of characters per group, separated by dashes. If not set, no grouping is applied.",
				Optional:    true,
			},
			"seed": schema.Int64Attribute{
				Description: "Optional seed for deterministic ID generation. When provided, the same seed will " +
					"always produce the same ID. WARNING: Seeded IDs are predictable and should not be used for security.",
				Optional: true,
			},
		},
	}
}

func (d *ProquintDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Provider configuration is not needed for this implementation
}

func (d *ProquintDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProquintDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Set defaults
	// Default length is 11 characters (2 words: "lusab-babad")
	// Each proquint word is 5 chars, plus 1 separator = 11 chars for 2 words
	length := int64(11)
	if !data.Length.IsNull() {
		length = data.Length.ValueInt64()
	}

	// Convert character length to byte length
	// Proquint: 2 bytes = 1 word (5 chars), separator between words
	// Approximate: (length + 1) / 6 * 2 bytes
	byteLength := int((length + 1) / 6 * 2)
	if byteLength < 2 {
		byteLength = 2 // Minimum 1 word
	}

	// Check if seed is provided
	var seed *int64
	if !data.Seed.IsNull() {
		seedVal := data.Seed.ValueInt64()
		seed = &seedVal
	}

	// Generate the Proquint
	id, err := idgen.GenerateProquint(byteLength, seed)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to generate Proquint",
			"Could not generate Proquint: "+err.Error(),
		)
		return
	}

	// Apply grouping if group_size is specified
	if !data.GroupSize.IsNull() {
		groupSize := int(data.GroupSize.ValueInt64())
		if groupSize > 0 {
			id = applyGrouping(id, groupSize)
		}
	}

	data.ID = types.StringValue(id)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
