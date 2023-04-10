from jose import jwt
from passlib.context import CryptContext
from fastapi import Depends
from fastapi.security import OAuth2PasswordBearer
import re
import json
import base64

from db import get_db

from typing import (
    Union,
    Optional,
    Dict
)
from datetime import datetime, timedelta
from sqlalchemy.orm import Session

from auth.redis_client import (
    redis_get,
    redis_set
)
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

async def validate_email(email: str) -> bool:
    '''
    Validates email format
    '''  
    if re.match(r"[^@]+@[^@]+\.[^@]+", email):  
        return True  
    return False   

async def is_user_email_in_db(
    email: str,
    db: Session
) -> bool:
    '''
    Checks whether there is a user in the database that has the
        provided email
    '''
    if db.query(UserModel).filter(
        UserModel.email == email
    ).first():
        return True
    return False

def create_token_sauce() -> str:
    '''
    Creates the random number key for the email
    '''
    return '123456'

async def create_verification_token(
    email: str
) -> str:
    '''
    Create the validation string
    '''
    token_sauce: str = create_token_sauce()
    validation_dict: Dict[str,str] = {
        email: token_sauce
    }

    json_str: str = json.dumps(validation_dict)
    encoded_str: str = base64.b64encode(json_str.encode('utf-8')).decode('utf-8')

    redis_set(
        key=email,
        value=token_sauce,
        expire_minutes=30
    )

    return encoded_str

async def validate_verification_token(
    token: str
) -> Optional[str]:
    try:
        token_dict: Dict[str,str] = json.loads(
            base64.b64decode(
                token.encode('utf-8')
            )
        )
    except:
        return
    
    user_email: str = tuple(token_dict.keys())[0]
    token_sauce: str = token_dict[user_email]
    redis_token_sauce: str = redis_get(user_email)

    if token_sauce != redis_token_sauce:
        return

    return user_email