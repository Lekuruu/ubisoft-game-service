package main

import (
	"sync"

	"github.com/lekuruu/ubisoft-game-service/cdkey"
	"github.com/lekuruu/ubisoft-game-service/common"
	"github.com/lekuruu/ubisoft-game-service/router"
)

func runService(wg *sync.WaitGroup, worker func()) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		worker()
	}()
}

func main() {
	router := router.Router{
		Host:   "127.0.0.1",
		Port:   40000,
		Logger: *common.CreateLogger("Router", common.DEBUG),
	}

	cdks := cdkey.CDKeyServer{
		Port:   44000,
		Logger: *common.CreateLogger("CDKeyServer", common.DEBUG),
	}

	var wg sync.WaitGroup

	runService(&wg, router.Serve)
	runService(&wg, cdks.Serve)

	wg.Wait()
}
