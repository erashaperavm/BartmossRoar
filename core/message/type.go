package message

// type.go declares types

type ID [32]byte
type Signal = byte
type Proto = byte

const (
	SYNC_SIG_KEY  Signal = iota
	SYNC_KEM_KEY
	SYNC_ENC_KEY
	SEND_MSG

	// DUMMY is unseen
)

const (
	HTTP Proto = iota
	PURE_TEXT
)

type OriginDataPack struct {
	/*
	SignalBitCode ciphertext
	receiver should try to decrypt it, then can transfer to handler
	*/
	SignalBitCode Signal

	/*
	TargetTag plaintext
	receiver should compute the shortest path to target
	*/
	TargetTag []ID

	/*
	LatestTopTimeStamp plaintext
	the latest receiver should update it 
	*/
	LatestTopTimeStamp uint64

	/*
	Payload ciphertext
	(blank)
	*/
	Payload Payload

	/*
	TempSig plaintext
	sender should use temp sig pub key to sign Payload + SignalBitCode
	*/
	TempSig []byte 
}

type Payload struct {
	/*
	TempSig
	the real send timestamp
	*/
	SendTimestamp uint64

	/*
	Protocol
	help target node handler to distinct them
	*/
	Protocol Proto

	/*
	Content
	meaningful data
	*/
	Content []byte

	/*
	TempPubSigKey
	the temp key from argon3
	*/
	TempPubSigKey [32]byte

	/*
	SigToTempPubSigKey
	use long-term id sig key to sign TempPubSigKey to promise it
	*/
	SigToTempPubSigKey  Signature
}

type Signature struct {
	SignerID ID // who sig it
	SigContent []byte // signature bytes
}