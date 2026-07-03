---
name: resik-architecture
description: System architecture and layer responsibilities
---

# Architecture

## Flow

REST / gRPC
    ↓
handler
    ↓
core
    ↓
repositories
    ↓
database

## Cross-cutting layers

- ctx → request lifecycle
- entities → shared models
- internal → framework internals
- pkg → reusable libraries
- utils → application helpers
