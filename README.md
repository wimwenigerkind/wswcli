## Installation

### Using Go
```bash
go install github.com/wimwenigerkind/wswcli@latest
```

### Using Homebrew (macOS/Linux)

```bash
brew install wimwenigerkind/tap/wswcli
```
or
```bash
brew tap wimwenigerkind/homebrew-tap
brew install wswcli
```

### Download Binary
Download the latest binary from the [releases page](https://github.com/wimwenigerkind/wswcli/releases).

## Usage

```bash
wswcli --help
```

### PatchVendor Command

Generate unified diff patches for Shopware vendor modifications:

```bash
# Basic usage
wswcli patchvendor source.php patched.php output.patch

# Interactive mode
wswcli patchvendor

# Directory processing
wswcli patchvendor vendor/shopware/core custom/patches patches/
```

#### Features
- **Directory processing**: Batch process entire directory structures  
- **Smart validation**: Comprehensive input validation with detailed error messages
- **Vendor path handling**: Automatic vendor path extraction and normalization
- **Interactive mode**: Guided workflow with helpful prompts and suggestions

For detailed documentation, see [docs/patchvendor.md](docs/patchvendor.md).

### TwigBlocks Command

Find duplicate Twig blocks in Shopware/Symfony projects to prevent template conflicts:

```bash
# Scan current directory
wswcli twigblocks .

# Scan specific project
wswcli twigblocks /path/to/project

# Generate JSON report
wswcli twigblocks . --output report.json

# CI/CD integration with Bitbucket format
wswcli twigblocks . --bitbucket
```

#### Features
- **Recursive scanning**: Finds all `*.html.twig` files in project directories
- **Duplicate detection**: Identifies blocks with same name/content across files
- **CI/CD ready**: Exit codes and multiple output formats for automation
- **Bitbucket integration**: Native support for Bitbucket Pipes reporting
- **Smart filtering**: Automatically ignores common build/cache directories

For detailed documentation, see [docs/twigblocks.md](docs/twigblocks.md).

## Development

### Prerequisites
- Go 1.21 or later

### Building
```bash
make build
```

### Testing
```bash
make test
```

### Creating a Release
```bash
make release
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.