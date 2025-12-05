FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o trace-api cmd/api/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/trace-api .

EXPOSE 8080
CMD ["./trace-api"]