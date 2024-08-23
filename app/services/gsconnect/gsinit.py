
from twisted.web.resource import Resource
from twisted.web.server import Request

import logging
import app

class GameServiceInit(Resource):
    isLeaf = True

    def __init__(self):
        self.children = {}
        self.logger = logging.getLogger(__name__)

    @staticmethod
    def get_argument(name: str, request: Request) -> str:
        return request.args.get(name.encode(), [b''])[0].decode()

    def render_GET(self, request: Request) -> bytes:
        user = self.get_argument('user', request) or 'Anonymous'
        product = self.get_argument('dp', request)
        games = app.config['gsconnect']['Games']

        if not (config_path := games.get(product)):
            self.logger.warning(f'Unsupported product: "{product}"')
            return b''

        with open(config_path, 'r') as file:
            config = file.read()

        self.logger.info(f'"{user}" is connecting to "{product}"')
        return config.encode()
