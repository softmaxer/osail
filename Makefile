BINARY=osail

build:
	@templ fmt .
	@templ generate
	@GOARCH=amd64 GOOS=darwin go build -o ./target/${BINARY}-darwin-amd64 .
	@GOARCH=arm64 GOOS=darwin go build -o ./target/${BINARY}-darwin-arm64 .
	@GOARCH=amd64 GOOS=linux go build -o ./target/${BINARY}-linux-amd64 .
	@GOARCH=amd64 GOOS=windows go build -o ./target/${BINARY}-windows-amd64 .

run: build
		@./target/${BINARY}-darwin-arm64

clean:
	rm -r ./views/*.go
	rm -r ./views/layout/*.go
	rm -r ./views/styles/*.go
	rm -r ./target/*
