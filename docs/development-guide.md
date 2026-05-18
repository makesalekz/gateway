# Руководство по разработке

## Предварительные требования

- **Go** 1.23+ (для billing, iam, notifications, rbac) или Go 1.21+ (для media, tenants)
- **Docker** и **Docker Compose**
- **Make** (GNU Make)
- **protoc** (Protocol Buffers compiler)
- **ffmpeg** (только для media — обработка видео/аудио)
- Доступ к GitLab `gitlab.calendaria.team` (приватные Go-модули)

## Быстрый старт

### 1. Клонирование

```bash
git clone <repo-url> umag
cd umag
```

### 2. Установка инструментов (один раз)

Для любого сервиса:

```bash
cd <service>
make init
```

Устанавливает: protoc-gen-go, protoc-gen-go-grpc, protoc-gen-go-http, protoc-gen-go-errors, protoc-gen-openapi, kratos CLI, wire, mockgen, widdershins.

### 3. Локальный запуск сервиса

```bash
cd <service>
make start    # Собрать и запустить через docker-compose.local.yml
```

Или без Docker:

```bash
cd <service>
make run      # Запуск через kratos CLI (требует локальный PostgreSQL, NATS, Consul, Vault)
```

### 4. Остановка

```bash
make stop
```

## Основные Make-команды

| Команда | Описание |
|---------|----------|
| `make init` | Установить dev-инструменты |
| `make run` | Запуск локально (kratos) |
| `make start` | Сборка + запуск Docker |
| `make stop` | Остановка Docker |
| `make db` | Создать базу данных |
| `make all` | api + config + generate + tidy |
| `make api` | Генерация Go-кода из proto |
| `make config` | Генерация conf.pb.go |
| `make ent` | Генерация Ent ORM кода |
| `make migrations` | Генерация Atlas-миграций |
| `make generate` | Ent + Wire кодогенерация |
| `make build` | Сборка бинарника |
| `make lint` | Линтер (golangci-lint) |
| `make test` | Тесты |
| `make race` | Тесты с race detector (10x) |
| `make cover` | Покрытие тестами |
| `make mock` | Генерация моков (mockgen) |

## Переменные окружения

### Общие для всех сервисов

| Переменная | Описание | Пример |
|-----------|----------|--------|
| `SERVICE_NAME` | Имя сервиса для Consul | `iam` |
| `CONSUL_ADDRESS` | Адрес Consul | `http://consul:8500` |
| `VAULT_ADDRESS` | Адрес Vault | `http://vault:8200` |
| `VAULT_ROLE_ID` | Vault AppRole Role ID | — |
| `VAULT_SECRET_ID` | Vault AppRole Secret ID | — |
| `OTLP_GRPC_ADDRESS` | OpenTelemetry collector | `otlp:4317` |
| `DEBUG` | Режим отладки | `true` |
| `AUTOMIGRATE` | Авто-миграция БД при старте | `true` |
| `ENT_LOGGING` | SQL query logging | `true` |

### Специфичные для сервисов

| Сервис | Переменная | Описание |
|--------|-----------|----------|
| media | `AWS_REGION` | AWS регион для S3 |
| media | `AWS_BUCKET` | Имя S3-бакета |
| media | `AWS_ACCESS_KEY_ID` | AWS credentials |
| media | `AWS_SECRET_ACCESS_KEY` | AWS credentials |
| notifications | `GOOGLE_APPLICATION_CREDENTIALS` | Путь к Firebase JSON |
| notifications | `FIREBASE_CONFIG` | Firebase config path |
| notifications | `AWS_REGION` | AWS для SES |
| iam | `TOKEN_DURATION` | Время жизни access token |
| iam | `REFRESH_TOKEN_DURATION` | Время жизни refresh token |
| iam | `BRAND_NAME` | SMS sender name |

## Workflow разработки

### Изменение API (proto)

1. Редактировать `.proto` файлы в `api/{service}/v1/`
2. `make api` — сгенерировать Go-код
3. Реализовать новые методы в `internal/service/`

### Изменение схемы БД

1. Создать/изменить файл в `ent/schema/`
2. `make ent` — сгенерировать Ent-код
3. `make migrations` — создать SQL-миграцию
4. `make hash` — обновить хеши миграций

### Добавление зависимости (DI)

1. Создать конструктор `NewXxx(deps) *Xxx`
2. Добавить в `ProviderSet` соответствующего пакета
3. Обновить `wire.go` при необходимости
4. `make generate` — перегенерировать wire_gen.go

### Тестирование

```bash
make test              # Быстрые тесты
make race              # С race detector
make cover             # С покрытием
```

Моки генерируются через `make mock` (mockgen).

## Конфигурации

Каждый сервис имеет 3 конфига:
- `config.example.yaml` — для локальной разработки без Docker
- `config.dev.yaml` — для Docker-окружения
- `config.stage.yaml` — для staging (секреты из Vault)

В production DSN базы данных и API-ключи загружаются из Vault.

## CI/CD

GitLab CI pipeline (в `ci/` каждого сервиса):
1. **lint** — golangci-lint
2. **build** — сборка Docker-образа
3. **deploy (dev)** — Atlas-миграции + деплой в K8s
4. **release** — тегирование версии
5. **deploy (stage)** — деплой на staging

## Helm-деплой

Каждый сервис имеет Helm-чарт в `helm/`:
- Deployment + Service
- ConfigMap + Secrets
- HPA (Horizontal Pod Autoscaler)
- Ingress (если нужен HTTP)

## Conventions

- **Module naming**: `gitlab.calendaria.team/services/{name}`
- **Error codes**: определены в `errors.proto` каждого сервиса
- **Soft delete**: через `SoftDeleteMixin` (поле `deleted_at`)
- **Pagination**: cursor-based (fromId/toId) + page-based
- **Logging**: структурированный JSON (Zap через Kratos)
- **Metrics**: Prometheus на `/metrics` HTTP endpoint
