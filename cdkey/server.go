package cdkey

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"

	"github.com/lekuruu/ubisoft-game-service/common"
)

type CDKeyServer struct {
	Port     uint16
	Logger   common.Logger
	Listener net.PacketConn
}

type Client struct {
	Address net.Addr
	Reader  bytes.Reader
	Server  *CDKeyServer
}

func (cdks *CDKeyServer) Serve() {
	listener, err := net.ListenPacket("udp", fmt.Sprintf(":%d", cdks.Port))

	if err != nil {
		log.Fatal(err)
	}

	cdks.Logger.Info(fmt.Sprintf("Listening on :%d", cdks.Port))
	cdks.Listener = listener

	defer listener.Close()

	for {
		buffer := make([]byte, PACKET_BUFFER_SIZE)
		_, addr, err := listener.ReadFrom(buffer)

		if err != nil {
			log.Fatal(err)
		}

		client := &Client{
			Reader:  *bytes.NewReader(buffer),
			Address: addr,
			Server:  cdks,
		}

		go cdks.HandleClient(client)
	}
}

func (cdks *CDKeyServer) HandleClient(client *Client) {
	defer cdks.HandlePanic(client)

	for {
		msg, err := ReadCDKeyMessage(client)

		if err == io.EOF {
			break
		}

		if msg == nil {
			continue
		}

		if err != nil {
			cdks.Logger.Error(fmt.Sprintf("Failed to parse header: %s", err))
			break
		}

		if msg.Type != 211 {
			cdks.Logger.Warning(fmt.Sprintf("Received message with unknown type '%d'", msg.Type))
			break
		}

		requestTypeString, err := common.GetStringListItem(msg.Data, 1)
		if err != nil {
			cdks.Logger.Error(fmt.Sprintf("Failed to parse message ID: %s", err))
			break
		}

		requestType, err := strconv.Atoi(requestTypeString)
		if err != nil {
			cdks.Logger.Error(fmt.Sprintf("Failed to parse message ID: %s", err))
			break
		}

		cdks.Logger.Debug(fmt.Sprintf("-> %v", msg.String()))
		handler, ok := CDKeyHandlers[requestType]

		if !ok {
			cdks.Logger.Warning(fmt.Sprintf("Couldn't find handler for type '%d'", msg.Type))
			continue
		}

		response, err := handler(msg, client)

		if err != nil {
			cdks.Logger.Error(fmt.Sprintf("Failed to handle message: %s", err))
			break
		}

		serialized, err := response.Serialize()

		if err != nil {
			cdks.Logger.Error(fmt.Sprintf("Failed to serialize message: %s", err))
			break
		}

		_, err = client.Server.Listener.WriteTo(serialized, client.Address)

		if err != nil {
			cdks.Logger.Error(fmt.Sprintf("Failed to send message: %s", err))
			break
		}

		cdks.Logger.Debug(fmt.Sprintf("<- %v", response.String()))
	}
}

func (cdks *CDKeyServer) HandlePanic(client *Client) {
	if r := recover(); r != nil {
		cdks.Logger.Error(fmt.Sprintf("Panic: %s", r))
	}
}
