from fastapi import (
    APIRouter,
    Depends,
    status,
    Response
)

from db import get_db
from match.models import Match as MatchModel
from match.schema import (
    MatchWithoutGoals,
    Match,
    MatchWithId
)
from match.utils import (
    delete_match_by_id,
    insert_match_into_db,
    get_matches_in_date_range,
    update_match_in_db
)
from auth.utils import is_current_user_admin
from exceptions import (
    COULD_NOT_CREATE_MATCH_EXCEPTION,
    USER_IS_NOT_ADMIN_EXCEPTION,
    COULD_NOT_DELETE_EXCEPTION,
    COULD_NOT_UPDATE_EXCEPTION
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

@router.post('/', response_model=MatchWithId)
async def create_match(
    match: MatchWithoutGoals,
    user_is_admin: bool = Depends(is_current_user_admin),
    db: Session = Depends(get_db)
) -> MatchWithId:
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

@router.delete('/{match_id}', status_code=status.HTTP_204_NO_CONTENT)
async def delete_match(
    match_id: int,
    user_is_admin: bool = Depends(is_current_user_admin),
    db: Session = Depends(get_db)
) -> Response:
    if not user_is_admin:
        raise USER_IS_NOT_ADMIN_EXCEPTION
    
    match_deleted: bool = await delete_match_by_id(match_id=match_id, db=db)
    if not match_deleted:
        raise COULD_NOT_DELETE_EXCEPTION('match')
    
@router.get('/full/', response_model=List[MatchWithId])
async def get_full_matches(
    start_date: Optional[date] = None,
    end_date: Optional[date] = None,
    user_is_admin: bool = Depends(is_current_user_admin),
    db: Session = Depends(get_db)
) -> List[MatchWithId]:
    if not user_is_admin:
        raise USER_IS_NOT_ADMIN_EXCEPTION
    
    matches: List[MatchModel] = await get_matches_in_date_range(
        start_date=start_date,
        end_date=end_date,
        db=db
    )

    return matches

@router.put('/', status_code=status.HTTP_202_ACCEPTED, response_model=MatchWithId)
async def update_match(
    match: MatchWithId,
    user_is_admin: bool = Depends(is_current_user_admin),
    db: Session = Depends(get_db)
) -> MatchWithId:
    if not user_is_admin:
        raise USER_IS_NOT_ADMIN_EXCEPTION
    
    updated_match: Optional[MatchModel] = await update_match_in_db(
        match_id=match.match_id,
        home=match.home,
        away=match.away,
        match_date=match.match_date,
        home_goals=match.home_goals,
        away_goals=match.away_goals,
        db=db
    )

    if not updated_match:
        raise COULD_NOT_UPDATE_EXCEPTION(what='match', why='database issue')
    
    return updated_match 

@router.post('/createTestMatch/', response_model=Match)
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