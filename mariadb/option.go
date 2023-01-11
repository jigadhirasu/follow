package mariadb

import (
	"regexp"
	"strings"

	"github.com/jigadhirasu/follow/types"
	"gorm.io/gorm"
)

type TTDI func(tags types.Tags, tx *gorm.DB) (int64, error)
type TXDI func(tx *gorm.DB) (int64, error)

type Option func(tx *gorm.DB) *gorm.DB

func TableName(tablename string) Option {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Table(tablename)
	}
}

func Where(kv KV) Option {
	return func(tx *gorm.DB) *gorm.DB {
		return kv.Where(tx)
	}
}

func Limit(size, index int) Option {
	return func(tx *gorm.DB) *gorm.DB {
		if size < 1 {
			size = 25
		}
		return tx.Limit(size).Offset(index * size)
	}
}

func OrderBy(sorts ...string) Option {
	return func(tx *gorm.DB) *gorm.DB {
		regexpField, _ := regexp.Compile(`^\w+!?$`)

		ss := []string{}
		for _, s := range sorts {
			if !regexpField.MatchString(s) {
				continue
			}

			ss = append(ss, strings.ReplaceAll(s, "!", " DESC"))
		}

		return tx.Order(strings.Join(ss, ", "))
	}
}
