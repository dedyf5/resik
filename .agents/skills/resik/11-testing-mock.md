---
name: resik-testing-mock
description: Testing and mock strategy
---

# Testing & Mock

## Rules

- Each module owns its own mock
- No global mock package
- Mock must live inside domain

## Reason

Supports modular monolith → microservice extraction without refactoring mocks.
