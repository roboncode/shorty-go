default: build

start:
	docker-compose down
	docker-compose up -d

stop:
	docker-compose down

mongo:
	docker-compose up -d db

build:
	docker build -t roboncode/urlshortener .

standalone:
	go build -o ./bin/urlshortener .
	cp config.* bin
	cp -rf public bin

run:
	cd ./bin && ./urlshortener

dev:
	go run main.go

test:
	go test ./...
