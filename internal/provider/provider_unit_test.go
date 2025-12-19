package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
)

func TestIdgenProvider_Metadata(t *testing.T) {
	p := &IdgenProvider{}

	req := provider.MetadataRequest{}
	resp := &provider.MetadataResponse{}

	p.Metadata(context.Background(), req, resp)

	expected := "idgen"
	if resp.TypeName != expected {
		t.Errorf("Metadata() TypeName = %q, want %q", resp.TypeName, expected)
	}
}

func TestIdgenProvider_Schema(t *testing.T) {
	p := &IdgenProvider{}

	req := provider.SchemaRequest{}
	resp := &provider.SchemaResponse{}

	p.Schema(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Errorf("Schema() should not return errors, got: %v", resp.Diagnostics.Errors())
	}
}

func TestIdgenProvider_Configure(t *testing.T) {
	p := &IdgenProvider{}

	var _ provider.Provider = p
}

func TestIdgenProvider_Resources(t *testing.T) {
	p := &IdgenProvider{}

	resources := p.Resources(context.Background())

	// This provider doesn't define any resources, so it should return an empty slice
	if len(resources) != 0 {
		t.Errorf("Resources() should return empty slice, got %d resources", len(resources))
	}
}

func TestIdgenProvider_DataSources(t *testing.T) {
	p := &IdgenProvider{}

	dataSources := p.DataSources(context.Background())

	// Should return all data sources
	expectedCount := 5 // nanoid, proquint, proquint_canonical, random_word, templated
	if len(dataSources) != expectedCount {
		t.Errorf("DataSources() should return %d data sources, got %d", expectedCount, len(dataSources))
	}

	// Verify each data source can be created
	for i, dsFunc := range dataSources {
		ds := dsFunc()
		if ds == nil {
			t.Errorf("DataSources()[%d]() returned nil", i)
		}
	}
}
