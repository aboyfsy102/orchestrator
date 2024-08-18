#!/bin/bash

# Create the dist directory if it doesn't exist
mkdir -p ../dist

# Clean up the dist directory
echo "Cleaning up the dist directory..."
rm -rf ../dist/*

# # Get current date in yyyymmdd format
# current_date=$(date +"%Y%m%d")

# Iterate through all directories in the current folder
for dir in */; do
    # Remove the trailing slash from the directory name
    dir_name=${dir%/}
    
    # # Generate a random number between 100 and 999
    # random_number=$(shuf -i 100-999 -n 1)
    
    # # Create the suffix
    # suffix="${current_date}-${random_number}"
    
    # Create a zip file for the directory contents with the new suffix
    echo "Creating zip file for $dir_name"
    (cd "$dir_name" && zip -r "../../dist/${dir_name}.zip" . -x "*.DS_Store" "*.git*")
done

echo "All directories have been zipped and placed in ../dist"