.DEFAULT_GOAL := test

GOPATH=$(shell pwd)/vendor:$(shell pwd)

clean:
	@echo "Clean ./bin"
	rm -rf bin pkg *.out

get: clean
	@echo "Get..."
	@GOPATH=$(GOPATH) go get github.com/hashicorp/terraform/plugin

build: get
	@echo "Build..."
	@GOPATH=$(GOPATH) go build -o bin/terraform-provider-keycloak -tags netgo

install: build
	cp bin/* ~/.terraform.d/plugins/

test:
	go test ./... -v
vet: 
	go vet ./...

stest: vet
	go test ./... -short
