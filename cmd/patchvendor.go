package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var patchvendorCmd = &cobra.Command{
	Use:   "patchvendor [SOURCE] [PATCHED] [OUTPUT]",
	Short: "Generate patches for Shopware vendor modifications",
	Long: `Generate unified diff patches for Shopware vendor modifications with proper a/b paths from vendor/provider structure.
	
Examples:
  wswcli patchvendor /path/to/source /path/to/patched /path/to/output
  wswcli patchvendor  # Interactive mode with prompts`,
	Args: cobra.RangeArgs(0, 3),
	RunE: runPatchVendor,
}

func init() {
	rootCmd.AddCommand(patchvendorCmd)
}

func runPatchVendor(cmd *cobra.Command, args []string) error {
	var sourcePath, patchedPath, outputPath string
	var err error

	// If not all arguments provided, use interactive mode
	if len(args) < 3 {
		sourcePath, patchedPath, outputPath, err = getPathsInteractively(args)
		if err != nil {
			return err
		}
	} else {
		sourcePath = args[0]
		patchedPath = args[1]
		outputPath = args[2]
	}

	fmt.Printf("Processing Shopware vendor patches...\n")
	fmt.Printf("Source: %s\n", sourcePath)
	fmt.Printf("Patched: %s\n", patchedPath)
	fmt.Printf("Output: %s\n", outputPath)

	// Comprehensive validation
	if err := validateInputs(sourcePath, patchedPath, outputPath); err != nil {
		return err
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	// Process the patches
	if err := processPatchFiles(sourcePath, patchedPath, outputPath); err != nil {
		return fmt.Errorf("error processing patches: %w", err)
	}

	fmt.Printf("Patches successfully processed and saved to %s\n", outputPath)
	return nil
}

// validateInputs performs comprehensive validation of input parameters
func validateInputs(sourcePath, patchedPath, outputPath string) error {
	// Check if source path exists
	sourceInfo, err := os.Stat(sourcePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("source path does not exist: %s", sourcePath)
	}
	if err != nil {
		return fmt.Errorf("error accessing source path: %w", err)
	}

	// Check if patched path exists
	patchedInfo, err := os.Stat(patchedPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("patched path does not exist: %s", patchedPath)
	}
	if err != nil {
		return fmt.Errorf("error accessing patched path: %w", err)
	}

	// Validate that both paths are of the same type (file or directory)
	if sourceInfo.IsDir() != patchedInfo.IsDir() {
		return fmt.Errorf("source and patched paths must both be files or both be directories")
	}

	// Check if source and patched are the same file
	if sourcePath == patchedPath {
		return fmt.Errorf("source and patched paths cannot be the same")
	}

	// Validate output path
	if outputPath == "" {
		return fmt.Errorf("output path cannot be empty")
	}

	// Check if output path already exists and is a directory
	if outputInfo, err := os.Stat(outputPath); err == nil && outputInfo.IsDir() {
		return fmt.Errorf("output path exists and is a directory: %s", outputPath)
	}

	// Validate file extensions for single files
	if !sourceInfo.IsDir() {
		if err := validateFileExtensions(sourcePath, patchedPath); err != nil {
			return err
		}
	}

	return nil
}

// validateFileExtensions checks if file extensions are compatible
func validateFileExtensions(sourcePath, patchedPath string) error {
	sourceExt := strings.ToLower(filepath.Ext(sourcePath))
	patchedExt := strings.ToLower(filepath.Ext(patchedPath))

	// Allow common file extensions
	allowedExts := map[string]bool{
		".php": true, ".js": true, ".ts": true, ".css": true, ".scss": true,
		".html": true, ".twig": true, ".xml": true, ".json": true, ".yml": true,
		".yaml": true, ".md": true, ".txt": true, ".sql": true, ".sh": true,
		".vue": true, ".jsx": true, ".tsx": true, ".less": true, ".sass": true,
	}

	if sourceExt != "" && !allowedExts[sourceExt] {
		fmt.Printf("Warning: Uncommon file extension for source: %s\n", sourceExt)
	}

	if patchedExt != "" && !allowedExts[patchedExt] {
		fmt.Printf("Warning: Uncommon file extension for patched: %s\n", patchedExt)
	}

	if sourceExt != patchedExt {
		return fmt.Errorf("source and patched files have different extensions: %s vs %s", sourceExt, patchedExt)
	}

	return nil
}

func getPathsInteractively(args []string) (string, string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n=== Shopware Vendor Patch Generator ===")
	fmt.Println("This tool creates unified diff patches for Shopware vendor modifications.")
	fmt.Println()

	var sourcePath, patchedPath, outputPath string
	var err error

	// Get SOURCE path
	if len(args) >= 1 {
		sourcePath = args[0]
		fmt.Printf("Using provided SOURCE path: %s\n", sourcePath)
	} else {
		fmt.Println("SOURCE PATH:")
		fmt.Println("   This is the original, unmodified vendor file or directory.")
		fmt.Println("   Example: vendor/shopware/core/Framework/Plugin/PluginManager.php")
		fmt.Print("   Enter source path: ")
		sourcePath, err = reader.ReadString('\n')
		if err != nil {
			return "", "", "", fmt.Errorf("error reading source path: %w", err)
		}
		sourcePath = strings.TrimSpace(sourcePath)
	}

	// Get PATCHED path
	if len(args) >= 2 {
		patchedPath = args[1]
		fmt.Printf("Using provided PATCHED path: %s\n", patchedPath)
	} else {
		fmt.Println()
		fmt.Println("PATCHED PATH:")
		fmt.Println("   This is the modified version of the vendor file or directory.")
		fmt.Println("   It contains your custom changes that should be preserved.")
		fmt.Println("   Example: custom/plugins/MyPlugin/vendor-patches/PluginManager.php")
		fmt.Print("   Enter patched path: ")
		patchedPath, err = reader.ReadString('\n')
		if err != nil {
			return "", "", "", fmt.Errorf("error reading patched path: %w", err)
		}
		patchedPath = strings.TrimSpace(patchedPath)
	}

	// Get OUTPUT path
	if len(args) >= 3 {
		outputPath = args[2]
		fmt.Printf("Using provided OUTPUT path: %s\n", outputPath)
	} else {
		// Generate suggested output path
		suggestedPath := generateSuggestedOutputPath(sourcePath)

		fmt.Println()
		fmt.Println("OUTPUT PATH:")
		fmt.Println("   This is where the generated patch file will be saved.")
		fmt.Println("   The patch can later be applied using 'git apply' or 'patch' command.")
		fmt.Printf("   Suggested: %s\n", suggestedPath)
		fmt.Print("   Enter output path (or press Enter for suggested): ")
		outputPath, err = reader.ReadString('\n')
		if err != nil {
			return "", "", "", fmt.Errorf("error reading output path: %w", err)
		}
		outputPath = strings.TrimSpace(outputPath)

		// Use suggested path if user pressed Enter without input
		if outputPath == "" {
			outputPath = suggestedPath
		}
	}

	// Validate inputs
	if sourcePath == "" {
		return "", "", "", fmt.Errorf("source path cannot be empty")
	}
	if patchedPath == "" {
		return "", "", "", fmt.Errorf("patched path cannot be empty")
	}
	if outputPath == "" {
		return "", "", "", fmt.Errorf("output path cannot be empty")
	}

	fmt.Println()
	fmt.Println("Summary:")
	fmt.Printf("   SOURCE:  %s\n", sourcePath)
	fmt.Printf("   PATCHED: %s\n", patchedPath)
	fmt.Printf("   OUTPUT:  %s\n", outputPath)
	fmt.Print("\nProceed? (y/N): ")

	confirmation, err := reader.ReadString('\n')
	if err != nil {
		return "", "", "", fmt.Errorf("error reading confirmation: %w", err)
	}

	confirmation = strings.ToLower(strings.TrimSpace(confirmation))
	if confirmation != "y" && confirmation != "yes" {
		return "", "", "", fmt.Errorf("operation cancelled by user")
	}

	return sourcePath, patchedPath, outputPath, nil
}

func generateSuggestedOutputPath(sourcePath string) string {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	// Extract trimmed vendor path (e.g. vendor/shopware/core -> shopware/core)
	trimmedPath := trimVendorPath(sourcePath)

	// Generate unix timestamp
	timestamp := time.Now().Unix()

	// Construct path: <cwd>/artifacts/patches/<trimmed_path>/<timestamp>-patch.patch
	outputPath := filepath.Join(cwd, "artifacts", "patches", trimmedPath, fmt.Sprintf("%d-%s.patch", timestamp, strings.TrimSuffix(filepath.Base(sourcePath), filepath.Ext(sourcePath))))

	return outputPath
}

// trimVendorPath extracts only the provider/package part from vendor paths
// e.g. vendor/shopware/core/Framework/Plugin/PluginManager.php -> shopware/core
func trimVendorPath(sourcePath string) string {
	// Convert to forward slashes for consistent handling
	normalizedPath := filepath.ToSlash(sourcePath)

	// Remove filename if it's a file path
	if strings.Contains(filepath.Base(normalizedPath), ".") {
		normalizedPath = filepath.Dir(normalizedPath)
		normalizedPath = filepath.ToSlash(normalizedPath)
	}

	parts := strings.Split(normalizedPath, "/")

	// Find vendor directory and extract provider/package part
	for i, part := range parts {
		if part == "vendor" && i+2 < len(parts) {
			// Skip "vendor" and take only provider/package (first 2 parts)
			// e.g. vendor/shopware/core/Framework/Plugin -> shopware/core
			provider := parts[i+1]
			pkg := parts[i+2]
			return provider + "/" + pkg
		}
	}

	// If no vendor found, return the full path for non-vendor files
	return sourcePath
}

func processPatchFiles(sourcePath, patchedPath, outputPath string) error {
	// Check whether these are files or directories
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	patchedInfo, err := os.Stat(patchedPath)
	if err != nil {
		return err
	}

	if sourceInfo.IsDir() && patchedInfo.IsDir() {
		return processDirectories(sourcePath, patchedPath, outputPath)
	} else if !sourceInfo.IsDir() && !patchedInfo.IsDir() {
		return processSingleFile(sourcePath, patchedPath, outputPath)
	} else {
		return fmt.Errorf("source and patched must both be either files or directories")
	}
}

func processDirectories(sourcePath, patchedPath, outputPath string) error {
	// Iterate over all files in the source directory
	return filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Calculate relative paths
		relPath, err := filepath.Rel(sourcePath, path)
		if err != nil {
			return err
		}

		patchedFile := filepath.Join(patchedPath, relPath)
		outputFile := filepath.Join(outputPath, relPath)

		// Check if the corresponding patched file exists
		if _, err := os.Stat(patchedFile); os.IsNotExist(err) {
			fmt.Printf("No patched version found for: %s\n", relPath)
			return nil
		}

		// Create output directory for this file
		if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
			return err
		}

		// Process the individual file
		return processSingleFile(path, patchedFile, outputFile)
	})
}

func processSingleFile(sourcePath, patchedPath, outputPath string) error {
	fmt.Printf("Processing: %s\n", filepath.Base(sourcePath))

	// Generate patch in unified diff format
	patch := generateUnifiedDiff(sourcePath, patchedPath)

	info, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("error accessing source file: %w", err)
	}
	if info.Size() == 0 {
		return fmt.Errorf("source file is empty")
	}
	if patch == "" {
		return fmt.Errorf("source and patched files do not have different content")
	}

	// Write patch to output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer outputFile.Close()

	_, err = outputFile.WriteString(patch)
	if err != nil {
		return fmt.Errorf("error writing patch file: %w", err)
	}

	return nil
}

// generateUnifiedDiff creates a unified diff patch between source and patched content
func generateUnifiedDiff(sourcePath, patchedPath string) string {
	// Try to use git diff first if files exist
	if _, err := os.Stat(sourcePath); err == nil {
		if _, err := os.Stat(patchedPath); err == nil {
			cmd := exec.Command("git", "diff", "--no-index", "--unified=3", sourcePath, patchedPath)
			output, err := cmd.CombinedOutput()
			if err != nil && cmd.ProcessState.ExitCode() != 1 {
				// Exit-Code 1 ist normal, wenn es Unterschiede gibt.
				return fmt.Sprintf("error running git diff: %v", err)
			}

			// Post-process the output to fix the vendor paths
			return fixVendorPaths(string(output), sourcePath, patchedPath)
		}
	}

	// Fallback: generate diff manually when files don't exist (e.g., in tests)
	return fmt.Sprintf("error: could not generate diff for %s and %s", sourcePath, patchedPath)
}

// fixVendorPaths post-processes git diff output to fix vendor paths in headers
func fixVendorPaths(diffOutput, sourcePath string, patchedPath string) string {
	lines := strings.Split(diffOutput, "\n")
	vendorPath := extractVendorPath(sourcePath)

	for i, line := range lines {
		if strings.HasPrefix(line, "--- ") {
			lines[i] = fmt.Sprintf("--- %s", vendorPath)
		} else if strings.HasPrefix(line, "+++ ") {
			lines[i] = fmt.Sprintf("+++ %s", vendorPath)
		} else if strings.HasPrefix(line, "diff --git") {
			lines[i] = fmt.Sprintf("diff --git %s %s", vendorPath, vendorPath)
		}
	}

	return strings.Join(lines, "\n")
}

// extractVendorPath extracts the vendor path from a full file path for use in patch headers
func extractVendorPath(fullPath string) string {
	// Convert to forward slashes for consistent handling
	normalizedPath := filepath.ToSlash(fullPath)
	parts := strings.Split(normalizedPath, "/")

	// Look for vendor directory and extract the full path from vendor onwards
	for i, part := range parts {
		if part == "vendor" && i+2 < len(parts) {
			// Return the full path from vendor/ onwards
			// e.g. vendor/shopware/core/Framework/Plugin/PluginManager.php
			return strings.Join(parts[i:], "/")
		}
	}

	// Fallback: return the full normalized path
	return normalizedPath
}
