import os

# DB Config
DB_HOST: str = os.getenv('DB_HOST', 'localhost')
DB_NAME: str = os.getenv('DB_NAME', 'db') 
DB_PASSWORD: str = os.getenv('DB_PASSWORD', 'password') 
DB_PORT: str = os.getenv('DB_PORT', 5432)
DB_TYPE: str = os.getenv('DB_TYPE', 'postgresql+psycopg2')
DB_USER: str = os.getenv('DB_USER', 'test')

USERNAME_MAX_LENGTH: int = 30