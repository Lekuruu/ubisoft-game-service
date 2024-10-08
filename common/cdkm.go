package common

import (
	"bytes"
	"fmt"
)

const CDKM_PACKET_BUFFER_SIZE = 512
const CDKM_HEADER_SIZE = 5

var blowfishKeyCd = []byte("SKJDHF$0maoijfn4i8$aJdnv1jaldifar93-AS_dfo;hjhC4jhflasnF3fnd")
var blowfishCD = NewBlowfishCipher(blowfishKeyCd)

type CDKeyMessage struct {
	Type uint8
	Size uint32
	Data []interface{}
}

func (msg *CDKeyMessage) Serialize() ([]byte, error) {
	dataList, err := SerializeDataList(msg.Data)
	if err != nil {
		return nil, err
	}

	encrypted, err := blowfishCD.Encrypt(dataList)
	if err != nil {
		return nil, err
	}

	msg.Size = uint32(len(encrypted))
	data := make([]byte, 0)
	data = append(data, WriteU8(msg.Type)...)
	data = append(data, WriteU32BE(msg.Size)...)
	data = append(data, encrypted...)
	return data, nil
}

func (msg *CDKeyMessage) String() string {
	return fmt.Sprintf(
		"CDKeyMessage{Size: %d, Type: %d, Data: %v}",
		msg.Size, msg.Type, msg.Data,
	)
}

func ReadCDKeyMessage(reader *bytes.Reader) (*CDKeyMessage, error) {
	header := make([]byte, CDKM_HEADER_SIZE)
	_, err := reader.Read(header)

	if err != nil {
		return nil, err
	}

	msg := &CDKeyMessage{
		Type: header[0],
		Size: ReadU32BE(header[1:5]),
	}

	if msg.Size == 0 || msg.Type == 0 {
		// Empty data, do nothing
		return nil, nil
	}

	if msg.Size > CDKM_PACKET_BUFFER_SIZE {
		return nil, fmt.Errorf("requested size too large: %d", msg.Size)
	}

	data := make([]byte, msg.Size)
	_, err = reader.Read(data)

	if err != nil {
		return nil, err
	}

	decrypted, err := blowfishCD.Decrypt(data)

	if err != nil {
		return nil, err
	}

	dataList, err := DeserializeDataList(decrypted)

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
