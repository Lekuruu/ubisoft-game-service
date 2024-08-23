
from __future__ import annotations
from typing import Callable

from app.services.router import RouterProtocol
from app.utils.pkc import RsaPublicKey
from app.constants import MessageType
from app.utils.gsm import Message
from app.utils import pkc, gsm

"""Global dictionary to store message handler functions"""
MessageTypeHandlers = {}

def register(type: MessageType):
    def decorator(func: Callable):
        MessageTypeHandlers[type] = func
        return func
    return decorator

@register(MessageType.STILLALIVE)
def still_alive(message: Message, client: RouterProtocol):
    client.send(bytes(message))
