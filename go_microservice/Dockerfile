FROM golang:1.17-alpine

# Устанавливаем переменные окружения
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы зависимостей и устанавливаем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем все файлы проекта в контейнер
COPY . .

# Собираем приложение Go
RUN go build -o main .

# Запускаем приложение по умолчанию
CMD ["./main"]