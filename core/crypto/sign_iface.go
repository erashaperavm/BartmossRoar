package crypto

type SignAlgo byte

const (
	SignMLDSA44 SignAlgo = iota
	SignMLDSA65
	SignMLDSA87
	SignEd25519
)

type Signer interface {
	// Sign 形参不含私钥，在 impl 时使用 struct 封装
	Sign(msg []byte) ([]byte, error)
	PublicKey() PublicKey
}

type PublicKey interface {
	Verify(msg, sig []byte) error

	// Bytes 序列化公钥
	Bytes() []byte

	Algorithm() SignAlgo
}

func NewSigner(algo SignAlgo, priv any) (Signer, error) {
	switch algo {

	}

	return nil, nil
}
