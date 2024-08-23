
from app.constants import MessageType
from typing import Callable

# Global dictionary to store message handler functions
MessageTypeHandlers = {}

def register(type: MessageType):
    def decorator(func: Callable):
        MessageTypeHandlers[type] = func
        return func
    return decorator
