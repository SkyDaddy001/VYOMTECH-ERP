#!/bin/bash

HANDLERS_DIR="./internal/handlers"

echo "Fixing Gin handler conversions..."

for file in "$HANDLERS_DIR"/*.go; do
    if grep -q "gin.Context" "$file" 2>/dev/null; then
        echo "Fixing: $(basename $file)"
        
        # Fix json.NewEncoder(w).Encode -> c.JSON  
        sed -i 's/json\.NewEncoder(w)\.Encode(/c.JSON(http.StatusOK, /g' "$file"
        
        # Remove orphaned variable references
        sed -i '/^[[:space:]]*w$/d' "$file"
        sed -i 's/, w)/)/g' "$file"
        sed -i 's/, w,/, /g' "$file"
        sed -i '/^[[:space:]]*w,/d' "$file"
    fi
done

echo "✅ Fixed all handler JSON encoding issues"
