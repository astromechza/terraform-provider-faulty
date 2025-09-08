// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &FaultyProvider{}

type FaultyProvider struct {
	version string
}

type FaultyProviderModel struct {
	RequiredBoolean types.Bool `tfsdk:"required_boolean"`
}

func (p *FaultyProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "faulty"
	resp.Version = p.version
}

func (p *FaultyProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"required_boolean": schema.BoolAttribute{
				MarkdownDescription: "This boolean must be set. If false, it will fail.",
				Required:            true,
			},
		},
	}
}

func (p *FaultyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data FaultyProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if !data.RequiredBoolean.ValueBool() {
		resp.Diagnostics.AddError("required_boolean_not_true", "Faulty provider requires required_boolean to be set to true.")
	}

	resp.DataSourceData = true
	resp.ResourceData = true
}

func (p *FaultyProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
	}
}

func (p *FaultyProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExampleDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &FaultyProvider{
			version: version,
		}
	}
}
