from typing import List
from fastapi import APIRouter, Depends, status, Body, Response, Query
from sqlalchemy.orm import Session

from db import get_db
from authentication.utils import is_current_user_admin
from http_exceptions import NOT_ADMIN_EXCEPTION, MATCH_NOT_FOUND_EXCEPTION, COULD_NOT_UPDATE_EXCEPTION, NOT_AUTHORIZED_EXCEPTION
from matches.schema import Match, DisplayMatch, MatchScore, MatchScoreWithId
from matches import utils

router = APIRouter(
    prefix = '/matches',
    tags = ['matches']
)

@router.get('/', response_model = List[DisplayMatch])
async def get_all_matches(database: Session = Depends(get_db)) -> List[DisplayMatch]:
    '''
    Returns a list of all matches
    '''
    matches = await utils.get_all_matches(database)
    return matches

@router.post('/', status_code = status.HTTP_201_CREATED, response_model = DisplayMatch)
async def create_new_match(
    match: Match,
    database: Session = Depends(get_db),
    is_admin: bool = Depends(is_current_user_admin)
) -> DisplayMatch:
    '''
    Inserts a new match into the db
    '''
    if not is_admin:
        raise NOT_ADMIN_EXCEPTION
    new_match = await utils.create_new_match(match, database)
    return new_match

@router.get('/{match_id}', response_model = DisplayMatch)
async def get_match_by_id(match_id: int, database: Session = Depends(get_db)) -> DisplayMatch:
    '''
    Gets a match by its ID
    '''
    match = await utils.get_match_by_id(match_id, database)
    if not match:
        raise MATCH_NOT_FOUND_EXCEPTION
    return match

@router.get('/get_matches_by_id/', response_model = List[DisplayMatch])
async def get_multiple_matches_by_id(
    match_ids: List[int] = Query(...),
    database: Session = Depends(get_db)
) -> List[DisplayMatch]:
    '''
    Retrieves multiple matches by given ids
    '''
    matches = await utils.get_multiple_matches_by_id(match_ids, database)
    if not matches:
        raise MATCH_NOT_FOUND_EXCEPTION
    return matches

@router.put('/update_matches_by_id', response_model = List[DisplayMatch])
async def update_multiple_matches_by_id(
    match_scores: List[MatchScoreWithId] = Body(...),
    database: Session = Depends(get_db),
    is_admin: bool = Depends(is_current_user_admin)
) -> List[DisplayMatch]:
    '''
    Updates multiple matches of specified IDs
    '''
    if not is_admin:
        raise NOT_ADMIN_EXCEPTION

    new_match_scores = await utils.update_multiple_matches_by_id(match_scores, database)
    if not new_match_scores:
        raise MATCH_NOT_FOUND_EXCEPTION
    
    return new_match_scores

@router.delete('/{match_id}', status_code = status.HTTP_204_NO_CONTENT, response_class = Response)
async def delete_user_by_id(
    match_id: int,
    database: Session = Depends(get_db),
    is_admin: Match = Depends(is_current_user_admin)
):
    '''
    Deletes a match from the database
    '''
    if not is_admin:
        raise NOT_AUTHORIZED_EXCEPTION

    return await utils.delete_match_by_id(match_id, database)