### builder
FROM golang:alpine as builder
LABEL MAINTAINER="2677759629@qq.com"

ENV GOPROXY="https://goproxy.cn,direct" \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
WORKDIR /app/
COPY . .
RUN go mod download
RUN go build -ldflags "-s -w" -o dubbo-go-pixiu ./cmd/pixiu/*.go


### alpine
FROM amd64/ubuntu:latest
COPY --from=builder /app/dubbo-go-pixiu /app/
COPY --from=builder /app/docker-entrypoint.sh /app/
COPY --from=builder /app/configs/* /etc/pixiu/
WORKDIR /app/
RUN ["chmod", "+x","./docker-entrypoint.sh"]
ENTRYPOINT ["./docker-entrypoint.sh"]