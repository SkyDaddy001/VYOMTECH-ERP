#!/bin/bash
# Kubernetes Setup Script for VYOM ERP

set -e

echo "🚀 Setting up VYOM ERP Kubernetes cluster..."

# 1. Create namespace
echo "📦 Creating vyom-erp namespace..."
kubectl create namespace vyom-erp --dry-run=client -o yaml | kubectl apply -f -

# 2. Create migrations ConfigMap
echo "📋 Creating migrations ConfigMap..."
kubectl create configmap migrations-configmap \
  --from-file=migrations/ \
  -n vyom-erp --dry-run=client -o yaml | kubectl apply -f -

# 3. Verify namespace
echo "✅ Checking namespace..."
kubectl get namespace vyom-erp

# 4. Verify ConfigMap
echo "✅ Checking ConfigMap..."
kubectl get configmap -n vyom-erp

echo ""
echo "🎉 Kubernetes setup complete!"
echo ""
echo "Next steps:"
echo "1. Build and push Docker images:"
echo "   podman build -t vyomtech/backend:latest ."
echo "   podman build -t vyomtech/frontend:latest -f Dockerfile.frontend frontend/"
echo "   podman push vyomtech/backend:latest"
echo "   podman push vyomtech/frontend:latest"
echo ""
echo "2. Deploy the complete stack:"
echo "   kubectl apply -f k8s/complete-deployment.yaml"
echo ""
echo "3. Monitor deployment:"
echo "   kubectl get pods -n vyom-erp -w"
