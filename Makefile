# Makefile
build:
	vgo build

test:
	vgo test -coverprofile cover.out $$(vgo list ./... | grep -v mock_)
	vgo tool cover -html cover.out -o cover.html
	open cover.html

docker-test:
	docker build -t duck8823/duci:test -f .duci/Dockerfile .
	docker run --rm \
	           -v /var/run/docker.sock:/var/run/docker.sock \
	           -v ${HOME}/.ssh:/root/.ssh:ro \
	           -v ${GOPATH}/pkg/mod/cache:/go/pkg/mod/cache \
	           duck8823/duci:test

clean:
	rm -f duci go.sum cover.out cover.html
