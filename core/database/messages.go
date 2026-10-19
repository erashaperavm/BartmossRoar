package database

import (
	"bytes"
	"encoding/gob"
	"time"
)

type Message struct {
	Identity                    ID
	From                        ID
	Timestamp                   int64  // unix timestamp
	Content                     []byte // valid payload
	OriginCompressedArchiveData []byte // origin decrypted compressed datapack (Optional) data
	OriginCompressedArchivePath string // origin compressed archive path
}

func (db *DB) RecordHistoryMsg(msg Message, ttl int) error {
	key := append([]byte("messages::"), msg.Identity[:]...)
	var value Message
	if msg.OriginCompressedArchivePath != "" {
		// has path, prefer path in advance
		value = Message{
			From:                        msg.From,
			Timestamp:                   msg.Timestamp,
			Content:                     msg.Content,
			OriginCompressedArchivePath: msg.OriginCompressedArchivePath,
		}
	} else {
		if msg.OriginCompressedArchiveData != nil {
			value = Message{
				From:                        msg.From,
				Timestamp:                   msg.Timestamp,
				Content:                     msg.Content,
				OriginCompressedArchiveData: msg.OriginCompressedArchiveData,
			}
		} else {
			value = Message{
				From:      msg.From,
				Timestamp: msg.Timestamp,
				Content:   msg.Content,
			}
		}
	}

	bytesVal, err := value.Bytes()
	if err != nil {
		return err
	}

	return db.Set(key, bytesVal, time.Duration(ttl)*time.Second)
}

func (db *DB) GetHistoryMsg(identity ID) (*Message, error) {
	key := append([]byte("messages::"), identity[:]...)
	bytesVal, err := db.Get(key)
	if err != nil {
		return nil, err
	}

	msg, err := FromGobBytesToMessage(bytesVal)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (db *DB) DeleteHistoryMsg(identity ID) error {
	key := append([]byte("messages::"), identity[:]...)
	return db.Delete(key)
}

func (db *DB) GetHistoryMsgList() [][]byte {
	var keyItems [][]byte
	db.ScanPrefix([]byte("messages::"), func(key, value []byte) error {
		keyItems = append(keyItems, key)
		return nil
	})

	return keyItems
}

func (m *Message) Bytes() ([]byte, error) {
	var buf bytes.Buffer
	// 创建一个 gob 编码器
	enc := gob.NewEncoder(&buf)

	// 序列化结构体
	err := enc.Encode(m)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// FromGobBytes 将 []byte 反序列化回 Message 结构体
func FromGobBytesToMessage(data []byte) (*Message, error) {
	var m Message
	buf := bytes.NewReader(data)
	// 创建一个 gob 解码器
	dec := gob.NewDecoder(buf)

	// 反序列化到结构体实例中
	err := dec.Decode(&m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}
