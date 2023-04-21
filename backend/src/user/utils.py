from fastapi import HTTPException, status
from sqlalchemy import and_

from points.utils import insert_points_into_db
from points.models import Points as PointsModel
from user.models import User as UserModel
from prediction.populate import populate_predictions_for_user

from sqlalchemy.orm import Session
from typing import Optional

async def insert_user_into_db(
    username: str,
    email: str,
    hashed_password: str,
    db: Session
) -> UserModel:
    '''
    Runs database query to insert a new user into the database
    When the first user is added, they become an admin
        - All other users by default will not be admins
    Users will be verified separately
    '''
    is_admin: bool = False
    if not db.query(UserModel).first():
        is_admin = True
        
    try:
        new_user: Optional[UserModel] = UserModel(
            username=username,
            email=email,
            hashed_password=hashed_password,
            is_admin=is_admin
        )
        db.add(new_user)
        db.commit()
        db.refresh(new_user)
    except Exception as e:
        print(e)
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail='Could not insert new user into database'
        )
    
    new_points: Optional[PointsModel] = await insert_points_into_db(
        user_id=new_user.user_id,
        db=db
    )
    if not new_points:
        db.delete(new_user)
        db.commit()
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail='Could not insert new user into database as could not create new points entry'
        )
    
    num_new_predictions: Optional[int] = await populate_predictions_for_user(
        user_id=new_user.user_id,
        db=db
    )
    if num_new_predictions is None:
        db.delete(new_user)
        db.delete(new_points)
        db.commit()
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail='Could not insert new user into database as could not create predictions'
        )
    
    return new_user

async def delete_user_by_id(
    user_id: int,
    db: Session
) -> None:
    '''
    Deletes a given user by ID
    '''
    admin_user: UserModel = db.query(UserModel).filter(
        and_(
            UserModel.user_id != user_id,
            UserModel.is_admin == True
        )
    ).first()
    if not admin_user:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail='User table needs at least one admin user.'
        )
    
    try:
        db.query(UserModel).filter(UserModel.user_id == user_id).delete()
        db.commit()
    except Exception as e:
        print(e)
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail='Could not delete user'
        )
    
async def change_username_by_id(
    user_id: int,
    new_username: str,
    db: Session
) -> Optional[str]:
    user: Optional[UserModel] = db.query(UserModel).filter(UserModel.user_id == user_id).first()
    if not user:
        return
    
    user.username = new_username
    db.commit()
    return new_username