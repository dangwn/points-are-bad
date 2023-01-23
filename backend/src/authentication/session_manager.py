from redis import Redis
import json
from typing import Optional, Dict, Any

from user.models import User as UserModel
from authentication.schema import SessionUser
from config import REDIS_HOST, REDIS_PORT, REDIS_PASSWORD, REDIS_DB

__all__ = [
    'create_user_session',
    'get_user_session',
    'delete_user_session'
]

redis_session: Redis = Redis(
    host=REDIS_HOST,
    port=REDIS_PORT,
    db=REDIS_DB,
    password=REDIS_PASSWORD
)

def create_redis_key(access_token: str, provider: str) -> str:
    '''
    Creates hashed session key from the provider and access token
    @TODO create one way hash
    '''
    return f'{provider}--{access_token}'

async def create_user_session(
    access_token: str,
    provider: str,
    user: UserModel,
    expire_seconds: int = 900
) -> Optional[SessionUser]:
    '''
    Creates session in redis with hashed key
    '''
    redis_key: str = create_redis_key(access_token, provider)
    
    try:
        session_data: Dict[str,Any] = {
            'id':user.id,
            'display_name':user.display_name,
            'email':user.email,
            'avatar':user.avatar,
            'provider':user.provider,
            'is_admin':user.is_admin
        }
        redis_session.set(
            name=redis_key,
            value=json.dumps(session_data),
            ex=expire_seconds
        )
        return SessionUser(**session_data)
    except Exception as e:
        print(e)
        return

async def get_user_session(
    access_token: str,
    provider: str
) -> Optional[SessionUser]:
    '''
    Gets the users session info
    '''
    redis_key: str = create_redis_key(access_token, provider)
    try:
        raw_session: bytes = redis_session.get(redis_key)
        if raw_session:
            session_data: Dict[str, str] = json.loads(raw_session)
            return SessionUser(**session_data)
        raise ValueError('Could not retrieve session with given provider and access_token')
    except Exception as e:
        print(e)
        return

async def delete_user_session(
    access_token: str,
    provider: str
) -> Optional[bool]:
    '''
    Deletes users session from redis
    '''
    redis_key: str = create_redis_key(access_token, provider)
    try: 
        redis_session.delete(redis_key)
        return True
    except Exception as e:
        return 
