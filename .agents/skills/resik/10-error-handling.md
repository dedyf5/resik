---
name: resik-error-handling
description: Standard error handling rules
---

# Error Handling

## Standard type

All layers MUST use:

*resPkg.Status

## Rules

- Repository returns Status
- Core returns Status
- Handler returns Status
- No raw error propagation

## Rule

Never return raw error across layers.
All errors must be normalized into Status.
