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


<pre> ```bash docker compose build --no-cache ``` </pre>


<pre> ```docker compose up ``` </pre>
-----------------------------------------------------


Проверка:

curl http://localhost:8080/api/v1/users


Что дальше (следующие шаги):

2. Context timeouts

3. Structured logging (zap / slog)

4. Graceful shutdown

5. CRUD операции

Dependency injection

OpenAPI / Swagger

Тесты (unit + integration)