from match.models import Match as MatchModel
from prediction.populate import populate_predictions_for_match

from datetime import date
from typing import Optional, List
from sqlalchemy.orm import Session, Query

async def get_matches_in_date_range(
    start_date: Optional[date],
    end_date: Optional[date],
    db: Session
) -> List[MatchModel]:
    matches_query: Query = db.query(MatchModel)
    if start_date:
        matches_query = matches_query.filter(MatchModel.match_date >= start_date)
    if end_date:
        matches_query = matches_query.filter(MatchModel.match_date < end_date)

    matches: List[MatchModel] = matches_query.order_by(MatchModel.match_date).all()
    return matches

async def insert_match_into_db(
    match_date: date,
    home: str,
    away: str,
    db: Session
) -> Optional[MatchModel]:
    try:
        new_match: MatchModel = MatchModel(
            match_date=match_date,
            home=home,
            away=away
        )
        db.add(new_match)
        db.commit()
        db.refresh(new_match)
    except Exception as e:
        print(e)
        return
    
    num_new_predictions: Optional[int] = await populate_predictions_for_match(
        match_id=new_match.match_id,
        db=db
    )
    if num_new_predictions is None:
        db.delete(new_match)
        db.commit()
        return
    
    return new_match

async def delete_match_by_id(
    match_id: int,
    db: Session
) -> bool:
    try:
        db.query(MatchModel).filter(MatchModel.match_id == match_id).delete()
        db.commit()
        return True
    except:
        return False
    
async def update_match_in_db(
    match_id: int,
    home: str,
    away: str,
    match_date: date,
    home_goals: int,
    away_goals: int,
    db: Session
) -> Optional[MatchModel]:
    match: Optional[MatchModel] = db.query(MatchModel).filter(
        MatchModel.match_id == match_id
    ).first()

    if not match:
        return
    
    try:
        match.home = home
        match.away = away
        match.match_date = match_date
        match.home_goals = home_goals
        match.away_goals = away_goals

        db.commit()

        return match
    except:
        return