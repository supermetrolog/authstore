.PHONY: run
run:
	go run cmd/main/main.go

.PHONY: build
build:
	(if exist "build" rd /q /s "build") && mkdir build && go build -o build/server.exe -v cmd/main/main.go

.PHONY: serve
serve: 
	./build/server.exe

.DEFAULT_GOAL := build