package generator

import (
	"crypto/rand"
	"log"
	"math/big"

	"github.com/sqids/sqids-go"
)

type ShortCodeGenerator struct {
	sqid *sqids.Sqids
}

// NewShortCodeGenerator создаёт новый экземпляр генератора
func NewShortCodeGenerator() *ShortCodeGenerator {
	s, err := sqids.New(sqids.Options{
		Alphabet:  "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
		MinLength: 6,
	})
	if err != nil {
		log.Fatalf("failed to init sqids: %v", err)
	}
	return &ShortCodeGenerator{sqid: s}
}

// Generate генерирует короткий код по числовому идентификатору
func (g *ShortCodeGenerator) Generate() string {
	n, err := rand.Int(rand.Reader, big.NewInt(1<<53)) // до ~9e15
	if err != nil {
		log.Printf("failed to generate random number: %v", err)
		return ""
	}

	code, err := g.sqid.Encode([]uint64{n.Uint64()})
	if err != nil {
		log.Printf("failed to encode sqid: %v", err)
		return ""
	}

	return code
}
