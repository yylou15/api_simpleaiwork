# ---------- build stage ----------
FROM golang:1.24-alpine AS builder

WORKDIR /app

# 先拷 go.mod / go.sum，利用缓存
COPY go.mod go.sum ./
RUN go mod download

# 再拷源码
COPY . .

# 编译当前目录的 main 包
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app .

# ---------- runtime stage ----------
FROM alpine:3.20

WORKDIR /app

# 基础运行依赖（HTTPS / 时区）
RUN apk add --no-cache ca-certificates tzdata

# 拷贝可执行文件
COPY --from=builder /app/app /app/app

EXPOSE 8080
CMD ["/app/app"]
