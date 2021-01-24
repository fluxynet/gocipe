COMMIT := `git rev-parse HEAD`
DATE := `date +%FT%T%z`

.PHONY: build

build: clean
	@go generate ./...
	@env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -extldflags '-static' -X=main.appCommit=$(COMMIT) -X=main.appBuilt=$(DATE)" -o build/gocipe github.com/fluxynet/gocipe/cmd

clean:
	@rm -f build/gocipe
	@rm -f rice-box.go
