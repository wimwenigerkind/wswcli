package cmd

import (
	"bufio"
	"fmt"
	"os"
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
	outputPath := filepath.Join(cwd, "artifacts", "patches", trimmedPath, fmt.Sprintf("%d-patch.patch", timestamp))

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

	// Read source file
	sourceContent, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("error reading source file: %w", err)
	}

	// Read patched file
	patchedContent, err := os.ReadFile(patchedPath)
	if err != nil {
		return fmt.Errorf("error reading patched file: %w", err)
	}

	// Generate patch in unified diff format
	patch := generateUnifiedDiff(sourcePath, patchedPath, string(sourceContent), string(patchedContent))

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
func generateUnifiedDiff(sourcePath, patchedPath, sourceContent, patchedContent string) string {
	// Extract vendor/provider path from source path
	vendorPath := extractVendorPath(sourcePath)

	// Split content into lines for line-based diff
	sourceLines := strings.Split(sourceContent, "\n")
	patchedLines := strings.Split(patchedContent, "\n")

	// Generate diff operations
	diffs := computeLineDiff(sourceLines, patchedLines)

	// Build unified diff format
	var result strings.Builder

	// Write headers
	result.WriteString(fmt.Sprintf("--- a/%s\n", vendorPath))
	result.WriteString(fmt.Sprintf("+++ b/%s\n", vendorPath))

	// Generate hunks from diffs
	hunks := generateHunks(diffs, sourceLines, patchedLines)
	for _, hunk := range hunks {
		result.WriteString(hunk)
	}

	return result.String()
}

// DiffOperation represents a single diff operation
type DiffOperation struct {
	Type        string // "equal", "delete", "insert"
	SourceStart int    // Starting line in source (0-based)
	SourceCount int    // Number of lines in source
	PatchStart  int    // Starting line in patch (0-based)
	PatchCount  int    // Number of lines in patch
	Lines       []string
}

// computeLineDiff computes line-based differences between source and patched content
func computeLineDiff(sourceLines, patchedLines []string) []DiffOperation {
	var operations []DiffOperation

	sourceIdx, patchIdx := 0, 0

	for sourceIdx < len(sourceLines) || patchIdx < len(patchedLines) {
		// Find the next difference
		equalStart := sourceIdx
		equalPatchStart := patchIdx

		// Skip equal lines
		for sourceIdx < len(sourceLines) && patchIdx < len(patchedLines) &&
			sourceLines[sourceIdx] == patchedLines[patchIdx] {
			sourceIdx++
			patchIdx++
		}

		// Add equal operation if we found equal lines
		if sourceIdx > equalStart {
			equalLines := sourceLines[equalStart:sourceIdx]
			operations = append(operations, DiffOperation{
				Type:        "equal",
				SourceStart: equalStart,
				SourceCount: len(equalLines),
				PatchStart:  equalPatchStart,
				PatchCount:  len(equalLines),
				Lines:       equalLines,
			})
		}

		// Handle differences
		if sourceIdx < len(sourceLines) || patchIdx < len(patchedLines) {
			deleteStart := sourceIdx
			insertStart := patchIdx

			// Find the next matching line or end of content
			nextMatch := findNextMatchingLine(sourceLines[sourceIdx:], patchedLines[patchIdx:])

			// Add delete operation
			if nextMatch.SourceOffset > 0 {
				deletedLines := sourceLines[sourceIdx : sourceIdx+nextMatch.SourceOffset]
				operations = append(operations, DiffOperation{
					Type:        "delete",
					SourceStart: deleteStart,
					SourceCount: len(deletedLines),
					PatchStart:  patchIdx,
					PatchCount:  0,
					Lines:       deletedLines,
				})
				sourceIdx += nextMatch.SourceOffset
			}

			// Add insert operation
			if nextMatch.PatchOffset > 0 {
				insertedLines := patchedLines[patchIdx : patchIdx+nextMatch.PatchOffset]
				operations = append(operations, DiffOperation{
					Type:        "insert",
					SourceStart: sourceIdx,
					SourceCount: 0,
					PatchStart:  insertStart,
					PatchCount:  len(insertedLines),
					Lines:       insertedLines,
				})
				patchIdx += nextMatch.PatchOffset
			}
		}
	}

	return operations
}

// MatchResult represents the result of finding the next matching line
type MatchResult struct {
	SourceOffset int
	PatchOffset  int
}

// findNextMatchingLine finds the next line that appears in both slices
func findNextMatchingLine(sourceLines, patchedLines []string) MatchResult {
	// Look for the next common line within a reasonable window
	maxLookAhead := 50 // Limit search to prevent performance issues

	for i := 0; i < len(sourceLines) && i < maxLookAhead; i++ {
		for j := 0; j < len(patchedLines) && j < maxLookAhead; j++ {
			if sourceLines[i] == patchedLines[j] {
				return MatchResult{SourceOffset: i, PatchOffset: j}
			}
		}
	}

	// No match found within window, consume all remaining lines
	return MatchResult{SourceOffset: len(sourceLines), PatchOffset: len(patchedLines)}
}

// generateHunks converts diff operations into unified diff hunks
func generateHunks(operations []DiffOperation, sourceLines, patchedLines []string) []string {
	var hunks []string

	if len(operations) == 0 {
		return hunks
	}

	// Group operations into hunks with context
	contextLines := 3
	var currentHunk strings.Builder
	var hunkSourceStart, hunkPatchStart int
	var hunkSourceCount, hunkPatchCount int
	var inHunk bool

	for i, op := range operations {
		switch op.Type {
		case "equal":
			if inHunk {
				// Add context lines after changes
				contextToAdd := min(contextLines, len(op.Lines))
				for j := 0; j < contextToAdd; j++ {
					currentHunk.WriteString(" " + op.Lines[j] + "\n")
				}
				hunkSourceCount += contextToAdd
				hunkPatchCount += contextToAdd

				// End hunk if we have enough context or this is the last operation
				if len(op.Lines) > contextLines || i == len(operations)-1 {
					hunkHeader := fmt.Sprintf("@@ -%d,%d +%d,%d @@\n",
						hunkSourceStart+1, hunkSourceCount, hunkPatchStart+1, hunkPatchCount)
					hunks = append(hunks, hunkHeader+currentHunk.String())
					currentHunk.Reset()
					inHunk = false
				}
			}

		case "delete", "insert":
			if !inHunk {
				// Start new hunk with context
				hunkSourceStart = max(0, op.SourceStart-contextLines)
				hunkPatchStart = max(0, op.PatchStart-contextLines)
				hunkSourceCount = 0
				hunkPatchCount = 0

				// Add context before changes
				contextStart := max(0, op.SourceStart-contextLines)
				contextEnd := op.SourceStart
				for j := contextStart; j < contextEnd && j < len(sourceLines); j++ {
					currentHunk.WriteString(" " + sourceLines[j] + "\n")
					hunkSourceCount++
					hunkPatchCount++
				}

				inHunk = true
			}

			// Add the actual changes
			if op.Type == "delete" {
				for _, line := range op.Lines {
					currentHunk.WriteString("-" + line + "\n")
				}
				hunkSourceCount += len(op.Lines)
			} else { // insert
				for _, line := range op.Lines {
					currentHunk.WriteString("+" + line + "\n")
				}
				hunkPatchCount += len(op.Lines)
			}
		}
	}

	// Add final hunk if still in progress
	if inHunk {
		hunkHeader := fmt.Sprintf("@@ -%d,%d +%d,%d @@\n",
			hunkSourceStart+1, hunkSourceCount, hunkPatchStart+1, hunkPatchCount)
		hunks = append(hunks, hunkHeader+currentHunk.String())
	}

	return hunks
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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
