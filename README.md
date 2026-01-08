REST 

-----------------------------------------------------

Что используется:
1. Golang
2. Slog
3. Chi
4. Goose
5. cleanenv
6. Postgres
7. Docker 

-----------------------------------------------------
Билд:

<pre> ```bash docker compose build --no-cache ``` </pre>


Запуск:

<pre> ```docker compose up ``` </pre>
-----------------------------------------------------

Проверка:

curl http://localhost:8080/api/v1/users



------------------------------------------------------------

Что дальше (следующие шаги):

2. Context timeouts

4. Graceful shutdown

5. CRUD операции

Dependency injection

OpenAPI / Swagger

Тесты (unit + integration)