FROM golang:1.22-alpine

RUN set -ex; \
    apk update

WORKDIR /app

COPY ./ /app

RUN go build -o /output ./cmd/app

EXPOSE 3000

CMD ["/output"]