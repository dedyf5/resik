---
name: resik-dependency-rules
description: Allowed and forbidden dependencies
---

# Dependency Rules

## Allowed flow

handler → core → repositories
handler → ctx
core → repositories
core → entities
repositories → entities
app → handler

## Rules

- No upward dependency
- No repository ↔ repository calls
- No handler → repository direct access
- No repository → core dependency

## Module independence rule

Every module must be independently movable to microservice.
