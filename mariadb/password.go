package mariadb

import (
	"crypto/sha256"
	"fmt"
)

func Password(password string) string {
	s := fmt.Sprintf("初一吃素%s初二出素%s初三粗促QQ!~", password[:4], password[4:])
	hash := sha256.New()
	hash.Write([]byte(s))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
