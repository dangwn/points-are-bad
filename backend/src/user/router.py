from fastapi import APIRouter, Depends, status, Response
from sqlalchemy.orm import Session
from typing import Optional, Dict, List

from user.schema import DisplayUser, CreateUser
from user.models import User as UserModel
from user.utils import (
    create_user,
    delete_user_by_id,
    get_user_by_email_and_provider,
    get_all_users,
    get_user_by_id,
    alter_admin_status_by_id
)
from authentication.schema import SessionUser, HeaderCredentials
from authentication.utils import get_current_user, is_current_user_admin
from authentication.auth_handler import pab_auth_handler
from authentication.session_manager import create_user_session, delete_user_session
from authentication.validate import validate_token
from http_exceptions import (
    NOT_AUTHORIZED_EXCEPTION,
    USERNAME_EXISTS_EXCEPTION,
    NOT_ADMIN_EXCEPTION,
    USER_NOT_FOUND_EXCEPION,
    COULD_NOT_UPDATE_EXCEPTION
)
from db import get_db

router: APIRouter = APIRouter(
    prefix='/user',
    tags=['user']
)

@router.get('/', response_model=DisplayUser)
async def display_current_user(
    current_user: SessionUser = Depends(get_current_user)
) -> DisplayUser:
    if not current_user:
        raise USER_NOT_FOUND_EXCEPION
    return current_user

@router.post('/', response_model=DisplayUser)
async def create_new_user(
    new_user: CreateUser,
    database: Session = Depends(get_db),
    credentials: HeaderCredentials = Depends(pab_auth_handler)
) -> DisplayUser:
    # See if access token is valid
    token_email: str = await validate_token(
        token=credentials.access_token,
        provider=credentials.provider
    )
    if not token_email:
        raise NOT_AUTHORIZED_EXCEPTION

    # See if a user with same email and provider exists
    user_already_exists: UserModel = await get_user_by_email_and_provider(
        email=token_email,
        provider=credentials.provider,
        database=database
    )
    if user_already_exists:
        raise USERNAME_EXISTS_EXCEPTION

    # Create the new user
    user: UserModel = await create_user(
        display_name=new_user.display_name,
        email=token_email,
        provider=credentials.provider,
        avatar=new_user.avatar,
        database=database
    )

    # Create user session
    await create_user_session(credentials.access_token, credentials.provider, user)

    return user

@router.delete('/', status_code=status.HTTP_204_NO_CONTENT, response_class=Response)
async def delete_current_user(
    database: Session = Depends(get_db),
    current_user: SessionUser = Depends(get_current_user),
    credentials: HeaderCredentials = Depends(pab_auth_handler)
) -> Response:
    await delete_user_session(credentials.access_token, credentials.provider)
    return await delete_user_by_id(current_user.id, database)

@router.get('/all', response_model=List[DisplayUser])
async def display_all_users(
    database: Session = Depends(get_db), 
    is_admin: bool = Depends(is_current_user_admin)
) -> List[DisplayUser]:
    '''
    Returns all users in the database (ADMIN ONLY)
    '''
    if not is_admin:
        raise NOT_ADMIN_EXCEPTION
    db_users: List[UserModel] =  await get_all_users(database)
    return db_users 

@router.get('/{user_id}', response_model=DisplayUser)
async def get_user(
    user_id: int, 
    database: Session = Depends(get_db),
    is_admin: bool = Depends(is_current_user_admin)
) -> DisplayUser:
    '''
    Returns a user by user id (ADMIN ONLY)
    '''
    if not is_admin:
        raise NOT_ADMIN_EXCEPTION
    db_user: Optional[UserModel] = await get_user_by_id(user_id, database)

    if not db_user:
        raise USER_NOT_FOUND_EXCEPION
    return db_user

@router.delete('/{user_id}', status_code=status.HTTP_204_NO_CONTENT, response_class=Response)
async def remove_user_by_id(
    user_id: int,
    database: Session = Depends(get_db),
    current_user: SessionUser = Depends(get_current_user)
) -> Response:
    '''
    Deletes a user from the database
    '''
    if (current_user.id != user_id and not current_user.is_admin) or current_user is None:
        raise NOT_AUTHORIZED_EXCEPTION

    return await delete_user_by_id(user_id, database)

@router.delete('/end_session', status_code=status.HTTP_204_NO_CONTENT, response_class=Response)
async def end_current_user_session(
    credentials: HeaderCredentials = Depends(pab_auth_handler)
) -> None:
    await delete_user_session(credentials.access_token, credentials.provider)

@router.put('/make_admin/{user_id}', status_code=status.HTTP_202_ACCEPTED, response_class=Response)
async def make_user_admin(
    user_id: int,
    database: Session = Depends(get_db),
    is_admin: bool = Depends(is_current_user_admin) 
) -> Response:
    '''
    Updates the DB to make the given user an admin user
    '''
    if not is_admin:
        raise NOT_ADMIN_EXCEPTION

    return await alter_admin_status_by_id(user_id, True, database)

@router.put('/revoke_admin/{user_id}', status_code=status.HTTP_202_ACCEPTED, response_class=Response)
async def revoke_user_admin(
    user_id: int,
    database: Session = Depends(get_db),
    current_user: SessionUser = Depends(get_current_user) 
) -> Response:
    '''
    Updates the DB to make the given user an admin user
    '''
    if not current_user.is_admin:
        raise NOT_ADMIN_EXCEPTION

    if current_user.id == user_id:
        raise COULD_NOT_UPDATE_EXCEPTION("current user's admin status")

    return await alter_admin_status_by_id(user_id, False, database)