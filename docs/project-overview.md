# Umag - Обзор проекта

## Назначение

Umag — платформа базовых микросервисов для SaaS-приложений. Обеспечивает управление пользователями, мультитенантность, биллинг, уведомления, медиафайлы и контроль доступа. На базе этих сервисов строятся продуктовые приложения (AIgenda/Qalai, BasQaru, Vibe, IdeasGen).

## Тип репозитория

**Monorepo** — 6 микросервисов + 1 shared-библиотека в одном рабочем пространстве.

## Технологический стек

| Категория | Технология | Версия |
|-----------|-----------|--------|
| Язык | Go | 1.21–1.23 |
| Фреймворк | go-kratos/kratos | v2.7.3 |
| ORM | entgo.io/ent | v0.14.x |
| База данных | PostgreSQL | 15.x |
| Messaging | NATS JetStream | — |
| Transport | gRPC + Protobuf | — |
| DI | Google Wire | — |
| Метрики | Prometheus | — |
| Трейсинг | OpenTelemetry (OTLP) | — |
| Service Discovery | HashiCorp Consul | — |
| Secrets | HashiCorp Vault | — |
| Кэш/Badges | Redis/Dragonfly | — |
| Контейнеризация | Docker + Helm (K8s) | — |
| CI/CD | GitLab CI | — |
| Миграции | Atlas (ent) | — |

## Сервисы

| Сервис | Модуль | Назначение |
|--------|--------|-----------|
| **billing** | `gitlab.calendaria.team/services/finance/billing` | Платежи (TipTopPay, Apple Store), подписки, продукты, инвойсы |
| **iam** | `gitlab.calendaria.team/services/iam` | Аутентификация (OTP, OAuth2), пользователи, credentials, settings, privacy |
| **media** | `gitlab.calendaria.team/services/media` | Загрузка/хранение файлов (S3), обработка изображений/видео/аудио |
| **notifications** | `gitlab.calendaria.team/services/notifications` | Push (FCM), Email (AWS SES), SMS (SMSC.kz), in-app уведомления |
| **rbac** | `gitlab.calendaria.team/services/rbac` | Роли, права, команды, проверка доступа |
| **tenants** | `gitlab.calendaria.team/services/tenants` | Организации, участники, группы, приглашения |
| **utils** | `gitlab.calendaria.team/services/utils` | Shared-библиотека: JWT, config, NATS, dialer, tracing, Redis badges |

## Архитектурный паттерн

Все сервисы следуют **Clean Architecture** (Kratos layout):

```
cmd/app/           → Bootstrap + Wire DI
internal/server/   → Transport (gRPC, HTTP, Cron)
internal/service/  → gRPC handlers (валидация, извлечение контекста)
internal/biz/      → Бизнес-логика (use cases)
internal/data/     → Репозитории (БД, внешние API)
ent/schema/        → Схема базы данных (Ent ORM)
api/{service}/v1/  → Protobuf-контракты
```

## Внешние зависимости

| Сервис | Провайдер | Назначение |
|--------|----------|-----------|
| Платежи | TipTopPay (KZ) | Прием карточных платежей |
| Платежи | Apple App Store | In-App Purchases (iOS) |
| Push | Firebase Cloud Messaging | Push-уведомления |
| Email | AWS SES | Отправка email |
| SMS | SMSC.kz | Отправка SMS |
| Storage | AWS S3 | Хранение медиафайлов |
| OAuth | Google, Sxodim | Внешняя аутентификация |

## Приложения на платформе

| App ID | Бренд | Описание |
|--------|-------|----------|
| `calendaria` | AIgenda | Календарь и планирование |
| `pms` | BasQaru | Project Management |
| `tickets` | Vibe | Мероприятия |
| `knowledge` | IdeasGen | База знаний |
