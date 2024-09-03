package router

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/lekuruu/ubisoft-game-service/common"
)

// A map to store the handlers for each message type
var RouterHandlers = map[uint8]func(*GSMessage, *Client) (*GSMessage, GSError){}
var LobbyHandlers = map[int]func(*GSMessage, *Client) (*GSMessage, GSError){}

func stillAlive(message *GSMessage, _ *Client) (*GSMessage, GSError) {
	return NewGSMessageFromRequest(message), nil
}

func handleKeyExchange(message *GSMessage, client *Client) (*GSMessage, GSError) {
	requestId, err := common.GetStringListItem(message.Data, 0)
	if err != nil {
		return nil, &RouterError{Message: err.Error()}
	}

	requestArgs, err := common.GetListItem(message.Data, 1)
	if err != nil {
		return nil, &RouterError{Message: err.Error()}
	}

	response := NewGSMessageFromRequest(message)
	response.Data = []interface{}{requestId}
	responseArgs := []interface{}{"1"}

	switch requestId {
	case "1":
		// RSA Encryption
		rsaBuffer, err := common.GetBinaryListItem(requestArgs, 2)
		if err != nil {
			return nil, &RouterError{Message: err.Error()}
		}

		client.GamePublicKey = common.RsaPublicKeyFromBuffer(rsaBuffer)
		privateKey, err := common.RsaKeygen()
		if err != nil {
			return nil, &RouterError{Message: err.Error()}
		}

		client.ServerPrivateKey = privateKey
		client.ServerPublicKey = &privateKey.PublicKey

		keyData := common.RsaPublicKeyToBuffer(&privateKey.PublicKey)
		responseArgs = append(responseArgs, fmt.Sprint(len(keyData)))
		responseArgs = append(responseArgs, keyData)

	case "2":
		// Blowfish encryption
		if client.GamePublicKey == nil {
			return nil, &RouterError{Message: "game public key not initialized"}
		}

		encryptedBlowfishKey, err := common.GetBinaryListItem(requestArgs, 2)
		if err != nil {
			return nil, &RouterError{Message: err.Error()}
		}

		blowfishKey, err := client.ServerPrivateKey.Decrypt(rand.Reader, encryptedBlowfishKey, nil)
		if err != nil {
			return nil, &RouterError{Message: err.Error()}
		}

		client.GameBlowfishKey = blowfishKey
		client.ServerBlowfishKey = common.BlowfishKeygen(16)

		encryptedKey, err := rsa.EncryptPKCS1v15(rand.Reader, client.GamePublicKey, client.ServerBlowfishKey)
		if err != nil {
			return nil, &RouterError{Message: err.Error()}
		}

		responseArgs = append(responseArgs, fmt.Sprint(len(encryptedKey)))
		responseArgs = append(responseArgs, encryptedKey)

	default:
		return nil, &RouterError{Message: "invalid request id"}
	}

	response.Data = append(response.Data, responseArgs)
	return response, nil
}

func handleLogin(message *GSMessage, client *Client) (*GSMessage, GSError) {
	username, err := common.GetStringListItem(message.Data, 0)
	if err != nil {
		return nil, &RouterError{Message: err.Error()}
	}

	version, err := common.GetStringListItem(message.Data, 2)
	if err != nil {
		return nil, &RouterError{Message: err.Error()}
	}

	public, err := common.GetBoolListItem(message.Data, 3)
	if err != nil {
		return nil, &RouterError{Message: err.Error()}
	}

	// TODO: Implement login validation
	// password, err := common.GetStringListItem(message.Data, 1)

	if player := client.Server.Players.ByName(username); player != nil {
		// Player already logged in
		return nil, &RouterError{
			Message:      "player already logged in",
			ResponseCode: ERRORROUTER_NOTDISCONNECTED,
		}
	}

	return nil, &RouterError{
		Message:      "player already logged in",
		ResponseCode: ERRORROUTER_NOTDISCONNECTED,
	}

	// Create initial player object
	player := &Player{
		Name:    username,
		Version: version,
		Info:    Info{Public: public},
	}

	// Add player to pending waitmodule logins
	ipAddress := strings.Split(client.Conn.RemoteAddr().String(), ":")[0]
	client.Server.Pending[ipAddress] = player

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

func handleWaitModuleJoin(message *GSMessage, client *Client) (*GSMessage, GSError) {
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

func handleWaitModuleLogin(message *GSMessage, client *Client) (*GSMessage, GSError) {
	username, err := common.GetStringListItem(message.Data, 0)
	if err != nil {
		return nil, &RouterError{Message: err.Error()}
	}

	if player := client.Server.Players.ByName(username); player != nil {
		// Player already logged in
		return nil, &RouterError{
			Message:      "player already logged in",
			ResponseCode: ERRORROUTER_NOTDISCONNECTED,
		}
	}

	ipAddress := strings.Split(client.Conn.RemoteAddr().String(), ":")[0]
	player, ok := client.Server.Pending[ipAddress]

	if !ok {
		return nil, &LobbyError{ErrorMessage: "ip not found in pending login list"}
	}

	if player.Name != username {
		return nil, &LobbyError{ErrorMessage: "username mismatch"}
	}

	// Remove pending login
	delete(client.Server.Pending, ipAddress)

	player.Client = *client
	client.Player = player
	client.Server.Players.Add(player)

	response := NewGSMessageFromRequest(message)
	response.Property = PROPERTY_GS
	response.Type = GSM_GSSUCCESS
	response.Data = []interface{}{common.WriteU8(GSM_LOGINWAITMODULE)}
	return response, nil
}

func handlePlayerInfo(message *GSMessage, client *Client) (*GSMessage, GSError) {
	targetName, err := common.GetStringListItem(message.Data, 0)
	if err != nil {
		return nil, &RouterError{Message: err.Error()}
	}

	player := client.Server.Players.ByName(targetName)
	if player == nil {
		return nil, &RouterError{
			ResponseCode: ERRORROUTER_NOTREGISTERED,
			Message:      "player was not found",
		}
	}

	if !player.Info.Public && player != client.Player {
		return nil, &RouterError{
			ResponseCode: ERRORROUTER_NOTREGISTERED,
			Message:      "player info is not public",
		}
	}

	playerData := []interface{}{
		player.Name, player.Info.Surname, player.Info.Firstname,
		player.Info.Country, player.Info.Email, "IRCID", player.IpAddress(),
	}

	response := NewGSMessageFromRequest(message)
	response.Type = GSM_GSSUCCESS
	response.Data = []interface{}{common.WriteU8(GSM_PLAYERINFO), playerData}
	return response, nil
}

func handleLobbyMessage(message *GSMessage, client *Client) (*GSMessage, GSError) {
	subTypeString, err := common.GetStringListItem(message.Data, 0)
	if err != nil {
		return nil, &LobbyError{ErrorMessage: err.Error()}
	}

	subType, err := strconv.Atoi(subTypeString)
	if err != nil {
		return nil, &LobbyError{ErrorMessage: err.Error()}
	}

	handler, ok := LobbyHandlers[subType]
	if !ok {
		return nil, &LobbyError{ErrorMessage: fmt.Sprintf("lobby handler for '%s' not found", subTypeString)}
	}

	return handler(message, client)
}

func handleLobbyLogin(message *GSMessage, client *Client) (*GSMessage, GSError) {
	requestArgs, err := common.GetListItem(message.Data, 1)
	if err != nil {
		return nil, &LobbyError{ErrorMessage: err.Error()}
	}

	gameName, err := common.GetStringListItem(requestArgs, 0)
	if err != nil {
		return nil, &LobbyError{ErrorMessage: err.Error()}
	}

	i := sort.SearchStrings(client.Server.Games, gameName)

	// Check if game is supported
	if i >= len(client.Server.Games) || client.Server.Games[i] != gameName {
		return nil, &LobbyError{ErrorMessage: "game not supported"}
	}

	client.Player.Game = gameName
	response := NewGSMessageFromRequest(message)
	response.Data = []interface{}{
		strconv.Itoa(GSM_GSSUCCESS),
		[]interface{}{strconv.Itoa(LOBBY_LOGIN)},
	}

	return response, nil
}

func handleFriendsLogin(message *GSMessage, client *Client) (*GSMessage, GSError) {
	status, err := common.GetU32ListItem(message.Data, 0)
	if err != nil {
		return nil, &RouterError{Message: err.Error()}
	}

	mood, err := common.GetU32ListItem(message.Data, 1)
	if err != nil {
		return nil, &RouterError{Message: err.Error()}
	}

	client.Player.Friends.Status = status
	client.Player.Friends.Mood = mood

	// TODO: Load relationships from database
	client.Player.Friends.List = NewPlayerCollection()
	client.Player.Friends.Ignored = NewPlayerCollection()

	response := NewGSMessageFromRequest(message)
	response.Type = GSM_GSSUCCESS
	response.Data = []interface{}{common.WriteU8(GSM_LOGINFRIENDS)}
	return response, nil
}

func handleMotdRequest(message *GSMessage, client *Client) (*GSMessage, GSError) {
	// language, err := common.GetStringListItem(message.Data, 0)
	response := NewGSMessageFromRequest(message)
	response.Type = GSM_GSSUCCESS
	response.Data = []interface{}{
		common.WriteU8(GSM_MOTD_REQUEST),
		[]interface{}{
			"Welcome to the server!",  // szUbiMOTD (UBI's MOTD)
			"This is a test message.", // szGameMOTD (Game's MOTD)
		},
	}
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
	RouterHandlers[GSM_MOTD_REQUEST] = handleMotdRequest

	LobbyHandlers[LOBBY_LOGIN] = handleLobbyLogin
}
