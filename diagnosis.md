# Cloudflare Ruleset Phantom Changes Bug Analysis

## Issue Summary
The `cloudflare_ruleset` resource exhibits phantom changes on every `terraform apply`, showing that the resource will be "updated in-place" even when no actual changes have been made to the configuration.

## Root Cause

The bug is located in the `ModifyPlan` function at `/home/david/code/terraform-provider-cloudflare/internal/services/ruleset/resource.go:322`.

### Code Analysis

```go
// Line 321-322 (BUGGY CODE)
planElements := make([]RulesetRulesModel, 0, len(state.Rules.Elements()))
diags = state.Rules.ElementsAs(ctx, &planElements, false)  // BUG: Using state instead of plan
```

The function is supposed to extract elements from the plan, but it's incorrectly extracting from the state twice:
1. First extraction at line 309: `state.Rules.ElementsAs(ctx, &stateElements, false)` - Correct
2. Second extraction at line 322: `state.Rules.ElementsAs(ctx, &planElements, false)` - Incorrect (should be from plan)

### Why This Causes Phantom Changes

The `ModifyPlan` function's purpose is to preserve computed rule IDs from the state when rules have matching `ref` values. The workflow should be:

1. Extract rules from state and build a map of IDs by ref
2. Extract rules from plan
3. For each plan rule with an unknown ID but known ref, copy the ID from state if a matching ref exists

However, because the function extracts state elements into `planElements`, it never actually processes the plan rules. This means:
- Plan rules keep their unknown IDs
- Terraform sees a difference between state (with IDs) and plan (without IDs)
- This triggers a perpetual "update needed" condition

## Resolution

### Fix
Change line 322 from:
```go
diags = state.Rules.ElementsAs(ctx, &planElements, false)
```

To:
```go
diags = plan.Rules.ElementsAs(ctx, &planElements, false)
```

### Complete Corrected Function Segment
```go
// Extract plan elements (corrected)
planElements := make([]RulesetRulesModel, 0, len(plan.Rules.Elements()))
diags = plan.Rules.ElementsAs(ctx, &planElements, false)
if diags != nil {
    resp.Diagnostics.Append(diags...)
    return
}

for i, rule := range planElements {
    // Do nothing if the rule's ID is a known planned value.
    if !rule.ID.IsUnknown() {
        continue
    }

    // If the rule's ref matches a rule in the state, populate the planned
    // value of its ID with the corresponding ID from the state.
    if ref := rule.Ref.ValueString(); ref != "" {
        if id, ok := ruleIDsByRef[ref]; ok {
            planElements[i].ID = id
        }
    }
}

// Convert back to list and set the plan
planList, diags := customfield.NewNestedObjectList(ctx, planElements)
if diags != nil {
    resp.Diagnostics.Append(diags...)
    return
}
plan.Rules = planList
```

Note: The loop also needs to be updated to modify `planElements[i]` directly rather than the loop variable `rule`, as Go's range loop creates copies.

## Testing

Two acceptance tests have been created to verify this bug:

1. `TestAccCloudflareRuleset_PhantomChanges` - Reproduces the exact scenario from the bug report
2. `TestAccCloudflareRuleset_ModifyPlanBug` - Focused test for the ModifyPlan function

Both tests will fail with the current implementation but should pass after applying the fix.

## Impact

This bug affects all users of the `cloudflare_ruleset` resource who use rules with computed IDs, causing:
- Unnecessary API calls on every terraform apply
- Confusing plan output showing changes when none exist
- Potential issues with automation that expects clean plans
- User workarounds with `lifecycle.ignore_changes` that shouldn't be necessary

The fix is minimal and low-risk, affecting only the plan modification logic without changing the actual resource CRUD operations.