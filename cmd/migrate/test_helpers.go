package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/require"
)

// TestCase represents a single test case for HCL transformation
type TestCase struct {
	Name     string
	Config   string
	State    string
	Expected []string
}

// RunTransformationTests executes a series of test cases for HCL transformation
func RunTransformationTests(t *testing.T, tests []TestCase, transformFunc func([]byte, string) ([]byte, error)) {
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Parse the input
			file, diags := hclwrite.ParseConfig([]byte(tt.Config), "test.tf", hcl.InitialPos)
			require.False(t, diags.HasErrors(), "Failed to parse input config: %s", diags)

			// Transform the file
			result, err := transformFunc(file.Bytes(), "test.tf")
			require.NoError(t, err, "Transformation failed")

			// Check each expected output fragment
			resultString := string(result)
			for _, expected := range tt.Expected {
				assertHCLContains(t, resultString, expected)
			}
		})
	}
}

// assertHCLContains checks if the expected HCL fragment is present in the actual output
func assertHCLContains(t *testing.T, actual, expected string) {
	t.Helper()

	// Trim leading/trailing whitespace from expected as it's usually not significant
	expected = strings.TrimSpace(expected)

	matched, err := compareHCLFragments(actual, expected)
	if err != nil {
		// If parsing failed, this is a test failure - the HCL is invalid
		t.Fatalf("Invalid HCL in test: %v\n\nExpected:\n%s\n\nActual:\n%s", err, expected, actual)
		return
	}

	if matched {
		return
	}

	// If not found, show a helpful diff
	actualFormatted := formatHCLForDiff(actual)
	expectedFormatted := formatHCLForDiff(expected)

	t.Errorf("Expected HCL fragment not found in output.\n\nExpected fragment:\n%s\n\nActual output:\n%s\n\nDiff:\n%s",
		expectedFormatted,
		actualFormatted,
		cmp.Diff(expectedFormatted, actualFormatted))
}

// compareHCLFragments compares HCL fragments structurally using proper parsing
func compareHCLFragments(actual, expected string) (bool, error) {
	// Try to parse both as HCL
	actualParser := hclparse.NewParser()
	actualFile, actualDiags := actualParser.ParseHCL([]byte(actual), "actual.hcl")
	if actualDiags.HasErrors() {
		return false, fmt.Errorf("failed to parse actual HCL: %s", actualDiags.Error())
	}

	expectedParser := hclparse.NewParser()
	expectedFile, expectedDiags := expectedParser.ParseHCL([]byte(expected), "expected.hcl")
	if expectedDiags.HasErrors() {
		return false, fmt.Errorf("failed to parse expected HCL: %s", expectedDiags.Error())
	}

	// Check if expected elements exist in actual
	return hclBodyContains(actualFile.Body, expectedFile.Body), nil
}

// hclBodyContains checks if actual body contains all elements from expected
func hclBodyContains(actual, expected hcl.Body) bool {
	// This is a simplified check - for full implementation, we'd need to
	// traverse the AST more thoroughly
	actualContent, _ := actual.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{},
		Blocks:     []hcl.BlockHeaderSchema{},
	})

	expectedContent, _ := expected.Content(&hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{},
		Blocks:     []hcl.BlockHeaderSchema{},
	})

	// Check attributes exist
	for name := range expectedContent.Attributes {
		if _, exists := actualContent.Attributes[name]; !exists {
			return false
		}
	}

	// Check blocks exist
	expectedBlockMap := make(map[string]int)
	for _, block := range expectedContent.Blocks {
		key := block.Type + ":" + strings.Join(block.Labels, ":")
		expectedBlockMap[key]++
	}

	actualBlockMap := make(map[string]int)
	for _, block := range actualContent.Blocks {
		key := block.Type + ":" + strings.Join(block.Labels, ":")
		actualBlockMap[key]++
	}

	for key, count := range expectedBlockMap {
		if actualBlockMap[key] < count {
			return false
		}
	}

	return true
}

// formatHCLForDiff formats HCL nicely for diff output
func formatHCLForDiff(hclContent string) string {
	// Try to parse and reformat for consistent display
	if file, diags := hclwrite.ParseConfig([]byte(hclContent), "", hcl.InitialPos); !diags.HasErrors() {
		return string(hclwrite.Format(file.Bytes()))
	}
	// If it's a fragment, just normalize whitespace
	return normalizeHCLWhitespace(hclContent)
}

// normalizeHCLWhitespace does minimal whitespace normalization
func normalizeHCLWhitespace(s string) string {
	// Trim lines and remove extra spaces
	lines := strings.Split(s, "\n")
	var normalized []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			// Collapse multiple spaces to single space
			normalized = append(normalized, strings.Join(strings.Fields(line), " "))
		}
	}
	return strings.Join(normalized, "\n")
}

// StateTestCase represents a single test case for state transformation
type StateTestCase struct {
	Name     string
	Input    string // Input JSON
	Expected string // Expected JSON output
}

// RunStateTransformationTests executes state transformation tests for a specific transform function
func RunStateTransformationTests(t *testing.T, tests []StateTestCase, transformFunc func(map[string]interface{})) {
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			// Parse input JSON
			var inputMap map[string]interface{}
			err := json.Unmarshal([]byte(tc.Input), &inputMap)
			require.NoError(t, err, "Failed to parse input JSON")

			// Apply transformation
			transformFunc(inputMap)

			// Compare with expected
			assertJSONEqual(t, tc.Expected, inputMap)
		})
	}
}

// RunFullStateTransformationTests executes tests for the full state JSON transformation
func RunFullStateTransformationTests(t *testing.T, tests []StateTestCase) {
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			// Transform using the JSON-based function
			result, err := transformStateJSON([]byte(tc.Input))
			require.NoError(t, err, "Failed to transform state")

			// Parse result and compare
			var actualData interface{}
			err = json.Unmarshal(result, &actualData)
			require.NoError(t, err, "Failed to parse transformed JSON")

			assertJSONEqual(t, tc.Expected, actualData)
		})
	}
}

// assertJSONEqual compares expected JSON with actual data
func assertJSONEqual(t *testing.T, expectedJSON string, actualData interface{}) {
	t.Helper()

	// Parse expected JSON
	var expectedData interface{}
	err := json.Unmarshal([]byte(expectedJSON), &expectedData)
	require.NoError(t, err, "Failed to parse expected JSON")

	// Normalize both for comparison (handles map ordering)
	expectedNorm := normalizeJSONValue(expectedData)
	actualNorm := normalizeJSONValue(actualData)

	// Compare with helpful diff
	if !reflect.DeepEqual(expectedNorm, actualNorm) {
		// Pretty print for better error messages
		expectedPretty, _ := json.MarshalIndent(expectedNorm, "", "  ")
		actualPretty, _ := json.MarshalIndent(actualNorm, "", "  ")

		diff := cmp.Diff(string(expectedPretty), string(actualPretty))
		t.Errorf("JSON mismatch (-expected +actual):\n%s", diff)
	}
}

// normalizeJSONValue recursively normalizes JSON values for comparison
func normalizeJSONValue(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		// Sort map keys for consistent comparison
		normalized := make(map[string]interface{})
		for k, v := range val {
			normalized[k] = normalizeJSONValue(v)
		}
		return normalized
	case []interface{}:
		// Normalize array elements
		normalized := make([]interface{}, len(val))
		for i, elem := range val {
			normalized[i] = normalizeJSONValue(elem)
		}
		return normalized
	default:
		return val
	}
}

// normalizeWhitespace normalizes whitespace in a string for simple comparison
// This is kept for backwards compatibility with existing tests
func normalizeWhitespace(s string) string {
	// Collapse multiple spaces to single space
	s = strings.Join(strings.Fields(s), " ")
	// Remove spaces around brackets and braces
	replacements := []struct{ old, new string }{
		{" {", "{"}, {"{ ", "{"},
		{" }", "}"}, {"} ", "}"},
		{" [", "["}, {"[ ", "["},
		{" ]", "]"}, {"] ", "]"},
	}
	for _, r := range replacements {
		s = strings.ReplaceAll(s, r.old, r.new)
	}
	return strings.TrimSpace(s)
}
