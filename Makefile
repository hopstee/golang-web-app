APP_BINPATH = bin/app
SCRIPT_BINPATH = bin/scripts

.PHONY: build
build: build-templ build-app

.PHONY: build-app
build-app:
	go build -o ${APP_BINPATH} cmd/app/main.go

.PHONY: build-templ
build-templ:
	templ generate

.PHONY: run
run: build
	${APP_BINPATH}

.PHONY: watch
watch:
	${MAKE} -j2 watch-app watch-templ

.PHONY: watch-app
watch-app:
	go run github.com/air-verse/air@latest \
	--build.cmd "${MAKE} build-app" \
	--build.bin "${APP_BINPATH}" \
	--build.include_ext "go" \
	--build.exclude_dir "bin, data"

.PHONY: watch-templ
watch-templ:
	templ generate \
	--watch \
	--proxy="http://localhost:8080" \
	--open-browser=false

.PHONY: build-scripts
build-scripts:
	go build -o ${SCRIPT_BINPATH} cmd/scripts/main.go

.PHONY: add-admin
add-admin: build-scripts
	${SCRIPT_BINPATH} create

.PHONY: upsert-admin
upsert-admin: build-scripts
	${SCRIPT_BINPATH} create --force