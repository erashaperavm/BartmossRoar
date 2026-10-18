package chacha20poly1305

import (
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/erashaperavm/BartmossRoar/core/crypto"
	"golang.org/x/crypto/chacha20poly1305"
)

/*
	use golang.org/x/crypto/chacha20poly1305
	ChaCha20Poly1305 (RFC 8439)
*/

type ChaCha20poly1305Cipher struct {
	aead cipher.AEAD
}

func (c2p1c *ChaCha20poly1305Cipher) Encrypt(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, c2p1c.aead.NonceSize(), c2p1c.aead.NonceSize()+len(plaintext)+c2p1c.aead.Overhead())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	// Encrypt the message and append the ciphertext to the nonce.
	encryptedMsg := c2p1c.aead.Seal(nonce, nonce, plaintext, nil)
	return encryptedMsg, nil
}

func (c2p1c *ChaCha20poly1305Cipher) Decrypt(ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < c2p1c.aead.NonceSize() {
		return nil, errors.New("chacha20poly1305 decrypt failed: ciphertext too short")
	}

	// Split nonce and ciphertext.
	nonce, ciphertext := ciphertext[:c2p1c.aead.NonceSize()], ciphertext[c2p1c.aead.NonceSize():]

	// Decrypt the message and check it wasn't tampered with.
	plaintext, err := c2p1c.aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("chacha20poly1305 decrypt failed: %v", err))
	}

	return plaintext, nil
}

func (c2p1c *ChaCha20poly1305Cipher) NonceSize() int {
	return chacha20poly1305.NonceSize
}

func (c2p1c *ChaCha20poly1305Cipher) Algorithm() crypto.EncAlgo {
	return crypto.EncChaCha20Poly1305
}

func New(key crypto.Key) (*ChaCha20poly1305Cipher, error) {
	aead, err := chacha20poly1305.New(key[:])
	if err != nil {
		return nil, errors.New(fmt.Sprintf("ChaCha20poly1305Cipher new failed: %v", err))
	}
	return &ChaCha20poly1305Cipher{aead: aead}, nil
}
