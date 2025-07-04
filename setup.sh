#!/bin/bash

# Go Project Template Setup Script
# This script helps you customize the go-cmd-project-platform-template for your own project

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to validate project name (should be valid for Go module names)
validate_project_name() {
    local name="$1"
    if [[ ! "$name" =~ ^[a-z][a-z0-9-]*$ ]]; then
        print_error "Project name should start with a lowercase letter and contain only lowercase letters, numbers, and hyphens"
        return 1
    fi
    return 0
}

# Function to validate GitHub username
validate_github_username() {
    local username="$1"
    if [[ ! "$username" =~ ^[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?$ ]]; then
        print_error "GitHub username should contain only alphanumeric characters and hyphens, and cannot start or end with a hyphen"
        return 1
    fi
    return 0
}

# Check if we're in the right directory
if [[ ! -f "go.mod" ]] || [[ ! -f "main.go" ]] || [[ ! -f ".goreleaser.yml" ]]; then
    print_error "This script must be run from the root of the go-cmd-project-platform-template project"
    exit 1
fi

# Check if template has already been customized
if ! grep -q "go-cmd-project-platform-template" go.mod; then
    print_warning "This template appears to have already been customized."
    read -p "Do you want to continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Setup cancelled."
        exit 0
    fi
fi

print_info "Welcome to the Go Project Template Setup!"
print_info "This script will help you customize this template for your new project."
echo

# Collect user input
while true; do
    read -p "Enter your project name (e.g., 'my-awesome-tool'): " PROJECT_NAME
    if validate_project_name "$PROJECT_NAME"; then
        break
    fi
done

read -p "Enter a short description for your project: " PROJECT_DESCRIPTION

while true; do
    read -p "Enter your GitHub username: " GITHUB_USERNAME
    if validate_github_username "$GITHUB_USERNAME"; then
        break
    fi
done

read -p "Enter your full name (for LICENSE): " AUTHOR_NAME
read -p "Enter your email (optional): " AUTHOR_EMAIL

# Optional: Homebrew tap repository name
read -p "Enter Homebrew tap repository name (default: homebrew-tap): " HOMEBREW_TAP_REPO
HOMEBREW_TAP_REPO=${HOMEBREW_TAP_REPO:-homebrew-tap}

echo
print_info "Configuration Summary:"
echo "  Project Name: $PROJECT_NAME"
echo "  Description: $PROJECT_DESCRIPTION"
echo "  GitHub Username: $GITHUB_USERNAME"
echo "  Author: $AUTHOR_NAME"
echo "  Email: $AUTHOR_EMAIL"
echo "  Homebrew Tap: $HOMEBREW_TAP_REPO"
echo

read -p "Is this correct? (Y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Nn]$ ]]; then
    print_info "Setup cancelled. Please run the script again."
    exit 0
fi

print_info "Starting project customization..."

# Create backup of original files
print_info "Creating backup of original files..."
mkdir -p .template-backup
cp go.mod .template-backup/
cp main.go .template-backup/
cp cmd/root.go .template-backup/
cp .goreleaser.yml .template-backup/
cp README.md .template-backup/
if [[ -f LICENSE ]]; then
    cp LICENSE .template-backup/
fi

# Replace placeholders in go.mod
print_info "Updating go.mod..."
sed -i.bak "s|github.com/wimwenigerkind/go-cmd-project-platform-template|github.com/$GITHUB_USERNAME/$PROJECT_NAME|g" go.mod
rm go.mod.bak

# Replace placeholders in main.go
print_info "Updating main.go..."
sed -i.bak "s|go-cmd-project-platform-template/cmd|github.com/$GITHUB_USERNAME/$PROJECT_NAME/cmd|g" main.go
rm main.go.bak

# Replace placeholders in cmd/root.go
print_info "Updating cmd/root.go..."
sed -i.bak "s|PLACEHOLDER_USE|$PROJECT_NAME|g" cmd/root.go
sed -i.bak "s|PLACEHOLDER_SHORT|$PROJECT_DESCRIPTION|g" cmd/root.go
sed -i.bak "s|PLACEHOLDER_LONG|$PROJECT_DESCRIPTION|g" cmd/root.go
rm cmd/root.bak

# Replace placeholders in .goreleaser.yml
print_info "Updating .goreleaser.yml..."
sed -i.bak "s|go-cmd-project-platform-template|$PROJECT_NAME|g" .goreleaser.yml
sed -i.bak "s|wimwenigerkind|$GITHUB_USERNAME|g" .goreleaser.yml
sed -i.bak "s|homebrew-tab|$HOMEBREW_TAP_REPO|g" .goreleaser.yml
sed -i.bak "s|\"go-cmd-project-platform-template\"|\"$PROJECT_DESCRIPTION\"|g" .goreleaser.yml
rm .goreleaser.yml.bak

# Replace placeholders in README.md
print_info "Updating README.md..."
cat > README.md << EOF
# $PROJECT_NAME

$PROJECT_DESCRIPTION

## Installation

### Using Go
\`\`\`bash
go install github.com/$GITHUB_USERNAME/$PROJECT_NAME@latest
\`\`\`

### Using Homebrew (macOS/Linux)
\`\`\`bash
brew tap $GITHUB_USERNAME/$HOMEBREW_TAP_REPO
brew install $PROJECT_NAME
\`\`\`

### Download Binary
Download the latest binary from the [releases page](https://github.com/$GITHUB_USERNAME/$PROJECT_NAME/releases).

## Usage

\`\`\`bash
$PROJECT_NAME --help
\`\`\`

## Development

### Prerequisites
- Go 1.24.3 or later

### Building
\`\`\`bash
make build
\`\`\`

### Testing
\`\`\`bash
make test
\`\`\`

### Creating a Release
\`\`\`bash
make release
\`\`\`

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
EOF

# Update LICENSE file if it exists
if [[ -f LICENSE ]]; then
    print_info "Updating LICENSE..."
    CURRENT_YEAR=$(date +%Y)
    if [[ -n "$AUTHOR_EMAIL" ]]; then
        AUTHOR_INFO="$AUTHOR_NAME <$AUTHOR_EMAIL>"
    else
        AUTHOR_INFO="$AUTHOR_NAME"
    fi
    
    # Update MIT License
    sed -i.bak "s|Copyright (c) [0-9]* .*|Copyright (c) $CURRENT_YEAR $AUTHOR_INFO|g" LICENSE
    rm LICENSE.bak
fi

# Clean up Go modules
print_info "Cleaning up Go modules..."
go mod tidy

# Initialize git repository if not already initialized
if [[ ! -d ".git" ]]; then
    print_info "Initializing Git repository..."
    git init
    git add .
    git commit -m "Initial commit: $PROJECT_NAME

Generated from go-cmd-project-platform-template"
else
    print_info "Git repository already exists. You may want to commit these changes manually."
fi

# Test build
print_info "Testing build..."
if make build; then
    print_success "Build successful!"
else
    print_error "Build failed. Please check the configuration."
    exit 1
fi

print_success "Project setup completed successfully!"
echo
print_info "Next steps:"
echo "  1. Review the generated files"
echo "  2. Update the project description in README.md if needed"
echo "  3. Add your project logic to cmd/root.go or create new commands"
echo "  4. Create a GitHub repository: https://github.com/new"
echo "  5. Push your code: git remote add origin https://github.com/$GITHUB_USERNAME/$PROJECT_NAME.git"
echo "  6. Create your first release to enable Homebrew installation"
echo
print_info "Template backup files are stored in .template-backup/ directory"
print_warning "You can delete this setup.sh script and .template-backup/ directory when you're satisfied with the setup"

# Offer to remove setup script
echo
read -p "Do you want to remove this setup script now? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    rm -rf .template-backup
    rm setup.sh
    print_success "Setup script and backup files removed."
fi

print_success "Happy coding! 🚀"