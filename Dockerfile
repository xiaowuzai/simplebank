# Build stage
FROM golang:1.21-alpine3.19 AS builder
WORKDIR /app
COPY . .
# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN GOPROXY='https://goproxy.cn,direct' go build -o main .
# RUN apk add git curl
# RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
# COPY --from=builder /app/migrate .
COPY app.env .
COPY start.sh .
COPY db/migration ./db/migration

EXPOSE 8080
CMD [ "/app/main" ]
#CMD 作为参数
ENTRYPOINT [ "/app/start.sh" ] 