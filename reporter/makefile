MAKEFLAGS += --no-print-directory
ifneq (,$(wildcard .env))
	include .env
	export
endif

init:
	go mod tidy
	@echo "✅ tidy executed"
	go install github.com/air-verse/air@latest
	@echo "✅ air installed"
	make docker_build
	@echo "✅ docker built"
docker_up:
	docker compose -f ./docker-compose.yaml up -d
docker_down:
	docker compose -f ./docker-compose.yaml down
run:
	go run main.go
run_dev:
	air
docker_build:
	docker compose -f ./docker-compose.yaml build --no-cache
prepare_go_mod:
	@bash prepare_mod.sh
prepare_go_mod_test:
	@bash prepare_mod.sh test
