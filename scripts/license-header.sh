#!/usr/bin/env bash

# License header management script for Go files
# Adds or updates Apache 2.0 license headers with SPDX identifier
#
# Assumption: License headers are ALWAYS in the first 2 lines of the file

set -euo pipefail

# Configuration
YEAR=$(date +%Y)
COPYRIGHT_HOLDER="LeafLock Security Solutions"
LICENSE_IDENTIFIER="Apache-2.0"

# Generate expected license header (with comment markers)
generate_header() {
    cat <<EOF
// Copyright ${YEAR} ${COPYRIGHT_HOLDER}
// SPDX-License-Identifier: ${LICENSE_IDENTIFIER}
EOF
}

# Check if file has correct license header (first 2 lines must match exactly)
has_correct_license_header() {
    local file="$1"
    local expected_header
    expected_header=$(generate_header)

    # Read first 2 lines of file
    local actual_header
    actual_header=$(head -n 2 "$file" 2>/dev/null || echo "")

    if [ "$actual_header" = "$expected_header" ]; then
        return 0  # Correct header
    else
        return 1  # Missing or outdated
    fi
}

# Check if first 2 lines look like a license header
has_any_license_header() {
    local file="$1"
    local first_two_lines
    first_two_lines=$(head -n 2 "$file" 2>/dev/null || echo "")

    # Check if first 2 lines contain copyright or SPDX markers
    if echo "$first_two_lines" | grep -qE "(Copyright|SPDX-License-Identifier)"; then
        return 0  # Has some kind of license header
    else
        return 1  # No license header
    fi
}

# Fix license header in file
# If file has existing license header (lines 1-2), replace it
# If file has no license header, prepend it (don't delete any code)
fix_license_header() {
    local file="$1"
    local temp_file="${file}.tmp"

    if has_any_license_header "$file"; then
        # File has existing license header in first 2 lines - replace it
        {
            generate_header
            echo ""
            tail -n +3 "$file"
        } > "$temp_file"
    else
        # File has no license header - prepend it (don't delete anything)
        {
            generate_header
            echo ""
            cat "$file"
        } > "$temp_file"
    fi

    mv "$temp_file" "$file"
}

# Check mode - verify files have correct headers
check_files() {
    local staged_only="${1:-false}"
    local exit_code=0
    local files_needing_fix=()

    if [ "$staged_only" = "true" ]; then
        # Check only staged Go files
        while IFS= read -r file; do
            [ -f "$file" ] || continue
            if ! has_correct_license_header "$file"; then
                files_needing_fix+=("$file")
                exit_code=1
            fi
        done < <(git diff --cached --name-only --diff-filter=ACM | grep '\.go$' || true)
    else
        # Check all Go files
        while IFS= read -r -d '' file; do
            if ! has_correct_license_header "$file"; then
                files_needing_fix+=("$file")
                exit_code=1
            fi
        done < <(find . -name "*.go" -not -path "./vendor/*" -print0)
    fi

    if [ ${#files_needing_fix[@]} -gt 0 ]; then
        echo "❌ Files with missing or outdated license headers (${#files_needing_fix[@]}):"
        printf '  %s\n' "${files_needing_fix[@]}"
        echo ""
        echo "Run 'make license-fix' to fix all headers"
    else
        if [ "$staged_only" = "true" ]; then
            echo "✓ All staged Go files have correct license headers"
        else
            echo "✓ All Go files have correct license headers"
        fi
    fi

    return $exit_code
}

# Fix mode - add or replace headers to match template
fix_headers() {
    local count=0

    while IFS= read -r -d '' file; do
        if ! has_correct_license_header "$file"; then
            fix_license_header "$file"
            echo "✓ Fixed license header: $file"
            ((count++))
        fi
    done < <(find . -name "*.go" -not -path "./vendor/*" -print0)

    echo ""
    if [ $count -eq 0 ]; then
        echo "All files already have correct license headers"
    else
        echo "Fixed license headers in $count file(s)"
    fi
}

# Show usage
usage() {
    cat <<EOF
Usage: $(basename "$0") [COMMAND]

Manage Apache 2.0 license headers in Go source files.

Assumption: License headers are always in the first 2 lines of files.

Commands:
  check          Check if all Go files have correct license headers (default)
  check-staged   Check only staged Go files (for pre-commit hook)
  fix            Fix all files (add missing or replace outdated headers)

Configuration (dynamic):
  Copyright year:    ${YEAR} (auto-detected from current year)
  Copyright holder:  ${COPYRIGHT_HOLDER}
  License:           ${LICENSE_IDENTIFIER}

Examples:
  $(basename "$0")              # Check all files
  $(basename "$0") check        # Check all files
  $(basename "$0") check-staged # Check only staged files (pre-commit)
  $(basename "$0") fix          # Fix all files with missing/outdated headers

Makefile targets:
  make license-check            # Check all files
  make license-fix              # Fix all files

EOF
}

# Main
main() {
    local command="${1:-check}"

    case "$command" in
        check)
            check_files false
            ;;
        check-staged)
            check_files true
            ;;
        fix)
            fix_headers
            ;;
        -h|--help|help)
            usage
            exit 0
            ;;
        *)
            echo "Error: Unknown command '$command'"
            echo ""
            usage
            exit 1
            ;;
    esac
}

main "$@"
