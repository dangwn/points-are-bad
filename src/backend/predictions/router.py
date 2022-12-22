from fastapi import APIRouter, Depends, status
from sqlalchemy.orm import Session
from typing import List

from db import get_db
from authentication.utils import get_current_user, is_current_user_admin
from predictions import utils
from user.models import User as UserModel
from http_exceptions import NOT_AUTHORIZED_EXCEPTION, NOT_FOUND_EXCEPTION, NOT_ADMIN_EXCEPTION, USER_NOT_FOUND_EXCEPION
from predictions.schema import MatchPrediction, MatchPredictionWithUserId, DisplayMatchPrediction

router = APIRouter(
    prefix = '/predictions',
    tags = ['predictions']
)

@router.get('/', response_model = List[DisplayMatchPrediction])
async def get_current_user_predictions(
    database: Session = Depends(get_db),
    current_user: UserModel = Depends(get_current_user)
) -> List[DisplayMatchPrediction]:
    '''
    Returns the predictions for a given user
    '''

    predictions = await utils.get_display_predictions(current_user.id, database)

    if predictions is None:
        raise NOT_FOUND_EXCEPTION(f'Predictions for user {current_user.id}')

    return predictions

@router.put('/', status_code = status.HTTP_202_ACCEPTED, response_model = List[DisplayMatchPrediction])
async def update_predictions(
    predictions: List[MatchPrediction],
    database: Session = Depends(get_db),
    current_user: UserModel = Depends(get_current_user)
) -> List[MatchPredictionWithUserId]:
    '''
    Updates predictions(s) for the current user
    '''
    if not current_user:
        raise USER_NOT_FOUND_EXCEPION

    updated_predictions = await utils.update_predictions(predictions, current_user.id, database)

    return updated_predictions

@router.get('/user/{user_id}', response_model = List[MatchPrediction])
async def get_user_predictions(
    user_id: int,
    database: Session = Depends(get_db),
    is_admin: UserModel = Depends(is_current_user_admin)
) -> List[MatchPrediction]:
    '''
    Returns the predictions for a given user
    '''
    if not is_admin:
        raise NOT_ADMIN_EXCEPTION

    predictions = await utils.get_predictions_by_user_id(user_id, database)

    if predictions is None:
        raise NOT_FOUND_EXCEPTION(f'Predictions for user {user_id}')

    return predictions

@router.get('/match/{match_id}', response_model = List[MatchPredictionWithUserId])
async def get_match_predictions(
    match_id: int,
    database: Session = Depends(get_db),
    is_admin: bool = Depends(is_current_user_admin)
) -> List[MatchPredictionWithUserId]:
    '''
    Returns a list of predictions for a given match
    '''
    if not is_admin:
        raise NOT_ADMIN_EXCEPTION

    predictions = await utils.get_predictions_by_match_id(match_id, database)
    
    if predictions is None:
        raise NOT_FOUND_EXCEPTION(f'Predictions for match {match_id}')

    return predictions

