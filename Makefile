.PHONY: build
build: templ-gen
	go build -o ./bin/app ./cmd/app/main.go

.PHONY: run
run: build
	./bin/app

.PHONY: build-scripts
build-scripts:
	go build -o ./bin/scripts ./cmd/scripts/main.go

.PHONY: add-admin
add-admin: build-scripts
	./bin/scripts create

.PHONY: upsert-admin
upsert-admin: build-scripts
	./bin/scripts create --force

.PHONY: templ-gen
templ-gen:
	templ generate internal/view