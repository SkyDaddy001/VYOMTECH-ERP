#!/bin/bash

HANDLERS_DIR="./internal/handlers"

echo "Fixing type assertions in handler files..."

for file in "$HANDLERS_DIR"/*.go; do
    # Fix double "Val" from type assertion replacement
    sed -i 's/\([a-zA-Z_][a-zA-Z0-9_]*\)ValVal,/\1,/g' "$file"
    
    # Fix missing type assertions  - add back the casting
    sed -i 's/, ok := c\.Get(/Val, ok := c.Get(/g' "$file"
    sed -i '/.*Val, ok := c\.Get/! {/.*userID.*/s/userID, ok :=/userID, ok :=/g}' "$file"
done

echo "✅ Fixed type assertions"
