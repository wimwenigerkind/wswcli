# TwigBlocks Command

The `twigblocks` command helps you find duplicate Twig blocks in your Shopware or Symfony projects. Duplicate blocks can cause template inheritance issues and unexpected behavior in your application.

## Overview

This command scans all `*.html.twig` files in a project directory (and subdirectories) to identify:
- Blocks with identical names within the same file (duplicate block definitions)
- Multiple block definitions that would cause template rendering conflicts

## Usage

### Basic Usage

```bash
# Scan current directory
wswcli twigblocks .

# Scan specific project directory
wswcli twigblocks /path/to/your/project

# Scan with verbose output
wswcli twigblocks /path/to/project --output report.json
```

### Command Options

| Flag | Short | Description |
|------|-------|-------------|
| `--bitbucket` | | Output report in Bitbucket Pipes format for CI/CD |
| `--output` | `-o` | Save detailed report to JSON file |

### Examples

```bash
# Basic scan
wswcli twigblocks .

# Generate JSON report
wswcli twigblocks . --output duplicate-blocks-report.json

# CI/CD integration with Bitbucket
wswcli twigblocks . --bitbucket > bitbucket-report.json
```

## What It Detects

### Duplicate Block Definitions Within Same File

Blocks with the same name defined multiple times in the same file, which causes template rendering conflicts:

```twig
<!-- File: templates/product/detail.html.twig -->
{% block product_title %}
    <h1>{{ product.name }}</h1>
{% endblock %}

<!-- Later in the same file -->
{% block product_title %}
    <h2>{{ product.title }}</h2>  <!-- This will cause a conflict! -->
{% endblock %}
```

**Note:** Blocks with the same name in different files are **NOT** reported as duplicates, since this is normal behavior in Twig template inheritance where child templates override parent blocks.

## Output Formats

### Standard Output

```
============================================================
TWIG BLOCK DUPLICATE ANALYSIS REPORT
============================================================
❌ Found 2 duplicate block groups:

1. Block: 'product_title'
   Hash: a1b2c3
   Occurrences: 2
   Files:
     - templates/product/detail.html.twig:15
       Content: {% block product_title %}
     - templates/category/listing.html.twig:23
       Content: {% block product_title %}

2. Block: 'sidebar'
   Hash: d4e5f6
   Occurrences: 1
   Files:
     - templates/base.html.twig:45
       Content: {% block sidebar %}

------------------------------------------------------------
Summary: 2 duplicate groups found in 15 files
Please review and consolidate duplicate blocks to avoid template conflicts.
```

### JSON Output (`--output report.json`)

```json
{
  "summary": {
    "files_scanned": 15,
    "duplicate_groups": 2,
    "status": "FAILED"
  },
  "duplicates": [
    {
      "block_name": "product_title",
      "hash": "a1b2c3",
      "count": 2,
      "files": [
        {
          "name": "product_title",
          "file": "templates/product/detail.html.twig",
          "line": 15,
          "content": "{% block product_title %}",
          "hash": "a1b2c3"
        },
        {
          "name": "product_title", 
          "file": "templates/category/listing.html.twig",
          "line": 23,
          "content": "{% block product_title %}",
          "hash": "a1b2c3"
        }
      ]
    }
  ],
  "files": [
    "templates/base.html.twig",
    "templates/product/detail.html.twig"
  ]
}
```

### Bitbucket Pipelines Format (`--bitbucket`)

Generates JUnit XML format that Bitbucket Pipelines automatically detects and displays in the Tests tab:

**Generated file:** `test-reports/twig-blocks-junit.xml`

```xml
<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="TwigBlockDuplicateAnalysis" tests="15" failures="2" errors="0" time="0">
  <testcase classname="TwigBlocks" name="TwigBlocks.templates.product.detail.html.twig" time="0">
    <failure message="Duplicate Twig blocks found" type="DuplicateBlockError">
Line 15: Duplicate block 'product_title' (appears 2 times in this file)
Line 23: Duplicate block 'product_title' (appears 2 times in this file)
    </failure>
  </testcase>
  <testcase classname="TwigBlocks" name="TwigBlocks.templates.base.html.twig" time="0"/>
  <!-- ... more test cases ... -->
</testsuite>
```

**Console output:**
```
Bitbucket test report generated: test-reports/twig-blocks-junit.xml
❌ FAILED: Found 2 duplicate block groups in 15 files
```

## CI/CD Integration

### Exit Codes

- `0`: No duplicates found (success)
- `1`: Duplicates found (failure)

### GitHub Actions

```yaml
name: Twig Block Analysis

on: [push, pull_request]

jobs:
  twig-analysis:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Install wswcli
        run: |
          curl -L https://github.com/wimwenigerkind/wswcli/releases/latest/download/wswcli_Linux_x86_64.tar.gz | tar xz
          sudo mv wswcli /usr/local/bin/
      
      - name: Check for duplicate Twig blocks
        run: wswcli twigblocks . --output twig-report.json
      
      - name: Upload report
        uses: actions/upload-artifact@v3
        if: failure()
        with:
          name: twig-duplicate-report
          path: twig-report.json
```

### Bitbucket Pipelines

```yaml
pipelines:
  default:
    - step:
        name: Twig Block Analysis
        image: alpine:latest
        script:
          - apk add --no-cache curl tar
          - curl -L https://github.com/wimwenigerkind/wswcli/releases/latest/download/wswcli_Linux_x86_64.tar.gz | tar xz
          - ./wswcli twigblocks . --bitbucket
        artifacts:
          - test-reports/**
```

**Note:** Bitbucket Pipelines will automatically detect the JUnit XML file in `test-reports/` and display failed tests in the Tests tab with detailed failure information.

### GitLab CI

```yaml
twig_analysis:
  stage: test
  image: alpine:latest
  before_script:
    - apk add --no-cache curl tar
    - curl -L https://github.com/wimwenigerkind/wswcli/releases/latest/download/wswcli_Linux_x86_64.tar.gz | tar xz
    - mv wswcli /usr/local/bin/
  script:
    - wswcli twigblocks . --output twig-report.json
  artifacts:
    reports:
      junit: twig-report.json
    when: always
    expire_in: 1 week
```

### Docker Usage

```bash
# Using the official Docker image
docker run --rm -v $(pwd):/workspace ghcr.io/wimwenigerkind/wswcli:latest twigblocks /workspace

# With output file
docker run --rm -v $(pwd):/workspace ghcr.io/wimwenigerkind/wswcli:latest twigblocks /workspace --output /workspace/report.json

# Bitbucket format
docker run --rm -v $(pwd):/workspace ghcr.io/wimwenigerkind/wswcli:latest twigblocks /workspace --bitbucket
```

## Best Practices

### 1. Regular Scanning

Run the command regularly as part of your development workflow:

```bash
# Add to your Makefile
check-twig:
	wswcli twigblocks . --output reports/twig-duplicates.json

# Add to package.json scripts
{
  "scripts": {
    "check:twig": "wswcli twigblocks . --output reports/twig-duplicates.json"
  }
}
```

### 2. Pre-commit Hooks

Add to `.pre-commit-config.yaml`:

```yaml
repos:
  - repo: local
    hooks:
      - id: twig-duplicates
        name: Check for duplicate Twig blocks
        entry: wswcli twigblocks .
        language: system
        pass_filenames: false
        files: '\.html\.twig$'
```

### 3. Ignore Patterns

The command automatically ignores common directories:
- `node_modules/`
- `vendor/`
- `var/`
- `cache/`
- `build/`
- Hidden directories (starting with `.`)

## Troubleshooting

### Common Issues

1. **No files found**: Ensure you're running the command in a directory containing `*.html.twig` files
2. **Permission errors**: Make sure the command has read access to all directories
3. **Large projects**: For very large projects, consider using the `--output` flag to save results to a file

### Performance

- The command is optimized for large codebases
- Scanning 1000+ files typically takes less than 10 seconds
- Memory usage scales linearly with the number of blocks found

### Limitations

- Only detects blocks in `*.html.twig` files
- Does not analyze block inheritance chains
- Simple content hashing (may not detect complex semantic duplicates)

## Integration with IDEs

### VS Code

Create a task in `.vscode/tasks.json`:

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Check Twig Duplicates",
      "type": "shell",
      "command": "wswcli",
      "args": ["twigblocks", ".", "--output", "twig-report.json"],
      "group": "test",
      "presentation": {
        "echo": true,
        "reveal": "always",
        "focus": false,
        "panel": "shared"
      },
      "problemMatcher": []
    }
  ]
}
```

### PhpStorm

Add as an External Tool:
1. Go to Settings → Tools → External Tools
2. Add new tool:
   - Name: "Check Twig Duplicates"
   - Program: `wswcli`
   - Arguments: `twigblocks $ProjectFileDir$ --output $ProjectFileDir$/twig-report.json`
   - Working directory: `$ProjectFileDir$`

## Related Commands

- [`patchvendor`](patchvendor.md) - Generate patches for vendor modifications
- Use together for comprehensive project quality checks

## Support

For issues, feature requests, or questions:
- GitHub Issues: https://github.com/wimwenigerkind/wswcli/issues
- Documentation: https://github.com/wimwenigerkind/wswcli#readme