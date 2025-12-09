package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure IdgenProvider satisfies various provider interfaces.
var _ provider.Provider = &IdgenProvider{}

// IdgenProvider defines the provider implementation.
type IdgenProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// IdgenProviderModel describes the provider data model.
type IdgenProviderModel struct{}

func (p *IdgenProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "idgen"
	resp.Version = p.version
}

func (p *IdgenProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The idgen provider offers flexible, human-friendly identifier generation for Terraform.",
	}
}

func (p *IdgenProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data IdgenProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (p *IdgenProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *IdgenProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewNanoIDDataSource,
		NewProquintDataSource,
		NewTemplatedDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &IdgenProvider{
			version: version,
		}
	}
}
