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
	GOOS=darwin GOARCH=amd64 go build -o ./bin/osx_urlshortener .
	GOOS=linux GOARCH=amd64 go build -o ./bin/urlshortener .
	GOOS=windows GOARCH=386 go build -o ./bin/urlshortener.exe .
	cp config.* ./bin
	cp -rf public ./bin

run:
	cd ./bin && ./urlshortener

run_osx:
	cd ./bin && ./osx_urlshortener

run_win:
	cd ./bin && ./urlshortener.exe

dev:
	go run main.go

test:
	go test ./...
