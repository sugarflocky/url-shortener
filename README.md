# url-shortener

Сервис коротких ссылок: принимает URL, возвращает код из 10 символов,
по коду редиректит на оригинальный адрес.

## Запуск

Через docker compose (приложение + postgres):

    docker compose up --build

Docker-образ один, хранилище выбирается аргументами при запуске:

    docker build -t url-shortener .
    docker run --rm -p 8080:8080 url-shortener                                    # память
    docker run --rm -p 8080:8080 url-shortener -storage=postgres -dsn="..."       # postgres

Или локально, без докера:

    go run ./cmd/shortener -storage=memory

Флаги:

    -storage  memory | postgres (по умолчанию memory)
    -dsn      строка подключения к postgres, обязательна при -storage=postgres
    -addr     адрес сервера, по умолчанию :8080

## API

Создать короткую ссылку:

    curl -X POST localhost:8080/shorten -d '{"url":"https://go.dev/"}'
    # {"code":"aB3_x9Kq2Z"}

Перейти по ней:

    curl -v localhost:8080/aB3_x9Kq2Z
    # HTTP 301, Location: https://go.dev/

Ошибки: 400 (кривой JSON или URL), 404 (неизвестный код), 500 (внутренние).

## Тесты

    go test ./...

Интеграционные тесты postgres запускаются только при заданной переменной
TEST_POSTGRES_DSN, иначе пропускаются:

    $env:TEST_POSTGRES_DSN = "postgres://postgres:devpass@localhost:5432/shortener"
    go test ./internal/storage/postgres/