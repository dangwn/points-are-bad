from fastapi.middleware.cors import CORSMiddleware

from api_client import APIClient
from rabbitmq_client import AsyncProducer
from config import (
    REQUEST_ORIGINS,
    RABBITMQ_CONNECTION_STRING,
    RABBITMQ_EXCHANGE_NAME,
    RABBITMQ_QUEUE_NAMES
)

from auth.router import router as AuthRouter
from match.router import router as MatchRouter
from points.router import router as PointsRouter
from prediction.router import router as PredictionRouter
from user.router import router as UserRouter

from typing import Dict

app: APIClient = APIClient()

app.include_router(router=AuthRouter)
app.include_router(router=MatchRouter)
app.include_router(router=PointsRouter)
app.include_router(router=PredictionRouter)
app.include_router(router=UserRouter)

app.add_middleware(
    CORSMiddleware,
    allow_origins=REQUEST_ORIGINS,
    allow_credentials=True,
    allow_methods=['*'],
    allow_headers=['*']
)

@app.on_event('startup')
async def app_startup():
    app.rabbitmq_producer: AsyncProducer = await AsyncProducer.startup(
        queue_names=RABBITMQ_QUEUE_NAMES,
        connection_string=RABBITMQ_CONNECTION_STRING,
        exchange_name=RABBITMQ_EXCHANGE_NAME
    )

@app.on_event('shutdown')
async def app_shutdown():
    await app.rabbitmq_producer.shutdown()

@app.get('/')
async def root() -> Dict[str,str]:
    return {'message':'We are go'}
