#Basic makefile

default: build

build: vet
	@go generate ./... && go build

doc:
	@godoc -http=:6060 -index

lint:
	@golint ./...

debug:
	@reflex -c reflex.conf

run: build
	./ginblog

test:
	@go test ./...

vet:
	@go vet ./...

clean:
	@rm -f ./blog
	@rm -f ./system/*.rice-box.go
