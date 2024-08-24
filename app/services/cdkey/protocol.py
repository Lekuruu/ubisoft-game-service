
from __future__ import annotations
from app.services.udp import BaseUdpProtocol, IPAddress
from app.utils.cdkm import CDKeyMessage

from .handlers import CDKeyHandlers

class CDKeyProtocol(BaseUdpProtocol):
    def on_data(self, data: bytes, address: IPAddress) -> None:
        request = CDKeyMessage.from_buffer(data)
        handler = CDKeyHandlers.get(request.request_type)

        if not handler:
            self.logger.warning(f"Unsupported request type: {request.request_type.name}")
            return

        if response := handler(request):
            self.transport.write(bytes(response), address)
