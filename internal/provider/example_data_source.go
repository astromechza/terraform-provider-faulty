// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ExampleDataSource{}

func NewExampleDataSource() datasource.DataSource {
	return &ExampleDataSource{}
}

type ExampleDataSource struct {
}

type FaultyDataSourceModel struct {
	Id              types.String `tfsdk:"id"`
	RequiredBoolean types.Bool   `tfsdk:"required_boolean"`
}

func (d *ExampleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_example"
}

func (d *ExampleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Faulty data source",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Computed identifier",
				Computed:            true,
			},
			"required_boolean": schema.BoolAttribute{
				MarkdownDescription: "Required boolean must be true otherwise an error will be returned.",
				Required:            true,
			},
		},
	}
}

func (d *ExampleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
}

func (d *ExampleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FaultyDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if !data.RequiredBoolean.ValueBool() {
		resp.Diagnostics.AddError("required_boolean_not_true", "Example data source requires required_boolean to be set to true.")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Id.IsUnknown() || data.Id.IsNull() {
		data.Id = types.StringValue(uuid.NewString())
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
