FROM golang:1.12.6-alpine AS build
MAINTAINER shunsuke maeda <duck8823@gmail.com>

RUN apk --update add --no-cache alpine-sdk ca-certificates \
 && update-ca-certificates

WORKDIR /workdir

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 make build

FROM scratch

WORKDIR /workdir

COPY --from=build /workdir/duci /usr/local/bin/duci
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080

ENTRYPOINT ["duci"]
CMD ["server"]