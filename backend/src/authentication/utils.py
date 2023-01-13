from fastapi.security import OAuth2PasswordBearer
from jose import JWTError, jwt
from sqlalchemy.orm import Session

from fastapi import Depends

from datetime import datetime, timedelta

from fastapi_jwt_auth import AuthJWT

from user.models import User as UserModel
from authentication.hash_brown import verify_password
from config import AUTH_SECRET_KEY, AUTH_ALGORITHM, AUTH_ACCESS_TOKEN_EXPIRE_MINUTES
from db import get_db
from http_exceptions import CREDENTIALS_EXCEPTION, USER_NOT_FOUND_EXCEPION, PASSWORD_INCORRECT_EXCEPTION


oauth2_scheme = OAuth2PasswordBearer(tokenUrl = 'login')


async def verify_user(username: str, password: str, database: Session) -> UserModel:
    '''
    Verifies a user against a given password, then returns the user
    '''
    user = database.query(UserModel).filter(UserModel.username == username).first()
    if not user:
        raise USER_NOT_FOUND_EXCEPION
    if not verify_password(password, user.password):
        raise PASSWORD_INCORRECT_EXCEPTION
    return user

async def get_current_user(
    Authorize: AuthJWT = Depends(),
    database: Session = Depends(get_db)
) -> UserModel:
    '''
    Gets the current user from a given jwt token
    '''
    Authorize.jwt_required()

    user_id = Authorize.get_jwt_subject()
    
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