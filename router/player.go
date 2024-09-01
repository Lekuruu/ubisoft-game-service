package router

import (
	"strconv"
	"strings"
)

type Info struct {
	Firstname string
	Surname   string
	Country   string
	Email     string
	Gender    uint8
	Public    bool
}

type Friends struct {
	Status  uint32
	Mood    uint32
	Ignored PlayerCollection
	List    PlayerCollection
}

type Player struct {
	Id      int
	Name    string
	Game    string
	Version string
	Info    Info
	Friends Friends
	Client
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
