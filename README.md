# Web app boilerplate

Light weight project.
Contains base architecture for API and web pages.

- Repositories (now SQLite, you can add PostgreSQL/MySQL)
- SErvices (buisnes logic)
- Handlers (REST API and Web)
- Notifications (now Telegram, you can add others)
- JWT-authentication
- Chi routing
- Static fails with templ and htmx

---

## Start project
1. Clone repo

```bash
git clone https://github.com/hopstee/golang-web-app.git
cd golang-web-app
```

2. Configuration

- config/config.yml
```yml
server:
  addr: ":8080"

database:
  driver: "sqlite"
  dsn: "./data/app.db"

authentication:
  jwt:
    secret: "${JWT_SECRET}"

telegram:
  token: "${TELEGRAM_TOKEN}"
  chat_id: "${TELEGRAM_CHAT_ID}"

static:
  dir: "./static"
```

- .env
```bash
JWT_SECRET=very-secret-secret

TELEGRAM_TOKEN=tg-token
TELEGRAM_CHAT_ID=chat-id
```

3. Run server
```bash
make run
```

---

## Add admin
```bash
make add-admin
```

---

## How to change SQLite on PostgreSQL / MySQL

1. Change `database.driver` and `database.dsn` in `config/config.yml`
2. Add handler in `internal/infrastructure/db.go` for new database
3. Add realisation for new database in `internal/repository/`
4. Add repos initialisation in `internal/infrastructure/repos.go`

---

## License

MIT License.