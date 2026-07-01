# reporter

сервис формирования выгрузок (отчётов).

## Локальный запуск

```
cp .env.example .env
make docker_up      # mysql + сервис (air)
make run            # или локально: go run main.go
```

Порты: gRPC `40021`, healthcheck `40020`.

## Как добавить новый тип выгрузки

Пример: код `transactions`.

1. Добавить код типа в `internal/defs/defs.go`.
2. Создать файл `internal/report/report_<code>.go` — реализовать интерфейс `ReportType`
   (`Code/Formatters/StorageCode/Init/Build`) и self-register через `init()` (`registerType`).
3. Добавить тип в `config.yaml` (`reports.<code>`): набор форматтеров и код стореджа.