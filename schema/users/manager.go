package users

import (
	"errors"
	"fmt"

	"github.com/jigadhirasu/follow/redis"
	"github.com/jigadhirasu/follow/types"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

var mgr *Manager

func Mgr() *Manager {
	if mgr != nil {
		return mgr
	}

	mgr = &Manager{
		key: "cache:account",
	}
	return mgr
}

type Manager struct {
	key string
}

func (h *Manager) Length() int {
	dbr := redis.Connent()
	return int(dbr.HLen(dbr.Context(), h.key).Val())
}

func (h *Manager) Load(tx *gorm.DB, idd ...string) []User {
	tx = tx.Table(TableName)
	if len(idd) > 0 {
		tx = tx.Where(`UUID IN ?`, idd)
	}

	bb := [][]byte{}
	tx.Pluck(`JSON_MERGE_PATCH(Doc, JSON_OBJECT(
		'Updater',Updater,
		'UpdatedAt',UNIX_TIMESTAMP(UpdatedAt)
	)) as Doc`, &bb)

	xx := []User{}
	data := []any{}
	for _, b := range bb {
		x := User{}
		types.STRUCT(b, &x)
		xx = append(xx, x)

		data = append(data, x.UUID, b)
	}

	if len(bb) > 0 {
		dbr := redis.Connent()
		if err := dbr.HSet(dbr.Context(), h.key, data...).Err(); err != nil {
			panic(err)
		}
	}

	return xx
}

func (h *Manager) Get(id string) (User, error) {
	dbr := redis.Connent()
	b, err := dbr.HGet(dbr.Context(), h.key, id).Bytes()
	if err != nil || len(b) < 1 {
		return User{}, errors.New("404, account")
	}

	p := User{}
	types.STRUCT(b, &p)
	return p, nil
}

func (h *Manager) Exists(idd ...string) error {
	dbr := redis.Connent()

	keys := dbr.HKeys(dbr.Context(), h.key).Val()
	for _, id := range idd {
		if !slices.Contains(keys, id) {
			return fmt.Errorf("404, item %s not found", id)
		}
	}
	return nil
}

func (h *Manager) Delete(id string) error {
	dbr := redis.Connent()
	b, err := dbr.HGet(dbr.Context(), h.key, id).Bytes()
	if err != nil || len(b) < 1 {
		return fmt.Errorf("404, account")
	}

	if err := dbr.HDel(dbr.Context(), h.key, id).Err(); err != nil {
		return fmt.Errorf("501 %s", err)
	}
	return nil
}
