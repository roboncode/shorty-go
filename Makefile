default: build

start:
	docker-compose down
	docker-compose up -d

stop:
	docker-compose down

mongo:
	docker-compose up -d db

build:
	docker build -t roboncode/shorty .

standalone:
	GOOS=darwin GOARCH=amd64 go build -o ./bin/shorty_osx .
	GOOS=linux GOARCH=amd64 go build -o ./bin/shorty .
	GOOS=windows GOARCH=386 go build -o ./bin/shorty.exe .
	cp config.* ./bin
	cp -rf public ./bin

run:
	cd ./bin && ./shorty

run_osx:
	cd ./bin && ./shorty_osx

run_win:
	cd ./bin && ./shorty.exe

dev:
	go run main.go

test:
	go test ./...
