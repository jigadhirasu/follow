package users

import "github.com/jigadhirasu/follow/mariadb"

type Doc struct {
	UUID []byte `gorm:"column:UUID; uniqueindex:uuid; type:varchar(40) AS (JSON_VALUE(Doc, '$.UUID'))"`
	mariadb.Pack
	Name string `gorm:"column:Name; type:varchar(40) AS (JSON_VALUE(Doc, '$.Name'))"`
}

func (Doc) TableName() string {
	return TableName
}
