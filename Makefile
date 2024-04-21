build:
	@templ fmt .
	@templ generate
	@go build .

run: build
		@./osail
	
clean:
	rm ./localflow
	rm test.db
	touch test.db
