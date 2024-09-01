package router

import (
	"strconv"
	"strings"
)

type Player struct {
	Client
	Id        int
	Name      string
	Firstname string
	Surname   string
	Country   string
	Email     string
	Game      string
	Version   string
	Public    bool
	Status    uint32
	Mood      uint32
}

func (player *Player) IpAddress() string {
	return strings.Split(player.Client.Conn.RemoteAddr().String(), ":")[0]
}

func (player *Player) Port() int {
	portString := strings.Split(player.Client.Conn.RemoteAddr().String(), ":")[1]
	port, err := strconv.Atoi(portString)

	if err != nil {
		return 0
	}

	return port
}
