# Umag — Документация проекта

## Обзор

| Параметр | Значение |
|----------|---------|
| **Тип** | Monorepo — 6 микросервисов + 1 shared-библиотека |
| **Язык** | Go 1.21–1.23 |
| **Фреймворк** | go-kratos v2.7.3 |
| **Архитектура** | Clean Architecture (Kratos layout) |
| **ORM** | Ent (PostgreSQL) |
| **Transport** | gRPC + Protobuf |
| **Messaging** | NATS JetStream |

## Сервисы

| Сервис | Порт (gRPC) | Назначение |
|--------|------------|-----------|
| [billing](#billing) | :9020 | Платежи, подписки, продукты, инвойсы |
| [iam](#iam) | :9000 | Аутентификация, пользователи, OAuth2 |
| [media](#media) | :9004 | Загрузка/хранение файлов (S3) |
| [notifications](#notifications) | :9003 | Push (FCM), Email (SES), SMS |
| [rbac](#rbac) | :9009 | Роли, права, команды |
| [tenants](#tenants) | :9007 | Организации, участники, приглашения |
| [utils](#utils) | — | Shared-библиотека |

---

## Документация

### Архитектура и обзор
- [Обзор проекта](./project-overview.md) — стек, сервисы, назначение
- [Архитектура интеграций](./integration-architecture.md) — взаимодействие между сервисами, NATS-очереди, авторизация
- [Структура исходного кода](./source-tree-analysis.md) — дерево каталогов, порты, критические файлы

### API и данные
- [API-контракты (gRPC)](./api-contracts.md) — все RPC-методы всех сервисов
- [Модели данных](./data-models.md) — Ent-схемы, таблицы, поля, связи

### Разработка
- [Руководство по разработке](./development-guide.md) — быстрый старт, команды, env vars, CI/CD

---

## Быстрый старт

```bash
cd <service>        # billing / iam / media / notifications / rbac / tenants
make init           # Установить инструменты (один раз)
make start          # Собрать и запустить в Docker
```

---

## Ключевые технологии

| Компонент | Технология |
|-----------|-----------|
| Service Discovery | Consul |
| Secrets | Vault (AppRole) |
| Database | PostgreSQL + Ent ORM + Atlas migrations |
| Messaging | NATS JetStream |
| Tracing | OpenTelemetry (OTLP) |
| Metrics | Prometheus |
| CI/CD | GitLab CI + Helm + K8s |
| Payments | TipTopPay (KZ), Apple App Store |
| Push | Firebase Cloud Messaging |
| Email | AWS SES |
| SMS | SMSC.kz |
| Storage | AWS S3 |
| Cache | Redis/Dragonfly |

---

## Для AI-assisted разработки

При создании PRD или планировании фич ссылайтесь на:
- **Full-stack фичи**: этот index + `integration-architecture.md`
- **API-фичи**: `api-contracts.md` + `data-models.md`
- **Новый сервис**: `project-overview.md` + `development-guide.md`
- **Изменение БД**: `data-models.md`

---

*Сгенерировано: 2026-05-16 | Scan level: exhaustive | Workflow: initial_scan*
