package hybrid

import (
	"errors"

	"github.com/cloudflare/circl/kem"
	"github.com/cloudflare/circl/kem/hybrid"
	"github.com/erashaperavm/BartmossRoar/core/crypto"
)

// 确保实现你的接口
var _ crypto.KEMScheme = (*X25519MLKEM768Scheme)(nil)

type X25519MLKEM768Pk struct {
	pk kem.PublicKey
}

type X25519MLKEM768Sk struct {
	sk kem.PrivateKey
}

type X25519MLKEM768Scheme struct {
	scheme kem.Scheme
}

func NewScheme() crypto.KEMScheme {
	return &X25519MLKEM768Scheme{scheme: hybrid.X25519MLKEM768()}
}

// ---------------- PublicKey 实现 ----------------

func (x2m7p *X25519MLKEM768Pk) MarshalBinary() ([]byte, error) {
	return x2m7p.pk.MarshalBinary()
}

func (x2m7p *X25519MLKEM768Pk) Scheme() crypto.KEMScheme {
	// 每次返回一个新的 wrapper
	return &X25519MLKEM768Scheme{scheme: x2m7p.pk.Scheme()}
}

func (x2m7p *X25519MLKEM768Pk) Equal(other crypto.KEMPublicKey) bool {
	// 先断言对方是不是同一个实现
	otherPk, ok := other.(*X25519MLKEM768Pk)
	if !ok {
		return false
	}
	return x2m7p.pk.Equal(otherPk.pk)
}

// ---------------- PrivateKey 实现 ----------------

func (x2m7s *X25519MLKEM768Sk) MarshalBinary() ([]byte, error) {
	return x2m7s.sk.MarshalBinary()
}

func (x2m7s *X25519MLKEM768Sk) Public() crypto.KEMPublicKey {
	// 包装底层返回的接口
	return &X25519MLKEM768Pk{pk: x2m7s.sk.Public()}
}

func (x2m7s *X25519MLKEM768Sk) Scheme() crypto.KEMScheme {
	return &X25519MLKEM768Scheme{scheme: x2m7s.sk.Scheme()}
}

// ---------------- Scheme 实现 ----------------

func (s *X25519MLKEM768Scheme) GenerateKeyPair() (crypto.KEMPublicKey, crypto.KEMPrivateKey, error) {
	pk, sk, err := s.scheme.GenerateKeyPair()
	if err != nil {
		return nil, nil, err
	}
	return &X25519MLKEM768Pk{pk: pk}, &X25519MLKEM768Sk{sk: sk}, nil
}

func (s *X25519MLKEM768Scheme) Encapsulate(pk crypto.KEMPublicKey) (ciphertext, sharedSecret []byte, err error) {
	// 必须把外层接口拆包成底层实现
	wrappedPk, ok := pk.(*X25519MLKEM768Pk)
	if !ok {
		return nil, nil, errors.New("invalid public key type")
	}
	return s.scheme.Encapsulate(wrappedPk.pk)
}

func (s *X25519MLKEM768Scheme) Decapsulate(sk crypto.KEMPrivateKey, ciphertext []byte) (sharedSecret []byte, err error) {
	wrappedSk, ok := sk.(*X25519MLKEM768Sk)
	if !ok {
		return nil, errors.New("invalid private key type")
	}
	return s.scheme.Decapsulate(wrappedSk.sk, ciphertext)
}

func (s *X25519MLKEM768Scheme) UnmarshalBinaryPublicKey(data []byte) (crypto.KEMPublicKey, error) {
	pk, err := s.scheme.UnmarshalBinaryPublicKey(data)
	if err != nil {
		return nil, err
	}
	return &X25519MLKEM768Pk{pk: pk}, nil
}

func (s *X25519MLKEM768Scheme) UnmarshalBinaryPrivateKey(data []byte) (crypto.KEMPrivateKey, error) {
	sk, err := s.scheme.UnmarshalBinaryPrivateKey(data)
	if err != nil {
		return nil, err
	}
	return &X25519MLKEM768Sk{sk: sk}, nil
}

func (s *X25519MLKEM768Scheme) Name() string {
	return s.scheme.Name()
}
