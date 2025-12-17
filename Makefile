dev: build
	doppler run -- air -c .air.toml

install:
	npm install
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

build:
	npx @tailwindcss/cli -i ./static/css/tailwind_input.css -o ./static/css/tailwind_output.css
	sqlc generate
	templ generate
	go build 
