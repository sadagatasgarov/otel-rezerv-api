build:
	@go build -o bin/hotel_rezerv

run: build
	@./bin/hotel_rezerv

seed:
	@go run scripts/seed.go

test:
	@go test -v ./...