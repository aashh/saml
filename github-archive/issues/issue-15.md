# Issue #15: Example code uses request body as randomness source

**State:** OPEN  
**Created:** 2026-02-08T05:49:57Z  
**Updated:** 2026-02-08T05:49:57Z  
**URL:** https://github.com/aashh/saml/issues/15

## Description

`example/service.go` `CreateLink()` reads 8 bytes from the HTTP request body and uses them as the shortlink ID:

```go
randomness := make([]byte, 8)
r.Body.Read(randomness)
l := Link{ShortLink: base64.RawURLEncoding.EncodeToString(randomness), ...}
```

The link IDs are derived from request body data, making them predictable. A client can compute the bytes needed to collide with any existing shortlink.

While this is example code, using `crypto/rand.Read` instead demonstrates the correct pattern for generating unpredictable identifiers.

Fix: [`fix/issue-15`](../tree/fix/issue-15)

