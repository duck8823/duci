FROM golang:1.11.4-alpine
MAINTAINER shunsuke maeda <duck8823@gmail.com>

RUN apk --update add --no-cache alpine-sdk

WORKDIR /workdir

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

ENTRYPOINT ["make"]
CMD ["test"]