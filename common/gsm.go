package common

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"io"
)

const MAX_PACKET_SIZE = 0x50000
const GSMSG_HEADER_SIZE = 6

const (
	GSM_PROPERTY_GS         = 0
	GSM_PROPERTY_GAME       = 1
	GSM_PROPERTY_GS_ENCRYPT = 2
)

type GSMessage struct {
	Size     uint32
	Property uint8
	Priority uint8
	Type     uint8
	Sender   uint8
	Receiver uint8
	Data     []interface{}
}

type GSClientState struct {
	GamePublicKey     *rsa.PublicKey
	GameBlowfishKey   []byte
	ServerPublicKey   *rsa.PublicKey
	ServerPrivateKey  *rsa.PrivateKey
	ServerBlowfishKey []byte
}

// Serialize a GSMessage to be sent to the client
func (msg *GSMessage) Serialize(state *GSClientState) ([]byte, error) {
	data, err := SerializeDataList(msg.Data)
	if err != nil {
		return nil, err
	}

	encrypted, err := EncryptDataList(data, msg.Property, state)
	if err != nil {
		return nil, err
	}

	msg.Size = uint32(len(encrypted) + GSMSG_HEADER_SIZE)
	header := make([]byte, GSMSG_HEADER_SIZE)
	header[0] = byte(msg.Size >> 16)
	header[1] = byte(msg.Size >> 8)
	header[2] = byte(msg.Size)
	header[3] &= 0x3F
	header[3] |= (msg.Property << 6)
	header[3] |= msg.Priority & 0x20
	header[4] = msg.Type
	header[5] &= 0xF
	header[5] |= 0x10 * msg.Sender
	header[5] &= 0xF0
	header[5] |= msg.Receiver & 0xF

	return append(header, encrypted...), nil
}

// Format a GSMessage to be logged
func (msg *GSMessage) String() string {
	return fmt.Sprintf(
		"GSMessage{Size: %d, Property: %d, Priority: %d, Type: %d, Sender: %d, Receiver: %d, Data: %v}",
		msg.Size, msg.Property, msg.Priority, msg.Type, msg.Sender, msg.Receiver, msg.Data,
	)
}

// Read a GSMessage from the client
func ReadGSMessage(reader io.Reader, state *GSClientState) (*GSMessage, error) {
	header := make([]byte, GSMSG_HEADER_SIZE)
	_, err := reader.Read(header)

	if err != nil {
		return nil, err
	}

	if len(header) != GSMSG_HEADER_SIZE {
		return nil, errors.New("invalid data size")
	}

	size := (int(header[0]) << 16) + (int(header[1]) << 8) + int(header[2])
	property := (header[3] >> 6)
	priority := (header[3] & 0x3F)
	msgType := (header[4])
	sender := (header[5] >> 4)
	receiver := (header[5] & 0x0F)

	if size > MAX_PACKET_SIZE {
		return nil, errors.New("requested packet size too large")
	}

	data := make([]byte, size-GSMSG_HEADER_SIZE)
	_, err = reader.Read(data)

	if err != nil {
		return nil, err
	}

	dataList, err := DecryptDataList(
		data,
		property,
		state,
	)

	if err != nil {
		return nil, err
	}

	return &GSMessage{
		Size:     uint32(size),
		Property: property,
		Priority: priority,
		Type:     msgType,
		Sender:   sender,
		Receiver: receiver,
		Data:     dataList,
	}, nil
}

// Create a new GSMessage from a request, which can be used to send a response
func NewGSMessageFromRequest(request *GSMessage) *GSMessage {
	return &GSMessage{
		Property: request.Property,
		Priority: request.Priority,
		Type:     request.Type,
		Sender:   request.Receiver,
		Receiver: request.Sender,
		Data:     request.Data,
	}
}

// Encrypt serialized data list
func EncryptDataList(data []byte, property uint8, state *GSClientState) ([]byte, error) {
	switch property {
	case GSM_PROPERTY_GS:
		return GSXOREncrypt(data), nil

	case GSM_PROPERTY_GS_ENCRYPT:
		cipher := NewBlowfishCipher(state.GameBlowfishKey)
		return cipher.Encrypt(data)

	default:
		return data, nil
	}
}

// Decrypt & deserialize data list
func DecryptDataList(data []byte, property uint8, state *GSClientState) ([]interface{}, error) {
	switch property {
	case GSM_PROPERTY_GS_ENCRYPT:
		if state.GameBlowfishKey == nil {
			return nil, errors.New("blowfish key not initialized")
		}
		cipher := NewBlowfishCipher(state.GameBlowfishKey)
		decrypted, err := cipher.Decrypt(data)
		if err != nil {
			return nil, err
		}
		return DeserializeDataList(decrypted)

	case GSM_PROPERTY_GS:
		decrypted := GSXORDecrypt(data)
		return DeserializeDataList(decrypted)

	default:
		return DeserializeDataList(data)
	}
}
