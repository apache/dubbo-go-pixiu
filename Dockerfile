### builder
FROM golang:alpine as builder
LABEL MAINTAINER="2677759629@qq.com"

ENV GOPROXY="https://goproxy.cn,direct" \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
WORKDIR /apps/
COPY . .
RUN go mod download 
RUN go build -ldflags "-s -w"  -o pixiu ./cmd/pixiu/*.go


### alpine
FROM amd64/ubuntu:latest
COPY --from=builder /apps/pixiu /app/
COPY --from=builder /apps/docker-entrypoint.sh /app/
COPY --from=builder /apps/conf/* /etc/pixiu/
WORKDIR /app/
RUN ["chmod", "+x","./docker-entrypoint.sh"]
ENTRYPOINT ["./docker-entrypoint.sh"]
