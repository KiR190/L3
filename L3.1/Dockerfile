FROM golang:1.25.3-alpine

WORKDIR /app

# 1. Копируем только файлы с зависимостями
COPY go.mod go.sum ./

# 2. Скачиваем зависимости (будет кэшироваться, если go.mod/go.sum не менялись)
RUN go mod download -x

# 3. Копируем весь код проекта
COPY ./ ./

# 4. Сборка бинарника
RUN go build -o main ./cmd/

# 5. Запуск приложения
CMD ["./main"]