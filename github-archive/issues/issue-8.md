# Issue #8: MemoryStore grows without bound

**State:** OPEN  
**Created:** 2026-02-08T05:49:50Z  
**Updated:** 2026-02-08T21:30:52Z  
**URL:** https://github.com/aashh/saml/issues/8

## Description

`samlidp.MemoryStore` grows unboundedly, which could lead to memory exhaustion in long-running processes or under sustained load.

While `samlidp` is intended for testing and demos, adding a configurable size cap with FIFO eviction improves robustness without adding complexity. The fix adds a `MaxItems` field (default 10,000) with oldest-first eviction on `Put`.

Fix: [`fix/issue-08`](../tree/fix/issue-08)

