FROM golang:1.12.6-alpine
MAINTAINER shunsuke maeda <duck8823@gmail.com>

RUN apk --update add --no-cache alpine-sdk

WORKDIR /workdir

# install golangci-lint
RUN curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

ENTRYPOINT ["make"]
CMD ["test"]