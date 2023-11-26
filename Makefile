build:
	@go build -o bin/hotel_rezerv

run: build
	@./bin/hotel_rezerv

test:
	@go test -v ./...