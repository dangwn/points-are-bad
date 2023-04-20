from fastapi import HTTPException, status
from datetime import date
from sqlalchemy import text, update

from config import NULL_PREDICTION_PENALTY
from points.models import Points as PointsModel

from sqlalchemy.orm import Session
from typing import Optional, List, Tuple, Dict, Any

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

async def update_all_points(
    db: Session
) -> bool:
    today_date: str = date.today().strftime('%Y-%m-%d')

    try:
        # Need to spark this bad boy
        results: List[Tuple] = db.execute(text(f'''
        SELECT *, RANK() OVER (ORDER BY points ASC, correct_scores DESC, largest_error ASC) as position
        FROM (    
            SELECT user_id, SUM(points) as points, COUNT(CASE WHEN points = 0 THEN 1 END) as correct_scores, MAX(points) as largest_error
            FROM (
                SELECT user_id, COALESCE(ABS(pred_hg-hg) + ABS(pred_ag-ag), {NULL_PREDICTION_PENALTY}) as points
                FROM (
                    SELECT user_id, pred_hg, pred_ag, home_goals as hg, away_goals as ag
                    FROM (  
                        SELECT "user_id", "home_goals" as "pred_hg", "away_goals" as "pred_ag", match_id
                        FROM predictions
                    ) as t1
                    JOIN matches on t1.match_id = matches.match_id
                    WHERE 
                        matches.match_date < DATE('{today_date}') AND
                        matches.home_goals IS NOT NULL AND
                        matches.away_goals IS NOT NULL
                ) as t2
            ) as t3
            GROUP BY user_id
        ) as t4
        ''')).all()
        results_keys: Tuple[str] = ('user_id', 'points', 'correct_scores', 'largest_error', 'position')

        new_points_data: List[Dict[str, Any]] = [dict(zip(results_keys, row)) for row in results]
        db.bulk_update_mappings(
            PointsModel, new_points_data
        )
        db.commit()
    except: 
        return False

    return True