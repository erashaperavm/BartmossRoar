package database

type Message struct {
	identity                ID
	from                    ID
	timestamp               int64  // unix timestamp
	content                 []byte // valid payload
	originCompressedArchive []byte // origin decrypted compressed datapack (Optional)
}
