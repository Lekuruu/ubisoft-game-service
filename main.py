
from app.logging import ConsoleLogger

import logging

logging.basicConfig(
    level=logging.DEBUG,
    handlers=[ConsoleLogger]
)
