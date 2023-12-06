include .env
export

run:
	go run ./src/api/*.go

migrate-up:
	 migrate -path=./migrations -database="$(DATABASE_URL)" up

migrate-down:
	 migrate -path=./migrations -database="$(DATABASE_URL)" down 

