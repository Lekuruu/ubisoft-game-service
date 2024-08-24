
from __future__ import annotations

from twisted.internet.protocol import Protocol, DatagramProtocol
from twisted.web.resource import Resource
from twisted.internet import reactor

from app.logging import ConsoleLogger
from app.services import (
    RouterProtocol,
    CDKeyProtocol,
    HttpServer,
    TcpServer,
    GSConnect
)

import logging
import app

logging.basicConfig(
    level=logging.DEBUG,
    handlers=[ConsoleLogger]
)

def listen_tcp(port: int, protocol: Protocol) -> None:
    TcpServer(
        protocol.__name__,
        port,
        protocol
    ).start()

def listen_http(port: int, resource: Resource) -> None:
    HttpServer(
        resource.__name__,
        port,
        resource()
    ).start()

def listen_udp(port: int, protocol: DatagramProtocol) -> None:
    reactor.listenUDP(
        port,
        protocol()
    )

def main() -> None:
    Services = app.config['services']
    listen_http(Services['GSConnect']['Port'], GSConnect)
    listen_tcp(Services['Router']["WaitModule"]['Port'], RouterProtocol)
    listen_tcp(Services['Router']['Port'], RouterProtocol)
    listen_udp(Services['CDKey']['Port'], CDKeyProtocol)
    reactor.run()

if __name__ == '__main__':
    main()
