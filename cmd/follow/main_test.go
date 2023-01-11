package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jigadhirasu/follow/env"
	"github.com/jigadhirasu/follow/mariadb"
	"github.com/jigadhirasu/follow/schema/users"
	"golang.org/x/exp/slices"
)

const host = "http://localhost:8080"

var chs = make(chan time.Duration, 1000)
var ns = time.Duration(0)
var times = int64(0)
var failed = 0
var wg = new(sync.WaitGroup)

func Subscribe(self, subsribe string) {
	start := time.Now()
	resp, err := http.Get(fmt.Sprintf("%s/subscribe/%s?id=%s", host, subsribe, self))
	if err != nil {
		panic(err)
	}
	if resp.StatusCode == http.StatusOK {
		chs <- time.Since(start)
		return
	}
	failed++
}
func Unsubscribe(self, unsubscribe string) {
	start := time.Now()
	resp, err := http.Get(fmt.Sprintf("%s/unsubscribe/%s?id=%s", host, unsubscribe, self))
	if err != nil {
		panic(err)
	}
	if resp.StatusCode == http.StatusOK {
		chs <- time.Since(start)
		return
	}
	failed++
}
func Fans(self string, size, page int) {
	start := time.Now()
	resp, err := http.Get(fmt.Sprintf("%s/fans?id=%s&size=%d&page=%d", host, self, size, page))
	if err != nil {
		panic(err)
	}
	if resp.StatusCode == http.StatusOK {
		chs <- time.Since(start)
		return
	}
	failed++
}
func Follows(self string, size, page int) {
	start := time.Now()
	resp, err := http.Get(fmt.Sprintf("%s/follows?id=%s&size=%d&page=%d", host, self, size, page))
	if err != nil {
		panic(err)
	}
	if resp.StatusCode == http.StatusOK {
		chs <- time.Since(start)
		return
	}
	failed++
}
func Friends(self string, size, page int) {
	start := time.Now()
	resp, err := http.Get(fmt.Sprintf("%s/friends?id=%s&size=%d&page=%d", host, self, size, page))
	if err != nil {
		panic(err)
	}
	if resp.StatusCode == http.StatusOK {
		chs <- time.Since(start)
		return
	}
	failed++
}

type ModidyRequest struct {
	Method string
	Self   string
	Target string
}
type QueryRequest struct {
	Method string
	Self   string
	Size   int
	Page   int
}

func TestMain(t *testing.T) {
	env.Develop()

	go func() {
		for {
			v, ok := <-chs
			if !ok {
				return
			}
			ns += v
			times++
			wg.Done()
		}
	}()

	rand.Seed(time.Now().Unix())

	db := mariadb.Connect(os.Getenv("DBNAME"))

	idd := []string{}
	db.Table(users.TableName).Pluck(`UUID`, &idd)
	count := len(idd)

	go main()

	<-time.After(time.Second)

	modifyCH := make(chan ModidyRequest)
	for i := 0; i < 10; i++ {
		go func() {
			for {
				m, ok := <-modifyCH
				if !ok {
					panic("ggg")
				}
				switch m.Method {
				case "unsubscribe":
					Unsubscribe(m.Self, m.Target)
				default:
					Subscribe(m.Self, m.Target)
				}
			}
		}()
	}

	queryCH := make(chan QueryRequest)
	for i := 0; i < 10; i++ {
		go func() {
			for {
				q, ok := <-queryCH
				if !ok {
					panic("ggg")
				}
				switch q.Method {
				case "friends":
					Friends(q.Self, q.Size, q.Page)
				case "follows":
					Follows(q.Self, q.Size, q.Page)
				default:
					Fans(q.Self, q.Size, q.Page)
				}
			}
		}()
	}

	x := 10000
	wg.Add(x)
	for i := 0; i < x; i++ {
		rd := rand.Intn(count)
		self := idd[rd]
		switch rd % 5 {
		case 0:
			tg := idd[rand.Intn(count)]
			modifyCH <- ModidyRequest{Method: "subscribe", Self: self, Target: tg}
		case 1:
			tg := idd[rand.Intn(count)]
			modifyCH <- ModidyRequest{Method: "unsubscribe", Self: self, Target: tg}
		case 2:
			size := rand.Intn(5) + 1
			page := rand.Intn(3)
			queryCH <- QueryRequest{Method: "fans", Self: self, Size: size, Page: page}
		case 3:
			size := rand.Intn(5) + 1
			page := rand.Intn(3)
			queryCH <- QueryRequest{Method: "follows", Self: self, Size: size, Page: page}
		default:
			size := rand.Intn(5) + 1
			page := rand.Intn(3)
			queryCH <- QueryRequest{Method: "friends", Self: self, Size: size, Page: page}
		}
	}

	wg.Wait()

	fmt.Println(ns, times, ns/time.Duration(times), failed)

}

// func TestData(t *testing.T) {
// 	env.Develop()

// 	rand.Seed(time.Now().Unix())

// 	// dbr := redis.Connent()
// 	// dbr.Set(dbr.Context(), "hhh", "hhh", 0)
// 	// if err := mariadb.Connect(os.Getenv("DBNAME")).AutoMigrate(new(users.Doc), new(follows.Follow)); err != nil {
// 	// 	panic(err)
// 	// }

// 	db := mariadb.Connect(os.Getenv("DBNAME"))

// 	idd := []string{}
// 	db.Table(users.TableName).Pluck(`UUID`, &idd)

// 	count := len(idd)
// 	// for i := 0; i < count; i++ {
// 	// 	idd = append(idd, uuid.New())
// 	// }

// 	// pp := []*mariadb.Pack{}
// 	ff := []*follows.Follow{}
// 	for _, id := range idd[1000:10000] {
// 		// pp = append(pp, &mariadb.Pack{
// 		// 	Doc: types.JSON(users.User{
// 		// 		UUID:       id,
// 		// 		Name:       fmt.Sprintf("user%05d", i+1),
// 		// 		Subscribes: []string{},
// 		// 		Follows:    []string{},
// 		// 	}),
// 		// 	Updater: "default",
// 		// })

// 		start := rand.Intn(count - 10)
// 		for _, fl := range idd[start : start+10] {
// 			if fl != id {
// 				ff = append(ff, &follows.Follow{Fanser: id, Follower: fl})
// 			}
// 		}
// 	}
// 	// db.Table(users.TableName).CreateInBatches(pp, 10000)
// 	if result := db.CreateInBatches(ff, 10000); result.Error != nil {
// 		panic(result.Error)
// 	}

// }

func TestSlice(t *testing.T) {

	ss := []string{"A", "B", "C", "D", "C", "B", "D"}

	for {
		idx := slices.Index(ss, "D")
		if idx < 0 {
			break
		}
		ss = slices.Delete(ss, idx, idx+1)
	}

	fmt.Println(ss, len(ss))
}

func TestChan(t *testing.T) {

	ch := make(chan int, 10)
	for i := 0; i < 10; i++ {
		ch <- i
	}

	go func() {
		for i := 0; i < 10; i++ {
			g := i + 10
			ch <- g
			fmt.Println(g)
		}
	}()

	for {
		v, ok := <-ch
		fmt.Println(v, ok)
		if !ok {
			break
		}
		<-time.After(time.Second)
	}

	close(ch)

}
