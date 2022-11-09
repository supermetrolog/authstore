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

MOCKS_DESTINATION=mocks
.PHONY: mocks
# put the files with interfaces you'd like to mock in prerequisites
# wildcards are allowed
mocks: internal/common/loggerinterface/loggerinterface.go pkg/closer/closer.go
	@echo "Generating mocks..."
	@rm -rf $(MOCKS_DESTINATION)
	@for file in $^; do mockgen -source=$$file -destination=$(MOCKS_DESTINATION)/$$file; done