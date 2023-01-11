package follows

const TableName = "follows"

type Follow struct {
	Fanser   string `gorm:"column:Fanser; uniqueindex:ff"`
	Follower string `gorm:"column:Follower; uniqueindex:ff; index;"`
}

func (Follow) TableName() string {
	return TableName
}
