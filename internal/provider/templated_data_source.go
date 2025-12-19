package provider

import (
	"bytes"
	"context"
	_ "embed"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/iilei/terraform-provider-idgen/internal/idgen"
)

//go:embed docs_embeds/templated_data_source.md
var templateFunctionsDocs string

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &TemplatedDataSource{}

func NewTemplatedDataSource() datasource.DataSource {
	return &TemplatedDataSource{}
}

// TemplatedDataSource defines the data source implementation.
type TemplatedDataSource struct{}

// TemplatedDataSourceModel describes the data source data model.
type TemplatedDataSourceModel struct {
	ID                types.String `tfsdk:"id"`
	Template          types.String `tfsdk:"template"`
	Proquint          types.Object `tfsdk:"proquint"`
	ProquintCanonical types.Object `tfsdk:"proquint_canonical"`
	NanoID            types.Object `tfsdk:"nanoid"`
	RandomWord        types.Object `tfsdk:"random_word"`
}

// ProquintConfig holds configuration for proquint generation
type ProquintConfig struct {
	Length    types.Int64  `tfsdk:"length"`
	Seed      types.String `tfsdk:"seed"`
	GroupSize types.Int64  `tfsdk:"group_size"`
}

// ProquintCanonicalConfig holds configuration for canonical proquint generation
type ProquintCanonicalConfig struct {
	Seed      types.String `tfsdk:"seed"`
	GroupSize types.Int64  `tfsdk:"group_size"`
}

// NanoIDConfig holds configuration for nanoid generation
type NanoIDConfig struct {
	Length    types.Int64  `tfsdk:"length"`
	Seed      types.String `tfsdk:"seed"`
	GroupSize types.Int64  `tfsdk:"group_size"`
	Alphabet  types.String `tfsdk:"alphabet"`
}

// RandomWordConfig holds configuration for random word generation
type RandomWordConfig struct {
	Seed     types.String `tfsdk:"seed"`
	Wordlist types.String `tfsdk:"wordlist"`
}

func (d *TemplatedDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_templated"
}

func (d *TemplatedDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	// Base attributes shared by most components
	baseAttributes := map[string]schema.Attribute{
		"seed": schema.StringAttribute{
			Optional:    true,
			Description: "Seed for deterministic generation",
		},
		"group_size": schema.Int64Attribute{
			Optional:    true,
			Description: "Number of characters per group separated by dashes",
		},
	}

	// Proquint schema (length + base)
	proquintAttributes := map[string]schema.Attribute{
		"length": schema.Int64Attribute{
			Optional:    true,
			Description: "Length of the generated Proquint (default: 11)",
		},
	}
	for k, v := range baseAttributes {
		proquintAttributes[k] = v
	}

	// Proquint canonical schema (only seed + group_size, seed required)
	proquintCanonicalAttributes := map[string]schema.Attribute{
		"seed": schema.StringAttribute{
			Required:            true,
			MarkdownDescription: "Seed value (IPv4 address, hex string, or integer) for canonical encoding",
		},
		"group_size": schema.Int64Attribute{
			Optional:    true,
			Description: "Number of characters per group separated by dashes",
		},
	}

	// NanoID schema (length + alphabet + base)
	nanoidAttributes := map[string]schema.Attribute{
		"length": schema.Int64Attribute{
			Optional:    true,
			Description: "Length of the generated NanoID (default: 21)",
		},
		"alphabet": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Alphabet preset (`alphanumeric`, `numeric`, `readable`) or custom alphabet string. Default: `alphanumeric`",
		},
	}
	for k, v := range baseAttributes {
		nanoidAttributes[k] = v
	}

	// Random word schema (only seed + wordlist, no length or group_size)
	randomWordAttributes := map[string]schema.Attribute{
		"seed": schema.StringAttribute{
			Optional:    true,
			Description: "Seed for deterministic word selection",
		},
		"wordlist": schema.StringAttribute{
			Optional:            true,
			MarkdownDescription: "Comma-separated custom word list (uses default 5-letter word list if omitted). See [random_word](./random_word) for more details about the word list limitations.",
		},
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: "Generates a templated identifier combining multiple ID types.\n\n" +
			"Use Go template syntax with `.proquint`, `.proquint_canonical`, `.nanoid`, and `.random_word` variables. " +
			"Example: `{{ .proquint }}-{{ .nanoid }}`\n\n" +
			templateFunctionsDocs,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The generated templated ID.",
				Computed:    true,
			},
			"template": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Go template string with `.proquint`, `.proquint_canonical`, `.nanoid`, and `.random_word` variables",
			},
			"proquint": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Proquint component configuration. See [proquint](./proquint) for more details.",
				Attributes:          proquintAttributes,
			},
			"proquint_canonical": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Canonical Proquint component (encodes IPv4 addresses or integers). See [proquint_canonical](./proquint_canonical) for more details.",
				Attributes:          proquintCanonicalAttributes,
			},
			"nanoid": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "NanoID component configuration. See [nanoid](./nanoid) for more details.",
				Attributes:          nanoidAttributes,
			},
			"random_word": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Random word component configuration. See [random_word](./random_word) for more details.",
				Attributes:          randomWordAttributes,
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

	// Generate ID components
	idComponents := make(map[string]string)

	// Generate proquint if configured
	if !data.Proquint.IsNull() {
		var config ProquintConfig
		resp.Diagnostics.Append(data.Proquint.As(ctx, &config, basetypes.ObjectAsOptions{})...)
		if !resp.Diagnostics.HasError() {
			id := generateProquint(config)
			idComponents["proquint"] = id
		}
	}

	// Generate proquint_canonical if configured
	if !data.ProquintCanonical.IsNull() {
		var config ProquintCanonicalConfig
		resp.Diagnostics.Append(data.ProquintCanonical.As(ctx, &config, basetypes.ObjectAsOptions{})...)
		if !resp.Diagnostics.HasError() {
			id := generateProquintCanonical(config, &resp.Diagnostics)
			idComponents["proquint_canonical"] = id
		}
	}

	// Generate nanoid if configured
	if !data.NanoID.IsNull() {
		var config NanoIDConfig
		resp.Diagnostics.Append(data.NanoID.As(ctx, &config, basetypes.ObjectAsOptions{})...)
		if !resp.Diagnostics.HasError() {
			id, err := generateNanoID(config, &resp.Diagnostics)
			if err != nil {
				resp.Diagnostics.AddError("Failed to generate NanoID", err.Error())
				return
			}
			idComponents["nanoid"] = id
		}
	}

	// Generate random_word if configured
	if !data.RandomWord.IsNull() {
		var config RandomWordConfig
		resp.Diagnostics.Append(data.RandomWord.As(ctx, &config, basetypes.ObjectAsOptions{})...)
		if !resp.Diagnostics.HasError() {
			id := generateRandomWord(config)
			idComponents["random_word"] = id
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Apply template with custom functions
	templateStr := data.Template.ValueString()
	tmpl, err := template.New("id").Funcs(templateFuncs()).Parse(templateStr)
	if err != nil {
		resp.Diagnostics.AddError("Invalid template", err.Error())
		return
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, idComponents); err != nil {
		resp.Diagnostics.AddError("Failed to execute template", err.Error())
		return
	}

	data.ID = types.StringValue(buf.String())

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Helper functions to generate IDs
func generateProquint(config ProquintConfig) string {
	length := 11
	if !config.Length.IsNull() {
		length = int(config.Length.ValueInt64())
	}

	var seed *int64
	if !config.Seed.IsNull() {
		seedVal, _ := idgen.StringToSeed(config.Seed.ValueString())
		seed = &seedVal
	}

	byteLength := (length + 1) / 6 * 2
	if byteLength < 2 {
		byteLength = 2
	}

	id, _ := idgen.GenerateProquint(byteLength, seed, false)

	// Determine group size (default to 5 for standard proquint format)
	groupSize := 5
	if !config.GroupSize.IsNull() {
		groupSize = int(config.GroupSize.ValueInt64())
	}

	// Remove all dashes and apply grouping
	if groupSize > 0 {
		id = strings.ReplaceAll(id, "-", "")
		id = idgen.ApplyGrouping(id, groupSize)
	}

	return id
}

func generateNanoID(config NanoIDConfig, diags *diag.Diagnostics) (string, error) {
	length := 21
	if !config.Length.IsNull() {
		length = int(config.Length.ValueInt64())
	}

	alphabet := idgen.Alphanumeric
	if !config.Alphabet.IsNull() {
		alphabetStr := config.Alphabet.ValueString()
		switch strings.ToLower(alphabetStr) {
		case "alphanumeric":
			alphabet = idgen.Alphanumeric
		case "numeric":
			alphabet = idgen.Numeric
		case "readable":
			alphabet = idgen.Readable
		default:
			alphabet = alphabetStr
		}
	}

	// Warn if alphabet contains dashes and grouping is enabled
	if !config.GroupSize.IsNull() && config.GroupSize.ValueInt64() > 0 {
		if strings.Contains(alphabet, "-") {
			diags.AddWarning(
				warningAlphabetContainsDashTitle,
				warningAlphabetContainsDashDetail,
			)
		}
	}

	var seed *int64
	if !config.Seed.IsNull() {
		seedVal, _ := idgen.StringToSeed(config.Seed.ValueString())
		seed = &seedVal
	}

	// Determine group size for length calculation
	groupSize := 0
	if !config.GroupSize.IsNull() {
		groupSize = int(config.GroupSize.ValueInt64())
	}

	// Generate the NanoID (grouping is applied internally if groupSize > 0)
	return idgen.GenerateNanoID(alphabet, length, seed, groupSize)
}

func generateProquintCanonical(config ProquintCanonicalConfig, diags *diag.Diagnostics) string {
	if config.Seed.IsNull() {
		diags.AddError("Seed required", "proquint_canonical requires a seed value")
		return ""
	}

	value, _, errMsg := stringToCanonicalValue(config.Seed.ValueString())
	if errMsg != "" {
		diags.AddError("Invalid seed for canonical proquint", errMsg)
		return ""
	}

	id, _ := idgen.GenerateCanonicalProquint(value)

	// Determine group size (default to 5 for standard proquint format)
	groupSize := 5
	if !config.GroupSize.IsNull() {
		groupSize = int(config.GroupSize.ValueInt64())
	}

	// Remove all dashes and apply grouping
	if groupSize > 0 {
		id = strings.ReplaceAll(id, "-", "")
		id = idgen.ApplyGrouping(id, groupSize)
	}

	return id
}

func generateRandomWord(config RandomWordConfig) string {
	seed := ""
	if !config.Seed.IsNull() {
		seed = config.Seed.ValueString()
	}

	var wordlist []string
	if !config.Wordlist.IsNull() {
		wordlist = parseWordlist(config.Wordlist.ValueString())
	}

	return idgen.GetWordBySeed(seed, wordlist)
}

// templateFuncs returns custom template functions for string manipulation.
// Functions are pipe-friendly: the piped value is the last parameter.
func templateFuncs() template.FuncMap {
	return template.FuncMap{
		// Case conversion
		"upper": strings.ToUpper,
		"lower": strings.ToLower,

		// String manipulation
		"replace": func(old, new, s string) string {
			return strings.ReplaceAll(s, old, new)
		},
		"prepend": func(prefix, s string) string {
			return prefix + s
		},
		"append": func(suffix, s string) string {
			return s + suffix
		},
		"substr": func(start, length int, s string) string {
			runes := []rune(s)
			if start < 0 || start >= len(runes) {
				return ""
			}
			end := start + length
			if end > len(runes) {
				end = len(runes)
			}
			return string(runes[start:end])
		},
		"trim": strings.TrimSpace,
		"trimPrefix": func(prefix, s string) string {
			return strings.TrimPrefix(s, prefix)
		},
		"trimSuffix": func(suffix, s string) string {
			return strings.TrimSuffix(s, suffix)
		},

		// Repetition and reversal
		"repeat": func(count int, s string) string {
			return strings.Repeat(s, count)
		},
		"reverse": func(s string) string {
			runes := []rune(s)
			for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
				runes[i], runes[j] = runes[j], runes[i]
			}
			return string(runes)
		},
	}
}
