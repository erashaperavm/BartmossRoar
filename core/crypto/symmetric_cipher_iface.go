package crypto

type EncAlgo byte

type Key [32]byte

const (
	EncChaCha20Poly1305 EncAlgo = iota
	EncXChaCha20Poly1305
)

// SymmetricCipher 是一个 AEAD 对称加密接口。
//
// 约定：
//   - Encrypt 返回自包含密文：nonce || ciphertext || tag
//   - Decrypt 从密文中读 nonce，解密失败返回 error，不返回明文
//   - 实现内部持有密钥，调用方不传密钥
type SymmetricCipher interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
	NonceSize() int
	Algorithm() EncAlgo
}

// NewSymmetricCipher 根据算法和密钥构造一个 SymmetricCipher。
// key 的长度和算法相关（ChaCha/XChaCha 都是 32 字节）。
func NewSymmetricCipher(algo EncAlgo, key []byte) (SymmetricCipher, error) {
	switch algo {

	}

	return nil, nil
}
