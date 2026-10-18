package mldsa

import (
	"errors"
	"fmt"

	"github.com/cloudflare/circl/sign/mldsa/mldsa44"
	"github.com/erashaperavm/BartmossRoar/core/crypto"
)

/*
	use CIRCL（cloudflare/circl/sign/mldsa）
	ml-dsa-44
*/

type Mldsa44Signer struct {
	sk mldsa44.PrivateKey
}

type Mldsa44Verifyer struct {
	pk mldsa44.PublicKey
}

func (m44s *Mldsa44Signer) Sign(msg []byte, ctx []byte) ([]byte, error) {
	sig := make([]byte, mldsa44.SignatureSize)

	err := mldsa44.SignTo(&m44s.sk, msg, ctx, false, sig)
	if err != nil {
		return nil, err
	}

	return sig, nil
}

func (m44s *Mldsa44Signer) PublicKey() crypto.PublicKey {
	return &Mldsa44Verifyer{*m44s.sk.Public().(*mldsa44.PublicKey)}
}

func (m44v *Mldsa44Verifyer) Verify(msg, ctx, sig []byte) bool {
	return mldsa44.Verify(&m44v.pk, msg, ctx, sig)
}

func (m44v *Mldsa44Verifyer) Bytes() ([]byte, error) {
	return m44v.pk.Bytes(), nil
}

func (m44v *Mldsa44Verifyer) Algorithm() crypto.SignAlgo {
	return crypto.SignMLDSA44
}

func New44() (crypto.Signer, crypto.PublicKey, error) {
	pk, sk, err := mldsa44.GenerateKey(nil)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("MLDSA44 New key failed: %v", err))
	} // rand 为 nil，则使用 crypto/rand.Reader

	return &Mldsa44Signer{sk: *sk}, &Mldsa44Verifyer{pk: *pk}, nil
}
