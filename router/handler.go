package router

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/lekuruu/ubisoft-game-service/common"
)

// A map to store the handlers for each message type
var RouterHandlers = map[uint8]func(*GSMessage, *Client) (*GSMessage, error){}

func stillAlive(message *GSMessage, _ *Client) (*GSMessage, error) {
	return NewGSMessageFromRequest(message), nil
}

func handleKeyExchange(message *GSMessage, client *Client) (*GSMessage, error) {
	requestId := message.Data[0].(string)
	requestArgs := message.Data[1].([]interface{})

	response := NewGSMessageFromRequest(message)
	response.Data = []interface{}{requestId}
	responseArgs := []interface{}{"1"}

	switch requestId {
	case "1":
		// RSA Encryption
		client.GamePublicKey = common.RsaPublicKeyFromBuffer(requestArgs[2].([]byte))
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

		encryptedBlowfishKey := requestArgs[2].([]byte)
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
	// username := message.Data[0].(string)
	// password := message.Data[1].(string)
	// game := message.Data[2].(string)
	// TODO: Implement login logic

	response := NewGSMessageFromRequest(message)
	response.Property = PROPERTY_GS
	response.Type = GSM_GSSUCCESS
	response.Data = []interface{}{common.WriteU8(GSM_LOGIN)}
	return response, nil
}

func init() {
	RouterHandlers[GSM_STILLALIVE] = stillAlive
	RouterHandlers[GSM_KEY_EXCHANGE] = handleKeyExchange
	RouterHandlers[GSM_LOGIN] = handleLogin
}
