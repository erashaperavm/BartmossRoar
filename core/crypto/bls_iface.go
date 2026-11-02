package crypto

import stdcrypto "crypto"

type BlsAlgo byte

const (
	Bls12381 BlsAlgo = iota
)

// Signature is the serialized form returned by a BLS implementation.
type Signature []byte

// BlsSigner signs messages and exposes the corresponding private-key operations.
type BlsSigner interface {
	Sign(msg []byte) (Signature, error)
	Equal(x stdcrypto.PrivateKey) bool
	MarshalBinary() ([]byte, error)
	Public() stdcrypto.PublicKey
	UnmarshalBinary(data []byte) error
	Validate() bool
}

type BlsPublicKey interface {
	Equal(x stdcrypto.PublicKey) bool
	MarshalBinary() ([]byte, error)
	UnmarshalBinary(data []byte) error
	Validate() bool
	Verify(msg []byte, sig Signature) bool
}

// BlsScheme provides key generation and signature aggregation for one algorithm.
type BlsScheme interface {
	Algorithm() BlsAlgo
	GenerateKeyPair() (BlsSigner, BlsPublicKey, error)
	Aggregate(sigs []Signature) (Signature, error)
	VerifyAggregate(publicKeys []BlsPublicKey, messages [][]byte, sig Signature) bool
}
