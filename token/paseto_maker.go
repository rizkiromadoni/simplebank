package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto      *paseto.V2
	symetricKey []byte
}

func NewPasetoMaker(symetricKey string) (Maker, error) {
	if len(symetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("symetric key must be %d length", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:      paseto.NewV2(),
		symetricKey: []byte(symetricKey),
	}

	return maker, nil
}

func (m *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	return m.paseto.Encrypt(m.symetricKey, payload, nil)
}

func (m *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	if err := m.paseto.Decrypt(token, m.symetricKey, payload, nil); err != nil {
		return nil, err
	}

	if err := payload.Valid(); err != nil {
		return nil, err
	}

	return payload, nil
}
