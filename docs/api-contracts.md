# API-контракты (gRPC)

## Обзор

Все сервисы предоставляют gRPC API (Protobuf v3). HTTP-серверы используются только для `/metrics` (Prometheus). Аутентификация — JWT middleware на каждом сервере.

---

## billing (finance-billing)

### Payments Service
| RPC | Описание |
|-----|----------|
| `CreatePayment(productId, amount, cryptogram, ip, name?, email?)` | Создание платежа через TipTopPay |
| `Complete3DS(transactionId, paRes)` | Завершение 3DS-верификации |
| `GetPaymentStatus(transactionId)` | Статус платежа |
| `CancelSubscription(userId, subscriptionId)` | Отмена подписки |
| `CheckWebhook(WebhookPayload)` | Webhook: проверка транзакции |
| `PaymentWebhook(WebhookPayload)` | Webhook: подтверждение оплаты |
| `RecurrentWebhook(WebhookPayload)` | Webhook: рекуррентный платеж |
| `PaymentCallback(data, sign)` | Legacy callback |

### Products Service
| RPC | Описание |
|-----|----------|
| `CreateProduct(ProductDto)` | Создание продукта |
| `UpdateProduct(id, ProductDto)` | Обновление продукта |
| `DeleteProduct(id)` | Удаление продукта |
| `GetProduct(id)` | Получение продукта |
| `ListProducts(pagination)` | Список продуктов |

### Invoices Service
| RPC | Описание |
|-----|----------|
| `GetInvoice(id)` | Получение инвойса |
| `ListInvoices(status?, productId?, paid?, subscriptionId?, pagination)` | Список инвойсов |
| `GetInvoiceReceipt(invoiceId)` | Данные чека |
| `GetInvoicePDF(invoiceId)` | PDF-чек |

### Items Service
| RPC | Описание |
|-----|----------|
| `CreateItem(ItemDto)` | Создание элемента |
| `UpdateItem(id, ItemDto)` | Обновление элемента |
| `DeleteItem(id)` | Удаление |
| `GetItem(id)` | Получение |
| `ListItems(pagination)` | Список |

### Subscriptions Service
| RPC | Описание |
|-----|----------|
| `GetSubscription(subscriptionId, withInvoices?)` | Получение подписки |
| `ListSubscriptions(withInvoices?, pagination)` | Список подписок |
| `GetSubscriptionStatus()` | Текущий статус подписки пользователя |

### AppleStore Service
| RPC | Описание |
|-----|----------|
| `ProcessServerNotification(signedPayload)` | Обработка Apple S2S notification |

---

## iam

### Auth Service
| RPC | Описание |
|-----|----------|
| `AuthByPhone(phone, isRegistrationNeeded?, isRegistration?, appSignature?)` | Аутентификация по телефону → OTP |
| `AuthByEmail(email, language?, isRegistrationNeeded?, isRegistration?)` | Аутентификация по email → OTP |
| `AuthByCode(userId, code)` | Верификация OTP → JWT tokens |
| `RefreshToken(tenantId)` | Обновление access token |

### Users Service
| RPC | Описание |
|-----|----------|
| `GetOwnProfile()` | Профиль текущего пользователя |
| `UpdateOwnProfile(phone?, email?, username?, name?, bio?, avatar?, timezone?)` | Обновление профиля |
| `DeleteOwnProfile()` | Запланировать удаление (30 дней) |
| `GetUserFull(userId)` | Полный профиль по ID |
| `GetUserByFilterFull(phone/email)` | Полный профиль по фильтру |
| `GetUser(userId)` | Краткий профиль |
| `GetUserByFilter(phone/email)` | Краткий профиль по фильтру |
| `GetUsers(ids[], phones[], emails[], withPrivacies?, withVerified?)` | Batch-получение |
| `ListUsers(ids[], search?, sort, paginate)` | Список пользователей |
| `UpdateUserLastSeen(userId, lastSeenTime)` | Обновление last_seen (s2s) |
| `BlockUser(userId)` | Блокировка |
| `UnblockUser(userId)` | Разблокировка |
| `DeleteUser(userId)` | Удаление пользователя |

### Credentials Service
| RPC | Описание |
|-----|----------|
| `ExternalAuth(authCode, provider)` | OAuth2 (Google/Sxodim) |
| `RefreshCredential(credentialId)` | Обновление OAuth-токена |
| `GetCredential(credentialId)` | Получение credential |
| `ListCredentials(provider?)` | Список credentials |
| `DeleteCredential(credentialId)` | Удаление credential |

### Settings Service
| RPC | Описание |
|-----|----------|
| `GetSettings()` | Настройки текущего пользователя |
| `UpdateSettings(map<string,string>)` | Обновление настроек |
| `GetUsersSettings(userIds[])` | Настройки нескольких пользователей |

### Privacy Service
| RPC | Описание |
|-----|----------|
| `GetPrivacy()` | Privacy-настройки текущего пользователя |
| `UpdatePrivacy(map<string,string>)` | Обновление privacy |
| `GetUsersPrivacies(ids[])` | Privacy нескольких пользователей |

---

## media

### MediaService
| RPC | Описание |
|-----|----------|
| `UploadMedia(content, fileName, filePath?, isPrivate?)` | Загрузка файла (до 2GB) |
| `GetMedia(mediaId)` | Получение медиа (public) |
| `GetPrivateMedia(mediaId)` | Получение медиа (private, pre-signed URL) |
| `GetMediaList(mediaIds[], ownOnly?)` | Список медиа по IDs |
| `DeleteAvatar(urls[])` | Удаление аватаров |

**Особенности:**
- HTTP-сервер принимает файлы через custom RequestDecoder (body = file, headers: X-File-Name, X-File-Path, X-File-Private)
- Max message size: 2GB (gRPC)
- Поддерживаемые типы: 36 MIME-типов (images, video, audio, documents, archives)
- Async processing: dimensions, thumbnails, duration — после ответа клиенту

---

## notifications

### Notifications Service
| RPC | Описание |
|-----|----------|
| `CreateNotifications(notifications[])` | Создание in-app уведомлений |
| `ListNotifications(language?, type?, paginate)` | Список уведомлений |
| `GetNotificationsCounters()` | Непрочитанные по типам |
| `DoActionOnNotification(notificationId, action, type)` | Действие (read) |

### Sender Service
| RPC | Описание |
|-----|----------|
| `CreateFcmDevice(language, token, oldToken?)` | Регистрация FCM-устройства |
| `DeleteFcmDevice(token)` | Удаление FCM-устройства |
| `PersonalSmsSender(phone, message, sender)` | Отправка SMS |
| `EmailSender(language?, type, emails[], data)` | Отправка email |

**Email types:** invite, confirm_email, new_user
**Языки:** en, ru, kk

---

## rbac

### Assigns Service
| RPC | Описание |
|-----|----------|
| `AssignRole(identityId, roleId, teamId?, resource?)` | Назначение роли |
| `AssignRoles(assigns[])` | Массовое назначение |
| `UnassignRole(assignId)` | Снятие роли |
| `ListAssigns(identityIds[], resourceTypes[], resources[])` | Список назначений |
| `GetAssign(assignId)` | Получение назначения |

### CheckPermissions Service
| RPC | Описание |
|-----|----------|
| `CheckPermissions(tenantId, permissions[], identities[], resources[], value?)` | Проверка прав → map<permission, {fields, resources}> |

### Permissions Service
| RPC | Описание |
|-----|----------|
| `CreatePermission(id, groupId, appId, name, description?, fields[])` | Создание permission |
| `UpdatePermission(permissionId, name?, description?, fields[])` | Обновление |
| `DeletePermission(permissionId)` | Удаление |
| `GetPermission(permissionId)` | Получение |
| `ListPermissions(appsIds[])` | Список по группам |

### Roles Service
| RPC | Описание |
|-----|----------|
| `CreateRole(name, description?, isSystem?, allow[], deny[])` | Создание роли |
| `UpdateRole(roleId, name?, description?, allow[], deny[])` | Обновление |
| `DeleteRole(roleId)` | Удаление (soft) |
| `GetRole(roleId)` | Получение |
| `ListRoles(search?, includeSystem?)` | Список ролей |
| `AddPermissionToRole(roleId, permissionId, deny?, fields[])` | Добавить permission к роли |
| `RemovePermissionFromRole(roleId, permissionId)` | Убрать permission |
| `ListRolePermissions(roleId)` | Permissions роли |

### Teams Service
| RPC | Описание |
|-----|----------|
| `CreateTeam(name, description?, parentId?)` | Создание команды |
| `UpdateTeam(teamId, name?, description?)` | Обновление |
| `DeleteTeam(teamId)` | Удаление (soft) |
| `GetTeam(teamId, withTree?)` | Получение (с поддеревом) |
| `GetTeams(teamIds[])` | Batch-получение |
| `ListTeams(parentId?, paginate)` | Список команд |

---

## tenants

### Tenants Service
| RPC | Описание |
|-----|----------|
| `CreateTenant(name, type?)` | Создание тенанта (default: PERSONAL) |
| `UpdateTenant(tenantId, name)` | Обновление |
| `DeleteTenant(tenantId)` | Удаление (soft) |
| `GetTenant(tenantId)` | Получение |
| `ListTenants(userId?, ownerId?, paginate)` | Список тенантов |
| `DeleteUsersTenants(usersIds[])` | Удаление PERSONAL тенантов пользователей |

### Members Service
| RPC | Описание |
|-----|----------|
| `CreateMembers(usersIds[])` | Добавление участников |
| `DeleteMember(memberId)` | Удаление участника |
| `GetMember(memberId)` | Получение участника |
| `GetShortMembers(identityIds[], withGroups?)` | По identity UUID |
| `GetMemberIdentities(userId, tenantId)` | Identities участника (member + groups) |
| `ListMembers(groupId?, search?, withGroups?, sort, paginate, excludeGroupId?)` | Список с обогащением |
| `CountMembers()` | Количество участников |

### Groups Service
| RPC | Описание |
|-----|----------|
| `CreateGroup(name, description?)` | Создание группы |
| `UpdateGroup(groupId, name?, description?)` | Обновление |
| `DeleteGroup(groupId)` | Удаление (soft) |
| `GetGroup(groupId)` | Получение |
| `ListGroups(search?, sort, paginate)` | Список групп |
| `AddMembersToGroup(groupId, membersIds[])` | Добавить участников |
| `RemoveMembersFromGroup(groupId, membersIds[])` | Убрать участников |

### Invites Service
| RPC | Описание |
|-----|----------|
| `CreateInvites(emails[], appId?, language?, roleId?, resource?, resourceId?)` | Создание приглашений |
| `CancelInvite(inviteId)` | Отмена приглашения |
| `DeleteInvite(inviteId)` | Удаление |
| `ListInvites(search?, status?, sort, paginate)` | Список приглашений |
| `AcceptInvite(inviteId, code)` | Принятие приглашения |
| `ShownInvite(inviteId, code)` | Просмотр приглашения |
| `DeclineInvite(inviteId, code)` | Отклонение |

---

## Общие Proto-модели (utils/v1)

| Message | Поля | Использование |
|---------|------|--------------|
| `EmptyRequest` | — | Запросы без параметров |
| `EmptyReply` | — | Ответы без данных |
| `PaginateRequest` | limit, fromId, toId, aroundId, fromDate, toDate, page, descending | Пагинация |
| `PaginateReply` | total?, fromId?, toId?, fromDate?, toDate? | Метаданные пагинации |
| `SortRequest` | field, descending | Сортировка |
| `ActorRequest` | actorId, tenantId | Контекст пользователя |
