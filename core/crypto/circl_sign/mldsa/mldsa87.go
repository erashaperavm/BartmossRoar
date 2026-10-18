package mldsa

import (
	"errors"
	"fmt"

	"github.com/cloudflare/circl/sign/mldsa/mldsa87"
	"github.com/erashaperavm/BartmossRoar/core/crypto"
)

/*
	use CIRCL（cloudflare/circl/sign/mldsa）
	ml-dsa-87
*/

type Mldsa87Signer struct {
	sk mldsa87.PrivateKey
}

type Mldsa87Verifyer struct {
	pk mldsa87.PublicKey
}

func (m87s *Mldsa87Signer) Sign(msg []byte, ctx []byte) ([]byte, error) {
	sig := make([]byte, mldsa87.SignatureSize)

	err := mldsa87.SignTo(&m87s.sk, msg, ctx, false, sig)
	if err != nil {
		return nil, err
	}

	return sig, nil
}

func (m87s *Mldsa87Signer) PublicKey() crypto.PublicKey {
	return &Mldsa87Verifyer{*m87s.sk.Public().(*mldsa87.PublicKey)}
}

func (m87v *Mldsa87Verifyer) Verify(msg, ctx, sig []byte) bool {
	ok := mldsa87.Verify(&m87v.pk, msg, ctx, sig)
	if !ok {
		return false
	}

	return true
}

func (m87v *Mldsa87Verifyer) Bytes() ([]byte, error) {
	return m87v.pk.Bytes(), nil
}

func (m87v *Mldsa87Verifyer) Algorithm() crypto.SignAlgo {
	return crypto.SignMLDSA87
}

func New87() (crypto.Signer, crypto.PublicKey, error) {
	pk, sk, err := mldsa87.GenerateKey(nil)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("MLDSA87 New key failed: %v", err))
	} // rand 为 nil，则使用 crypto/rand.Reader

	return &Mldsa87Signer{sk: *sk}, &Mldsa87Verifyer{pk: *pk}, nil
}
