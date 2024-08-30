package gsnat

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/lekuruu/ubisoft-game-service/common"
)

type GSNatServer struct {
	Port     uint16
	Logger   common.Logger
	Listener net.PacketConn
}

type Client struct {
	Address net.Addr
	Reader  bytes.Reader
	Server  *GSNatServer
}

func (gsn *GSNatServer) Serve() {
	listener, err := net.ListenPacket("udp", fmt.Sprintf(":%d", gsn.Port))

	if err != nil {
		log.Fatal(err)
	}

	gsn.Logger.Info(fmt.Sprintf("Listening on :%d", gsn.Port))
	gsn.Listener = listener

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
			Server:  gsn,
		}

		go gsn.HandleClient(client)
	}
}

func (cdks *GSNatServer) HandleClient(client *Client) {
	defer cdks.HandlePanic(client)

	for {
		srp, err := ReadSRPPacket(client)

		if srp == nil {
			continue
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			cdks.Logger.Error(fmt.Sprintf("Failed to parse packet: %s", err))
			return
		}

		cdks.Logger.Debug(fmt.Sprintf("-> %s", srp.String()))
		HandlePacket(client, srp)
	}
}

func (cdks *GSNatServer) HandlePanic(client *Client) {
	if r := recover(); r != nil {
		cdks.Logger.Error(fmt.Sprintf("Panic: %s", r))
	}
}

func HandlePacket(client *Client, packet *SRPPacket) {
	if packet.Flags&FLAGS_SYN == 0 {
		return
	}

	response := NewSRPPacketFromRequest(packet)

	_, err := client.Server.Listener.WriteTo(response.Serialize(), client.Address)
	if err != nil {
		client.Server.Logger.Error(fmt.Sprintf("Failed to send packet: %s", err))
	}

	client.Server.Logger.Debug(fmt.Sprintf("<- %s", response.String()))
}
