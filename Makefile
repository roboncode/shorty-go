default: build

start:
	docker-compose up -d

stop:
	docker-compose down

mongo:
	docker-compose up -d db

build:
	docker build -t roboncode/urlshortener .

standalone:
	go build -o ./bin/urlshortener .

run:
	./bin/urlshortener