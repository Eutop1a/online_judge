FROM golang:latest
LABEL authors="eutop1a"

# 设置环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    MYSQL_HOST=172.17.0.1 MYSQL_PORT=33061 MYSQL_USER=root MYSQL_PASSWORD=123456 MYSQL_DBNAME=OnlineJudge \
    REDIS_HOST=172.17.0.1 REDIS_PORT=63790 REDIS_PASSWORD=0 REDIS_DB=0

# 设置工作目录
WORKDIR /go/src/app

COPY . .

EXPOSE 65533

RUN go mod download

# 构建 Go 项目
RUN go build -o OnlineJudge .


CMD ["./OnlineJudge"]

