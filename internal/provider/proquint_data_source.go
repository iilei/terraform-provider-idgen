package provider

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/iilei/terraform-provider-idgen/internal/idgen"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ProquintDataSource{}

// ProquintDataSource defines the data source implementation.
type ProquintDataSource struct{}

// ProquintDataSourceModel describes the data source data model.
type ProquintDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	Seed      types.String `tfsdk:"seed"`
	Length    types.Int64  `tfsdk:"length"`
	GroupSize types.Int64  `tfsdk:"group_size"`
}

// NewProquintDataSource returns a new Proquint data source.
func NewProquintDataSource() datasource.DataSource {
	return &ProquintDataSource{}
}

// Metadata returns the data source type name.
func (d *ProquintDataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_proquint"
}

// Schema defines the data source attributes.
func (d *ProquintDataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Generates a random Proquint identifier with configurable length.\n\n" +
			"This data source generates proquint-formatted IDs with random content. " +
			"For canonical encoding of IPv4 addresses or uint32 integers, use `idgen_proquint_canonical` instead.\n\n" +
			"**Security Notice:** When using `seed`, IDs become deterministic and predictable. " +
			"Never use seeded IDs for security tokens, passwords, or cryptographic purposes.",
		Attributes: map[string]schema.Attribute{
			attrID: schema.StringAttribute{
				Description: "The generated Proquint.",
				Computed:    true,
			},
			attrLength: schema.Int64Attribute{
				MarkdownDescription: "The length of the generated ID in characters.\n\n" +
					"Common values: `11` (2 words, e.g., `lusab-babad`), `17` (3 words), `23` (4 words).",
				Required: true,
			},
			attrGroupSize: schema.Int64Attribute{
				Description: "Number of characters per group, separated by dashes. If not set, no grouping is applied.",
				Optional:    true,
			},
			attrSeed: schema.StringAttribute{
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

// Configure is a no-op: the data source needs no provider configuration.
func (d *ProquintDataSource) Configure(
	ctx context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
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

	// Convert character length to byte length (2 bytes per 5-letter word, one separator between words).
	byteLength := proquintByteLength(length)

	// Check if seed is provided
	// Seed parsing intentionally stays on the int64 path (idgen.StringToSeed) so that
	// existing configurations keep producing the same IDs as v1.0.x. Large values are
	// handled by idgen_proquint_big instead.
	var seed *big.Int
	var directEncode bool
	if !data.Seed.IsNull() {
		seedVal, shouldDirectEncode := stringToSeed(data.Seed.ValueString())
		seed = big.NewInt(seedVal)
		directEncode = shouldDirectEncode

		// Warn if using direct encoding with non-canonical length
		if shouldDirectEncode {
			// Direct encoding only applies to uint32-range seeds, which are canonically 11 characters.
			const canonicalLength = int64(11)

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

	// Determine group size (default to 5 for standard proquint format)
	groupSize := 5
	if !data.GroupSize.IsNull() {
		groupSize = int(data.GroupSize.ValueInt64())
	}

	// Remove all dashes and apply grouping
	if groupSize > 0 {
		id = strings.ReplaceAll(id, "-", "")
		id = idgen.ApplyGrouping(id, groupSize)
	}

	data.ID = types.StringValue(id)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
