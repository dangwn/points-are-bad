import aio_pika

from exceptions import ClientNotInitializedException, QueueNotFoundError

from aio_pika.abc import (
    AbstractRobustConnection,
    AbstractRobustChannel,
    AbstractQueue,
    AbstractExchange
)
from custom_types import ConsumerCallback
from typing import Optional, Dict, Union, List

class AsyncClient:
    def __init__(
        self,
        queue_names: Union[List[str], str],
    ):
        self.queue_names: Union[List[str], str] = queue_names

        self.connection: Optional[AbstractRobustConnection] = None
        self.channel: Optional[AbstractRobustChannel] = None

    async def shutdown(self) -> None:
        if (self.connection):
            await self.channel.close()
            await self.connection.close()


class AsyncConsumer(AsyncClient):
    def __init__(self, queue_name: str):
        super().__init__(queue_names=queue_name)

        self.queue: Optional[AbstractQueue] = None
        self.consumer_callback: Optional[ConsumerCallback] = None

    @classmethod
    async def startup(
        cls,
        queue_name: str,
        connection_string: str,
        consumer_callback: ConsumerCallback
    ):
        self = AsyncConsumer(
            queue_name=queue_name,
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
    
    async def consume_messages(self) -> None:
        if not self.queue:
            raise ClientNotInitializedException('Queue not initialized')
        
        await self.queue.consume(self.consumer_callback, no_ack=False)


class AsyncProducer(AsyncClient):
    def __init__(self, queue_names: Union[List[str], str]):
        super().__init__(queue_names=([queue_names] if isinstance(queue_names, str) else queue_names))

        self.exchange: Optional[AbstractExchange] = None
        self.queues: Optional[Dict[str, AbstractQueue]] = None

    @classmethod
    async def startup(
        cls,
        queue_names: Union[List[str],str],
        connection_string: str,
        exchange_name: Optional[str] = None
    ):
        self = AsyncProducer(
            queue_names=queue_names,
        )

        self.connection = await aio_pika.connect_robust(
            url=connection_string
        )
        self.channel = await self.connection.channel()
        
        self.queues = {
            queue_name: await self.channel.declare_queue(
                name=queue_name,
                durable=True
            ) for queue_name in self.queue_names
        }

        if exchange_name:
            self.exchange = await self.channel.declare_exchange(name=exchange_name, durable=True)
        else:
            self.exchange = self.channel.default_exchange

        return self

    async def send_message(self, msg: str, queue_name: str) -> Dict[str,str]:
        if queue_name not in self.queue_names:
            raise QueueNotFoundError(f'Could not find queue "{queue_name}". Available queues: {",".join(self.queue_names)}.')
        if not self.channel:
            raise ClientNotInitializedException('Channel not initialized')
        
        await self.exchange.publish(
            aio_pika.Message(
                body=msg.encode()
            ),
            routing_key=queue_name
        )
        return {'status':'success'}