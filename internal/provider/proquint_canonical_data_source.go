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
var _ datasource.DataSource = &ProquintCanonicalDataSource{}

func NewProquintCanonicalDataSource() datasource.DataSource {
	return &ProquintCanonicalDataSource{}
}

// ProquintCanonicalDataSource defines the data source implementation.
type ProquintCanonicalDataSource struct{}

// ProquintCanonicalDataSourceModel describes the data source data model.
type ProquintCanonicalDataSourceModel struct {
	ID    types.String `tfsdk:"id"`
	Value types.String `tfsdk:"value"`
}

func (d *ProquintCanonicalDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_proquint_canonical"
}

func (d *ProquintCanonicalDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Generates a canonical Proquint identifier from an IPv4 address or unsigned integer.\n\n" +
			"This data source implements the canonical proquint encoding as described in the " +
			"[original specification](https://arxiv.org/html/0901.4016). " +
			"It directly encodes the provided value as a proquint.\n\n" +
			"**Output Length:**\n\n" +
			"- **32-bit values** (IPv4, uint32 0-4294967295): 11 characters (2 proquint words)\n" +
			"- **64-bit values** (uint64 4294967296+): 23 characters (4 proquint words)\n\n" +
			"**Limitations:**\n\n" +
			"- **IPv6 not supported**: The original proquint specification focuses on 32-bit values. " +
			"IPv6 addresses (128-bit) must be manually converted to their integer representation before encoding.\n\n" +
			"**Use Cases:**\n\n" +
			"- Convert IP addresses to memorable identifiers\n" +
			"- Encode integer values as human-readable proquints\n" +
			"- Generate deterministic identifiers from numeric data\n\n" +
			"**Security Notice:** Canonical proquints are deterministic encodings of the input value. " +
			"They should not be used for security tokens or secrets.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The generated canonical Proquint identifier. Length varies by input:\n\n" +
					"- 11 characters for 32-bit values (IPv4, uint32)\n" +
					"- 23 characters for 64-bit values (uint64)",
				Computed: true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "The value to encode as a proquint. Accepts:\n\n" +
					"- **IPv4 address** (e.g., `127.0.0.1`)~>11 chars\n" +
					"- **Hexadecimal string** (e.g., `0x7f000001` or `7f000001`)~>11 or 23 chars\n" +
					"- **uint32 integer** (0-4294967295)~>11 chars\n" +
					"- **uint64 integer** (4294967296-18446744073709551615)~>23 chars\n\n" +
					"**Note:** IPv6 addresses are not supported directly. The canonical encoding follows the original " +
					"specification which focuses on 32-bit values. To encode IPv6 (128-bit), convert it to its " +
					"integer representation first.\n\n" +
					"Examples:\n" +
					"- `127.0.0.1`~>`lusab-babad` (11 chars)\n" +
					"- `0x7f000001`~>`lusab-babad` (11 chars, hex format)\n" +
					"- `2130706433`~>`lusab-babad` (11 chars, decimal)\n" +
					"- `0xffffffff`~>`zuzuz-zuzuz` (11 chars, max uint32)\n" +
					"- `0x7fffffffffffffff`~>`luzuz-zuzuz-zuzuz-zuzuz` (23 chars, max int64)",
				Required: true,
			},
		},
	}
}

func (d *ProquintCanonicalDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Provider configuration is not needed for this implementation
}

func (d *ProquintCanonicalDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ProquintCanonicalDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Parse the value and check if it's valid for canonical encoding
	value, _, errMsg := stringToCanonicalValue(data.Value.ValueString())

	if errMsg != "" {
		resp.Diagnostics.AddError(
			"Invalid value for canonical encoding",
			fmt.Sprintf(
				"The value '%s' cannot be canonically encoded as a proquint.\n\n"+
					"Error: %s\n\n"+
					"Canonical encoding accepts:\n"+
					"  - IPv4 addresses (e.g., 127.0.0.1)~>11 chars\n"+
					"  - Hexadecimal strings (e.g., 0x7f000001 or 7f000001)~>11 or 23 chars\n"+
					"  - Unsigned integers 0-4294967295~>11 chars\n"+
					"  - Unsigned integers 4294967296-18446744073709551615~>23 chars\n\n"+
					"For generating proquint-formatted IDs from arbitrary strings, use the 'idgen_proquint' data source instead.",
				data.Value.ValueString(),
				errMsg,
			),
		)
		return
	}

	// Generate the canonical proquint
	id, err := idgen.GenerateCanonicalProquint(value)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to generate canonical Proquint",
			"Could not generate Proquint: "+err.Error(),
		)
		return
	}

	data.ID = types.StringValue(id)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
