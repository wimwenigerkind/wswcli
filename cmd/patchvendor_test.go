package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateInputs(t *testing.T) {
	// Create temporary test files
	tempDir := t.TempDir()
	sourceFile := filepath.Join(tempDir, "source.php")
	patchedFile := filepath.Join(tempDir, "patched.php")
	outputFile := filepath.Join(tempDir, "output.patch")

	// Create test files
	if err := os.WriteFile(sourceFile, []byte("<?php echo 'test';"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(patchedFile, []byte("<?php echo 'modified';"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		sourcePath  string
		patchedPath string
		outputPath  string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid inputs",
			sourcePath:  sourceFile,
			patchedPath: patchedFile,
			outputPath:  outputFile,
			expectError: false,
		},
		{
			name:        "Source file does not exist",
			sourcePath:  filepath.Join(tempDir, "nonexistent.php"),
			patchedPath: patchedFile,
			outputPath:  outputFile,
			expectError: true,
			errorMsg:    "source path does not exist",
		},
		{
			name:        "Patched file does not exist",
			sourcePath:  sourceFile,
			patchedPath: filepath.Join(tempDir, "nonexistent.php"),
			outputPath:  outputFile,
			expectError: true,
			errorMsg:    "patched path does not exist",
		},
		{
			name:        "Same source and patched paths",
			sourcePath:  sourceFile,
			patchedPath: sourceFile,
			outputPath:  outputFile,
			expectError: true,
			errorMsg:    "source and patched paths cannot be the same",
		},
		{
			name:        "Empty output path",
			sourcePath:  sourceFile,
			patchedPath: patchedFile,
			outputPath:  "",
			expectError: true,
			errorMsg:    "output path cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateInputs(tt.sourcePath, tt.patchedPath, tt.outputPath)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %s", err.Error())
				}
			}
		})
	}
}

func TestValidateFileExtensions(t *testing.T) {
	tests := []struct {
		name        string
		sourcePath  string
		patchedPath string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Same PHP extensions",
			sourcePath:  "test.php",
			patchedPath: "test2.php",
			expectError: false,
		},
		{
			name:        "Same JS extensions",
			sourcePath:  "test.js",
			patchedPath: "test2.js",
			expectError: false,
		},
		{
			name:        "Different extensions",
			sourcePath:  "test.php",
			patchedPath: "test.js",
			expectError: true,
			errorMsg:    "different extensions",
		},
		{
			name:        "No extensions",
			sourcePath:  "test",
			patchedPath: "test2",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFileExtensions(tt.sourcePath, tt.patchedPath)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %s", err.Error())
				}
			}
		})
	}
}

func TestTrimVendorPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Shopware core file",
			input:    "vendor/shopware/core/Framework/Plugin/PluginManager.php",
			expected: "shopware/core",
		},
		{
			name:     "Symfony component",
			input:    "vendor/symfony/console/Command/Command.php",
			expected: "symfony/console",
		},
		{
			name:     "Deep nested file",
			input:    "vendor/doctrine/orm/lib/Doctrine/ORM/EntityManager.php",
			expected: "doctrine/orm",
		},
		{
			name:     "Directory without vendor",
			input:    "src/Controller/TestController.php",
			expected: "src/Controller/TestController.php",
		},
		{
			name:     "Simple filename",
			input:    "test.php",
			expected: "test.php",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := trimVendorPath(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestExtractVendorPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Shopware core file",
			input:    "vendor/shopware/core/Framework/Plugin/PluginManager.php",
			expected: "vendor/shopware/core/Framework/Plugin/PluginManager.php",
		},
		{
			name:     "Symfony component",
			input:    "vendor/symfony/console/Command/Command.php",
			expected: "vendor/symfony/console/Command/Command.php",
		},
		{
			name:     "Directory without vendor",
			input:    "src/Controller/TestController.php",
			expected: "src/Controller/TestController.php",
		},
		{
			name:     "Simple filename",
			input:    "test.php",
			expected: "test.php",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractVendorPath(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestGenerateUnifiedDiff(t *testing.T) {
	tests := []struct {
		name           string
		sourcePath     string
		sourceContent  string
		patchedContent string
		expectContains []string
	}{
		{
			name:       "Simple line removal",
			sourcePath: "vendor/shopware/core/Test.php",
			sourceContent: `<?php
class Test {
    /**
     * Old comment
     */
    public function test() {
        return true;
    }
}`,
			patchedContent: `<?php
class Test {
    public function test() {
        return true;
    }
}`,
			expectContains: []string{
				"--- a/vendor/shopware/core",
				"+++ b/vendor/shopware/core",
				"-    /**",
				"-     * Old comment",
				"-     */",
			},
		},
		{
			name:       "Line addition",
			sourcePath: "vendor/symfony/console/Command.php",
			sourceContent: `<?php
class Command {
    public function execute() {
        return 0;
    }
}`,
			patchedContent: `<?php
class Command {
    public function execute() {
        $this->validate();
        return 0;
    }
}`,
			expectContains: []string{
				"--- a/vendor/symfony/console",
				"+++ b/vendor/symfony/console",
				"+        $this->validate();",
			},
		},
		{
			name:       "Multiple changes",
			sourcePath: "vendor/doctrine/orm/Entity.php",
			sourceContent: `<?php
class Entity {
    private $id;
    
    public function getId() {
        return $this->id;
    }
    
    public function setId($id) {
        $this->id = $id;
    }
}`,
			patchedContent: `<?php
class Entity {
    private $id;
    private $name;
    
    public function getId() {
        return $this->id;
    }
    
    public function getName() {
        return $this->name;
    }
    
    public function setId($id) {
        $this->id = $id;
    }
}`,
			expectContains: []string{
				"--- a/vendor/doctrine/orm",
				"+++ b/vendor/doctrine/orm",
				"+    private $name;",
				"+    public function getName() {",
				"+        return $this->name;",
				"+    }",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateUnifiedDiff(tt.sourcePath, "", tt.sourceContent, tt.patchedContent)

			for _, expected := range tt.expectContains {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected diff to contain '%s', but it didn't.\nFull diff:\n%s", expected, result)
				}
			}
		})
	}
}

func TestComputeLineDiff(t *testing.T) {
	tests := []struct {
		name         string
		sourceLines  []string
		patchedLines []string
		expectOps    int // Expected number of operations
	}{
		{
			name:         "Identical content",
			sourceLines:  []string{"line1", "line2", "line3"},
			patchedLines: []string{"line1", "line2", "line3"},
			expectOps:    1, // One equal operation
		},
		{
			name:         "Complete replacement",
			sourceLines:  []string{"old1", "old2"},
			patchedLines: []string{"new1", "new2"},
			expectOps:    2, // Delete and insert operations
		},
		{
			name:         "Line addition",
			sourceLines:  []string{"line1", "line3"},
			patchedLines: []string{"line1", "line2", "line3"},
			expectOps:    3, // Equal, insert, equal
		},
		{
			name:         "Line removal",
			sourceLines:  []string{"line1", "line2", "line3"},
			patchedLines: []string{"line1", "line3"},
			expectOps:    3, // Equal, delete, equal
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operations := computeLineDiff(tt.sourceLines, tt.patchedLines)
			if len(operations) != tt.expectOps {
				t.Errorf("Expected %d operations, got %d", tt.expectOps, len(operations))
			}
		})
	}
}
