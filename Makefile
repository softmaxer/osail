build:
	@templ fmt .
	@templ generate
	@go build .

run: build
		@./localflow
	
clean:
	rm ./localflow
	rm test.db
	touch test.db
