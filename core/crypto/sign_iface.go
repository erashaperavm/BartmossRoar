package crypto

type SignAlgo byte

const (
	SignMLDSA44 SignAlgo = iota
	SignMLDSA65
	SignMLDSA87
	SignEd25519
)

type Signer interface {
	// Sign 形参不含私钥，在 impl 时使用 struct 封装，ctx 目前用不到(ctx -> plaintext)
	// 确保 string(ctx) 安全
	Sign(msg []byte, ctx []byte) ([]byte, error)
	PublicKey() PublicKey
}

type PublicKey interface {
	Verify(msg, ctx, sig []byte) bool

	// Bytes 序列化公钥
	Bytes() ([]byte, error)

	Algorithm() SignAlgo
}
