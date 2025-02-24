
FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server cmd/server.go

FROM gcr.io/distroless/base

COPY --from=builder /app/server .

COPY user_creds.db .

COPY cert.pem .
COPY key.pem .

CMD ["./server"]
