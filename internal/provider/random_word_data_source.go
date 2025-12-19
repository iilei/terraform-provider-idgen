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
var _ datasource.DataSource = &RandomWordDataSource{}

func NewRandomWordDataSource() datasource.DataSource {
	return &RandomWordDataSource{}
}

// RandomWordDataSource defines the data source implementation.
type RandomWordDataSource struct{}

// RandomWordDataSourceModel describes the data source data model.
type RandomWordDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	Seed     types.String `tfsdk:"seed"`
	Wordlist types.String `tfsdk:"wordlist"`
}

func (d *RandomWordDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_random_word"
}

func (d *RandomWordDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Generates a random word identifier\n\n" +
			"**Important:** The bundled word list contains only 20 five-letter words, providing limited randomness. " +
			"This list is intentionally small and serves merely as an example. It will not be maintained or expanded " +
			"for the reasons described in the [word list philosophy](https://github.com/iilei/terraform-provider-idgen/tree/master/internal/data/five_letter_words.txt). " +
			"For more control, provide your own custom `wordlist` with sufficient entropy for your needs.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The generated RandomWord.",
				Computed:    true,
			},
			"seed": schema.StringAttribute{
				MarkdownDescription: "Optional seed for deterministic word selection. The same seed always produces the same word.\n\n" +
					"- **Numeric strings** (e.g., `\"1\"`, `\"2\"`, `\"3\"`) - Used directly as the index selector, producing words in sequential alphabetical order\n" +
					"- **Non-numeric strings** (e.g., `\"sha256:2d711...\"`, `\"sha256:a1fce...\"`) - Hashed to produce a distributed index selector for more varied word selection\n" +
					"- **Omitted** - Generates a random word that changes on each Terraform apply\n\n" +
					"**WARNING:** Seeded identifiers are predictable and should never be used for passwords, tokens, or any security-sensitive values.",
				Optional: true,
			},
			"wordlist": schema.StringAttribute{
				MarkdownDescription: "Optional custom word list to select words from. " +
					"Provide a comma-separated string of words (e.g., `apple,banana,cherry,date`). " +
					"If omitted, the default five-letter word list is used.",
				Optional: true,
			},
		},
	}
}

func (d *RandomWordDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Provider configuration is not needed for this implementation
}

// parseWordlist parses a comma-separated string into a slice of non-blank words
func parseWordlist(input string) []string {
	if input == "" {
		return nil
	}

	parts := strings.Split(input, ",")
	words := make([]string, 0, len(parts))

	for _, word := range parts {
		word = strings.TrimSpace(word)
		if word != "" {
			words = append(words, word)
		}
	}

	return words
}

func (d *RandomWordDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RandomWordDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var seed string
	if data.Seed.IsNull() {
		seed = ""
	} else {
		seed = data.Seed.ValueString()
	}

	var wordlist []string

	if !data.Wordlist.IsNull() {
		wordlist = parseWordlist(data.Wordlist.ValueString())
	}

	// Generate the RandomWord
	id := idgen.GetWordBySeed(seed, wordlist)

	data.ID = types.StringValue(id)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
