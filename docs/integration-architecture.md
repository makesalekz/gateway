# Архитектура интеграций между сервисами

## Обзор

Сервисы взаимодействуют через **gRPC** (синхронно) и **NATS JetStream** (асинхронно). Все gRPC-вызовы проходят через Consul service discovery. Аутентификация между сервисами — s2s JWT-токены (HS256, 60 мин TTL).

## Граф зависимостей

```
                    ┌──────────┐
                    │  billing │
                    └────┬─────┘
                         │ NATS: role events
                         ▼
┌────────┐  gRPC   ┌──────────┐  gRPC   ┌──────────┐
│  iam   │◄────────│ tenants  │────────►│   rbac   │
└───┬────┘         └──────────┘         └──────────┘
    │                    │
    │ gRPC               │ NATS: email
    ▼                    ▼
┌────────────────┐  ┌──────────────────┐
│    media       │  │  notifications   │
└────────────────┘  └──────────────────┘
```

## Матрица интеграций

### gRPC-вызовы (синхронные)

| Вызывающий | Вызываемый | Методы | Назначение |
|-----------|-----------|--------|-----------|
| tenants | iam | GetUser, GetUsers, ListUsers | Обогащение данных участников |
| tenants | rbac | AssignRoles | Назначение ролей при создании тенанта/инвайта |
| notifications | iam | GetUsers, GetUsersSettings | Данные пользователей, настройки уведомлений |
| notifications | chats* | CountUnreadMessages | Подсчет непрочитанных для badges |
| notifications | events* | GetEventsCount | Подсчет ожидающих событий для badges |
| notifications | contacts* | GetContactsByUserID, GetBatchContactLabels | Разрешение имен контактов |
| iam | tenants | CreateTenants, GetUserTenants, GetMemberIdentities, DeleteUsersTenants | Управление тенантами пользователя |
| iam | notifications | PersonalSmsSender, PersonalEmailSender | OTP через SMS/Email |
| iam | contacts* | GetIncomingRelations, DeleteUsersDataInContacts | Данные контактов |
| iam | chats* | DeleteUsersDataInChats | Удаление данных при удалении аккаунта |
| iam | events* | DeleteUsersDataInEvents, DisconnectExternalCalendarsBulk | Удаление данных/отключение календарей |
| iam | media | DeleteAvatar | Удаление аватара при обновлении |
| iam | websockets* | GetUserStatus, ListUsersStatuses | Онлайн-статус пользователя |
| billing | — | (нет gRPC-вызовов к другим сервисам в этом репо) | — |

> *Сервисы chats, events, contacts, websockets — внешние, не входят в данный monorepo.

### NATS JetStream (асинхронные)

| Публикатор | Очередь/Топик | Подписчик | Payload | Назначение |
|-----------|--------------|-----------|---------|-----------|
| tenants | `notifications.email` | notifications | EmailDetails | Email-приглашения в тенант |
| iam | `contacts.confirmed_phone` | contacts* | {userId, phone} | Телефон подтвержден |
| iam | `contacts.confirmed_emails` | contacts* | {userId, email} | Email подтвержден |
| iam | `events.default_calendars` | events* | {userId, tenantId} | Создание дефолтных календарей |
| iam | `notifications.delete_tokens` | notifications | {userId} | Удаление FCM-токенов при удалении аккаунта |
| notifications | `fcm` (local) | notifications | FirebaseNotification | Push-уведомления |
| notifications | `fcm_silent` (local) | notifications | FirebaseNotification | Silent push (badge update) |
| notifications | `email` (local) | notifications | EmailDetails | Отправка email |
| notifications | `delete_tokens` (local) | notifications | {userId} | Удаление токенов устройств |
| media | `delete` (local) | media | mediaID | Удаление медиа из S3 |
| media | `delete_list` (local) | media | []mediaID | Массовое удаление медиа |
| media | `delete_record` (local) | media | []mediaID | Повторное удаление записей БД |
| rbac | `role_assign` (local) | rbac | AssignRoleDto | Логирование назначения роли |
| rbac | `role_unassign` (local) | rbac | assignID | Логирование снятия роли |
| billing | `{item.topic_name}` | external* | RefreshItems | Активация ресурсов при оплате |
| billing | `{item.topic_name}_revoke` | external* | RefreshItems | Отзыв ресурсов при истечении |
| billing (remote) | `finance-billing.*` | rbac | subscription event | Premium/Expired → управление ролями |

## Аутентификация и авторизация

### Поток аутентификации

1. Клиент вызывает `iam.AuthByPhone/AuthByEmail` → получает userID
2. Клиент отправляет OTP-код через `iam.AuthByCode` → получает accessToken + refreshToken
3. accessToken (JWT HS256) содержит: userId, tenantId, memberId, groupsIds, issuer="iam"
4. refreshToken — долгоживущий JWT (30 дней по умолчанию)

### Авторизация запросов

1. Каждый gRPC-сервер имеет JWT middleware (`u_auth.Server`)
2. Из JWT извлекаются: actorId, tenantId, identities
3. Для проверки прав: сервис вызывает `rbac.CheckPermissions(tenantId, identities, permissions)`
4. RBAC возвращает map разрешенных permissions с полями и ресурсами

### Service-to-Service (s2s)

1. Исходящий middleware (`u_auth.Client`) генерирует s2s JWT (issuer=serviceName, audience=targetService)
2. Принимающий middleware (`s2s.Server`) валидирует s2s JWT
3. Секреты JWT хранятся в Vault (global или per-service)

## Протокол передачи контекста

Metadata gRPC headers:
- `x-md-global-actor-id` — ID пользователя
- `x-md-global-tenant-id` — ID тенанта
- `x-md-global-app-id` — ID приложения
- `x-md-global-identities` — UUID-идентификаторы (member + groups)

## Инфраструктурная интеграция

| Компонент | Адрес (dev) | Назначение |
|-----------|-------------|-----------|
| Consul | consul:8500 | Service discovery, KV-конфигурация |
| Vault | vault:8200 | Секреты (DB DSN, JWT, API keys) |
| NATS | nats:4222 | Messaging (JetStream) |
| PostgreSQL | postgres_db:5432 | Каждый сервис — своя БД |
| Redis/Dragonfly | dragonfly:6379 | Badges (notifications) |
| OTLP Collector | otlp:4317 | Distributed tracing |
