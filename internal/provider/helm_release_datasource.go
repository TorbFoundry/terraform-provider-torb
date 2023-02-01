package provider

import (
	"context"
	"os/exec"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jeremywohl/flatten"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &HelmReleaseDataSource{}

func NewHelmReleaseDataSource() datasource.DataSource {
	return &HelmReleaseDataSource{}
}

type HelmReleaseDataSource struct{}

type HelmReleaseDataSourceModel struct {
	ReleaseName types.String `tfsdk:"release_name"`
	Namespace   types.String `tfsdk:"namespace"`
	Values      types.String `tfsdk:"values"`
	Id          types.String `tfsdk:"id"`
}

func (d *HelmReleaseDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_helm_release"
}

func (d *HelmReleaseDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Fetches information about a given Helm release.",

		Attributes: map[string]tfsdk.Attribute{
			"release_name": {
				MarkdownDescription: "The name of the release to reference.",
				Optional:            false,
				Type:                types.StringType,
				Required:            true,
			},
			"namespace": {
				MarkdownDescription: "The name of the namespace under which the release happened.",
				Optional:            true,
				Type:                types.StringType,
				Computed:            true,
				PlanModifiers: tfsdk.AttributePlanModifiers{
					stringDefault("default"),
				},
			},
			"id": {
				MarkdownDescription: "Example identifier",
				Type:                types.StringType,
				Computed:            true,
			},
			"values": {
				MarkdownDescription: "The values of the release.",
				Type:                types.StringType,
				Computed:            true,
			},
		},
	}, nil
}

func (d *HelmReleaseDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
}

func (d *HelmReleaseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data HelmReleaseDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	release_name := data.ReleaseName.ValueString()
	namespace := data.Namespace.ValueString()

	res, err := readReleaseFromHelm(ctx, release_name, namespace)

	if err != nil {
		resp.Diagnostics.AddError("Failed to read values from Helm", err.Error())
	}

	if resp.Diagnostics.HasError() {
		return
	}

	var flat_json string

	if res == "" || res == "null" {
		flat_json = "{}"
	} else {
		flat_json, err = flatten.FlattenString(res, "", flatten.DotStyle)
	}

	if err != nil {
		resp.Diagnostics.AddError("Failed to flatten nested JSON", err.Error())
	}

	data.Id = types.StringValue(strconv.FormatInt(time.Now().Unix(), 10))
	data.Values = types.StringValue(flat_json)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func readReleaseFromHelm(ctx context.Context, release_name string, namespace string) (string, error) {
	out, err := exec.Command("helm", "get", "--namespace", namespace, "values", release_name, "-o", "json").CombinedOutput()

	str_out := string(out)

	if err != nil {
		return "", err
	}

	if str_out == "null\n" {
		str_out = "{}"
	}

	return str_out, nil
}
