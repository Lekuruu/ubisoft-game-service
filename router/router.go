package router

import (
	"crypto/rsa"
	"fmt"
	"log"
	"net"

	"github.com/lekuruu/ubisoft-game-service/common"
)

type Router struct {
	Host   string
	Port   uint16
	Logger common.Logger
}

type Client struct {
	Conn              net.Conn
	GamePublicKey     *rsa.PublicKey
	GameBlowfishKey   []byte
	ServerPublicKey   *rsa.PublicKey
	ServerPrivateKey  *rsa.PrivateKey
	ServerBlowfishKey []byte
}

func (router *Router) Serve() {
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
	defer router.OnClose(conn)
	router.Logger.Info(fmt.Sprintf("-> <%s>", conn.RemoteAddr()))

	client := &Client{Conn: conn}

	for {
		msg, err := ReadGSMessage(client)

		if err != nil {
			router.Logger.Error(fmt.Sprintf("Failed to parse header: %s", err))
			break
		}

		router.Logger.Debug(fmt.Sprintf("-> %v", msg.String()))
		handler, ok := RouterHandlers[msg.Type]

		if !ok {
			router.Logger.Warning(fmt.Sprintf("Couldn't find handler for type '%d'", msg.Type))
			break
		}

		response, err := handler(msg, client)

		if err != nil {
			router.Logger.Error(fmt.Sprintf("Failed to handle message: %s", err))
			break
		}

		router.Logger.Debug(fmt.Sprintf("<- %v", response.String()))
		serialized, err := response.Serialize(client)

		if err != nil {
			router.Logger.Error(fmt.Sprintf("Failed to serialize message: %s", err))
			break
		}

		_, err = conn.Write(serialized)

		if err != nil {
			router.Logger.Error(fmt.Sprintf("Failed to send message: %s", err))
			break
		}
	}
}

func (router *Router) OnClose(conn net.Conn) {
	router.Logger.Info(fmt.Sprintf("-> <%s> Disconnected", conn.RemoteAddr()))
	conn.Close()
}
