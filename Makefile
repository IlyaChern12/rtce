# Makefile

# подъем всего
up:
	docker compose -f docker/docker_compose.yml up -d

# остановка всего
down:
	docker compose -f docker/docker_compose.yml down

# локальный запуск
run:
	go run ./cmd/editor-server

# тесты
test:
	go test ./...

# линт
lint:
	golangci-lint run