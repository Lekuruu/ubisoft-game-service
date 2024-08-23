
from app.services.server import Server, HttpServer
from app.services.router import RouterProtocol
from app.services.gsconnect import GSConnect
from app.logging import ConsoleLogger
from twisted.internet import reactor

import logging
import app

logging.basicConfig(
    level=logging.DEBUG,
    handlers=[ConsoleLogger]
)

def main():
    HttpServer('GSConnect', app.config['services']['GSConnect']['Port'], GSConnect()).start()
    Server('Router', app.config['services']['Router']['Port'], RouterProtocol).start()
    reactor.run()

if __name__ == '__main__':
    main()
