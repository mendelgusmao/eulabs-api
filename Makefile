all: run

build:
	docker build . -t eulabs-api:latest

run: build
	docker-compose up 

stop:
	docker-compose down
