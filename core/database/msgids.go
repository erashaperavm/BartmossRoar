package database

import (
	"bytes"
	"encoding/gob"
)

type MsgIDsInWindow struct {
	// should think carefully to avoid 合伙伪造，可靠的办法是依赖 channel 随机性，只认 channel 成员的 PoT
	head      []byte // origin id msg head
	signature []byte // origin id msg head signature
}

func (db *DB) RecordMsgIDsInWindow(identity ID, miw MsgIDsInWindow, ttl int) error {}

func (db *DB) GetMsgIDsInWindow() (map[ID]MsgIDsInWindow, error) {
	items := make(map[ID]MsgIDsInWindow)
	prefix := []byte("msgids::")
	err := db.ScanPrefix(prefix, func(key, value []byte) error {
		valBytes, err := FromGobBytesToMsgIDsInWindow(value)
		if err != nil {
			return err
		}
		if valBytes == nil {
			return nil
		}

		items[ID(key[:])] = *valBytes
		return nil
	})

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (miw *MsgIDsInWindow) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	// 创建一个 gob 编码器
	enc := gob.NewEncoder(&buf)

	// 序列化结构体
	err := enc.Encode(miw)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// FromGobBytesToMsgIDsInWindow 将 []byte 反序列化回 MsgIDsInWindow 结构体
func FromGobBytesToMsgIDsInWindow(data []byte) (*MsgIDsInWindow, error) {
	var miw MsgIDsInWindow
	buf := bytes.NewReader(data)
	// 创建一个 gob 解码器
	dec := gob.NewDecoder(buf)

	// 反序列化到结构体实例中
	err := dec.Decode(&miw)
	if err != nil {
		return nil, err
	}

	return &miw, nil
}
