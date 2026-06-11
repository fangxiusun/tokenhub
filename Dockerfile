# Build stage for frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app/web/default

# Copy package files
COPY web/default/package.json web/default/package-lock.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY web/default/ .

# Build frontend
RUN npm run build

# Build stage for backend
FROM golang:1.22-alpine AS backend-builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Copy built frontend
COPY --from=frontend-builder /app/web/default/dist ./web/default/dist

# Build backend (frontend is embedded into binary)
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary (includes embedded frontend)
COPY --from=backend-builder /app/server .

# Create data directory
RUN mkdir -p /data

# Environment variables
ENV GIN_MODE=release
ENV PORT=3000
ENV SQL_DSN=""
ENV REDIS_CONN_STRING=""
ENV TZ=Asia/Shanghai

EXPOSE 3000

VOLUME ["/data"]

CMD ["./server"]
