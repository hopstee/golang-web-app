# Web app boilerplate

Light weight project.
Contains base architecture for API and web pages.

- Repositories (PostgreSQL)
- Services (buisnes logic)
- Handlers (REST API and Web)
- KVStore (cache - Redis)
- Notifications (now Telegram, you can add others)
- JWT-authentication
- Chi routing
- Static fails with [templ](https://templ.guide/) and [htmx](https://htmx.org/)

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
  base_url: "http://localhost:8080"

database:
  driver: "postgres"
  dsn: "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

kvstore:
  type: "redis"
  redis:
    addr: "${REDIS_HOST}:${REDIS_PORT}"
    password: "${REDIS_PASSWORD}"
    db: ${REDIS_DB}

filestorage:
  type: "local"
  local:
    base_path: "./static/uploads"

authentication:
  jwt:
    secret: "${JWT_SECRET}"

telegram:
  token: "${TELEGRAM_TOKEN}"
  chat_id: "${TELEGRAM_CHAT_ID}"

static:
  dir: "./static"

schemas:
  schema: "./internal/schemas/schema.v1.json"
  pages: "./internal/schemas/pages.v1.json"
  layouts: "./internal/schemas/layouts.v1.json"
  blocks: "./internal/schemas/blocks.v1.json"
  modules: "./internal/schemas/modules.v1.json"
  shared: "./internal/schemas/shared.v1.json"
```

- .env
```bash
JWT_SECRET=very-secret-secret

TELEGRAM_TOKEN=tg-token
TELEGRAM_CHAT_ID=chat-id

DB_HOST=localhost
DB_PORT=5432
DB_USER=appuser
DB_PASSWORD=apppass
DB_NAME=appdb

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=redispass
REDIS_DB=0
```

3. Run app in dev mode
```bash
make watch
```

---

## Add admin
```bash
make add-admin
```

---

## Generate schemas
```bash
make get-schemas
```

---

## Static entities schema generator

### Declaration

| Directive                        | Description                                                   |
| ---------------------------------| --------------------------------------------------------------|
| `@type: page`                    | Entity type (layout, block, page, module, shared)             |
| `@id: id`                        | Page ID (index)                                               |
| `@title: text`                   | Readable title                                                |
| `@layout: id`                    | Connect layout to the page                                    |
| `@blocks: block1, block2, ...`   | List of blocks for this page                                  |
| `@shared: shared1, shared2, ...` | List of shared components for this page (e.g. contacts block) |
| `@field ...`                     | Fields description (with support for nesting and `[]`)        |

---

## License

MIT License.