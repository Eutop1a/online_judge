# 基于 golang 官方镜像
FROM golang:1.22-alpine AS builder

# 设置工作目录
WORKDIR /app

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 将 go.mod 和 go.sum 文件复制到容器中
COPY go.mod .
COPY go.sum .

# 下载依赖
RUN go mod download && go mod tidy


# 将项目文件添加到容器中
COPY . .

# 编译项目
RUN chmod 777 ./wait-for-it.sh &&  \
       go build -o judgement ./app/judgement/cmd/main.go && \
       go build -o online-judge main.go

# 基于 centos 官方镜像
FROM centos:latest
#FROM scratch

# 设置工作目录
WORKDIR /app

# 从builder镜像中把配置文件拷贝到当前目录
COPY --from=builder /app/conf /app/conf

# 拷贝可执行文件
COPY --from=builder /app/online-judge /app/online-judge
COPY --from=builder /app/judgement /app/judgement
COPY --from=builder /app/wait-for-it.sh /app/wait-for-it.sh

# 假设你的项目需要运行在 65533 端口
EXPOSE 65533
EXPOSE 8082

# 启动你的应用程序
CMD ["sh", "-c", "./judgement & ./online-judge"]


