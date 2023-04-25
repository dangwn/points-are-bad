from fastapi import FastAPI
from AsyncPikaClient import AsyncProducer

from typing import Optional

class APIClient(FastAPI):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)

        self.rabbitmq_producer: Optional[AsyncProducer] = None