from datetime import date
from typing import Optional, List
from sqlalchemy.orm import Session, Query
from match.models import Match as MatchModel

async def get_matches_in_date_range(
    start_date: Optional[date],
    end_date: Optional[date],
    db: Session
) -> List[MatchModel]:
    matches_query: Query = db.query(MatchModel)
    if start_date:
        matches_query = matches_query.filter(MatchModel.match_date >= start_date)
    if end_date:
        matches_query = matches_query.filter(MatchModel.match_date <= end_date)

    matches: List[MatchModel] = matches_query.order_by(MatchModel.match_date).all()
    return matches

async def insert_match_into_db(
    match_date: date,
    home: str,
    away: str,
    db: Session
) -> MatchModel:
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
    
    return new_match