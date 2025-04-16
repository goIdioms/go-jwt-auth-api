.PHONY: build dev clean deps run

deps:
	go mod tidy
	go mod vendor

build: deps
	docker-compose build

dev: deps
	docker-compose up -d
	air

run:
	go run main.go

clean:
	docker-compose down
	rm -rf tmp/main tmp/build-errors.log
	go clean