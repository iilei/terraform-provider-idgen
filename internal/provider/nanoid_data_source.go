package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/iilei/terraform-provider-idgen/internal/idgen"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &NanoIDDataSource{}

func NewNanoIDDataSource() datasource.DataSource {
	return &NanoIDDataSource{}
}

// NanoIDDataSource defines the data source implementation.
type NanoIDDataSource struct{}

// NanoIDDataSourceModel describes the data source data model.
type NanoIDDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	Length    types.Int64  `tfsdk:"length"`
	Alphabet  types.String `tfsdk:"alphabet"`
	GroupSize types.Int64  `tfsdk:"group_size"`
	Seed      types.String `tfsdk:"seed"`
}

func (d *NanoIDDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nanoid"
}

func (d *NanoIDDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Generates a NanoID identifier.\n\n" +
			"**Security Notice:** When using `seed`, IDs become deterministic and predictable. " +
			"Never use seeded IDs for security tokens, passwords, or cryptographic purposes.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The generated NanoID.",
				Computed:    true,
			},
			"length": schema.Int64Attribute{
				Description: "The length of the generated ID. Defaults to 21.",
				Optional:    true,
			},
			"alphabet": schema.StringAttribute{
				Description: "The alphabet to use for ID generation. Can be 'alphanumeric' (a-zA-Z0-9), 'numeric' (0-9), " +
					"'readable' (excludes 0/O, 1/l/I), 'less_confusable' (lowercase only, excludes 0/O, 1/l/I), " +
					"'least_confusable' (most distinct characters, excludes 0/O, 1/l/I, 2/Z, 5/S, 6/G, 8/B), " +
					"or a custom string of characters.",
				Optional: true,
			},
			"group_size": schema.Int64Attribute{
				Description: "Number of characters per group, separated by dashes. If not set, no grouping is applied.",
				Optional:    true,
			},
			"seed": schema.StringAttribute{
				MarkdownDescription: "Optional seed for deterministic ID generation. Behavior:\n\n" +
					"- **Integer** - parsed and used as random seed\n" +
					"- **Text string** - hashed deterministically and used as random seed\n" +
					"- **Omitted** - cryptographically random (different each apply)\n\n" +
					"**WARNING:** Seeded IDs are deterministic and should not be used for security tokens or secrets.",
				Optional: true,
			},
		},
	}
}

func (d *NanoIDDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Provider configuration is not needed for this implementation
}

func (d *NanoIDDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data NanoIDDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Set defaults
	length := int(21)
	if !data.Length.IsNull() {
		length = int(data.Length.ValueInt64())
	}

	// Validate length
	if !validateLength(int64(length), &resp.Diagnostics) {
		return
	}

	alphabet := idgen.Readable
	if !data.Alphabet.IsNull() {
		alphabetStr := data.Alphabet.ValueString()
		switch alphabetStr {
		case "alphanumeric":
			alphabet = idgen.Alphanumeric
		case "numeric":
			alphabet = idgen.Numeric
		case "readable":
			alphabet = idgen.Readable
		default:
			// Custom alphabet
			alphabet = alphabetStr
		}
	}

	// Warn if alphabet contains dashes and grouping is enabled
	if !data.GroupSize.IsNull() && data.GroupSize.ValueInt64() > 0 {
		if strings.Contains(alphabet, "-") {
			resp.Diagnostics.AddWarning(
				warningAlphabetContainsDashTitle,
				warningAlphabetContainsDashDetail,
			)
		}
	}

	// Check if seed is provided
	var seed *int64
	if !data.Seed.IsNull() {
		seedVal, _ := stringToSeed(data.Seed.ValueString())
		seed = &seedVal
	}

	// Determine group size for length calculation
	groupSize := 0
	if !data.GroupSize.IsNull() {
		groupSize = int(data.GroupSize.ValueInt64())
	}

	// Generate the NanoID (grouping is applied internally if groupSize > 0)
	id, err := idgen.GenerateNanoID(alphabet, length, seed, groupSize)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to generate NanoID",
			"Could not generate NanoID: "+err.Error(),
		)
		return
	}

	data.ID = types.StringValue(id)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
