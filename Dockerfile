# Базовый образ с Go
FROM golang:1.24-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Устанавливаем protoc и инструменты
RUN apk add --no-cache protobuf
RUN go install github.com/air-verse/air@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Копируем go.mod и go.sum, подтягиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Указываем порты
EXPOSE 8080 50051

# Запускаем через Air
CMD ["air"]