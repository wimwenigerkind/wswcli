package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestFindTwigFiles(t *testing.T) {
	// Create temporary directory structure
	tempDir := t.TempDir()

	// Create test files
	testFiles := []string{
		"template.html.twig",
		"subdir/another.html.twig",
		"subdir/nested/deep.html.twig",
		"regular.html",                   // Should be ignored
		"template.twig",                  // Should be ignored (not .html.twig)
		"node_modules/ignored.html.twig", // Should be ignored (node_modules)
		".hidden/ignored.html.twig",      // Should be ignored (hidden dir)
	}

	for _, file := range testFiles {
		fullPath := filepath.Join(tempDir, file)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte("test content"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Test findTwigFiles
	found, err := findTwigFiles(tempDir)
	if err != nil {
		t.Fatalf("findTwigFiles failed: %v", err)
	}

	// Should find exactly 3 .html.twig files (excluding ignored ones)
	if len(found) != 3 {
		t.Errorf("Expected 3 files, got %d: %v", len(found), found)
	}

	// Check that all found files end with .html.twig
	for _, file := range found {
		if !strings.HasSuffix(file, ".html.twig") {
			t.Errorf("Found file doesn't end with .html.twig: %s", file)
		}
	}
}

func TestExtractBlocksFromFile(t *testing.T) {
	// Create temporary file with Twig blocks
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.html.twig")

	content := `
{% extends "base.html.twig" %}

{% block title %}Page Title{% endblock %}

{% block content %}
    <div class="content">
        {% block inner_content %}
            Default inner content
        {% endblock %}
    </div>
{% endblock %}

{% block sidebar %}
    Sidebar content
{% endblock %}

{# This is a comment with {% block fake_block %} - should be ignored #}
`

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Test block extraction
	blockRegex := regexp.MustCompile(`{%\s*block\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*(?:[^%]*)%}`)
	blocks, err := extractBlocksFromFile(testFile, blockRegex)
	if err != nil {
		t.Fatalf("extractBlocksFromFile failed: %v", err)
	}

	// Should find 4 blocks: title, content, inner_content, sidebar
	expectedBlocks := []string{"title", "content", "inner_content", "sidebar"}
	if len(blocks) != len(expectedBlocks) {
		t.Errorf("Expected %d blocks, got %d. Found blocks: %v", len(expectedBlocks), len(blocks), getBlockNames(blocks))
	}

	// Check block names
	foundNames := make(map[string]bool)
	for _, block := range blocks {
		foundNames[block.Name] = true

		// Verify file path is correct
		if block.File != testFile {
			t.Errorf("Block file path incorrect: expected %s, got %s", testFile, block.File)
		}

		// Verify line number is positive
		if block.Line <= 0 {
			t.Errorf("Block line number should be positive, got %d", block.Line)
		}
	}

	for _, expected := range expectedBlocks {
		if !foundNames[expected] {
			t.Errorf("Expected block '%s' not found", expected)
		}
	}
}

func TestFindDuplicateBlocks(t *testing.T) {
	// Create test blocks with duplicates within same file
	blocks := []TwigBlock{
		{Name: "content", File: "file1.twig", Line: 5, Content: "{% block content %}", Hash: "hash1"},
		{Name: "content", File: "file1.twig", Line: 15, Content: "{% block content %}", Hash: "hash1"}, // Duplicate in same file
		{Name: "sidebar", File: "file1.twig", Line: 10, Content: "{% block sidebar %}", Hash: "hash2"},
		{Name: "title", File: "file2.twig", Line: 1, Content: "{% block title %}", Hash: "hash3"},
		{Name: "content", File: "file2.twig", Line: 8, Content: "{% block content %}", Hash: "hash4"}, // Same name, different file - should NOT be duplicate
		{Name: "footer", File: "file3.twig", Line: 5, Content: "{% block footer %}", Hash: "hash5"},
		{Name: "footer", File: "file3.twig", Line: 10, Content: "{% block footer %}", Hash: "hash5"}, // Another duplicate in same file
	}

	duplicates := findDuplicateBlocks(blocks)

	// Should find duplicates only within same files
	if len(duplicates) != 2 {
		t.Errorf("Expected 2 duplicate groups (content in file1, footer in file3), got %d", len(duplicates))
	}

	// Check that content block duplicates are found in file1
	var contentDuplicates *DuplicateGroup
	var footerDuplicates *DuplicateGroup
	for i := range duplicates {
		if duplicates[i].BlockName == "content" {
			contentDuplicates = &duplicates[i]
		}
		if duplicates[i].BlockName == "footer" {
			footerDuplicates = &duplicates[i]
		}
	}

	if contentDuplicates == nil {
		t.Error("Expected to find content block duplicates in file1")
	} else {
		if contentDuplicates.Count != 2 {
			t.Errorf("Expected 2 duplicate content blocks in file1, got %d", contentDuplicates.Count)
		}
		// Verify all blocks are from the same file
		for _, block := range contentDuplicates.Files {
			if block.File != "file1.twig" {
				t.Errorf("Expected all duplicate blocks to be from file1.twig, got %s", block.File)
			}
		}
	}

	if footerDuplicates == nil {
		t.Error("Expected to find footer block duplicates in file3")
	} else {
		if footerDuplicates.Count != 2 {
			t.Errorf("Expected 2 duplicate footer blocks in file3, got %d", footerDuplicates.Count)
		}
		// Verify all blocks are from the same file
		for _, block := range footerDuplicates.Files {
			if block.File != "file3.twig" {
				t.Errorf("Expected all duplicate blocks to be from file3.twig, got %s", block.File)
			}
		}
	}
}

func TestGenerateContentHash(t *testing.T) {
	// Test that identical content produces same hash
	content1 := "{% block content %}"
	content2 := "{% block content %}"
	content3 := "{% block content with class %}"

	hash1 := generateContentHash(content1)
	hash2 := generateContentHash(content2)
	hash3 := generateContentHash(content3)

	if hash1 != hash2 {
		t.Error("Identical content should produce same hash")
	}

	if hash1 == hash3 {
		t.Error("Different content should produce different hash")
	}
}

func TestBitbucketReportGeneration(t *testing.T) {
	// Create test duplicates
	duplicates := []DuplicateGroup{
		{
			BlockName: "content",
			Hash:      "hash1",
			Count:     2,
			Files: []TwigBlock{
				{Name: "content", File: "file1.twig", Line: 5, Content: "{% block content %}"},
				{Name: "content", File: "file2.twig", Line: 3, Content: "{% block content %}"},
			},
		},
	}

	allFiles := []string{"file1.twig", "file2.twig"}

	// Set bitbucket format and project path for testing
	originalBitbucket := bitbucketFormat
	originalProjectPath := projectPath
	defer func() {
		bitbucketFormat = originalBitbucket
		projectPath = originalProjectPath
	}()

	bitbucketFormat = true
	projectPath = "."

	// Capture output by redirecting stdout
	// Note: In a real test, you might want to use a more sophisticated approach
	// to capture and verify the JSON output

	// For now, just test that the function doesn't panic
	err := generateBitbucketReport(duplicates, allFiles)
	if err != nil {
		t.Errorf("generateBitbucketReport failed: %v", err)
	}
}

func TestSaveJSONReport(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "report.json")

	// Set output file for testing
	originalOutputFile := outputFile
	defer func() {
		outputFile = originalOutputFile
	}()
	outputFile = outputPath

	duplicates := []DuplicateGroup{
		{
			BlockName: "content",
			Hash:      "hash1",
			Count:     2,
			Files: []TwigBlock{
				{Name: "content", File: "file1.twig", Line: 5, Content: "{% block content %}"},
				{Name: "content", File: "file2.twig", Line: 3, Content: "{% block content %}"},
			},
		},
	}

	allFiles := []string{"file1.twig", "file2.twig"}

	err := saveJSONReport(duplicates, allFiles)
	if err != nil {
		t.Fatalf("saveJSONReport failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}

	// Verify JSON content
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	var report map[string]interface{}
	if err := json.Unmarshal(content, &report); err != nil {
		t.Fatalf("Failed to parse JSON report: %v", err)
	}

	// Check report structure
	if summary, ok := report["summary"].(map[string]interface{}); ok {
		if filesScanned, ok := summary["files_scanned"].(float64); !ok || int(filesScanned) != 2 {
			t.Error("Incorrect files_scanned in report")
		}
		if duplicateGroups, ok := summary["duplicate_groups"].(float64); !ok || int(duplicateGroups) != 1 {
			t.Error("Incorrect duplicate_groups in report")
		}
	} else {
		t.Error("Missing or invalid summary in report")
	}
}

func TestIntegrationWithRealFiles(t *testing.T) {
	// Create a temporary project structure with real Twig files
	tempDir := t.TempDir()

	// Create test files with duplicate blocks within same file
	files := map[string]string{
		"templates/base.html.twig": `
{% block title %}Default Title{% endblock %}
{% block content %}{% endblock %}
`,
		"templates/page.html.twig": `
{% extends "base.html.twig" %}
{% block title %}Page Title{% endblock %}
{% block content %}Page content{% endblock %}
`,
		"templates/duplicate.html.twig": `
{% block title %}Default Title{% endblock %}
{% block title %}Duplicate Title{% endblock %}
{% block sidebar %}Sidebar{% endblock %}
`,
	}

	for filePath, content := range files {
		fullPath := filepath.Join(tempDir, filePath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Test the full workflow
	twigFiles, err := findTwigFiles(tempDir)
	if err != nil {
		t.Fatalf("findTwigFiles failed: %v", err)
	}

	if len(twigFiles) != 3 {
		t.Errorf("Expected 3 Twig files, got %d", len(twigFiles))
	}

	allBlocks, err := extractBlocksFromFiles(twigFiles)
	if err != nil {
		t.Fatalf("extractBlocksFromFiles failed: %v", err)
	}

	duplicates := findDuplicateBlocks(allBlocks)

	// Should find duplicates for 'title' block within duplicate.html.twig
	if len(duplicates) != 1 {
		t.Errorf("Expected 1 duplicate group, got %d", len(duplicates))
	}

	if len(duplicates) > 0 {
		titleDuplicates := &duplicates[0]
		if titleDuplicates.BlockName != "title" {
			t.Errorf("Expected duplicate block name 'title', got '%s'", titleDuplicates.BlockName)
		}
		if titleDuplicates.Count != 2 {
			t.Errorf("Expected 2 duplicate title blocks, got %d", titleDuplicates.Count)
		}
		// Verify all blocks are from the same file
		for _, block := range titleDuplicates.Files {
			if !strings.Contains(block.File, "duplicate.html.twig") {
				t.Errorf("Expected all duplicate blocks to be from duplicate.html.twig, got %s", block.File)
			}
		}
	}
}

// Helper function to get block names for debugging
func getBlockNames(blocks []TwigBlock) []string {
	var names []string
	for _, block := range blocks {
		names = append(names, block.Name)
	}
	return names
}
