import os

API_HOST_NAME = os.getenv('API_HOST_NAME', '0.0.0.0')
API_PORT = os.getenv('API_PORT', 8000)

DB_HOST = os.getenv('DB_HOST', 'localhost')
DB_NAME = os.getenv('DB_NAME', 'db') 
DB_PASSWORD = os.getenv('DB_PASSWORD', 'password') 
DB_PORT = os.getenv('DB_PORT', 5432)
DB_TYPE = os.getenv('DB_TYPE', 'postgresql+psycopg2')
DB_USER = os.getenv('DB_USER', 'test') 

USERNAME_MAX_LENGTH = 30

REDIS_HOST = os.getenv('REDIS_HOST', 'localhost')
REDIS_PORT = int(os.getenv('REDIS_PORT', 6379))
REDIS_DB = int(os.getenv('REDIS_DB', 0))
REDIS_PASSWORD = os.getenv('REDIS_PASSWORD', None)

GOOGLE_CLIENT_ID = os.getenv(
    'GOOGLE_CLIENT_ID',
    ''
)

if GOOGLE_CLIENT_ID == '':
    raise Exception('Google client ID not set')