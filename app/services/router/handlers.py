
from __future__ import annotations
from typing import Callable, List, TYPE_CHECKING

from app.utils.pkc import RsaPublicKey
from app.constants import MessageType
from app.utils.gsm import Message
from app.utils import pkc, gsm

import app

if TYPE_CHECKING:
    from app.services.router import RouterProtocol

RouterHandlers = {}
WaitModuleHandlers = {}

def register(
    type: MessageType,
    *handler_dicts
) -> Callable:
    def decorator(func: Callable) -> Callable:
        for handlers in handler_dicts:
            handlers[type] = func
        return func
    return decorator

@register(MessageType.STILLALIVE, RouterHandlers, WaitModuleHandlers)
def still_alive(message: Message, client: RouterProtocol):
    client.send_message(gsm.StillaliveResponse(
        client,
        message.header,
        message.data
    ))

@register(MessageType.KEY_EXCHANGE, RouterHandlers, WaitModuleHandlers)
def key_exchange(message: Message, client: RouterProtocol):
    request_id = message.data.lst[0]

    match request_id:
        case '1':
            client.game_pubkey = RsaPublicKey.from_buf(message.data.lst[1][2]).to_pubkey()
            pub_key, priv_key = pkc.keygen()
            client.sv_pubkey = pub_key
            client.sv_privkey = priv_key
            response = gsm.KeyExchangeResponse(client, message.header, message.data)

        case '2':
            enc_bf_key = bytes(message.data.lst[1][2])
            bf_key = pkc.decrypt(enc_bf_key, client.sv_privkey)
            client.game_bf_key = bf_key
            response = gsm.KeyExchangeResponse(client, message.header, message.data)

        case _:
            raise NotImplementedError(f"Unknown requestId: {request_id}")

    client.send_message(response)

@register(MessageType.LOGIN, RouterHandlers)
def do_login(message: Message, client: RouterProtocol):
    # TODO: Implement actual login
    username = message.data.lst[0]
    password = message.data.lst[1]
    game = message.data.lst[2]

    response = gsm.LoginResponse(client, message.header, message.data)
    client.send_message(response)

@register(MessageType.JOINWAITMODULE, RouterHandlers)
def wm_join_request(message: Message, client: RouterProtocol):
    response = gsm.JoinWaitModuleResponse(
        client,
        message.header,
        message.data,
        (
            app.config["services"]["Router"]["WaitModule"]["IP"],
            app.config["services"]["Router"]["WaitModule"]["Port"]
        )
    )

    client.send_message(response)

@register(MessageType.LOGINWAITMODULE, WaitModuleHandlers)
def wm_login(message: Message, client: RouterProtocol):
    # TODO: Implement actual login
    username = message.data.lst[0]

    response = gsm.LoginWaitModuleResponse(client, message.header, message.data)
    client.send_message(response)
