FROM golang:1.11.4-alpine
MAINTAINER shunsuke maeda <duck8823@gmail.com>

RUN apk --update add --no-cache alpine-sdk

WORKDIR /go/src/github.com/duck8823/duci

ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

ENTRYPOINT ["make"]
CMD ["test"]