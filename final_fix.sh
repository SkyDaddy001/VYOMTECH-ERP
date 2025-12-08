#!/bin/bash

HANDLERS_DIR="./internal/handlers"

echo "Final handler fixes..."

for file in "$HANDLERS_DIR"/*.go; do
    # Replace *mux.Router parameters with *gin.Engine
    sed -i 's/router \*mux\.Router/router *gin.Engine/g' "$file"
    sed -i 's/r \*mux\.Router/r *gin.Engine/g' "$file"
    
    # Remove router.HandleFunc calls - they will be replaced with route registration
    sed -i '/router\.HandleFunc/d' "$file"
    sed -i '/r\.HandleFunc/d' "$file"
    
    # Replace mux.Vars(r) calls that remain
    sed -i 's/mux\.Vars(r)/\/\/ TODO: Use c.Param() in handler/g' "$file"
done

echo "✅ Final fixes complete"
