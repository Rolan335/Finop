.PHONY: generate run test

generate:
	go generate cmd/main.go

run:
	docker compose up

test:
	go test -v ./test