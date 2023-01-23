from sqlalchemy import and_
from sqlalchemy.orm import Session
from typing import List, Optional

from user.models import User as UserModel
from predictions.populate import populate_predictions
from points.models import PlayerPoints as PlayerPointsModel
from http_exceptions import COULD_NOT_UPDATE_EXCEPTION, USER_NOT_FOUND_EXCEPION

async def create_user(
    display_name: str,
    email: str,
    provider: str,
    avatar: str,
    database: Session
) -> UserModel:
    # The first user added to the db is an admin
    is_admin: bool = False
    if not database.query(UserModel).first():
        is_admin = True

    try:
        new_user: UserModel = UserModel(
            display_name=display_name,
            email=email,
            avatar=avatar,
            provider=provider,
            is_admin=is_admin
        )
        database.add(new_user)
        database.commit()
        database.refresh(new_user)
    except:
        raise COULD_NOT_UPDATE_EXCEPTION('users table when creating new user')

    # Add fresh score as well as user
    try:
        player_score = PlayerPointsModel(new_user.id)
        database.add(player_score)
        database.commit()
        database.refresh(player_score)
    except:
        # Delete user if points could not be created
        try:
            await delete_user_by_id(new_user.id, database)
        except:
            raise COULD_NOT_UPDATE_EXCEPTION('user table when deleting user')
        raise COULD_NOT_UPDATE_EXCEPTION('points table')

    try:
        await populate_predictions(database, user_id = new_user.id)
    except:
        # Delete user if predictions could not be created
        try:
            await delete_user_by_id(new_user.id, database)
        except:
            raise COULD_NOT_UPDATE_EXCEPTION('user table when deleting user')
        raise COULD_NOT_UPDATE_EXCEPTION('predictions table')

    return new_user

async def get_user_by_id(
    user_id: int,
    database: Session
) -> Optional[UserModel]:
    return database.query(UserModel).filter(
        UserModel.id == user_id
    ).first()

async def get_user_by_email_and_provider(
    email: str,
    provider: str,
    database: Session
) -> Optional[UserModel]:
    return database.query(UserModel).filter(
        and_(
            UserModel.email == email,
            UserModel.provider == provider
        )
    ).first()

async def get_all_users(
    database: Session
) -> List[UserModel]:
    users: List[UserModel] = database.query(UserModel).all()
    return users

async def delete_user_by_id(
    user_id: int, 
    database: Session
) -> None:
    '''
    Deletes a user in the database given their ID
    Note: Cannot delete user if they are the final admin user
    '''
    # Check to make sure that there is at least one admin user remaining
    admin_user: UserModel = database.query(UserModel).filter(
        and_(
            UserModel.id != user_id,
            UserModel.is_admin == True
        )
    ).first()

    if not admin_user:
        raise COULD_NOT_UPDATE_EXCEPTION('users table as there needs to be at least one admin user')
    
    try:
        database.query(UserModel).filter(UserModel.id == user_id).delete()
        database.commit()
    except:
        raise COULD_NOT_UPDATE_EXCEPTION(f'user table when deleting user with id: {user_id}')

async def alter_admin_status_by_id(
    user_id: int,
    admin_status: bool,
    database: Session
) -> None:
    '''
    Alter user's admin status
    '''
    user: UserModel = database.query(UserModel).filter(UserModel.id == user_id).first()
    if not user:
        raise USER_NOT_FOUND_EXCEPION

    try:
        user.is_admin = admin_status
        database.commit()
    except:
        raise COULD_NOT_UPDATE_EXCEPTION(f"user {user.id}'s admin status") 