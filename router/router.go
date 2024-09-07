package router

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/lekuruu/ubisoft-game-service/common"
)

type Router struct {
	Host    string
	Port    uint16
	Games   []string
	Logger  common.Logger
	Players PlayerCollection
	Pending map[string]*Player
}

type Client struct {
	Conn   net.Conn
	Server *Router
	Player *Player
	State  *common.GSClientState
}

func (router *Router) Serve() {
	router.Players = NewPlayerCollection()
	router.Pending = make(map[string]*Player)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", router.Host, router.Port))

	if err != nil {
		log.Fatal(err)
	}

	router.Logger.Info(fmt.Sprintf("Listening on %s:%d", router.Host, router.Port))

	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go router.HandleClient(conn)
	}
}

func (router *Router) HandleClient(conn net.Conn) {
	router.Logger.Info(fmt.Sprintf("-> <%s>", conn.RemoteAddr()))

	client := &Client{
		Conn:   conn,
		Server: router,
		State:  &common.GSClientState{},
	}

	defer router.OnDisconnect(client)

	for {
		msg, err := common.ReadGSMessage(client.Conn, client.State)

		if err == io.EOF {
			// Client disconnected
			break
		}

		if err != nil {
			router.Logger.Error(fmt.Sprintf("Failed to parse header: %s", err))
			break
		}

		router.Logger.Debug(fmt.Sprintf("-> %v", msg.String()))
		handler, ok := RouterHandlers[msg.Type]

		if !ok {
			router.Logger.Warning(fmt.Sprintf("Couldn't find handler for type '%d'", msg.Type))
			continue
		}

		response, gsError := handler(msg, client)

		if gsError != nil {
			router.Logger.Error(gsError.Error())
			response = gsError.Response(msg)
		}

		if response == nil {
			// No response & no error
			continue
		}

		serialized, err := response.Serialize(client.State)

		if err != nil {
			router.Logger.Error(fmt.Sprintf("Failed to serialize message: %s", err))
			break
		}

		_, err = conn.Write(serialized)

		if err != nil {
			router.Logger.Error(fmt.Sprintf("Failed to send message: %s", err))
			break
		}

		router.Logger.Debug(fmt.Sprintf("<- %v", response.String()))
	}
}

func (router *Router) OnDisconnect(client *Client) {
	if r := recover(); r != nil {
		router.Logger.Error(fmt.Sprintf("Panic: %s", r))
	}

	if client.Player != nil {
		router.Players.Remove(client.Player)
	}

	router.Logger.Info(fmt.Sprintf("-> <%s> Disconnected", client.Conn.RemoteAddr()))
	client.Conn.Close()
}
