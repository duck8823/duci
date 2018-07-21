FROM golang:1.10-alpine
MAINTAINER shunsuke maeda <duck8823@gmail.com>

RUN apk --update add --no-cache alpine-sdk
RUN adduser duci --disabled-password --home /home/duci

WORKDIR /go/src/github.com/duck8823/duci

ADD . .

RUN chown -R duci:duci /go/src

USER duci
ENV HOME=/home/duci

RUN go get golang.org/x/vgo

ENV CC=gcc
ENV DOCKER_HOST=tcp://host.docker.internal:2375

ENTRYPOINT ["vgo"]
CMD ["test", "./..."]