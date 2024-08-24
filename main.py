
from app.services.router import RouterProtocol
from app.services.gsconnect import GSConnect
from app.services.cdkey import CDKeyProtocol
from app.services.http import HttpServer
from app.services.tcp import TcpServer
from app.logging import ConsoleLogger
from twisted.internet import reactor

import logging
import app

logging.basicConfig(
    level=logging.DEBUG,
    handlers=[ConsoleLogger]
)

def main() -> None:
    Services = app.config['services']
    HttpServer('GSConnect', Services['GSConnect']['Port'], GSConnect()).start()
    TcpServer('WaitModule', Services['Router']["WaitModule"]['Port'], RouterProtocol).start()
    TcpServer('Router', Services['Router']['Port'], RouterProtocol).start()
    TcpServer('CDKey', Services['CDKey']['Port'], CDKeyProtocol).start()
    reactor.run()

if __name__ == '__main__':
    main()
