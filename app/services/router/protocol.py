
from __future__ import annotations
from rsa import PublicKey, PrivateKey
from typing import Callable, Dict

from app.constants import MessageType, GSMSG_HEADER_SIZE
from app.services.tcp import BaseTcpProtocol, IPAddress
from app.utils.gsm import Message

from .handlers import RouterHandlers

class RouterProtocol(BaseTcpProtocol):
    Handlers: Dict[MessageType, Callable] = RouterHandlers

    def __init__(self, address: IPAddress) -> None:
        super().__init__(address)
        self.game_pubkey: PublicKey | None = None
        self.sv_privkey: PrivateKey | None = None
        self.sv_pubkey: PublicKey | None = None
        self.game_bf_key: bytes | None = None
        self.sv_bf_key: bytes | None = None

    def send_message(self, msg: Message) -> None:
        self.logger.debug(f'<- {msg}')
        self.send(bytes(msg))

    def on_data(self, data: bytes) -> None:
        while self.buffer:
            # Parse message header
            msg = Message.from_bytes(
                self.buffer,
                self.game_bf_key
            )

            # Handle message
            self.handle_message(msg)

    def handle_message(self, msg: Message) -> None:
        self.logger.debug(f'-> {msg}')

        # Reset packet buffer
        self.buffer = self.buffer[msg.header.size:]

        if not (handler := self.Handlers.get(msg.header.type)):
            self.logger.warning(f'Unsupported message type: "{msg.header.type.name}"')
            return

        return handler(msg, self)
