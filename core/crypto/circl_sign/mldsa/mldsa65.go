package mldsa

import (
	"errors"
	"fmt"

	"github.com/cloudflare/circl/sign/mldsa/mldsa65"
	"github.com/erashaperavm/BartmossRoar/core/crypto"
)

/*
	use CIRCL（cloudflare/circl/sign/mldsa）
	ml-dsa-65
*/

type Mldsa65Signer struct {
	sk mldsa65.PrivateKey
}

type Mldsa65Verifyer struct {
	pk mldsa65.PublicKey
}

func (m65s *Mldsa65Signer) Sign(msg []byte, ctx []byte) ([]byte, error) {
	sig := make([]byte, mldsa65.SignatureSize)

	err := mldsa65.SignTo(&m65s.sk, msg, ctx, false, sig)
	if err != nil {
		return nil, err
	}

	return sig, nil
}

func (m65s *Mldsa65Signer) PublicKey() crypto.PublicKey {
	return &Mldsa65Verifyer{*m65s.sk.Public().(*mldsa65.PublicKey)}
}

func (m65v *Mldsa65Verifyer) Verify(msg, ctx, sig []byte) bool {
	ok := mldsa65.Verify(&m65v.pk, msg, ctx, sig)
	if !ok {
		return false
	}

	return true
}

func (m65v *Mldsa65Verifyer) Bytes() ([]byte, error) {
	return m65v.pk.Bytes(), nil
}

func (m65v *Mldsa65Verifyer) Algorithm() crypto.SignAlgo {
	return crypto.SignMLDSA65
}

func New65() (crypto.Signer, crypto.PublicKey, error) {
	pk, sk, err := mldsa65.GenerateKey(nil)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("MLDSA65 New key failed: %v", err))
	} // rand 为 nil，则使用 crypto/rand.Reader

	return &Mldsa65Signer{sk: *sk}, &Mldsa65Verifyer{pk: *pk}, nil
}
