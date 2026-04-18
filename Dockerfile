# --- Build Stage ---
FROM golang:1.24-alpine AS builder

# Disable CGO for a static binary which is faster to run and guarantees alpine compatibility
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /src
COPY go.mod ./
# COPY go.sum ./  # We don't have dependencies, but if we did, we'd copy it here

# Copy the source code
COPY . .

# Build the binary
RUN go build -ldflags="-s -w" -o hypersweep cmd/hypersweep/main.go

# --- Run Stage ---
FROM alpine:latest

# We need CA certificates to verify HTTPS connections
RUN apk --no-cache add ca-certificates

# Copy the compiled binary from the builder stage
COPY --from=builder /src/hypersweep /usr/local/bin/hypersweep

# When GitHub Actions uses this, it mounts the repo code under /github/workspace
# and sets it as the working directory. So we don't need to specify WORKDIR here.
ENTRYPOINT ["hypersweep"]
