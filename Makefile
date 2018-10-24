# Makefile
build:
	go build

cross-build:
	GOOS=windows go build
	GOOS=linux go build
	GOOS=darwin go build

test:
	go test -coverprofile cover.out $$(go list ./... | grep -v mock_)
	go tool cover -html cover.out -o cover.html
	open cover.html

docker-test:
	docker build -t duck8823/duci:test -f .duci/Dockerfile .
	docker run --rm \
	           -v ${GOPATH}/pkg/mod/cache:/go/pkg/mod/cache \
	           duck8823/duci:test

clean:
	rm -f duci duci.exe go.sum cover.out cover.html
