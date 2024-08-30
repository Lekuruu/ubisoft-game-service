package main

import (
	"strings"
	"sync"

	"github.com/lekuruu/ubisoft-game-service/cdkey"
	"github.com/lekuruu/ubisoft-game-service/common"
	"github.com/lekuruu/ubisoft-game-service/gsconnect"
	"github.com/lekuruu/ubisoft-game-service/gsnat"
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

	gsc := gsconnect.GSConnect{
		Host:   "127.0.0.1",
		Port:   80,
		Games:  make(map[string]string),
		Logger: *common.CreateLogger("GSConnect", common.DEBUG),
	}

	nat := gsnat.GSNatServer{
		Port:   7781,
		Logger: *common.CreateLogger("GSNatServer", common.DEBUG),
	}

	scct := []string{
		"[Servers]",
		"RouterIP0=127.0.0.1",
		"RouterPort0=40000",
		"IRCIP0=127.0.0.1",
		"IRCPort0=6668",
		"CDKeyServerIP0=127.0.0.1",
		"CDKeyServerPort0=44000",
		"ProxyIP0=127.0.0.1",
		"ProxyPort0=4040",
		"NATServerIP0=127.0.0.1",
		"NATServerPort0=7781",
	}

	gsc.Games["SPLINTERCELL3PCADVERS"] = strings.Join(scct, "\n")
	gsc.Games["SPLINTERCELL3PCCOOP"] = strings.Join(scct, "\n")
	gsc.Games["SPLINTERCELL3PS2US"] = strings.Join(scct, "\n")
	gsc.Games["SPLINTERCELL3PC"] = strings.Join(scct, "\n")
	gsc.Games["HEROES_5"] = strings.Join(scct, "\n")

	var wg sync.WaitGroup

	runService(&wg, router.Serve)
	runService(&wg, cdks.Serve)
	runService(&wg, nat.Serve)
	runService(&wg, gsc.Serve)

	wg.Wait()
}

// TODO: Add configuration for services
