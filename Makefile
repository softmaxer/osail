build:
	@templ fmt .
	@templ generate
	@go build .

run: build
		@./osail
	
clean:
	rm ./osail
	rm test.db
	touch test.db
