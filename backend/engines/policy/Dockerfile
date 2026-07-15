FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /policy ./cmd/policy

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /policy /policy
COPY --from=builder /app/policies /policies
EXPOSE 8083
ENTRYPOINT ["/policy"]
