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

const (
	seedFormatHex     = "hex"
	seedFormatDecimal = "decimal"
)

type bigProquintDataSource struct{}

type bigProquintDataSourceModel struct {
	Seed       types.String `tfsdk:"seed"`
	SeedFormat types.String `tfsdk:"seed_format"`
	ID         types.String `tfsdk:"id"`
	Length     types.Int64  `tfsdk:"length"`
}

// NewBigProquintDataSource returns a new idgen_proquint_big data source.
func NewBigProquintDataSource() datasource.DataSource {
	return &bigProquintDataSource{}
}

func (d *bigProquintDataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = "idgen_proquint_big"
}

func (d *bigProquintDataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Encodes a large unsigned integer, such as a SHA-256 digest, directly as a Proquint ID.\n\n" +
			"The seed is encoded as big-endian bytes, zero-padded or truncated (keeping the least significant bytes) to fit `length`.",
		Attributes: map[string]schema.Attribute{
			attrSeed: schema.StringAttribute{
				MarkdownDescription: "The unsigned integer to encode, as a string. How it is read depends on `seed_format`:\n\n" +
					"- `hex` (default): hexadecimal digits, with or without a `0x` prefix, e.g. the output of Terraform's `sha256()` " +
					"(`ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad`)\n" +
					"- `decimal`: decimal digits, e.g. `115792089237316195423570985008687907853269984665640564039457584007913129639935`\n\n" +
					"Negative values are rejected.",
				Required: true,
			},
			"seed_format": schema.StringAttribute{
				MarkdownDescription: "How `seed` is interpreted: `hex` (default) or `decimal`. " +
					"A string of only the digits 0-9 is valid in both formats with different values, so the format is never guessed.",
				Optional: true,
			},
			attrLength: schema.Int64Attribute{
				MarkdownDescription: "The output length in characters. For SHA-256-sized values (32 bytes), use `95` characters. " +
					"If the seed needs more bytes than `length` provides, the most significant bytes are dropped and a warning is emitted.",
				Required: true,
			},
			attrID: schema.StringAttribute{
				MarkdownDescription: "The generated Proquint ID.",
				Computed:            true,
			},
		},
	}
}

// parseBigSeed parses s as an unsigned integer in the given format ("hex" or "decimal").
// It returns an error message suitable for a diagnostic when parsing fails.
func parseBigSeed(s, format string) (*big.Int, string) {
	digits := s
	base := 16

	switch format {
	case seedFormatHex:
		if strings.HasPrefix(digits, "0x") || strings.HasPrefix(digits, "0X") {
			digits = digits[2:]
		}
	case seedFormatDecimal:
		base = 10
	default:
		return nil, fmt.Sprintf("seed_format must be %q or %q, got %q.", seedFormatHex, seedFormatDecimal, format)
	}

	if strings.HasPrefix(digits, "-") {
		return nil, fmt.Sprintf("The seed %q is negative; only unsigned integers can be encoded.", s)
	}

	// big.Int.SetString also accepts a leading "+"; only plain digits are allowed here.
	value, ok := new(big.Int).SetString(digits, base)
	if !ok || strings.HasPrefix(digits, "+") {
		return nil, fmt.Sprintf("The seed %q is not a valid %s integer.", s, format)
	}

	return value, ""
}

func (d *bigProquintDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data bigProquintDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	length := data.Length.ValueInt64()
	if !validateLength(length, &resp.Diagnostics) {
		return
	}

	byteLength := proquintByteLength(length)

	format := seedFormatHex
	if !data.SeedFormat.IsNull() {
		format = data.SeedFormat.ValueString()
	}

	seed, errMsg := parseBigSeed(data.Seed.ValueString(), format)
	if errMsg != "" {
		resp.Diagnostics.AddError("Invalid Seed", errMsg)
		return
	}

	if needed := len(seed.Bytes()); needed > byteLength {
		resp.Diagnostics.AddWarning(
			"Seed Truncated",
			fmt.Sprintf(
				"The seed needs %d bytes, but length=%d only holds %d bytes. "+
					"The %d most significant bytes are dropped, so different seeds can produce the same ID. "+
					"Use length=%d to encode the full value.",
				needed, length, byteLength, needed-byteLength, (needed+needed%2)*3-1,
			),
		)
	}

	id, err := idgen.GenerateProquint(byteLength, seed, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to generate Proquint",
			"Could not generate Proquint: "+err.Error(),
		)
		return
	}
	data.ID = types.StringValue(id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
