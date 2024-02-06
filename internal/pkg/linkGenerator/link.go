package linkGenerator

import (
	"crypto/sha256"
	"math/big"
)

type LinkHash interface {
	GenLink(longLink string) string
}

type BaseLink struct {
	Alphabet string
	Length   int
}

func NewLinkHash(alphabet string, length int) *BaseLink {
	return &BaseLink{
		Alphabet: alphabet,
		Length:   length,
	}
}

func (l *BaseLink) GenLink(longLink string) string {
	hash := sha256.Sum256([]byte(longLink))
	base := big.NewInt(int64(l.Length))

	bufferInt := new(big.Int).SetBytes(hash[:])
	var encoded string

	for bufferInt.Cmp(big.NewInt(0)) > 0 && len(encoded) < 10 {
		mod := new(big.Int)
		bufferInt.DivMod(bufferInt, base, mod)
		encoded += string(l.Alphabet[mod.Int64()])
	}

	return encoded
}
