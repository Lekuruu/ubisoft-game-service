
from app.services.server import Server, HttpServer
from app.services.router import RouterProtocol
from app.services.gsconnect import GSConnect
from app.logging import ConsoleLogger
from twisted.internet import reactor

import logging

logging.basicConfig(
    level=logging.DEBUG,
    handlers=[ConsoleLogger]
)

def main():
    HttpServer('GSConnect', 80, GSConnect()).start()
    Server('Router', 40000, RouterProtocol).start()
    reactor.run()

if __name__ == '__main__':
    main()
