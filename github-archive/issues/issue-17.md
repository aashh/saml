# Issue #17: TimeNow and RandReader are mutable globals

**State:** CLOSED  
**Created:** 2026-02-08T05:49:59Z  
**Updated:** 2026-02-08T05:54:29Z  
**URL:** https://github.com/aashh/saml/issues/17

## Description

`saml.TimeNow` and `saml.RandReader` are exported, mutable package-level variables. Any code in the same process can override them -- freezing time to disable expiration checks, or replacing the random source with something predictable.

This is a real design tension: the mutability exists to support testing (mocking time and randomness), and it's used pervasively across the codebase. Moving these to instance fields would touch nearly every type and function in the library. That's a v3-scale API break for a threat model that mostly requires a compromised dependency in the same process (at which point you have bigger problems).

Not planning to change this.

