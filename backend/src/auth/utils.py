from jose import jwt
from passlib.context import CryptContext
from fastapi import Depends
from fastapi.security import OAuth2PasswordBearer

from db import get_db

from typing import (
    Union,
    Optional
)
from datetime import datetime, timedelta
from sqlalchemy.orm import Session
from user.models import User as UserModel
from config import ACCESS_TOKEN_SECRET

pwd_context: CryptContext = CryptContext(
    schemes='bcrypt',
    deprecated='auto'
)

oauth2_scheme: OAuth2PasswordBearer = OAuth2PasswordBearer(tokenUrl='login')

async def get_current_user(
    access_token: str = Depends(oauth2_scheme),
    db: Session = Depends(get_db)
) -> Optional[UserModel]:
    '''
    Retrieves the current user from db based off access token
    '''
    user_id: int = int(jwt.decode(
        token=access_token,
        key=ACCESS_TOKEN_SECRET,
        algorithms=['HS256']
    )['sub'])
    user: Optional[UserModel] = db.query(UserModel).filter(
        UserModel.user_id == user_id
    ).first()

    return user

async def is_current_user_admin(
    current_user: Optional[UserModel] = Depends(get_current_user)
) -> bool:
    '''
    Returns whether the current user is an admin user
    '''
    if not current_user:
        return False
    return current_user.is_admin
    

def generate_jwt_token(
    subject: str,
    expire_time: Union[datetime, timedelta],
    secret_key: str = ACCESS_TOKEN_SECRET,
    algorithm: str = 'HS256'
) -> str:
    '''
    Generates a JWT token
    '''
    # Change expire time into datetime if it's a timedelta
    if type(expire_time) == timedelta:
        expire_time = datetime.utcnow() + expire_time

    jwt_token: str = jwt.encode(
        {
            'sub': subject,
            'exp': expire_time
        },
        key=secret_key,
        algorithm=algorithm
    )

    return jwt_token

def hash_password(
    password: str
) -> str:
    '''
    Generates a password hash
    '''
    return pwd_context.hash(password)

def verify_password(
    plain_password: str,
    hashed_password: str
) -> bool:
    '''
    Verifies a plain text password against a hashed password
    '''
    return pwd_context.verify(
        plain_password,
        hashed_password
    )