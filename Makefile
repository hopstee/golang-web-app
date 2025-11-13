APP_BINPATH = bin/app
SCRIPT_BINPATH = bin/scripts
WEBAPP_DIR = webapp
DOCKER_IMAGE = golang_web_app:latest

.PHONY: all
all: build

.PHONY: build
build: build-templ build-app build-webapp gen-schemas

.PHONY: build-app
build-app:
	go build -o ${APP_BINPATH} cmd/app/main.go

.PHONY: build-templ
build-templ:
	templ generate

.PHONY: build-webapp
build-webapp:
	cd ${WEBAPP_DIR} && npm ci && npm run build

.PHONY: run
run: build
	${APP_BINPATH}

.PHONY: watch
watch:
	${MAKE} -j3 watch-app watch-templ watch-webapp

.PHONY: watch-app
watch-app:
	go run github.com/air-verse/air@latest \
	--build.cmd "${MAKE} build-app" \
	--build.bin "${APP_BINPATH}" \
	--build.include_ext "go" \
	--build.exclude_dir "bin, data, webapp"

.PHONY: watch-templ
watch-templ:
	templ generate \
	--watch \
	--proxy="http://localhost:8080" \
	--open-browser=false

.PHONY: watch-webapp
watch-webapp:
	cd ${WEBAPP_DIR} && npm ci && npm run dev

.PHONY: build-scripts
build-scripts:
	go build -o ${SCRIPT_BINPATH} cmd/scripts/main.go

.PHONY: add-admin
add-admin: build-scripts
	${SCRIPT_BINPATH} create

.PHONY: upsert-admin
upsert-admin: build-scripts
	${SCRIPT_BINPATH} create --force

.PHONY: gen-schemas
gen-schemas: build-scripts
	${SCRIPT_BINPATH} generate

.PHONY: docker-build
docker-build: build
	docker build --no-cache -f docker/Dockerfile -t ${DOCKER_IMAGE} .

.PHONY: docker-run
docker-run:
	docker run --rm -it \
		-p 8080:8080 \
		-v $(PWD)/uploads:/app/uploads \
		-v $(PWD)/config/config.yml:/app/config/config.yml:ro \
		${DOCKER_IMAGE}

.PHONY: docker-up
docker-up: docker-build
	docker-compose -f docker/docker-compose.yml --env-file .env up -d

.PHONY: docker-down
docker-down:
	docker-compose -f docker/docker-compose.yml --env-file .env down -v

.PHONY: docker-dev-up
docker-dev-up: docker-build
	docker-compose -f docker/docker-compose.dev.yml --env-file .env up -d

.PHONY: docker-dev-down
docker-dev-down:
	docker-compose -f docker/docker-compose.dev.yml --env-file .env down -v