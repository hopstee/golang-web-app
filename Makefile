.PHONY: build
build:
	templ generate internal/view
	go build -o ./bin/app ./cmd/app/main.go

.PHONY: run
run: build
	./bin/app

.PHONY: build-scripts
build-scripts:
	go build -o ./bin/scripts ./cmd/scripts/main.go

.PHONY: add-admin
run-scripts: build-scripts
	./bin/scripts create

.PHONY: upsert-admin
run-scripts: build-scripts
	./bin/scripts create --force