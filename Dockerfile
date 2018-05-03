FROM golang:1.10-alpine
MAINTAINER shunsuke maeda <duck8823@gmail.com>

RUN apk --update add --no-cache git

WORKDIR /go/src/github.com/duck8823/webhook-proxy

ADD . .

RUN go get -u github.com/golang/dep/cmd/dep && \
    dep ensure

RUN go build

EXPOSE 8080

CMD ["./webhook-proxy"]