package mariadb

import (
	"bytes"

	"github.com/jigadhirasu/follow/types"
	"gorm.io/gorm"
)

func Count(opts ...Option) func(tx *gorm.DB) (int64, error) {
	return func(tx *gorm.DB) (int64, error) {
		for _, opt := range opts {
			tx = opt(tx)
		}
		var count int64
		tx = tx.Count(&count)
		return count, tx.Error
	}
}

func Find(opts ...Option) func(tx *gorm.DB) (types.Bytes, error) {
	return func(tx *gorm.DB) (types.Bytes, error) {
		for _, opt := range opts {
			tx = opt(tx)
		}

		bb := [][]byte{}
		tx = tx.Pluck(`JSON_MERGE_PATCH(Doc, JSON_OBJECT('CreatedAt',CreatedAt,'UpdatedAt',UNIX_TIMESTAMP(UpdatedAt))) as Doc`, &bb)
		if tx.Error != nil {
			return nil, tx.Error
		}
		out := bytes.Join(bb, []byte(","))
		out = append([]byte("["), out...)
		out = append(out, []byte("]")...)
		return out, nil
	}
}
