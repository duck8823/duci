FROM golang:1.11-alpine AS build
MAINTAINER shunsuke maeda <duck8823@gmail.com>

RUN apk --update add --no-cache alpine-sdk

WORKDIR /go/src/github.com/duck8823/duci

ADD . .

ENV GO111MODULE=on

RUN go build

FROM alpine

WORKDIR /root/
COPY --from=build /go/src/github.com/duck8823/duci/duci .

EXPOSE 8080

CMD ["./duci"]