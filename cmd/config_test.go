package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Test loading default config when no file exists
	t.Run("DefaultConfig", func(t *testing.T) {
		// Ensure no .wswcli file exists
		os.Remove(".wswcli")

		config, err := LoadConfig()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if config.PatchVendor.PatchOutputDir != "artifacts/patches" {
			t.Errorf("Expected default patch_output_dir 'artifacts/patches', got '%s'", config.PatchVendor.PatchOutputDir)
		}

	})

	// Test loading custom config
	t.Run("CustomConfig", func(t *testing.T) {
		// Create test config file
		configContent := `[patchvendor]
patch_output_dir = "build/patches"
`
		err := os.WriteFile(".wswcli", []byte(configContent), 0644)
		if err != nil {
			t.Fatalf("Failed to create test config: %v", err)
		}
		defer os.Remove(".wswcli")

		config, err := LoadConfig()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if config.PatchVendor.PatchOutputDir != "build/patches" {
			t.Errorf("Expected patch_output_dir 'build/patches', got '%s'", config.PatchVendor.PatchOutputDir)
		}
	})
}

func TestGetConfiguredOutputPath(t *testing.T) {
	config := &Config{
		PatchVendor: PatchVendorConfig{
			PatchOutputDir: "build/patches",
		},
	}

	sourcePath := "vendor/shopware/core/Framework/Plugin/PluginManager.php"
	outputPath := config.GetConfiguredOutputPath(sourcePath)

	// Check that the path contains the configured output directory
	if !filepath.IsAbs(outputPath) {
		cwd, _ := os.Getwd()
		expectedPrefix := filepath.Join(cwd, "build/patches")
		if !filepath.HasPrefix(outputPath, expectedPrefix) {
			t.Errorf("Expected output path to start with '%s', got '%s'", expectedPrefix, outputPath)
		}
	}

	// Check that it contains the vendor path
	if !strings.Contains(outputPath, "shopware/core") {
		t.Errorf("Expected output path to contain 'shopware/core', got '%s'", outputPath)
	}

	// Check that it ends with .patch
	if !strings.HasSuffix(outputPath, ".patch") {
		t.Errorf("Expected output path to end with '.patch', got '%s'", outputPath)
	}
}

func TestCreateExampleConfig(t *testing.T) {
	// Remove existing config if any
	os.Remove(".wswcli")

	err := CreateExampleConfig()
	if err != nil {
		t.Fatalf("Expected no error creating config, got %v", err)
	}
	defer os.Remove(".wswcli")

	// Check if file was created
	if _, err := os.Stat(".wswcli"); os.IsNotExist(err) {
		t.Error("Expected .wswcli file to be created")
	}

	// Check if file contains expected content
	content, err := os.ReadFile(".wswcli")
	if err != nil {
		t.Fatalf("Failed to read created config file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "[patchvendor]") {
		t.Error("Expected config file to contain [patchvendor] section")
	}

	if !strings.Contains(contentStr, "patch_output_dir") {
		t.Error("Expected config file to contain patch_output_dir setting")
	}
}
