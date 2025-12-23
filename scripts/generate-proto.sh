#!/bin/bash

# Proto generation script for VYOMTECH-ERP
# Generates Go code from .proto files using protoc

set -e

PROTO_DIR="api/proto"
GO_OUT_DIR="api/pb"

echo "🔧 Generating protobuf code..."

# Install protoc if not present
if ! command -v protoc &> /dev/null; then
    echo "❌ protoc is not installed. Please install protobuf compiler."
    echo "   Windows: choco install protoc"
    echo "   macOS: brew install protobuf"
    echo "   Linux: apt-get install protobuf-compiler"
    exit 1
fi

# Check for protoc-gen-go plugin
if ! command -v protoc-gen-go &> /dev/null; then
    echo "📦 Installing protoc-gen-go plugin..."
    go install github.com/golang/protobuf/protoc-gen-go@latest
fi

# Check for protoc-gen-go-grpc plugin
if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "📦 Installing protoc-gen-go-grpc plugin..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# Create output directory
mkdir -p "$GO_OUT_DIR"

# Generate code from proto files
echo "📝 Generating RBAC proto..."
protoc -I "$PROTO_DIR" \
    --go_out="$GO_OUT_DIR" \
    --go-grpc_out="$GO_OUT_DIR" \
    "$PROTO_DIR/rbac/rbac.proto"

echo "📝 Generating Audit proto..."
protoc -I "$PROTO_DIR" \
    --go_out="$GO_OUT_DIR" \
    --go-grpc_out="$GO_OUT_DIR" \
    "$PROTO_DIR/audit/audit.proto"

echo "✅ Proto code generation completed!"
echo "   Generated files in: $GO_OUT_DIR"
echo ""
echo "📚 Proto files:"
find "$PROTO_DIR" -name "*.proto" -exec echo "   - {}" \;
