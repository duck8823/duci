FROM golang:1.10-alpine
MAINTAINER shunsuke maeda <duck8823@gmail.com>

RUN apk --update add --no-cache alpine-sdk

WORKDIR /go/src/github.com/duck8823/duci

ADD . .

RUN go get golang.org/x/vgo

ENV CC=gcc

RUN vgo install

EXPOSE 8080

CMD ["duci"]