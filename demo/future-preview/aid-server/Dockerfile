# 第一階段: 編譯階段
FROM golang:1.22 AS builder

WORKDIR /app

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN make swag
RUN make build-linux

# 第二階段: 執行階段

FROM debian:stable-slim

WORKDIR /app

COPY --from=builder /app/bin/aid /app/bin/aid

EXPOSE 8080

CMD ["./bin/aid", "server"]