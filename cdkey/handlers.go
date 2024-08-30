package cdkey

import (
	"strconv"

	"github.com/lekuruu/ubisoft-game-service/router"
)

// A map to store the handlers for each message type
var CDKeyHandlers = map[int]func(*CDKeyMessage, *Client) (*CDKeyMessage, error){}

func handleChallenge(msg *CDKeyMessage, client *Client) (*CDKeyMessage, error) {
	response := NewCDKeyMessageFromRequest(msg)
	hash := []byte{
		0x00, 0x11, 0x22, 0x33,
		0x44, 0x55, 0x66, 0x77,
		0x88, 0x99, 0xaa, 0xbb,
		0xcc, 0xdd, 0xee, 0xff,
		0x01, 0x02, 0x03, 0x04,
	}
	response.Data[3] = []interface{}{
		strconv.Itoa(router.GSM_GSSUCCESS),
		hash,
	}
	return response, nil
}

func handleActivation(msg *CDKeyMessage, client *Client) (*CDKeyMessage, error) {
	response := NewCDKeyMessageFromRequest(msg)
	activationId := []byte{
		0x33, 0x33, 0x33, 0x33, 0x33, 0x33,
		0x33, 0x33, 0x33, 0x33, 0x33,
	}
	buffer := []byte{
		0x44, 0x44, 0x44, 0x44, 0x44, 0x44,
		0x44, 0x44, 0x44, 0x44, 0x44,
	}
	response.Data[3] = []interface{}{
		strconv.Itoa(router.GSM_GSSUCCESS),
		[]interface{}{activationId, buffer},
	}
	return response, nil
}

func init() {
	CDKeyHandlers[CDKM_CHALLENGE] = handleChallenge
	CDKeyHandlers[CDKM_ACTIVATION] = handleActivation
}
