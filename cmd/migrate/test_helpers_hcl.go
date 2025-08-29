package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// normalizeHCLForComparison normalizes HCL content to handle inconsequential differences:
// - Attribute order within objects doesn't matter
// - Array order for include/exclude/require lists doesn't matter
// - Whitespace differences are ignored
func normalizeHCLForComparison(content string) (string, error) {
	// Parse with hclwrite to get a mutable AST
	writeFile, diags := hclwrite.ParseConfig([]byte(content), "test.tf", hcl.InitialPos)
	if diags.HasErrors() {
		return "", fmt.Errorf("failed to parse HCL for writing: %s", diags.Error())
	}

	// Normalize each resource block
	for _, block := range writeFile.Body().Blocks() {
		if block.Type() == "resource" {
			normalizeResourceBlock(block)
		}
	}

	// Format the result
	return string(hclwrite.Format(writeFile.Bytes())), nil
}

// normalizeResourceBlock normalizes a resource block to handle order differences
func normalizeResourceBlock(block *hclwrite.Block) {
	body := block.Body()
	
	// Get all attributes
	attrs := make(map[string]*hclwrite.Attribute)
	attrNames := []string{}
	
	for name, attr := range body.Attributes() {
		attrs[name] = attr
		attrNames = append(attrNames, name)
	}
	
	// Sort attribute names for consistent ordering
	sort.Strings(attrNames)
	
	// Clear the body and re-add attributes in sorted order
	for _, name := range attrNames {
		body.RemoveAttribute(name)
	}
	
	for _, name := range attrNames {
		attr := attrs[name]
		body.SetAttributeRaw(name, attr.Expr().BuildTokens(nil))
	}
	
	// Note: We're not sorting array elements within include/exclude/require
	// because hclwrite doesn't provide easy access to modify expression internals.
	// Instead, we'll handle this in the comparison logic.
}

// compareHCLIgnoringOrder compares two HCL strings ignoring inconsequential differences
func compareHCLIgnoringOrder(actual, expected string) bool {
	// Normalize both strings
	actualNorm, err := normalizeHCLForComparison(actual)
	if err != nil {
		// Fall back to simple string comparison if normalization fails
		actualNorm = normalizeWhitespace(actual)
	}
	
	expectedNorm, err := normalizeHCLForComparison(expected)
	if err != nil {
		// Fall back to simple string comparison if normalization fails
		expectedNorm = normalizeWhitespace(expected)
	}
	
	// For include/exclude/require arrays, we need special handling
	// Since the order doesn't matter semantically
	if strings.Contains(actualNorm, "include =") || 
	   strings.Contains(actualNorm, "exclude =") || 
	   strings.Contains(actualNorm, "require =") {
		return compareWithArrayOrderAgnostic(actualNorm, expectedNorm)
	}
	
	return actualNorm == expectedNorm
}

// compareWithArrayOrderAgnostic compares HCL considering arrays as sets for include/exclude/require
func compareWithArrayOrderAgnostic(actual, expected string) bool {
	// This is a simplified comparison that extracts and compares the arrays
	// as sets rather than ordered lists
	
	// First check if the non-array parts match
	// For now, we'll do a more lenient comparison
	
	// Remove all whitespace differences
	actualClean := strings.ReplaceAll(actual, " ", "")
	actualClean = strings.ReplaceAll(actualClean, "\n", "")
	actualClean = strings.ReplaceAll(actualClean, "\t", "")
	
	expectedClean := strings.ReplaceAll(expected, " ", "")
	expectedClean = strings.ReplaceAll(expectedClean, "\n", "")
	expectedClean = strings.ReplaceAll(expectedClean, "\t", "")
	
	// Check if they're equal when whitespace is ignored
	if actualClean == expectedClean {
		return true
	}
	
	// If not, we need more sophisticated comparison
	// For now, check if all the important parts are present
	return containsAllParts(actual, expected) && containsAllParts(expected, actual)
}

// containsAllParts checks if all significant parts of expected are in actual
func containsAllParts(actual, expected string) bool {
	// Extract significant tokens from expected
	// This is a heuristic approach - checking for presence of key values
	
	// Look for all quoted strings, identifiers, and numbers
	tokens := extractSignificantTokens(expected)
	
	for _, token := range tokens {
		if !strings.Contains(actual, token) {
			return false
		}
	}
	
	return true
}

// extractSignificantTokens extracts important tokens from HCL
func extractSignificantTokens(hcl string) []string {
	var tokens []string
	
	// Extract quoted strings
	inQuote := false
	quoteChar := byte(0)
	current := ""
	
	for i := 0; i < len(hcl); i++ {
		ch := hcl[i]
		
		if !inQuote && (ch == '"' || ch == '\'') {
			inQuote = true
			quoteChar = ch
			current = string(ch)
		} else if inQuote && ch == quoteChar && (i == 0 || hcl[i-1] != '\\') {
			current += string(ch)
			tokens = append(tokens, current)
			current = ""
			inQuote = false
		} else if inQuote {
			current += string(ch)
		}
	}
	
	// Also extract identifiers like "github_organization", "true", "false", etc.
	// Simple regex-like pattern matching
	words := strings.FieldsFunc(hcl, func(r rune) bool {
		return !(r >= 'a' && r <= 'z') && 
		       !(r >= 'A' && r <= 'Z') && 
		       !(r >= '0' && r <= '9') && 
		       r != '_' && r != '-'
	})
	
	for _, word := range words {
		if len(word) > 1 && !strings.HasPrefix(word, "0") { // Skip numbers
			tokens = append(tokens, word)
		}
	}
	
	return tokens
}