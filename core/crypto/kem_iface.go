package crypto

type KemAlgo byte

const (
	X25519MLKEM768 KemAlgo = iota
)

// KEMPublicKey 封装公钥，需支持序列化以便传输
type KEMPublicKey interface {
	MarshalBinary() ([]byte, error)
	Scheme() KEMScheme
	Equal(key KEMPublicKey) bool
}

// KEMPrivateKey 封装私钥
type KEMPrivateKey interface {
	MarshalBinary() ([]byte, error)
	Public() KEMPublicKey
	Scheme() KEMScheme
}

// KEMScheme 是算法的核心操作接口
type KEMScheme interface {
	// GenerateKeyPair 生成新的密钥对
	GenerateKeyPair() (KEMPublicKey, KEMPrivateKey, error)
	// Encapsulate 用接收方公钥生成共享密钥与密文
	Encapsulate(pk KEMPublicKey) (ciphertext, sharedSecret []byte, err error)
	// Decapsulate 用私钥从密文恢复共享密钥
	Decapsulate(sk KEMPrivateKey, ciphertext []byte) (sharedSecret []byte, err error)
	// UnmarshalBinaryPublicKey 从字节反序列化密钥
	UnmarshalBinaryPublicKey([]byte) (KEMPublicKey, error)
	UnmarshalBinaryPrivateKey([]byte) (KEMPrivateKey, error)
	// Name 返回算法名称，便于日志与协商
	Name() string
}

// KEM 门面接口，便于上层代码像使用 Hasher 一样使用
type KEM interface {
	Scheme() KEMScheme
	GenerateKeyPair() (KEMPublicKey, KEMPrivateKey, error)
	Algorithm() KemAlgo
}
