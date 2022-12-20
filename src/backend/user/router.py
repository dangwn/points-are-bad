from fastapi import APIRouter, Depends, status, Response
from sqlalchemy.orm import Session
from typing import List

from db import get_db
from user import schema
from user import utils
from user.models import User as UserModel
from user.checker import check_if_email_exists, check_if_username_exists
from authentication.hash_brown import get_password_hash

from authentication.utils import get_current_user, is_current_user_admin
from http_exceptions import EMAIL_EXISTS_EXCEPTION, NOT_ADMIN_EXCEPTION, USERNAME_EXISTS_EXCEPTION, NOT_AUTHORIZED_EXCEPTION


router = APIRouter(
    prefix = '/user',
    tags = ['user']
)

@router.post('/', status_code = status.HTTP_201_CREATED)
async def create_new_user(request: schema.User, database: Session = Depends(get_db)):
    '''
    Adds new user to database, checking to see if any users already have the same username or email address
    '''
    user = await check_if_email_exists(
        request.email,
        database
    )
    if user:
        raise USERNAME_EXISTS_EXCEPTION
    
    user = await check_if_username_exists(
        request.username,
        database
    )
    if user:
        raise EMAIL_EXISTS_EXCEPTION

    request.password = get_password_hash(request.password)

    new_user = await utils.create_user(
        request, database
    )
    return new_user


@router.get('/', response_model = List[schema.DisplayUser])
async def get_all_users(
    database: Session = Depends(get_db), 
    is_admin: bool = Depends(is_current_user_admin)
) -> List[schema.DisplayUser]:
    '''
    Returns all users in the database (ADMIN ONLY)
    '''
    if not is_admin:
        raise NOT_ADMIN_EXCEPTION
    db_users =  await utils.get_all_users(database)
    return db_users 

@router.get('/{user_id}', response_model = schema.DisplayUser)
async def get_user(
    user_id: int, 
    database: Session = Depends(get_db),
    is_admin = Depends(is_current_user_admin)
):
    '''
    Returns a user by user id (ADMIN ONLY)
    '''
    if not is_admin:
        raise NOT_ADMIN_EXCEPTION
    db_user = await utils.get_user_by_id(user_id, database)
    return db_user
    
@router.delete('/{user_id}', status_code = status.HTTP_204_NO_CONTENT, response_class = Response)
async def delete_user_by_id(
    user_id: int,
    database: Session = Depends(get_db),
    current_user: UserModel = Depends(get_current_user)
):
    '''
    Deletes a user from the database
    '''
    if (current_user.id != user_id and not current_user.is_admin) or current_user is None:
        raise NOT_AUTHORIZED_EXCEPTION

    return await utils.delete_user_by_id(user_id, database)