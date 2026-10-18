package chacha20poly1305

import (
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/erashaperavm/BartmossRoar/core/crypto"
	"golang.org/x/crypto/chacha20poly1305"
)

// todo import "golang.org/x/crypto/chacha20poly1305"

/*
	use golang.org/x/crypto/chacha20poly1305
	XChaCha20Poly1305 (RFC 8439)
*/

type XChaCha20poly1305Cipher struct {
	aead cipher.AEAD
}

func (xc2p1c *XChaCha20poly1305Cipher) Encrypt(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, xc2p1c.aead.NonceSize(), xc2p1c.aead.NonceSize()+len(plaintext)+xc2p1c.aead.Overhead())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	// Encrypt the message and append the ciphertext to the nonce.
	encryptedMsg := xc2p1c.aead.Seal(nonce, nonce, plaintext, nil)
	return encryptedMsg, nil
}

func (xc2p1c *XChaCha20poly1305Cipher) Decrypt(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < xc2p1c.aead.NonceSize() {
		return nil, errors.New("xchacha20poly1305 decrypt failed: ciphertext too short")
	}

	// Split nonce and ciphertext.
	nonce, ciphertext := ciphertext[:xc2p1c.aead.NonceSize()], ciphertext[xc2p1c.aead.NonceSize():]

	// Decrypt the message and check it wasn't tampered with.
	plaintext, err := xc2p1c.aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("xchacha20poly1305 decrypt failed: %v", err))
	}

	return plaintext, nil
}

func (xc2p1c *XChaCha20poly1305Cipher) NonceSize() int {
	return chacha20poly1305.NonceSizeX
}

func (xc2p1c *XChaCha20poly1305Cipher) Algorithm() crypto.EncAlgo {
	return crypto.EncXChaCha20Poly1305
}

func NewX(key crypto.Key) (*XChaCha20poly1305Cipher, error) {
	aead, err := chacha20poly1305.NewX(key[:])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("XChaCha20poly1305Cipher new failed: %v", err))
	}
	return &XChaCha20poly1305Cipher{aead: aead}, nil
}
