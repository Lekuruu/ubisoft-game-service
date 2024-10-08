package main

import (
	"flag"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/lekuruu/ubisoft-game-service/cdkey"
	"github.com/lekuruu/ubisoft-game-service/common"
	"github.com/lekuruu/ubisoft-game-service/gsconnect"
	"github.com/lekuruu/ubisoft-game-service/gsnat"
	"github.com/lekuruu/ubisoft-game-service/irc"
	"github.com/lekuruu/ubisoft-game-service/proxy"
	"github.com/lekuruu/ubisoft-game-service/router"
)

type Config struct {
	Web struct {
		Host string
		Port int
	}
	Router struct {
		Host string
		Port int
	}
	Proxy struct {
		Host string
		Port int
	}
	IRC struct {
		Host string
		Port int
	}
	NAT struct {
		Port int
	}
	CDKey struct {
		Port int
	}
	Games        []string
	ExternalHost string
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
		fmt.Sprintf("IRCPort0=%d", c.IRC.Port),
		fmt.Sprintf("ProxyIP0=%s", c.ExternalHost),
		fmt.Sprintf("ProxyPort0=%d", c.Proxy.Port),
	}

	games := make(map[string]string)

	for _, game := range c.Games {
		games[game] = strings.Join(config, "\n")
	}

	return games
}

func loadConfig() (*Config, error) {
	var config Config

	flag.StringVar(&config.Web.Host, "web-host", "0.0.0.0", "Web server host")
	flag.IntVar(&config.Web.Port, "web-port", 80, "Web server port")

	flag.StringVar(&config.Router.Host, "router-host", "0.0.0.0", "Router server host")
	flag.IntVar(&config.Router.Port, "router-port", 40000, "Router server port")

	flag.StringVar(&config.Proxy.Host, "proxy-host", "0.0.0.0", "Proxy server host")
	flag.IntVar(&config.Proxy.Port, "proxy-port", 4040, "Proxy server port")

	flag.StringVar(&config.IRC.Host, "irc-host", "0.0.0.0", "IRC server host")
	flag.IntVar(&config.IRC.Port, "irc-port", 6668, "IRC server port")

	flag.IntVar(&config.NAT.Port, "nat-port", 45000, "NAT server port")
	flag.IntVar(&config.CDKey.Port, "cdkey-port", 44000, "CDKey server port")

	flag.StringVar(&config.ExternalHost, "external-host", "127.0.0.1", "External host address")
	flag.Parse()

	// Default games list
	config.Games = []string{
		"SPLINTERCELL3PCADVERS",
		"SPLINTERCELL3PCCOOP",
		"SPLINTERCELL3PS2US",
		"SPLINTERCELL3PC",
		"HEROES_5",
	}

	// TODO: Move supported games into database
	games := flag.Args()

	if len(games) != 0 {
		// Overwrite default games list, if
		// games are provided as arguments
		config.Games = games
	}

	sort.Strings(config.Games)
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
		Port:   uint16(config.Router.Port),
		Logger: *common.CreateLogger("Router", common.DEBUG),
		Games:  config.Games,
	}

	proxy := proxy.Proxy{
		Host:   config.Proxy.Host,
		Port:   uint16(config.Proxy.Port),
		Logger: *common.CreateLogger("Proxy", common.DEBUG),
	}

	irc := irc.IRCServer{
		Host:   config.IRC.Host,
		Port:   uint16(config.IRC.Port),
		Logger: *common.CreateLogger("IRC", common.DEBUG),
	}

	cdks := cdkey.CDKeyServer{
		Port:   uint16(config.CDKey.Port),
		Logger: *common.CreateLogger("CDKey", common.DEBUG),
	}

	nat := gsnat.GSNatServer{
		Port:   uint16(config.NAT.Port),
		Logger: *common.CreateLogger("GSNat", common.DEBUG),
	}

	var wg sync.WaitGroup

	runService(&wg, router.Serve)
	runService(&wg, proxy.Serve)
	runService(&wg, cdks.Serve)
	runService(&wg, irc.Serve)
	runService(&wg, nat.Serve)
	runService(&wg, gsc.Serve)

	wg.Wait()
}
