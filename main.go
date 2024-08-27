package main

import (
	"github.com/lekuruu/ubisoft-game-service/common"
	"github.com/lekuruu/ubisoft-game-service/router"
)

func main() {
	router := router.Router{
		Host:   "127.0.0.1",
		Port:   40000,
		Logger: *common.CreateLogger("Router", common.DEBUG),
	}

	router.Serve()
}
