package mariadb

import (
	"fmt"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

type KV map[string]any

func (kv KV) Where(tx *gorm.DB) *gorm.DB {

	isDep, _ := regexp.Compile(`^\$\.`)

	for k, v := range kv {
		o := k
		if ok := isDep.MatchString(k); ok {
			k = fmt.Sprintf("JSON_EXTRACT(Doc,'%s')", k)
		}

		switch t := v.(type) {
		case int:
			tx = tx.Where(k+" = ?", k, t)
		case []string:
			tx = tx.Where(k+" IN ?", t)
		case string:
			if p := kv.EQ(t); p != nil {
				tx = tx.Where(k+" = ?", p)
			}
			if p := kv.LIKE(t); p != nil {
				tx = tx.Where(k+" LIKE ?", p)
			}
			if p := kv.LT(t); p != nil {
				tx = tx.Where(k+" < ?", p)
			}
			if p := kv.LTE(t); p != nil {
				tx = tx.Where(k+" <= ?", p)
			}
			if p := kv.GT(t); p != nil {
				tx = tx.Where(k+" > ?", p)
			}
			if p := kv.GTE(t); p != nil {
				tx = tx.Where(k+" >= ?", p)
			}
			if p := kv.IN(t); p != nil {
				tx = tx.Where(k+" IN ?", p)
			}
			if p := kv.NOT(t); p != nil {
				tx = tx.Where(k+" NOT ?", p)
			}
			if p := kv.SEARCH(t); p != nil {
				tx = tx.Where(fmt.Sprintf(`JSON_SEARCH(%s, 'one', ?) IS NOT NULL`, k), p)
			}
			if p := kv.LENGTH(t); p != nil {
				tx = tx.Where(fmt.Sprintf(`COALESCE(JSON_LENGTH(Doc, '%s'), 0) = ?`, o), p)
			}

		}
	}

	return tx
}

func (KV) EQ(str string) interface{} {
	if ok, _ := regexp.MatchString(":", str); !ok {
		return str
	}
	if ok, _ := regexp.MatchString("^EQ:", strings.ToUpper(str)); ok {
		return str[3:]
	}
	return nil
}
func (KV) LIKE(str string) interface{} {
	if ok, _ := regexp.MatchString("^LIKE:", strings.ToUpper(str)); ok {
		str = strings.ReplaceAll(str, "$", "%")
		return str[5:]
	}
	return nil
}

func (KV) LT(str string) interface{} {
	if ok, _ := regexp.MatchString("^LT:", strings.ToUpper(str)); ok {
		return str[3:]
	}
	return nil
}
func (KV) LTE(str string) interface{} {
	if ok, _ := regexp.MatchString("^LTE:", strings.ToUpper(str)); ok {
		return str[4:]
	}
	return nil
}
func (KV) GT(str string) interface{} {
	if ok, _ := regexp.MatchString("^GT:", strings.ToUpper(str)); ok {
		return str[3:]
	}
	return nil
}
func (KV) GTE(str string) interface{} {
	if ok, _ := regexp.MatchString("^GTE:", strings.ToUpper(str)); ok {
		return str[4:]
	}
	return nil
}
func (KV) IN(str string) interface{} {
	if ok, _ := regexp.MatchString("^IN:", strings.ToUpper(str)); ok {
		str = str[3:]
		return strings.Split(str, ",")
	}
	return nil
}
func (KV) NOT(str string) interface{} {
	if ok, _ := regexp.MatchString("^NOT:", strings.ToUpper(str)); ok {
		return str[4:]
	}
	return nil
}

func (KV) SEARCH(str string) interface{} {
	if ok, _ := regexp.MatchString("^SEARCH:", strings.ToUpper(str)); ok {
		return str[7:]
	}
	return nil
}

func (KV) LENGTH(str string) interface{} {
	if ok, _ := regexp.MatchString("^LENGTH:", strings.ToUpper(str)); ok {
		return str[7:]
	}
	return nil
}
