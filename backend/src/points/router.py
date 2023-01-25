from fastapi import APIRouter, Depends, status, Body
from sqlalchemy.orm import Session
from typing import List

from db import get_db
from authentication.utils import get_current_user, is_current_user_admin
from points import utils
from user.models import User as UserModel
from points.schema import DayScore, DayScoreWithUserId, PlayerPosition, PlayerPositionWithDisplayName
from http_exceptions import USER_NOT_FOUND_EXCEPION, COULD_NOT_UPDATE_EXCEPTION, NOT_ADMIN_EXCEPTION, INVALID_QUERY_PARAMETERS_EXCEPTION

router = APIRouter(
    prefix = '/points',
    tags = ['points']
)

@router.get('/', response_model = DayScore)
async def get_current_player_score(
    database: Session = Depends(get_db),
    current_user: UserModel = Depends(get_current_user)
):
    '''
    Retrieves the current user's score
    '''
    user_score = await utils.get_player_score_by_id(current_user.id, database)
    if not user_score:
        raise USER_NOT_FOUND_EXCEPION
    return user_score

@router.get('/all', response_model = List[DayScore])
async def get_all_player_scores(
    database: Session = Depends(get_db),
    is_admin: bool = Depends(is_current_user_admin)
):
    '''
    Retrieves all user scores from database
    '''
    if not is_admin:
        raise NOT_ADMIN_EXCEPTION
    db_scores = await utils.get_all_player_scores(database)
    return db_scores

@router.get('/player/{user_id}', response_model = DayScore)
async def get_player_score_by_id(
    user_id: int,
    database: Session = Depends(get_db),
    is_admin: bool = Depends(is_current_user_admin)
):
    '''
    Retrieves a user's score from database if that user exists
    '''
    if not is_admin:
        raise NOT_ADMIN_EXCEPTION

    user_score = await utils.get_player_score_by_id(user_id, database)
    if not user_score:
        raise USER_NOT_FOUND_EXCEPION
    return user_score

@router.put('/player/{user_id}', status_code = status.HTTP_202_ACCEPTED, response_model = DayScore)
async def update_user_score(
    user_id: int, 
    score: DayScore = Body(...),
    database: Session = Depends(get_db),
    is_admin: bool = Depends(is_current_user_admin)
):
    '''
    Updates a users score by user id
    '''
    if not is_admin:
        raise NOT_ADMIN_EXCEPTION

    user_score = await utils.update_player_score_by_id(user_id, score, database)
    if not user_score:
        raise COULD_NOT_UPDATE_EXCEPTION('player score')
        
    return user_score
    
@router.put('/update_multiple_by_id', status_code = status.HTTP_202_ACCEPTED, response_model = List[DayScore])
async def update_user_scores_by_id(
    scores: List[DayScoreWithUserId] = Body(...),
    database: Session = Depends(get_db),
    is_admin: bool = Depends(is_current_user_admin)
):
    '''
    Updates multiple users scores
    '''
    if not is_admin:
        raise NOT_ADMIN_EXCEPTION

    updated_scores = await utils.update_multiple_scores(scores, database)
    if not updated_scores:
        raise COULD_NOT_UPDATE_EXCEPTION('points table with new scores')

    return updated_scores

@router.get('/position', response_model = PlayerPosition)
async def get_current_user_position(
    database: Session = Depends(get_db),
    current_user: UserModel = Depends(get_current_user)
):
    '''
    Gets the current users position and score
    '''
    player_score = await utils.get_user_position(
        current_user.id,
        database
    )
    return player_score

@router.get('/leaderboard', response_model = List[PlayerPositionWithDisplayName])
async def get_leaderboard(
    table_start: int,
    table_end: int,
    database: Session = Depends(get_db)
):
    if table_start < 1 or table_start >= table_end:
        raise INVALID_QUERY_PARAMETERS_EXCEPTION('Invalid table start/end. Table start must be greater than 0 and less than table end.')

    leadboard = await utils.get_leaderboard(table_start, table_end, database)

    return leadboard