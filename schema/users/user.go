package users

const TableName = "users"

type User struct {
	UUID       string
	Name       string
	Subscribes []string // 追蹤者
	Follows    []string // 訂閱者
	Updater    string   `json:",omitempty"`
	UpdatedAt  int64    `json:",omitempty"`
}

func (User) TableName() string {
	return TableName
}

// return modify before and after
func (h User) Merge(n User) (User, User) {
	o := h
	if n.Name != "" {
		h.Name = n.Name
	}
	if len(n.Subscribes) > 0 {
		h.Subscribes = n.Subscribes
	}
	return o, h
}
