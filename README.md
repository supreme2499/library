# Library

___
Сервис принимает запросы на 8000 порту
### Установка
1. Клонируйте репозиторий
```sh
https://github.com/supreme2499/library.git
cd library
```
2. Настройте .env файл в корне проекта(пример):
```
ENV=local

POSTGRES_STORAGE_URL=postgres://postgres:1234@postgres:5432/music?sslmode=disable
POSTGRES_MIGRATIONS_PATH=migrations

HTTP_SERVER_ADDRESS=0.0.0.0:8000
HTTP_SERVER_TIMEOUT=4s
HTTP_SERVER_IDLE_TIMEOUT=60s
HTTP_SERVER_WITH_TIMEOUT=10s

INFO_URL="http://host.docker.internal:8082"
```
3. Запустите сервис
```
docker-compose build
docker-compose up -d
```
___
### Документация

***Swagger:*** http://localhost:8000/swagger/index.html#/ (при включенном приложении)