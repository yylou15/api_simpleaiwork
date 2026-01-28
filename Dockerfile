# ---------- build stage ----------
FROM golang:1.24-alpine AS builder

WORKDIR /app

# 先拷 go.mod / go.sum，利用缓存
COPY go.mod go.sum ./
RUN go mod download

# 再拷源码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o server ./cmd/server
# ↑ 如果你不是 ./cmd/server，改成你的 main 包路径

# ---------- runtime stage ----------
FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/server /app/server

EXPOSE 8080
CMD ["/app/server"]
