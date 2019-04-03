default: build

start:
	docker-compose up -d

stop:
	docker-compose down

mongo:
	docker-compose up -d db

build:
	docker build -t roboncode/shorturls .

