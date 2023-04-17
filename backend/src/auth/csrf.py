from pydantic import BaseModel

from config import (
    CSRF_TOKEN_LIFETIME_DAYS,
    CSRF_TOKEN_SECRET
)

class CsrfSettings(BaseModel):
    '''
    Settings for CSRF tokens
    '''
    secret_key: str = CSRF_TOKEN_SECRET
    max_age: int = 3600 * 24 * CSRF_TOKEN_LIFETIME_DAYS

