---
name: resik-core-layer
description: Core business logic rules
---

# Core Layer

## Purpose

- Business logic
- API contracts (request / response)
- Use case implementation

## CRUD organization rule

All CRUD must be in ONE file:

organization.go:
- CreateOrganization
- UpdateOrganization
- DeleteOrganization
- GetOrganization
- ListOrganization

## Do NOT:

- split CRUD into multiple files
