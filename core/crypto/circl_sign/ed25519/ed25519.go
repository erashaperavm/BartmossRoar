package ed25519

import (
	"errors"
	"fmt"

	"github.com/cloudflare/circl/sign/ed25519"
	"github.com/erashaperavm/BartmossRoar/core/crypto"
)

/*
	use CIRCL（cloudflare/circl/sign/ed25519）
	ed25519
*/

type Ed25519Signer struct {
	sk ed25519.PrivateKey
}

type Ed25519Verifyer struct {
	pk ed25519.PublicKey
}

func (e2s *Ed25519Signer) Sign(msg []byte, ctx []byte) ([]byte, error) {
	// the pack we use don't propagate en error but panic, we should check in advance
	givenLen := len(e2s.sk)
	if givenLen != ed25519.PrivateKeySize {
		return nil, errors.New(fmt.Sprintf("Ed25519 Sign failed: len(sk) != PrivateKeySize: %d", givenLen))
	}
	givenCtxLen := len(string(ctx))
	if givenCtxLen > ed25519.ContextMaxSize {
		return nil, errors.New(fmt.Sprintf("Ed25519 Sign failed: len(string(ctx)) > ed25519.ContextMaxSize: %d", givenCtxLen))
	}

	// SignPh allow blank ctx, but SignWithCtx not
	sig := ed25519.SignPh(e2s.sk, msg, string(ctx))

	return sig, nil
}

func (e2s *Ed25519Signer) PublicKey() crypto.PublicKey {
	return &Ed25519Verifyer{pk: e2s.sk.Public().(ed25519.PublicKey)}
}

func (e2v *Ed25519Verifyer) Verify(msg, ctx, sig []byte) bool {
	// the pack we use don't propagate en error but panic, we should check in advance
	givenCtxLen := len(string(ctx))
	if givenCtxLen > ed25519.ContextMaxSize {
		return false
	}

	// VerifyPh allow blank ctx, but VerifyWithCtx not
	return ed25519.VerifyPh(e2v.pk, msg, sig, string(ctx))
}

func (e2v *Ed25519Verifyer) Bytes() ([]byte, error) {
	return e2v.pk.MarshalBinary()
}

func (e2v *Ed25519Verifyer) Algorithm() crypto.SignAlgo {
	return crypto.SignEd25519
}

func New() (crypto.Signer, crypto.PublicKey, error) {
	pk, sk, err := ed25519.GenerateKey(nil)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Ed25519 New key failed: %v", err))
	} // rand 为 nil，则使用 crypto/rand.Reader

	return &Ed25519Signer{sk: sk}, &Ed25519Verifyer{pk: pk}, nil
}
