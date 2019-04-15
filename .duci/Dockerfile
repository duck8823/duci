FROM golang:1.12.4-alpine
MAINTAINER shunsuke maeda <duck8823@gmail.com>

RUN apk --update add --no-cache alpine-sdk

WORKDIR /workdir

# install gometalinter for lint
RUN curl -s -L https://git.io/vp6lP | sh -s -- -b /usr/local/bin

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

ENTRYPOINT ["make"]
CMD ["test"]