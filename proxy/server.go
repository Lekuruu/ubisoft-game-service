package proxy

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/lekuruu/ubisoft-game-service/common"
)

type Proxy struct {
	Host   string
	Port   uint16
	Logger common.Logger
}

func (server *Proxy) Serve() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Host, server.Port))

	if err != nil {
		log.Fatal(err)
	}

	server.Logger.Info(fmt.Sprintf("Listening on %s:%d", server.Host, server.Port))

	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go server.HandleClient(conn)
	}
}

func (server *Proxy) HandleClient(conn net.Conn) {
	defer server.OnDisconnect(conn)

	server.Logger.Info(fmt.Sprintf("-> <%s>", conn.RemoteAddr()))
	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			break
		}

		server.Logger.Debug(fmt.Sprintf("<- %s", line))
	}
}

func (server *Proxy) OnDisconnect(conn net.Conn) {
	if r := recover(); r != nil {
		server.Logger.Error(fmt.Sprintf("Panic: %s", r))
	}

	server.Logger.Info(fmt.Sprintf("-> <%s> Disconnected", conn.RemoteAddr()))
	conn.Close()
}
