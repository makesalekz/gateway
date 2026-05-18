---
name: "Endpoint Request"
about: "Request a new API endpoint or modify an existing one"
title: "[ENDPOINT] "
labels: ["endpoint", "mobile"]
---

## Endpoint

**Method:** GET / POST / PUT / DELETE
**Path:** `/api/v1/...`

## Description

Что должен делать endpoint.

## Request

```json
{
  "field": "type — описание"
}
```

## Response

```json
{
  "field": "type — описание"
}
```

## Headers

- `Authorization: Bearer <token>` (required)
- Other headers...

## Backend Service

Какой gRPC сервис вызывать: products / warehouse / sales / orders / agents / stores / platform-billing / iam

## Notes

Доп. требования, edge cases, pagination и т.д.
