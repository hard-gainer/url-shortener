docker-up:
	docker-compose up -d

docker-down:
	docker-compose down -v

test-coverage:
	go test -v -cover ./...

.PHONY: test-coverage docker-up docker-down