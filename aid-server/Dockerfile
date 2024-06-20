# 第一階段: 編譯階段
FROM golang:1.22 AS builder

WORKDIR /app

COPY . .

RUN make build-linux

# 第二階段: 執行階段

FROM alpine:3.14

WORKDIR /app

COPY --from=builder /app/bin/aid /app/bin/aid

# 如果你的應用程序需要任何其他依賴項,請在此處安裝
ENV GO_ENV=release
# set Timezone as cloud run always run in UTC + 0
RUN apk add --no-cache tzdata
ENV TZ=Asia/Taipei

EXPOSE 8080

CMD ["./bin/aid", "server"]