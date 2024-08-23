
from __future__ import annotations

from app.services.protocol import BaseProtocol
from twisted.internet.protocol import Factory
from twisted.internet import reactor

import logging

class Server(Factory):
    def __init__(
        self,
        name: str,
        port: int,
        protocol: BaseProtocol
    ) -> None:
        self.name = name
        self.port = port
        self.protocol = protocol
        self.logger = logging.getLogger(name)

    def start(self) -> None:
        reactor.listenTCP(self.port, self)

    def startFactory(self):
        self.logger.info(f'Listening on {self.port}...')

    def stopFactory(self):
        self.logger.warning(f'Stopping server...')

    def buildProtocol(self, address):
        return self.protocol(address)
