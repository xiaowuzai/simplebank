# Build stage
FROM golang:1.21-alpine3.19 AS builder
WORKDIR /app
COPY . .
# ENV GOPROXY https://goproxy.cn
RUN go build -o main .

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .

EXPOSE 8080
CMD [ "/app/main" ]