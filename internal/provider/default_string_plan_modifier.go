package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// stringDefaultModifier is a plan modifier that sets a default value for a
// types.StringType attribute when it is not configured. The attribute must be
// marked as Optional and Computed. When setting the state during the resource
// Create, Read, or Update methods, this default value must also be included or
// the Terraform CLI will generate an error.
type stringDefaultModifier struct {
	Default string
}

func (m stringDefaultModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %s", m.Default)
}

func (m stringDefaultModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `%s`", m.Default)
}

func (m stringDefaultModifier) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
	// If the value is unknown or known, do not set default value.
	if !req.AttributePlan.IsNull() {
		return
	}

	// types.String must be the attr.Value produced by the attr.Type in the schema for this attribute
	// for generic plan modifiers, use
	// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework/tfsdk#ConvertValue
	// to convert into a known type.
	var str types.String
	diags := tfsdk.ValueAs(ctx, req.AttributePlan, &str)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	resp.AttributePlan = types.StringValue(m.Default)
}

func stringDefault(defaultValue string) stringDefaultModifier {
	return stringDefaultModifier{
		Default: defaultValue,
	}
}
