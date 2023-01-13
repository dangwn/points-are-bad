import uvicorn
from app import app
import config

def main():
    uvicorn.run(app, host = config.API_HOST_NAME, port = config.API_PORT)

if __name__ == '__main__':
    main()