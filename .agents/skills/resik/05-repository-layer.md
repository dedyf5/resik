---
name: resik-repository-layer
description: Repository rules and behavior
---

# Repository Layer

## Purpose

Data access layer (DB / cache / external systems)

## Rules

- One method = one operation
- No repository calling repository
- No dependency on core
- No transaction handling

## Transaction rule

Transaction MUST be handled in core.

## Param rule

Use param struct only when:

- query is complex
- reusable
- logically grouped
- likely to evolve

Do NOT create param just because arguments are many.
