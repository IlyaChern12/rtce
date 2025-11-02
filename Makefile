# Makefile

# подъем всего
up:
	docker compose -f docker/docker_compose.yml up -d $(BUILD)

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
	docker run --rm -v "$(PWD)/internal/db/migrations":/migrations \
	--network=docker_rtce_net \
	migrate/migrate \
	-path=/migrations/ \
	-database "postgres://rtce:rtcepass@db:5432/rtce_dev?sslmode=disable" up

# откат миграций
migrate-down:
	echo "y" | docker run --rm -i -v "$(PWD)/internal/db/migrations":/migrations \
	--network=docker_rtce_net \
	migrate/migrate \
	-path=/migrations/ \
	-database "postgres://rtce:rtcepass@db:5432/rtce_dev?sslmode=disable" down