package provider

import (
	"strings"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/iilei/terraform-provider-idgen/internal/idgen"
)

func TestTemplateFuncs(t *testing.T) {
	funcs := templateFuncs()

	t.Run("upper function", func(t *testing.T) {
		if upperFunc, ok := funcs["upper"]; ok {
			result := upperFunc.(func(string) string)("hello")
			if result != "HELLO" {
				t.Errorf("upper(\"hello\") = %q, want \"HELLO\"", result)
			}
		} else {
			t.Error("upper function not found")
		}
	})

	t.Run("lower function", func(t *testing.T) {
		if lowerFunc, ok := funcs["lower"]; ok {
			result := lowerFunc.(func(string) string)("HELLO")
			if result != "hello" {
				t.Errorf("lower(\"HELLO\") = %q, want \"hello\"", result)
			}
		} else {
			t.Error("lower function not found")
		}
	})

	t.Run("replace function", func(t *testing.T) {
		if replaceFunc, ok := funcs["replace"]; ok {
			result := replaceFunc.(func(string, string, string) string)("old", "new", "old text old")
			if result != "new text new" {
				t.Errorf("replace(\"old\", \"new\", \"old text old\") = %q, want \"new text new\"", result)
			}
		} else {
			t.Error("replace function not found")
		}
	})

	t.Run("prepend function", func(t *testing.T) {
		if prependFunc, ok := funcs["prepend"]; ok {
			result := prependFunc.(func(string, string) string)("pre-", "text")
			if result != "pre-text" {
				t.Errorf("prepend(\"pre-\", \"text\") = %q, want \"pre-text\"", result)
			}
		} else {
			t.Error("prepend function not found")
		}
	})

	t.Run("append function", func(t *testing.T) {
		if appendFunc, ok := funcs["append"]; ok {
			result := appendFunc.(func(string, string) string)("-suf", "text")
			if result != "text-suf" {
				t.Errorf("append(\"-suf\", \"text\") = %q, want \"text-suf\"", result)
			}
		} else {
			t.Error("append function not found")
		}
	})

	t.Run("substr function", func(t *testing.T) {
		if substrFunc, ok := funcs["substr"]; ok {
			fn := substrFunc.(func(int, int, string) string)

			// Normal case
			result := fn(1, 3, "hello")
			if result != "ell" {
				t.Errorf("substr(1, 3, \"hello\") = %q, want \"ell\"", result)
			}

			// Start out of bounds
			result = fn(10, 3, "hello")
			if result != "" {
				t.Errorf("substr(10, 3, \"hello\") = %q, want \"\"", result)
			}

			// Length exceeds string
			result = fn(2, 10, "hello")
			if result != "llo" {
				t.Errorf("substr(2, 10, \"hello\") = %q, want \"llo\"", result)
			}

			// Negative start
			result = fn(-1, 3, "hello")
			if result != "" {
				t.Errorf("substr(-1, 3, \"hello\") = %q, want \"\"", result)
			}
		} else {
			t.Error("substr function not found")
		}
	})

	t.Run("repeat function", func(t *testing.T) {
		if repeatFunc, ok := funcs["repeat"]; ok {
			result := repeatFunc.(func(int, string) string)(3, "ha")
			if result != "hahaha" {
				t.Errorf("repeat(3, \"ha\") = %q, want \"hahaha\"", result)
			}
		} else {
			t.Error("repeat function not found")
		}
	})

	t.Run("reverse function", func(t *testing.T) {
		if reverseFunc, ok := funcs["reverse"]; ok {
			result := reverseFunc.(func(string) string)("hello")
			if result != "olleh" {
				t.Errorf("reverse(\"hello\") = %q, want \"olleh\"", result)
			}
		} else {
			t.Error("reverse function not found")
		}
	})

	t.Run("trim function", func(t *testing.T) {
		if trimFunc, ok := funcs["trim"]; ok {
			result := trimFunc.(func(string) string)("  hello  ")
			if result != "hello" {
				t.Errorf("trim(\"  hello  \") = %q, want \"hello\"", result)
			}
		} else {
			t.Error("trim function not found")
		}
	})

	t.Run("trimPrefix function", func(t *testing.T) {
		if trimPrefixFunc, ok := funcs["trimPrefix"]; ok {
			result := trimPrefixFunc.(func(string, string) string)("pre-", "pre-text")
			if result != "text" {
				t.Errorf("trimPrefix(\"pre-\", \"pre-text\") = %q, want \"text\"", result)
			}
		} else {
			t.Error("trimPrefix function not found")
		}
	})

	t.Run("trimSuffix function", func(t *testing.T) {
		if trimSuffixFunc, ok := funcs["trimSuffix"]; ok {
			result := trimSuffixFunc.(func(string, string) string)("-suf", "text-suf")
			if result != "text" {
				t.Errorf("trimSuffix(\"-suf\", \"text-suf\") = %q, want \"text\"", result)
			}
		} else {
			t.Error("trimSuffix function not found")
		}
	})
}

func TestTemplateFuncIntegration(t *testing.T) {
	// Test that all functions work in a template
	tmplText := `{{upper "hello"}} {{lower "WORLD"}} {{prepend "pre-" "text"}} {{append "-suf" "text"}} {{repeat 2 "x"}} {{reverse "abc"}}`

	tmpl, err := template.New("test").Funcs(templateFuncs()).Parse(tmplText)
	if err != nil {
		t.Fatalf("template parse error: %v", err)
	}

	var result strings.Builder
	err = tmpl.Execute(&result, nil)
	if err != nil {
		t.Fatalf("template execute error: %v", err)
	}

	expected := "HELLO world pre-text text-suf xx cba"
	if result.String() != expected {
		t.Errorf("template result = %q, want %q", result.String(), expected)
	}
}

func TestGenerateProquintCanonical_ErrorPaths(t *testing.T) {
	// These tests target the missing coverage in generateProquintCanonical

	t.Run("null seed error", func(t *testing.T) {
		var diags diag.Diagnostics

		// Create config with null seed
		config := ProquintCanonicalConfig{
			Seed:      types.StringNull(),
			GroupSize: types.Int64Value(5),
		}

		result := generateProquintCanonical(config, &diags)

		if result != "" {
			t.Errorf("Expected empty result for null seed, got: %q", result)
		}

		if !diags.HasError() {
			t.Error("Expected diagnostic error for null seed")
		}

		// Check error message
		if len(diags.Errors()) > 0 {
			errMsg := diags.Errors()[0].Summary()
			if errMsg != "Seed required" {
				t.Errorf("Expected error 'Seed required', got: %q", errMsg)
			}
		}
	})

	t.Run("invalid seed error", func(t *testing.T) {
		var diags diag.Diagnostics

		// Create config with invalid seed (non-IP, non-integer)
		config := ProquintCanonicalConfig{
			Seed:      types.StringValue("invalid-seed-value"),
			GroupSize: types.Int64Value(5),
		}

		result := generateProquintCanonical(config, &diags)

		if result != "" {
			t.Errorf("Expected empty result for invalid seed, got: %q", result)
		}

		if !diags.HasError() {
			t.Error("Expected diagnostic error for invalid seed")
		}

		// Check error message
		if len(diags.Errors()) > 0 {
			errMsg := diags.Errors()[0].Summary()
			if errMsg != "Invalid seed for canonical proquint" {
				t.Errorf("Expected error 'Invalid seed for canonical proquint', got: %q", errMsg)
			}
		}
	})

	t.Run("valid seed success", func(t *testing.T) {
		var diags diag.Diagnostics

		// Create config with valid seed
		config := ProquintCanonicalConfig{
			Seed:      types.StringValue("127.0.0.1"),
			GroupSize: types.Int64Value(5),
		}

		result := generateProquintCanonical(config, &diags)

		if result == "" {
			t.Error("Expected non-empty result for valid seed")
		}

		if diags.HasError() {
			t.Errorf("Unexpected error: %v", diags.Errors())
		}

		// Should be a valid proquint
		if !strings.Contains(result, "-") {
			t.Errorf("Expected grouped proquint with dashes, got: %q", result)
		}
	})

	t.Run("null group size uses default", func(t *testing.T) {
		var diags diag.Diagnostics

		// Create config with null group size
		config := ProquintCanonicalConfig{
			Seed:      types.StringValue("42"),
			GroupSize: types.Int64Null(),
		}

		result := generateProquintCanonical(config, &diags)

		if result == "" {
			t.Error("Expected non-empty result")
		}

		if diags.HasError() {
			t.Errorf("Unexpected error: %v", diags.Errors())
		}
	})
}

func TestGenerateNanoID_AlphabetCases(t *testing.T) {
	// Test the missing coverage cases for alphabet selection in generateNanoID

	t.Run("alphanumeric alphabet case", func(t *testing.T) {
		var diags diag.Diagnostics

		config := NanoIDConfig{
			Length:   types.Int64Value(10),
			Alphabet: types.StringValue("alphanumeric"), // This should trigger the "alphanumeric" case
			Seed:     types.StringValue("test-seed"),
		}

		result, err := generateNanoID(config, &diags)

		if err != nil {
			t.Fatalf("generateNanoID failed: %v", err)
		}

		if result == "" {
			t.Error("Expected non-empty result")
		}

		// Verify all characters are from alphanumeric alphabet
		for _, char := range result {
			if !strings.ContainsRune(idgen.Alphanumeric, char) {
				t.Errorf("Character %c not in alphanumeric alphabet", char)
			}
		}

		if diags.HasError() {
			t.Errorf("Unexpected diagnostics errors: %v", diags.Errors())
		}
	})

	t.Run("numeric alphabet case", func(t *testing.T) {
		var diags diag.Diagnostics

		config := NanoIDConfig{
			Length:   types.Int64Value(8),
			Alphabet: types.StringValue("numeric"), // This should trigger the "numeric" case
			Seed:     types.StringValue("test-seed"),
		}

		result, err := generateNanoID(config, &diags)

		if err != nil {
			t.Fatalf("generateNanoID failed: %v", err)
		}

		if result == "" {
			t.Error("Expected non-empty result")
		}

		// Verify all characters are from numeric alphabet
		for _, char := range result {
			if !strings.ContainsRune(idgen.Numeric, char) {
				t.Errorf("Character %c not in numeric alphabet", char)
			}
		}

		if diags.HasError() {
			t.Errorf("Unexpected diagnostics errors: %v", diags.Errors())
		}
	})

	t.Run("custom alphabet default case", func(t *testing.T) {
		var diags diag.Diagnostics

		customAlphabet := "ABCDEF123456"
		config := NanoIDConfig{
			Length:   types.Int64Value(6),
			Alphabet: types.StringValue(customAlphabet), // This should trigger the default case
			Seed:     types.StringValue("test-seed"),
		}

		result, err := generateNanoID(config, &diags)

		if err != nil {
			t.Fatalf("generateNanoID failed: %v", err)
		}

		if result == "" {
			t.Error("Expected non-empty result")
		}

		// Verify all characters are from custom alphabet
		for _, char := range result {
			if !strings.ContainsRune(customAlphabet, char) {
				t.Errorf("Character %c not in custom alphabet %s", char, customAlphabet)
			}
		}

		if diags.HasError() {
			t.Errorf("Unexpected diagnostics errors: %v", diags.Errors())
		}
	})

	t.Run("case insensitive alphabet matching", func(t *testing.T) {
		// Test that "ALPHANUMERIC" (uppercase) still matches the alphanumeric case
		var diags diag.Diagnostics

		config := NanoIDConfig{
			Length:   types.Int64Value(5),
			Alphabet: types.StringValue("ALPHANUMERIC"), // Uppercase should still match
			Seed:     types.StringValue("test-seed"),
		}

		result, err := generateNanoID(config, &diags)

		if err != nil {
			t.Fatalf("generateNanoID failed: %v", err)
		}

		// Should use alphanumeric alphabet (case-insensitive match)
		for _, char := range result {
			if !strings.ContainsRune(idgen.Alphanumeric, char) {
				t.Errorf("Character %c not in alphanumeric alphabet", char)
			}
		}

		// Test "NUMERIC" as well
		config.Alphabet = types.StringValue("NUMERIC")
		result2, err2 := generateNanoID(config, &diags)

		if err2 != nil {
			t.Fatalf("generateNanoID failed: %v", err2)
		}

		for _, char := range result2 {
			if !strings.ContainsRune(idgen.Numeric, char) {
				t.Errorf("Character %c not in numeric alphabet", char)
			}
		}
	})
}

func TestGenerateNanoID_DashWarning(t *testing.T) {
	// Test the missing coverage for dash warning in generateNanoID
	t.Run("alphabet with dash and grouping triggers warning", func(t *testing.T) {
		var diags diag.Diagnostics

		config := NanoIDConfig{
			Length:    types.Int64Value(10),
			Alphabet:  types.StringValue("ABC-DEF123"), // Contains dash
			GroupSize: types.Int64Value(3),             // Grouping enabled
			Seed:      types.StringValue("test-seed"),
		}

		result, err := generateNanoID(config, &diags)

		if err != nil {
			t.Fatalf("generateNanoID failed: %v", err)
		}

		if result == "" {
			t.Error("Expected non-empty result")
		}

		// Should have a warning diagnostic
		if !diags.HasError() && len(diags.Warnings()) == 0 {
			t.Error("Expected warning diagnostic for alphabet containing dash")
		}

		// Check warning message
		warnings := diags.Warnings()
		if len(warnings) > 0 {
			if warnings[0].Summary() != "Alphabet contains dash character" {
				t.Errorf("Expected warning 'Alphabet contains dash character', got: %q", warnings[0].Summary())
			}
		}
	})

	t.Run("alphabet with dash but no grouping - no warning", func(t *testing.T) {
		var diags diag.Diagnostics

		config := NanoIDConfig{
			Length:    types.Int64Value(10),
			Alphabet:  types.StringValue("ABC-DEF123"), // Contains dash
			GroupSize: types.Int64Null(),               // No grouping
			Seed:      types.StringValue("test-seed"),
		}

		result, err := generateNanoID(config, &diags)

		if err != nil {
			t.Fatalf("generateNanoID failed: %v", err)
		}

		if result == "" {
			t.Error("Expected non-empty result")
		}

		// Should NOT have a warning diagnostic
		if len(diags.Warnings()) > 0 {
			t.Errorf("Unexpected warning diagnostic: %v", diags.Warnings())
		}
	})
}

func TestGenerateProquint_EdgeCases(t *testing.T) {
	// Test missing coverage in generateProquint
	t.Run("very small length calculation", func(t *testing.T) {
		config := ProquintConfig{
			Length:    types.Int64Value(1), // Very small length
			Seed:      types.StringValue("test-seed"),
			GroupSize: types.Int64Value(5),
		}

		result := generateProquint(config)

		if result == "" {
			t.Error("Expected non-empty result even with small length")
		}
	})

	t.Run("zero group size skips grouping logic", func(t *testing.T) {
		config := ProquintConfig{
			Length:    types.Int64Value(11),
			Seed:      types.StringValue("test-seed"),
			GroupSize: types.Int64Value(0), // Zero group size
		}

		result := generateProquint(config)

		if result == "" {
			t.Error("Expected non-empty result")
		}

		// With group size 0, the if condition groupSize > 0 is false,
		// so the grouping logic (removing dashes and re-applying) is skipped
		// The result should be whatever the raw proquint generation returns
		t.Logf("Result with group size 0: %s", result)
	})
}
