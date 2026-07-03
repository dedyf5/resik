---
name: resik-pkg-utils
description: Difference between pkg and utils
---

# pkg

Pure reusable Go libraries.

Rules:
- Must NOT depend on Resik
- Prefer std library only
- Portable across projects

---

# utils

Resik-specific helpers.

Rules:
- Can depend on pkg
- Can depend on third-party libs
- Used to reduce duplication in application logic
