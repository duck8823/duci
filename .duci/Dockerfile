FROM golang:1.11.4-alpine
MAINTAINER shunsuke maeda <duck8823@gmail.com>

RUN apk --update add --no-cache alpine-sdk

WORKDIR /go/src/github.com/duck8823/duci

ADD . .

ENV CC=gcc
ENV CI=duci

ENV GO111MODULE=on

ENTRYPOINT ["make"]
CMD ["test"]