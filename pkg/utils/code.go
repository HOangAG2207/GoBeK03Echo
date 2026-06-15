package pkgutils

import (
	"bytes"
	"crypto/rand"
	"math/big"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

//go:generate mockery --name CodeGenerator --filename code.go --output ./mocks
type CodeGenerator interface {
	GenerateCode(length int) (string, error)
}
type randomCodeGenerator struct {
}

func NewCodeGenerator() CodeGenerator {
	return &randomCodeGenerator{}
}
func (r *randomCodeGenerator) GenerateCode(length int) (string, error) {
	var strBuilder bytes.Buffer

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		strBuilder.WriteByte(charset[randomIndex.Int64()])
	}
	return strBuilder.String(), nil
}
