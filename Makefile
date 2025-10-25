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
	cd ./webapp && \
	npm ci && \
	npm run dev

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