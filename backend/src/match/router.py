from fastapi import (
    APIRouter,
    Depends,
    status,
    HTTPException
)

from db import get_db
from match.models import Match as MatchModel
from match.schema import (
    MatchWithoutGoals,
    Match
)
from match.utils import (
    insert_match_into_db,
    get_matches_in_date_range
)
from auth.utils import is_current_user_admin
from exceptions import (
    COULD_NOT_CREATE_MATCH_EXCEPTION,
    USER_IS_NOT_ADMIN_EXCEPTION
)

from sqlalchemy.orm import Session
from typing import (
    Optional,
    List
)
from datetime import date 

router: APIRouter = APIRouter(
    prefix='/match'
)

@router.get('/', response_model=List[MatchWithoutGoals])
async def get_matches_without_goals(
    start_date: Optional[date] = None,
    end_date: Optional[date] = None,
    db: Session = Depends(get_db)
) -> List[MatchWithoutGoals]:
    matches: List[MatchModel] = await get_matches_in_date_range(
        start_date=start_date,
        end_date=end_date,
        db=db
    )

    return matches

@router.post('/', response_model=Match)
async def create_match(
    match: MatchWithoutGoals,
    user_is_admin: bool = Depends(is_current_user_admin),
    db: Session = Depends(get_db)
) -> Match:
    if not user_is_admin:
        raise USER_IS_NOT_ADMIN_EXCEPTION
    new_match: Optional[MatchModel] = await insert_match_into_db(
        match_date=match.match_date,
        home=match.home,
        away=match.away,
        db=db
    )

    if not new_match:
        raise COULD_NOT_CREATE_MATCH_EXCEPTION
    return new_match

@router.post('/createTestMatch', response_model=Match)
async def create_test_match(
    match: MatchWithoutGoals,
    db: Session = Depends(get_db)
) -> Match:
    new_match: Optional[MatchModel] = await insert_match_into_db(
        match_date=match.match_date,
        home=match.home,
        away=match.away,
        db=db
    )

    if not new_match:
        raise COULD_NOT_CREATE_MATCH_EXCEPTION
    return new_match