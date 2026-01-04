REST 

Что используется:
1. Golang
2. Chi
3. Postgres
4. Docker
5. Goose
6. Slog
7. 
-----------------------------------------------------

Запуск:
1. 

docker compose build --no-cache
docker compose up
-----------------------------------------------------


Проверка:

docker compose logs -f api
curl http://localhost:8080/api/v1/users


Что дальше (следующие шаги):

1. Добавить миграции (golang-migrate)

2. Context timeouts

3. Structured logging (zap / slog)

4. Graceful shutdown

5. CRUD операции

Dependency injection

OpenAPI / Swagger

Тесты (unit + integration)