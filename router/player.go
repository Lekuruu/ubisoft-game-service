package router

import "strings"

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
}

func (player *Player) IpAddress() string {
	return strings.Split(player.Client.Conn.RemoteAddr().String(), ":")[0]
}
