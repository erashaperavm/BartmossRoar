package database

import (
	"fmt"
	"time"
)

type DedupDB struct {
	database *DB
}

type DedupRecord struct {
	DoneMsgHeadHash [32]byte
}

func (db *DedupDB) RecordDoneMsgHeadHash(hash DedupRecord) error {
	err := db.database.Set(Dedup, hash.DoneMsgHeadHash[:], []byte(""), 15*time.Second)
	if err != nil {
		return fmt.Errorf("record done message head hash: %w", err)
	}

	return nil
}

func (db *DedupDB) IsDoneMsgHeadHashExist(hash DedupRecord) (bool, error) {
	exists, err := db.database.Has(Dedup, hash.DoneMsgHeadHash[:])
	if err != nil {
		return false, fmt.Errorf("check done message head hash: %w", err)
	}
	return exists, nil
}
