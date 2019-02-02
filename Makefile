.DEFAULT_GOAL := test

clean:
	@echo "Clean ./bin"
	rm -rf bin pkg *.out

build:
	@echo "Build..."
	go build -o bin/terraform-provider-keycloak -tags netgo

install: build
	cp bin/* ~/.terraform.d/plugins/

test:
	go test ./... -v
vet: 
	go vet ./...

stest: vet
	go test ./... -short
