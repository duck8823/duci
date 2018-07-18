# Makefile
build:
	vgo build

test:
	vgo test -coverprofile cover.out $$(vgo list ./... | grep -v mocks)
	vgo tool cover -html cover.out -o cover.html
	open cover.html

clean:
	rm -f duci go.sum cover.out cover.html
