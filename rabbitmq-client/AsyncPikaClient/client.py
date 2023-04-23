import aio_pika

from exceptions import ClientNotInitializedException

from aio_pika.connection import Connection
from aio_pika.channel import Channel
from aio_pika.queue import Queue
from custom_types import ConsumerCallback
from typing import Optional, Dict

class AsyncPikaClient:
    def __init__(
        self, 
        connection_string: str,
        queue_name: str,
    ):
        self.connection_string: str = connection_string
        self.queue_name: str = queue_name

        self.connection: Optional[Connection] = None
        self.channel: Optional[Channel] = None
        self.queue: Optional[Queue] = None
        self.consumer_callback: Optional[ConsumerCallback] = None

    @classmethod
    async def startup(
        cls,
        connection_string: str, 
        queue_name: str,
        consumer_callback: ConsumerCallback
    ):
        self = AsyncPikaClient(
            connection_string=connection_string,
            queue_name=queue_name
        )
        self.connection = await aio_pika.connect_robust(
            url=connection_string
        )
        self.channel = await self.connection.channel()
        self.queue = await self.channel.declare_queue(
            queue_name,
            durable=True
        )
        self.consumer_callback = consumer_callback

        return self

    async def shutdown(self) -> None:
        if (self.connection):
            await self.channel.close()
            await self.connection.close()

    async def send_message(self, msg: str) -> Dict[str,str]:
        if not (self.channel):
            raise ClientNotInitializedException('Channel not initialized')
        await self.channel.default_exchange.publish(
            aio_pika.Message(
                body=msg.encode()
            ),
            routing_key=self.queue_name
        )
        return {'status':'success'}
    
    async def consume_messages(self) -> None:
        if not self.queue:
            raise ClientNotInitializedException('Queue not initialized')
        
        await self.queue.consume(self.consumer_callback, no_ack=False)