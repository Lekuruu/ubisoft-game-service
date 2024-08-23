
from __future__ import annotations
from rsa import PublicKey, PrivateKey

from app.services.protocol import BaseProtocol, IPAddress
from app.utils.gsm import Message, GSMessageBundle
from app.constants import GSMSG_HEADER_SIZE

from .handlers import MessageTypeHandlers

class RouterProtocol(BaseProtocol):
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
        if len(data) < GSMSG_HEADER_SIZE:
            # Wait for next buffer
            return

        # Parse message header
        msg = Message.from_bytes(data, self.game_bf_key)
        self.logger.debug(f'-> {msg}')

        if msg.header.size >= len(data):
            self.handle_message(msg)
            return

    def handle_message(self, msg: Message) -> None:
        # Reset packet buffer
        self.buffer = self.buffer[msg.header.size:]

        if not (handler := MessageTypeHandlers.get(msg.header.type)):
            self.logger.warning(f'Unsupported message type: "{msg.header.type.name}"')
            return

        return handler(self, msg)

    def handle_message_bundle(self, msg: Message) -> None:
        bundle = GSMessageBundle(
            msg,
            self.buffer,
            self.game_bf_key
        )

        for msg in bundle.messages:
            self.handle_message(msg)
