from aio_pika.message import IncomingMessage
from typing import Any, Coroutine, Callable

ConsumerCallback = Callable[[IncomingMessage], Coroutine[Any, Any, None]]