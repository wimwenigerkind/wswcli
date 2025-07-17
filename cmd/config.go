package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config represents the configuration for wswcli
type Config struct {
	PatchVendor PatchVendorConfig `ini:"patchvendor"`
}

// PatchVendorConfig represents the patchvendor specific configuration
type PatchVendorConfig struct {
	PatchOutputDir string `ini:"patch_output_dir"`
}

// LoadConfig loads configuration from .wswcli file in current directory
func LoadConfig() (*Config, error) {
	configPath := ".wswcli"

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return default config if file doesn't exist
		return &Config{
			PatchVendor: PatchVendorConfig{
				PatchOutputDir: "artifacts/patches",
			},
		}, nil
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}
	defer file.Close()

	config := &Config{
		PatchVendor: PatchVendorConfig{
			PatchOutputDir: "artifacts/patches",
		},
	}

	scanner := bufio.NewScanner(file)
	var currentSection string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		// Check for section headers
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.ToLower(strings.Trim(line, "[]"))
			continue
		}

		// Parse key-value pairs
		if strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Remove quotes if present
			if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
				(strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
				value = value[1 : len(value)-1]
			}

			// Apply configuration based on section
			switch currentSection {
			case "patchvendor":
				switch key {
				case "patch_output_dir":
					config.PatchVendor.PatchOutputDir = value
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	return config, nil
}

// GetConfiguredOutputPath generates output path based on configuration
func (c *Config) GetConfiguredOutputPath(sourcePath string) string {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	// Extract trimmed vendor path
	trimmedPath := trimVendorPath(sourcePath)

	// Use configured output directory
	outputDir := c.PatchVendor.PatchOutputDir
	if !filepath.IsAbs(outputDir) {
		outputDir = filepath.Join(cwd, outputDir)
	}

	// Generate filename
	filename := fmt.Sprintf("%s.patch", strings.TrimSuffix(filepath.Base(sourcePath), filepath.Ext(sourcePath)))

	// Construct full path
	outputPath := filepath.Join(outputDir, trimmedPath, filename)

	return outputPath
}

// CreateExampleConfig creates an example .wswcli configuration file
func CreateExampleConfig() error {
	configContent := `# wswcli Configuration File
# This file configures project-specific settings for wswcli commands

[patchvendor]
# Directory where patch files will be saved (relative to project root)
patch_output_dir = "artifacts/patches"

# Example configuration for different project structures:
# patch_output_dir = "build/patches"
`

	file, err := os.Create(".wswcli")
	if err != nil {
		return fmt.Errorf("error creating config file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(configContent)
	if err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}
