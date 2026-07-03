---
name: resik-overview
description: High-level overview of Resik architecture and mental model
---

# Resik Overview

Resik is a Go backend boilerplate designed for:

- Modular monolith architecture
- Microservice-ready extraction
- REST + gRPC support
- Explicit ownership and dependency boundaries

## Core Idea

Resik is NOT just folder structure.

It is a **rule system for code generation and architecture consistency**.

## Main Layers

- app → bootstrap layer (REST / gRPC server)
- handler → transport adapter
- core → business logic + API contracts
- repositories → persistence layer
- entities → shared models
- ctx → request-scoped context
- internal → framework internals
- pkg → reusable libraries
- utils → application helpers
