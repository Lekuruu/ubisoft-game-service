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
		blowfishKey, err := client.ServerPrivateKey.Decrypt(nil, encryptedBlowfishKey, nil)
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
	}

	response.Data = append(response.Data, responseArgs)
	return response, nil
}

func init() {
	RouterHandlers[GSM_KEY_EXCHANGE] = handleKeyExchange
}
