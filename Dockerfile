FROM golang:1.19-alpine

RUN set -ex; \
    apk update

COPY src ./src

WORKDIR src/rest-api

RUN go build -o /output

EXPOSE 3000

CMD ["/output"]