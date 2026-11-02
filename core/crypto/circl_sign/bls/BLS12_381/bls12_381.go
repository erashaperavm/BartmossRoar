package BLS12_381

import (
	stdcrypto "crypto"
	"crypto/rand"
	"errors"
	"fmt"
	"io"

	circlbls "github.com/cloudflare/circl/sign/bls"
	"github.com/erashaperavm/BartmossRoar/core/crypto"
)

const ikmSize = 32

type Scheme struct{}

type Signer struct {
	key *circlbls.PrivateKey[circlbls.G1]
}

type PublicKey struct {
	key *circlbls.PublicKey[circlbls.G1]
}

func (Scheme) Algorithm() crypto.BlsAlgo {
	return crypto.Bls12381
}

func (Scheme) GenerateKeyPair() (crypto.BlsSigner, crypto.BlsPublicKey, error) {
	ikm := make([]byte, ikmSize)
	if _, err := io.ReadFull(rand.Reader, ikm); err != nil {
		return nil, nil, fmt.Errorf("generate BLS12-381 key material: %w", err)
	}
	defer clear(ikm)

	sk, err := circlbls.KeyGen[circlbls.G1](ikm, nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("generate BLS12-381 key: %w", err)
	}

	signer := &Signer{key: sk}
	return signer, &PublicKey{key: sk.PublicKey()}, nil
}

func (s *Signer) Sign(msg []byte) (crypto.Signature, error) {
	if s == nil || s.key == nil || !s.key.Validate() {
		return nil, circlbls.ErrInvalidKey
	}
	return circlbls.Sign(s.key, msg), nil
}

func (s *Signer) Equal(x stdcrypto.PrivateKey) bool {
	return s != nil && s.key != nil && s.key.Equal(x)
}

func (s *Signer) MarshalBinary() ([]byte, error) {
	if s == nil || s.key == nil {
		return nil, circlbls.ErrInvalidKey
	}
	return s.key.MarshalBinary()
}

func (s *Signer) Public() stdcrypto.PublicKey {
	if s == nil || s.key == nil {
		return nil
	}
	return &PublicKey{key: s.key.PublicKey()}
}

func (s *Signer) UnmarshalBinary(data []byte) error {
	if s == nil {
		return errors.New("unmarshal BLS12-381 private key into nil signer")
	}
	key := new(circlbls.PrivateKey[circlbls.G1])
	if err := key.UnmarshalBinary(data); err != nil {
		return err
	}
	s.key = key
	return nil
}

func (s *Signer) Validate() bool {
	return s != nil && s.key != nil && s.key.Validate()
}

func (p *PublicKey) Equal(x stdcrypto.PublicKey) bool {
	return p != nil && p.key != nil && p.key.Equal(x)
}

func (p *PublicKey) MarshalBinary() ([]byte, error) {
	if p == nil || p.key == nil {
		return nil, circlbls.ErrInvalidKey
	}
	return p.key.MarshalBinary()
}

func (p *PublicKey) UnmarshalBinary(data []byte) error {
	if p == nil {
		return errors.New("unmarshal BLS12-381 public key into nil key")
	}
	key := new(circlbls.PublicKey[circlbls.G1])
	if err := key.UnmarshalBinary(data); err != nil {
		return err
	}
	if !key.Validate() {
		return circlbls.ErrInvalidKey
	}
	p.key = key
	return nil
}

func (p *PublicKey) Validate() bool {
	return p != nil && p.key != nil && p.key.Validate()
}

func (p *PublicKey) Verify(msg []byte, sig crypto.Signature) bool {
	return p != nil && p.key != nil && circlbls.Verify(p.key, msg, sig)
}

func (Scheme) Aggregate(sigs []crypto.Signature) (crypto.Signature, error) {
	encoded := make([]circlbls.Signature, len(sigs))
	for i, sig := range sigs {
		encoded[i] = sig
	}
	aggregated, err := circlbls.Aggregate(circlbls.G1{}, encoded)
	if err != nil {
		return nil, fmt.Errorf("aggregate BLS12-381 signatures: %w", err)
	}
	return aggregated, nil
}

func (Scheme) VerifyAggregate(
	publicKeys []crypto.BlsPublicKey,
	messages [][]byte,
	sig crypto.Signature,
) bool {
	if len(publicKeys) != len(messages) || len(publicKeys) == 0 {
		return false
	}

	keys := make([]*circlbls.PublicKey[circlbls.G1], len(publicKeys))
	for i, publicKey := range publicKeys {
		key, ok := publicKey.(*PublicKey)
		if !ok || key == nil || key.key == nil {
			return false
		}
		keys[i] = key.key
	}
	return circlbls.VerifyAggregate(keys, messages, sig)
}

var (
	_ crypto.BlsSigner    = (*Signer)(nil)
	_ crypto.BlsPublicKey = (*PublicKey)(nil)
	_ crypto.BlsScheme    = Scheme{}
)
