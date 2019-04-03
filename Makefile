default: build

start:
	docker-compose up -d

stop:
	docker-compose down

build:
	docker build -t roboncode/shorturls .