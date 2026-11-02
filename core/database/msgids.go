package database

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

type MsgIDsDB struct {
	database *DB
}

type MsgIDsRecord struct {
	// should think carefully to avoid 合伙伪造，可靠的办法是依赖 channel 随机性，只认 channel 成员的 PoT
	Head      []byte // origin id msg Head
	Signature []byte // origin id msg Head Signature
}

func (db *MsgIDsDB) RecordMsgIDsInWindow(miw MsgIDsRecord) error {
	err := db.database.Set(MsgIDs, miw.Head, miw.Signature, 60*time.Minute)
	if err != nil {
		return fmt.Errorf("record message ID in window: %w", err)
	}

	return nil
}

func (db *MsgIDsDB) GetMsgIDsInWindow() (map[ID]MsgIDsRecord, error) {
	items := make(map[ID]MsgIDsRecord)
	err := db.database.ScanPrefix(MsgIDs, func(key, value []byte) error {
		valBytes, err := FromGobBytesToMsgIDsInWindow(value)
		if err != nil {
			return fmt.Errorf("decode message ID window record: %w", err)
		}
		if valBytes == nil {
			return nil
		}

		items[ID(key[:])] = *valBytes
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("list message IDs in window: %w", err)
	}

	return items, nil
}

func (miw *MsgIDsRecord) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	// 创建一个 gob 编码器
	enc := gob.NewEncoder(&buf)

	// 序列化结构体
	err := enc.Encode(miw)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrEncode, err)
	}

	return buf.Bytes(), nil
}

// FromGobBytesToMsgIDsInWindow 将 []byte 反序列化回 MsgIDsInWindow 结构体
func FromGobBytesToMsgIDsInWindow(data []byte) (*MsgIDsRecord, error) {
	var miw MsgIDsRecord
	buf := bytes.NewReader(data)
	// 创建一个 gob 解码器
	dec := gob.NewDecoder(buf)

	// 反序列化到结构体实例中
	err := dec.Decode(&miw)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDecode, err)
	}

	return &miw, nil
}
