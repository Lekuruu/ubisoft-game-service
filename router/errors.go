package router

import (
	"fmt"
	"strconv"

	"github.com/lekuruu/ubisoft-game-service/common"
)

// Extension of the "error" interface, that contains an internal message
// and a response error code, that will be sent to the client
type GSError interface {
	Code() int
	Error() string
	Response(request *GSMessage) *GSMessage
}

type RouterError struct {
	Message      string
	ResponseCode int
}

func (e *RouterError) Error() string {
	return fmt.Sprintf("RouterError: '%s' (%d)", e.Message, e.ResponseCode)
}

func (e *RouterError) Code() int {
	return e.ResponseCode
}

func (e *RouterError) Response(request *GSMessage) *GSMessage {
	return &GSMessage{
		Type:     GSM_GSFAIL,
		Property: request.Property,
		Priority: request.Priority,
		Sender:   request.Receiver,
		Receiver: request.Sender,
		Data: []interface{}{
			common.WriteU8(request.Type),
			[]interface{}{common.WriteU32(e.Code())},
		},
	}
}

type LobbyError struct {
	Message      string
	ResponseCode int
}

func (e *LobbyError) Error() string {
	return fmt.Sprintf("LobbyError: '%s' (%d)", e.Message, e.ResponseCode)
}

func (e *LobbyError) Code() int {
	return e.ResponseCode
}

func (e *LobbyError) Response(request *GSMessage) *GSMessage {
	subTypeString, _ := common.GetStringListItem(request.Data, 0)
	subType, _ := strconv.Atoi(subTypeString)

	return &GSMessage{
		Type:     GSM_LOBBY_MSG,
		Property: request.Property,
		Priority: request.Priority,
		Sender:   request.Receiver,
		Receiver: request.Sender,
		Data: []interface{}{
			strconv.Itoa(GSM_GSFAIL),
			[]interface{}{
				strconv.Itoa(subType),
				[]interface{}{strconv.Itoa(e.Code())},
			},
		},
	}
}

// Errors returned by the login service
const (
	ERRORROUTER_UNKNOWNERROR            = 0
	ERRORROUTER_NOTREGISTERED           = 1
	ERRORROUTER_PASSWORDNOTCORRECT      = 2
	ERRORROUTER_NOTDISCONNECTED         = 3
	ERRORROUTER_ARENANOTAVAILABLE       = 4
	ERRORROUTER_FRIENDSNOTAVAILABLE     = 5
	ERRORROUTER_NAMEALREADYUSED         = 6
	ERRORROUTER_PLAYERNOTCONNECTED      = 7
	ERRORROUTER_PLAYERNOTREGISTERED     = 8
	ERRORROUTER_PLAYERCONNECTED         = 9
	ERRORROUTER_PLAYERALREADYREGISTERED = 10
	ERRORROUTER_CLIENTVERSIONTOOOLD     = 11
	ERRORROUTER_DBINBACKUPMODE          = 12
	ERRORROUTER_DBPROBLEM               = 13
	ERRORROUTER_CLIENTINCOMPATIBLE      = 50
)

// Errors returned by the friends service
const (
	ERRORFRIENDS_FRIENDNOTEXIST         = 14
	ERRORFRIENDS_NOTINARENA             = 15
	ERRORFRIENDS_PLAYERNOTONLINE        = 16
	ERRORFRIENDS_NOTINSESSION           = 17
	ERRORFRIENDS_PLAYERIGNORE           = 18
	ERRORFRIENDS_ALREADYCONNECTED       = 19
	ERRORFRIENDS_NOMOREPLAYERS          = 20
	ERRORFRIENDS_NOPLAYERSCORE          = 47
	ERRORFRIENDS_SEARCHPLAYERFINISHED   = 48
	ERRORFRIENDS_PLAYERSTATUSCOREONLINE = 56
)

// Errors for Secure Accounts
const (
	ERRORSECURE_USERNAMEEXISTS    = 1
	ERRORSECURE_USERNAMEMALFORMED = 2
	ERRORSECURE_USERNAMEFORBIDDEN = 3
	ERRORSECURE_INVALIDACCOUNT    = 4
	ERRORSECURE_USERNAMERESERVED  = 5
	ERRORSECURE_PASSWORDMALFORMED = 11
	ERRORSECURE_PASSWORDFORBIDDEN = 13
	ERRORSECURE_INVALIDPASSWORD   = 15
	ERRORSECURE_DATABASEFAILED    = 100
	ERRORSECURE_BANNEDACCOUNT     = 501
	ERRORSECURE_BLOCKEDACCOUNT    = 502
	ERRORSECURE_LOCKEDACCOUNT     = 512
)

// Errors returned by the lobby service
const (
	ERRORLOBBYSRV_UNKNOWNERROR                    = 0
	ERRORLOBBYSRV_GROUPNOTEXIST                   = 1
	ERRORLOBBYSRV_GAMENOTALLOWED                  = 2
	ERRORLOBBYSRV_SPECTATORNOTALLOWED             = 4
	ERRORLOBBYSRV_NOMOREPLAYERS                   = 5
	ERRORLOBBYSRV_NOMORESPECTATORS                = 6
	ERRORLOBBYSRV_NOMOREMEMBERS                   = 7
	ERRORLOBBYSRV_MEMBERNOTREGISTERED             = 8
	ERRORLOBBYSRV_GAMEINPROGRESS                  = 9
	ERRORLOBBYSRV_WRONGGAMEVERSION                = 10
	ERRORLOBBYSRV_PASSWORDNOTCORRECT              = 11
	ERRORLOBBYSRV_ALREADYINGROUP                  = 12
	ERRORLOBBYSRV_NOTMASTER                       = 13
	ERRORLOBBYSRV_NOTINGROUP                      = 14
	ERRORLOBBYSRV_MINPLAYERSNOTREACH              = 15
	ERRORLOBBYSRV_CONNECTADDCONNECTION            = 16
	ERRORLOBBYSRV_CONNECTSENDLOGINMSG             = 17
	ERRORLOBBYSRV_ERRORLOGINMESSAGE               = 18
	ERRORLOBBYSRV_NOHOSTLOBBYSERVER               = 19
	ERRORLOBBYSRV_LOBBYSRVDISCONNECTED            = 20
	ERRORLOBBYSRV_INVALIDGROUPNAME                = 21
	ERRORLOBBYSRV_INVALIDGAMETYPE                 = 22
	ERRORLOBBYSRV_NOMOREGAMEMODULE                = 23
	ERRORLOBBYSRV_CREATENOTALLOWED                = 24
	ERRORLOBBYSRV_GROUPCLOSE                      = 25
	ERRORLOBBYSRV_WRONGGROUPTYPE                  = 26
	ERRORLOBBYSRV_MEMBERNOTFOUND                  = 27
	ERRORLOBBYSRV_MATCHNOTEXIST                   = 30
	ERRORLOBBYSRV_MATCHNOTFINISHED                = 31
	ERRORLOBBYSRV_GAMENOTINITIATED                = 32
	ERRORLOBBYSRV_BEGINALREADYDONE                = 33
	ERRORLOBBYSRV_MATCHALREADYFINISHEDFORYOU      = 34
	ERRORLOBBYSRV_MATCHSCORESSUBMISSIONEVENTFAIL  = 35
	ERRORLOBBYSRV_MATCHSCORESSUBMISSIONALREDYSENT = 36
	ERRORLOBBYSRV_MATCHRESULTSPROCESSNOTFINISHED  = 37
	ERRORLOBBYSRV_MEMBERBANNED                    = 38
	ERRORLOBBYSRV_PASSPORTFAIL                    = 39
	ERRORLOBBYSRV_NOTCREATOR                      = 40
	ERRORLOBBYSRV_GAMENOTFINISHED                 = 41
	ERRORLOBBYSRV_PASSPORTTIMEOUT                 = 42
	ERRORLOBBYSRV_PASSPORTNOTFOUND                = 43
	ERRORLOBBYSRV_GROUPALREADYEXIST               = 44
)

// Errors returned by the session service
const (
	ERRORARENA_SESSIONEXIST             = 21
	ERRORARENA_GAMENOTALLOWED           = 22
	ERRORARENA_NUMBERPLAYER             = 23
	ERRORARENA_NUMBERSPECTATOR          = 24
	ERRORARENA_VISITORNOTALLOWED        = 25
	ERRORARENA_NOTREGISTERED            = 26
	ERRORARENA_NOMOREPLAYERS            = 27
	ERRORARENA_NOMORESPECTATORS         = 28
	ERRORARENA_PLAYERNOTREGISTERED      = 29
	ERRORARENA_SESSIONNOTAVAILABLE      = 30
	ERRORARENA_SESSIONINPROCESS         = 31
	ERRORARENA_BADGAMEVERSION           = 32
	ERRORARENA_PASSWORDNOTCORRECT       = 33
	ERRORARENA_ALREADYINSESSION         = 34
	ERRORARENA_NOTMASTER                = 35
	ERRORARENA_NOTINSESSION             = 36
	ERRORARENA_MINPLAYERS               = 37
	ERRORARENA_ADMINGAMEDOESNOTEXIST    = 38
	ERRORARENA_ADMINSESSIONDOESNOTEXIST = 39
	ERRORARENA_CONNECTADDCONNECTION     = 40
	ERRORARENA_CONNECTSENDLOGINMSG      = 41
	ERRORARENA_ERRORLOGINMESSAGE        = 42
	ERRORARENA_NOHOSTARENA              = 43
	ERRORARENA_ARENADISCONNECTED        = 44
	ERRORARENA_INVALIDGROUPNAME         = 45
	ERRORARENA_INVALIDGAMETYPE          = 46
	ERRORARENA_NOMOREGAMEMODULE         = 47
	ERRORARENA_PASSPORTLABELNOTFOUND    = 48
	ERRORARENA_PASSPORTFAIL             = 49
	ERRORARENA_CREATENOTALLOWED         = 50
	ERRORARENA_INVALIDSESSIONTYPE       = 51
	ERRORARENA_SESSIONCLOSE             = 52
	ERRORARENA_NOTCREATOR               = 53
	ERRORARENA_DEDICATEDSERVERONLY      = 54
)
