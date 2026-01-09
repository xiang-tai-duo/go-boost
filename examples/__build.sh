#!/bin/bash

# Create output directory
mkdir -p .tmp

# Iterate through all .go files
for file in *.go; do
    # Get filename without extension
    name=$(basename "$file" .go)
    
    # Compile Go file
    echo "Compiling $file..."
    go build -o .tmp/$name "$file"
    
    # Check if compilation succeeded
    if [ $? -eq 0 ]; then
        echo "✓ Successfully compiled $name"
    else
        echo "✗ Failed to compile $name"
    fi
done

echo "\nBuild completed! Executables are in .tmp/ directory."