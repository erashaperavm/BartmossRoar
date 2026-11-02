package database

import (
	"errors"
	"fmt"
)

var (
	ErrNotOpen = errors.New("database: not open")
	ErrEncode  = errors.New("database: encode record")
	ErrDecode  = errors.New("database: decode record")
)

// ErrUnknownTableKind means a key prefix is not one of this package's logical tables.
type ErrUnknownTableKind struct {
	kind TableKind
}

func (e *ErrUnknownTableKind) Error() string {
	// TableKind is a byte prefix; print it numerically so non-printable values are clear.
	return fmt.Sprintf("database: unknown table kind %d (0x%02x)", e.kind, e.kind)
}
