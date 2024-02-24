# Test task BackDev

**Используемые технологии:**

- Go
- JWT
- MongoDB

## Запуск

Клонируем репозиторий

```sh
git clone https://github.com/v7ktory/test-exercise.git
cd test
```

команда make собирает приложение и запускает в докер контейнере
команда make stop останавливает приложение и удаляет контейнеры

```sh
make
make stop
```

если не хотим запускать приложение в докере, тогда в .env меняем MONGO_HOSTS="mongodb" на MONGO_HOSTS="localhost" и пишем команду

```sh
go run ./cmd/app/main.go
```

## Проверка приложения через постман

**В body вводим email и password и получаем userID**
![Screenshot](screenshot/screen1.png)
**В параметре запроса вводим полученный userID и в body указываем наши email и password**
**Затем в header получаем access токен**
![Screenshot](screenshot/screen2.png)
**В параметре запроса вводим полученный userID**
**Копируем полученный access токен и вставляем его в header Authorization**
![Screenshot](screenshot/screen3.png)
