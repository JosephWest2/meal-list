dev: build
	doppler run -- air -c .air.toml

install:
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

build:
	sqlc generate
	templ generate
	go build 
