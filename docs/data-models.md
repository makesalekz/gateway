# Модели данных

## Обзор

Все сервисы используют **Ent ORM** с PostgreSQL. Каждый сервис имеет свою изолированную базу данных. Миграции управляются через **Atlas**.

---

## billing (БД: billing)

### Bundle
| Поле | Тип | Ограничения |
|------|-----|------------|
| id | int64 | PK, immutable |
| product_id | int64 | FK→Product, immutable |
| item_id | int64 | FK→Item, immutable |
| amount | float | default: 1 |
| created_at | time | immutable, default: now |
| updated_at | time | auto-update |
| deleted_at | time? | soft delete |

### Invoice
| Поле | Тип | Ограничения |
|------|-----|------------|
| id | int64 | PK |
| user_id | int64 | immutable |
| tenant_id | int64 | immutable |
| app_id | string | immutable |
| product_id | int64 | FK→Product, immutable |
| amount | int64 | immutable |
| price | decimal | numeric |
| currency | string | default: "KZT", max: 3 |
| status | enum | CREATED/PAID/CANCELED_BY_USER/CANCELED_BY_VENDOR/FAILED/REJECTED/REVOKED |
| paid_at | time? | |
| paid_till | time? | |
| is_revoked | bool | default: false |
| revoked_at | time? | |
| is_revoked_processed | bool | default: false |
| is_paid_at_processed | bool | default: false |
| is_paid_till_processed | bool | default: false |
| subscription_id | int64? | FK→Subscriptions |
| external_transaction_id | string? | |
| payment_provider | enum | APP_STORE/ONE_VISION_PAYMENT/TIP_TOP_PAYMENT |
| is_trial | bool | default: false |
| payment_profile_id | int64? | FK→PaymentProfile |
| ttp_subscription_id | string? | |

### Item
| Поле | Тип | Ограничения |
|------|-----|------------|
| id | int64 | PK |
| name | string | |
| description | string | |
| topic_name | string? | unique |
| created_at/updated_at/deleted_at | time | mixins |

### PaymentProfile
| Поле | Тип | Ограничения |
|------|-----|------------|
| id | int64 | PK |
| user_id | int64 | immutable |
| pan_masked | string | immutable |
| holder | string | immutable |
| email | string | default: "" |
| phone | string | default: "" |
| user_token | string | default: "" |
| recurrent_token | string? | unique |
| created_at/updated_at/deleted_at | time | mixins |

### Product
| Поле | Тип | Ограничения |
|------|-----|------------|
| id | int64 | PK |
| app_id | string | |
| name | string | |
| description | string | |
| price | decimal | numeric |
| currency | string | default: "KZT" |
| is_active | bool | default: true |
| is_limited | bool | default: false |
| limited_till | time? | |
| left | int64 | default: 0 (stock) |
| is_unique | bool | default: false |
| unique_limit | int64 | max: 100 |
| is_expiring | bool | default: false |
| expiring_time | time? | |
| payment_model | enum | ONE_TIME/RECURRENT |
| product_period | enum | day/week/month/year/unlimited |
| created_at/updated_at/deleted_at | time | mixins |

### ProductReservation
| Поле | Тип | Ограничения |
|------|-----|------------|
| id | int64 | PK |
| product_id | int64 | FK→Product |
| invoice_id | int64 | FK→Invoice |
| user_id | int64 | immutable |
| reserved_quantity | int64 | positive, default: 1 |
| status | enum | PENDING/COMPLETED/EXPIRED/CANCELLED |
| expiration_time | time | default: now + 15min |
| created_at/updated_at | time | mixins |

### Subscriptions
| Поле | Тип | Ограничения |
|------|-----|------------|
| id | int64 | PK |
| user_id | int64 | immutable |
| tenant_id | int64 | immutable |
| app_id | string | immutable |
| product_id | int64 | FK→Product |

---

## iam (БД: iam)

### User
| Поле | Тип | Ограничения |
|------|-----|------------|
| id | int64 | PK, immutable |
| phone | string? | unique |
| email | string? | unique |
| username | string? | unique, 3-30 chars |
| name | string | default: "" |
| bio | string | default: "" |
| avatar | string? | |
| timezone | string | default: "UTC" |
| is_active | bool | default: false |
| phone_verified | bool | default: false |
| email_verified | bool | default: false |
| last_login_at | time | default: now |
| last_seen | time? | |
| default_tenant_id | int64? | |
| is_blocked | bool | default: false |
| created_at/updated_at | time | |
| deleted_at | time? | soft delete |
| remove_at | time? | scheduled deletion (30 дней) |

### UserSettings
| Поле | Тип | Ограничения |
|------|-----|------------|
| id | int | PK |
| user_id | int64 | FK→User |
| setting | enum | LANGUAGE/THEME/NOTIFICATION_SOUND_ENABLED/NOTIFICATION_VIBRATION_ENABLED/EVENTS_CHAT_ENABLED |
| value | string | |
| updated_at | time | |
| **Index** | unique(user_id, setting) | |

### UserCredentials
| Поле | Тип | Ограничения |
|------|-----|------------|
| id | int | PK |
| user_id | int64 | FK→User |
| external_user_id | int64? | |
| mail | string? | |
| phone | string? | |
| display_name | string? | |
| provider | enum? | Calendaria/Google/Outlook/Apple/Sxodim |
| access_token | string | |
| token_type | string? | |
| refresh_token | string? | |
| expires_at | time? | |
| created_at/updated_at/deleted_at | time | mixins |

### UserPrivacy
| Поле | Тип | Ограничения |
|------|-----|------------|
| id | int | PK |
| user_id | int64 | FK→User |
| setting | enum | MY_LAST_ACTIONS/MY_PROFILE_IMAGE/MY_EVENTS/GROUP_CHAT_INVITE/EVENT_INVITE/MY_SLOTS/SLOTS_DETAILS/LAST_VISIT |
| option | enum | ALL/MY_CONTACTS/NO_ONE |
| updated_at | time | |
| **Index** | unique(user_id, setting) | |

### OneTimePassword
| Поле | Тип | Ограничения |
|------|-----|------------|
| id | int | PK |
| user_id | int64 | FK→User |
| code | string | max: 6 |
| type | enum | EMAIL/PHONE |
| is_used | bool | default: false |
| expires_at | time | (5 мин) |
| created_at | time | |
| failed_attempts | int64 | default: 0, max: 5 |

---

## media (БД: media)

### Media
| Поле | Тип | Ограничения |
|------|-----|------------|
| id | int | PK |
| owner_id | int64 | immutable |
| file_name | string | immutable |
| extension | string | immutable, 2-10 chars |
| path | string | immutable, unique |
| url | string? | |
| size | int32 | bytes |
| width | int32? | pixels |
| height | int32? | pixels |
| duration | float32? | seconds |
| thumbnail_url | string? | |
| thumbnail_path | string? | |
| is_activated | bool | default: false |
| is_private | bool | default: false |
| created_at | time | immutable |
| uploaded_at | time? | |
| deleted_at | time? | soft delete |
| **Index** | url | |

---

## notifications (БД: notifications)

### Device
| Поле | Тип | Ограничения |
|------|-----|------------|
| id | int | PK |
| user_id | int64 | positive |
| token | string | unique, immutable, min: 1 |
| language | string | default: "en" |
| created_at | time | immutable |

### Notification
| Поле | Тип | Ограничения |
|------|-----|------------|
| id | int | PK |
| user_id | int64 | positive |
| type | enum | common/calendar/event/contact/tasks/projects/chat |
| title | string | not empty |
| text | string | not empty |
| event_id | int64? | |
| contact_id | int64? | |
| task_id | int64? | |
| project_id | int64? | |
| notification_data_id | int64? | FK→NotificationData |
| created_at | time | |
| **Index** | (user_id, type) | |

### NotificationData
| Поле | Тип | Ограничения |
|------|-----|------------|
| id | int64 | PK |
| type | string? | |
| event | string? | JSON |
| member | string? | JSON |
| chat | string? | JSON |
| message | string? | JSON |
| contact | string? | JSON |
| task | string? | JSON |
| project | string? | JSON |
| target_user_id | int64? | |
| metadata | string? | JSON |
| plural_count | int64? | |

### LastReadNotification
| Поле | Тип | Ограничения |
|------|-----|------------|
| id | int | PK |
| user_id | int64 | immutable |
| type | enum | (notification types) |
| last_read_id | int64 | |
| **Index** | unique(user_id, type) | |

---

## rbac (БД: rbac)

### Permission
| Поле | Тип | Ограничения |
|------|-----|------------|
| id | string | PK (format: group.action), immutable |
| group_id | string | FK→PermissionGroup, immutable |
| name | string | max: 32, not empty |
| description | string | default: "" |
| app_id | string | max: 10, immutable |
| fields | JSON []string | default: [] |

### PermissionGroup
| Поле | Тип | Ограничения |
|------|-----|------------|
| id | string | PK, immutable |
| app_id | string | max: 10, immutable |
| name | string | max: 16, not empty |

### Role
| Поле | Тип | Ограничения |
|------|-----|------------|
| id | int64 | PK, immutable |
| name | string | max: 32, not empty |
| description | string | default: "" |
| tenant_id | int64 | immutable, default: 0 (system) |
| is_system | bool | immutable, default: false |
| created_at/updated_at | time | |
| deleted_at | time? | soft delete |

### RolePermission
| Поле | Тип | Ограничения |
|------|-----|------------|
| id | int | PK |
| tenant_id | int64 | immutable, default: 0 |
| role_id | int64 | FK→Role, immutable |
| permission_id | string | FK→Permission, immutable |
| deny | bool | default: false |
| fields | JSON []string | |
| value | int64 | default: 0 |
| **Index** | unique(role_id, permission_id) | |

### Team
| Поле | Тип | Ограничения |
|------|-----|------------|
| id | int64 | PK |
| tenant_id | int64 | |
| parent_id | int64? | FK→Team (self-ref) |
| parents_ids | bigint[] | PostgreSQL array |
| name | string | |
| description | string | default: "" |
| created_at/updated_at | time | |
| deleted_at | time? | soft delete |

### ResourceAccess
| Поле | Тип | Ограничения |
|------|-----|------------|
| id | int | PK |
| tenant_id | int64 | immutable |
| resource_type | string? | FK→ResourceType, immutable |
| resource_id | int64? | immutable |
| identity_id | string | immutable, default: "" |
| role_id | int64 | FK→Role, immutable |
| **Index** | unique(tenant_id, role_id, identity_id, resource_type, resource_id) WHERE resource_id IS NOT NULL |
| **Index** | unique(tenant_id, role_id, identity_id) WHERE resource_id IS NULL |

### ResourceType
| Поле | Тип | Ограничения |
|------|-----|------------|
| id | string | PK, immutable |
| description | string | default: "" |

---

## tenants (БД: tenants)

### Tenant
| Поле | Тип | Ограничения |
|------|-----|------------|
| id | int64 | PK, immutable |
| owner_id | int64 | |
| name | string | |
| type | enum | PERSONAL/BUSINESS, immutable, default: "PERSONAL" |
| created_at/updated_at | time | |
| deleted_at | time? | soft delete |
| **Index** | (owner_id, type) | |

### Member
| Поле | Тип | Ограничения |
|------|-----|------------|
| id | int | PK |
| identity_id | UUID | immutable, unique |
| tenant_id | int64 | FK→Tenant, immutable |
| user_id | int64 | immutable |
| created_at | time | immutable |
| deleted_at | time? | soft delete |
| **Index** | unique(tenant_id, user_id) | |
| **Edge** | groups (M2M with Group) | |

### Group
| Поле | Тип | Ограничения |
|------|-----|------------|
| id | int | PK |
| identity_id | UUID | immutable, unique |
| tenant_id | int64 | FK→Tenant, immutable |
| name | string | not empty |
| description | string | |
| created_at/updated_at | time | mixins |
| deleted_at | time? | soft delete |
| **Index** | unique(tenant_id, name) | |
| **Edge** | members (M2M with Member) | |

### Invite
| Поле | Тип | Ограничения |
|------|-----|------------|
| id | int | PK |
| tenant_id | int64 | FK→Tenant, immutable |
| code | UUID | immutable |
| email | string | immutable |
| user_id | int64? | |
| status | enum | NEW/SENT/SHOWN/ACCEPTED/DECLINED/CANCELED |
| role_id | int64 | optional |
| resource | string | optional |
| resource_id | int64 | optional |
| created_at/updated_at | time | |
| **Index** | unique(tenant_id, user_id) | |
| **Index** | unique(tenant_id, email) | |

---

## Общие паттерны

### SoftDeleteMixin
Все сущности с пометкой "soft delete" используют общий mixin:
- Добавляет поле `deleted_at` (nullable time)
- Interceptor: WHERE deleted_at IS NULL (автоматически)
- Hook: DELETE → UPDATE SET deleted_at = now()
- Обход: `SkipSoftDelete(ctx)` в контексте

### CreateUpdateMixin
- `created_at` — immutable, default: time.Now
- `updated_at` — auto-updates on mutation
