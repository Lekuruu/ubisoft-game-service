
from __future__ import annotations

from twisted.internet.protocol import Factory
from app.services.tcp import BaseTcpProtocol
from twisted.web.resource import Resource
from twisted.internet import reactor
from twisted.web.server import Site

import logging

class Server(Factory):
    def __init__(
        self,
        name: str,
        port: int,
        protocol: BaseTcpProtocol
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

class HttpServer(Site):
    def __init__(
        self,
        name: str,
        port: int,
        root: Resource
    ) -> None:
        super().__init__(root)
        self.name = name
        self.port = port
        self.logger = logging.getLogger(name)

    def start(self) -> None:
        reactor.listenTCP(self.port, self)

    def startFactory(self):
        self.logger.info(f'Listening on {self.port}...')

    def stopFactory(self):
        self.logger.warning(f'Stopping server...')
