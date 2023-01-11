package users

import (
	"errors"

	"github.com/jigadhirasu/follow/mariadb"
	"github.com/jigadhirasu/follow/types"
	"github.com/jigadhirasu/follow/uuid"
	"gorm.io/gorm"
)

func Create(h *User) func(tags types.Tags, tx *gorm.DB) error {
	return func(tags types.Tags, tx *gorm.DB) error {
		h.UUID = uuid.New()
		if h.Name == "" {
			h.Name = "Guest"
		}
		if h.Subscribes == nil {
			h.Subscribes = []string{}
		}
		if h.Follows == nil {
			h.Follows = []string{}
		}

		pack := &mariadb.Pack{
			Doc:     types.JSON(h),
			Updater: tags.String("User"),
		}

		tx.Table(TableName).Create(pack)
		if tx.Error != nil {
			return tx.Error
		}

		Mgr().Load(tx, h.UUID)
		return nil
	}
}

func Update(h *User) func(tags types.Tags, tx *gorm.DB) error {
	return func(tags types.Tags, tx *gorm.DB) error {
		if len(h.UUID) == 0 {
			return errors.New("uuid invalid")
		}

		o, err := Mgr().Get(h.UUID)
		if err != nil {
			return err
		}

		_, next := o.Merge(*h)
		if result := tx.Table(TableName).Where("UUID = ?", h.UUID).Update("Doc", types.JSON(next)); result.Error != nil {
			return result.Error
		}

		Mgr().Load(tx, h.UUID)
		return nil
	}
}
