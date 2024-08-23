
from twisted.web.resource import Resource
from twisted.web.server import Request

import logging

class GameServiceInit(Resource):
    isLeaf = True

    # TODO: Move this to seperate config file
    Config = {
        'SPLINTERCELL3COOP': {
            'Router': [{
                'IP': '127.0.0.1',
                'Port': 40000
            }],
            'IRC': [{
                'IP': '127.0.0.1',
                'Port': 6668
            }],
            'CDKeyServer': [{
                'IP': '127.0.0.1',
                'Port': 44000
            }],
            'Proxy': [{
                'IP': '127.0.0.1',
                'Port': 4040
            }],
            'NATServer': [{
                'IP': '127.0.0.1',
                'Port': 7781
            }]
        }
    }

    def __init__(self):
        self.children = {}
        self.logger = logging.getLogger(__name__)

    def render_GET(self, request: Request) -> bytes:
        user = self.get_argument('user', request) or 'Anonymous'
        product = self.get_argument('dp', request)

        if not (config := self.Config.get(product)):
            self.logger.warning(f'Unsupported product: "{product}"')
            return b''

        self.logger.info(f'"{user}" is connecting to "{product}"')
        return self.config_to_ini(config).encode()

    @staticmethod
    def get_argument(name: str, request: Request) -> str:
        return request.args.get(name.encode(), [b''])[0].decode()

    @staticmethod
    def config_to_ini(config: dict, name: str = 'Servers') -> str:
        result = (
            [f'[{name}]']
            if name else []
        )

        servers = {
            f'{service_name}{key}{index}': value
            for service_name, servers in config.items()
            for index, server in enumerate(servers)
            for key, value in server.items()
        }

        for key, value in servers.items():
            if type(value) != list:
                result.append(f'{key}={value}')
                continue

            for item in value:
                result.append(f'{key}={item}')

        return '\n'.join(result)
