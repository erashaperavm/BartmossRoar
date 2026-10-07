package database

type MsgIDsInWindow struct {
	// should think carefully to avoid 合伙伪造，可靠的办法是依赖 channel 随机性，只认 channel 成员的 PoT
	identity  ID
	head      []byte // origin id msg head
	signature []byte // origin id msg head signature
}
