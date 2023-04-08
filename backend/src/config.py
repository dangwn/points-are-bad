import os
from typing import List

# App Config
API_HOST_NAME: str = os.getenv('API_HOST_NAME', '0.0.0.0')
API_PORT: int = os.getenv('API_PORT', 8000)
REQUEST_ORIGINS: List[str] = [
    'http://localhost',
    'http://localhost:3000',
    f'http://localhost:{API_PORT}'
]

# DB Config
DB_HOST: str = os.getenv('DB_HOST', 'localhost')
DB_NAME: str = os.getenv('DB_NAME', 'db') 
DB_PASSWORD: str = os.getenv('DB_PASSWORD', 'password') 
DB_PORT: str = os.getenv('DB_PORT', 5432)
DB_TYPE: str = os.getenv('DB_TYPE', 'postgresql+psycopg2')
DB_USER: str = os.getenv('DB_USER', 'test')

# Auth Config
ACCESS_TOKEN_LIFETIME_MINUTES: int = 20
ACCESS_TOKEN_SECRET: str = 'access_secret'
REFRESH_TOKEN_LIFETIME_DAYS: int = 30
REFRESH_TOKEN_SECRET: str = 'refresh_secret'
REFRESH_TOKEN_COOKIE_KEY: str = 'X-Refresh-Token'
CSRF_TOKEN_LIFETIME_DAYS: int = 30
CSRF_TOKEN_SECRET: str = 'csrf_secret'
CSRF_TOKEN_COOKIE_KEY: str = 'X-CSRF-Token'

# User Config
USERNAME_MIN_LENGTH: int = 3
USERNAME_MAX_LENGTH: int = 30
