
from twisted.web.resource import Resource
from twisted.web.server import Request

from .gsinit import GameServiceInit

class Root(Resource):
    isLeaf = True

    def render_GET(self, request: Request) -> bytes:
        return (
            b"<html><head>"
            b"<title>connect</title>"
            b"</head></html>"
        )

class GSConnect(Resource):
    isLeaf = False

    def __init__(self):
        self.children = {}
        self.putChild(b'gsinit.php', GameServiceInit())
        self.putChild(b'', Root())
