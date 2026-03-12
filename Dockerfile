# Build stage
FROM golang:alpine AS builder

WORKDIR /build

# Install git (needed for some Go modules)
RUN apk add --no-cache git

# Set GOTOOLCHAIN to auto to allow downloading newer Go versions
ENV GOTOOLCHAIN=auto

# Copy go mod files
COPY go.mod ./
COPY go.sum* ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o erd-viewer ./cmd/erd-viewer

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from build stage
COPY --from=builder /build/erd-viewer .

EXPOSE 8080

ENTRYPOINT ["/app/erd-viewer"]
