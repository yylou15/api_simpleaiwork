# -------- build stage --------
FROM golang:1.24-alpine AS builder

WORKDIR /app

# 先拷依赖文件，利用缓存
COPY go.mod go.sum ./
RUN go mod download

# 再拷源码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./

# -------- runtime stage --------
FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/app /app/app

EXPOSE 8080
CMD ["/app/app"]
