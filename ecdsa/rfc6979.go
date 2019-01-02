package ecdsa

import (
	"bytes"
	"crypto"
	"crypto/elliptic"
	"crypto/hmac"
	"math/big"
)

func mac(hashFunc crypto.Hash, k, m []byte) []byte {
	h := hmac.New(hashFunc.New, k)
	h.Write(m)

	return h.Sum(nil)
}

// https://tools.ietf.org/html/rfc6979#section-2.3.2
func bitsToInt(b []byte, qlen int) *big.Int {
	v := new(big.Int).SetBytes(b)
	vlen := len(b) * 8

	if vlen > qlen {
		v = new(big.Int).Rsh(v, uint(vlen-qlen))
	}

	return v
}

// https://tools.ietf.org/html/rfc6979#section-2.3.3
func intToOctets(v *big.Int, rlen int) []byte {
	ret := v.Bytes()

	if len(ret) < rlen {
		return append(make([]byte, rlen-len(ret)), ret...)
	}

	if len(ret) > rlen {
		panic("ecdsa.intToOctets: len(ret) is larger than rlen")
	}

	return ret
}

// https://tools.ietf.org/html/rfc6979#section-2.3.4
func bitsToOctets(b []byte, q *big.Int, qlen, rlen int) []byte {
	z1 := bitsToInt(b, qlen)
	z2 := new(big.Int).Sub(z1, q)

	if z2.Sign() < 0 {
		return intToOctets(z1, rlen)
	}

	return intToOctets(z2, rlen)
}

// https://git.io/fpxpl
func hashToInt(hash []byte, curve elliptic.Curve) *big.Int {
	orderBits := curve.Params().N.BitLen()
	orderBytes := (orderBits + 7) / 8
	if len(hash) > orderBytes {
		hash = hash[:orderBytes]
	}

	ret := new(big.Int).SetBytes(hash)
	rem := len(hash)*8 - orderBits
	if rem > 0 {
		ret.Rsh(ret, uint(rem))
	}

	return ret
}

var one = big.NewInt(1)

// https://tools.ietf.org/html/rfc6979#section-2.4
func sign(priv *PrivateKey, hash []byte, h crypto.Hash) *Signature {
	curve := priv.Curve

	q := curve.Params().N
	x := priv.D

	hlen := h.Size()
	qlen := q.BitLen()
	rlen := (qlen + 7) >> 3

	b := append(intToOctets(x, rlen), bitsToOctets(hash, q, qlen, rlen)...)

	V := bytes.Repeat([]byte{0x01}, hlen)
	K := bytes.Repeat([]byte{0x00}, hlen)
	K = mac(h, K, append(append(V, 0x00), b...))
	V = mac(h, K, V)
	K = mac(h, K, append(append(V, 0x01), b...))
	V = mac(h, K, V)

	for {
		var T []byte

		for len(T) < qlen/8 {
			V = mac(h, K, V)
			T = append(T, V...)
		}

		k := bitsToInt(T, qlen)

		if k.Cmp(one) >= 0 && k.Cmp(q) < 0 {
			r, _ := curve.ScalarBaseMult(k.Bytes())
			r.Mod(r, q)

			if r.Sign() == 0 {
				continue
			}

			e := hashToInt(hash, curve)
			i := new(big.Int).ModInverse(k, q)

			s := new(big.Int)
			s.Mul(x, r)
			s.Add(s, e)
			s.Mul(s, i)
			s.Mod(s, q)

			if s.Sign() != 0 {
				return &Signature{r, s}
			}
		}

		K = mac(h, K, append(V, 0x00))
		V = mac(h, K, V)
	}
}
