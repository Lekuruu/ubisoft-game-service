package router

import (
	"errors"
	"fmt"

	"github.com/lekuruu/ubisoft-game-service/common"
)

const GSMSG_HEADER_SIZE = 6

const (
	GSM_NEWUSERREQUEST         = 1
	GSM_CONNECTIONREQUEST      = 2
	GSM_PLAYERNEW              = 3
	GSM_DISCONNECTION          = 4
	GSM_PLAYERREMOVED          = 5
	GSM_EVENT_UDPCONNECT       = 6
	GSM_NEWS                   = 7
	GSM_SEARCHPLAYER           = 8
	GSM_REMOVEACCOUNT          = 9
	GSM_SERVERSLIST            = 11
	GSM_SESSIONLIST            = 13
	GSM_PLAYERLIST             = 15
	GSM_GETGROUPINFO           = 16
	GSM_GROUPINFO              = 17
	GSM_GETPLAYERINFO          = 18
	GSM_PLAYERINFO             = 19
	GSM_CHATALL                = 20
	GSM_CHATLIST               = 21
	GSM_CHATSESSION            = 22
	GSM_CHAT                   = 24
	GSM_CREATESESSION          = 26
	GSM_SESSIONNEW             = 27
	GSM_JOINSESSION            = 28
	GSM_JOINNEW                = 31
	GSM_LEAVESESSION           = 32
	GSM_JOINLEAVE              = 33
	GSM_SESSIONREMOVE          = 34
	GSM_GSSUCCESS              = 38
	GSM_GSFAIL                 = 39
	GSM_BEGINGAME              = 40
	GSM_UPDATEPLAYERINFO       = 45
	GSM_MASTERCHANGED          = 48
	GSM_UPDATESESSIONSTATE     = 51
	GSM_URGENTMESSAGE          = 52
	GSM_NEWWAITMODULE          = 54
	GSM_KILLMODULE             = 55
	GSM_STILLALIVE             = 58
	GSM_PING                   = 59
	GSM_PLAYERKICK             = 60
	GSM_PLAYERMUTE             = 61
	GSM_ALLOWGAME              = 62
	GSM_FORBIDGAME             = 63
	GSM_GAMELIST               = 64
	GSM_UPDATEADVERTISMEMENTS  = 65
	GSM_UPDATENEWS             = 66
	GSM_VERSIONLIST            = 67
	GSM_UPDATEVERSIONS         = 68
	GSM_UPDATEDISTANTROUTERS   = 70
	GSM_ADMINLOGIN             = 71
	GSM_STAT_PLAYER            = 72
	GSM_STAT_GAME              = 73
	GSM_UPDATEFRIEND           = 74
	GSM_ADDFRIEND              = 75
	GSM_DELFRIEND              = 76
	GSM_LOGINWAITMODULE        = 77
	GSM_LOGINFRIENDS           = 78
	GSM_ADDIGNOREFRIEND        = 79
	GSM_DELIGNOREFRIEND        = 80
	GSM_STATUSCHANGE           = 81
	GSM_JOINARENA              = 82
	GSM_LEAVEARENA             = 83
	GSM_IGNORELIST             = 84
	GSM_IGNOREFRIEND           = 85
	GSM_GETARENA               = 86
	GSM_GETSESSION             = 87
	GSM_PAGEPLAYER             = 88
	GSM_FRIENDLIST             = 89
	GSM_PEERMSG                = 90
	GSM_PEERPLAYER             = 91
	GSM_DISCONNECTFRIENDS      = 92
	GSM_JOINWAITMODULE         = 93
	GSM_LOGINSESSION           = 94
	GSM_DISCONNECTSESSION      = 95
	GSM_PLAYERDISCONNECT       = 96
	GSM_ADVERTISEMENT          = 97
	GSM_MODIFYUSER             = 98
	GSM_STARTGAME              = 99
	GSM_CHANGEVERSION          = 100
	GSM_PAGER                  = 101
	GSM_LOGIN                  = 102
	GSM_PHOTO                  = 103
	GSM_LOGINARENA             = 104
	GSM_SQLCREATE              = 106
	GSM_SQLSELECT              = 107
	GSM_SQLDELETE              = 108
	GSM_SQLSET                 = 109
	GSM_SQLSTAT                = 110
	GSM_SQLQUERY               = 111
	GSM_ROUTEURLIST            = 127
	GSM_DISTANCEVECTOR         = 131
	GSM_WRAPPEDMESSAGE         = 132
	GSM_CHANGEFRIEND           = 133
	GSM_NEWRELFRIEND           = 134
	GSM_DELRELFRIEND           = 135
	GSM_NEWIGNOREFRIEND        = 136
	GSM_DELETEIGNOREFRIEND     = 137
	GSM_ARENACONNECTION        = 138
	GSM_ARENADISCONNECTION     = 139
	GSM_ARENAWAITMODULE        = 140
	GSM_ARENANEW               = 141
	GSM_NEWBASICGROUP          = 143
	GSM_ARENAREMOVED           = 144
	GSM_DELETEBASICGROUP       = 145
	GSM_SESSIONSBEGIN          = 146
	GSM_GROUPDATA              = 148
	GSM_ARENA_MESSAGE          = 151
	GSM_ARENALISTREQUEST       = 157
	GSM_ROUTERPLAYERNEW        = 158
	GSM_BASEGROUPREQUEST       = 159
	GSM_UPDATEPLAYERPING       = 166
	GSM_UPDATEGROUPSIZE        = 169
	GSM_SLEEP                  = 179
	GSM_WAKEUP                 = 180
	GSM_SYSTEMPAGE             = 181
	GSM_SESSIONOPEN            = 189
	GSM_SESSIONCLOSE           = 190
	GSM_LOGINCLANMANAGER       = 192
	GSM_DISCONNECTCLANMANAGER  = 193
	GSM_CLANMANAGERPAGE        = 194
	GSM_UPDATECLANPLAYER       = 195
	GSM_PLAYERCLANS            = 196
	GSM_GETPERSISTANTGROUPINFO = 199
	GSM_UPDATEGROUPPING        = 202
	GSM_DEFERREDGAMESTARTED    = 203
	GSM_PROXY_HANDLER          = 204
	GSM_BEGINCLIENTHOSTGAME    = 205
	GSM_LOBBY_MSG              = 209
	GSM_LOBBYSERVERLOGIN       = 210
	GSM_SETGROUPSZDATA         = 211
	GSM_GROUPSZDATA            = 212
	GSM_KEY_EXCHANGE           = 219
	GSM_REQUESTPORTID          = 221
)

const (
	TARGET_R   = 1
	TARGET_S   = 2
	TARGET_W   = 3
	TARGET_P   = 4
	TARGET_AP  = 5
	TARGET_B   = 6
	TARGET_LP  = 7
	TARGET_UNK = 8
	TARGET_G   = 9
	TARGET_A   = 10
)

const (
	PROPERTY_GS         = 0
	PROPERTY_GAME       = 1
	PROPERTY_GS_ENCRYPT = 2
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

// Serialize a GSMessage to be sent to the client
func (msg *GSMessage) Serialize(client *Client) ([]byte, error) {
	data, err := common.SerializeDataList(msg.Data)
	if err != nil {
		return nil, err
	}

	encrypted := EncryptDataList(data, msg.Property, client)
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
func ReadGSMessage(client *Client) (*GSMessage, error) {
	header := make([]byte, GSMSG_HEADER_SIZE)
	_, err := client.Conn.Read(header)

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

	data := make([]byte, size-GSMSG_HEADER_SIZE)
	_, err = client.Conn.Read(data)

	if err != nil {
		return nil, err
	}

	dataList, err := DecryptDataList(
		data,
		property,
		client,
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
func EncryptDataList(data []byte, property uint8, client *Client) []byte {
	switch property {
	case PROPERTY_GS:
		return common.GSXOREncrypt(data)

	case PROPERTY_GS_ENCRYPT:
		cipher := common.NewBlowfishCipher(client.GameBlowfishKey)
		return cipher.Encrypt(data)

	default:
		return data
	}
}

// Decrypt & deserialize data list
func DecryptDataList(data []byte, property uint8, client *Client) ([]interface{}, error) {
	switch property {
	case PROPERTY_GS:
		decrypted := common.GSXORDecrypt(data)
		return common.DeserializeDataList(decrypted)

	case PROPERTY_GS_ENCRYPT:
		cipher := common.NewBlowfishCipher(client.GameBlowfishKey)
		decrypted := cipher.Decrypt(data)
		return common.DeserializeDataList(decrypted)

	default:
		return common.DeserializeDataList(data)
	}
}
