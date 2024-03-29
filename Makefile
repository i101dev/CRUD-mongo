build:
	@go build -o bin/api

run: build
	@./bin/api

seed:
	@go run scripts/seed.go

test:
	@go test -v ./...


docker:
	echo "building docker file"
	@docker build -t hotel-api .
	echo "running API inside Docker container"
	@docker run -p 5000:5000 hotel-api