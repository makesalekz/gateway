# Анализ структуры исходного кода

## Корневая структура

```
umag/
├── billing/              # Сервис биллинга и платежей
├── iam/                  # Identity & Access Management
├── media/                # Хранение и обработка медиафайлов
├── notifications/        # Push/Email/SMS уведомления
├── rbac/                 # Role-Based Access Control
├── tenants/              # Мультитенантность и организации
├── utils/                # Shared-библиотека утилит
├── docs/                 # Документация проекта (этот каталог)
├── _bmad/                # BMad workflow конфигурация
└── _bmad-output/         # Артефакты планирования
```

## Общая структура каждого микросервиса

Все 6 сервисов следуют единому layout (Kratos Clean Architecture):

```
{service}/
├── api/{service}/v1/     # Protobuf-контракты
│   ├── models.proto      # Общие модели данных
│   ├── errors.proto      # Коды ошибок
│   ├── *.proto           # RPC-сервисы
│   └── *.pb.go           # Сгенерированный Go-код
├── cmd/app/              # Точка входа
│   ├── main.go           # Bootstrap приложения
│   ├── wire.go           # Wire DI-определения
│   └── wire_gen.go       # Сгенерированный DI-код
├── configs/              # Конфигурации по окружениям
│   ├── config.example.yaml
│   ├── config.dev.yaml
│   └── config.stage.yaml
├── ent/                  # Ent ORM
│   ├── schema/           # Определения сущностей
│   └── ...               # Сгенерированный код
├── internal/             # Внутренняя логика (не экспортируется)
│   ├── biz/              # Бизнес-логика (use cases)
│   ├── conf/             # Proto-конфигурация (conf.proto)
│   ├── data/             # Репозитории и внешние клиенты
│   ├── server/           # Серверы (gRPC, HTTP, Cron)
│   └── service/          # gRPC-хендлеры
├── third_party/          # Vendor proto-файлы (google, validate)
├── ci/                   # GitLab CI pipelines
├── helm/                 # Helm-чарты для K8s
├── doc/                  # Документация API (сгенерированная)
├── Dockerfile            # Multi-stage build
├── docker-compose.yml    # Production deploy
├── docker-compose.local.yml  # Локальная разработка
├── Makefile              # Команды сборки/тестирования
├── go.mod / go.sum       # Go модули
├── openapi.yaml          # OpenAPI spec (stub)
└── README.md             # README сервиса
```

## Структура utils (shared-библиотека)

```
utils/
├── api/utils/v1/         # Общие Proto-модели (EmptyRequest, PaginateRequest, etc.)
├── v1/                   # Первое поколение (устаревший dialer)
│   ├── config/           # Config (Consul + Vault)
│   ├── jwt/              # JWT-процессор
│   ├── dialer/           # gRPC-диалер (deprecated)
│   ├── nats/             # NATS (простой pub/sub)
│   ├── log/              # Zap logger
│   ├── error/            # PostgreSQL error checks
│   ├── pagination/       # Cursor-пагинация
│   └── middlewares/      # Auth, Metrics, Error middlewares
├── v2/                   # Второе поколение (JetStream, s2s, tracing)
│   ├── auth/             # Context helpers (actorId, tenantId)
│   ├── jwt/              # IJwtProcessor interface
│   ├── dialer/           # Connection manager с состояниями
│   ├── nats/             # JetStream queues
│   ├── tracing/          # OpenTelemetry tracer
│   ├── zap/              # Logger
│   ├── uuid/             # UUID с embedded actorID
│   ├── struc/            # Shared structs/enums
│   └── middlewares/      # s2s JWT + BffMeta
├── v3/                   # Третье поколение (per-service JWT)
│   ├── jwt/              # IJwtSecret
│   ├── dialer/           # Per-endpoint JWT from Vault
│   └── middlewares/      # Separate BFF/S2S servers
├── v4/                   # Четвертое поколение (текущее, interface-first)
│   ├── config/           # IConfig interface
│   ├── jwt/              # IJwtProcessor + IConfig
│   ├── dialer/           # Interface params, 20MB msg size
│   ├── nats/             # Options pattern + PubDelayed
│   ├── tracing/          # ITracer interface
│   ├── badge/            # IBadgeClient (Redis)
│   └── middlewares/      # Interface-based auth
├── go.mod
└── Makefile
```

## Критические директории по сервисам

### billing
- `internal/biz/payments.go` — логика платежей TipTopPay
- `internal/biz/apple-store.go` — Apple IAP обработка
- `internal/biz/invoice-manager.go` — создание инвойсов с валидацией
- `internal/data/ttp_client.go` — HTTP-клиент TipTopPay
- `ent/schema/` — 7 сущностей (Bundle, Invoice, Item, PaymentProfile, Product, ProductReservation, Subscriptions)
- `messages/` — структуры NATS-сообщений

### iam
- `internal/biz/auth.go` — OTP-аутентификация, генерация JWT
- `internal/biz/users.go` — CRUD пользователей, scheduled deletion
- `internal/biz/credentials.go` — OAuth2 (Google, Sxodim)
- `internal/data/integration/` — Google/Sxodim OAuth gateways
- `internal/data/remote/` — вызовы к tenants, notifications, contacts, chats, events, media, websockets
- `internal/server/cron.go` — cron для удаления аккаунтов
- `ent/schema/` — 5 сущностей (User, UserSettings, UserCredentials, UserPrivacy, OneTimePassword)

### media
- `internal/biz/media.go` — upload, async processing, private media
- `internal/data/s3.go` — AWS S3 клиент (upload, delete, pre-signed URLs)
- `internal/data/video_processor.go` — ffmpeg для видео
- `internal/data/audio_processor.go` — ffmpeg для аудио
- `ent/schema/` — 1 сущность (Media) с SoftDeleteMixin

### notifications
- `internal/biz/fcm.go` — FCM push (основная логика)
- `internal/biz/email.go` — Email через AWS SES
- `internal/biz/sms.go` — SMS через SMSC.kz
- `internal/data/localizer.go` — i18n (ru/en/kk)
- `locales/` — JSON-файлы локализации
- `templates/` — HTML email-шаблоны
- `ent/schema/` — 4 сущности (Device, Notification, NotificationData, LastReadNotification)

### rbac
- `internal/biz/check_permissions.go` — логика проверки прав
- `internal/biz/assigned_roles.go` — назначение ролей
- `internal/biz/paid_content.go` — обработка подписок из billing
- `ent/schema/` — 8 сущностей (Permission, PermissionGroup, Role, RolePermission, Team, ResourceAccess, ResourceType, TeamIdentityRole)

### tenants
- `internal/biz/invites.go` — приглашения с email-нотификацией
- `internal/biz/members.go` — управление участниками
- `internal/biz/groups.go` — группы внутри тенанта
- `internal/data/remote/` — вызовы к iam, rbac
- `ent/schema/` — 4 сущности (Tenant, Member, Group, Invite)

## Порты сервисов

| Сервис | HTTP | gRPC | DB Name |
|--------|------|------|---------|
| iam | 8000 | 9000 | iam |
| media | 8004 | 9004 | media |
| notifications | 8003 | 9003 | notifications |
| tenants | 8007 | 9007 | tenants |
| rbac | 8009 | 9009 | rbac |
| billing | 8020 | 9020 | billing |
