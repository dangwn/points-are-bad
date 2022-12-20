from fastapi.security import OAuth2PasswordBearer
from jose import JWTError, jwt
from sqlalchemy.orm import Session

from fastapi import Depends

from datetime import datetime, timedelta

from user.models import User as UserModel
from authentication.hash_brown import verify_password
from config import AUTH_SECRET_KEY, AUTH_ALGORITHM, AUTH_ACCESS_TOKEN_EXPIRE_MINUTES
from db import get_db
from http_exceptions import CREDENTIALS_EXCEPTION, USER_NOT_FOUND_EXCEPION, PASSWORD_INCORRECT_EXCEPTION


oauth2_scheme = OAuth2PasswordBearer(tokenUrl = 'login')


async def verify_user(database, username: str, password: str) -> UserModel:
    '''
    Verifies a user against a given password, then returns the user
    '''
    user = database.query(UserModel).filter(UserModel.username == username).first()
    if not user:
        raise USER_NOT_FOUND_EXCEPION
    if not verify_password(password, user.password):
        raise PASSWORD_INCORRECT_EXCEPTION
    return user

def create_access_token(
    user_id: int,
    expires_delta: timedelta = timedelta(minutes = AUTH_ACCESS_TOKEN_EXPIRE_MINUTES)
) -> str:
    '''
    Creates access token for user authentication
    '''
    # Subject has to be a string
    data_to_encode = {'sub':str(user_id)}
    expire = datetime.utcnow() + expires_delta

    data_to_encode.update({'exp':expire})
    return jwt.encode(data_to_encode, AUTH_SECRET_KEY, algorithm = AUTH_ALGORITHM)

async def get_current_user(
    token: str = Depends(oauth2_scheme),
    database: Session = Depends(get_db)
) -> UserModel:
    '''
    Gets the current user from a given jwt token
    '''
    # Try retrieving user ID from token payload
    try:
        payload = jwt.decode(
            token, AUTH_SECRET_KEY, algorithms = [AUTH_ALGORITHM]
        )
        user_id = payload.get('sub')
        if user_id is None:
            raise CREDENTIALS_EXCEPTION
        user_id = int(user_id)
    except JWTError as e:
        raise CREDENTIALS_EXCEPTION 
    
    user = database.query(UserModel).get(user_id)
    if user is None:
        raise USER_NOT_FOUND_EXCEPION
    return user
    
async def is_current_user_admin(
    current_user: UserModel = Depends(get_current_user)
) -> bool:
    '''
    Whether the current user is an admin or not
    '''
    return current_user.is_admin

async def is_user_authorized(
    user_id: int,
    current_user: UserModel
):
    '''
    Users can access their own resources, admins can access all resources
    '''
    is_admin = await is_current_user_admin(current_user)
    
    if is_admin or user_id == current_user.id:
        return True
    return False