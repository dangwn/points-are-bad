import uvicorn
from app import app
import config

if __name__ == '__main__':
    uvicorn.run(
        app=app,
        host=config.API_HOST_NAME,
        port=config.API_PORT
    )