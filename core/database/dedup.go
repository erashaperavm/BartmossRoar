package database

import (
	"time"
)

type Dedup struct {
	DoneMsgHeadHash [32]byte
}

func (db *DB)RecordDoneMsgHeadHash(hash Dedup) error {
	err := db.Set(hash.DoneMsgHeadHash[:], []byte(""), 15 *time.Second)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB)IsDoneMsgHeadHashExist(hash Dedup) (bool, error) {
	return db.Has(hash.DoneMsgHeadHash[:])
}