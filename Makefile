build:
	@go build -o bin/hotel_rezerv

run: build
	@./bin/hotel_rezerv

seed:
	@go run scripts/seed.go

docker:
	@docker build -t sadagatasgarov/hotel-rezerv .
	@docker push sadagatasgarov/hotel-rezerv
	@docker compose up

test:
	@go test -v ./... -count=1