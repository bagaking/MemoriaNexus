#!/bin/bash

# Function to count files in a directory
count_files() {
    local dir_path="$1"
    find "$dir_path" -type f | wc -l | xargs
}

# Function to check if a directory should be ignored
should_ignore() {
    local item="$1"
    local ignore_dirs=("${!2}")

    for ignore_dir in "${ignore_dirs[@]}"; do
        if [[ "$item" == $ignore_dir ]]; then
            return 0
        fi
        if [[ "$item" == $ignore_dir* ]]; then
            return 0
        fi
    done
    return 1
}

# Function to get file statistics (line count and character count)
get_file_stats() {
    local file_path="$1"
    if [[ ! -s "$file_path" ]]; then
        echo "(0 chars)"
    else
        local line_count=$(wc -l < "$file_path" | xargs)
        local char_count=$(wc -m < "$file_path" | xargs)
        echo "(${line_count} lines, ${char_count} chars)"
    fi
}

# Function to recursively generate markdown for directory structure
generate_tree() {
    local dir_path="$1"
    local indent="$2"
    local ignore_dirs=("${!3}")

    # First, list all directories
    for item in "$dir_path"/*; do
        # Skip hidden files and directories
        [[ "$(basename "$item")" =~ ^\..* ]] && continue

        # Check if the directory should be ignored
        if [[ -d "$item" ]]; then
            if should_ignore "$(basename "$item")" ignore_dirs[@]; then
                continue
            fi
            # Count the number of files in the directory
            local file_count=$(count_files "$item")
            # Print the directory as a header with file count
            echo "${indent}- $(basename "$item")/ (${file_count} files)"
            # Recurse into the directory
            generate_tree "$item" "  $indent" ignore_dirs[@]
        fi
    done

    # Then, list all files
    for item in "$dir_path"/*; do
        # Skip hidden files and directories
        [[ "$(basename "$item")" =~ ^\..* ]] && continue

        # Check if the file should be ignored
        if [[ -f "$item" ]]; then
            if should_ignore "$(basename "$item")" ignore_dirs[@]; then
                continue
            fi
            # Get file statistics
            local stats=$(get_file_stats "$item")
            # Print the file as a list item with stats
            if [[ -z "$stats" ]]; then
                echo "${indent}- $(basename "$item")"
            else
                echo "${indent}- $(basename "$item") $stats"
            fi
        fi
    done
}

# Read ignore directories from arguments
IFS=',' read -r -a ignore_dirs <<< "$1"

# Start generating the tree from the current directory
echo "# Project Structure"
generate_tree "." "" ignore_dirs[@]