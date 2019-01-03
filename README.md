## Overview

This program was designed to build Android app releases on CI without having to
store and manage any certificates or private keys.

Instead, certificate authority is derived from `SEED_CA` environment variable,
which is an arbitrary string. You can generate one with, for example,
`head -c 30 /dev/urandom | base64`.

Running the program with the environment variable set will deterministically
spawn PEM-encoded `ca.key` and `ca.crt` files in the current working directory
and print `ca.key\nca.crt\n` to standard output.

## Rationale

Most CI dispatchers support protected variables, environment variables that can
only be seen by people with ownership rights over designated repository.

Safely provisioning a file to a stateless environment such as CI would require
setting protected variables anyway, but the sheer complexity and coupling
potential of that alone is going to be a major centralizing/fork-hostile power.

This work is part of the effort to make infrastructure concrete, stateless,
and forkable just like projects that it serves. All effective build state is
encapsulated into a tiny concrete string, with no additional setup involved
other than setting it to some (any) value.

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

ECDSA was chosen because it has better compatibility, especially with older
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

## Precautions

Please make sure to back up your seed. If you lose it, it won't be possible to
push any updates to Android apps that were signed with the lost seed, other uses
may share similar consequences.

Sign with smartcards for high-threat models. Be aware that whoever hosts CI
dispatcher and builders will be able to access the seed and forge signatures.

## Stability

*This program has not stabilized yet. Once it does, it will be moved from
`unstable` namespace and will be effectively set in stone.*

Output of this program must never change between versions. If that ever happens,
please file a bug report. Non-compatible versions of this program will go to a
separate repository.

`golang.org/x/crypto` is pinned down. If any breaking changes happen to Go
standard library, which is the only effectively unpinned dependency, the latest
good versions of required broken packages will be vendored into this project.

The most likely thing to warrant vendoring is `crypto/ecdsa` key generation.

## Adopters

This is a list of projects that use Seed CA (if you do, send a merge request):

- [Noise](https://gitlab.com/prism-break/noise), rebranded and deblobbed Signal
client for Android (used to sign the `.apk` and to authenticate F-Droid repo)

## Thanks

Coda Hale, for [implementing RFC 6979 in Go](https://github.com/codahale/rfc6979).
That's what local `ecdsa` package is based on.

Thomas Pornin, for designing a safe reverse-compatible deterministic scheme for
ECDSA signatures. None of this would be possible without that great and
relatively recent work.
