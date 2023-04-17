from fastapi import HTTPException, status
from sqlalchemy import func

from points.models import Points as PointsModel

from sqlalchemy.orm import Session, Query
from sqlalchemy.sql.expression import Select
from typing import Optional, List

async def insert_points_into_db(
    user_id: int,
    db: Session
) -> PointsModel:
    '''
    Runs SQL query to insert a new player points into the db
    '''
    try:
        new_points: PointsModel = PointsModel(
            user_id=user_id
        )
        db.add(new_points)
        db.commit()
        db.refresh(new_points)
    except Exception as e:
        print(e)
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail='Could not insert player points into database'
        )
    
    return new_points

async def get_points_by_user_id(
    user_id: int,
    db: Session
) -> Optional[PointsModel]:
    '''
    Returns a user's points
    '''
    user_points: Optional[PointsModel] = db.query(PointsModel).filter(
        PointsModel.user_id == user_id
    ).first()

    if not user_points:
        return
    return user_points

async def get_leaderboard(
    limit: int,
    offset: int, 
    db: Session
) -> List[PointsModel]:
    leaderboard: Optional[List[PointsModel]] = db.query(
        PointsModel
    ).order_by(
        PointsModel.position
    ).offset(
        offset=offset
    ).limit(
        limit=limit
    ).all()

    if not leaderboard:
        return
    return leaderboard