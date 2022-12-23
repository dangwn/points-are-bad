from fastapi import APIRouter, Depends, status, Response
from sqlalchemy.orm import Session
from typing import List

from db import get_db
from user.schema import User, DisplayUser
from user import utils
from user.models import User as UserModel
from user import checker
from authentication.hash_brown import get_password_hash

from authentication.utils import get_current_user, is_current_user_admin
from http_exceptions import EMAIL_EXISTS_EXCEPTION, NOT_ADMIN_EXCEPTION, USERNAME_EXISTS_EXCEPTION, NOT_AUTHORIZED_EXCEPTION, COULD_NOT_UPDATE_EXCEPTION


router = APIRouter(
    prefix = '/user',
    tags = ['user']
)

@router.post('/', status_code = status.HTTP_201_CREATED, response_model = DisplayUser)
async def create_new_user(request: User, database: Session = Depends(get_db)):
    '''
    Adds new user to database, checking to see if any users already have the same username or email address
    '''
    username_exists = await checker.username_exists(
        request.email,
        database
    )
    if username_exists:
        raise USERNAME_EXISTS_EXCEPTION
    
    email_exists = await checker.email_exists(
        request.username,
        database
    )
    if email_exists:
        raise EMAIL_EXISTS_EXCEPTION

    first_user = await checker.is_first_user(
        database
    )
    if first_user:
        is_admin = True
    else:
        is_admin = False

    request.password = get_password_hash(request.password)

    new_user = await utils.create_user(
        request, is_admin, database
    )
    return new_user

@router.get('/', response_model = DisplayUser)
async def display_current_user(
    current_user: UserModel = Depends(get_current_user)
) -> DisplayUser:
    return current_user

@router.delete('/', status_code = status.HTTP_204_NO_CONTENT, response_class = Response)
async def delete_current_user(
    database: Session = Depends(get_db),
    current_user: UserModel = Depends(get_current_user)
):
    return await utils.delete_user_by_id(current_user.id, database)

@router.get('/all', response_model = List[DisplayUser])
async def get_all_users(
    database: Session = Depends(get_db), 
    is_admin: bool = Depends(is_current_user_admin)
) -> List[DisplayUser]:
    '''
    Returns all users in the database (ADMIN ONLY)
    '''
    if not is_admin:
        raise NOT_ADMIN_EXCEPTION
    db_users =  await utils.get_all_users(database)
    return db_users 

@router.get('/{user_id}', response_model = DisplayUser)
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

@router.put('/make_admin/{user_id}', status_code = status.HTTP_202_ACCEPTED, response_class = Response)
async def make_user_admin(
    user_id: int,
    database: Session = Depends(get_db),
    is_admin: bool = Depends(is_current_user_admin) 
):
    '''
    Updates the DB to make the given user an admin user
    '''
    if not is_admin:
        raise NOT_ADMIN_EXCEPTION

    return await utils.alter_admin_status_by_id(user_id, True, database)

@router.put('/revoke_admin/{user_id}', status_code = status.HTTP_202_ACCEPTED, response_class = Response)
async def revoke_user_admin(
    user_id: int,
    database: Session = Depends(get_db),
    current_user: UserModel = Depends(get_current_user) 
):
    '''
    Updates the DB to make the given user an admin user
    '''
    if not current_user.is_admin:
        raise NOT_ADMIN_EXCEPTION

    if current_user.id == user_id:
        raise COULD_NOT_UPDATE_EXCEPTION("current user's admin status")

    return await utils.alter_admin_status_by_id(user_id, False, database)