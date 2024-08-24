
from twisted.web.server import Request, NOT_DONE_YET
from twisted.web.resource import Resource

import logging
import app

class GameServiceInit(Resource):
    isLeaf = True

    def __init__(self):
        self.children = {}
        self.logger = logging.getLogger('GSConnect')

    @staticmethod
    def argument(name: str, request: Request) -> str:
        return request.args.get(name.encode(), [b''])[0].decode()

    @staticmethod
    def write_and_disconnect(request: Request, response: bytes):
        request.write(response)
        request.finish()
        request.transport.loseConnection()

    def game_config(self, product: str, user: str) -> str:
        games = app.config['services']['GSConnect']['Games']

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
            # Load game config
            config = self.game_config(product, user)
        except Exception as e:
            self.logger.error(f'Failed to load game config: {e}')
            return b''

        # Write the response, and immediately disconnect
        self.write_and_disconnect(request, config.encode())
        return NOT_DONE_YET
