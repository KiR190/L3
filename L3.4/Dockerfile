FROM golang:1.25-alpine

WORKDIR /app

# Устанавливаем пакет с системными шрифтами (опционально, если нужен DejaVu)
RUN apk add --no-cache ttf-dejavu

# 1. Копируем только файлы с зависимостями
COPY go.mod go.sum ./

# 2. Скачиваем зависимости (будет кэшироваться, если go.mod/go.sum не менялись)
RUN go mod download -x

# 3. Копируем весь код проекта
COPY ./ ./

# 4. Создаём папку assets и копируем Roboto-Bold.ttf внутрь
RUN mkdir -p assets
COPY Roboto-Bold.ttf assets/Roboto-Bold.ttf

# 5. Сборка бинарника
RUN go build -o main ./cmd/

# 6. Запуск приложения
CMD ["./main"]