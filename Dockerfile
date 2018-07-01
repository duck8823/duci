FROM golang:1.10-alpine
MAINTAINER shunsuke maeda <duck8823@gmail.com>

RUN apk --update add --no-cache git

WORKDIR /go/src/github.com/duck8823/minimal-ci

ADD . .

RUN go get golang.org/x/vgo

RUN vgo install

EXPOSE 8080

CMD ["minimal-ci"]