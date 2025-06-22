# ----------------------------------------
# Stage 1: Build Application
# ----------------------------------------
FROM golang:1.24.3-alpine AS builder

# Install build tools and dependencies
RUN apk add --no-cache git gcc musl-dev make protoc protobuf protobuf-dev

# Install required Go tools (pinned versions)
RUN go install github.com/google/wire/cmd/wire@v0.6.0 && \
    go install github.com/golang/mock/mockgen@latest && \
    go install github.com/swaggo/swag/cmd/swag@v1.16.4 && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.6 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1 && \
    go install github.com/favadi/protoc-go-inject-tag@latest

# Set working directory and enable module caching
WORKDIR /app
ENV GOPATH=/go
RUN mkdir -p /go/pkg/mod && chmod -R 777 /go

# Copy dependency files and download modules
COPY go.mod go.sum Makefile ./
RUN go mod download

# Copy source code and build
COPY . .

RUN make generate
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o resik .

# ----------------------------------------
# Stage 2: Runtime Image
# ----------------------------------------
FROM alpine:3.21

# Install minimal runtime dependencies and healthcheck tools
RUN apk add --no-cache ca-certificates tzdata && \
    apk add --no-cache curl wget && \
    # Install grpcurl - adjust version and architecture as needed
    GRPCURL_VERSION=1.9.3 && \
    wget "https://github.com/fullstorydev/grpcurl/releases/download/v${GRPCURL_VERSION}/grpcurl_${GRPCURL_VERSION}_linux_x86_64.tar.gz" -O /tmp/grpcurl.tar.gz && \
    tar -xzf /tmp/grpcurl.tar.gz -C /usr/local/bin grpcurl && \
    rm /tmp/grpcurl.tar.gz && \
    update-ca-certificates

# Create non-root user and group
ARG APP_USER=resik-user
ARG APP_GROUP=resik-group
ARG APP_UID=15001
ARG APP_GID=15001

RUN addgroup -g ${APP_GID} -S ${APP_GROUP} && \
    adduser -u ${APP_UID} -S -D -H -G ${APP_GROUP} ${APP_USER} && \
    mkdir -p /opt/resik/app && \
    chown -R ${APP_USER}:${APP_GROUP} /opt/resik

WORKDIR /opt/resik

# Copy artifacts from builder stage
COPY --from=builder --chown=${APP_USER}:${APP_GROUP} /app/resik .
COPY --from=builder --chown=${APP_USER}:${APP_GROUP} /app/app/rest/docs ./app/rest/docs
COPY --from=builder --chown=${APP_USER}:${APP_GROUP} /app/static ./static

# Copy .proto files needed for grpcurl health checks
COPY --from=builder --chown=${APP_USER}:${APP_GROUP} /app/app/grpc/handler/health/health.proto ./app/grpc/handler/health/health.proto
COPY --from=builder --chown=${APP_USER}:${APP_GROUP} /app/app/grpc/proto/status/status.proto ./app/grpc/proto/status/status.proto
COPY --from=builder --chown=${APP_USER}:${APP_GROUP} /app/core/health/response/healthz.proto ./core/health/response/healthz.proto
COPY --from=builder --chown=${APP_USER}:${APP_GROUP} /app/core/health/response/readyz.proto ./core/health/response/readyz.proto

# Switch to non-root user and set entrypoint
USER ${APP_USER}
EXPOSE 8081 8071
ENTRYPOINT ["./resik"]
