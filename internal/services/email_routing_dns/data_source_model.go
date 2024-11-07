// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package email_routing_dns

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v3"
	"github.com/cloudflare/cloudflare-go/v3/email_routing"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EmailRoutingDNSResultDataSourceEnvelope struct {
	Result EmailRoutingDNSDataSourceModel `json:"result,computed"`
}

type EmailRoutingDNSDataSourceModel struct {
	ZoneID     types.String      `tfsdk:"zone_id" path:"zone_id,required"`
	Created    timetypes.RFC3339 `tfsdk:"created" json:"created,optional" format:"date-time"`
	ID         types.String      `tfsdk:"id" json:"id,optional"`
	Modified   timetypes.RFC3339 `tfsdk:"modified" json:"modified,optional" format:"date-time"`
	Name       types.String      `tfsdk:"name" json:"name,optional"`
	Status     types.String      `tfsdk:"status" json:"status,optional"`
	Tag        types.String      `tfsdk:"tag" json:"tag,optional"`
	Enabled    types.Bool        `tfsdk:"enabled" json:"enabled,computed_optional"`
	SkipWizard types.Bool        `tfsdk:"skip_wizard" json:"skip_wizard,computed_optional"`
}

func (m *EmailRoutingDNSDataSourceModel) toReadParams(_ context.Context) (params email_routing.DNSGetParams, diags diag.Diagnostics) {
	params = email_routing.DNSGetParams{
		ZoneID: cloudflare.F(m.ZoneID.ValueString()),
	}

	return
}
