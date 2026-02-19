# Docker Multi-Instance Environment

Запуск нескольких изолированных инстансов Laravel приложения.

## Быстрый старт

**1. Скопировать конфиг:**
```bash
cp .env.example .env
```

**2. Изменить в .env:**
- `SPACE_NAME` - имя инстанса для контейнеров (например: `madwave_prod`)
- `APP_DIR_NAME` - название папки приложения (например: `madwave_ru`)
- `APP_DOMAIN` - домен (например: `madwave.ru`)
- Все порты (каждый инстанс должен использовать свои)
- Настройки БД (MYSQL_DATABASE, MYSQL_USER, MYSQL_PASSWORD)

**3. Запустить:**
```bash
make start
```

## Основные команды

```bash
make start        # Запустить инстанс
make stop         # Остановить инстанс
make restart      # Перезапустить инстанс
make build        # Пересобрать контейнеры без кеша
make clean        # Удалить инстанс со всеми данными 
```

## Ключевые параметры .env

```env
# Идентификация
SPACE_NAME=madwave_prod        # Для названий контейнеров
APP_DIR_NAME=madwave_ru        # Папка с приложением
APP_DOMAIN=madwave.ru          # Домен

# Порты (меняйте для каждого инстанса)
NGINX_HTTP_PORT=80
NGINX_HTTPS_PORT=443
REDIS_PORT=6379
MARIADB_PORT=3306
PHP_PORT=9003
WEBSOCKET_PORT=6001
PMA_PORT=8081
ELASTICSEARCH_HOST_HTTP_PORT=9200

# База данных
MYSQL_DATABASE=madwave
MYSQL_USER=madwave
MYSQL_PASSWORD=secret
```

## Примеры конфигураций

**Production:**
```env
SPACE_NAME=madwave_prod
APP_DIR_NAME=madwave_ru
NGINX_HTTP_PORT=80
MARIADB_PORT=3306
```

**Staging на том же сервере:**
```env
SPACE_NAME=madwave_staging
APP_DIR_NAME=madwave_ru
NGINX_HTTP_PORT=81
MARIADB_PORT=33061
```

