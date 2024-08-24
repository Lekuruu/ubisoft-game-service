
from app.services.router import RouterProtocol, WaitModuleProtocol
from app.services.server import Server, HttpServer
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
    Services = app.config['services']
    Server('WaitModule', Services['Router']["WaitModule"]['Port'], WaitModuleProtocol).start()
    Server('Router', Services['Router']['Port'], RouterProtocol).start()
    HttpServer('GSConnect', Services['GSConnect']['Port'], GSConnect()).start()
    reactor.run()

if __name__ == '__main__':
    main()
