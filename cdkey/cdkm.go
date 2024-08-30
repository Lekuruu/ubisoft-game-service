package cdkey

import (
	"fmt"

	"github.com/lekuruu/ubisoft-game-service/common"
)

const CDKM_HEADER_SIZE = 5

const (
	CDKM_E_PLAYER_UNKNOWN = 0
	CDKM_E_PLAYER_INVALID = 1
	CDKM_E_PLAYER_VALID   = 2
)

const (
	CDKM_CHALLENGE       = 1
	CDKM_ACTIVATION      = 2
	CDKM_AUTH            = 3
	CDKM_VALIDATION      = 4
	CDKM_PLAYER_STATUS   = 5
	CDKM_DISCONNECT_USER = 6
	CDKM_STILL_ALIVE     = 7
)

var BlowfishKey = []byte("SKJDHF$0maoijfn4i8$aJdnv1jaldifar93-AS_dfo;hjhC4jhflasnF3fnd")
var Blowfish = common.NewBlowfishCipher(BlowfishKey)

type CDKeyMessage struct {
	Type uint8
	Size uint32
	Data []interface{}
}

func (msg *CDKeyMessage) Serialize() ([]byte, error) {
	dataList, err := common.SerializeDataList(msg.Data)
	if err != nil {
		return nil, err
	}

	encrypted, err := Blowfish.Encrypt(dataList)
	if err != nil {
		return nil, err
	}

	msg.Size = uint32(len(encrypted))
	header := make([]byte, CDKM_HEADER_SIZE)
	header = append(header, common.WriteU8(msg.Type)...)
	header = append(header, common.WriteU32BE(msg.Size)...)
	header = append(header, encrypted...)
	return header, nil
}

func (msg *CDKeyMessage) String() string {
	return fmt.Sprintf(
		"CDKeyMessage{Size: %d, Type: %d, Data: %v}",
		msg.Size, msg.Type, msg.Data,
	)
}

func ReadCDKeyMessage(client *Client) (*CDKeyMessage, error) {
	header := make([]byte, CDKM_HEADER_SIZE)
	_, err := client.Reader.Read(header)

	if err != nil {
		return nil, err
	}

	msg := &CDKeyMessage{
		Type: header[0],
		Size: common.ReadU32BE(header[1:5]),
	}

	if msg.Size == 0 && msg.Type == 0 {
		// Empty data, do nothing
		return nil, nil
	}

	data := make([]byte, msg.Size)
	_, err = client.Reader.Read(data)

	if err != nil {
		return nil, err
	}

	decrypted, err := Blowfish.Decrypt(data)

	if err != nil {
		return nil, err
	}

	dataList, err := common.DeserializeDataList(decrypted)

	if err != nil {
		return nil, err
	}

	msg.Data = dataList
	return msg, nil
}

func NewCDKeyMessageFromRequest(request *CDKeyMessage) *CDKeyMessage {
	return &CDKeyMessage{
		Type: request.Type,
		Size: request.Size,
		Data: []interface{}{
			request.Data[0],
			request.Data[1],
			request.Data[2],
			[]interface{}{},
		},
	}
}
