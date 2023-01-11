package env

import "os"

func Develop() {
	os.Setenv("MARIADB_HOST", "34.80.181.3:3306")
	os.Setenv("MARIADB_USER", "4dingdev")
	os.Setenv("MARIADB_PASS", "4ding@dev")
	os.Setenv("PREFIX_DBNAME", "f")
	os.Setenv("DBNAME", "focus")
	os.Setenv("REDIS_ADDR", "34.80.181.3:6379")
	os.Setenv("REDIS_DB", "14")
	os.Setenv("HTTP_PORT", "8080")
}
