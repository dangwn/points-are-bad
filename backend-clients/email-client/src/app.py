import asyncio
from fastapi import (
    FastAPI,
    status,
    Response
)

from rabbitmq_client import AsyncConsumer
from config import (
    RABBITMQ_CONNECTION_STRING,
    RABBITMQ_QUEUE_NAME
)
from utils import email_callback

from asyncio import AbstractEventLoop, Task

app: FastAPI = FastAPI()

@app.on_event('startup')
async def app_startup():
    app.rabbitmq_client: AsyncConsumer = await AsyncConsumer.startup(
        queue_name=RABBITMQ_QUEUE_NAME,
        connection_string=RABBITMQ_CONNECTION_STRING,
        consumer_callback=email_callback
    )

    loop: AbstractEventLoop = asyncio.get_running_loop()
    task: Task = loop.create_task(app.rabbitmq_client.consume_messages())
    await task

@app.get('/', status_code=status.HTTP_204_NO_CONTENT)
async def get_messages() -> Response:
    await app.rabbitmq_client.consume_messages()

@app.on_event('shutdown')
async def app_shutdown():
    await app.rabbitmq_client.shutdown()