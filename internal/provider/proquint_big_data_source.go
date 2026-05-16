package provider

import (
	"context"
	"math/big"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/iilei/terraform-provider-idgen/internal/idgen"
)

type bigProquintDataSource struct{}

func NewBigProquintDataSource() datasource.DataSource {
	return &bigProquintDataSource{}
}

func (d *bigProquintDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "idgen_proquint_big"
}

func (d *bigProquintDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Generates a Proquint ID from large integer seeds, including 256-bit values such as SHA-256 digests represented in hexadecimal.",
		Attributes: map[string]schema.Attribute{
			"seed": schema.StringAttribute{
				MarkdownDescription: "The seed value as a string. Supports:\n\n" +
					"- decimal integers (e.g., `115792089237316195423570985008687907853269984665640564039457584007913129639935`)\n" +
					"- hexadecimal integers with `0x` prefix (e.g., `0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF`)\n" +
					"- hexadecimal integers without prefix (e.g., `FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF`)\n\n" +
					"This is useful for encoding SHA-256 values as deterministic proquints.",
				Required: true,
			},
			"length": schema.Int64Attribute{
				MarkdownDescription: "The output length in characters. For SHA-256-sized values (32 bytes), use `95` characters.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The generated Proquint ID.",
				Computed:            true,
			},
		},
	}
}

func (d *bigProquintDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data struct {
		Seed   types.String `tfsdk:"seed"`
		Length types.Int64  `tfsdk:"length"`
		ID     types.String `tfsdk:"id"`
	}
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	length := data.Length.ValueInt64()
	byteLength := int((length + 1) / 6 * 2)
	if byteLength < 2 {
		byteLength = 2
	}

	// Convert seed to big.Int (supports decimal and hexadecimal values)
	seedStr := data.Seed.ValueString()
	seed, ok := new(big.Int).SetString(seedStr, 10)
	if !ok {
		hexStr := seedStr
		if strings.HasPrefix(strings.ToLower(seedStr), "0x") {
			hexStr = seedStr[2:]
		}
		if strings.ContainsAny(hexStr, "abcdefABCDEF") {
			seed, ok = new(big.Int).SetString(hexStr, 16)
		}
	}
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Seed",
			"The provided seed is not a valid integer (decimal or hexadecimal).",
		)
		return
	}

	// Generate Proquint ID
	proquint, err := idgen.GenerateProquint(byteLength, seed, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to generate Proquint",
			"Could not generate Proquint: "+err.Error(),
		)
		return
	}
	data.ID = types.StringValue(proquint)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}
