## Usage

This program was designed to build Android app releases on CI without having to
store and manage any certificates or private keys.

Instead, certificate authority is derived from `SEED_CA` environment variable,
which is an arbitrary string. You can generate one with, for example:

```
head -c 30 /dev/urandom | base64
```

Running the program with the environment variable set will deterministically
spawn PEM-encoded `ca.key` and `ca.crt` files in the current working directory
and print file names to standard output.

## Rationale

Most CI dispatchers support protected variables, environment variables that can
only be seen by people with ownership rights over designated repository.

As a rule, safely provisioning a file to a stateless environment such as CI would
require setting protected variables anyway, but the sheer complexity of that alone
is going to be a major fork-hostile power.

This project makes forks as friendly as possible, encapsulating all of effective
build state into a concrete string, with no additional setup involved other than
setting it to some (any) value.

## Cryptography

Seed is stretched to 73 bytes via [Argon2id KDF][argon2] with parameters t = 2,
m = 1536MB, p = 4 and salt `jcD_:}VIy'_q)$>^@=q&2ywx6`. Resulting byte buffer is
then used as an entropy pool to generate ECDSA private key on NIST P-521 curve
via Go `crypto/ecdsa` package.

Private key self-signs SHA-512 hash of a CA certificate following deterministic
[RFC 6979][rfc6979] scheme. Not before date matches Unix epoch, not after date
is not specified (see [RFC 5280 4.1.2.5][rfc5280-4.1.2.5]), CA subject is X.250
compatible with every field being blank.

### Why not Argon2d?

Timing attacks are not a threat for this program, so Argon2d could have been
a better choice. Unfortunately, `golang.org/x/crypto/argon2` package doesn't
expose Argon2d, only Argon2i and Argon2id.

### Why not EdDSA?

EdDSA over Curve25519 is a way more misimplementation-resistant signature scheme
than ECDSA over any of the NIST curves. Please prefer the former when designing
new formats and protocols.

ECDSA was chosen because it is has better compatibility, especially with older
software. Also, this program being non-interactive and deterministic protects
against most attack vectors on ECDSA.

Among NIST curves, P-521 seems to be the safest bet, even with DJB saying
[it's the only one using a "nice" prime][djb-p521] (grep for `fair`).

[argon2]: https://tools.ietf.org/html/draft-irtf-cfrg-argon2-04
[rfc6979]: https://tools.ietf.org/html/rfc6979
[rfc5280-4.1.2.5]: https://tools.ietf.org/html/rfc5280#section-4.1.2.5
[djb-p521]: https://blog.cr.yp.to/20140323-ecdsa.html

## Parameters

Some projects might warrant different KDF parameters, curve selection, or X.509
certificate fields. If you think that applies to your use case, fork this repo,
change parameters in `main.go` file, commit, and push to something like GitHub
or GitLab.

## Stability

*This program has not stabilized yet. Once it is, it will be moved from unstable
namespace and will be effectively set in stone.*

Output of this program must never change between versions. If that ever happens,
please file a bug report. Non-compatible versions of this program will go to a
separate repository.

`golang.org/x/crypto` dependency is pinned down. If any breaking changes happen
in Go standard library, which is the only unpinned dependency, we will vendor in
the last good version of required packages.

Most likely thing that might warrant vendoring is `crypto/ecdsa` key generation.

## Thanks

Coda Hale, for [implementing RFC 6979 in Go](https://github.com/codahale/rfc6979).
That's what local `ecdsa` package is based on.

Thomas Pornin, for designing a safe reverse-compatible deterministic scheme for
ECDSA signatures. None of this would be possible without that great and
relatively recent work.
