#!/bin/bash

HANDLERS_DIR="/c/Users/Skydaddy/Desktop/VYOM - ERP/internal/handlers"

# Find all files using gin.Context but without the gin import
find "$HANDLERS_DIR" -name "*.go" ! -name "*_test.go" | while read file; do
  if grep -q "(c \*gin.Context)" "$file" && ! grep -q "github.com/gin-gonic/gin" "$file"; then
    echo "Adding gin import to: $(basename $file)"
    
    # Add gin import after "net/http"
    sed -i '/^[[:space:]]*"net\/http"$/a\\t"github.com/gin-gonic/gin"' "$file"
  fi
done

echo "Fixed gin imports"
