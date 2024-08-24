
from __future__ import annotations

from app.utils.gsm import MessageType
from app.utils.blowfish import Cipher
from app.utils.data import List, Bin
from app.utils import serialization

from dataclasses import dataclass
from enum import IntEnum

CDKM_HEADER_SIZE = 5
"""Length of `CDKeyMessage` header in bytes."""

class CDKeyPlayerStatus(IntEnum):
    """Player status."""
    E_PLAYER_UNKNOWN = 0
    E_PLAYER_INVALID = 1
    E_PLAYER_VALID = 2

class RequestType(IntEnum):
    """CD-Key service requests."""
    CHALLENGE = 1
    ACTIVATION = 2
    AUTH = 3
    VALIDATION = 4
    PLAYER_STATUS = 5
    DISCONNECT_USER = 6
    STILL_ALIVE = 7

BLOWFISH = Cipher(bytes("SKJDHF$0maoijfn4i8$aJdnv1jaldifar93-AS_dfo;hjhC4jhflasnF3fnd", 'utf8'))

@dataclass
class CDKeyMessage:
    type: int
    size: int
    data: List
    message_id: int
    request_type: RequestType
    unknown: int = 0
    inner_data: List = List()

    @classmethod
    def from_buffer(cls, bts: bytes) -> "CDKeyMessage":
        data = List.from_buf(bytearray(BLOWFISH.decrypt(bts[CDKM_HEADER_SIZE:])))
        request_type = RequestType(int(data.lst[1]))

        return cls(
            type=bts[0],
            size=serialization.read_u32_be(bts[1:CDKM_HEADER_SIZE]),
            data=data,
            message_id=int(data.lst[0]),
            request_type=request_type,
            unknown=int(data.lst[2]) if request_type != RequestType.STILL_ALIVE else 0,
            inner_data=data.lst[3] if request_type != RequestType.STILL_ALIVE else List()
        )

@dataclass
class Response:
    """Base class for CDKM responses."""
    type: int
    request_type: RequestType
    size: int
    data: List

    def __bytes__(self):
        """Serializes the response into a CDKeyMessage buffer."""
        buf = bytearray()
        buf.append(self.type)
        data = bytearray(bytes(self.data))
        data.pop(0)
        data.pop()
        data = BLOWFISH.encrypt(bytes(data))
        self.size = len(data)
        buf.extend(serialization.write_u32_be(self.size))
        buf.extend(data)
        return bytes(buf)

    @classmethod
    def from_request(cls, req: CDKeyMessage) -> "Response":
        return cls(
            req.type,
            req.request_type,
            size=0,
            data=List([
                str(req.message_id),
                str(req.request_type.value),
                str(req.unknown),
                []
            ])
        )

@dataclass
class ChallengeResponse(Response):
    def __post_init__(self):
        msg_type = MessageType.GSSUCCESS
        hash = b'\x00\x11\x22\x33\x44\x55\x66\x77\x88\x99\xaa\xbb\xcc\xdd\xee\xff\x01\x02\x03\x04'
        res_data = [bytes(Bin(hash))]
        self.data.lst[3].append(str(msg_type.value))
        self.data.lst[3].append(res_data)

@dataclass
class ActivationResponse(Response):
    def __post_init__(self):
        msg_type = MessageType.GSSUCCESS
        activation_id = b'\x33\x33\x33\x33\x33\x33\x33\x33\x33\x33\x33'
        buf2 = b'\x44\x44\x44\x44\x44\x44\x44\x44\x44\x44\x44'
        res_data = [bytes(Bin(activation_id)), bytes(Bin(buf2))]
        self.data.lst[3].append(str(msg_type.value))
        self.data.lst[3].append(res_data)

@dataclass
class AuthResponse(Response):
    def __post_init__(self):
        msg_type = MessageType.GSSUCCESS
        auth_id = b'\x55\x55\x55\x55\x55\x55\x55\x55\x55\x55\x55\x55\x55\x55\x55'
        res_data = [bytes(Bin(auth_id))]
        self.data.lst[3].append(str(msg_type.value))
        self.data.lst[3].append(res_data)

@dataclass
class ValidationResponse(Response):
    def __post_init__(self):
        msg_type = MessageType.GSSUCCESS
        status = CDKeyPlayerStatus.E_PLAYER_VALID
        buf = b'\x66\x66\x66\x66\x66\x66\x66\x66\x66\x66\x66'
        res_data = [str(status.value), bytes(Bin(buf))]
        self.data.lst[3].append(str(msg_type.value))
        self.data.lst[3].append(res_data)
