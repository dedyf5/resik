---
name: resik-anti-patterns
description: Forbidden patterns in Resik
---

# Anti Patterns

Never:

- Repository calling repository
- Repository handling transactions
- Handler calling repository directly
- Splitting CRUD into many files
- Cross-module SQL joins across unrelated domains
- Exposing entities directly in API contracts
- Global mock package
- Business logic in handler
