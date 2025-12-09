package provider

import (
	"context"
	"fmt"

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
	Seed      types.String `tfsdk:"seed"`
}

func (d *ProquintDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_proquint"
}

func (d *ProquintDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Generates a random Proquint identifier with configurable length.\n\n" +
			"This data source generates proquint-formatted IDs with random content. " +
			"For canonical encoding of IPv4 addresses or uint32 integers, use `idgen_proquint_canonical` instead.\n\n" +
			"**Security Notice:** When using `seed`, IDs become deterministic and predictable. " +
			"Never use seeded IDs for security tokens, passwords, or cryptographic purposes.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The generated Proquint.",
				Computed:    true,
			},
			"length": schema.Int64Attribute{
				MarkdownDescription: "The length of the generated ID in characters.\n\n" +
					"Common values: `11` (2 words, e.g., `lusab-babad`), `17` (3 words), `23` (4 words).",
				Required: true,
			},
			"group_size": schema.Int64Attribute{
				Description: "Number of characters per group, separated by dashes. If not set, no grouping is applied.",
				Optional:    true,
			},
			"seed": schema.StringAttribute{
				MarkdownDescription: "Optional seed for deterministic random generation. Accepts any string value:\n\n" +
					"- **Text string** (e.g., `my-app-seed-42`) - hashed deterministically and used as random seed\n" +
					"- **Integer string** (e.g., `12345`) - parsed and used as random seed\n" +
					"- **Omitted** - uses cryptographically secure random generation (different each apply)\n\n" +
					"**Note:** For canonical encoding of IPv4 addresses or uint32/uint64 integers, use `idgen_proquint_canonical` instead.\n\n" +
					"**WARNING:** Seeded IDs are deterministic and should not be used for security tokens or secrets.",
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

	// Length is now required (no default)
	length := data.Length.ValueInt64()

	// Validate length
	if !validateLength(length, &resp.Diagnostics) {
		return
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
	var directEncode bool
	if !data.Seed.IsNull() {
		seedVal, shouldDirectEncode := stringToSeed(data.Seed.ValueString())
		seed = &seedVal
		directEncode = shouldDirectEncode

		// Warn if using direct encoding with non-canonical length
		if shouldDirectEncode {
			// Determine what the canonical length would be
			canonicalLength := int64(11) // default for uint32
			if seedVal > 0xFFFFFFFF {
				canonicalLength = 23 // uint64 range
			}

			if length != canonicalLength {
				resp.Diagnostics.AddWarning(
					"Non-Canonical Length for Direct Encoding",
					fmt.Sprintf(
						"The seed value '%s' will be canonically encoded to %d characters, but length=%d was requested. "+
							"The output will be %s to match your requested length. "+
							"Consider using idgen_proquint_canonical for canonical encoding without specifying length, "+
							"or adjust length to %d for the standard canonical output.",
						data.Seed.ValueString(),
						canonicalLength,
						length,
						map[bool]string{true: "truncated", false: "zero-padded"}[length < canonicalLength],
						canonicalLength,
					),
				)
			}
		}
	}

	// Generate the Proquint
	id, err := idgen.GenerateProquint(byteLength, seed, directEncode)
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
