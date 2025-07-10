package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// TwigBlock represents a Twig block found in a file
type TwigBlock struct {
	Name    string `json:"name"`
	File    string `json:"file"`
	Line    int    `json:"line"`
	Content string `json:"content"`
	Hash    string `json:"hash"`
}

// DuplicateGroup represents a group of duplicate blocks
type DuplicateGroup struct {
	BlockName string      `json:"block_name"`
	Hash      string      `json:"hash"`
	Count     int         `json:"count"`
	Files     []TwigBlock `json:"files"`
}

// BitbucketReport represents the Bitbucket Pipes report format
type BitbucketReport struct {
	Title       string                `json:"title"`
	Details     string                `json:"details"`
	Result      string                `json:"result"`
	Data        []BitbucketData       `json:"data"`
	Annotations []BitbucketAnnotation `json:"annotations,omitempty"`
}

type BitbucketData struct {
	Title string `json:"title"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type BitbucketAnnotation struct {
	Path       string `json:"path"`
	Line       int    `json:"line"`
	Message    string `json:"message"`
	Severity   string `json:"severity"`
	Type       string `json:"type"`
	ExternalID string `json:"external_id,omitempty"`
}

var (
	bitbucketFormat bool
	projectPath     string
	outputFile      string
)

var twigblocksCmd = &cobra.Command{
	Use:   "twigblocks [PATH]",
	Short: "Find duplicate Twig blocks in *.html.twig files",
	Long: `Find duplicate Twig blocks in *.html.twig files across a project.

This command scans all *.html.twig files in the specified directory (and subdirectories)
to find duplicate block definitions. This helps prevent errors in Shopware and Symfony
projects where duplicate blocks can cause template inheritance issues.

Examples:
  wswcli twigblocks .                    # Scan current directory
  wswcli twigblocks /path/to/project     # Scan specific project
  wswcli twigblocks . --bitbucket        # Output in Bitbucket format
  wswcli twigblocks . --output report.json  # Save report to file`,
	Args: cobra.MaximumNArgs(1),
	RunE: runTwigBlocks,
}

func init() {
	rootCmd.AddCommand(twigblocksCmd)
	twigblocksCmd.Flags().BoolVar(&bitbucketFormat, "bitbucket", false, "Output in Bitbucket Pipes format")
	twigblocksCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file for the report (JSON format)")
}

func runTwigBlocks(cmd *cobra.Command, args []string) error {
	// Determine project path
	if len(args) > 0 {
		projectPath = args[0]
	} else {
		projectPath = "."
	}

	// Validate project path
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return fmt.Errorf("project path does not exist: %s", projectPath)
	}

	fmt.Printf("Scanning for duplicate Twig blocks in: %s\n", projectPath)

	// Find all *.html.twig files
	twigFiles, err := findTwigFiles(projectPath)
	if err != nil {
		return fmt.Errorf("error finding Twig files: %w", err)
	}

	if len(twigFiles) == 0 {
		fmt.Println("No *.html.twig files found in the specified directory.")
		return nil
	}

	fmt.Printf("Found %d *.html.twig files\n", len(twigFiles))

	// Extract blocks from all files
	allBlocks, err := extractBlocksFromFiles(twigFiles)
	if err != nil {
		return fmt.Errorf("error extracting blocks: %w", err)
	}

	fmt.Printf("Found %d total blocks\n", len(allBlocks))

	// Find duplicates
	duplicates := findDuplicateBlocks(allBlocks)

	// Generate and output report
	if err := generateReport(duplicates, twigFiles); err != nil {
		return fmt.Errorf("error generating report: %w", err)
	}

	// Exit with error code if duplicates found (for CI/CD)
	if len(duplicates) > 0 {
		os.Exit(1)
	}

	return nil
}

// findTwigFiles recursively finds all *.html.twig files in the given directory
func findTwigFiles(rootPath string) ([]string, error) {
	var twigFiles []string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden directories and common build/cache directories
		if info.IsDir() {
			name := info.Name()
			// Skip hidden directories, but not if it's the root path we're scanning
			if path != rootPath && (strings.HasPrefix(name, ".") ||
				name == "node_modules" ||
				name == "vendor" ||
				name == "var" ||
				name == "cache" ||
				name == "build") {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if file matches *.html.twig pattern
		if strings.HasSuffix(strings.ToLower(path), ".html.twig") {
			twigFiles = append(twigFiles, path)
		}

		return nil
	})

	return twigFiles, err
}

// extractBlocksFromFiles extracts all Twig blocks from the given files
func extractBlocksFromFiles(files []string) ([]TwigBlock, error) {
	var allBlocks []TwigBlock

	// Regex to match Twig block definitions
	// Matches: {% block blockname %} or {% block blockname with optional attributes %}
	// But not inside comments {# ... #}
	blockRegex := regexp.MustCompile(`{%\s*block\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*(?:[^%]*)%}`)

	for _, file := range files {
		blocks, err := extractBlocksFromFile(file, blockRegex)
		if err != nil {
			return nil, fmt.Errorf("error processing file %s: %w", file, err)
		}
		allBlocks = append(allBlocks, blocks...)
	}

	return allBlocks, nil
}

// extractBlocksFromFile extracts Twig blocks from a single file
func extractBlocksFromFile(filename string, blockRegex *regexp.Regexp) ([]TwigBlock, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var blocks []TwigBlock
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		// Skip lines that are comments
		if strings.Contains(line, "{#") && strings.Contains(line, "#}") {
			continue
		}

		// Find all block matches in this line
		matches := blockRegex.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) >= 2 {
				blockName := match[1]

				// Extract the full block content for hashing
				blockContent := extractBlockContent(file, line, blockName)

				block := TwigBlock{
					Name:    blockName,
					File:    filename,
					Line:    lineNumber,
					Content: strings.TrimSpace(line),
					Hash:    generateContentHash(blockContent),
				}
				blocks = append(blocks, block)
			}
		}
	}

	return blocks, scanner.Err()
}

// extractBlockContent attempts to extract the full content of a block for better duplicate detection
func extractBlockContent(file *os.File, blockLine, blockName string) string {
	// For now, just use the block declaration line
	// In a more sophisticated implementation, we could parse the entire block content
	// until the matching {% endblock %}
	return strings.TrimSpace(blockLine)
}

// generateContentHash generates a hash of the block content for duplicate detection
func generateContentHash(content string) string {
	// Simple hash based on content (normalized)
	normalized := strings.ReplaceAll(strings.TrimSpace(content), " ", "")
	normalized = strings.ToLower(normalized)
	return fmt.Sprintf("%x", len(normalized)) // Simple hash for now
}

// findDuplicateBlocks identifies blocks with the same name within the same file
func findDuplicateBlocks(blocks []TwigBlock) []DuplicateGroup {
	// Group blocks by file first, then by block name within each file
	fileGroups := make(map[string][]TwigBlock)
	for _, block := range blocks {
		fileGroups[block.File] = append(fileGroups[block.File], block)
	}

	var duplicates []DuplicateGroup

	// Check each file for duplicate block names
	for fileName, fileBlocks := range fileGroups {
		// Group blocks by name within this file
		blockGroups := make(map[string][]TwigBlock)
		for _, block := range fileBlocks {
			blockGroups[block.Name] = append(blockGroups[block.Name], block)
		}

		// Check for duplicates within this file
		for blockName, groupBlocks := range blockGroups {
			if len(groupBlocks) > 1 {
				// Multiple blocks with same name in same file - this is a problem
				duplicates = append(duplicates, DuplicateGroup{
					BlockName: blockName,
					Hash:      fmt.Sprintf("file-%s", fileName),
					Count:     len(groupBlocks),
					Files:     groupBlocks,
				})
			}
		}
	}

	// Sort duplicates by block name for consistent output
	sort.Slice(duplicates, func(i, j int) bool {
		if duplicates[i].BlockName == duplicates[j].BlockName {
			return duplicates[i].Files[0].File < duplicates[j].Files[0].File
		}
		return duplicates[i].BlockName < duplicates[j].BlockName
	})

	return duplicates
}

// generateReport generates and outputs the duplicate blocks report
func generateReport(duplicates []DuplicateGroup, allFiles []string) error {
	if bitbucketFormat {
		return generateBitbucketReport(duplicates, allFiles)
	}
	return generateStandardReport(duplicates, allFiles)
}

// generateStandardReport generates a human-readable report
func generateStandardReport(duplicates []DuplicateGroup, allFiles []string) error {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("TWIG BLOCK DUPLICATE ANALYSIS REPORT")
	fmt.Println(strings.Repeat("=", 60))

	if len(duplicates) == 0 {
		fmt.Println("✅ No duplicate blocks found!")
		fmt.Printf("Scanned %d files successfully.\n", len(allFiles))
		return nil
	}

	fmt.Printf("❌ Found %d duplicate block groups:\n\n", len(duplicates))

	for i, group := range duplicates {
		fmt.Printf("%d. Block: '%s'\n", i+1, group.BlockName)
		fmt.Printf("   Hash: %s\n", group.Hash)
		fmt.Printf("   Occurrences: %d\n", group.Count)
		fmt.Println("   Files:")

		for _, block := range group.Files {
			relPath, _ := filepath.Rel(projectPath, block.File)
			fmt.Printf("     - %s:%d\n", relPath, block.Line)
			fmt.Printf("       Content: %s\n", block.Content)
		}
		fmt.Println()
	}

	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("Summary: %d duplicate groups found in %d files\n", len(duplicates), len(allFiles))
	fmt.Println("Please review and consolidate duplicate blocks to avoid template conflicts.")

	// Save to output file if specified
	if outputFile != "" {
		return saveJSONReport(duplicates, allFiles)
	}

	return nil
}

// generateBitbucketReport generates a JUnit XML compatible test report for Bitbucket Pipelines
func generateBitbucketReport(duplicates []DuplicateGroup, allFiles []string) error {
	// Create test-reports directory
	if err := os.MkdirAll("test-reports", 0755); err != nil {
		return fmt.Errorf("error creating test-reports directory: %w", err)
	}

	// Generate JUnit XML format
	junitXML := generateJUnitXML(duplicates, allFiles)

	// Write to test-reports/twig-blocks-junit.xml
	outputPath := "test-reports/twig-blocks-junit.xml"
	if err := os.WriteFile(outputPath, []byte(junitXML), 0644); err != nil {
		return fmt.Errorf("error writing JUnit XML report: %w", err)
	}

	fmt.Printf("Bitbucket test report generated: %s\n", outputPath)

	// Also output summary to stdout
	if len(duplicates) == 0 {
		fmt.Println("✅ PASSED: No duplicate Twig blocks found")
	} else {
		fmt.Printf("❌ FAILED: Found %d duplicate block groups in %d files\n", len(duplicates), len(allFiles))
	}

	return nil
}

// generateJUnitXML creates JUnit XML format for Bitbucket test reporting
func generateJUnitXML(duplicates []DuplicateGroup, allFiles []string) string {
	totalTests := len(allFiles)
	failures := len(duplicates)

	xml := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="TwigBlockDuplicateAnalysis" tests="%d" failures="%d" errors="0" time="0">
`, totalTests, failures)

	// Add test cases for each file
	fileStatus := make(map[string]bool)
	for _, group := range duplicates {
		for _, block := range group.Files {
			fileStatus[block.File] = true // Mark as failed
		}
	}

	// Generate test cases for all files
	for _, file := range allFiles {
		relPath, _ := filepath.Rel(projectPath, file)
		testName := fmt.Sprintf("TwigBlocks.%s", strings.ReplaceAll(relPath, "/", "."))

		if fileStatus[file] {
			// Failed test case
			var failureDetails []string
			for _, group := range duplicates {
				for _, block := range group.Files {
					if block.File == file {
						failureDetails = append(failureDetails,
							fmt.Sprintf("Line %d: Duplicate block '%s' (appears %d times in this file)",
								block.Line, group.BlockName, group.Count))
					}
				}
			}

			xml += fmt.Sprintf(`  <testcase classname="TwigBlocks" name="%s" time="0">
    <failure message="Duplicate Twig blocks found" type="DuplicateBlockError">
%s
    </failure>
  </testcase>
`, testName, strings.Join(failureDetails, "\n"))
		} else {
			// Passed test case
			xml += fmt.Sprintf(`  <testcase classname="TwigBlocks" name="%s" time="0"/>
`, testName)
		}
	}

	xml += "</testsuite>"
	return xml
}

// saveJSONReport saves the report in JSON format to the specified file
func saveJSONReport(duplicates []DuplicateGroup, allFiles []string) error {
	report := map[string]interface{}{
		"summary": map[string]interface{}{
			"files_scanned":    len(allFiles),
			"duplicate_groups": len(duplicates),
			"status":           map[bool]string{true: "PASSED", false: "FAILED"}[len(duplicates) == 0],
		},
		"duplicates": duplicates,
		"files":      allFiles,
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(report); err != nil {
		return fmt.Errorf("error writing JSON report: %w", err)
	}

	fmt.Printf("Report saved to: %s\n", outputFile)
	return nil
}
