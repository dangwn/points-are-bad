import redis
from datetime import timedelta
from config import (
    REDIS_DB,
    REDIS_HOST,
    REDIS_PASSWORD,
    REDIS_PORT
)

redis_client: redis.Redis = redis.Redis(
    host=REDIS_HOST,
    port=REDIS_PORT,
    db=REDIS_DB,
    password=REDIS_PASSWORD
)

def redis_get(
    key: str
) -> str:
    return redis_client.get(name=key).decode('utf-8')

def redis_set(
    key: str,
    value: str,
    expire_minutes: int
) -> bool:
    return redis_client.set(
        name=key,
        value=value,
        ex=timedelta(minutes=expire_minutes)
    )
