package refresh_token

import (
	"crypto/sha256"
	"encoding/hex"
)

type RefreshTokenHasher struct {
	pepper []byte
}

func NewRefreshTokenHasher(pepper string) *RefreshTokenHasher {
	return &RefreshTokenHasher{pepper: []byte(pepper)}
}

func (h *RefreshTokenHasher) Hash(token string) string {

	b := append([]byte(token), h.pepper...)
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
