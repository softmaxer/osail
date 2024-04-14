build:
	@templ fmt .
	@templ generate
	@go build .

run: build
		@./localflow
