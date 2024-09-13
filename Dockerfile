FROM golang:1.23-alpine AS builder

ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY .. .

RUN go build -o /go/bin/go-pandoc ./

FROM registry.cn-chengdu.aliyuncs.com/custom-gq/libreoffice-pandoc:latest

ENV PANDOC_DEFAULT_DATA_DIR=/app/data

WORKDIR /app

COPY --from=builder /go/bin/go-pandoc /usr/local/bin/go-pandoc

COPY ./manifest /app/manifest

EXPOSE 80

ENTRYPOINT ["/usr/local/bin/go-pandoc"]