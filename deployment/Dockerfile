# This assumes you're using a multi-stage build with Go
# Start from the official Golang image to build our source code.
FROM golang:1.22 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum file
COPY go.mod go.sum ./

# Download all the dependencies
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o memorianexus ./cmd
# RUN go build -o memorianexus ./cmd
# alpine 缺少动态链接库，编译时要禁止 cgo，如果依赖 c，需要在 Alpine 容器中安装 libc6-compat 或 glibc 兼容库，或者换下镜像

# Start a new stage from scratch for a lightweight final image
FROM alpine:latest

WORKDIR /app/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/memorianexus /app/

RUN chmod +x /app/memorianexus

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD /app/memorianexus

# 在 docker-compose 时，配置了 network 但仍然找不到 mysql host 的情况下，可以使用 dockerize 等待 MySQL 服务启动
# 但这样 build 时间比较长，而且并非正式环境使用，需要单独构建。稍微简单一点，用 makefile 好了
#RUN apk add --no-cache openssl
#ENV DOCKERIZE_VERSION v0.6.1
#RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
#    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
#    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz
#CMD dockerize -wait tcp://mysql:3306 -timeout 30s /app/memorianexus