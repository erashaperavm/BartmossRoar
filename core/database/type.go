package database

type TableKind = byte

const (
	Dedup TableKind = iota
	Messages
	MsgIDs
	Nodes
	PoT
)
