package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jigadhirasu/follow/env"
	"github.com/jigadhirasu/follow/mariadb"
	"github.com/jigadhirasu/follow/schema/follows"
	"github.com/jigadhirasu/follow/schema/users"
	"github.com/jigadhirasu/follow/types"
)

var couter = make(chan int, 100)

//	func getTicket(ctx *gin.Context) {
//		<-couter
//	}
func relTicket(ctx *gin.Context) {
	go func() { couter <- 1 }()
}

func main() {

	env.Develop()

	go func() {
		for i := 0; i < 100; i++ {
			couter <- 1
		}
	}()
	// dbm := mariadb.Connect(os.Getenv("DBNAME"))
	// if err := dbm.AutoMigrate(new(users.Doc), new(follows.Follow)); err != nil {
	// 	panic(err)
	// }

	// users.Mgr().Load(dbm)

	router := gin.Default()
	router.Use(func(ctx *gin.Context) {
		select {
		case <-couter:
			ctx.Next()
		case <-time.After(time.Second * 2):
			ctx.JSON(504, "timeout")
			return
		}
	})

	router.POST("/register/:name", func(ctx *gin.Context) {
		db := mariadb.Connect(os.Getenv("DBNAME"))
		tags := types.Ctx2Tags(ctx)

		name := ctx.Param("name")
		u := users.User{Name: name}

		tx := db.Begin()
		defer tx.Rollback()
		if err := users.Create(&u)(tags, tx); err != nil {
			ctx.JSON(500, err)
			return
		}
		tx.Commit()

		ctx.SetCookie("user_cookie", u.UUID, 3600, "/", "localhost", false, true)
		ctx.JSON(200, "registed")
		ctx.Next()
	}, relTicket)

	router.GET("/subscribe/:uuid", func(ctx *gin.Context) {
		tags := types.Ctx2Tags(ctx)

		subscribe := ctx.Param("uuid")

		db := mariadb.Connect(os.Getenv("DBNAME"))
		tx := db.Begin()
		defer tx.Rollback()
		f := &follows.Follow{Fanser: tags.String("user"), Follower: subscribe}
		if result := tx.Create(f); result.Error != nil {
			ctx.JSON(500, result.Error)
			return
		}

		tx.Commit()
		ctx.JSON(200, "subscribed")
		ctx.Next()
	}, relTicket)

	router.GET("/unsubscribe/:uuid", func(ctx *gin.Context) {
		tags := types.Ctx2Tags(ctx)

		subscribe := ctx.Param("uuid")

		db := mariadb.Connect(os.Getenv("DBNAME"))
		tx := db.Begin()
		defer tx.Rollback()
		result := tx.Delete(new(follows.Follow), "Fanser = ? AND Follower = ?", tags.String("user"), subscribe)
		if result.Error != nil {
			ctx.JSON(500, result.Error)
			return
		}

		tx.Commit()
		ctx.JSON(200, "unsubscribed")
		ctx.Next()
	}, relTicket)

	router.GET("/follows", func(ctx *gin.Context) {
		tags := types.Ctx2Tags(ctx)

		size, _ := strconv.Atoi(ctx.Query("size"))
		page, _ := strconv.Atoi(ctx.Query("page"))

		db := mariadb.Connect(os.Getenv("DBNAME")).Debug()

		bb := [][]byte{}
		db.Table(users.TableName).Where(fmt.Sprintf(`
		UUID IN (SELECT Follower FROM %s WHERE Fanser = ? ORDER BY Follower)
		`, follows.TableName), tags.String("user")).Offset(page*size).Limit(size).Pluck(`
		JSON_ARRAYAGG(Doc)
		`, &bb)

		uu := []users.User{}
		if len(bb) < 1 {
			ctx.JSON(200, uu)
			return
		}
		types.STRUCT(bb[0], &uu)
		ctx.JSON(200, uu)
		ctx.Next()
	}, relTicket)

	router.GET("/fans", func(ctx *gin.Context) {
		tags := types.Ctx2Tags(ctx)

		size, _ := strconv.Atoi(ctx.Query("size"))
		page, _ := strconv.Atoi(ctx.Query("page"))

		db := mariadb.Connect(os.Getenv("DBNAME")).Debug()
		bb := [][]byte{}
		db.Table(users.TableName).Where(fmt.Sprintf(`
		UUID IN (SELECT Fanser FROM %s WHERE Follower = ? ORDER BY Fanser)
		`, follows.TableName), tags.String("user")).Offset(page*size).Limit(size).Pluck(`
		JSON_ARRAYAGG(Doc)
		`, &bb)

		uu := []users.User{}
		if len(bb) < 1 {
			ctx.JSON(200, uu)
			return
		}
		types.STRUCT(bb[0], &uu)
		ctx.JSON(200, uu)
		ctx.Next()
	}, relTicket)

	router.GET("/friends", func(ctx *gin.Context) {
		tags := types.Ctx2Tags(ctx)

		size, _ := strconv.Atoi(ctx.Query("size"))
		page, _ := strconv.Atoi(ctx.Query("page"))

		db := mariadb.Connect(os.Getenv("DBNAME")).Debug()

		ff := []string{}
		db.Table(follows.TableName).Where("Fanser = ?", tags.String("user")).Pluck(`Follower`, &ff)

		bb := [][]byte{}
		db.Table(users.TableName).Where(fmt.Sprintf(`
		UUID IN (
			SELECT Fanser FROM %s WHERE Follower = ? AND Fanser IN (?) ORDER BY Fanser
		)`, follows.TableName), tags.String("user"), ff).Offset(page*size).Limit(size).Pluck(`
		JSON_ARRAYAGG(Doc)
		`, &bb)

		uu := []users.User{}
		if len(bb) < 1 {
			ctx.JSON(200, uu)
			return
		}

		types.STRUCT(bb[0], &uu)
		ctx.JSON(200, uu)
		ctx.Next()
	}, relTicket)

	log.Println("service listen at :8080")
	router.Run(":8080")
}
