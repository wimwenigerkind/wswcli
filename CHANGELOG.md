# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2.4.0] - 2025-07-21

### Added
- **Bootstrap 4 to 5 Migration command** (`bs-4-to-5`) for automated migration of HTML and Twig templates

## [2.3.0] - 2025-07-17

### Added
- Configuration support for the `patchvendor` command, enabling flexible patch generation and improved diff path handling

### Changed

- Enhanced vendor path validation in patch tests to ensure correct relative paths and proper `a/b` prefix formatting
- Refactored patchvendor test cases for more accurate path expectations and robust validation

## [2.2.0] - 2025-07-10

### Added
- **TwigBlocks command** for finding duplicate Twig blocks in Shopware/Symfony projects
- Recursive scanning of `*.html.twig` files with smart directory filtering
- Detection of duplicate block definitions within the same file (prevents template conflicts)
- Multiple output formats: human-readable, JSON, and JUnit XML for CI/CD
- Bitbucket Code Insights integration with Reports API and annotations
- File links in test reports for easy navigation to problematic files
- Smart directory filtering (ignores `node_modules`, `vendor`, `cache`, etc.)
- CI/CD integration with proper exit codes for automation
- Comprehensive test suite for TwigBlocks functionality
- Support for both relative (`.`) and absolute path scanning
- Nightly release pipeline with automated daily builds
- Nightly Docker images with `nightly`, `nightly-amd64`, `nightly-arm64` tags
- Manual trigger support for nightly releases via GitHub Actions
- Smart nightly release logic (only releases when there are new commits)
- Separate GoReleaser configuration for nightly builds

### Changed
- Improved directory scanning logic to properly handle root directory (`.`) paths
- Enhanced JUnit XML output with clickable file links for Bitbucket integration

### Fixed
- Fixed issue where scanning with relative path (`.`) would not find any files
- Resolved Docker authentication for GitHub Container Registry in CI/CD
- Fixed GoReleaser configuration compatibility with version 2

## [2.1.0] - 2025-07-10

### Added
- Docker support with multi-architecture images (AMD64, ARM64)
- Multi-platform Docker manifests for automatic platform selection
- Enhanced Dockerfile with security improvements (non-root user, CA certificates)
- GoReleaser configuration updated to version 2

### Changed
- Improved documentation structure and content
- Updated GitHub Actions workflow to use GoReleaser v2 (goreleaser-action@v6)

### Fixed
- GitHub Actions release workflow compatibility with GoReleaser configuration version 2
- Added missing `packages: write` permission for Docker image publishing to GitHub Container Registry
- Docker authentication issue preventing image publishing to GitHub Container Registry

## [2.0.0] - 2025-07-05

### Added
- Homebrew installation support via custom tap
- Installation instructions for multiple package managers

### Changed
- **BREAKING**: Refactored `generateUnifiedDiff` to use system git diff for improved performance
- Enhanced error handling and performance optimizations

### Fixed
- Performance improvements in diff generation

## [1.0.0] - 2025-07-04

### Added
- Initial stable release of wswcli
- PatchVendor command for generating unified diff patches
- Directory processing capabilities
- Interactive mode with guided workflow
- Smart validation with comprehensive error messages
- Vendor path handling and normalization
- Support for Shopware vendor modifications

### Fixed
- Fixed grep command compatibility issue with unrecognized option
- Improved PHP class definition formatting in test files

### Changed
- Project renamed from original name to wswcli
- Updated all configurations and documentation to reflect new project name
