package ecdsa

import (
	"encoding/asn1"
	"math/big"
)

type Signature struct {
	R, S *big.Int
}

// Tests if two Signatures are equal or not.
func (sig *Signature) Equal(sig2 *Signature) bool {
	return sig.R.Cmp(sig2.R) == 0 && sig.S.Cmp(sig2.S) == 0
}

// Marshals Signature to its standard ASN.1 DER representation.
func (sig *Signature) Marshal() ([]byte, error) {
	return asn1.Marshal(*sig)
}
