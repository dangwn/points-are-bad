import asyncio
import logging
import signal

from .parse_request import parse_raw_request, RequestDict
from .callback import run_callbacks

from typing import (
    Callable,
    List,
    Optional, 
    Tuple
)
from types import FrameType
from .custom_types import (
    EventType,
    DecoratedCallback
)

HANDLED_SIGNALS: Tuple[signal.Signals] = (
    signal.SIGINT,  # Unix signal 2. Sent by Ctrl+C.
    signal.SIGTERM  # Unix signal 15. Sent by `kill <pid>`.
)

logging.basicConfig(
    level=logging.INFO,
    format='%(levelname)s - %(asctime)s - %(message)s',
    datefmt='%H:%M:%S'
)
clientLogger: logging.Logger = logging.getLogger('AsyncClient')

class AsyncServer:
    def __init__(
        self,
        host: str,
        port: int = 8008
    ) -> None:
        self.host = host
        self.port = port
        self.server: Optional[asyncio.AbstractServer] = None

        self.should_exit: bool = False
        self.force_exit: bool = False

        self.startup_callbacks: List[Callable[..., None]] = []
        self.loop_callbacks: List[Callable[..., None]] = []
        self.shutdown_callbacks: List[Callable[..., None]] = []

    @property    
    def server_is_initialized(self) -> bool:
        return isinstance(self.server, asyncio.AbstractServer)
    
    def _add_on_startup(self, callback: Callable[..., None]) -> None:
        self.startup_callbacks.append(callback)

    def _add_on_loop(self, callback: Callable[..., None]) -> None:
        self.loop_callbacks.append(callback)

    def _add_on_shutdown(self, callback: Callable[..., None]) -> None:
        self.shutdown_callbacks.append(callback)

    async def handle_request(
        self,
        reader: asyncio.StreamReader,
        writer: asyncio.StreamWriter
    ) -> None:
        request: bytes = b''
        while not request.endswith(b'\r\n\r\n'):
            chunk: bytes = await reader.read(1024)
            request += chunk

        request_dict: RequestDict = parse_raw_request(request)
        response: str = ''
        response_code: str = ''
        if request_dict['method'] == 'GET' and request_dict['path'] == '/':
            response_code = '200 OK'
            response = f'HTTP/1.1 {response_code}\r\nContent-Type: text/plain\r\n\r\nHealthy'
        else:
            response_code = '404 Not Found'
            response = f'HTTP/1.1 {response_code}\r\nContent-Type: text/plain\r\n\r\n404 Not Found'

        logging.info(f'{request_dict["method"]} {request_dict["path"]} - {response_code}')

        writer.write(response.encode())
        await writer.drain()

        writer.close()
    
    def run(self) -> None:
        return asyncio.run(self.serve())
    
    async def serve(self) -> None:
        self.install_signal_handlers()

        await self.startup()
        if not self.server_is_initialized:
            return

        while not self.should_exit and not self.force_exit:
            await run_callbacks(self.loop_callbacks)
            await asyncio.sleep(0.1)
            
        if not self.force_exit:
            await self.shutdown()

    async def startup(self) -> None:
        clientLogger.info('Starting server...')
        self.server = await asyncio.start_server(
            client_connected_cb=self.handle_request,
            host=self.host,
            port=self.port
        )
        await self.server.start_serving()

        await run_callbacks(self.startup_callbacks)

        clientLogger.info('Server started successfully!')

    async def shutdown(self) -> None:
        clientLogger.info('Shutting down server...')
        
        await run_callbacks(self.shutdown_callbacks)
        
        if self.server_is_initialized:
            server_loop = self.server.get_loop()
            self.server.close()

            if not server_loop.is_closed():
                await self.server.wait_closed()

            await asyncio.sleep(0.1)
            self.server = None
        
        clientLogger.info('Server shut down!')
    
    #=====================================================================
    # Pinched from uvicorn
    # Allows for keyboard interupts without leaving subprocesses hanging
    def install_signal_handlers(self) -> None:
        loop: asyncio.AbstractEventLoop = asyncio.get_event_loop()

        try:
            for sig in HANDLED_SIGNALS:
                loop.add_signal_handler(sig, self.handle_exit, sig, None)
        except NotImplementedError: 
            # Windows
            for sig in HANDLED_SIGNALS:
                signal.signal(sig, self.handle_exit)

    def handle_exit(self, sig: int, frame: Optional[FrameType]) -> None:
        if self.should_exit and sig == signal.SIGINT:
            self.force_exit = True
        else:
            self.should_exit = True
    #====================================================================

    def on_event(
        self, 
        event_type: EventType
    ) -> Callable[[DecoratedCallback], DecoratedCallback]:
        def decorator(func: Callable[..., None]) -> Callable[..., None]:
            if event_type == 'startup':
                self._add_on_startup(func)
            elif event_type == 'loop':
                self._add_on_loop(func)
            elif event_type == 'shutdown':
                self._add_on_shutdown(func)
            else:
                raise ValueError('Event type not in specified event types: "startup", "loop", "shutdown".')
            return func
        return decorator
    