# Use official Go image for building
FROM golang:alpine3.22 AS builder
WORKDIR /app
RUN apk --no-cache add gcc musl-dev sqlite-dev
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o palworld-helper ./cmd/

# Use a smaller base image for the final container
FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup
WORKDIR /root/
COPY --from=builder /app/palworld-helper .
COPY --from=builder /app/web/static ./app/web/static/
COPY --from=builder /app/web/templates ./app/web/templates/
RUN chown -R appuser:appgroup .
USER appuser
EXPOSE 8080
CMD ["./palworld-helper"]