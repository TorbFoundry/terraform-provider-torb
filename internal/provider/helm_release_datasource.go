package provider

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

	res, err := readReleaseFromHelm(ctx, data.ReleaseName.String(), data.Namespace.String())
	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.

	if err != nil {
		resp.Diagnostics.AddError("Failed to read values from Helm", err.Error())
	}

	if resp.Diagnostics.HasError() {
		return
	}

	data.Id = types.StringValue(strconv.FormatInt(time.Now().Unix(), 10))
	data.Values = types.StringValue(res)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func readReleaseFromHelm(ctx context.Context, release_name string, namespace string) (string, error) {
	cmd := exec.Command("helm", "get", "--namespace", namespace, "values", release_name, "-o", "json")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	if err := cmd.Start(); err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadAll(stdOut)

	tflog.Warn(ctx, fmt.Sprintf("%s", bytes))

	if err != nil {
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			err_msg := fmt.Sprintf("Exit code is %d, %s\n", exitError.ExitCode(), bytes)
			return "", fmt.Errorf(err_msg)
		}
	}

	return string(bytes), nil
}
