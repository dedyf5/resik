---
name: resik-ctx
description: Request-scoped context wrapper
---

# ctx (Request Context)

## Purpose

Holds request-scoped data for a single request lifecycle.

## Lifecycle rule

Only handler creates ctx.Ctx

core and repository only consume it

## Contains

- Context
- Language
- Logger
- JWT claims

## Rule

ctx must NOT store global or singleton state
