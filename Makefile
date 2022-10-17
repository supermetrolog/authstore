.PHONY: run
run:
	go run cmd/main/main.go

.PHONY: build
build:
	(if exist "build" rd /q /s "build") && mkdir build && go build -o build/server.exe -v cmd/main/main.go

.PHONY: serve
serve: 
	./build/server.exe

.PHONY: migrate-up
migrate-up: 
	goose -dir ./db/migrations mysql "root:@/auth?parseTime=true" up

.PHONY: migrate-down
migrate-down: 
	goose -dir ./db/migrations mysql "root:@/auth?parseTime=true" down

.DEFAULT_GOAL := build

# create goose migration in directory example
# goose -dir ./db/migrations create add_name_column_in_tree_table sql