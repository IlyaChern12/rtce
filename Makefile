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

# миграции
# создание таблицы юзеров
migrate-up:
	migrate -path internal/db/migrations -database "postgres://rtce:rtcepass@localhost:5432/rtce_dev?sslmode=disable" up

# cброс таблицы юзеров
migrate-down:
	migrate -path internal/db/migrations -database "postgres://rtce:rtcepass@localhost:5432/rtce_dev?sslmode=disable" down