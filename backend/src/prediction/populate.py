from user.models import User
from match.models import Match
from prediction.models import Prediction

from sqlalchemy.orm import Session
from typing import (
    List,
    Optional
)

async def populate_predictions_for_match(
    match_id: int,
    db: Session
) -> Optional[int]:
    user_ids: List[int] = db.query(User.user_id).all()
    if len(user_ids) == 0:
        return 0
    if not user_ids:
        return
    
    user_predictions: List[Prediction] = [Prediction(
        user_id=user_id[0], 
        match_id=match_id
    ) for user_id in user_ids]

    try:
        db.add_all(user_predictions)
        db.commit()
    except:
        return
    
    return len(user_predictions)

async def populate_predictions_for_user(
    user_id: int,
    db: Session
) -> Optional[int]:
    match_ids: List[int] = db.query(Match.match_id).all()
    if len(match_ids) == 0:
        return 0
    if not match_ids:
        return
    
    match_predictions: List[Prediction] = [Prediction(
        user_id=user_id, 
        match_id=match_id[0]
    ) for match_id in match_ids]
    
    try:
        db.add_all(match_predictions)
        db.commit()
    except:
        return
    
    return len(match_ids)