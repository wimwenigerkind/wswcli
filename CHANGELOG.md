# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
