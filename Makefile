.PHONY generate run test

generate:
	go generate cmd/main.go

run:
	docker compose up

test:
	docker compose up
	go test -v ./test