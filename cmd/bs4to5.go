package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// BootstrapMigration represents a single migration rule
type BootstrapMigration struct {
	Name        string
	Pattern     *regexp.Regexp
	Replacement string
	Description string
}

var (
	dryRun bool
)

var bs4to5Cmd = &cobra.Command{
	Use:   "bs-4-to-5 [PATH]",
	Short: "Migrate Bootstrap 4 to Bootstrap 5 in HTML and Twig files",
	Long: `Migrate Bootstrap 4 to Bootstrap 5 classes and components in HTML and Twig files.

This command scans HTML (.html) and Twig (.html.twig, .twig) files in the specified
directory and updates Bootstrap 4 classes to their Bootstrap 5 equivalents.

The migration includes:
- Form component updates (custom-control to form-check, etc.)
- Button and utility class changes
- Grid system updates
- Color and badge class updates
- Component-specific migrations

Examples:
  wswcli bs-4-to-5 .                    # Migrate current directory (recursive)
  wswcli bs-4-to-5 /path/to/templates   # Migrate specific directory (recursive)
  wswcli bs-4-to-5 . --dry-run          # Preview changes without applying`,
	Args: cobra.MaximumNArgs(1),
	RunE: runBS4to5Migration,
}

func init() {
	rootCmd.AddCommand(bs4to5Cmd)
	bs4to5Cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without applying them")
}

func runBS4to5Migration(cmd *cobra.Command, args []string) error {
	// Determine project path
	projectPath := "."
	if len(args) > 0 {
		projectPath = args[0]
	}

	// Validate project path
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", projectPath)
	}

	if dryRun {
		fmt.Printf("🔍 DRY RUN: Previewing Bootstrap 4 to 5 migration in: %s\n", projectPath)
	} else {
		fmt.Printf("🚀 Migrating Bootstrap 4 to 5 in: %s\n", projectPath)
	}

	// Find all relevant files (always recursive)
	files, err := findTemplateFiles(projectPath, true)
	if err != nil {
		return fmt.Errorf("error finding template files: %w", err)
	}

	if len(files) == 0 {
		fmt.Println("No HTML or Twig files found in the specified directory.")
		return nil
	}

	fmt.Printf("Found %d template files\n", len(files))

	// Initialize migration rules
	migrations := getBootstrapMigrations()

	// Process each file
	totalChanges := 0
	for _, file := range files {
		changes, err := processFile(file, migrations)
		if err != nil {
			fmt.Printf("❌ Error processing %s: %v\n", file, err)
			continue
		}
		totalChanges += changes
		if changes > 0 {
			fmt.Printf("✅ %s: %d changes\n", file, changes)
		}
	}

	if dryRun {
		fmt.Printf("\n🔍 DRY RUN COMPLETE: Would make %d changes across %d files\n", totalChanges, len(files))
		fmt.Println("Run without --dry-run to apply changes")
	} else {
		fmt.Printf("\n✅ MIGRATION COMPLETE: Made %d changes across %d files\n", totalChanges, len(files))
	}

	return nil
}

// findTemplateFiles recursively finds all HTML and Twig files
func findTemplateFiles(rootPath string, recursive bool) ([]string, error) {
	var files []string

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories we don't want to process
		if info.IsDir() {
			name := info.Name()
			// Skip if not recursive and it's not the root path
			if !recursive && path != rootPath {
				return filepath.SkipDir
			}
			// Skip common directories
			if strings.HasPrefix(name, ".") ||
				name == "node_modules" ||
				name == "vendor" ||
				name == "var" ||
				name == "cache" ||
				name == "build" {
				return filepath.SkipDir
			}
			return nil
		}

		// Check for HTML and Twig files
		lowerPath := strings.ToLower(path)
		if strings.HasSuffix(lowerPath, ".html") ||
			strings.HasSuffix(lowerPath, ".html.twig") ||
			strings.HasSuffix(lowerPath, ".twig") {
			files = append(files, path)
		}

		return nil
	}

	err := filepath.Walk(rootPath, walkFunc)
	return files, err
}

// processFile applies Bootstrap migrations to a single file
func processFile(filename string, migrations []BootstrapMigration) (int, error) {
	// Read file content
	content, err := os.ReadFile(filename)
	if err != nil {
		return 0, err
	}

	originalContent := string(content)
	modifiedContent := originalContent
	changeCount := 0

	// Apply all migrations
	for _, migration := range migrations {
		matches := migration.Pattern.FindAllString(modifiedContent, -1)
		if len(matches) > 0 {
			newContent := migration.Pattern.ReplaceAllString(modifiedContent, migration.Replacement)
			if newContent != modifiedContent {
				changeCount += len(matches)
				modifiedContent = newContent
				if dryRun {
					fmt.Printf("  📝 %s: %s (%d matches)\n", migration.Name, migration.Description, len(matches))
				}
			}
		}
	}

	// Write back to file if changes were made and not in dry-run mode
	if changeCount > 0 && !dryRun {
		err = os.WriteFile(filename, []byte(modifiedContent), 0644)
		if err != nil {
			return 0, fmt.Errorf("error writing file: %w", err)
		}
	}

	return changeCount, nil
}

// getBootstrapMigrations returns all Bootstrap 4 to 5 migration rules
func getBootstrapMigrations() []BootstrapMigration {
	return []BootstrapMigration{
		// Form Controls
		{
			Name:        "Custom Checkbox",
			Pattern:     regexp.MustCompile(`custom-control\s+custom-checkbox`),
			Replacement: "form-check",
			Description: "custom-control custom-checkbox → form-check",
		},
		{
			Name:        "Custom Radio",
			Pattern:     regexp.MustCompile(`custom-control\s+custom-radio`),
			Replacement: "form-check",
			Description: "custom-control custom-radio → form-check",
		},
		{
			Name:        "Custom Switch",
			Pattern:     regexp.MustCompile(`custom-control\s+custom-switch`),
			Replacement: "form-check form-switch",
			Description: "custom-control custom-switch → form-check form-switch",
		},
		{
			Name:        "Custom Control Input",
			Pattern:     regexp.MustCompile(`custom-control-input`),
			Replacement: "form-check-input",
			Description: "custom-control-input → form-check-input",
		},
		{
			Name:        "Custom Control Label",
			Pattern:     regexp.MustCompile(`custom-control-label`),
			Replacement: "form-check-label",
			Description: "custom-control-label → form-check-label",
		},
		{
			Name:        "Custom Select",
			Pattern:     regexp.MustCompile(`custom-select`),
			Replacement: "form-select",
			Description: "custom-select → form-select",
		},
		{
			Name:        "Custom Range",
			Pattern:     regexp.MustCompile(`custom-range`),
			Replacement: "form-range",
			Description: "custom-range → form-range",
		},
		{
			Name:        "Custom File",
			Pattern:     regexp.MustCompile(`custom-file`),
			Replacement: "form-control",
			Description: "custom-file → form-control (requires additional styling)",
		},

		// Grid System
		{
			Name:        "No Gutters",
			Pattern:     regexp.MustCompile(`no-gutters`),
			Replacement: "g-0",
			Description: "no-gutters → g-0",
		},

		// Buttons
		{
			Name:        "Button Block",
			Pattern:     regexp.MustCompile(`btn-block`),
			Replacement: "d-grid",
			Description: "btn-block → d-grid (may need gap-* utilities)",
		},

		// Badges
		{
			Name:        "Badge Pill",
			Pattern:     regexp.MustCompile(`badge-pill`),
			Replacement: "rounded-pill",
			Description: "badge-pill → rounded-pill",
		},
		{
			Name:        "Badge Primary",
			Pattern:     regexp.MustCompile(`badge-primary`),
			Replacement: "bg-primary",
			Description: "badge-primary → bg-primary",
		},
		{
			Name:        "Badge Secondary",
			Pattern:     regexp.MustCompile(`badge-secondary`),
			Replacement: "bg-secondary",
			Description: "badge-secondary → bg-secondary",
		},
		{
			Name:        "Badge Success",
			Pattern:     regexp.MustCompile(`badge-success`),
			Replacement: "bg-success",
			Description: "badge-success → bg-success",
		},
		{
			Name:        "Badge Danger",
			Pattern:     regexp.MustCompile(`badge-danger`),
			Replacement: "bg-danger",
			Description: "badge-danger → bg-danger",
		},
		{
			Name:        "Badge Warning",
			Pattern:     regexp.MustCompile(`badge-warning`),
			Replacement: "bg-warning",
			Description: "badge-warning → bg-warning",
		},
		{
			Name:        "Badge Info",
			Pattern:     regexp.MustCompile(`badge-info`),
			Replacement: "bg-info",
			Description: "badge-info → bg-info",
		},
		{
			Name:        "Badge Light",
			Pattern:     regexp.MustCompile(`badge-light`),
			Replacement: "bg-light",
			Description: "badge-light → bg-light",
		},
		{
			Name:        "Badge Dark",
			Pattern:     regexp.MustCompile(`badge-dark`),
			Replacement: "bg-dark",
			Description: "badge-dark → bg-dark",
		},

		// Close Button
		{
			Name:        "Close Button",
			Pattern:     regexp.MustCompile(`class="close"`),
			Replacement: `class="btn-close"`,
			Description: "close → btn-close",
		},

		// Tables
		{
			Name:        "Table Header Light",
			Pattern:     regexp.MustCompile(`thead-light`),
			Replacement: "table-light",
			Description: "thead-light → table-light",
		},
		{
			Name:        "Table Header Dark",
			Pattern:     regexp.MustCompile(`thead-dark`),
			Replacement: "table-dark",
			Description: "thead-dark → table-dark",
		},
		{
			Name:        "Card Deck",
			Pattern:     regexp.MustCompile(`card-deck`),
			Replacement: "row row-cols-1 row-cols-md-3 g-4", // Common replacement
			Description: "card-deck → row row-cols-* g-* (card-deck removed, use grid)",
		},
		{
			Name:        "Card Columns",
			Pattern:     regexp.MustCompile(`card-columns`),
			Replacement: "row row-cols-1 row-cols-md-2 row-cols-xl-3",
			Description: "card-columns → row row-cols-* (card-columns removed, use grid)",
		},
		{
			Name:        "Form Group with Custom Control",
			Pattern:     regexp.MustCompile(`form-group`),
			Replacement: "mb-3",
			Description: "form-group → mb-3 (or other margin utility)",
		},
		{
			Name:        "Left alignment",
			Pattern:     regexp.MustCompile(`text-left`),
			Replacement: "text-start",
			Description: "text-left → text-start",
		},
		{
			Name:        "Right alignment",
			Pattern:     regexp.MustCompile(`text-right`),
			Replacement: "text-end",
			Description: "text-right → text-end",
		},
		{
			Name:        "Float Left",
			Pattern:     regexp.MustCompile(`float-left`),
			Replacement: "float-start",
			Description: "float-left → float-start",
		},
		{
			Name:        "Float Right",
			Pattern:     regexp.MustCompile(`float-right`),
			Replacement: "float-end",
			Description: "float-right → float-end",
		},
		{
			Name:        "Border Left",
			Pattern:     regexp.MustCompile(`border-left`),
			Replacement: "border-start",
			Description: "border-left → border-start",
		},
		{
			Name:        "Border Right",
			Pattern:     regexp.MustCompile(`border-right`),
			Replacement: "border-end",
			Description: "border-right → border-end",
		},
		{
			Name:        "Margin Left",
			Pattern:     regexp.MustCompile(`\bml-(\d+)\b`),
			Replacement: "ms-$1",
			Description: "ml-* → ms-*",
		},
		{
			Name:        "Margin Right",
			Pattern:     regexp.MustCompile(`\bmr-(\d+)\b`),
			Replacement: "me-$1",
			Description: "mr-* → me-*",
		},
		{
			Name:        "Padding Left",
			Pattern:     regexp.MustCompile(`\bpl-(\d+)\b`),
			Replacement: "ps-$1",
			Description: "pl-* → ps-*",
		},
		{
			Name:        "Padding Right",
			Pattern:     regexp.MustCompile(`\bpr-(\d+)\b`),
			Replacement: "pe-$1",
			Description: "pr-* → pe-*",
		},
		{
			Name:        "Data Attributes",
			Pattern:     regexp.MustCompile(`data-(toggle|target|dismiss|slide|dropdown|toggle|bs-tooltip|bs-popover|bs-modal)`),
			Replacement: "data-bs-$1",
			Description: "data-* → data-bs-* (essential for JS components)",
		},
		{
			Name:        "Form Row",
			Pattern:     regexp.MustCompile(`form-row`),
			Replacement: "row g-2",
			Description: "form-row → row g-2 (gap utilities)",
		},
		{
			Name:        "Dropdown Menu Right",
			Pattern:     regexp.MustCompile(`dropdown-menu-right`),
			Replacement: "dropdown-menu-end",
			Description: "dropdown-menu-right → dropdown-menu-end",
		},
		{
			Name:        "Rounded Left",
			Pattern:     regexp.MustCompile(`rounded-left`),
			Replacement: "rounded-start",
			Description: "rounded-left → rounded-start",
		},
		{
			Name:        "Rounded Right",
			Pattern:     regexp.MustCompile(`rounded-right`),
			Replacement: "rounded-end",
			Description: "rounded-right → rounded-end",
		},
		{
			Name:        "Input Group Structure",
			Pattern:     regexp.MustCompile(`input-group-(append|prepend)`),
			Replacement: "<!-- TODO: input-group-append/prepend removed - Restructure completely: https://getbootstrap.com/docs/5.2/forms/input-group/ -->",
			Description: "input-group-append/prepend removed - Requires manual restructuring",
		},
		{
			Name:        "Custom File Input",
			Pattern:     regexp.MustCompile(`custom-file-input`),
			Replacement: "<!-- TODO: custom-file-input removed - Use new file input: https://getbootstrap.com/docs/5.2/forms/form-control/#file-input -->",
			Description: "Custom file input removed - Requires complete redesign",
		},
		{
			Name:        "Media Body",
			Pattern:     regexp.MustCompile(`media-body`),
			Replacement: "flex-grow-1",
			Description: "media-body → flex-grow-1",
		},
		{
			Name:        "JavaScript Initialization",
			Pattern:     regexp.MustCompile(`\$\('\.(dropdown|tooltip|popover|modal|carousel|collapse)'\).*?\(\);?`),
			Replacement: "<!-- TODO: jQuery initialization removed - Use new API: https://getbootstrap.com/docs/5.2/getting-started/javascript/#usage -->",
			Description: "jQuery initialization → Vanilla JS initialization",
		},
		{
			Name:        "Breadcrumb Separator",
			Pattern:     regexp.MustCompile(`breadcrumb-item`),
			Replacement: "breadcrumb-item <!-- TODO: Remove manual separators - Now via CSS: https://getbootstrap.com/docs/5.2/components/breadcrumb/#changing-the-separator -->",
			Description: "Breadcrumb separators now via CSS (remove manual separators)",
		},
		{
			Name:        "Close Button Content",
			Pattern:     regexp.MustCompile(`&times;`),
			Replacement: "<!-- TODO: Replace &times; with SVG icon: https://getbootstrap.com/docs/5.2/components/close-button/ -->",
			Description: "&times; must be replaced with SVG icon",
		},
		{
			Name:        "Navbar Toggle Icon",
			Pattern:     regexp.MustCompile(`navbar-toggler-icon`),
			Replacement: "navbar-toggler-icon <!-- TODO: Use SVG icon instead: https://getbootstrap.com/docs/5.2/components/navbar/#toggler -->",
			Description: "Navbar toggle requires SVG instead of icon font",
		},
		{
			Name:        "Form Validation",
			Pattern:     regexp.MustCompile(`(was|is)-(valid|invalid)`),
			Replacement: "has-$2",
			Description: "Validation classes updated to has-* prefix",
		},
		{
			Name:        "Text Monospace",
			Pattern:     regexp.MustCompile(`text-monospace`),
			Replacement: "font-monospace",
			Description: "text-monospace → font-monospace",
		},
		{
			Name:        "Font Weight",
			Pattern:     regexp.MustCompile(`font-weight-(\w+)`),
			Replacement: "fw-$1",
			Description: "font-weight-* → fw-*",
		},
		{
			Name:        "Small Rounded",
			Pattern:     regexp.MustCompile(`rounded-sm`),
			Replacement: "rounded-1",
			Description: "rounded-sm → rounded-1",
		},
		{
			Name:        "Large Rounded",
			Pattern:     regexp.MustCompile(`rounded-lg`),
			Replacement: "rounded-3",
			Description: "rounded-lg → rounded-3",
		},
		{
			Name:        "Responsive Tables",
			Pattern:     regexp.MustCompile(`table-responsive(-\w+)?`),
			Replacement: "table-responsive$1 <!-- TODO: Use wrapper div instead: https://getbootstrap.com/docs/5.2/content/tables/#responsive-tables -->",
			Description: "Responsive tables now require wrapper div",
		},
		{
			Name:        "Tooltip/Popover Positioning",
			Pattern:     regexp.MustCompile(`bs-tooltip-[top|bottom|left|right]`),
			Replacement: "<!-- TODO: Update placement attributes (e.g., 'left' → 'start'): https://getbootstrap.com/docs/5.2/components/tooltips/#position -->",
			Description: "Positioning attributes changed (e.g., 'left' → 'start')",
		},
		{
			Name:        "Input Group Text",
			Pattern:     regexp.MustCompile(`input-group-text`),
			Replacement: "input-group-text <!-- TODO: Restructure without wrapper div: https://getbootstrap.com/docs/5.2/forms/input-group/ -->",
			Description: "input-group-text no longer needs wrapper div",
		},
	}
}
