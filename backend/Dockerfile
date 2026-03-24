FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/api ./cmd/api

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /bin/api /app/api
COPY --from=builder /app/config /app/config

EXPOSE 8080
ENV PORT=8080
CMD ["/app/api"]
