from typing import List
from sqlalchemy.orm import Session

from user.models import User as UserModel
from user.schema import User
from points.models import PlayerPoints as PlayerPointsModel
from predictions.populate import populate_predictions
from http_exceptions import USER_NOT_FOUND_EXCEPION, COULD_NOT_UPDATE_EXCEPTION

async def create_user(request: User, database: Session) -> UserModel:
    '''
    Creates new user and player points in database
    '''
    user = UserModel(
        username = request.username, 
        email = request.email, 
        password = request.password, 
        is_admin = request.is_admin
    )
    
    try:
        database.add(user)
        database.commit()
        database.refresh(user)
    except:
        raise COULD_NOT_UPDATE_EXCEPTION('user table')

    # Add fresh score as well as user
    try:
        player_score = PlayerPointsModel(user.id)
        database.add(player_score)
        database.commit()
        database.refresh(player_score)
    except:
        # Delete user if points could not be created
        try:
            _ = await delete_user_by_id(user.id, database)
        except:
            raise COULD_NOT_UPDATE_EXCEPTION('user table when deleting user')
        raise COULD_NOT_UPDATE_EXCEPTION('points table')

    try:
        _ = await populate_predictions(database, user_id = user.id)
    except:
        # Delete user if predictions could not be created
        try:
            _ = await delete_user_by_id(user.id, database)
        except:
            raise COULD_NOT_UPDATE_EXCEPTION('user table when deleting user')
        raise COULD_NOT_UPDATE_EXCEPTION('predictions table')

    return user


async def get_all_users(database: Session) -> List[UserModel]:
    '''
    Gets all users in the user table
    '''
    users = database.query(UserModel).all()
    return users

async def get_user_by_id(user_id: int, database: Session) -> UserModel:
    '''
    Returns a user based on their user id
    '''
    user = database.query(UserModel).get(user_id)
    if not user:
        raise USER_NOT_FOUND_EXCEPION
    return user

async def get_user_by_username(username: str, database: Session) -> UserModel:
    '''
    Returns a user based on their username
    '''
    user = database.query(UserModel).filter(UserModel.username == username).first()
    if not user:
        raise USER_NOT_FOUND_EXCEPION
    return user

async def delete_user_by_id(user_id: int, database: Session) -> None:
    database.query(UserModel).filter(UserModel.id == user_id).delete()
    database.commit()

