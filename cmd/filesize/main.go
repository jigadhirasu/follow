package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

type Manager struct {
	wg   *sync.WaitGroup
	data *sync.Map
}

func main() {

	mgr := &Manager{
		wg:   new(sync.WaitGroup),
		data: new(sync.Map),
	}

	mgr.wg.Add(1)
	go traverseDir(mgr, os.Getenv("SEARCH_DIR"))

	mgr.wg.Wait()

	mgr.data.Range(func(key, value any) bool {
		fmt.Printf("%s: %d \n", key, value)
		return true
	})

}

func traverseDir(mgr *Manager, dir string) {
	ff, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	mgr.wg.Add(len(ff))
	for _, f := range ff {
		name := f.Name()

		if f.IsDir() {
			go traverseDir(mgr, dir+"/"+name)
			continue
		}

		fi, _ := f.Info()
		idx := strings.LastIndex(fi.Name(), ".")
		sub := "file"
		if idx > 0 {
			sub = fi.Name()[idx:]
		}
		data, ok := mgr.data.Load(sub)
		size := fi.Size()
		if ok {
			size += data.(int64)
		}
		// fmt.Println(sub, size, fi.Name())
		mgr.data.Store(sub, size)
		mgr.wg.Done()
	}

	mgr.wg.Done()
}
