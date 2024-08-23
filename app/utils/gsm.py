
from __future__ import annotations
from typing import List as TypedList
from dataclasses import dataclass

from app.utils import serialization, gsxor, pkc
from app.utils.pkc import RsaPublicKey
from app.utils.blowfish import Cipher
from app.utils.data import List
from app.services import router
from app.constants import (
    GSMSG_HEADER_SIZE,
    MessageProperty,
    MessageTarget,
    MessageType
)

@dataclass
class GSMessageHeader:
    """Header for `GSMessage` and `GSEncryptMessage`"""
    size: int
    property: MessageProperty
    priority: int
    type: MessageType
    sender: MessageTarget
    receiver: MessageTarget

    @classmethod
    def from_bytes(cls, bts: bytes):
        return cls(
            (bts[0] << 16) + (bts[1] << 8) + bts[2],
            MessageProperty(bts[3] >> 6),
            bts[3] & 0x3F,
            MessageType(bts[4]),
            MessageTarget(bts[5] >> 4),
            MessageTarget(bts[5] & 0x0F)
        )

    def __bytes__(self):
        result = bytearray(GSMSG_HEADER_SIZE)
        size = serialization.write_u24_be(self.size)
        result[0] = size[0]
        result[1] = size[1]
        result[2] = size[2]
        result[3] &= 0x1F
        result[3] |= self.property.value << 6
        result[3] |= self.priority & 0x20
        result[4] = self.type.value
        result[5] &= 0xF
        result[5] |= 0x10 * self.sender.value
        result[5] &= 0xF0
        result[5] |= self.receiver.value & 0xF
        return bytes(result)

@dataclass
class Message:
    """Common message implementation"""
    header: GSMessageHeader
    data: List | None = None

    @classmethod
    def from_bytes(cls, bts: bytes, blowfish_key: bytes):
        header = GSMessageHeader(bts[:GSMSG_HEADER_SIZE])
        data = None

        match header.property:
            case MessageProperty.GS:
                if header.size > GSMSG_HEADER_SIZE:
                    dec = gsxor.decrypt(bts[GSMSG_HEADER_SIZE:header.size])
                    data: List = List.from_buf(bytearray(dec))

            case MessageProperty.GS_ENCRYPT:
                dec = Cipher(blowfish_key).decrypt(bts[GSMSG_HEADER_SIZE:header.size])
                data: List = List.from_buf(bytearray(dec))

            case MessageProperty.GAME:
                pass

        return cls(header, data)

@dataclass
class GSMessageBundle:
    """Packet containing 2 or more GS messages"""
    messages: TypedList[Message]

    def from_bytes(cls, first: Message, bts: bytes, blowfish_key: bytes):
        messages = [first]

        while len(bts) > 0:
            msg = Message.from_bytes(bts, blowfish_key)
            messages.append(msg)
            bts = bts[msg.header.size:]

        return cls(messages)

@dataclass
class GSMResponse:
    """Base class for GS message responses"""
    header: GSMessageHeader
    client: router.RouterProtocol
    data: List | None = None

    def __bytes__(self):
        if self.data is None:
            return bytes(self.header)

        bts = bytearray()
        data = bytearray(bytes(self.data))
        data.pop(0)
        data.pop()

        match self.header.property:
            case MessageProperty.GS:
                dl = gsxor.encrypt(bytes(dl))
                self.header.size = GSMSG_HEADER_SIZE + len(dl)
            case MessageProperty.GS_ENCRYPT:
                raise NotImplementedError("GS_ENCRYPT message serialization unsupported.")

        bts += bytes(self.header)
        bts += bytes(data)
        return bytes(bts)

@dataclass
class KeyExchangeResponse(GSMResponse):
    """Response to `KEY_EXCHANGE` messages"""
    def __post_init__(self):
        assert self.blowfish_key is not None
        assert self.header.type == MessageType.KEY_EXCHANGE
        request_id = int(self.data.lst[0])

        match request_id:
            case 1:
                self.data = List(['1', ['1']])
                pub_key: RsaPublicKey = RsaPublicKey.from_pubkey(self.client.sv_pubkey)
                buf = bytes(pub_key)
                self.data.lst[1].append(str(len(buf)))
                self.data.lst[1].append(buf)

            case 2:
                self.data = List(['2', ['1']])
                bf_key = Cipher.keygen(16)
                self.client.sv_bf_key = bf_key
                enc_key = pkc.encrypt(bf_key, self.client.game_pubkey)
                self.data.lst[1].append(str(len(enc_key)))
                self.data.lst[1].append(enc_key)

            case 3:
                raise NotImplementedError("KEY_EXCHANGE disconnections are not implemented.")

            case _:
                raise BufferError(f"KEY_EXCHANGE request with id={request_id}.")

class LoginResponse(GSMResponse):
    """Response to `LOGIN` messages"""
    def __init__(self, req: Message):
        assert req.header.type == MessageType.LOGIN
        super().__init__(req)
        self.header.property = MessageProperty.GS
        self.header.type = MessageType.GSSUCCESS
        msg_id = MessageType.LOGIN.value
        self.data = List([msg_id.to_bytes(1, 'little')])

class JoinWaitModuleResponse(GSMResponse):
    """Response to `JOINWAITMODULE` messages"""
    def __init__(self, req: Message, wait_module: tuple[str, int]):
        assert req.header.type == MessageType.JOINWAITMODULE
        super().__init__(req)
        self.header.property = MessageProperty.GS
        self.header.type = MessageType.GSSUCCESS
        msg_id = MessageType.JOINWAITMODULE.value
        self.data = List([
            msg_id.to_bytes(1, 'little'),
            [wait_module[0], serialization.write_u32(wait_module[1])]
        ])

class LoginWaitModuleResponse(GSMResponse):
    """Response to `LOGINWAITMODULE` messages"""
    def __init__(self, req: Message):
        assert req.header.type == MessageType.LOGINWAITMODULE
        super().__init__(req)
        self.header.property = MessageProperty.GS
        self.header.type = MessageType.GSSUCCESS
        msg_id = MessageType.LOGINWAITMODULE.value
        self.data = List([msg_id.to_bytes(1, 'little')])

class PlayerInfoResponse(GSMResponse):
    """Response to `PLAYERINFO` messages"""
    def __init__(self, req: Message):
        assert req.header.type == MessageType.PLAYERINFO
        super().__init__(req)
        self.header.property = MessageProperty.GS
        self.header.type = MessageType.GSSUCCESS
        msg_id = MessageType.PLAYERINFO.value
        player_data = ['findme1', 'findme2', 'findme3', 'findme4', 'findme5', 'findme6', 'findme7']
        self.data = List([msg_id.to_bytes(1, 'little'), player_data])