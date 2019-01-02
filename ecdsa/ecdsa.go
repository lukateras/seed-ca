// ECDSA implementation with deterministic signatures (RFC 6979).
package ecdsa

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"io"
)

type PrivateKey struct {
	*ecdsa.PrivateKey
}

func GenerateKey(curve elliptic.Curve, rand io.Reader) (*PrivateKey, error) {
	priv, err := ecdsa.GenerateKey(curve, rand)
	return &PrivateKey{priv}, err
}

func (priv *PrivateKey) Sign(_ io.Reader, hash []byte, opts crypto.SignerOpts) ([]byte, error) {
	return sign(priv, hash, opts.HashFunc()).Marshal()
}

// Signs a message's hash (specified by crypto.Hash) using private key priv.
// Returns a deterministic signature.
func Sign(priv *PrivateKey, message []byte, h crypto.Hash) *Signature {
	m := h.New()
	m.Write(message)

	return sign(priv, m.Sum(nil), h)
}
