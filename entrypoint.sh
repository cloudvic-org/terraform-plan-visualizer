#!/bin/bash

# Terraform Plan Visualizer GitHub Action Entrypoint

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
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

# Get inputs
PLAN_FILE="$1"
OUTPUT_FILE="$2"
UPLOAD_ARTIFACT="$3"

# Validate inputs
if [ -z "$PLAN_FILE" ]; then
    print_error "plan-file input is required"
    exit 1
fi

if [ ! -f "$PLAN_FILE" ]; then
    print_error "Plan file '$PLAN_FILE' does not exist"
    exit 1
fi

# Set default output file if not provided
if [ -z "$OUTPUT_FILE" ]; then
    OUTPUT_FILE="terraform-plan-visualization.html"
fi

print_status "Terraform Plan Visualizer Action"
print_status "Plan file: $PLAN_FILE"
print_status "Output file: $OUTPUT_FILE"
print_status "Upload artifact: $UPLOAD_ARTIFACT"

# Generate the visualization
print_status "Generating visualization..."
terraform-plan-visualizer -i "$PLAN_FILE" -o "$OUTPUT_FILE"

if [ $? -eq 0 ]; then
    print_success "Visualization generated successfully: $OUTPUT_FILE"
    
    # Set output for GitHub Actions
    echo "html-file=$OUTPUT_FILE" >> $GITHUB_OUTPUT
    
    # Upload as artifact if requested
    if [ "$UPLOAD_ARTIFACT" = "true" ]; then
        print_status "Uploading artifact..."
        # Note: This would need to be handled by the calling workflow
        # as we can't directly upload artifacts from within a Docker action
        print_warning "Artifact upload requested but must be handled by the calling workflow"
    fi
else
    print_error "Failed to generate visualization"
    exit 1
fi

print_success "Action completed successfully!"
