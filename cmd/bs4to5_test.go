package cmd

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestFindTemplateFiles(t *testing.T) {
	// Create temporary directory structure
	tempDir := t.TempDir()

	// Create test files
	testFiles := []string{
		"template.html",
		"template.html.twig",
		"template.twig",
		"subdir/another.html",
		"subdir/another.html.twig",
		"subdir/nested/deep.twig",
		"regular.txt",               // Should be ignored
		"node_modules/ignored.html", // Should be ignored (node_modules)
		".hidden/ignored.html.twig", // Should be ignored (hidden dir)
		"vendor/ignored.twig",       // Should be ignored (vendor)
		"var/cache/ignored.html",    // Should be ignored (var)
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

	// Test findTemplateFiles with recursive=true
	found, err := findTemplateFiles(tempDir, true)
	if err != nil {
		t.Fatalf("findTemplateFiles failed: %v", err)
	}

	// Should find 6 template files: .html, .html.twig, .twig (excluding ignored ones)
	expectedCount := 6
	if len(found) != expectedCount {
		t.Errorf("Expected %d files, got %d: %v", expectedCount, len(found), found)
	}

	// Check that all found files have correct extensions
	validExtensions := []string{".html", ".html.twig", ".twig"}
	for _, file := range found {
		hasValidExtension := false
		for _, ext := range validExtensions {
			if strings.HasSuffix(strings.ToLower(file), ext) {
				hasValidExtension = true
				break
			}
		}
		if !hasValidExtension {
			t.Errorf("Found file doesn't have valid extension: %s", file)
		}
	}

	// Test findTemplateFiles with recursive=false
	found, err = findTemplateFiles(tempDir, false)
	if err != nil {
		t.Fatalf("findTemplateFiles failed: %v", err)
	}

	// Should find only root level files: 3 files
	expectedCount = 3
	if len(found) != expectedCount {
		t.Errorf("Expected %d files with recursive=false, got %d: %v", expectedCount, len(found), found)
	}
}

func TestBootstrapMigrations(t *testing.T) {
	migrations := getBootstrapMigrations()

	// Test that we have migrations defined
	if len(migrations) == 0 {
		t.Fatal("No bootstrap migrations found")
	}

	// Test specific migrations
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Custom Checkbox",
			input:    `<div class="custom-control custom-checkbox">`,
			expected: `<div class="form-check">`,
		},
		{
			name:     "Custom Radio",
			input:    `<div class="custom-control custom-radio">`,
			expected: `<div class="form-check">`,
		},
		{
			name:     "Custom Switch",
			input:    `<div class="custom-control custom-switch">`,
			expected: `<div class="form-check form-switch">`,
		},
		{
			name:     "Custom Control Input",
			input:    `<input class="custom-control-input" type="checkbox">`,
			expected: `<input class="form-check-input" type="checkbox">`,
		},
		{
			name:     "Custom Control Label",
			input:    `<label class="custom-control-label">`,
			expected: `<label class="form-check-label">`,
		},
		{
			name:     "Custom Select",
			input:    `<select class="custom-select">`,
			expected: `<select class="form-select">`,
		},
		{
			name:     "No Gutters",
			input:    `<div class="row no-gutters">`,
			expected: `<div class="row g-0">`,
		},
		{
			name:     "Button Block",
			input:    `<button class="btn btn-primary btn-block">`,
			expected: `<button class="btn btn-primary d-grid">`,
		},
		{
			name:     "Badge Pill",
			input:    `<span class="badge badge-primary badge-pill">`,
			expected: `<span class="badge bg-primary rounded-pill">`,
		},
		{
			name:     "Close Button",
			input:    `<button class="close">`,
			expected: `<button class="btn-close">`,
		},
		{
			name:     "Text Left",
			input:    `<div class="text-left">`,
			expected: `<div class="text-start">`,
		},
		{
			name:     "Text Right",
			input:    `<div class="text-right">`,
			expected: `<div class="text-end">`,
		},
		{
			name:     "Float Left",
			input:    `<div class="float-left">`,
			expected: `<div class="float-start">`,
		},
		{
			name:     "Float Right",
			input:    `<div class="float-right">`,
			expected: `<div class="float-end">`,
		},
		{
			name:     "Margin Left",
			input:    `<div class="ml-3">`,
			expected: `<div class="ms-3">`,
		},
		{
			name:     "Margin Right",
			input:    `<div class="mr-2">`,
			expected: `<div class="me-2">`,
		},
		{
			name:     "Padding Left",
			input:    `<div class="pl-4">`,
			expected: `<div class="ps-4">`,
		},
		{
			name:     "Padding Right",
			input:    `<div class="pr-1">`,
			expected: `<div class="pe-1">`,
		},
		{
			name:     "Form Group",
			input:    `<div class="form-group">`,
			expected: `<div class="mb-3">`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input
			changeCount := 0

			// Apply all migrations to see if our test case gets transformed
			for _, migration := range migrations {
				matches := migration.Pattern.FindAllString(result, -1)
				if len(matches) > 0 {
					newResult := migration.Pattern.ReplaceAllString(result, migration.Replacement)
					if newResult != result {
						changeCount += len(matches)
						result = newResult
					}
				}
			}

			if result != tc.expected {
				t.Errorf("Migration failed for %s:\nInput:    %s\nExpected: %s\nActual:   %s", tc.name, tc.input, tc.expected, result)
			}

			if changeCount == 0 {
				t.Errorf("No changes applied for test case: %s", tc.name)
			}
		})
	}
}

func TestProcessFileWithDryRun(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.html")

	// Create test content with Bootstrap 4 classes
	content := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
    <div class="custom-control custom-checkbox">
        <input type="checkbox" class="custom-control-input" id="check1">
        <label class="custom-control-label" for="check1">Check me</label>
    </div>
    <div class="row no-gutters">
        <div class="col text-left">
            <button class="btn btn-primary btn-block">Button</button>
        </div>
    </div>
</body>
</html>`

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	migrations := getBootstrapMigrations()

	// Test with dry-run (should not modify file)
	originalDryRun := dryRun
	dryRun = true
	defer func() { dryRun = originalDryRun }()

	changes, err := processFile(testFile, migrations)
	if err != nil {
		t.Fatalf("processFile failed: %v", err)
	}

	// Should detect changes
	if changes == 0 {
		t.Error("Expected to detect changes in file")
	}

	// File content should remain unchanged in dry-run mode
	afterContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(afterContent) != content {
		t.Error("File content should not change in dry-run mode")
	}
}

func TestProcessFileWithActualChanges(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.html")

	// Create test content with Bootstrap 4 classes
	originalContent := `<div class="custom-control custom-checkbox">
    <input type="checkbox" class="custom-control-input" id="check1">
    <label class="custom-control-label" for="check1">Check me</label>
</div>
<div class="row no-gutters text-left">
    <button class="btn btn-primary btn-block">Button</button>
</div>`

	if err := os.WriteFile(testFile, []byte(originalContent), 0644); err != nil {
		t.Fatal(err)
	}

	migrations := getBootstrapMigrations()

	// Test with actual changes (not dry-run)
	originalDryRun := dryRun
	dryRun = false
	defer func() { dryRun = originalDryRun }()

	changes, err := processFile(testFile, migrations)
	if err != nil {
		t.Fatalf("processFile failed: %v", err)
	}

	// Should make changes
	if changes == 0 {
		t.Error("Expected to make changes to file")
	}

	// File content should be modified
	afterContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	afterString := string(afterContent)
	if afterString == originalContent {
		t.Error("File content should be changed when not in dry-run mode")
	}

	// Check specific transformations
	expectedTransformations := map[string]string{
		"custom-control custom-checkbox": "form-check",
		"custom-control-input":           "form-check-input",
		"custom-control-label":           "form-check-label",
		"no-gutters":                     "g-0",
		"text-left":                      "text-start",
		"btn-block":                      "d-grid",
	}

	for oldClass, newClass := range expectedTransformations {
		if strings.Contains(afterString, oldClass) {
			t.Errorf("Old class '%s' should have been replaced", oldClass)
		}
		if !strings.Contains(afterString, newClass) {
			t.Errorf("New class '%s' should be present in output", newClass)
		}
	}
}

func TestTwigFileProcessing(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.html.twig")

	// Create Twig template with Bootstrap 4 classes
	content := `{% extends 'base.html.twig' %}

{% block content %}
    <div class="container">
        <div class="custom-control custom-switch">
            <input type="checkbox" class="custom-control-input" id="switch1">
            <label class="custom-control-label" for="switch1">{{ 'toggle.label'|trans }}</label>
        </div>
        
        {% for item in items %}
            <div class="row no-gutters">
                <div class="col text-right mr-3">
                    <span class="badge badge-success badge-pill">{{ item.status }}</span>
                </div>
            </div>
        {% endfor %}
        
        <select class="custom-select">
            {% for option in options %}
                <option value="{{ option.value }}">{{ option.label }}</option>
            {% endfor %}
        </select>
    </div>
{% endblock %}`

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	migrations := getBootstrapMigrations()

	// Test processing Twig file
	originalDryRun := dryRun
	dryRun = false
	defer func() { dryRun = originalDryRun }()

	changes, err := processFile(testFile, migrations)
	if err != nil {
		t.Fatalf("processFile failed for Twig file: %v", err)
	}

	if changes == 0 {
		t.Error("Expected to make changes to Twig file")
	}

	// Check that Twig syntax is preserved
	afterContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	afterString := string(afterContent)

	// Verify Twig syntax is intact
	twigSyntax := []string{
		"{% extends",
		"{% block content %}",
		"{% endblock %}",
		"{{ 'toggle.label'|trans }}",
		"{% for item in items %}",
		"{% endfor %}",
		"{{ item.status }}",
		"{{ option.value }}",
		"{{ option.label }}",
	}

	for _, syntax := range twigSyntax {
		if !strings.Contains(afterString, syntax) {
			t.Errorf("Twig syntax should be preserved: %s", syntax)
		}
	}

	// Verify Bootstrap classes were transformed
	expectedTransformations := []string{
		"form-check form-switch",
		"form-check-input",
		"form-check-label",
		"g-0",
		"text-end",
		"me-3",
		"bg-success",
		"rounded-pill",
		"form-select",
	}

	for _, newClass := range expectedTransformations {
		if !strings.Contains(afterString, newClass) {
			t.Errorf("Expected Bootstrap 5 class should be present: %s", newClass)
		}
	}
}

func TestComplexBootstrapMigration(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "complex.html")

	// Create complex HTML with multiple Bootstrap 4 elements
	content := `<!DOCTYPE html>
<html>
<body>
    <div class="container">
        <!-- Form Controls -->
        <div class="custom-control custom-checkbox mb-2">
            <input type="checkbox" class="custom-control-input" id="check1">
            <label class="custom-control-label" for="check1">Checkbox</label>
        </div>
        <div class="custom-control custom-radio mb-2">
            <input type="radio" class="custom-control-input" id="radio1" name="radio">
            <label class="custom-control-label" for="radio1">Radio</label>
        </div>
        <div class="custom-control custom-switch mb-3">
            <input type="checkbox" class="custom-control-input" id="switch1">
            <label class="custom-control-label" for="switch1">Switch</label>
        </div>
        
        <!-- Grid and Layout -->
        <div class="row no-gutters">
            <div class="col-6 text-left pl-3">
                <h3>Left Column</h3>
            </div>
            <div class="col-6 text-right pr-3">
                <h3>Right Column</h3>
            </div>
        </div>
        
        <!-- Buttons and Badges -->
        <button class="btn btn-primary btn-block mt-3">Block Button</button>
        <div class="mt-2">
            <span class="badge badge-primary badge-pill mr-2">Primary</span>
            <span class="badge badge-success mr-2">Success</span>
            <span class="badge badge-warning">Warning</span>
        </div>
        
        <!-- Form Elements -->
        <div class="form-group mt-3">
            <label>Select Options</label>
            <select class="custom-select">
                <option>Choose...</option>
                <option value="1">Option 1</option>
            </select>
        </div>
        
        <!-- Float and Close -->
        <div class="float-left ml-2">
            <p>Floated content</p>
            <button type="button" class="close" aria-label="Close">
                <span aria-hidden="true">&times;</span>
            </button>
        </div>
    </div>
</body>
</html>`

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	migrations := getBootstrapMigrations()

	// Process the file
	originalDryRun := dryRun
	dryRun = false
	defer func() { dryRun = originalDryRun }()

	changes, err := processFile(testFile, migrations)
	if err != nil {
		t.Fatalf("processFile failed: %v", err)
	}

	// Should make many changes
	expectedMinChanges := 15
	if changes < expectedMinChanges {
		t.Errorf("Expected at least %d changes, got %d", expectedMinChanges, changes)
	}

	// Read result and verify transformations
	afterContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	afterString := string(afterContent)

	// Verify Bootstrap 5 classes are present
	bs5Classes := []string{
		"form-check",
		"form-check-input",
		"form-check-label",
		"form-switch",
		"g-0",
		"text-start",
		"text-end",
		"ps-3",
		"pe-3",
		"d-grid",
		"bg-primary",
		"bg-success",
		"bg-warning",
		"rounded-pill",
		"me-2",
		"mb-3",
		"form-select",
		"float-start",
		"ms-2",
		"btn-close",
	}

	for _, class := range bs5Classes {
		if !strings.Contains(afterString, class) {
			t.Errorf("Expected Bootstrap 5 class not found: %s", class)
		}
	}

	// Verify Bootstrap 4 classes are removed
	bs4Classes := []string{
		"custom-control custom-checkbox",
		"custom-control custom-radio",
		"custom-control custom-switch",
		"custom-control-input",
		"custom-control-label",
		"no-gutters",
		"text-left",
		"text-right",
		"pl-3",
		"pr-3",
		"btn-block",
		"badge-primary",
		"badge-success",
		"badge-warning",
		"badge-pill",
		"mr-2",
		"form-group",
		"custom-select",
		"float-left",
		"ml-2",
	}

	for _, class := range bs4Classes {
		// Use word boundaries to avoid false positives (e.g., "close" in "btn-close")
		pattern := `\b` + regexp.QuoteMeta(class) + `\b`
		matched, _ := regexp.MatchString(pattern, afterString)
		if matched {
			t.Errorf("Bootstrap 4 class should have been replaced: %s", class)
		}
	}
}

func TestRegexPatterns(t *testing.T) {
	// Test that our regex patterns work correctly
	testCases := []struct {
		pattern     string
		input       string
		shouldMatch bool
	}{
		{
			pattern:     `custom-control\s+custom-checkbox`,
			input:       `class="custom-control custom-checkbox"`,
			shouldMatch: true,
		},
		{
			pattern:     `custom-control\s+custom-checkbox`,
			input:       `class="custom-control  custom-checkbox"`, // Multiple spaces
			shouldMatch: true,
		},
		{
			pattern:     `custom-control\s+custom-checkbox`,
			input:       `class="custom-checkbox"`, // Missing custom-control
			shouldMatch: false,
		},
		{
			pattern:     `custom-control-input`,
			input:       `class="custom-control-input"`,
			shouldMatch: true,
		},
		{
			pattern:     `ml-(\d+)`,
			input:       `class="ml-3"`,
			shouldMatch: true,
		},
		{
			pattern:     `ml-(\d+)`,
			input:       `class="ml-auto"`, // Should not match non-numeric
			shouldMatch: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.pattern+"_"+tc.input, func(t *testing.T) {
			regex := regexp.MustCompile(tc.pattern)
			matches := regex.MatchString(tc.input)

			if matches != tc.shouldMatch {
				t.Errorf("Pattern %s with input %s: expected match=%v, got match=%v",
					tc.pattern, tc.input, tc.shouldMatch, matches)
			}
		})
	}
}

func TestTwigCommentsInTwigFile(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.html.twig")

	// Create Twig file with migrations that have Twig comments (use input-group-append which fully replaces with comment)
	content := `{% extends 'base.html.twig' %}

{% block content %}
    <div class="container">
        <div class="input-group-append">
            <span class="input-group-text">Button</span>
        </div>
    </div>
{% endblock %}`

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	migrations := getBootstrapMigrations()

	// Process the Twig file
	originalDryRun := dryRun
	dryRun = false
	defer func() { dryRun = originalDryRun }()

	changes, err := processFile(testFile, migrations)
	if err != nil {
		t.Fatalf("processFile failed for Twig file: %v", err)
	}

	if changes == 0 {
		t.Error("Expected to make changes to Twig file")
	}

	// Read result
	afterContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	afterString := string(afterContent)

	// Should contain Twig comments {# #}, not HTML comments <!-- -->
	if !strings.Contains(afterString, "{# TODO:") {
		t.Errorf("Twig file should contain Twig comment format {# #}, got: %s", afterString)
	}

	if strings.Contains(afterString, "<!-- TODO:") {
		t.Error("Twig file should not contain HTML comment format <!-- -->")
	}
}

func TestTwigCommentsInHTMLFile(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.html")

	// Create HTML file with migrations that have Twig comments (use input-group-append which fully replaces with comment)
	content := `<!DOCTYPE html>
<html>
<body>
    <div class="container">
        <div class="input-group-append">
            <span class="input-group-text">Button</span>
        </div>
    </div>
</body>
</html>`

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	migrations := getBootstrapMigrations()

	// Process the HTML file
	originalDryRun := dryRun
	dryRun = false
	defer func() { dryRun = originalDryRun }()

	changes, err := processFile(testFile, migrations)
	if err != nil {
		t.Fatalf("processFile failed for HTML file: %v", err)
	}

	if changes == 0 {
		t.Error("Expected to make changes to HTML file")
	}

	// Read result
	afterContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}

	afterString := string(afterContent)

	// Should contain HTML comments <!-- -->, not Twig comments {# #}
	if !strings.Contains(afterString, "<!--  TODO:") {
		t.Errorf("HTML file should contain HTML comment format <!--  TODO:, got: %s", afterString)
	}

	if strings.Contains(afterString, "{# TODO:") {
		t.Error("HTML file should not contain Twig comment format {# #}")
	}
}

func TestIntegrationFullWorkflow(t *testing.T) {
	// Test the complete workflow from finding files to applying migrations
	tempDir := t.TempDir()

	// Create a realistic project structure
	files := map[string]string{
		"templates/layout/base.html.twig": `<!DOCTYPE html>
<html>
<head><title>{% block title %}{% endblock %}</title></head>
<body class="bg-light">
    {% block content %}{% endblock %}
</body>
</html>`,
		"templates/forms/contact.html.twig": `{% extends 'layout/base.html.twig' %}
{% block content %}
<div class="container">
    <div class="row no-gutters">
        <div class="col-md-6 offset-md-3">
            <form class="form-group">
                <div class="custom-control custom-checkbox">
                    <input type="checkbox" class="custom-control-input" id="newsletter">
                    <label class="custom-control-label" for="newsletter">Subscribe</label>
                </div>
                <button type="submit" class="btn btn-primary btn-block mt-3">Submit</button>
            </form>
        </div>
    </div>
</div>
{% endblock %}`,
		"public/index.html": `<!DOCTYPE html>
<html>
<body>
    <div class="text-left">
        <span class="badge badge-success badge-pill">Welcome</span>
        <button class="close ml-2" type="button">&times;</button>
    </div>
</body>
</html>`,
		"components/card.html": `<div class="card">
    <div class="card-body text-right pr-3">
        <select class="custom-select">
            <option>Choose...</option>
        </select>
    </div>
</div>`,
	}

	// Create files
	for filePath, content := range files {
		fullPath := filepath.Join(tempDir, filePath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Test finding template files
	templateFiles, err := findTemplateFiles(tempDir, true)
	if err != nil {
		t.Fatalf("findTemplateFiles failed: %v", err)
	}

	expectedFileCount := 4 // All files should be found
	if len(templateFiles) != expectedFileCount {
		t.Errorf("Expected %d template files, got %d", expectedFileCount, len(templateFiles))
	}

	// Test processing all files
	migrations := getBootstrapMigrations()
	totalChanges := 0

	originalDryRun := dryRun
	dryRun = false
	defer func() { dryRun = originalDryRun }()

	for _, file := range templateFiles {
		changes, err := processFile(file, migrations)
		if err != nil {
			t.Errorf("Error processing %s: %v", file, err)
			continue
		}
		totalChanges += changes
	}

	// Should make changes across files
	if totalChanges == 0 {
		t.Error("Expected to make changes across template files")
	}

	// Verify specific files were transformed correctly
	// Check contact.html.twig
	contactFile := filepath.Join(tempDir, "templates/forms/contact.html.twig")
	contactContent, err := os.ReadFile(contactFile)
	if err != nil {
		t.Fatal(err)
	}

	contactString := string(contactContent)

	// Should have Bootstrap 5 classes
	requiredClasses := []string{"form-check", "form-check-input", "form-check-label", "g-0", "d-grid", "mb-3"}
	for _, class := range requiredClasses {
		if !strings.Contains(contactString, class) {
			t.Errorf("Contact template should contain Bootstrap 5 class: %s", class)
		}
	}

	// Should not have Bootstrap 4 classes
	forbiddenClasses := []string{"custom-control", "custom-control-input", "custom-control-label", "no-gutters", "btn-block", "form-group"}
	for _, class := range forbiddenClasses {
		if strings.Contains(contactString, class) {
			t.Errorf("Contact template should not contain Bootstrap 4 class: %s", class)
		}
	}

	// Verify Twig syntax is preserved
	twigElements := []string{"{% extends", "{% block content %}", "{% endblock %}"}
	for _, element := range twigElements {
		if !strings.Contains(contactString, element) {
			t.Errorf("Twig syntax should be preserved: %s", element)
		}
	}
}
