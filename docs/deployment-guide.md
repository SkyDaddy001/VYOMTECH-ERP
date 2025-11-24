# Deployment Guide

This guide covers deploying the Multi-Tenant AI Call Center to production environments.

## Prerequisites

- Kubernetes cluster (v1.27+)
- Helm 3.12+
- MySQL 8.0+
- Redis 7+
- Docker registry access
- Domain name and SSL certificate

## Infrastructure Requirements

### Minimum Hardware Requirements

| Component | CPU | Memory | Storage | Replicas |
|-----------|-----|--------|---------|----------|
| API Server | 2 vCPU | 4GB | 50GB | 3 |
| MySQL | 4 vCPU | 16GB | 500GB SSD | 1 (with replicas) |
| Redis | 2 vCPU | 8GB | 100GB | 1 (with replicas) |
| Monitoring | 1 vCPU | 2GB | 50GB | 1 |

### Network Requirements

- Inbound: 80/443 (HTTP/HTTPS)
- Database: 3306 (MySQL), 6379 (Redis)
- Internal: Pod-to-pod communication

## Deployment Steps

### 1. Prepare Secrets

Create Kubernetes secrets for sensitive data:

```bash
# Database credentials
kubectl create secret generic db-secret \
  --from-literal=host=mysql-service \
  --from-literal=password=<strong-password>

# JWT secret (generate a secure random string)
kubectl create secret generic jwt-secret \
  --from-literal=secret=<jwt-secret>

# OpenAI API key
kubectl create secret generic openai-secret \
  --from-literal=api-key=<openai-api-key>

# Email service credentials (if using SMTP)
kubectl create secret generic email-secret \
  --from-literal=smtp-host=<smtp-host> \
  --from-literal=smtp-port=<smtp-port> \
  --from-literal=smtp-user=<smtp-user> \
  --from-literal=smtp-password=<smtp-password>
```

### 2. Deploy Database

Deploy MySQL with persistent storage:

```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install mysql bitnami/mysql \
  --set auth.rootPassword=<root-password> \
  --set auth.database=callcenter \
  --set auth.username=callcenter_user \
  --set auth.password=<db-password> \
  --set persistence.size=500Gi \
  --set metrics.enabled=true
```

### 3. Deploy Redis

Deploy Redis for caching and sessions:

```bash
helm install redis bitnami/redis \
  --set auth.password=<redis-password> \
  --set persistence.size=100Gi \
  --set metrics.enabled=true
```

### 4. Deploy Application

Apply Kubernetes manifests:

```bash
# Apply configuration
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secrets.yaml
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/deployment.yaml

# Wait for rollout
kubectl rollout status deployment/callcenter-api
```

### 5. Run Database Migrations

Execute database migrations:

```bash
# Get MySQL pod name
MYSQL_POD=$(kubectl get pods -l app=mysql -o jsonpath='{.items[0].metadata.name}')

# Run migrations
kubectl exec -it $MYSQL_POD -- mysql -u root -p callcenter < migrations/001_initial_schema.sql
```

### 6. Configure Ingress

Set up ingress for external access:

```yaml
# ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: callcenter-ingress
  annotations:
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
spec:
  tls:
  - hosts:
    - api.yourdomain.com
    - app.yourdomain.com
    secretName: callcenter-tls
  rules:
  - host: api.yourdomain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: callcenter-api-service
            port:
              number: 80
  - host: app.yourdomain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: frontend-service
            port:
              number: 80
```

```bash
kubectl apply -f ingress.yaml
```

### 7. Deploy Monitoring

Deploy Prometheus and Grafana:

```bash
# Add Prometheus community helm repo
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts

# Install Prometheus
helm install prometheus prometheus-community/prometheus \
  --set server.persistentVolume.size=50Gi

# Install Grafana
helm install grafana stable/grafana \
  --set persistence.enabled=true \
  --set persistence.size=10Gi \
  --set adminPassword=<admin-password>
```

Apply custom monitoring configuration:

```bash
kubectl apply -f monitoring/
```

## Configuration

### Environment Variables

Configure the application through ConfigMap:

```yaml
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: callcenter-config
data:
  DB_NAME: "callcenter"
  DB_USER: "callcenter_user"
  REDIS_HOST: "redis-service"
  REDIS_PORT: "6379"
  LOG_LEVEL: "info"
  AI_CACHE_TTL: "1800"
  RATE_LIMIT_REQUESTS: "1000"
  RATE_LIMIT_WINDOW: "60"
  MAX_CONCURRENT_CALLS: "50"
  SESSION_TIMEOUT: "3600"
  PASSWORD_RESET_EXPIRY: "3600"
  SMTP_HOST: "smtp.gmail.com"
  SMTP_PORT: "587"
  EMAIL_FROM: "noreply@yourdomain.com"
```

### Scaling

Configure horizontal pod autoscaling:

```yaml
# hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: callcenter-api-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: callcenter-api
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

```bash
kubectl apply -f hpa.yaml
```

## Monitoring and Observability

### Health Checks

Monitor application health:

```bash
# Check pod status
kubectl get pods

# Check service endpoints
kubectl get endpoints

# View logs
kubectl logs -f deployment/callcenter-api

# Check resource usage
kubectl top pods
```

### Metrics Dashboard

Access Grafana dashboard:

```bash
# Get Grafana admin password
kubectl get secret --namespace default grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo

# Port forward to access Grafana
kubectl port-forward svc/grafana 3000:80

# Access at http://localhost:3000
```

Import the provided dashboard JSON for call center metrics.

### Alerting

Configure alerts for critical issues:

- High error rates (>5%)
- Database connection failures
- Memory/CPU usage >90%
- Pod restarts >5/minute

## Backup and Recovery

### Database Backup

Set up automated MySQL backups:

```bash
# Create backup cron job
kubectl apply -f - <<EOF
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: mysql-backup
spec:
  schedule: "0 2 * * *"  # Daily at 2 AM
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: backup
            image: mysql:8.0
            command:
            - /bin/bash
            - -c
            - mysqldump -h mysql-service -u root -p\$MYSQL_ROOT_PASSWORD callcenter > /backup/backup.sql
            env:
            - name: MYSQL_ROOT_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: mysql-secret
                  key: mysql-root-password
            volumeMounts:
            - name: backup-storage
              mountPath: /backup
          volumes:
          - name: backup-storage
            persistentVolumeClaim:
              claimName: backup-pvc
          restartPolicy: OnFailure
EOF
```

### Disaster Recovery

Implement backup recovery procedures:

1. **Database Recovery**
   ```bash
   # Restore from backup
   kubectl exec -it mysql-pod -- mysql -u root -p callcenter < /backup/backup.sql
   ```

2. **Application Rollback**
   ```bash
   # Rollback deployment
   kubectl rollout undo deployment/callcenter-api
   ```

## Security Hardening

### Network Policies

Implement network segmentation:

```yaml
# network-policy.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: callcenter-network-policy
spec:
  podSelector:
    matchLabels:
      app: callcenter-api
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: frontend
    ports:
    - protocol: TCP
      port: 8080
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: mysql
    ports:
    - protocol: TCP
      port: 3306
  - to:
    - podSelector:
        matchLabels:
          app: redis
    ports:
    - protocol: TCP
      port: 6379
  - to: []  # Allow external API calls
    ports:
    - protocol: TCP
      port: 443
```

### Security Scanning

Regular security scans:

```bash
# Scan container images
trivy image callcenter/api:latest

# Scan Kubernetes manifests
kube-score score k8s/

# Security audit
kube-bench run
```

## Troubleshooting

### Common Issues

1. **Pods not starting**
   ```bash
   kubectl describe pod <pod-name>
   kubectl logs <pod-name>
   ```

2. **Database connection issues**
   ```bash
   kubectl exec -it mysql-pod -- mysql -u root -p -e "SHOW PROCESSLIST;"
   ```

3. **High memory usage**
   ```bash
   kubectl top pods
   kubectl describe hpa callcenter-api-hpa
   ```

4. **Slow API responses**
   ```bash
   # Check Prometheus metrics
   kubectl port-forward svc/prometheus 9090:9090
   # Access at http://localhost:9090
   ```

### Performance Tuning

Optimize for high load:

1. **Database optimization**
   - Add indexes for frequently queried columns
   - Configure connection pooling
   - Enable query caching

2. **Application optimization**
   - Implement response caching
   - Use connection pooling
   - Optimize goroutine usage

3. **Infrastructure optimization**
   - Configure resource limits
   - Use node affinity
   - Implement pod disruption budgets

## Maintenance Procedures

### Regular Tasks

- **Weekly**: Review logs and metrics
- **Monthly**: Update dependencies and security patches
- **Quarterly**: Performance testing and optimization

### Update Procedures

1. **Application updates**
   ```bash
   # Update image tag
   kubectl set image deployment/callcenter-api api=callcenter/api:v1.2.0

   # Monitor rollout
   kubectl rollout status deployment/callcenter-api
   ```

2. **Database schema updates**
   ```bash
   # Apply migrations
   kubectl exec -it mysql-pod -- mysql -u root -p callcenter < migrations/002_new_feature.sql
   ```

## Support and Monitoring

### Monitoring Dashboards

Key metrics to monitor:

- API response times and error rates
- Database connection pool usage
- AI API call latency and success rates
- Agent status and call queue lengths
- Resource utilization (CPU, memory, disk)

### Alert Thresholds

Configure alerts for:

- Error rate > 5%
- Response time > 2 seconds (p95)
- Database connections > 90% capacity
- Memory usage > 85%
- Disk usage > 80%

### Incident Response

1. **Detection**: Alerts trigger notification
2. **Assessment**: Check monitoring dashboards
3. **Containment**: Scale resources or rollback
4. **Recovery**: Apply fixes and test
5. **Post-mortem**: Document lessons learned

This deployment guide ensures a robust, scalable, and secure production environment for the Multi-Tenant AI Call Center system.
