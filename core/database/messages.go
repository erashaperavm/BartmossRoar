package database

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

type MessagesDB struct {
	database *DB
}

type MessageRecord struct {
	Identity                    ID
	From                        ID
	Timestamp                   int64  // unix timestamp
	Content                     []byte // valid payload
	OriginCompressedArchiveData []byte // origin decrypted compressed datapack (Optional) data
	OriginCompressedArchivePath string // origin compressed archive path
}

func (db *MessagesDB) RecordHistoryMsg(msg MessageRecord, ttl int) error {
	var value MessageRecord
	if msg.OriginCompressedArchivePath != "" {
		// has path, prefer path in advance
		value = MessageRecord{
			From:                        msg.From,
			Timestamp:                   msg.Timestamp,
			Content:                     msg.Content,
			OriginCompressedArchivePath: msg.OriginCompressedArchivePath,
		}
	} else {
		if msg.OriginCompressedArchiveData != nil {
			value = MessageRecord{
				From:                        msg.From,
				Timestamp:                   msg.Timestamp,
				Content:                     msg.Content,
				OriginCompressedArchiveData: msg.OriginCompressedArchiveData,
			}
		} else {
			value = MessageRecord{
				From:      msg.From,
				Timestamp: msg.Timestamp,
				Content:   msg.Content,
			}
		}
	}

	bytesVal, err := value.Bytes()
	if err != nil {
		return fmt.Errorf("encode history message: %w", err)
	}

	if err := db.database.Set(Messages, msg.Identity[:], bytesVal, time.Duration(ttl)*time.Second); err != nil {
		return fmt.Errorf("store history message: %w", err)
	}
	return nil
}

func (db *MessagesDB) GetHistoryMsg(identity ID) (*MessageRecord, error) {
	bytesVal, err := db.database.Get(Messages, identity[:])
	if err != nil {
		return nil, fmt.Errorf("load history message: %w", err)
	}

	msg, err := FromGobBytesToMessage(bytesVal)
	if err != nil {
		return nil, fmt.Errorf("decode history message: %w", err)
	}

	return msg, nil
}

func (db *MessagesDB) DeleteHistoryMsg(identity ID) error {
	if err := db.database.Delete(Messages, identity[:]); err != nil {
		return fmt.Errorf("delete history message: %w", err)
	}
	return nil
}

func (db *MessagesDB) GetHistoryMsgIDList() ([][]byte, error) {
	var keyItems [][]byte
	err := db.database.ScanPrefix(Messages, func(key, value []byte) error {
		keyItems = append(keyItems, key)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("list history message IDs: %w", err)
	}

	return keyItems, nil
}

func (m *MessageRecord) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	// 创建一个 gob 编码器
	enc := gob.NewEncoder(&buf)

	// 序列化结构体
	err := enc.Encode(m)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrEncode, err)
	}

	return buf.Bytes(), nil
}

// FromGobBytesToMessage 将 []byte 反序列化回 MessageRecord 结构体
func FromGobBytesToMessage(data []byte) (*MessageRecord, error) {
	var m MessageRecord
	buf := bytes.NewReader(data)
	// 创建一个 gob 解码器
	dec := gob.NewDecoder(buf)

	// 反序列化到结构体实例中
	err := dec.Decode(&m)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDecode, err)
	}

	return &m, nil
}
