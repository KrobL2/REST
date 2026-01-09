# Базовый образ для сборки
FROM golang:1.24-alpine AS builder 

# Рабочий каталог внутри контейнера
WORKDIR /app

# Копирует файлы зависимостей
COPY go.mod go.sum ./

# Загружает зависимости
RUN go mod download

# Копирует весь исходный код в контейнер
COPY . .

# Компилирует приложение
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

FROM gcr.io/distroless/base-debian12
WORKDIR /app

COPY --from=builder /app/server .
COPY ./migrations ./migrations

# Порт который будет прослушивать контейнер во время выполнения
EXPOSE 8080

# Команда, которая будет выполняться при запуске контейнера
CMD ["./server"]
