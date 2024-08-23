
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
    def argument(name: str, request: Request) -> str:
        return request.args.get(name.encode(), [b''])[0].decode()

    def game_config(self, product: str, user: str) -> str:
        games = app.config['gsconnect']['Games']

        if not (config_path := games.get(product)):
            self.logger.warning(f'Unsupported product: "{product}"')
            return ''

        with open(config_path, 'r') as file:
            config = file.read()

        self.logger.info(f'"{user}" is connecting to "{product}"')
        return config

    def render_GET(self, request: Request):
        user = self.argument('user', request) or 'Anonymous'
        product = self.argument('dp', request)

        try:
            config = self.game_config(product, user)
            return config.encode()
        except Exception as e:
            self.logger.error(f'Failed to load game config: {e}')
            return b''
