package shortid

import (
	"crypto/rand"
	"time"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const base = uint64(len(alphabet))

type Generator struct {
	nodeID uint8
	length int
	maxVal uint64
}

func NewGenerator(nodeID int, length int) *Generator {
	if nodeID < 0 {
		nodeID = 0
	}
	if nodeID > 15 {
		nodeID = 15
	}
	if length < 7 {
		length = 7
	}

	var maxVal uint64 = 1
	for i := 0; i < length; i++ {
		maxVal *= base
	}

	return &Generator{
		nodeID: uint8(nodeID),
		length: length,
		maxVal: maxVal,
	}
}

func (g *Generator) Generate() (string, error) {
	ts := uint64(time.Now().Unix()) & 0xFFFFFFFF

	randBits, err := randomBits28()
	if err != nil {
		return "", err
	}

	value := (uint64(g.nodeID&0xF) << 60) | (ts << 28) | randBits
	value = value % g.maxVal

	return g.encode(value), nil
}

func (g *Generator) encode(value uint64) string {
	buf := make([]byte, g.length)
	for i := g.length - 1; i >= 0; i-- {
		buf[i] = alphabet[value%base]
		value /= base
	}
	return string(buf)
}

func randomBits28() (uint64, error) {
	var b [4]byte
	if _, err := rand.Read(b[:]); err != nil {
		return 0, err
	}
	val := uint64(b[0])<<24 | uint64(b[1])<<16 | uint64(b[2])<<8 | uint64(b[3])
	return val & 0x0FFFFFFF, nil
}
