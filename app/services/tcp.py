
from __future__ import annotations

from twisted.internet.address import IPv4Address, IPv6Address
from twisted.internet.protocol import Protocol
from twisted.internet.error import ConnectionDone
from twisted.python.failure import Failure

import logging

IPAddress = IPv4Address | IPv6Address

class BaseTcpProtocol(Protocol):
    """
    Base protocol class that includes basic logging, as well as
    providing functions that follow conventional naming schemes.
    """

    def __init__(self, address: IPAddress) -> None:
        self.host = address.host
        self.port = address.port

        self.logger = logging.getLogger(self.host)
        self.buffer = bytearray()
        self.busy = False

    def on_connect(self) -> None:
        ...

    def on_disconnect(self, reason: Failure) -> None:
        ...

    def on_data(self, data: bytes) -> None:
        ...

    def send(self, data: bytes):
        try:
            self.transport.write(data)
        except Exception as e:
            self.logger.error(
                f'Could not write to transport layer: {e}',
                exc_info=True
            )

    def disconnect(self) -> None:
        self.logger.debug(f'-> Closing connection...')
        self.transport.loseConnection()
        self.on_disconnect(ConnectionDone())

    def connectionMade(self):
        self.logger.debug(f'-> <{self.host}:{self.port}>')
        self.on_connect()

    def connectionLost(self, reason: Failure = ConnectionDone()) -> None:
        if reason.type is ConnectionDone:
            self.logger.debug(f'-> Connection done.')
            self.on_disconnect(reason)
            return

        self.logger.warning(f'-> Lost connection "{reason.getErrorMessage()}".')
        self.on_disconnect(reason)

    def dataReceived(self, data: bytes) -> None:
        if self.busy:
            self.buffer += data
            return

        try:
            self.busy = True
            self.buffer += data
            self.on_data(self.buffer)
        except Exception as e:
            self.logger.error(
                f'Error while processing data: {e}',
                exc_info=True
            )
        finally:
            self.busy = False
