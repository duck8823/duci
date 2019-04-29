# Makefile
build:
	go build -ldflags "-s -w -X github.com/duck8823/duci/application.version=$$(git describe --tags)"

cross-build:
	for os in darwin linux windows; do \
	    for arch in amd64 386; do \
	        GOOS=$$os GOARCH=$$arch go build \
	          -ldflags "-s -w -X github.com/duck8823/duci/application.version=$$(git describe --tags)" \
	          -o dist/duci_$$os_$$arch; \
	    done; \
	done

docker-build:
	docker build -t duck8823/duci:$$(git describe --tags) .

lint:
	golangci-lint run \
      --disable-all \
      --enable=gofmt \
      --enable=vet \
      --enable=gocyclo \
      --enable=golint \
      --enable=ineffassign \
      --enable=misspell \
      --deadline=5m

test:
	go test -coverprofile cover.out $$(go list ./... | grep -v mock_)
	go tool cover -html cover.out -o cover.html

test-in-docker:
	docker build -t duck8823/duci:test -f .duci/Dockerfile .
	docker run --rm \
	           duck8823/duci:test test

clean:
	rm -fr duci duci.exe cover.out cover.html dist
