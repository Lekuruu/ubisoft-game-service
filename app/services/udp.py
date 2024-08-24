
from __future__ import annotations

from twisted.internet.address import IPv4Address, IPv6Address
from twisted.internet.protocol import DatagramProtocol

import logging

IPAddress = IPv4Address | IPv6Address

class BaseUdpProtocol(DatagramProtocol):
    def __init__(self) -> None:
        self.logger = logging.getLogger(self.logPrefix())

    def on_data(self, data: bytes, address: IPAddress) -> None:
        ...

    def startProtocol(self) -> None:
        address = self.transport.getHost()
        self.logger.info(f'Listening on {address.port}...')

    def stopProtocol(self) -> None:
        self.logger.info('Stopping server...')

    def datagramReceived(self, datagram: bytes, address: IPAddress) -> None:
        self.on_data(datagram, address)
