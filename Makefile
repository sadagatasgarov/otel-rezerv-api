build:
	@go build -o bin/hotel_rezerv

run: build
	@./bin/hotel_rezerv

seed:
	@go run scripts/seed.go

docker:
	@docker build -t sadagatasgarov/hotel-rezerv:latest .
	@docker push sadagatasgarov/hotel-rezerv:latest
	@docker compose up -d

test:
	@go test -v ./... -count=1


git:
	@git add .
	@git commit -m"Duzelis"
	@git push -u origin hazir