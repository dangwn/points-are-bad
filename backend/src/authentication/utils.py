from fastapi import Depends
from db import get_db
from sqlalchemy.orm import Session
from sqlalchemy import and_

from authentication.session_manager import *
from authentication.validate import validate_token
from authentication.auth_handler import pab_auth_handler
from authentication.schema import SessionUser, HeaderCredentials
from user.models import User as UserModel

from typing import Optional, Dict

async def get_user_from_token_and_provider(
    email: str,
    provider: str,
    database: Session
) -> Optional[UserModel]:
    user: UserModel = database.query(UserModel).filter(
        and_(
            UserModel.email == email,
            UserModel.provider == provider
        )
    ).first()
    return user

async def get_current_user(
    auth_data: HeaderCredentials = Depends(pab_auth_handler),
    database: Session = Depends(get_db)
) -> Optional[SessionUser]:
    '''
    Returns the current user based on provider and access token
    -----
    Returns none if user doesn't exist or token is invalid
    '''
    access_token: str = auth_data.access_token
    provider: str = auth_data.provider

    # Get user session if it exists and return
    session_user: Optional[SessionUser] = await get_user_session(access_token, provider)
    if session_user:
        return session_user

    # See if validation token is valid, if not return none
    token_email: Optional[str] = await validate_token(access_token, provider)
    if not token_email:
        return
        
    # Check is user exists, if they do, start a new session in redis
    user: Optional[UserModel] = await get_user_from_token_and_provider(
        token_email, 
        provider, 
        database
    )
    if not user:
        return
    session_user: Optional[SessionUser] = await create_user_session(
        access_token,
        provider,
        user
    )
    return session_user

async def is_current_user_admin(
    current_user: SessionUser = Depends(get_current_user)
) -> bool:
    '''
    Returns whether the current user is an admin
    '''
    if current_user:
        return current_user.is_admin
    return False