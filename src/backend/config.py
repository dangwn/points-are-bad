import os

API_HOST_NAME = os.getenv('API_HOST_NAME','0.0.0.0')
API_PORT = os.getenv('API_PORT',8000)

DB_HOST = os.getenv('DB_HOST','localhost')
DB_NAME = os.getenv('DB_NAME','db') 
DB_PASSWORD = os.getenv('DB_PASSWORD','password') 
DB_PORT = os.getenv('DB_PORT',5432)
DB_TYPE = os.getenv('DB_TYPE','postgresql+psycopg2')
DB_USER = os.getenv('DB_USER','test') 

USERNAME_MAX_LENGTH = 30

AUTH_SECRET_KEY = os.getenv(
    'AUTH_SECRET_KEY',
    'b227c4c88cb59cce057eeb97728ce835ddd41e6a86af8a9e1f8a1ebe2a006196'
)
AUTH_ALGORITHM = os.getenv('AUTH_ALGORITHM', 'HS256')
AUTH_ACCESS_TOKEN_EXPIRE_MINUTES = 15