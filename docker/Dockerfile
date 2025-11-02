FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# cборка
RUN go build -o /app/rtce ./cmd/editor-server

FROM alpine:latest

WORKDIR /app

# сертификаты для https
RUN apk --no-cache add ca-certificates

# бинарь после сборки
COPY --from=builder /app/rtce .

EXPOSE 8080

CMD ["./rtce"]