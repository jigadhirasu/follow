package uuid

import (
	"crypto/rand"
	"encoding/hex"
)

// UUID Represents a UUID
type UUID [16]byte

func (u *UUID) String() string {
	result := hex.EncodeToString(u[:])
	result = result[:8] + "-" + result[8:12] + "-" + result[12:16] + "-" + result[16:20] + "-" + result[20:]
	return result
}
func (u *UUID) setVersion(v byte) {
	u[6] = (v << 4) | (u[6] & 0xf)
}

func (u *UUID) setVariant() {
	// Clear the first two bits of the byte
	u[8] = u[8] & 0x3f
	// Set the first two bits of the byte to 10
	u[8] = u[8] | 0x80
}

func (u *UUID) Hex() string {
	return hex.EncodeToString(u[:])
}

func New() string {
	id := Gen()
	return id.String()
}

func Gen() *UUID {
	bytes := make([]byte, 16)
	_, _ = rand.Read(bytes)

	result := &UUID{}
	for i, v := range bytes {
		result[i] = v
	}

	result.setVersion(4)
	result.setVariant()
	return result
}

func FromHex(str string) *UUID {
	bytes, err := hex.DecodeString(str)
	if err != nil {
		panic(err)
	}

	result := &UUID{}
	for i, v := range bytes {
		result[i] = v
	}

	return result
}
