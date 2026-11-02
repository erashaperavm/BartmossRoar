package database

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/erashaperavm/BartmossRoar/core/crypto"
)

type PoTDB struct {
	database *DB
}

type PoTRecord struct {
	Id        ID
	Signature crypto.Signature
	MsgRecord MsgIDsRecord
	GroupID   []ID
}

func (db *PoTDB) RecordPoT(pot PoTRecord) error {
	b, err := pot.Bytes()
	if err != nil {
		return fmt.Errorf("encode PoT record: %w", err)
	}

	if err := db.database.Set(PoT, pot.Id[:], b, 60*time.Minute); err != nil {
		return fmt.Errorf("store PoT record: %w", err)
	}

	return nil
}

func (db *PoTDB) GetPoTs() (map[ID]PoTRecord, error) {
	items := make(map[ID]PoTRecord)
	err := db.database.ScanPrefix(PoT, func(key, value []byte) error {
		valBytes, err := FromGobBytesToPoTRecord(value)
		if err != nil {
			return fmt.Errorf("decode PoT record: %w", err)
		}
		if valBytes == nil {
			return nil
		}

		items[ID(key[:])] = *valBytes
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("list PoT records: %w", err)
	}

	return items, nil
}

func (p *PoTRecord) Bytes() ([]byte, error) {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrEncode, err)
	}

	return buf.Bytes(), nil
}

func FromGobBytesToPoTRecord(data []byte) (*PoTRecord, error) {
	var pot PoTRecord
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&pot)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDecode, err)
	}

	return &pot, nil
}
