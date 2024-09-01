package router

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/lekuruu/ubisoft-game-service/common"
)

// A map to store the handlers for each message type
var RouterHandlers = map[uint8]func(*GSMessage, *Client) (*GSMessage, error){}
var LobbyHandlers = map[int]func(*GSMessage, *Client) (*GSMessage, error){}

func stillAlive(message *GSMessage, _ *Client) (*GSMessage, error) {
	return NewGSMessageFromRequest(message), nil
}

func handleKeyExchange(message *GSMessage, client *Client) (*GSMessage, error) {
	requestId, err := common.GetStringListItem(message.Data, 0)
	if err != nil {
		return nil, err
	}

	requestArgs, err := common.GetListItem(message.Data, 1)
	if err != nil {
		return nil, err
	}

	response := NewGSMessageFromRequest(message)
	response.Data = []interface{}{requestId}
	responseArgs := []interface{}{"1"}

	switch requestId {
	case "1":
		// RSA Encryption
		rsaBuffer, err := common.GetBinaryListItem(requestArgs, 2)
		if err != nil {
			return nil, err
		}

		client.GamePublicKey = common.RsaPublicKeyFromBuffer(rsaBuffer)
		privateKey, err := common.RsaKeygen()
		if err != nil {
			return nil, err
		}

		client.ServerPrivateKey = privateKey
		client.ServerPublicKey = &privateKey.PublicKey

		keyData := common.RsaPublicKeyToBuffer(&privateKey.PublicKey)
		responseArgs = append(responseArgs, fmt.Sprint(len(keyData)))
		responseArgs = append(responseArgs, keyData)

	case "2":
		// Blowfish encryption
		if client.GamePublicKey == nil {
			return nil, errors.New("game public key not initialized")
		}

		encryptedBlowfishKey, err := common.GetBinaryListItem(requestArgs, 2)
		if err != nil {
			return nil, err
		}

		blowfishKey, err := client.ServerPrivateKey.Decrypt(rand.Reader, encryptedBlowfishKey, nil)
		if err != nil {
			return nil, err
		}

		client.GameBlowfishKey = blowfishKey
		client.ServerBlowfishKey = common.BlowfishKeygen(16)

		encryptedKey, err := rsa.EncryptPKCS1v15(rand.Reader, client.GamePublicKey, client.ServerBlowfishKey)
		if err != nil {
			return nil, err
		}

		responseArgs = append(responseArgs, fmt.Sprint(len(encryptedKey)))
		responseArgs = append(responseArgs, encryptedKey)

	default:
		return nil, errors.New("invalid request id")
	}

	response.Data = append(response.Data, responseArgs)
	return response, nil
}

func handleLogin(message *GSMessage, client *Client) (*GSMessage, error) {
	username, err := common.GetStringListItem(message.Data, 0)
	if err != nil {
		return nil, err
	}

	if player := client.Server.Players.ByName(username); player != nil {
		// Player already logged in
		return NewGSErrorMessage(ERRORROUTER_NOTDISCONNECTED, message), nil
	}

	// TODO: Implement login validation
	// password, err := common.GetStringListItem(message.Data, 1)
	// game, err := common.GetStringListItem(message.Data, 2)

	// Setup pending waitmodule login
	ipAddress := strings.Split(client.Conn.RemoteAddr().String(), ":")[0]
	client.Server.Pending[ipAddress] = username

	// Remove pending login after 5 seconds
	time.AfterFunc(5*time.Second, func() {
		delete(client.Server.Pending, ipAddress)
	})

	response := NewGSMessageFromRequest(message)
	response.Property = PROPERTY_GS
	response.Type = GSM_GSSUCCESS
	response.Data = []interface{}{common.WriteU8(GSM_LOGIN)}
	return response, nil
}

func handleWaitModuleJoin(message *GSMessage, client *Client) (*GSMessage, error) {
	// NOTE: The WaitModule server is not implemented in this project yet, so
	//		 that's why we are sending the router's host and port as the
	//		 WaitModule connection info.
	response := NewGSMessageFromRequest(message)
	response.Type = GSM_GSSUCCESS
	response.Data = []interface{}{common.WriteU8(GSM_JOINWAITMODULE)}

	waitModuleInfo := []interface{}{}
	waitModuleInfo = append(waitModuleInfo, client.Server.Host)
	waitModuleInfo = append(waitModuleInfo, common.WriteU32(client.Server.Port))
	response.Data = append(response.Data, waitModuleInfo)

	return response, nil
}

func handleWaitModuleLogin(message *GSMessage, client *Client) (*GSMessage, error) {
	username, err := common.GetStringListItem(message.Data, 0)
	if err != nil {
		return nil, err
	}

	if player := client.Server.Players.ByName(username); player != nil {
		// Player already logged in
		return NewGSErrorMessage(ERRORROUTER_NOTDISCONNECTED, message), nil
	}

	ipAddress := strings.Split(client.Conn.RemoteAddr().String(), ":")[0]
	pendingUsername, ok := client.Server.Pending[ipAddress]

	if !ok {
		return nil, errors.New("ip not found in pending login list")
	}

	if pendingUsername != username {
		return nil, errors.New("username mismatch")
	}

	// Remove pending login
	delete(client.Server.Pending, ipAddress)

	player := &Player{
		Client: *client,
		Name:   username,
		Nick:   username,
	}

	client.Player = player
	client.Server.Players.Add(player)

	response := NewGSMessageFromRequest(message)
	response.Property = PROPERTY_GS
	response.Type = GSM_GSSUCCESS
	response.Data = []interface{}{common.WriteU8(GSM_LOGINWAITMODULE)}
	return response, nil
}

func handlePlayerInfo(message *GSMessage, client *Client) (*GSMessage, error) {
	targetName, err := common.GetStringListItem(message.Data, 0)
	if err != nil {
		return nil, err
	}

	player := client.Server.Players.ByName(targetName)
	if player == nil {
		return NewGSErrorMessage(ERRORROUTER_NOTREGISTERED, message), nil
	}

	response := NewGSMessageFromRequest(message)
	response.Type = GSM_GSSUCCESS
	playerData := []interface{}{player.Nick, player.Name, "findme3", "findme4", "findme5", "findme6", "findme7"}
	response.Data = []interface{}{common.WriteU8(GSM_PLAYERINFO), playerData}
	return response, nil
}

func handleLobbyMessage(message *GSMessage, client *Client) (*GSMessage, error) {
	lobbyMessageTypeString, err := common.GetStringListItem(message.Data, 0)
	if err != nil {
		return nil, err
	}

	lobbyMessageType, err := strconv.Atoi(lobbyMessageTypeString)
	if err != nil {
		return nil, err
	}

	handler, ok := LobbyHandlers[lobbyMessageType]
	if !ok {
		return nil, fmt.Errorf("lobby handler for '%s' not found", lobbyMessageTypeString)
	}

	return handler(message, client)
}

func handleLobbyLogin(message *GSMessage, client *Client) (*GSMessage, error) {
	requestArgs, err := common.GetListItem(message.Data, 1)
	if err != nil {
		return nil, err
	}

	gameName, err := common.GetStringListItem(requestArgs, 0)
	if err != nil {
		return nil, err
	}

	i := sort.SearchStrings(client.Server.Games, gameName)

	// Check if game is supported
	if i >= len(client.Server.Games) || client.Server.Games[i] != gameName {
		return nil, errors.New("game not supported")
	}

	client.Player.Game = gameName
	response := NewGSMessageFromRequest(message)
	response.Data = []interface{}{
		strconv.Itoa(GSM_GSSUCCESS),
		[]interface{}{strconv.Itoa(LOBBY_LOGIN)},
	}

	return response, nil
}

func handleFriendsLogin(message *GSMessage, client *Client) (*GSMessage, error) {
	response := NewGSMessageFromRequest(message)
	response.Type = GSM_GSSUCCESS
	response.Data = []interface{}{common.WriteU8(GSM_LOGINFRIENDS)}
	return response, nil
}

func init() {
	RouterHandlers[GSM_STILLALIVE] = stillAlive
	RouterHandlers[GSM_KEY_EXCHANGE] = handleKeyExchange
	RouterHandlers[GSM_LOGIN] = handleLogin
	RouterHandlers[GSM_JOINWAITMODULE] = handleWaitModuleJoin
	RouterHandlers[GSM_LOGINWAITMODULE] = handleWaitModuleLogin
	RouterHandlers[GSM_PLAYERINFO] = handlePlayerInfo
	RouterHandlers[GSM_LOBBY_MSG] = handleLobbyMessage
	RouterHandlers[GSM_LOGINFRIENDS] = handleFriendsLogin

	LobbyHandlers[LOBBY_LOGIN] = handleLobbyLogin
}
