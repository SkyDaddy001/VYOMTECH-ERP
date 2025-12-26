# VYOM ERP Frontend - Next.js Application
# Optimized multi-stage build for production deployment

# Stage 1: Build
FROM node:22-alpine AS builder

WORKDIR /app

# Copy package files only first (better layer caching)
COPY package*.json ./

# Install all dependencies (including dev dependencies for build)
RUN npm ci --prefer-offline --no-audit

# Copy source code
COPY . .

# Build the Next.js application
RUN npm run build

# Prune dependencies to production only
RUN npm prune --production

# Stage 2: Runtime
FROM node:22-alpine AS runtime

WORKDIR /app

# Install dumb-init for proper signal handling
RUN apk add --no-cache dumb-init ca-certificates

# Copy package files
COPY package*.json ./


# Copy node_modules from builder (already pruned)
COPY --from=builder /app/node_modules ./node_modules

# Copy built Next.js application from builder
COPY --from=builder /app/.next ./.next

# Copy public directory
COPY --from=builder /app/public ./public

# Create non-root user for security
RUN addgroup -g 1001 -S nodejs && adduser -S nextjs -u 1001

# Change ownership
RUN chown -R nextjs:nodejs /app

# Switch to non-root user
USER nextjs

# Expose port
EXPOSE 3000

# Set environment to production
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

# Health check
HEALTHCHECK --interval=15s --timeout=10s --retries=5 \
  CMD node -e "require('http').get('http://localhost:3000', (r) => {if (r.statusCode !== 200) throw new Error(r.statusCode)})"

# Use dumb-init to handle signals properly
ENTRYPOINT ["dumb-init", "--"]

# Start the Next.js application
CMD ["npm", "start"]
