RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color


start: ## Запустить инстанс
	docker compose up -d

stop: ## Остановить инстанс
	docker compose down

restart: ## Перезапустить инстанс
	docker compose restart

build: ## Только пересобрать контейнеры
	docker compose build --no-cache

rebuild: ## Пересобрать и пересоздать контейнеры (сохраняет данные)
	docker compose down
	docker compose build --no-cache
	docker compose up -d

clean: ## Удалить инстанс и ВСЕ данные
	@echo "$(RED)ВНИМАНИЕ! Это удалит все данные инстанса!$(NC)"
	@read -p "Продолжить? (y/n): " confirm; \
	if [ "$$confirm" = "y" ]; then \
		SPACE_NAME=$$(grep SPACE_NAME .env | cut -d '=' -f2); \
		docker compose down -v; \
		rm -rf ./nginx/logs-$$SPACE_NAME; \
		rm -rf ./redis/data-$$SPACE_NAME; \
		rm -rf ./mariadb/data-$$SPACE_NAME; \
		rm -rf ./mariadb/logs-$$SPACE_NAME; \
		rm -rf ./elasticsearch/data-$$SPACE_NAME; \
		echo "$(GREEN)✓ Инстанс удален$(NC)"; \
	else \
		echo "$(YELLOW)Отменено$(NC)"; \
	fi

php-exec: ## Войти в php контейнер
	docker compose exec php bash

db-exec: ## Войти в db контейнер
	docker compose exec mariadb bash