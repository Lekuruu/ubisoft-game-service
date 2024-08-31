package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/lekuruu/ubisoft-game-service/cdkey"
	"github.com/lekuruu/ubisoft-game-service/common"
	"github.com/lekuruu/ubisoft-game-service/gsconnect"
	"github.com/lekuruu/ubisoft-game-service/gsnat"
	"github.com/lekuruu/ubisoft-game-service/router"
	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Web struct {
		Host string `toml:"Host"`
		Port int    `toml:"Port"`
	} `toml:"Web"`
	Router struct {
		Host string `toml:"Host"`
		Port uint16 `toml:"Port"`
	} `toml:"Router"`
	NAT struct {
		Port uint16 `toml:"Port"`
	} `toml:"NAT"`
	CDKey struct {
		Port uint16 `toml:"Port"`
	} `toml:"CDKey"`
	Games        []string `toml:"Games"`
	ExternalHost string   `toml:"ExternalHost"`
}

func (c *Config) createGameConfig() map[string]string {
	config := []string{
		"[Servers]",
		fmt.Sprintf("RouterIP0=%s", c.ExternalHost),
		fmt.Sprintf("RouterPort0=%d", c.Router.Port),
		fmt.Sprintf("CDKeyServerIP0=%s", c.ExternalHost),
		fmt.Sprintf("CDKeyServerPort0=%d", c.CDKey.Port),
		fmt.Sprintf("NATServerIP0=%s", c.ExternalHost),
		fmt.Sprintf("NATServerPort0=%d", c.NAT.Port),
		fmt.Sprintf("IRCIP0=%s", c.ExternalHost),
		fmt.Sprintf("IRCPort0=%d", 6668),
		fmt.Sprintf("ProxyIP0=%s", c.ExternalHost),
		fmt.Sprintf("ProxyPort0=%d", 4040),
	}

	games := make(map[string]string)

	for _, game := range c.Games {
		games[game] = strings.Join(config, "\n")
	}

	return games
}

func loadConfig() (*Config, error) {
	file, err := os.ReadFile("config.toml")
	if err != nil {
		return nil, err
	}

	var config Config

	err = toml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func runService(wg *sync.WaitGroup, worker func()) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		worker()
	}()
}

func main() {
	config, err := loadConfig()

	if err != nil {
		fmt.Println("Failed to load 'config.toml' file:", err)
		return
	}

	gsc := gsconnect.GSConnect{
		Host:   config.Web.Host,
		Port:   config.Web.Port,
		Games:  config.createGameConfig(),
		Logger: *common.CreateLogger("GSConnect", common.DEBUG),
	}

	router := router.Router{
		Host:   config.Router.Host,
		Port:   config.Router.Port,
		Logger: *common.CreateLogger("Router", common.DEBUG),
	}

	cdks := cdkey.CDKeyServer{
		Port:   config.CDKey.Port,
		Logger: *common.CreateLogger("CDKeyServer", common.DEBUG),
	}

	nat := gsnat.GSNatServer{
		Port:   config.NAT.Port,
		Logger: *common.CreateLogger("GSNatServer", common.DEBUG),
	}

	var wg sync.WaitGroup

	runService(&wg, router.Serve)
	runService(&wg, cdks.Serve)
	runService(&wg, nat.Serve)
	runService(&wg, gsc.Serve)

	wg.Wait()
}
