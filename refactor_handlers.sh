#!/bin/bash

# Refactor Go handlers from gorilla/mux to Gin
# This script automates the conversion of all handler files

PROJECT_DIR="/c/Users/Skydaddy/Desktop/VYOM - ERP"
HANDLERS_DIR="$PROJECT_DIR/internal/handlers"

echo "=== Starting Handler Refactoring from gorilla/mux to Gin ==="
echo "Processing directory: $HANDLERS_DIR"

# Step 1: Replace gorilla/mux import with gin
echo ""
echo "Step 1: Replacing imports..."
find "$HANDLERS_DIR" -name "*.go" -type f ! -name "*_test.go" -exec grep -l "github.com/gorilla/mux" {} \; | while read file; do
  echo "  Processing: $(basename $file)"
  
  # Replace the import
  sed -i 's|"github.com/gorilla/mux"|"github.com/gin-gonic/gin"|g' "$file"
done

# Step 2: Convert method signatures from http.ResponseWriter, *http.Request to *gin.Context
echo ""
echo "Step 2: Converting method signatures..."
find "$HANDLERS_DIR" -name "*.go" -type f ! -name "*_test.go" -exec grep -l "(w http.ResponseWriter, r \*http.Request)" {} \; | while read file; do
  echo "  Processing: $(basename $file)"
  
  # Replace method signatures
  sed -i 's/(w http\.ResponseWriter, r \*http\.Request)/(c *gin.Context)/g' "$file"
done

# Step 3: Replace mux.Vars with gin c.Param
echo ""
echo "Step 3: Converting URL parameter extraction..."
find "$HANDLERS_DIR" -name "*.go" -type f ! -name "*_test.go" -exec grep -l "mux.Vars" {} \; | while read file; do
  echo "  Processing: $(basename $file)"
  
  # This requires more complex handling - we'll do pattern-by-pattern replacement
  # Replace: vars := mux.Vars(r); id := vars["id"] with id := c.Param("id")
  sed -i '/vars := mux\.Vars(r)/d' "$file"
  
  # Replace patterns like: vars["fieldname"] with c.Param("fieldname")
  sed -i 's/vars\["\([^"]*\)"\]/c.Param("\1")/g' "$file"
done

# Step 4: Replace json.NewDecoder with c.ShouldBindJSON
echo ""
echo "Step 4: Converting JSON decoding..."
find "$HANDLERS_DIR" -name "*.go" -type f ! -name "*_test.go" -exec grep -l "json.NewDecoder(r.Body).Decode" {} \; | while read file; do
  echo "  Processing: $(basename $file)"
  
  # Replace JSON decode pattern
  sed -i 's/json\.NewDecoder(r\.Body)\.Decode(&\([^)]*\))/c.ShouldBindJSON(\&\1)/g' "$file"
  
  # Update error checks
  sed -i 's/if err := c\.ShouldBindJSON/if err := c.ShouldBindJSON/g' "$file"
done

# Step 5: Replace http.Error with c.JSON responses
echo ""
echo "Step 5: Converting error responses..."
find "$HANDLERS_DIR" -name "*.go" -type f ! -name "*_test.go" -exec grep -l "http.Error" {} \; | while read file; do
  echo "  Processing: $(basename $file)"
  
  # Replace http.Error pattern - this is complex, will need manual review
  # c.JSON(statusCode, gin.H{"error": "message"})
  sed -i 's/http\.Error(w, "\([^"]*\)", http\.StatusBadRequest)/c.JSON(http.StatusBadRequest, gin.H{"error": "\1"})/g' "$file"
  sed -i 's/http\.Error(w, "\([^"]*\)", http\.StatusUnauthorized)/c.JSON(http.StatusUnauthorized, gin.H{"error": "\1"})/g' "$file"
  sed -i 's/http\.Error(w, "\([^"]*\)", http\.StatusInternalServerError)/c.JSON(http.StatusInternalServerError, gin.H{"error": "\1"})/g' "$file"
  sed -i 's/http\.Error(w, "\([^"]*\)", http\.StatusNotFound)/c.JSON(http.StatusNotFound, gin.H{"error": "\1"})/g' "$file"
done

# Step 6: Replace w.Header().Set and w.WriteHeader with Gin responses
echo ""
echo "Step 6: Converting response headers and status codes..."
find "$HANDLERS_DIR" -name "*.go" -type f ! -name "*_test.go" -exec grep -l "w.Header()" {} \; | while read file; do
  echo "  Processing: $(basename $file)"
  
  # Remove Content-Type headers (Gin handles this)
  sed -i '/w\.Header()\.Set("Content-Type", "application\/json")/d' "$file"
  
  # Remove SetHeader calls
  sed -i '/w\.Header()\.Add/d' "$file"
done

find "$HANDLERS_DIR" -name "*.go" -type f ! -name "*_test.go" -exec grep -l "w.WriteHeader" {} \; | while read file; do
  echo "  Processing: $(basename $file)"
  
  # This is complex - will need to be replaced in context
  sed -i 's/w\.WriteHeader(http\.StatusCreated)/\/\/ Status set by c.JSON/g' "$file"
  sed -i 's/w\.WriteHeader(http\.StatusOK)/\/\/ Status set by c.JSON/g' "$file"
done

# Step 7: Replace json.NewEncoder with c.JSON
echo ""
echo "Step 7: Converting JSON encoding..."
find "$HANDLERS_DIR" -name "*.go" -type f ! -name "*_test.go" -exec grep -l "json.NewEncoder" {} \; | while read file; do
  echo "  Processing: $(basename $file)"
  
  # Replace json.NewEncoder(w).Encode() with c.JSON()
  sed -i 's/json\.NewEncoder(w)\.Encode(\(.*\))/c.JSON(http.StatusOK, \1)/g' "$file"
done

# Step 8: Replace context extraction
echo ""
echo "Step 8: Converting context extraction..."
find "$HANDLERS_DIR" -name "*.go" -type f ! -name "*_test.go" -exec grep -l "r.Context().Value" {} \; | while read file; do
  echo "  Processing: $(basename $file)"
  
  # Replace r.Context().Value with c.Get
  sed -i 's/r\.Context()\.Value(\([^)]*\))/c.Get(\1)/g' "$file"
  
  # Replace c.Request.Context() with c.Request.Context()
  sed -i 's/r\.Context()/c.Request.Context()/g' "$file"
done

# Step 9: Remove unused imports
echo ""
echo "Step 9: Cleaning up unused imports..."
find "$HANDLERS_DIR" -name "*.go" -type f ! -name "*_test.go" -exec grep -l "encoding/json" {} \; | while read file; do
  echo "  Processing: $(basename $file)"
  
  # Remove json import if no longer needed (will need manual verification)
  # For now, keep it as it may still be used
  true
done

echo ""
echo "=== Refactoring Complete ==="
echo "Please review the changes and run: go build ./cmd/main.go"
