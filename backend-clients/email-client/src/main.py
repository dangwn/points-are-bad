import uvicorn

from app import app
from config import (
    CONSUMER_API_HOST,
    CONSUMER_API_PORT
)

if __name__ == '__main__':
    uvicorn.run(
        app=app,
        host=CONSUMER_API_HOST,
        port=CONSUMER_API_PORT
    )