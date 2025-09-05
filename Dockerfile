# Use official Go image for building
FROM golang:1.22-alpine AS builder

# Set working directory
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy all source code (maintaining directory structure)
COPY app/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o palworld-helper .

# Use a smaller base image for the final container
FROM alpine:latest

# Install ca-certificates for HTTPS requests if needed
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /root/

# Copy the binary from builder stage (correct path)
COPY --from=builder /app/palworld-helper .

# Copy static files if needed
COPY --from=builder /app/app/web/statics ./app/web/statics/
COPY --from=builder /app/app/web/templates ./app/web/templates/

# Change ownership to non-root user
RUN chown -R appuser:appgroup .

# Switch to non-root user
USER appuser

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./palworld-helper"]