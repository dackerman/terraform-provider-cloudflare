// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package api_token

import (
	"encoding/json"
	"fmt"

	"github.com/cloudflare/terraform-provider-cloudflare/internal/apijson"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type APITokenResultEnvelope struct {
	Result APITokenModel `json:"result"`
}

type APITokenModel struct {
	ID         types.String              `tfsdk:"id" json:"id,computed"`
	Name       types.String              `tfsdk:"name" json:"name,required"`
	Policies   *[]*APITokenPoliciesModel `tfsdk:"policies" json:"policies,required"`
	ExpiresOn  timetypes.RFC3339         `tfsdk:"expires_on" json:"expires_on,optional" format:"date-time"`
	NotBefore  timetypes.RFC3339         `tfsdk:"not_before" json:"not_before,optional" format:"date-time"`
	Condition  *APITokenConditionModel   `tfsdk:"condition" json:"condition,optional"`
	Status     types.String              `tfsdk:"status" json:"status,computed_optional"`
	IssuedOn   timetypes.RFC3339         `tfsdk:"issued_on" json:"issued_on,computed" format:"date-time"`
	LastUsedOn timetypes.RFC3339         `tfsdk:"last_used_on" json:"last_used_on,computed" format:"date-time"`
	ModifiedOn timetypes.RFC3339         `tfsdk:"modified_on" json:"modified_on,computed" format:"date-time"`
	Value      types.String              `tfsdk:"value" json:"value,computed,no_refresh"`
}

func (m APITokenModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}

func (m APITokenModel) MarshalJSONForUpdate(state APITokenModel) (data []byte, err error) {
	return apijson.MarshalForUpdate(m, state)
}

// PolicyResources is a custom type that handles JSON-encoded nested resources (v4 compatibility)
type PolicyResources map[string]types.String

// MarshalJSONWithState implements apijson.CustomMarshaler to handle JSON-encoded nested resources
func (r PolicyResources) MarshalJSONWithState(plan interface{}, state interface{}) ([]byte, error) {
	// Extract the actual map from the interface
	planResources, ok := plan.(PolicyResources)
	if !ok {
		// If it's a pointer, dereference it
		if ptr, ok := plan.(*PolicyResources); ok && ptr != nil {
			planResources = *ptr
		} else {
			return json.Marshal(plan) // Fallback to regular marshaling
		}
	}
	
	result := make(map[string]interface{})
	
	for key, val := range planResources {
		if !val.IsNull() && !val.IsUnknown() {
			strVal := val.ValueString()
			
			// Try to unmarshal as JSON object (following v4 provider pattern)
			var nestedObj map[string]string
			if err := json.Unmarshal([]byte(strVal), &nestedObj); err == nil {
				// Successfully unmarshaled as JSON - use the nested object
				result[key] = nestedObj
			} else {
				// Not JSON or unmarshaling failed - use as plain string
				result[key] = strVal
			}
		}
	}
	
	return json.Marshal(result)
}

// UnmarshalJSON implements json.Unmarshaler to handle nested resources coming from the API
func (r *PolicyResources) UnmarshalJSON(data []byte) error {
	// Check if receiver is nil
	if r == nil {
		return fmt.Errorf("cannot unmarshal into nil PolicyResources")
	}
	
	// Handle null or empty JSON
	if len(data) == 0 || string(data) == "null" || string(data) == "{}" {
		// Initialize as empty map
		if *r == nil {
			*r = make(PolicyResources)
		}
		return nil
	}
	
	// First unmarshal into a generic map to inspect the structure
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	
	// Initialize the map if nil
	if *r == nil {
		*r = make(PolicyResources)
	}
	
	// Process each entry in the raw map
	for key, val := range raw {
		switch v := val.(type) {
		case string:
			// Simple string value
			(*r)[key] = types.StringValue(v)
		case map[string]interface{}:
			// Nested object - need to marshal it back to JSON string for storage
			// This maintains compatibility with the v4 provider's jsonencode pattern
			jsonBytes, err := json.Marshal(v)
			if err != nil {
				return err
			}
			(*r)[key] = types.StringValue(string(jsonBytes))
		default:
			// For any other type, convert to string representation
			(*r)[key] = types.StringValue(fmt.Sprintf("%v", v))
		}
	}
	
	return nil
}

type APITokenPoliciesModel struct {
	ID               types.String                              `tfsdk:"id" json:"id,computed,force_encode,encode_state_for_unknown"`
	Effect           types.String                              `tfsdk:"effect" json:"effect,required"`
	PermissionGroups *[]*APITokenPoliciesPermissionGroupsModel `tfsdk:"permission_groups" json:"permission_groups,required"`
	Resources        *PolicyResources                          `tfsdk:"resources" json:"resources,required"`
}

type APITokenPoliciesPermissionGroupsModel struct {
	ID   types.String                               `tfsdk:"id" json:"id,required"`
	Meta *APITokenPoliciesPermissionGroupsMetaModel `tfsdk:"meta" json:"meta,optional"`
	Name types.String                               `tfsdk:"name" json:"name,computed"`
}

type APITokenPoliciesPermissionGroupsMetaModel struct {
	Key   types.String `tfsdk:"key" json:"key,optional"`
	Value types.String `tfsdk:"value" json:"value,optional"`
}

type APITokenConditionModel struct {
	RequestIP *APITokenConditionRequestIPModel `tfsdk:"request_ip" json:"request_ip,optional"`
}

type APITokenConditionRequestIPModel struct {
	In    *[]types.String `tfsdk:"in" json:"in,optional"`
	NotIn *[]types.String `tfsdk:"not_in" json:"not_in,optional"`
}
