package mariadb

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var mgr = new(sync.Map)

func Connect(dbname string) *gorm.DB {
	dbname = os.Getenv("PREFIX_DBNAME") + "_" + dbname
	dbi, ok := mgr.Load(dbname)
	if ok {
		return dbi.(*gorm.DB)
	}

	GenerateDatabase(dbname)

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MARIADB_USER"),
		os.Getenv("MARIADB_PASS"),
		os.Getenv("MARIADB_HOST"),
		dbname,
	)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         255,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	maxidle, _ := strconv.Atoi(os.Getenv("MARIADB_MAX_IDLE"))
	if maxidle < 1 {
		maxidle = 4
	}
	maxopen, _ := strconv.Atoi(os.Getenv("MARIADB_MAX_OPEN"))
	if maxopen < 1 {
		maxopen = 16
	}

	sqlDB.SetMaxIdleConns(maxidle)
	sqlDB.SetMaxOpenConns(maxopen)
	sqlDB.SetConnMaxLifetime(time.Hour)

	mgr.Store(dbname, db)
	return db
}

// GenerateDatabase 初始化資料庫
func GenerateDatabase(dbnames ...string) {

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MARIADB_USER"),
		os.Getenv("MARIADB_PASS"),
		os.Getenv("MARIADB_HOST"),
	)

	db, err := gorm.Open(mysql.New(mysql.Config{DSN: dsn}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	result := []string{}
	db.Raw(`SHOW DATABASES`).Scan(&result)

	set := map[string]bool{}
	// 指定要新建的DB
	for _, dbname := range dbnames {
		set[dbname] = true
	}
	// 移除已經存在的DB
	for _, dbname := range result {
		set[dbname] = false
	}
	// 逐一建立新DB
	tx := db.Begin()
	defer tx.Rollback()
	for dbname, gen := range set {
		if gen {
			tx.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci", dbname))
			if tx.Error != nil {
				panic(tx.Error)
			}
		}
	}
	tx.Commit()
}
