package mariadb

import (
	"time"

	"github.com/jigadhirasu/follow/types"
)

type Pack struct {
	Doc       types.Bytes `gorm:"column:Doc"`                                  // 資料主體
	Updater   string      `gorm:"column:Updater; type:varchar(40)"`            // 建立者
	UpdatedAt time.Time   `gorm:"column:UpdatedAt; default:CURRENT_TIMESTAMP"` // 建立時間
}
