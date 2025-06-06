package spacetimedb

import (
	"fmt"
	"log"
)

type CompressableQueryUpdate struct {
	Update *QueryUpdate
}

func (it *CompressableQueryUpdate) Deserialize(reader *BinaryReader) error {
	unionType := reader.ReadU8()
	switch unionType {
	case 0x00:
		uncompressed := &QueryUpdate{}
		if err := uncompressed.Deserialize(reader); err != nil {
			return fmt.Errorf("failed to deserialize uncompressed query update: %w", err)
		}
		it.Update = uncompressed
	case 0x01:
		log.Println("CompressableQueryUpdate.Deserialize: type 0x01 is not implemented yet. Brotli compression.")
	case 0x02:
		log.Println("CompressableQueryUpdate.Deserialize: type 0x02 is not implemented yet. Gzip compression.")
	}

	return nil
}
