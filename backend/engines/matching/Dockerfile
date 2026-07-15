FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /matching ./cmd/matching

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /matching /matching
EXPOSE 8089
ENTRYPOINT ["/matching"]
