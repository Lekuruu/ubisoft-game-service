
from __future__ import annotations
from typing import Callable
from app.utils.cdkm import (
    ValidationResponse,
    ActivationResponse,
    ChallengeResponse,
    AuthResponse,
    CDKeyMessage,
    RequestType,
    Response
)

CDKeyHandlers = {}

def register(type: RequestType) -> Callable:
    def decorator(func: Callable) -> Callable:
        CDKeyHandlers[type] = func
        return func
    return decorator

@register(RequestType.CHALLENGE)
def handle_challenge(request: CDKeyMessage) -> Response:
    return ChallengeResponse.from_request(request)

@register(RequestType.ACTIVATION)
def handle_activation(request: CDKeyMessage) -> Response:
    return ActivationResponse.from_request(request)

@register(RequestType.AUTH)
def handle_auth(request: CDKeyMessage) -> Response:
    return AuthResponse.from_request(request)

@register(RequestType.VALIDATION)
def handle_validation(request: CDKeyMessage) -> Response:
    return ValidationResponse.from_request(request)

@register(RequestType.STILL_ALIVE)
def handle_still_alive(request: CDKeyMessage) -> Response:
    pass
