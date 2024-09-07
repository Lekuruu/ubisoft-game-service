package common

import (
	"io"

	"gopkg.in/irc.v4"
)

var blowfishKeyIrc = []byte{
	0x06, 0xE2, 0xC8, 0x46, 0x01, 0x90, 0x55, 0x7C,
	0x3C, 0xA1, 0xCD, 0xA3, 0xE3, 0xA1, 0x10, 0x6C,
}

var blowfishIrc = NewBlowfishCipher(blowfishKeyIrc)

func ReadIrcRequestRaw(reader io.Reader) (string, error) {
	sizeBytes := make([]byte, 2)
	_, err := reader.Read(sizeBytes)

	if err != nil {
		return "", err
	}

	size := ReadU16BE(sizeBytes)
	if size == 0 {
		return "", nil
	}

	data := make([]byte, size)
	_, err = reader.Read(data)

	if err != nil {
		return "", err
	}

	decrypted, err := blowfishIrc.Decrypt(data)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

func WriteIrcResponseRaw(writer io.Writer, data string) error {
	encrypted, err := blowfishIrc.Encrypt([]byte(data))
	if err != nil {
		return err
	}

	size := len(encrypted)
	sizeBytes := WriteU16BE(size)

	_, err = writer.Write(sizeBytes)
	if err != nil {
		return err
	}

	_, err = writer.Write(encrypted)
	return err
}

func ReadIrcRequest(reader io.Reader) (*irc.Message, error) {
	data, err := ReadIrcRequestRaw(reader)
	if err != nil {
		return nil, err
	}

	msg, err := irc.ParseMessage(data)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func WriteIrcResponse(writer io.Writer, msg *irc.Message) error {
	data := msg.String()
	return WriteIrcResponseRaw(writer, data)
}
