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
ACCESS_TOKEN_LIFETIME_MINUTES: int = int(os.getenv('ACCESS_TOKEN_LIFETIME_MINUTES', 20))
ACCESS_TOKEN_SECRET: str = os.getenv('ACCESS_TOKEN_SECRET', 'access_secret')
REFRESH_TOKEN_LIFETIME_DAYS: int = int(os.getenv('REFRESH_TOKEN_LIFETIME_DAYS', 30))
REFRESH_TOKEN_SECRET: str = os.getenv('REFRESH_TOKEN_SECRET', 'refresh_secret')
REFRESH_TOKEN_COOKIE_KEY: str = os.getenv('REFRESH_TOKEN_COOKIE_KEY', 'X-Refresh-Token')
CSRF_TOKEN_LIFETIME_DAYS: int = int(os.getenv('CSRF_TOKEN_LIFETIME_DAYS', 30))
CSRF_TOKEN_SECRET: str = os.getenv('CSRF_TOKEN_SECRET', 'csrf_secret')
CSRF_TOKEN_COOKIE_KEY: str = os.getenv('CSRF_TOKEN_COOKIE_KEY', 'X-CSRF-Token')

# Redis Config
REDIS_HOST: str = os.getenv('REDIS_HOST', 'localhost')
REDIS_PORT: int = int(os.getenv('REDIS_PORT', 6379))
REDIS_DB: int = int(os.getenv('REDIS_DB', 0))
REDIS_PASSWORD: str = os.getenv('REDIS_PASSWORD', None)

# User Config
USERNAME_MIN_LENGTH: int = 3
USERNAME_MAX_LENGTH: int = 30

# Points config
NULL_PREDICTION_PENALTY: int = int(os.getenv('NULL_PREDICTION_PENALTY', 10))